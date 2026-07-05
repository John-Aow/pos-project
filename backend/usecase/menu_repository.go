package usecase

import (
	"context"

	"pos-project/backend/entity"
)

// MenuManagementRepository describes menu catalog persistence.
type MenuManagementRepository interface {
	ListCategories(ctx context.Context) ([]entity.Category, error)
	ListMenuItems(ctx context.Context) ([]entity.MenuItem, error)
	ListAvailableMenuItems(ctx context.Context) ([]entity.MenuItem, error)
	GetMenuItemByID(ctx context.Context, id int64) (entity.MenuItem, error)
	CreateMenuItem(ctx context.Context, item entity.MenuItem) (entity.MenuItem, error)
	UpdateMenuItem(ctx context.Context, item entity.MenuItem) (entity.MenuItem, error)
	UpdateMenuItemPrice(ctx context.Context, id int64, priceCents int64) (entity.MenuItem, error)
	UpdateMenuItemStock(ctx context.Context, id int64, stockQuantity int) (entity.MenuItem, error)
	DeactivateMenuItem(ctx context.Context, id int64) (entity.MenuItem, error)
}

// MenuRepository preserves the existing usecase dependency name.
type MenuRepository = MenuManagementRepository
