package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"pos-project/backend/entity"
)

func TestInventoryUsecaseListLowStock(t *testing.T) {
	t.Parallel()

	settings := entity.InventorySettings{
		ID:                1,
		LowStockThreshold: 4,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}
	items := []entity.MenuItem{
		{ID: 1, Name: "Tea", StockQuantity: 2, IsActive: true, IsAvailable: true},
		{ID: 2, Name: "Soup", StockQuantity: 4, IsActive: true, IsAvailable: true},
	}

	repo := &inventoryRepositoryStub{
		settings: settings,
		items:    items,
	}
	uc := NewInventoryUsecase(repo)

	overview, err := uc.ListLowStock(context.Background())
	if err != nil {
		t.Fatalf("ListLowStock returned error: %v", err)
	}
	if overview.Settings.LowStockThreshold != 4 {
		t.Fatalf("unexpected threshold: %#v", overview.Settings)
	}
	if len(overview.Items) != 2 {
		t.Fatalf("unexpected items: %#v", overview.Items)
	}
	if repo.lastThreshold != 4 {
		t.Fatalf("expected repository to receive threshold 4, got %d", repo.lastThreshold)
	}
}

func TestInventoryUsecaseGetLowStockItems(t *testing.T) {
	t.Parallel()

	settings := entity.InventorySettings{
		ID:                1,
		LowStockThreshold: 2,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}
	repo := &inventoryRepositoryStub{
		settings: settings,
		items: []entity.MenuItem{
			{ID: 1, Name: "Chili Oil", StockQuantity: 1, IsActive: true, IsAvailable: true},
		},
	}
	uc := NewInventoryUsecase(repo)

	overview, err := uc.GetLowStockItems(context.Background())
	if err != nil {
		t.Fatalf("GetLowStockItems returned error: %v", err)
	}
	if overview.Settings.LowStockThreshold != 2 {
		t.Fatalf("unexpected threshold: %#v", overview.Settings)
	}
	if len(overview.Items) != 1 || overview.Items[0].Name != "Chili Oil" {
		t.Fatalf("unexpected low-stock items: %#v", overview.Items)
	}
}

func TestInventoryUsecaseUpdateLowStockThreshold(t *testing.T) {
	t.Parallel()

	repo := &inventoryRepositoryStub{
		settings: entity.InventorySettings{
			ID:                1,
			LowStockThreshold: 5,
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
		},
	}
	uc := NewInventoryUsecase(repo)

	updated, err := uc.UpdateLowStockThreshold(context.Background(), UpdateLowStockThresholdInput{
		LowStockThreshold: 7,
		Actor:             "manager-1",
		Reason:            "seasonal demand",
	})
	if err != nil {
		t.Fatalf("UpdateLowStockThreshold returned error: %v", err)
	}
	if updated.LowStockThreshold != 7 {
		t.Fatalf("unexpected updated settings: %#v", updated)
	}
	if repo.settings.LowStockThreshold != 7 {
		t.Fatalf("expected repository settings to update, got %#v", repo.settings)
	}
}

func TestInventoryUsecaseValidationErrors(t *testing.T) {
	t.Parallel()

	uc := NewInventoryUsecase(nil)

	if _, err := uc.ListLowStock(context.Background()); err == nil {
		t.Fatal("expected missing repository error")
	}
	if _, err := uc.UpdateLowStockThreshold(context.Background(), UpdateLowStockThresholdInput{LowStockThreshold: -1, Actor: "manager-1"}); err == nil {
		t.Fatal("expected negative threshold validation error")
	}
	if _, err := uc.UpdateLowStockThreshold(context.Background(), UpdateLowStockThresholdInput{LowStockThreshold: 1, Actor: "manager-1"}); err == nil {
		t.Fatal("expected repository required error")
	}
}

func TestInventoryUsecaseRepositoryErrors(t *testing.T) {
	t.Parallel()

	repo := &inventoryRepositoryStub{
		settingsErr: errors.New("settings unavailable"),
		itemsErr:    errors.New("items unavailable"),
	}
	uc := NewInventoryUsecase(repo)

	if _, err := uc.ListLowStock(context.Background()); err == nil {
		t.Fatal("expected settings error")
	}

	repo.settingsErr = nil
	if _, err := uc.GetLowStockItems(context.Background()); err == nil {
		t.Fatal("expected items error")
	}

	if _, err := uc.UpdateLowStockThreshold(context.Background(), UpdateLowStockThresholdInput{
		LowStockThreshold: 3,
		Actor:             "",
	}); err == nil {
		t.Fatal("expected actor validation error")
	}
}

type inventoryRepositoryStub struct {
	settings      entity.InventorySettings
	items         []entity.MenuItem
	lastThreshold int
	settingsErr   error
	itemsErr      error
}

func (s *inventoryRepositoryStub) GetInventorySettings(context.Context) (entity.InventorySettings, error) {
	if s.settingsErr != nil {
		return entity.InventorySettings{}, s.settingsErr
	}
	if s.settings.ID == 0 {
		return entity.InventorySettings{}, errors.New("inventory settings not found")
	}
	return s.settings, nil
}

func (s *inventoryRepositoryStub) UpdateInventorySettings(_ context.Context, settings entity.InventorySettings) (entity.InventorySettings, error) {
	s.settings = entity.InventorySettings{
		ID:                1,
		LowStockThreshold: settings.LowStockThreshold,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}
	return s.settings, nil
}

func (s *inventoryRepositoryStub) UpdateMenuItemStock(context.Context, int64, int) (entity.MenuItem, error) {
	return entity.MenuItem{}, nil
}

func (s *inventoryRepositoryStub) ListLowStockMenuItems(_ context.Context, threshold int) ([]entity.MenuItem, error) {
	if s.itemsErr != nil {
		return nil, s.itemsErr
	}
	s.lastThreshold = threshold
	return s.items, nil
}
