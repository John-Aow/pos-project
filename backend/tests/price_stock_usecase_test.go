package tests

import (
	"context"
	"errors"
	"testing"
	"time"

	"pos-project/backend/entity"
	"pos-project/backend/usecase"
)

func TestMenuUsecaseUpdateMenuItemPrice(t *testing.T) {
	t.Parallel()

	menuRepo := &fakePriceStockMenuRepository{
		menuItem: entity.MenuItem{
			ID:                41,
			Name:              "Fried Rice",
			Description:       "Classic fried rice",
			PriceCents:        10000,
			StockQuantity:     7,
			LowStockThreshold: 3,
			IsActive:          true,
			IsAvailable:       true,
			CategoryID:        2,
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
		},
	}
	auditRepo := &fakeAuditRepository{}
	priceHistoryRepo := &fakePriceHistoryRepository{}
	auditLogRepo := &fakeMenuAuditLogRepository{}
	uc := usecase.NewMenuUsecase(menuRepo, auditRepo, priceHistoryRepo, auditLogRepo)

	historicalBillItem := menuRepo.menuItem
	updated, err := uc.UpdateMenuItemPrice(context.Background(), usecase.UpdateMenuItemPriceInput{
		ID:         41,
		PriceCents: 12500,
		Actor:      "manager-1",
		Reason:     "menu update",
	})
	if err != nil {
		t.Fatalf("UpdateMenuItemPrice returned error: %v", err)
	}
	if updated.PriceCents != 12500 {
		t.Fatalf("expected updated price 12500, got %d", updated.PriceCents)
	}
	if historicalBillItem.PriceCents != 10000 {
		t.Fatalf("expected historical snapshot to keep old price, got %d", historicalBillItem.PriceCents)
	}
	if len(auditRepo.entries) != 1 {
		t.Fatalf("expected one audit entry, got %d", len(auditRepo.entries))
	}
	if auditRepo.entries[0].ActionType != entity.AuditActionPrice {
		t.Fatalf("expected price audit entry, got %s", auditRepo.entries[0].ActionType)
	}
	if len(priceHistoryRepo.entries) != 1 {
		t.Fatalf("expected one price history entry, got %d", len(priceHistoryRepo.entries))
	}
	if priceHistoryRepo.entries[0].PreviousPriceCents != 10000 || priceHistoryRepo.entries[0].NewPriceCents != 12500 {
		t.Fatalf("unexpected price history entry: %#v", priceHistoryRepo.entries[0])
	}
	if len(auditLogRepo.entries) != 1 {
		t.Fatalf("expected one menu audit log entry, got %d", len(auditLogRepo.entries))
	}
	if auditLogRepo.entries[0].ActionType != entity.AuditActionPrice {
		t.Fatalf("expected price menu audit log entry, got %s", auditLogRepo.entries[0].ActionType)
	}
}

func TestMenuUsecaseUpdateMenuItemStock(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		initialActive bool
		initialStock  int
		nextStock     int
		wantAvailable bool
	}{
		{
			name:          "zero stock disables availability",
			initialActive: true,
			initialStock:  4,
			nextStock:     0,
			wantAvailable: false,
		},
		{
			name:          "positive stock restores availability for active item",
			initialActive: true,
			initialStock:  0,
			nextStock:     9,
			wantAvailable: true,
		},
		{
			name:          "inactive item stays unavailable",
			initialActive: false,
			initialStock:  2,
			nextStock:     6,
			wantAvailable: false,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			menuRepo := &fakePriceStockMenuRepository{
				menuItem: entity.MenuItem{
					ID:                52,
					Name:              "Spring Roll",
					Description:       "Crispy starter",
					PriceCents:        5000,
					StockQuantity:     tc.initialStock,
					LowStockThreshold: 3,
					IsActive:          tc.initialActive,
					IsAvailable:       tc.initialActive && tc.initialStock > 0,
					CategoryID:        5,
					CreatedAt:         time.Now(),
					UpdatedAt:         time.Now(),
				},
			}
			auditRepo := &fakeAuditRepository{}
			auditLogRepo := &fakeMenuAuditLogRepository{}
			uc := usecase.NewMenuUsecase(menuRepo, auditRepo, nil, auditLogRepo)

			updated, err := uc.UpdateMenuItemStock(context.Background(), usecase.UpdateMenuItemStockInput{
				ID:            52,
				StockQuantity: tc.nextStock,
				Actor:         "manager-1",
				Reason:        "restock check",
			})
			if err != nil {
				t.Fatalf("UpdateMenuItemStock returned error: %v", err)
			}
			if updated.StockQuantity != tc.nextStock {
				t.Fatalf("expected stock %d, got %d", tc.nextStock, updated.StockQuantity)
			}
			if updated.IsAvailable != tc.wantAvailable {
				t.Fatalf("expected availability %v, got %v", tc.wantAvailable, updated.IsAvailable)
			}
			if len(auditRepo.entries) != 1 {
				t.Fatalf("expected one audit entry, got %d", len(auditRepo.entries))
			}
			if auditRepo.entries[0].ActionType != entity.AuditActionStock {
				t.Fatalf("expected stock audit entry, got %s", auditRepo.entries[0].ActionType)
			}
			if len(auditLogRepo.entries) != 1 {
				t.Fatalf("expected one menu audit log entry, got %d", len(auditLogRepo.entries))
			}
			if auditLogRepo.entries[0].ActionType != entity.AuditActionStock {
				t.Fatalf("expected stock menu audit log entry, got %s", auditLogRepo.entries[0].ActionType)
			}
		})
	}
}

func TestMenuUsecaseUpdateMenuItemPriceValidation(t *testing.T) {
	t.Parallel()

	uc := usecase.NewMenuUsecase(&fakePriceStockMenuRepository{}, &fakeAuditRepository{})

	_, err := uc.UpdateMenuItemPrice(context.Background(), usecase.UpdateMenuItemPriceInput{
		ID:         99,
		PriceCents: 0,
		Actor:      "manager-1",
	})
	if err == nil {
		t.Fatal("expected validation error for zero price")
	}
}

func TestMenuUsecaseUpdateMenuItemPriceRejectsNoChange(t *testing.T) {
	t.Parallel()

	uc := usecase.NewMenuUsecase(&fakePriceStockMenuRepository{
		menuItem: entity.MenuItem{
			ID:                99,
			Name:              "Soup",
			Description:       "Tom yum",
			PriceCents:        8800,
			StockQuantity:     2,
			LowStockThreshold: 3,
			IsActive:          true,
			IsAvailable:       true,
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
		},
	}, &fakeAuditRepository{}, &fakePriceHistoryRepository{})

	_, err := uc.UpdateMenuItemPrice(context.Background(), usecase.UpdateMenuItemPriceInput{
		ID:         99,
		PriceCents: 8800,
		Actor:      "manager-1",
	})
	if err == nil {
		t.Fatal("expected validation error when price does not change")
	}
}

type fakePriceStockMenuRepository struct {
	menuItem entity.MenuItem

	getCalls   int
	priceCalls int
	stockCalls int
	lastPrice  int64
	lastStock  int
}

func (f *fakePriceStockMenuRepository) ListCategories(context.Context) ([]entity.Category, error) {
	return nil, nil
}

func (f *fakePriceStockMenuRepository) ListMenuItems(context.Context) ([]entity.MenuItem, error) {
	return nil, nil
}

func (f *fakePriceStockMenuRepository) ListAvailableMenuItems(context.Context) ([]entity.MenuItem, error) {
	return nil, nil
}

func (f *fakePriceStockMenuRepository) GetMenuItemByID(context.Context, int64) (entity.MenuItem, error) {
	f.getCalls++
	if f.menuItem.ID == 0 {
		return entity.MenuItem{}, errors.New("not found")
	}
	return f.menuItem, nil
}

func (f *fakePriceStockMenuRepository) CreateMenuItem(context.Context, entity.MenuItem) (entity.MenuItem, error) {
	return entity.MenuItem{}, nil
}

func (f *fakePriceStockMenuRepository) UpdateMenuItem(context.Context, entity.MenuItem) (entity.MenuItem, error) {
	return entity.MenuItem{}, nil
}

func (f *fakePriceStockMenuRepository) UpdateMenuItemPrice(_ context.Context, id int64, priceCents int64) (entity.MenuItem, error) {
	f.priceCalls++
	f.lastPrice = priceCents
	f.menuItem.ID = id
	f.menuItem.PriceCents = priceCents
	return f.menuItem, nil
}

func (f *fakePriceStockMenuRepository) UpdateMenuItemStock(_ context.Context, id int64, stockQuantity int) (entity.MenuItem, error) {
	f.stockCalls++
	f.lastStock = stockQuantity
	f.menuItem.ID = id
	f.menuItem.StockQuantity = stockQuantity
	f.menuItem.IsAvailable = f.menuItem.IsActive && stockQuantity > 0
	return f.menuItem, nil
}

func (f *fakePriceStockMenuRepository) DeactivateMenuItem(context.Context, int64) (entity.MenuItem, error) {
	return entity.MenuItem{}, nil
}

type fakePriceHistoryRepository struct {
	entries []entity.PriceHistory
}

func (f *fakePriceHistoryRepository) RecordPriceHistory(_ context.Context, history entity.PriceHistory) error {
	f.entries = append(f.entries, history)
	return nil
}

type fakeMenuAuditLogRepository struct {
	entries []entity.MenuAuditLog
}

func (f *fakeMenuAuditLogRepository) RecordMenuAuditLog(_ context.Context, log entity.MenuAuditLog) error {
	f.entries = append(f.entries, log)
	return nil
}
