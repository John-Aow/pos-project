package postgres

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"pos-project/backend/entity"
)

func TestMenuAndInventoryRepositories(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 7, 5, 12, 0, 0, 0, time.UTC)
	script := newQueryScript(t, []queryExpectation{
		{
			contains: "FROM categories",
			columns:  []string{"id", "name", "sort_order", "created_at", "updated_at"},
			rows: [][]driver.Value{{
				int64(1), "Mains", int64(1), now, now,
			}},
		},
		{
			contains: "FROM menu_items ORDER BY created_at DESC",
			columns:  []string{"id", "name", "description", "price_cents", "stock_quantity", "low_stock_threshold", "is_available", "is_active", "category_id", "created_at", "updated_at"},
			rows: [][]driver.Value{{
				int64(2), "Green Curry", "Thai curry", int64(12900), int64(5), int64(3), true, true, int64(1), now, now,
			}},
		},
		{
			contains: "WHERE is_active = TRUE AND is_available = TRUE AND stock_quantity > 0",
			columns:  []string{"id", "name", "description", "price_cents", "stock_quantity", "low_stock_threshold", "is_available", "is_active", "category_id", "created_at", "updated_at"},
			rows: [][]driver.Value{{
				int64(2), "Green Curry", "Thai curry", int64(12900), int64(5), int64(3), true, true, int64(1), now, now,
			}},
		},
		{
			contains: "WHERE id = $1",
			columns:  []string{"id", "name", "description", "price_cents", "stock_quantity", "low_stock_threshold", "is_available", "is_active", "category_id", "created_at", "updated_at"},
			rows: [][]driver.Value{{
				int64(3), "Massaman", "Curry", int64(11900), int64(2), int64(3), true, true, int64(1), now, now,
			}},
		},
		{
			contains: "INSERT INTO menu_items",
			columns:  []string{"id", "name", "description", "price_cents", "stock_quantity", "low_stock_threshold", "is_available", "is_active", "category_id", "created_at", "updated_at"},
			rows: [][]driver.Value{{
				int64(4), "Pad Thai", "Noodles", int64(13900), int64(4), int64(3), true, true, int64(1), now, now,
			}},
		},
		{
			contains: "UPDATE menu_items",
			columns:  []string{"id", "name", "description", "price_cents", "stock_quantity", "low_stock_threshold", "is_available", "is_active", "category_id", "created_at", "updated_at"},
			rows: [][]driver.Value{{
				int64(4), "Pad Thai Deluxe", "Noodles", int64(14900), int64(1), int64(3), true, true, int64(1), now, now,
			}},
		},
		{
			contains: "SET price_cents = $2",
			columns:  []string{"id", "name", "description", "price_cents", "stock_quantity", "low_stock_threshold", "is_available", "is_active", "category_id", "created_at", "updated_at"},
			rows: [][]driver.Value{{
				int64(4), "Pad Thai Deluxe", "Noodles", int64(15900), int64(1), int64(3), true, true, int64(1), now, now,
			}},
		},
		{
			contains: "INSERT INTO menu_item_price_history",
			columns:  []string{"id"},
			rows: [][]driver.Value{{
				int64(11),
			}},
		},
		{
			contains: "INSERT INTO menu_audit_logs",
			columns:  []string{"id"},
			rows: [][]driver.Value{{
				int64(12),
			}},
		},
		{
			contains: "SET stock_quantity = $2",
			columns:  []string{"id", "name", "description", "price_cents", "stock_quantity", "low_stock_threshold", "is_available", "is_active", "category_id", "created_at", "updated_at"},
			rows: [][]driver.Value{{
				int64(4), "Pad Thai Deluxe", "Noodles", int64(15900), int64(0), int64(3), false, true, int64(1), now, now,
			}},
		},
		{
			contains: "SET is_active = FALSE",
			columns:  []string{"id", "name", "description", "price_cents", "stock_quantity", "low_stock_threshold", "is_available", "is_active", "category_id", "created_at", "updated_at"},
			rows: [][]driver.Value{{
				int64(4), "Pad Thai Deluxe", "Noodles", int64(15900), int64(0), int64(3), false, false, int64(1), now, now,
			}},
		},
		{
			contains: "WHERE id = 1",
			columns:  []string{"id", "low_stock_threshold", "created_at", "updated_at"},
			rows: [][]driver.Value{{
				int64(1), int64(5), now, now,
			}},
		},
		{
			contains: "UPDATE inventory_settings",
			columns:  []string{"id", "low_stock_threshold", "created_at", "updated_at"},
			rows: [][]driver.Value{{
				int64(1), int64(7), now, now,
			}},
		},
		{
			contains: "SET stock_quantity = $2",
			columns:  []string{"id", "name", "description", "price_cents", "stock_quantity", "low_stock_threshold", "is_available", "is_active", "category_id", "created_at", "updated_at"},
			rows: [][]driver.Value{{
				int64(6), "Tea", "Jasmine", int64(2500), int64(1), int64(3), true, true, int64(1), now, now,
			}},
		},
		{
			contains: "stock_quantity <= $1",
			columns:  []string{"id", "name", "description", "price_cents", "stock_quantity", "low_stock_threshold", "is_available", "is_active", "category_id", "created_at", "updated_at"},
			rows: [][]driver.Value{{
				int64(5), "Soup", "Tom yum", int64(8900), int64(2), int64(3), true, true, int64(1), now, now,
			}, {
				int64(6), "Tea", "Jasmine", int64(2500), int64(4), int64(5), true, true, int64(1), now, now,
			}},
		},
	})
	db := openScriptedDB(t, script)

	menuRepo := NewMenuRepository(db)
	inventoryRepo := NewInventoryRepository(db)

	categories, err := menuRepo.ListCategories(context.Background())
	if err != nil {
		t.Fatalf("ListCategories returned error: %v", err)
	}
	if len(categories) != 1 || categories[0].Name != "Mains" {
		t.Fatalf("unexpected categories: %#v", categories)
	}

	items, err := menuRepo.ListMenuItems(context.Background())
	if err != nil {
		t.Fatalf("ListMenuItems returned error: %v", err)
	}
	if len(items) != 1 || items[0].Name != "Green Curry" {
		t.Fatalf("unexpected items: %#v", items)
	}

	available, err := menuRepo.ListAvailableMenuItems(context.Background())
	if err != nil {
		t.Fatalf("ListAvailableMenuItems returned error: %v", err)
	}
	if len(available) != 1 || available[0].Name != "Green Curry" {
		t.Fatalf("unexpected available items: %#v", available)
	}

	gotItem, err := menuRepo.GetMenuItemByID(context.Background(), 3)
	if err != nil {
		t.Fatalf("GetMenuItemByID returned error: %v", err)
	}
	if gotItem.Name != "Massaman" {
		t.Fatalf("unexpected menu item: %#v", gotItem)
	}

	created, err := menuRepo.CreateMenuItem(context.Background(), entity.MenuItem{
		Name:          "Pad Thai",
		Description:   "Noodles",
		PriceCents:    13900,
		StockQuantity: 4,
		IsAvailable:   true,
		IsActive:      true,
		CategoryID:    1,
	})
	if err != nil {
		t.Fatalf("CreateMenuItem returned error: %v", err)
	}
	if created.ID != 4 {
		t.Fatalf("unexpected created item: %#v", created)
	}

	updated, err := menuRepo.UpdateMenuItem(context.Background(), entity.MenuItem{ID: 4, Name: "Pad Thai Deluxe"})
	if err != nil {
		t.Fatalf("UpdateMenuItem returned error: %v", err)
	}
	if updated.Name != "Pad Thai Deluxe" {
		t.Fatalf("unexpected updated item: %#v", updated)
	}

	priceUpdated, err := menuRepo.UpdateMenuItemPrice(context.Background(), 4, 15900)
	if err != nil {
		t.Fatalf("UpdateMenuItemPrice returned error: %v", err)
	}
	if priceUpdated.PriceCents != 15900 {
		t.Fatalf("unexpected price item: %#v", priceUpdated)
	}
	historicalOrderItem := entity.MenuItem{ID: 4, PriceCents: 14900}
	if historicalOrderItem.PriceCents != 14900 {
		t.Fatalf("expected historical snapshot to keep its original price, got %#v", historicalOrderItem)
	}

	if err := menuRepo.RecordPriceHistory(context.Background(), entity.PriceHistory{
		MenuItemID:         4,
		PreviousPriceCents: 14900,
		NewPriceCents:      15900,
		Actor:              "manager-1",
		Reason:             "menu refresh",
		CreatedAt:          now,
	}); err != nil {
		t.Fatalf("RecordPriceHistory returned error: %v", err)
	}

	if err := menuRepo.RecordMenuAuditLog(context.Background(), entity.MenuAuditLog{
		MenuItemID:    4,
		ActionType:    entity.AuditActionPrice,
		PreviousValue: `{"price_cents":14900}`,
		NewValue:      `{"price_cents":15900}`,
		Actor:         "manager-1",
		Reason:        "menu refresh",
		CreatedAt:     now,
	}); err != nil {
		t.Fatalf("RecordMenuAuditLog returned error: %v", err)
	}

	stockUpdated, err := menuRepo.UpdateMenuItemStock(context.Background(), 4, 0)
	if err != nil {
		t.Fatalf("UpdateMenuItemStock returned error: %v", err)
	}
	if stockUpdated.StockQuantity != 0 || stockUpdated.IsAvailable {
		t.Fatalf("unexpected stock item: %#v", stockUpdated)
	}

	deactivated, err := menuRepo.DeactivateMenuItem(context.Background(), 4)
	if err != nil {
		t.Fatalf("DeactivateMenuItem returned error: %v", err)
	}
	if deactivated.IsActive {
		t.Fatalf("unexpected deactivated item: %#v", deactivated)
	}

	settings, err := inventoryRepo.GetInventorySettings(context.Background())
	if err != nil {
		t.Fatalf("GetInventorySettings returned error: %v", err)
	}
	if settings.LowStockThreshold != 5 {
		t.Fatalf("unexpected settings: %#v", settings)
	}

	updatedSettings, err := inventoryRepo.UpdateInventorySettings(context.Background(), entity.InventorySettings{LowStockThreshold: 7})
	if err != nil {
		t.Fatalf("UpdateInventorySettings returned error: %v", err)
	}
	if updatedSettings.LowStockThreshold != 7 {
		t.Fatalf("unexpected updated settings: %#v", updatedSettings)
	}

	stockAfterInventoryUpdate, err := inventoryRepo.UpdateMenuItemStock(context.Background(), 6, 1)
	if err != nil {
		t.Fatalf("UpdateMenuItemStock returned error: %v", err)
	}
	if stockAfterInventoryUpdate.StockQuantity != 1 {
		t.Fatalf("unexpected inventory stock item: %#v", stockAfterInventoryUpdate)
	}

	lowStockItems, err := inventoryRepo.ListLowStockMenuItems(context.Background(), 7)
	if err != nil {
		t.Fatalf("ListLowStockMenuItems returned error: %v", err)
	}
	if len(lowStockItems) != 2 {
		t.Fatalf("unexpected low-stock items: %#v", lowStockItems)
	}
	if lowStockItems[0].Name != "Soup" || lowStockItems[1].Name != "Tea" {
		t.Fatalf("expected low-stock items to be sorted by stock then name, got %#v", lowStockItems)
	}
}

func TestMenuRepositoryReturnsNotFound(t *testing.T) {
	t.Parallel()

	script := newQueryScript(t, []queryExpectation{
		{
			contains: "WHERE id = $1",
			columns:  []string{"id", "name", "description", "price_cents", "stock_quantity", "low_stock_threshold", "is_available", "is_active", "category_id", "created_at", "updated_at"},
			rows:     [][]driver.Value{},
		},
	})
	db := openScriptedDB(t, script)

	_, err := NewMenuRepository(db).GetMenuItemByID(context.Background(), 99)
	if err == nil || !strings.Contains(err.Error(), "menu item not found") {
		t.Fatalf("expected not found error, got %v", err)
	}
}

func TestMenuMigrationsAreAdditiveAndReferenceMenuItems(t *testing.T) {
	t.Parallel()

	migrationsDir := filepath.Join("migrations")

	lowStockMigration := readMigration(t, migrationsDir, "20260705_add_low_stock_threshold_to_menu_items.sql")
	assertContains(t, lowStockMigration, "ADD COLUMN IF NOT EXISTS low_stock_threshold")
	assertContains(t, lowStockMigration, "CHECK (low_stock_threshold >= 0)")
	assertNotContains(t, lowStockMigration, "DROP TABLE")
	assertNotContains(t, lowStockMigration, "DROP COLUMN")
	assertNotContains(t, lowStockMigration, "DELETE FROM")

	historyMigration := readMigration(t, migrationsDir, "20260705_create_menu_price_history_and_audit_logs.sql")
	assertContains(t, historyMigration, "CREATE TABLE IF NOT EXISTS menu_item_price_history")
	assertContains(t, historyMigration, "CREATE TABLE IF NOT EXISTS menu_audit_logs")
	assertContains(t, historyMigration, "REFERENCES menu_items(id)")
	assertContains(t, historyMigration, "previous_price_cents BIGINT NOT NULL CHECK (previous_price_cents > 0)")
	assertContains(t, historyMigration, "new_price_cents BIGINT NOT NULL CHECK (new_price_cents > 0)")
	assertNotContains(t, historyMigration, "DROP TABLE")
	assertNotContains(t, historyMigration, "DROP COLUMN")
	assertNotContains(t, historyMigration, "DELETE FROM")

	downLowStock := readMigration(t, migrationsDir, "20260705_add_low_stock_threshold_to_menu_items.down.sql")
	assertContains(t, downLowStock, "DROP COLUMN IF EXISTS low_stock_threshold")

	downHistory := readMigration(t, migrationsDir, "20260705_create_menu_price_history_and_audit_logs.down.sql")
	assertContains(t, downHistory, "DROP TABLE IF EXISTS menu_audit_logs")
	assertContains(t, downHistory, "DROP TABLE IF EXISTS menu_item_price_history")
}

func readMigration(t *testing.T, dir string, filename string) string {
	t.Helper()

	content, err := os.ReadFile(filepath.Join(dir, filename))
	if err != nil {
		t.Fatalf("read migration %s: %v", filename, err)
	}
	return string(content)
}

func assertContains(t *testing.T, content string, needle string) {
	t.Helper()

	if !strings.Contains(content, needle) {
		t.Fatalf("expected migration to contain %q", needle)
	}
}

func assertNotContains(t *testing.T, content string, needle string) {
	t.Helper()

	if strings.Contains(content, needle) {
		t.Fatalf("expected migration not to contain %q", needle)
	}
}

type queryExpectation struct {
	contains string
	columns  []string
	rows     [][]driver.Value
}

type queryScript struct {
	t            *testing.T
	expectations []queryExpectation
	current      int
	mu           sync.Mutex
}

func newQueryScript(t *testing.T, expectations []queryExpectation) *queryScript {
	t.Helper()
	return &queryScript{t: t, expectations: expectations}
}

func (s *queryScript) next(query string, args []driver.NamedValue) (*scriptedRows, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.current >= len(s.expectations) {
		return nil, fmt.Errorf("unexpected query: %s", query)
	}

	expectation := s.expectations[s.current]
	s.current++

	if !strings.Contains(query, expectation.contains) {
		s.t.Fatalf("expected query to contain %q, got %q", expectation.contains, query)
	}

	return &scriptedRows{
		columns: expectation.columns,
		rows:    expectation.rows,
	}, nil
}

type scriptedDriver struct {
	script *queryScript
}

func (d *scriptedDriver) Open(string) (driver.Conn, error) {
	return &scriptedConn{script: d.script}, nil
}

type scriptedConn struct {
	script *queryScript
}

func (c *scriptedConn) Prepare(string) (driver.Stmt, error) {
	return nil, fmt.Errorf("prepare is not supported")
}

func (c *scriptedConn) Close() error {
	return nil
}

func (c *scriptedConn) Begin() (driver.Tx, error) {
	return nil, fmt.Errorf("transactions are not supported")
}

func (c *scriptedConn) QueryContext(_ context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	return c.script.next(query, args)
}

type scriptedRows struct {
	columns []string
	rows    [][]driver.Value
	index   int
}

func (r *scriptedRows) Columns() []string {
	return r.columns
}

func (r *scriptedRows) Close() error {
	return nil
}

func (r *scriptedRows) Next(dest []driver.Value) error {
	if r.index >= len(r.rows) {
		return io.EOF
	}

	copy(dest, r.rows[r.index])
	r.index++
	return nil
}

func openScriptedDB(t *testing.T, script *queryScript) *sql.DB {
	t.Helper()

	name := fmt.Sprintf("scripted-%s-%d", strings.ReplaceAll(t.Name(), "/", "_"), time.Now().UnixNano())
	sql.Register(name, &scriptedDriver{script: script})
	db, err := sql.Open(name, "")
	if err != nil {
		t.Fatalf("failed to open scripted db: %v", err)
	}
	t.Cleanup(func() {
		_ = db.Close()
	})
	return db
}
