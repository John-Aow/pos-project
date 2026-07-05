# Data Model: Staff Order and Billing Management

## ActiveOrder

Represents an order from Feature 1 that is visible to staff and may need preparation, checkout, cancellation, or split billing.

**Fields**
- `id`
- `status`
- `table_reference`
- `customer_reference`
- `items`
- `notes`
- `created_at`
- `updated_at`

**Validation**
- Must reference an existing Feature 1 order.
- Order detail lines shown to staff must come from the order snapshot and must not be recomputed from current catalog state.

**Relationships**
- Has one active `Bill` when ready for billing.
- Receives realtime visibility through Feature 1 `order_notifier` events.

## Bill

Financial record associated with an order.

**Fields**
- `id`
- `order_id`
- `parent_bill_id`
- `status`
- `total_cents`
- `currency`
- `cancelled_reason`
- `cancelled_by`
- `cancelled_at`
- `checked_out_by`
- `checked_out_at`
- `split_strategy`
- `created_at`
- `updated_at`

**Validation**
- `order_id` must reference an existing order.
- `total_cents` must be zero or greater.
- Cancellation requires a non-empty `cancelled_reason` and `cancelled_by`.
- A checkout can only complete a bill that is active and not already cancelled, completed, or split.
- A bill with child sub-bills preserves its original total.

**Relationships**
- Belongs to one `ActiveOrder`.
- May have many child `SubBill` records through `parent_bill_id`.
- Has many `BillAuditLog` records.

## SubBill

A bill created from splitting an original bill.

**Fields**
- `id`
- `order_id`
- `parent_bill_id`
- `status`
- `total_cents`
- `currency`
- `split_strategy`
- `line_allocations`
- `created_at`
- `updated_at`

**Validation**
- `parent_bill_id` is required.
- Split must create at least one sub-bill.
- Sum of all sub-bill totals must equal the original bill total before commit.
- Even splits allocate remainder cents deterministically.
- By-item splits must allocate valid existing order item lines.
- Custom splits must reject zero or negative sub-bill totals unless the original total is zero.

**Relationships**
- Belongs to one parent `Bill`.
- Shares the parent bill's `order_id`.

## BillAuditLog

Records sensitive bill lifecycle actions.

**Fields**
- `id`
- `bill_id`
- `action`
- `actor_id`
- `reason`
- `metadata`
- `created_at`

**Validation**
- `bill_id`, `action`, `actor_id`, and `created_at` are required.
- `reason` is required for cancellation.
- Split audit metadata must include resulting child bill references.

**Relationships**
- Belongs to one `Bill`.
- References staff identity through `actor_id`.

## BillHistoryFilter

Query input for staff bill history.

**Fields**
- `date_from`
- `date_to`
- `status`

**Validation**
- Date range uses the organization business timezone.
- `date_from` must be before or equal to `date_to` when both are provided.
- Status filter supports at least `completed` and `cancelled`.
- Unknown statuses return validation errors or an empty result according to API contract.

## SplitStrategy

Input model for split operations.

**Fields**
- `strategy`
- `parts`
- `item_allocations`
- `custom_amounts`
- `actor_id`

**Validation**
- `strategy` must be one of `even`, `by_item`, or `custom`.
- `parts` must be greater than zero for even split.
- `item_allocations` are required for by-item split.
- `custom_amounts` are required for custom split.
- Calculated sub-bill total must equal original total exactly.

## State Rules

- Active bills can transition to `completed`, `cancelled`, or `split`.
- Cancelled bills remain visible in historical records and keep cancellation reason.
- Completed bills move to historical records.
- Split parent bills remain preserved and point to sub-bills.
- Sub-bills may be checked out independently after creation.
- Billing actions must never delete existing bills, sub-bills, audit logs, or order snapshots.
