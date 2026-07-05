# Contract: Menu, Pricing, and Stock Management API

## Customer-Facing Read Contract

### `GET /api/menu`

Returns only customer-orderable menu items.

**Response fields**
- `id`
- `name`
- `description`
- `price`
- `price_cents`
- `stock_quantity`
- `is_available`
- `category`

**Rules**
- Unavailable or out-of-stock items must not be returned as orderable.
- The returned price is the current active price for new orders.

## Manager Menu Contract

### `GET /api/manager/menu`

Returns the full managed menu, including active state and stock information.

**Response fields**
- `id`
- `name`
- `description`
- `price_cents`
- `stock_quantity`
- `low_stock_threshold`
- `is_available`
- `is_active`
- `category_id`
- `created_at`
- `updated_at`

### `POST /api/manager/menu`

Creates a new menu item.

**Required fields**
- `name`
- `description`
- `price`
- `stock_quantity`
- `low_stock_threshold`
- `category_id`

### `PATCH /api/manager/menu/:id`

Updates menu item details such as name, description, or active state.

### `PATCH /api/manager/menu/:id/price`

Updates the current price for future orders.

**Required fields**
- `price_cents`
- `actor`

**Rules**
- Historical bills retain the old price snapshot.
- Every accepted price update creates a `menu_item_price_history` record with `menu_item_id`, previous price, new price, actor, and timestamp.
- Previously created `OrderItem.unit_price` values must not be updated by this endpoint.

### `PATCH /api/manager/menu/:id/stock`

Updates stock quantity.

**Required fields**
- `stock_quantity`
- `actor`

**Rules**
- When stock becomes `0`, customer ordering must treat the item as unavailable.
- Stock changes create a `menu_audit_logs` entry with actor, action type, affected item, previous value, new value, and timestamp.

### `GET /api/manager/menu/low-stock`

Returns items at or below their configured low-stock threshold.

**Response fields**
- `items`
- `items[].id`
- `items[].name`
- `items[].stock_quantity`
- `items[].low_stock_threshold`
- `items[].is_active`

**Rules**
- Items are sorted by stock quantity, then name.
- Only active menu items are returned.
- An item is low stock when `stock_quantity <= low_stock_threshold`.

### `PATCH /api/manager/settings/low-stock-threshold`

Updates the default warning threshold used when creating or bulk-managing menu items.

**Required fields**
- `low_stock_threshold`
- `actor`

**Optional fields**
- `reason`

**Response fields**
- `threshold`

## Price History Contract

`menu_item_price_history` records must include:
- `id`
- `menu_item_id`
- `previous_price_cents`
- `new_price_cents`
- `actor`
- `created_at`

**Rules**
- `menu_item_id` must reference `menu_items.id`.
- A price history row is written for every accepted price update.
- History rows are append-only and must not be used to recompute old order item prices.

## Audit Contract

`menu_audit_logs` records must include:
- `id`
- `menu_item_id`
- actor
- action type
- previous value
- new value
- reason when provided or required
- timestamp

**Rules**
- `menu_item_id` must reference `menu_items.id`.
- Create, update, price update, stock update, deactivate, and threshold changes must be auditable.
- Deactivation is a soft delete only: it marks the item unavailable/inactive for ordering and never removes historical records.

## Response Expectations

- Input validation failures return a clear error message.
- Concurrency conflicts return a conflict response rather than silently overwriting the latest stock or price state.
- Historical records never recompute old prices from the current menu item.
