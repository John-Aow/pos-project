# Quickstart: Manager Menu, Pricing, and Stock Management

## Purpose

Validate that managers can maintain the menu, adjust pricing and stock, and see the customer-facing effects without changing historical bills.

## Prerequisites

- A running frontend app
- A running backend app
- Seed data containing at least one category and a few menu items
- At least one existing historical bill for price-retention validation
- Feature 1 menu/order tables are present before applying Feature 3 migrations
- The `low_stock_threshold`, `menu_item_price_history`, and `menu_audit_logs` migrations have been reviewed before apply

## Validation Commands

From `backend/`:

```bash
go test ./... -coverprofile=coverage.out
```

From `frontend/`:

```bash
npm run typecheck
npm run lint
npm run test -- --coverage
npm run build
```

Backend and frontend coverage must each remain at or above 80%.

## Validation Scenarios

### 1. Review and apply additive migrations

1. Review the migration that adds `low_stock_threshold` to `menu_items`.
2. Review the migration that creates `menu_item_price_history` and `menu_audit_logs`.
3. Confirm both migrations are additive and do not edit, delete, or rewrite existing Feature 1 schema/data.
4. Apply the migrations in a local database.
5. Confirm `menu_item_price_history.menu_item_id` and `menu_audit_logs.menu_item_id` reference `menu_items.id`.

**Expected outcome**: Migrations run successfully, preserve existing menu data, and expose valid foreign keys.

### 2. Create and edit a menu item

1. Open the manager menu screen.
2. Create a new menu item with name, price, description, and initial stock.
3. Edit the item’s description or name.
4. Confirm the item remains visible in manager tools and appears in the customer-facing menu only when active and in stock.

**Expected outcome**: The item is saved, editable, and visible in the correct context.

### 3. Deactivate without deleting

1. Deactivate an existing item from the manager UI or API.
2. Confirm it is not customer-orderable.
3. Query manager menu records and confirm the menu item record still exists.
4. Confirm a `menu_audit_logs` entry was created for the deactivation.

**Expected outcome**: Deactivation is a soft delete only; no menu item or historical record is physically removed.

### 4. Change price for future orders

1. Update the price of an existing item.
2. Save the change through the pricing editor or API.
3. Confirm a `menu_item_price_history` row records the previous price and new price.
4. Start a new order after the change.
5. Confirm the new order uses the updated price.
6. Open an older historical bill or order item for the same item.
7. Confirm the historical `OrderItem.unit_price` still shows the old price snapshot.

**Expected outcome**: New orders use the updated price; old bills keep the original total.

### 5. Set stock to zero and verify availability

1. Reduce an item’s stock quantity to zero.
2. Save the change through the stock editor or API.
3. Confirm `is_available` becomes false automatically.
4. Open the customer-facing menu.
5. Confirm the item is shown as unavailable and cannot be ordered.

**Expected outcome**: Zero stock removes the item from customer ordering automatically.

### 6. Review low-stock warnings

1. Set a few items below their configured `low_stock_threshold`.
2. Open the low-stock warning view.
3. Adjust an item's threshold and confirm the warning list updates.
4. Switch back to the menu view using the in-app navigation.

**Expected outcome**: The warning list shows only active items where `stock_quantity <= low_stock_threshold`.

### 7. Verify audit tracking

1. Create, edit, price-update, stock-update, and deactivate a menu item.
2. Review the audit trail for the item.
3. Confirm each manager action is recorded with actor, action type, affected values, and time.

**Expected outcome**: Catalog changes are traceable for later review.

## Verification Notes

- Confirm manager actions are recorded in the audit trail.
- Confirm the experience is usable on mobile, tablet, desktop, and WebView-sized viewports.
- Confirm the viewport meta tag includes WebView-safe scaling settings.
- Confirm deactivating an item does not erase historical order or bill records.
- Confirm price history is written for every accepted price update.
- Confirm test coverage reports stay at or above 80% for backend and frontend.

## Validation Results

- Backend validation passed on 2026-07-05 with `go test ./... -coverprofile=coverage.out`.
- Backend statement coverage: 81.9%.
- Frontend validation passed on 2026-07-05 with `npm run typecheck`, `npm run lint`, `npm run test -- --coverage`, and `npm run build`.
- Frontend statement coverage: 89.72%.
- Responsive manager menu and low-stock viewport tests passed for mobile, tablet, and desktop layouts.
- No additional migration edits were introduced during the polish phase.
