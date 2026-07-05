# Research: Manager Menu, Pricing, and Stock Management

## Decision 1: Use a monorepo with isolated frontend and backend roots

- Decision: Keep the customer/manager UI under `frontend/` and the Go backend under `backend/`, with only shared contracts in `packages/` when necessary.
- Rationale: This matches the constitution, preserves independent build/test paths, and keeps ownership clear for UI and server code.
- Alternatives considered: A single shared app root was rejected because it blurs deployment and testing boundaries; a multi-repo split was rejected because it adds unnecessary coordination overhead for a store product that shares one domain.

## Decision 2: Model price as a snapshot at order time, not a live lookup

- Decision: Use the current `MenuItem.price` for new orders, but store an immutable price snapshot on order items and bills.
- Rationale: This prevents later price edits from changing historical bills while still letting new orders pick up the latest price immediately.
- Alternatives considered: Live price lookup on order history was rejected because it would retroactively rewrite completed financial records; duplicating price history only in audit logs was rejected because billing records still need their own stable totals.

## Decision 3: Treat stock reaching zero as an automatic availability change

- Decision: `stock_quantity <= 0` makes a menu item unavailable to customers automatically, even if the item remains active for manager maintenance.
- Rationale: Customers should never be able to order something with no stock, and the UI should reflect the inventory state without extra manual steps.
- Alternatives considered: Manual availability toggles were rejected because they create drift between stock and ordering availability; background alerts only were rejected because they do not block bad orders.

## Decision 4: Use a low-stock threshold as a configurable inventory setting

- Decision: Maintain a configurable low-stock threshold that drives manager warnings across all managed items.
- Rationale: Managers need a simple operational warning that can be adjusted to match ordering cadence and supplier lead times.
- Alternatives considered: Per-item thresholds were rejected for this feature because the acceptance criteria only require a list of items below a threshold; hard-coded warning levels were rejected because stores vary in replenishment strategy.

## Decision 5: Expose manager changes through audited write actions and read APIs

- Decision: Provide read endpoints for menu and low-stock views, and write endpoints for price, stock, deactivation, and threshold changes, with audit logging on the write path.
- Rationale: The feature needs a clear manager workflow while preserving traceability for operational and financial review.
- Alternatives considered: Direct database updates from the UI were rejected because they bypass validation and auditability; batching all changes into one bulk endpoint was rejected because it makes auditing and failure handling less precise.

## Decision 6: Keep stock updates transactionally consistent with ordering

- Decision: Treat stock quantity as the authoritative inventory value and apply stock-affecting writes inside a transaction so customer availability cannot drift during concurrent activity.
- Rationale: The feature must prevent stock from becoming negative or customer availability from contradicting the saved quantity when orders and manager updates happen close together.
- Alternatives considered: Eventual consistency was rejected because inventory and availability are operationally visible immediately; optimistic-only updates were rejected because they can leave a store with conflicting stock views under load.
