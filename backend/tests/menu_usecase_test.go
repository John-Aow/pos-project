package tests

import (
	"context"
	"errors"
	"testing"
	"time"

	"pos-project/backend/entity"
	"pos-project/backend/usecase"
)

type fakeMenuRepository struct {
	menuItem entity.MenuItem
	items    []entity.MenuItem
	created  entity.MenuItem
	updated  entity.MenuItem

	getByIDCalls       int
	listAvailableCalls int
	listMenuCalls      int
	createCalls        int
	updateCalls        int
	deactivateCalls    int
	lastCreateInput    entity.MenuItem
	lastUpdateInput    entity.MenuItem
	lastDeactivateID   int64
	returnNotFound     bool
}

func (f *fakeMenuRepository) ListCategories(context.Context) ([]entity.Category, error) {
	return nil, nil
}

func (f *fakeMenuRepository) ListMenuItems(context.Context) ([]entity.MenuItem, error) {
	f.listMenuCalls++
	return append([]entity.MenuItem(nil), f.items...), nil
}

func (f *fakeMenuRepository) ListAvailableMenuItems(context.Context) ([]entity.MenuItem, error) {
	f.listAvailableCalls++
	return append([]entity.MenuItem(nil), f.items...), nil
}

func (f *fakeMenuRepository) GetMenuItemByID(context.Context, int64) (entity.MenuItem, error) {
	f.getByIDCalls++
	if f.returnNotFound {
		return entity.MenuItem{}, errors.New("not found")
	}
	return f.menuItem, nil
}

func (f *fakeMenuRepository) CreateMenuItem(_ context.Context, item entity.MenuItem) (entity.MenuItem, error) {
	f.createCalls++
	f.lastCreateInput = item
	if f.created.ID == 0 {
		f.created = item
		f.created.ID = 101
	}
	return f.created, nil
}

func (f *fakeMenuRepository) UpdateMenuItem(_ context.Context, item entity.MenuItem) (entity.MenuItem, error) {
	f.updateCalls++
	f.lastUpdateInput = item
	if f.updated.ID == 0 {
		f.updated = item
	}
	return f.updated, nil
}

func (f *fakeMenuRepository) UpdateMenuItemPrice(context.Context, int64, int64) (entity.MenuItem, error) {
	return entity.MenuItem{}, nil
}

func (f *fakeMenuRepository) UpdateMenuItemStock(context.Context, int64, int) (entity.MenuItem, error) {
	return entity.MenuItem{}, nil
}

func (f *fakeMenuRepository) DeactivateMenuItem(_ context.Context, id int64) (entity.MenuItem, error) {
	f.deactivateCalls++
	f.lastDeactivateID = id
	item := f.menuItem
	item.IsActive = false
	item.RecalculateAvailability()
	return item, nil
}

type fakeAuditRepository struct {
	entries []entity.MenuItemAuditEntry
}

func (f *fakeAuditRepository) RecordMenuItemAuditEntry(_ context.Context, entry entity.MenuItemAuditEntry) error {
	f.entries = append(f.entries, entry)
	return nil
}

func TestMenuUsecaseCreateMenuItem(t *testing.T) {
	t.Parallel()

	menuRepo := &fakeMenuRepository{}
	auditRepo := &fakeAuditRepository{}
	uc := usecase.NewMenuUsecase(menuRepo, auditRepo)

	created, err := uc.CreateMenuItem(context.Background(), usecase.CreateMenuItemInput{
		Name:          "Green Curry",
		Description:   "Thai curry with rice",
		PriceCents:    12900,
		StockQuantity: 12,
		CategoryID:    7,
		Actor:         "manager-1",
		Reason:        "new menu item",
	})
	if err != nil {
		t.Fatalf("CreateMenuItem returned error: %v", err)
	}
	if created.ID != 101 {
		t.Fatalf("expected created id 101, got %d", created.ID)
	}
	if !created.CustomerOrderable() {
		t.Fatalf("expected created item to be customer orderable")
	}
	if menuRepo.createCalls != 1 {
		t.Fatalf("expected one create call, got %d", menuRepo.createCalls)
	}
	if len(auditRepo.entries) != 1 {
		t.Fatalf("expected one audit entry, got %d", len(auditRepo.entries))
	}
	if auditRepo.entries[0].ActionType != entity.AuditActionCreate {
		t.Fatalf("expected create audit entry, got %s", auditRepo.entries[0].ActionType)
	}
}

func TestMenuUsecaseCreateMenuItemValidation(t *testing.T) {
	t.Parallel()

	uc := usecase.NewMenuUsecase(&fakeMenuRepository{}, &fakeAuditRepository{})

	_, err := uc.CreateMenuItem(context.Background(), usecase.CreateMenuItemInput{
		Description:   "missing name",
		PriceCents:    900,
		StockQuantity: 1,
		CategoryID:    1,
		Actor:         "manager-1",
	})
	if err == nil {
		t.Fatal("expected validation error")
	}
}

func TestMenuUsecaseUpdateMenuItem(t *testing.T) {
	t.Parallel()

	menuRepo := &fakeMenuRepository{
		menuItem: entity.MenuItem{
			ID:            22,
			Name:          "Pad Thai",
			Description:   "Noodles",
			PriceCents:    12900,
			StockQuantity: 5,
			IsActive:      true,
			CategoryID:    3,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		},
		updated: entity.MenuItem{
			ID:            22,
			Name:          "Pad Thai Special",
			Description:   "Noodles",
			PriceCents:    13900,
			StockQuantity: 8,
			IsActive:      true,
			CategoryID:    3,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		},
	}
	auditRepo := &fakeAuditRepository{}
	uc := usecase.NewMenuUsecase(menuRepo, auditRepo)

	updated, err := uc.UpdateMenuItem(context.Background(), usecase.UpdateMenuItemInput{
		ID:            22,
		Name:          "Pad Thai Special",
		PriceCents:    13900,
		StockQuantity: 8,
		CategoryID:    3,
		IsActive:      true,
		Actor:         "manager-1",
		Reason:        "menu refresh",
	})
	if err != nil {
		t.Fatalf("UpdateMenuItem returned error: %v", err)
	}
	if updated.Name != "Pad Thai Special" {
		t.Fatalf("expected updated name, got %s", updated.Name)
	}
	if menuRepo.getByIDCalls != 1 || menuRepo.updateCalls != 1 {
		t.Fatalf("expected one get and one update call, got get=%d update=%d", menuRepo.getByIDCalls, menuRepo.updateCalls)
	}
	if len(auditRepo.entries) != 1 || auditRepo.entries[0].ActionType != entity.AuditActionUpdate {
		t.Fatalf("expected update audit entry, got %+v", auditRepo.entries)
	}
}

func TestMenuUsecaseDeactivateMenuItem(t *testing.T) {
	t.Parallel()

	menuRepo := &fakeMenuRepository{
		menuItem: entity.MenuItem{
			ID:            31,
			Name:          "Iced Tea",
			Description:   "Cold tea",
			PriceCents:    4500,
			StockQuantity: 9,
			IsActive:      true,
			IsAvailable:   true,
			CategoryID:    4,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		},
	}
	auditRepo := &fakeAuditRepository{}
	uc := usecase.NewMenuUsecase(menuRepo, auditRepo)

	deactivated, err := uc.DeactivateMenuItem(context.Background(), 31, "manager-1", "seasonal change")
	if err != nil {
		t.Fatalf("DeactivateMenuItem returned error: %v", err)
	}
	if deactivated.IsActive {
		t.Fatalf("expected item to be inactive")
	}
	if menuRepo.deactivateCalls != 1 || menuRepo.lastDeactivateID != 31 {
		t.Fatalf("expected deactivate call for id 31, got calls=%d id=%d", menuRepo.deactivateCalls, menuRepo.lastDeactivateID)
	}
	if len(auditRepo.entries) != 1 || auditRepo.entries[0].ActionType != entity.AuditActionDeactivate {
		t.Fatalf("expected deactivate audit entry, got %+v", auditRepo.entries)
	}
}
