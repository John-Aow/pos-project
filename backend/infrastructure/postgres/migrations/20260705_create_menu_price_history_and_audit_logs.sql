-- Add price history and manager audit logs for menu maintenance.
-- Review before applying: additive tables only, no existing table rewrites.

CREATE TABLE IF NOT EXISTS menu_item_price_history (
    id BIGSERIAL PRIMARY KEY,
    menu_item_id BIGINT NOT NULL REFERENCES menu_items(id) ON UPDATE CASCADE ON DELETE RESTRICT,
    previous_price_cents BIGINT NOT NULL CHECK (previous_price_cents > 0),
    new_price_cents BIGINT NOT NULL CHECK (new_price_cents > 0),
    actor TEXT NOT NULL,
    reason TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CHECK (previous_price_cents <> new_price_cents)
);

CREATE INDEX IF NOT EXISTS idx_menu_item_price_history_menu_item_id
    ON menu_item_price_history(menu_item_id);

CREATE INDEX IF NOT EXISTS idx_menu_item_price_history_created_at
    ON menu_item_price_history(created_at);

CREATE TABLE IF NOT EXISTS menu_audit_logs (
    id BIGSERIAL PRIMARY KEY,
    menu_item_id BIGINT NOT NULL REFERENCES menu_items(id) ON UPDATE CASCADE ON DELETE RESTRICT,
    action_type TEXT NOT NULL,
    previous_value JSONB,
    new_value JSONB,
    actor TEXT NOT NULL,
    reason TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_menu_audit_logs_menu_item_id
    ON menu_audit_logs(menu_item_id);

CREATE INDEX IF NOT EXISTS idx_menu_audit_logs_created_at
    ON menu_audit_logs(created_at);
