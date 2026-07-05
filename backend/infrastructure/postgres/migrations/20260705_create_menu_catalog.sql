-- Create menu catalog, inventory settings, and audit history tables.
-- Down migration statements are included for reference alongside the up migration.

CREATE TABLE IF NOT EXISTS categories (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    sort_order INTEGER NOT NULL DEFAULT 0 CHECK (sort_order >= 0)
);

CREATE TABLE IF NOT EXISTS menu_items (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    price_cents BIGINT NOT NULL CHECK (price_cents > 0),
    stock_quantity INTEGER NOT NULL DEFAULT 0 CHECK (stock_quantity >= 0),
    is_available BOOLEAN NOT NULL DEFAULT FALSE,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    category_id BIGINT NOT NULL REFERENCES categories(id) ON UPDATE CASCADE ON DELETE RESTRICT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CHECK (stock_quantity > 0 OR is_available = FALSE)
);

CREATE INDEX IF NOT EXISTS idx_menu_items_category_id ON menu_items(category_id);
CREATE INDEX IF NOT EXISTS idx_menu_items_is_available ON menu_items(is_available);
CREATE INDEX IF NOT EXISTS idx_menu_items_stock_quantity ON menu_items(stock_quantity);

CREATE TABLE IF NOT EXISTS inventory_settings (
    id BIGSERIAL PRIMARY KEY,
    low_stock_threshold INTEGER NOT NULL DEFAULT 5 CHECK (low_stock_threshold >= 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

INSERT INTO inventory_settings (id, low_stock_threshold)
VALUES (1, 5)
ON CONFLICT (id) DO NOTHING;

CREATE TABLE IF NOT EXISTS menu_item_audit_entries (
    id BIGSERIAL PRIMARY KEY,
    menu_item_id BIGINT NOT NULL REFERENCES menu_items(id) ON UPDATE CASCADE ON DELETE RESTRICT,
    action_type TEXT NOT NULL,
    previous_value JSONB,
    new_value JSONB,
    actor TEXT NOT NULL,
    reason TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_menu_item_audit_entries_menu_item_id ON menu_item_audit_entries(menu_item_id);
CREATE INDEX IF NOT EXISTS idx_menu_item_audit_entries_created_at ON menu_item_audit_entries(created_at);
