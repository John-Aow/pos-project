# Quickstart: Manager Menu, Pricing, and Stock Management

## Purpose

Validate that managers can maintain the menu, adjust pricing and stock, and see the customer-facing effects without changing historical bills.

## Prerequisites

- A running frontend app
- A running backend app
- Seed data containing at least one category and a few menu items
- At least one existing historical bill for price-retention validation

## Validation Scenarios

### 1. Create and edit a menu item

1. Open the manager menu screen.
2. Create a new menu item with name, price, description, and initial stock.
3. Edit the item’s description or name.
4. Confirm the item remains visible in manager tools and appears in the customer-facing menu only when active and in stock.

**Expected outcome**: The item is saved, editable, and visible in the correct context.

### 2. Change price for future orders

1. Update the price of an existing item.
2. Start a new order after the change.
3. Confirm the new order uses the updated price.
4. Open an older historical bill for the same item.
5. Confirm the historical bill still shows the old price snapshot.

**Expected outcome**: New orders use the updated price; old bills keep the original total.

### 3. Set stock to zero and verify availability

1. Reduce an item’s stock quantity to zero.
2. Open the customer-facing menu.
3. Confirm the item is shown as unavailable and cannot be ordered.

**Expected outcome**: Zero stock removes the item from customer ordering automatically.

### 4. Review low-stock warnings

1. Set a few items below the configured threshold.
2. Open the low-stock warning view.
3. Adjust the threshold and confirm the warning list updates.

**Expected outcome**: The warning list shows only items at or below the current threshold.

## Verification Notes

- Confirm manager actions are recorded in the audit trail.
- Confirm the experience is usable on mobile, tablet, desktop, and WebView-sized viewports.
- Confirm deactivating an item does not erase historical order or bill records.
