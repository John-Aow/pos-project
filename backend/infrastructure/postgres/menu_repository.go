package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"pos-project/backend/entity"
)

// MenuRepository persists menu catalog records in Postgres.
type MenuRepository struct {
	db *sql.DB
}

// NewMenuRepository constructs a Postgres-backed menu repository.
func NewMenuRepository(db *sql.DB) *MenuRepository {
	return &MenuRepository{db: db}
}

// ListCategories returns categories ordered by sort order and name.
func (r *MenuRepository) ListCategories(ctx context.Context) ([]entity.Category, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, name, sort_order, created_at, updated_at
		FROM categories
		ORDER BY sort_order ASC, name ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []entity.Category
	for rows.Next() {
		var category entity.Category
		if err := rows.Scan(&category.ID, &category.Name, &category.SortOrder, &category.CreatedAt, &category.UpdatedAt); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	return categories, rows.Err()
}

// ListMenuItems returns the full managed catalog.
func (r *MenuRepository) ListMenuItems(ctx context.Context) ([]entity.MenuItem, error) {
	return r.listMenuItems(ctx, `SELECT id, name, description, price_cents, stock_quantity, low_stock_threshold, is_available, is_active, category_id, created_at, updated_at FROM menu_items ORDER BY created_at DESC, id DESC`)
}

// ListAvailableMenuItems returns items customers can order.
func (r *MenuRepository) ListAvailableMenuItems(ctx context.Context) ([]entity.MenuItem, error) {
	return r.listMenuItems(ctx, `
		SELECT id, name, description, price_cents, stock_quantity, low_stock_threshold, is_available, is_active, category_id, created_at, updated_at
		FROM menu_items
		WHERE is_active = TRUE AND is_available = TRUE AND stock_quantity > 0
		ORDER BY name ASC
	`)
}

// GetMenuItemByID returns one menu item by id.
func (r *MenuRepository) GetMenuItemByID(ctx context.Context, id int64) (entity.MenuItem, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, name, description, price_cents, stock_quantity, low_stock_threshold, is_available, is_active, category_id, created_at, updated_at
		FROM menu_items
		WHERE id = $1
	`, id)

	return scanMenuItem(row)
}

// CreateMenuItem inserts a new menu item.
func (r *MenuRepository) CreateMenuItem(ctx context.Context, item entity.MenuItem) (entity.MenuItem, error) {
	now := time.Now().UTC()
	row := r.db.QueryRowContext(ctx, `
		INSERT INTO menu_items (name, description, price_cents, stock_quantity, low_stock_threshold, is_available, is_active, category_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, name, description, price_cents, stock_quantity, low_stock_threshold, is_available, is_active, category_id, created_at, updated_at
	`, item.Name, item.Description, item.PriceCents, item.StockQuantity, item.LowStockThreshold, item.IsAvailable, item.IsActive, item.CategoryID, now, now)

	return scanMenuItem(row)
}

// UpdateMenuItem updates name, description, stock, active state, and category.
func (r *MenuRepository) UpdateMenuItem(ctx context.Context, item entity.MenuItem) (entity.MenuItem, error) {
	row := r.db.QueryRowContext(ctx, `
		UPDATE menu_items
		SET name = $2,
		    description = $3,
		    price_cents = $4,
		    stock_quantity = $5,
		    low_stock_threshold = $6,
		    is_available = $7,
		    is_active = $8,
		    category_id = $9,
		    updated_at = NOW()
		WHERE id = $1
		RETURNING id, name, description, price_cents, stock_quantity, low_stock_threshold, is_available, is_active, category_id, created_at, updated_at
	`, item.ID, item.Name, item.Description, item.PriceCents, item.StockQuantity, item.LowStockThreshold, item.IsAvailable, item.IsActive, item.CategoryID)

	return scanMenuItem(row)
}

// UpdateMenuItemPrice updates only the current price for future orders.
func (r *MenuRepository) UpdateMenuItemPrice(ctx context.Context, id int64, priceCents int64) (entity.MenuItem, error) {
	row := r.db.QueryRowContext(ctx, `
		UPDATE menu_items
		SET price_cents = $2,
		    updated_at = NOW()
		WHERE id = $1
		RETURNING id, name, description, price_cents, stock_quantity, low_stock_threshold, is_available, is_active, category_id, created_at, updated_at
	`, id, priceCents)

	return scanMenuItem(row)
}

// UpdateMenuItemStock updates stock quantity and derived availability.
func (r *MenuRepository) UpdateMenuItemStock(ctx context.Context, id int64, stockQuantity int) (entity.MenuItem, error) {
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

// RecordPriceHistory appends an immutable price change record.
func (r *MenuRepository) RecordPriceHistory(ctx context.Context, history entity.PriceHistory) error {
	if err := history.Validate(); err != nil {
		return err
	}

	row := r.db.QueryRowContext(ctx, `
		INSERT INTO menu_item_price_history (
			menu_item_id,
			previous_price_cents,
			new_price_cents,
			actor,
			reason,
			created_at
		)
		VALUES ($1, $2, $3, $4, NULLIF($5, ''), $6)
		RETURNING id
	`, history.MenuItemID, history.PreviousPriceCents, history.NewPriceCents, history.Actor, history.Reason, history.CreatedAt)

	return row.Scan(&history.ID)
}

// RecordMenuAuditLog appends a manager audit log row.
func (r *MenuRepository) RecordMenuAuditLog(ctx context.Context, log entity.MenuAuditLog) error {
	if err := log.Validate(); err != nil {
		return err
	}

	row := r.db.QueryRowContext(ctx, `
		INSERT INTO menu_audit_logs (
			menu_item_id,
			action_type,
			previous_value,
			new_value,
			actor,
			reason,
			created_at
		)
		VALUES (
			$1,
			$2,
			NULLIF($3, '')::jsonb,
			NULLIF($4, '')::jsonb,
			$5,
			NULLIF($6, ''),
			$7
		)
		RETURNING id
	`, log.MenuItemID, log.ActionType, log.PreviousValue, log.NewValue, log.Actor, log.Reason, log.CreatedAt)

	return row.Scan(&log.ID)
}

// DeactivateMenuItem marks an item inactive without removing it from history.
func (r *MenuRepository) DeactivateMenuItem(ctx context.Context, id int64) (entity.MenuItem, error) {
	row := r.db.QueryRowContext(ctx, `
		UPDATE menu_items
		SET is_active = FALSE,
		    is_available = FALSE,
		    updated_at = NOW()
		WHERE id = $1
		RETURNING id, name, description, price_cents, stock_quantity, low_stock_threshold, is_available, is_active, category_id, created_at, updated_at
	`, id)

	return scanMenuItem(row)
}

func (r *MenuRepository) listMenuItems(ctx context.Context, query string) ([]entity.MenuItem, error) {
	rows, err := r.db.QueryContext(ctx, query)
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

func scanMenuItem(row scanner) (entity.MenuItem, error) {
	var item entity.MenuItem
	if err := row.Scan(
		&item.ID,
		&item.Name,
		&item.Description,
		&item.PriceCents,
		&item.StockQuantity,
		&item.LowStockThreshold,
		&item.IsAvailable,
		&item.IsActive,
		&item.CategoryID,
		&item.CreatedAt,
		&item.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.MenuItem{}, fmt.Errorf("menu item not found")
		}
		return entity.MenuItem{}, err
	}
	return item, nil
}

func scanMenuItemRows(rows *sql.Rows) (entity.MenuItem, error) {
	return scanMenuItem(rows)
}

type scanner interface {
	Scan(dest ...any) error
}
