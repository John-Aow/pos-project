# Contract: Menu, Pricing, and Stock Management API

## Customer-Facing Read Contract

### `GET /api/menu`

Returns only customer-orderable menu items.

**Response fields**
- `id`
- `name`
- `description`
- `price`
- `stock_quantity`
- `is_available`
- `category`

**Rules**
- Unavailable or out-of-stock items must not be returned as orderable.
- The returned price is the current active price for new orders.

## Manager Menu Contract

### `GET /api/manager/menu`

Returns the full managed menu, including active state and stock information.

### `POST /api/manager/menu`

Creates a new menu item.

**Required fields**
- `name`
- `description`
- `price`
- `stock_quantity`
- `category_id`

### `PATCH /api/manager/menu/:id`

Updates menu item details such as name, description, or active state.

### `PATCH /api/manager/menu/:id/price`

Updates the current price for future orders.

**Rules**
- Historical bills retain the old price snapshot.

### `PATCH /api/manager/menu/:id/stock`

Updates stock quantity.

**Rules**
- When stock becomes `0`, customer ordering must treat the item as unavailable.

### `GET /api/manager/menu/low-stock`

Returns items at or below the current low-stock threshold.

### `PATCH /api/manager/settings/low-stock-threshold`

Updates the warning threshold used for low-stock reporting.

## Audit Contract

Write actions must record:
- actor
- action type
- affected item
- previous value
- new value
- reason when provided or required
- timestamp

## Response Expectations

- Input validation failures return a clear error message.
- Concurrency conflicts return a conflict response rather than silently overwriting the latest stock or price state.
- Historical records never recompute old prices from the current menu item.
