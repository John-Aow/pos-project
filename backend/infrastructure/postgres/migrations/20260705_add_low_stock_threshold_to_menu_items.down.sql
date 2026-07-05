DROP INDEX IF EXISTS idx_menu_items_low_stock_threshold;

ALTER TABLE menu_items
    DROP COLUMN IF EXISTS low_stock_threshold;
