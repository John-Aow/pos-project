package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"pos-project/backend/entity"
)

func TestMenuUsecaseListMethods(t *testing.T) {
	t.Parallel()

	item := entity.MenuItem{
		ID:            1,
		Name:          "Green Curry",
		Description:   "Thai curry",
		PriceCents:    12900,
		StockQuantity: 4,
		IsActive:      true,
		IsAvailable:   true,
		CategoryID:    2,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	repo := &menuRepositoryStub{
		menuItems:      []entity.MenuItem{item},
		availableItems: []entity.MenuItem{item},
	}

	uc := NewMenuUsecase(repo, &auditRepositoryStub{})

	customerItems, err := uc.ListCustomerMenu(context.Background())
	if err != nil {
		t.Fatalf("ListCustomerMenu returned error: %v", err)
	}
	if len(customerItems) != 1 || customerItems[0].ID != item.ID {
		t.Fatalf("unexpected customer items: %#v", customerItems)
	}

	managerItems, err := uc.ListManagerMenu(context.Background())
	if err != nil {
		t.Fatalf("ListManagerMenu returned error: %v", err)
	}
	if len(managerItems) != 1 || managerItems[0].ID != item.ID {
		t.Fatalf("unexpected manager items: %#v", managerItems)
	}
}

func TestMenuUsecaseCreateUpdateAndAudit(t *testing.T) {
	t.Parallel()

	now := time.Now()
	repo := &menuRepositoryStub{
		item: entity.MenuItem{
			ID:            10,
			Name:          "Pad Thai",
			Description:   "Rice noodles",
			PriceCents:    12000,
			StockQuantity: 5,
			IsAvailable:   true,
			IsActive:      true,
			CategoryID:    3,
			CreatedAt:     now,
			UpdatedAt:     now,
		},
	}
	auditRepo := &auditRepositoryStub{}
	uc := NewMenuUsecase(repo, auditRepo)

	created, err := uc.CreateMenuItem(context.Background(), CreateMenuItemInput{
		Name:          "Green Curry",
		Description:   "Thai curry",
		PriceCents:    13900,
		StockQuantity: 6,
		CategoryID:    2,
		Actor:         "manager-1",
		Reason:        "new menu",
	})
	if err != nil {
		t.Fatalf("CreateMenuItem returned error: %v", err)
	}
	if created.Name != "Green Curry" || !created.IsAvailable {
		t.Fatalf("unexpected created item: %#v", created)
	}
	if len(auditRepo.entries) != 1 || auditRepo.entries[0].ActionType != entity.AuditActionCreate {
		t.Fatalf("expected create audit entry, got %#v", auditRepo.entries)
	}

	repo.item = entity.MenuItem{
		ID:            10,
		Name:          "Pad Thai",
		Description:   "Rice noodles",
		PriceCents:    12000,
		StockQuantity: 5,
		IsAvailable:   true,
		IsActive:      true,
		CategoryID:    3,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	updated, err := uc.UpdateMenuItem(context.Background(), UpdateMenuItemInput{
		ID:            10,
		Name:          "Pad Thai Deluxe",
		Description:   "Rice noodles with prawns",
		PriceCents:    13500,
		StockQuantity: 1,
		CategoryID:    4,
		Actor:         "manager-1",
		Reason:        "refresh",
	})
	if err != nil {
		t.Fatalf("UpdateMenuItem returned error: %v", err)
	}
	if updated.Name != "Pad Thai Deluxe" || updated.PriceCents != 13500 || updated.StockQuantity != 1 {
		t.Fatalf("unexpected updated item: %#v", updated)
	}

	priceUpdated, err := uc.UpdateMenuItemPrice(context.Background(), UpdateMenuItemPriceInput{
		ID:         10,
		PriceCents: 15000,
		Actor:      "manager-1",
	})
	if err != nil {
		t.Fatalf("UpdateMenuItemPrice returned error: %v", err)
	}
	if priceUpdated.PriceCents != 15000 {
		t.Fatalf("expected updated price, got %#v", priceUpdated)
	}

	stockUpdated, err := uc.UpdateMenuItemStock(context.Background(), UpdateMenuItemStockInput{
		ID:            10,
		StockQuantity: 0,
		Actor:         "manager-1",
		Reason:        "sold out",
	})
	if err != nil {
		t.Fatalf("UpdateMenuItemStock returned error: %v", err)
	}
	if stockUpdated.StockQuantity != 0 || stockUpdated.IsAvailable {
		t.Fatalf("expected zero-stock item to be unavailable, got %#v", stockUpdated)
	}

	deactivated, err := uc.DeactivateMenuItem(context.Background(), 10, "manager-1", "remove")
	if err != nil {
		t.Fatalf("DeactivateMenuItem returned error: %v", err)
	}
	if deactivated.IsActive {
		t.Fatalf("expected item to be deactivated, got %#v", deactivated)
	}
}

func TestMenuUsecaseValidationErrors(t *testing.T) {
	t.Parallel()

	uc := NewMenuUsecase(nil, nil)

	if _, err := uc.ListCustomerMenu(context.Background()); err == nil {
		t.Fatal("expected list customer menu error")
	}
	if _, err := uc.ListManagerMenu(context.Background()); err == nil {
		t.Fatal("expected list manager menu error")
	}
	if _, err := uc.CreateMenuItem(context.Background(), CreateMenuItemInput{}); err == nil {
		t.Fatal("expected create validation error")
	}
	if _, err := uc.UpdateMenuItem(context.Background(), UpdateMenuItemInput{}); err == nil {
		t.Fatal("expected update validation error")
	}
	if _, err := uc.UpdateMenuItemPrice(context.Background(), UpdateMenuItemPriceInput{}); err == nil {
		t.Fatal("expected price validation error")
	}
	if _, err := uc.UpdateMenuItemStock(context.Background(), UpdateMenuItemStockInput{}); err == nil {
		t.Fatal("expected stock validation error")
	}
	if _, err := uc.DeactivateMenuItem(context.Background(), 0, "", ""); err == nil {
		t.Fatal("expected deactivate validation error")
	}
	if _, err := NewMenuUsecase(nil, nil).CreateMenuItem(context.Background(), CreateMenuItemInput{
		Name:          "Soup",
		Description:   "Hot",
		PriceCents:    1000,
		StockQuantity: 1,
		CategoryID:    1,
		Actor:         "manager-1",
	}); err == nil {
		t.Fatal("expected missing repository error after validation")
	}
	if _, err := NewMenuUsecase(nil, nil).UpdateMenuItem(context.Background(), UpdateMenuItemInput{
		ID:            1,
		Name:          "Soup",
		Description:   "Hot",
		PriceCents:    1000,
		StockQuantity: 1,
		CategoryID:    1,
		Actor:         "manager-1",
	}); err == nil {
		t.Fatal("expected missing repository error on update")
	}
	if _, err := NewMenuUsecase(nil, nil).UpdateMenuItemPrice(context.Background(), UpdateMenuItemPriceInput{
		ID:         1,
		PriceCents: 1000,
		Actor:      "manager-1",
	}); err == nil {
		t.Fatal("expected missing repository error on price update")
	}
	if _, err := NewMenuUsecase(nil, nil).UpdateMenuItemStock(context.Background(), UpdateMenuItemStockInput{
		ID:            1,
		StockQuantity: 1,
		Actor:         "manager-1",
	}); err == nil {
		t.Fatal("expected missing repository error on stock update")
	}
	if _, err := NewMenuUsecase(nil, nil).DeactivateMenuItem(context.Background(), 1, "manager-1", "remove"); err == nil {
		t.Fatal("expected missing repository error on deactivate")
	}
	if err := validateReasonIfProvided(""); err != nil {
		t.Fatalf("expected optional reason to pass, got %v", err)
	}
}

func TestMenuUsecaseAuditDisabled(t *testing.T) {
	t.Parallel()

	repo := &menuRepositoryStub{
		item: entity.MenuItem{
			ID:            20,
			Name:          "Soup",
			Description:   "Tom yum",
			PriceCents:    8900,
			StockQuantity: 2,
			IsActive:      true,
			IsAvailable:   true,
			CategoryID:    1,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		},
	}

	uc := NewMenuUsecase(repo, nil)
	if _, err := uc.UpdateMenuItemPrice(context.Background(), UpdateMenuItemPriceInput{
		ID:         20,
		PriceCents: 9000,
		Actor:      "manager-1",
	}); err != nil {
		t.Fatalf("expected price update to succeed without audit repo, got %v", err)
	}
}

func TestMenuUsecaseOptionalReposAndHelperBranches(t *testing.T) {
	t.Parallel()

	now := time.Now()
	repo := &menuRepositoryStub{
		item: entity.MenuItem{
			ID:            40,
			Name:          "Laksa",
			Description:   "Noodle soup",
			PriceCents:    9900,
			StockQuantity: 8,
			IsActive:      true,
			IsAvailable:   true,
			CategoryID:    6,
			CreatedAt:     now,
			UpdatedAt:     now,
		},
	}
	auditRepo := &auditRepositoryStub{}
	priceHistoryRepo := &priceHistoryRepositoryStub{}
	auditLogRepo := &menuAuditLogRepositoryStub{}

	uc := NewMenuUsecase(repo, auditRepo, priceHistoryRepo, auditLogRepo)
	if uc.priceHistoryRepo == nil || uc.auditLogRepo == nil {
		t.Fatal("expected optional repositories to be wired")
	}

	if _, err := uc.UpdateMenuItemPrice(context.Background(), UpdateMenuItemPriceInput{
		ID:         40,
		PriceCents: 10900,
		Actor:      "manager-1",
		Reason:     "promo update",
	}); err != nil {
		t.Fatalf("UpdateMenuItemPrice returned error: %v", err)
	}
	if len(priceHistoryRepo.entries) != 1 {
		t.Fatalf("expected one price history entry, got %d", len(priceHistoryRepo.entries))
	}
	if len(auditLogRepo.entries) != 1 {
		t.Fatalf("expected one menu audit log entry, got %d", len(auditLogRepo.entries))
	}

	if _, err := uc.DeactivateMenuItem(context.Background(), 40, "manager-1", ""); err == nil {
		t.Fatal("expected deactivate validation error for empty reason")
	}

	if err := uc.recordPriceHistory(context.Background(), 40, 10900, 11900, "manager-1", "test"); err != nil {
		t.Fatalf("recordPriceHistory returned error: %v", err)
	}
	if err := uc.recordMenuAuditLog(context.Background(), 40, entity.AuditActionStock, repo.item, repo.item, "manager-1", "test"); err != nil {
		t.Fatalf("recordMenuAuditLog returned error: %v", err)
	}
	if err := uc.recordPriceHistory(context.Background(), 40, 11900, 11900, "manager-1", "test"); err == nil {
		t.Fatal("expected helper validation error for unchanged price")
	}
	if err := validateReason(""); err == nil {
		t.Fatal("expected empty reason validation error")
	}
	if err := validateReason("needed"); err != nil {
		t.Fatalf("expected valid reason, got %v", err)
	}
}

type menuRepositoryStub struct {
	menuItems      []entity.MenuItem
	availableItems []entity.MenuItem
	item           entity.MenuItem
}

func (s *menuRepositoryStub) ListCategories(context.Context) ([]entity.Category, error) {
	return nil, nil
}

func (s *menuRepositoryStub) ListMenuItems(context.Context) ([]entity.MenuItem, error) {
	return s.menuItems, nil
}

func (s *menuRepositoryStub) ListAvailableMenuItems(context.Context) ([]entity.MenuItem, error) {
	return s.availableItems, nil
}

func (s *menuRepositoryStub) GetMenuItemByID(context.Context, int64) (entity.MenuItem, error) {
	if s.item.ID == 0 {
		return entity.MenuItem{}, errors.New("menu item not found")
	}
	return s.item, nil
}

func (s *menuRepositoryStub) CreateMenuItem(context.Context, entity.MenuItem) (entity.MenuItem, error) {
	return entity.MenuItem{
		ID:            11,
		Name:          "Green Curry",
		Description:   "Thai curry",
		PriceCents:    13900,
		StockQuantity: 6,
		IsAvailable:   true,
		IsActive:      true,
		CategoryID:    2,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}, nil
}

func (s *menuRepositoryStub) UpdateMenuItem(context.Context, entity.MenuItem) (entity.MenuItem, error) {
	item := s.item
	item.Name = "Pad Thai Deluxe"
	item.Description = "Rice noodles with prawns"
	item.PriceCents = 13500
	item.StockQuantity = 1
	item.CategoryID = 4
	item.IsAvailable = true
	item.UpdatedAt = time.Now()
	return item, nil
}

func (s *menuRepositoryStub) UpdateMenuItemPrice(context.Context, int64, int64) (entity.MenuItem, error) {
	item := s.item
	item.PriceCents = 15000
	item.UpdatedAt = time.Now()
	return item, nil
}

func (s *menuRepositoryStub) UpdateMenuItemStock(context.Context, int64, int) (entity.MenuItem, error) {
	item := s.item
	item.StockQuantity = 0
	item.IsAvailable = false
	item.UpdatedAt = time.Now()
	return item, nil
}

func (s *menuRepositoryStub) DeactivateMenuItem(context.Context, int64) (entity.MenuItem, error) {
	item := s.item
	item.IsActive = false
	item.IsAvailable = false
	item.UpdatedAt = time.Now()
	return item, nil
}

type auditRepositoryStub struct {
	entries []entity.MenuItemAuditEntry
}

func (s *auditRepositoryStub) RecordMenuItemAuditEntry(_ context.Context, entry entity.MenuItemAuditEntry) error {
	s.entries = append(s.entries, entry)
	return nil
}

type priceHistoryRepositoryStub struct {
	entries []entity.PriceHistory
}

func (s *priceHistoryRepositoryStub) RecordPriceHistory(_ context.Context, history entity.PriceHistory) error {
	s.entries = append(s.entries, history)
	return nil
}

type menuAuditLogRepositoryStub struct {
	entries []entity.MenuAuditLog
}

func (s *menuAuditLogRepositoryStub) RecordMenuAuditLog(_ context.Context, log entity.MenuAuditLog) error {
	s.entries = append(s.entries, log)
	return nil
}
