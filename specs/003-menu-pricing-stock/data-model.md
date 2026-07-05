# Data Model: Manager Menu, Pricing, and Stock Management

## MenuItem

Represents a menu item that managers maintain and customers may order when available.

**Fields**
- `id`
- `name`
- `description`
- `price`
- `stock_quantity`
- `is_available`
- `is_active`
- `category_id`
- `created_at`
- `updated_at`

**Validation**
- Name is required and must be unique within the active menu catalog.
- Price must be non-negative.
- Stock quantity must be zero or greater.
- Deactivated items remain in history but are not customer-orderable.

**Relationships**
- Belongs to one `Category`.
- Is referenced by order items through historical price snapshots, not by live price lookup.

## Category

Groups menu items for customer browsing and manager maintenance.

**Fields**
- `id`
- `name`
- `sort_order`

**Validation**
- Name is required.
- Sort order controls display ordering among categories.

**Relationships**
- Has many `MenuItem` records.

## InventorySettings

Stores configurable stock-warning behavior.

**Fields**
- `id`
- `low_stock_threshold`
- `updated_at`

**Validation**
- Threshold must be zero or greater.

**Relationships**
- Applied across menu items when producing low-stock warnings.

## MenuItemAuditEntry

Records manager changes for menu, price, stock, and deactivation actions.

**Fields**
- `id`
- `menu_item_id`
- `action_type`
- `previous_value`
- `new_value`
- `actor`
- `reason`
- `created_at`

**Validation**
- Each audit entry must identify the actor and action time.
- Reason is required for actions that change availability or price policy, where the UX asks for it.

## OrderItem Snapshot Dependency

Existing order items must preserve `unit_price` as an immutable snapshot taken at order time.

**Fields**
- `order_id`
- `menu_item_id`
- `quantity`
- `unit_price`
- `subtotal`
- `note`

**Validation**
- `unit_price` must never be recalculated from current `MenuItem.price` once the order exists.

**Relationships**
- Belongs to one order and one menu item reference.
- Protects historical bills from later menu price changes.

## State Rules

- If `stock_quantity` is `0`, `is_available` must be false for customer ordering.
- If `stock_quantity` rises above `0` and the item is active, the item may become orderable again.
- Deactivation removes the item from active customer ordering but does not erase historical records.
- Price updates affect only new orders created after the change is accepted.
