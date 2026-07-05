package usecase

import (
	"context"
	"errors"

	"pos-project/backend/entity"
)

// LowStockOverview bundles the current threshold with the items that need attention.
type LowStockOverview struct {
	Settings entity.InventorySettings
	Items    []entity.MenuItem
}

// UpdateLowStockThresholdInput captures the manager input for threshold changes.
type UpdateLowStockThresholdInput struct {
	LowStockThreshold int
	Actor             string
	Reason            string
}

// InventoryUsecase coordinates low-stock warnings and threshold maintenance.
type InventoryUsecase struct {
	inventoryRepo InventoryRepository
}

// NewInventoryUsecase constructs an inventory usecase.
func NewInventoryUsecase(inventoryRepo InventoryRepository) *InventoryUsecase {
	return &InventoryUsecase{inventoryRepo: inventoryRepo}
}

// ListLowStock returns the current threshold and all items at or below it.
func (u *InventoryUsecase) ListLowStock(ctx context.Context) (LowStockOverview, error) {
	return u.GetLowStockItems(ctx)
}

// GetLowStockItems returns the current threshold and all items at or below it.
func (u *InventoryUsecase) GetLowStockItems(ctx context.Context) (LowStockOverview, error) {
	if u.inventoryRepo == nil {
		return LowStockOverview{}, errors.New("inventory repository is required")
	}

	settings, err := u.inventoryRepo.GetInventorySettings(ctx)
	if err != nil {
		return LowStockOverview{}, err
	}
	items, err := u.inventoryRepo.ListLowStockMenuItems(ctx, settings.LowStockThreshold)
	if err != nil {
		return LowStockOverview{}, err
	}

	return LowStockOverview{
		Settings: settings,
		Items:    items,
	}, nil
}

// UpdateLowStockThreshold persists a manager-controlled warning threshold.
func (u *InventoryUsecase) UpdateLowStockThreshold(ctx context.Context, input UpdateLowStockThresholdInput) (entity.InventorySettings, error) {
	if err := validateActor(input.Actor); err != nil {
		return entity.InventorySettings{}, err
	}
	if err := validateReasonIfProvided(input.Reason); err != nil {
		return entity.InventorySettings{}, err
	}
	if input.LowStockThreshold < 0 {
		return entity.InventorySettings{}, errors.New("low stock threshold must be zero or greater")
	}
	if u.inventoryRepo == nil {
		return entity.InventorySettings{}, errors.New("inventory repository is required")
	}

	settings := entity.InventorySettings{LowStockThreshold: input.LowStockThreshold}
	if err := settings.Validate(); err != nil {
		return entity.InventorySettings{}, err
	}

	return u.inventoryRepo.UpdateInventorySettings(ctx, settings)
}
