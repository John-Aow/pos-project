# Contract: Staff Order and Billing API

## Active Orders

### `GET /api/staff/orders`

Returns incoming and active orders visible to staff.

**Query fields**
- `status` optional active/incoming status filter

**Response fields**
- `orders`
- `orders[].id`
- `orders[].status`
- `orders[].table_reference`
- `orders[].customer_reference`
- `orders[].items`
- `orders[].notes`
- `orders[].bill`
- `orders[].created_at`
- `orders[].updated_at`

**Rules**
- The response reflects Feature 1 order snapshots.
- Orders should be sortable by newest or operational priority.
- The endpoint is the recovery source after websocket reconnect.

## Checkout

### `POST /api/staff/bills/:bill_id/checkout`

Marks an active bill completed after payment.

**Required fields**
- `actor_id`

**Response fields**
- `bill`
- `status`

**Rules**
- Reject checkout if the bill is already completed, cancelled, or split.
- Completed bills must appear in bill history.
- Duplicate submissions must not create duplicate completed states.

## Cancellation

### `POST /api/staff/bills/:bill_id/cancel`

Cancels an active bill with a required reason.

**Required fields**
- `actor_id`
- `reason`

**Response fields**
- `bill`
- `audit_log`

**Rules**
- Empty or whitespace-only reason returns validation failure.
- Cancellation records `cancelled_reason`, `cancelled_by`, and an audit log.
- Cancelled bills remain visible in history.

## Split

### `POST /api/staff/bills/:bill_id/split`

Splits an active bill into child sub-bills.

**Required fields**
- `actor_id`
- `strategy`

**Strategy fields**
- `parts` for `even`
- `item_allocations` for `by_item`
- `custom_amounts` for `custom`

**Response fields**
- `parent_bill`
- `sub_bills`
- `audit_log`

**Rules**
- `strategy` must be `even`, `by_item`, or `custom`.
- Split into zero bills is rejected.
- Sum of `sub_bills[].total_cents` must equal the original bill total before commit.
- Sub-bill insert and parent status update must occur in one database transaction.
- Even split remainder cents are allocated deterministically.

## History

### `GET /api/staff/bills/history`

Returns completed, cancelled, and split bill records for staff review.

**Query fields**
- `date_from`
- `date_to`
- `status`

**Response fields**
- `bills`
- `bills[].id`
- `bills[].order_id`
- `bills[].parent_bill_id`
- `bills[].status`
- `bills[].total_cents`
- `bills[].cancelled_reason`
- `bills[].cancelled_by`
- `bills[].created_at`
- `bills[].updated_at`

**Rules**
- Date filters use the organization business timezone.
- Status filter supports at least `completed` and `cancelled`.
- Empty result sets return an empty list, not an error.

## Error Expectations

- Input validation failures return a clear error message and do not call usecase mutation logic.
- Business rule failures return a clear conflict or validation response.
- All mutation failures leave bill state unchanged unless the full transaction commits.
