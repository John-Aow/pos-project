# Quickstart: Staff Order and Billing Management

## Purpose

Validate that staff can monitor incoming orders, check out bills, cancel bills with reasoned audit logs, split bills without losing total value, and review bill history across responsive staff layouts.

## Prerequisites

- Feature 1 order schema and `orders` table exist.
- Feature 1 `order_notifier` is available for realtime order events.
- The Feature 2 migration for `bills` and `bill_audit_logs` has been reviewed before apply.
- A running backend app.
- A running frontend app.
- Seed data with at least one active order containing multiple item lines and a known bill total.

## Validation Commands

From `backend/`:

```bash
go test ./...
```

From `frontend/`:

```bash
npm run typecheck
npm run lint
npm run test -- --coverage
npm run build
```

Coverage reports must show at least 80% for backend and frontend.

## Validation Scenarios

### 1. Apply reviewed migration

1. Review the new migration file for `bills` and `bill_audit_logs`.
2. Confirm it only adds new Feature 2 schema and does not edit existing Feature 1 tables or applied migrations.
3. Apply the migration in a local database.
4. Confirm `bills.order_id` references `orders.id`.

**Expected outcome**: Migration runs successfully and foreign keys to orders are valid.

### 2. See incoming orders in real time

1. Start the staff active orders view.
2. Submit a new customer order through the existing order flow or seed script.
3. Confirm the order appears without manual refresh.
4. Open the order and verify items, quantities, notes, and table/customer reference.

**Expected outcome**: Staff see the new order promptly and details match the order snapshot.

### 3. Check out a bill

1. Select an active bill.
2. Submit checkout as a staff actor.
3. Refresh active orders and history.

**Expected outcome**: The bill is completed once, removed from active billing, and visible in history with completed status.

### 4. Cancel a bill with a required reason

1. Open the cancel bill modal.
2. Attempt to confirm without a reason.
3. Enter a reason and confirm.
4. Inspect bill history and audit log records.

**Expected outcome**: Empty reason is blocked. Valid cancellation records reason, actor, status, and audit log.

### 5. Split a bill

1. Open the split bill modal for a bill with multiple items.
2. Test even split where total cents do not divide evenly.
3. Test by-item split.
4. Test custom split with a mismatched total.
5. Test zero sub-bill count.

**Expected outcome**: Valid splits create sub-bills whose totals equal the original. Mismatched totals and zero-bill splits are rejected. Parent status update and sub-bill creation commit together.

### 6. Filter bill history

1. Complete and cancel at least one bill each.
2. Open bill history.
3. Filter by date range.
4. Filter by completed and cancelled status.

**Expected outcome**: Only matching records appear, and empty filters return an empty list without error.

### 7. Verify responsive staff workflows

1. Test active orders at mobile >=360px, tablet >=768px, and desktop >=1024px.
2. Confirm desktop uses table layout and mobile/tablet use card layout.
3. Test split modal on mobile and desktop/tablet.
4. Test in a WebView-sized viewport with portrait and landscape orientation.

**Expected outcome**: Staff can complete view, checkout, cancel, split, and history workflows without layout overlap, inaccessible controls, or browser-chrome-only navigation.

## Verification Notes

- Mutation endpoints must validate input before calling usecase logic.
- Usecase tests should use mock or fake repositories.
- PostgreSQL split implementation must wrap child inserts and parent update in one transaction.
- Websocket events are refresh hints; API reads remain the source of truth.
- Historical and audit records must never be deleted during validation.
