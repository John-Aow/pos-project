package postgres

import (
	"context"
	"database/sql"
	"time"

	"pos-project/backend/entity"
)

// InventoryRepository persists low-stock threshold settings and stock queries.
type InventoryRepository struct {
	db *sql.DB
}

// NewInventoryRepository constructs a Postgres-backed inventory repository.
func NewInventoryRepository(db *sql.DB) *InventoryRepository {
	return &InventoryRepository{db: db}
}

// GetInventorySettings returns the singleton low-stock threshold configuration.
func (r *InventoryRepository) GetInventorySettings(ctx context.Context) (entity.InventorySettings, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, low_stock_threshold, created_at, updated_at
		FROM inventory_settings
		WHERE id = 1
	`)

	var settings entity.InventorySettings
	if err := row.Scan(&settings.ID, &settings.LowStockThreshold, &settings.CreatedAt, &settings.UpdatedAt); err != nil {
		return entity.InventorySettings{}, err
	}

	return settings, nil
}

// UpdateInventorySettings updates the singleton low-stock threshold.
func (r *InventoryRepository) UpdateInventorySettings(ctx context.Context, settings entity.InventorySettings) (entity.InventorySettings, error) {
	now := time.Now().UTC()
	row := r.db.QueryRowContext(ctx, `
		UPDATE inventory_settings
		SET low_stock_threshold = $2,
		    updated_at = $3
		WHERE id = 1
		RETURNING id, low_stock_threshold, created_at, updated_at
	`, 1, settings.LowStockThreshold, now)

	var updated entity.InventorySettings
	if err := row.Scan(&updated.ID, &updated.LowStockThreshold, &updated.CreatedAt, &updated.UpdatedAt); err != nil {
		return entity.InventorySettings{}, err
	}

	return updated, nil
}

// UpdateMenuItemStock updates item stock and derived availability.
func (r *InventoryRepository) UpdateMenuItemStock(ctx context.Context, id int64, stockQuantity int) (entity.MenuItem, error) {
	row := r.db.QueryRowContext(ctx, `
		UPDATE menu_items
		SET stock_quantity = $2,
		    is_available = CASE WHEN is_active = TRUE AND $2 > 0 THEN TRUE ELSE FALSE END,
		    updated_at = NOW()
		WHERE id = $1
		RETURNING id, name, description, price_cents, stock_quantity, low_stock_threshold, is_available, is_active, category_id, created_at, updated_at
	`, id, stockQuantity)

	return scanMenuItem(row)
}

// ListLowStockMenuItems returns items at or below the provided threshold.
func (r *InventoryRepository) ListLowStockMenuItems(ctx context.Context, threshold int) ([]entity.MenuItem, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, name, description, price_cents, stock_quantity, low_stock_threshold, is_available, is_active, category_id, created_at, updated_at
		FROM menu_items
		WHERE is_active = TRUE AND stock_quantity <= $1
		ORDER BY stock_quantity ASC, name ASC
	`, threshold)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []entity.MenuItem
	for rows.Next() {
		item, err := scanMenuItemRows(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, rows.Err()
}
