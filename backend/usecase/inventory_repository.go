package usecase

import (
	"context"

	"pos-project/backend/entity"
)

// InventoryRepository describes stock and threshold persistence.
type InventoryRepository interface {
	GetInventorySettings(ctx context.Context) (entity.InventorySettings, error)
	UpdateInventorySettings(ctx context.Context, settings entity.InventorySettings) (entity.InventorySettings, error)
	UpdateMenuItemStock(ctx context.Context, id int64, stockQuantity int) (entity.MenuItem, error)
	ListLowStockMenuItems(ctx context.Context, threshold int) ([]entity.MenuItem, error)
}
