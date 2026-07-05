package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"pos-project/backend/entity"
)

// MenuUsecase coordinates catalog maintenance workflows.
type MenuUsecase struct {
	menuRepo         MenuRepository
	auditRepo        AuditRepository
	priceHistoryRepo PriceHistoryRepository
	auditLogRepo     AuditLogRepository
}

// NewMenuUsecase constructs a menu usecase.
func NewMenuUsecase(menuRepo MenuRepository, auditRepo AuditRepository, extras ...any) *MenuUsecase {
	uc := &MenuUsecase{
		menuRepo:  menuRepo,
		auditRepo: auditRepo,
	}
	for _, extra := range extras {
		switch repo := extra.(type) {
		case PriceHistoryRepository:
			uc.priceHistoryRepo = repo
		case AuditLogRepository:
			uc.auditLogRepo = repo
		}
	}
	return uc
}

// CreateMenuItemInput captures the minimum fields required to create a menu item.
type CreateMenuItemInput struct {
	Name              string
	Description       string
	PriceCents        int64
	StockQuantity     int
	LowStockThreshold int
	CategoryID        int64
	Actor             string
	Reason            string
}

// UpdateMenuItemInput captures editable menu item fields.
type UpdateMenuItemInput struct {
	ID                int64
	Name              string
	Description       string
	PriceCents        int64
	StockQuantity     int
	LowStockThreshold int
	CategoryID        int64
	IsActive          bool
	Actor             string
	Reason            string
}

// UpdateMenuItemPriceInput captures the fields for a price-only update.
type UpdateMenuItemPriceInput struct {
	ID         int64
	PriceCents int64
	Actor      string
	Reason     string
}

// UpdateMenuItemStockInput captures the fields for a stock-only update.
type UpdateMenuItemStockInput struct {
	ID            int64
	StockQuantity int
	Actor         string
	Reason        string
}

// ListCustomerMenu returns only items available for customer ordering.
func (u *MenuUsecase) ListCustomerMenu(ctx context.Context) ([]entity.MenuItem, error) {
	if u.menuRepo == nil {
		return nil, errors.New("menu repository is required")
	}

	return u.menuRepo.ListAvailableMenuItems(ctx)
}

// ListManagerMenu returns the full managed catalog.
func (u *MenuUsecase) ListManagerMenu(ctx context.Context) ([]entity.MenuItem, error) {
	if u.menuRepo == nil {
		return nil, errors.New("menu repository is required")
	}

	return u.menuRepo.ListMenuItems(ctx)
}

// CreateMenuItem validates and persists a new menu item.
func (u *MenuUsecase) CreateMenuItem(ctx context.Context, input CreateMenuItemInput) (entity.MenuItem, error) {
	if err := validateActor(input.Actor); err != nil {
		return entity.MenuItem{}, err
	}
	if err := validateReasonIfProvided(input.Reason); err != nil {
		return entity.MenuItem{}, err
	}

	item := entity.MenuItem{
		Name:              input.Name,
		Description:       input.Description,
		PriceCents:        input.PriceCents,
		StockQuantity:     input.StockQuantity,
		LowStockThreshold: input.LowStockThreshold,
		IsActive:          true,
		CategoryID:        input.CategoryID,
	}
	item.RecalculateAvailability()

	if err := item.Validate(); err != nil {
		return entity.MenuItem{}, err
	}
	if u.menuRepo == nil {
		return entity.MenuItem{}, errors.New("menu repository is required")
	}

	created, err := u.menuRepo.CreateMenuItem(ctx, item)
	if err != nil {
		return entity.MenuItem{}, err
	}

	if err := u.recordAudit(ctx, created.ID, entity.AuditActionCreate, nil, created, input.Actor, input.Reason); err != nil {
		return entity.MenuItem{}, err
	}
	if err := u.recordMenuAuditLog(ctx, created.ID, entity.AuditActionCreate, nil, created, input.Actor, input.Reason); err != nil {
		return entity.MenuItem{}, err
	}

	return created, nil
}

// UpdateMenuItem updates the editable catalog fields for an existing item.
func (u *MenuUsecase) UpdateMenuItem(ctx context.Context, input UpdateMenuItemInput) (entity.MenuItem, error) {
	if input.ID <= 0 {
		return entity.MenuItem{}, errors.New("menu item id is required")
	}
	if err := validateActor(input.Actor); err != nil {
		return entity.MenuItem{}, err
	}
	if err := validateReasonIfProvided(input.Reason); err != nil {
		return entity.MenuItem{}, err
	}
	if u.menuRepo == nil {
		return entity.MenuItem{}, errors.New("menu repository is required")
	}

	current, err := u.menuRepo.GetMenuItemByID(ctx, input.ID)
	if err != nil {
		return entity.MenuItem{}, err
	}
	previous := current

	if input.Name != "" {
		current.Name = input.Name
	}
	if input.Description != "" {
		current.Description = input.Description
	}
	if input.PriceCents > 0 {
		current.PriceCents = input.PriceCents
	}
	if input.StockQuantity >= 0 {
		current.StockQuantity = input.StockQuantity
	}
	if input.LowStockThreshold >= 0 {
		current.LowStockThreshold = input.LowStockThreshold
	}
	if input.CategoryID > 0 {
		current.CategoryID = input.CategoryID
	}
	current.RecalculateAvailability()

	if err := current.Validate(); err != nil {
		return entity.MenuItem{}, err
	}

	updated, err := u.menuRepo.UpdateMenuItem(ctx, current)
	if err != nil {
		return entity.MenuItem{}, err
	}

	if err := u.recordAudit(ctx, updated.ID, entity.AuditActionUpdate, previous, updated, input.Actor, input.Reason); err != nil {
		return entity.MenuItem{}, err
	}
	if err := u.recordMenuAuditLog(ctx, updated.ID, entity.AuditActionUpdate, previous, updated, input.Actor, input.Reason); err != nil {
		return entity.MenuItem{}, err
	}

	return updated, nil
}

// DeactivateMenuItem marks an item inactive without deleting history.
func (u *MenuUsecase) DeactivateMenuItem(ctx context.Context, id int64, actor, reason string) (entity.MenuItem, error) {
	if id <= 0 {
		return entity.MenuItem{}, errors.New("menu item id is required")
	}
	if err := validateActor(actor); err != nil {
		return entity.MenuItem{}, err
	}
	if err := validateReason(reason); err != nil {
		return entity.MenuItem{}, err
	}
	if u.menuRepo == nil {
		return entity.MenuItem{}, errors.New("menu repository is required")
	}

	current, err := u.menuRepo.GetMenuItemByID(ctx, id)
	if err != nil {
		return entity.MenuItem{}, err
	}

	current.IsActive = false
	current.RecalculateAvailability()

	updated, err := u.menuRepo.DeactivateMenuItem(ctx, id)
	if err != nil {
		return entity.MenuItem{}, err
	}

	if updated.ID == 0 {
		updated = current
	}

	if err := u.recordAudit(ctx, updated.ID, entity.AuditActionDeactivate, current, updated, actor, reason); err != nil {
		return entity.MenuItem{}, err
	}
	if err := u.recordMenuAuditLog(ctx, updated.ID, entity.AuditActionDeactivate, current, updated, actor, reason); err != nil {
		return entity.MenuItem{}, err
	}

	return updated, nil
}

// UpdateMenuItemPrice changes the live price for future orders.
func (u *MenuUsecase) UpdateMenuItemPrice(ctx context.Context, input UpdateMenuItemPriceInput) (entity.MenuItem, error) {
	if input.ID <= 0 {
		return entity.MenuItem{}, errors.New("menu item id is required")
	}
	if err := validateActor(input.Actor); err != nil {
		return entity.MenuItem{}, err
	}
	if err := validateReasonIfProvided(input.Reason); err != nil {
		return entity.MenuItem{}, err
	}
	if input.PriceCents <= 0 {
		return entity.MenuItem{}, errors.New("menu item price must be greater than zero")
	}
	if u.menuRepo == nil {
		return entity.MenuItem{}, errors.New("menu repository is required")
	}

	current, err := u.menuRepo.GetMenuItemByID(ctx, input.ID)
	if err != nil {
		return entity.MenuItem{}, err
	}
	if current.PriceCents == input.PriceCents {
		return entity.MenuItem{}, errors.New("menu item price must change")
	}

	updated, err := u.menuRepo.UpdateMenuItemPrice(ctx, input.ID, input.PriceCents)
	if err != nil {
		return entity.MenuItem{}, err
	}

	if err := u.recordPriceHistory(ctx, updated.ID, current.PriceCents, updated.PriceCents, input.Actor, input.Reason); err != nil {
		return entity.MenuItem{}, err
	}
	if err := u.recordAudit(ctx, updated.ID, entity.AuditActionPrice, current, updated, input.Actor, input.Reason); err != nil {
		return entity.MenuItem{}, err
	}
	if err := u.recordMenuAuditLog(ctx, updated.ID, entity.AuditActionPrice, current, updated, input.Actor, input.Reason); err != nil {
		return entity.MenuItem{}, err
	}

	return updated, nil
}

// UpdateMenuItemStock changes stock and derived availability.
func (u *MenuUsecase) UpdateMenuItemStock(ctx context.Context, input UpdateMenuItemStockInput) (entity.MenuItem, error) {
	if input.ID <= 0 {
		return entity.MenuItem{}, errors.New("menu item id is required")
	}
	if err := validateActor(input.Actor); err != nil {
		return entity.MenuItem{}, err
	}
	if err := validateReasonIfProvided(input.Reason); err != nil {
		return entity.MenuItem{}, err
	}
	if input.StockQuantity < 0 {
		return entity.MenuItem{}, errors.New("menu item stock quantity must be zero or greater")
	}
	if u.menuRepo == nil {
		return entity.MenuItem{}, errors.New("menu repository is required")
	}

	current, err := u.menuRepo.GetMenuItemByID(ctx, input.ID)
	if err != nil {
		return entity.MenuItem{}, err
	}

	updated, err := u.menuRepo.UpdateMenuItemStock(ctx, input.ID, input.StockQuantity)
	if err != nil {
		return entity.MenuItem{}, err
	}

	if err := u.recordAudit(ctx, updated.ID, entity.AuditActionStock, current, updated, input.Actor, input.Reason); err != nil {
		return entity.MenuItem{}, err
	}
	if err := u.recordMenuAuditLog(ctx, updated.ID, entity.AuditActionStock, current, updated, input.Actor, input.Reason); err != nil {
		return entity.MenuItem{}, err
	}

	return updated, nil
}

func (u *MenuUsecase) recordPriceHistory(
	ctx context.Context,
	menuItemID int64,
	previousPriceCents int64,
	newPriceCents int64,
	actor string,
	reason string,
) error {
	if u.priceHistoryRepo == nil {
		return nil
	}

	history := entity.PriceHistory{
		MenuItemID:         menuItemID,
		PreviousPriceCents: previousPriceCents,
		NewPriceCents:      newPriceCents,
		Actor:              actor,
		Reason:             reason,
		CreatedAt:          time.Now().UTC(),
	}
	if err := history.Validate(); err != nil {
		return err
	}

	return u.priceHistoryRepo.RecordPriceHistory(ctx, history)
}

func (u *MenuUsecase) recordMenuAuditLog(
	ctx context.Context,
	menuItemID int64,
	actionType string,
	previous any,
	next any,
	actor string,
	reason string,
) error {
	if u.auditLogRepo == nil {
		return nil
	}

	log := entity.MenuAuditLog{
		MenuItemID:    menuItemID,
		ActionType:    actionType,
		PreviousValue: marshalAuditValue(previous),
		NewValue:      marshalAuditValue(next),
		Actor:         actor,
		Reason:        reason,
		CreatedAt:     time.Now().UTC(),
	}
	if err := log.Validate(); err != nil {
		return err
	}

	return u.auditLogRepo.RecordMenuAuditLog(ctx, log)
}

func (u *MenuUsecase) recordAudit(
	ctx context.Context,
	menuItemID int64,
	actionType string,
	previous any,
	next any,
	actor string,
	reason string,
) error {
	if u.auditRepo == nil {
		return nil
	}

	entry := entity.MenuItemAuditEntry{
		MenuItemID:    menuItemID,
		ActionType:    actionType,
		PreviousValue: marshalAuditValue(previous),
		NewValue:      marshalAuditValue(next),
		Actor:         actor,
		Reason:        reason,
	}
	if entry.Reason == "" && actionType != entity.AuditActionCreate {
		entry.Reason = reason
	}

	return u.auditRepo.RecordMenuItemAuditEntry(ctx, entry)
}

func marshalAuditValue(value any) string {
	if value == nil {
		return ""
	}

	payload, err := json.Marshal(value)
	if err != nil {
		return fmt.Sprintf("%v", value)
	}

	return string(payload)
}

func validateActor(actor string) error {
	if actor == "" {
		return errors.New("actor is required")
	}
	return nil
}

func validateReason(reason string) error {
	if reason == "" {
		return errors.New("reason is required")
	}
	return nil
}

func validateReasonIfProvided(reason string) error {
	return nil
}
