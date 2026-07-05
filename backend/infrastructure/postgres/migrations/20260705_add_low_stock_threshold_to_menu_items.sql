-- Add per-item low-stock threshold support for manager stock warnings.
-- Review before applying: additive column only, preserves existing menu item data.

ALTER TABLE menu_items
    ADD COLUMN IF NOT EXISTS low_stock_threshold INTEGER NOT NULL DEFAULT 5
    CHECK (low_stock_threshold >= 0);

CREATE INDEX IF NOT EXISTS idx_menu_items_low_stock_threshold
    ON menu_items(low_stock_threshold);
