# Feature Specification: Staff Order and Billing Management

**Feature Branch**: `002-staff-order-billing`

**Created**: 2026-07-05

**Status**: Draft

**Input**: User description: "Feature: Staff Order and Billing Management. As a staff member, I want to view incoming orders and manage billing (checkout, cancel, split bill), so that I can serve customers accurately and handle payments correctly."

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Monitor Incoming Orders (Priority: P1)

A staff member views the live list of incoming and active orders, then opens an order
to see its details, including ordered items, quantities, notes, and the related
table or customer reference.

**Why this priority**: Staff cannot prepare, charge, or resolve orders accurately
without seeing new work arrive promptly and understanding what each order contains.

**Independent Test**: Can be fully tested by placing a new customer order and
confirming that it appears in the staff list without manual refresh, then opening the
order to verify the details match the original submission.

**Acceptance Scenarios**:

1. **Given** a new customer order is submitted, **When** the staff member views the
   incoming order list, **Then** the new order appears without needing a manual
   refresh.
2. **Given** an order appears in the staff list, **When** the staff member opens it,
   **Then** the order details show the items, quantities, notes, and table or customer
   reference.

---

### User Story 2 - Check Out Or Cancel Bills (Priority: P1)

A staff member closes a bill after payment or cancels a bill when necessary, with a
required reason and confirmation step for cancellation.

**Why this priority**: Closing or cancelling a bill is the core billing responsibility
and directly affects customer service, payment handling, and financial accuracy.

**Independent Test**: Can be fully tested by selecting an active bill, completing the
checkout flow so it moves to completed history, and separately cancelling a bill with
a reason and confirmation.

**Acceptance Scenarios**:

1. **Given** an active bill is ready for payment, **When** the staff member checks it
   out, **Then** the bill is marked completed and becomes available in historical
   records.
2. **Given** a bill is selected for cancellation, **When** the staff member confirms
   the action with a reason, **Then** the bill is cancelled and the reason is
   recorded.
3. **Given** a cancellation reason has not been provided, **When** the staff member
   attempts to cancel a bill, **Then** the cancellation is blocked with a clear
   message.

---

### User Story 3 - Split Bills For Shared Payments (Priority: P2)

A staff member splits a bill into multiple sub-bills, either by item or by amount, so
that separate payments can be handled fairly and accurately.

**Why this priority**: Split billing is a common service need and must preserve the
original total while allowing multiple people to pay their share.

**Independent Test**: Can be fully tested by splitting a bill into two or more
sub-bills and verifying that the total of the sub-bills matches the original bill
total.

**Acceptance Scenarios**:

1. **Given** a bill contains multiple items or a total amount to divide, **When** the
   staff member splits it into two or more sub-bills, **Then** the sub-bills are
   created successfully.
2. **Given** a bill has been split, **When** the staff member reviews the result,
   **Then** the total of the sub-bills equals the original bill total.
3. **Given** the staff member chooses an even or custom split, **When** the split is
   confirmed, **Then** each sub-bill reflects the chosen split method.

---

### User Story 4 - Review Historical Bills (Priority: P2)

A staff member searches and filters historical bills by date range and status to find
past paid or cancelled records.

**Why this priority**: Staff need a reliable way to review completed work, answer
customer questions, and audit past bill activity.

**Independent Test**: Can be fully tested by completing or cancelling bills, then
searching historical records by date and status to confirm the expected records
appear.

**Acceptance Scenarios**:

1. **Given** historical bills exist, **When** the staff member filters by date range,
   **Then** only matching bills appear.
2. **Given** historical bills include paid and cancelled records, **When** the staff
   member filters by status, **Then** only bills with the selected status appear.

### Edge Cases

- A new order arrives while a staff member is already viewing the live list.
- An order appears in the live list before its details are opened.
- A staff member attempts to check out a bill that has already been completed.
- A staff member attempts to cancel a bill without entering a reason.
- A staff member confirms cancellation, then changes their mind before completion.
- A bill is split in a way that would leave remaining balances that do not add up.
- A bill contains notes or item lines that must be preserved across split bills.
- A staff member searches for historical bills across a date range with no results.
- A staff member filters by a status that does not exist in the current records.
- The staff interface is used on mobile, tablet, desktop, portrait orientation,
  landscape orientation, or inside a native app WebView.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The system MUST show staff a live list of incoming and active orders.
- **FR-002**: The staff list MUST update when new orders are submitted without
  requiring a manual refresh.
- **FR-003**: Staff MUST be able to open an order and view its items, quantities,
  notes, and table or customer reference.
- **FR-004**: Staff MUST be able to close an active bill after payment.
- **FR-005**: Closing a bill MUST move it to completed history.
- **FR-006**: Staff MUST be able to cancel a bill only after providing a reason and
  confirming the action.
- **FR-007**: Cancelled bills MUST retain the cancellation reason and remain visible
  in historical records.
- **FR-008**: Staff MUST be able to split one bill into multiple sub-bills.
- **FR-009**: The system MUST support splitting a bill by item or by amount.
- **FR-010**: The system MUST support even and custom bill splits.
- **FR-011**: The total of all sub-bills MUST equal the original bill total.
- **FR-012**: Staff MUST be able to search historical bills by date range.
- **FR-013**: Staff MUST be able to filter historical bills by status, including at
  least paid and cancelled.
- **FR-014**: The system MUST record who cancelled or split a bill, when the action
  occurred, and why.
- **FR-015**: The system MUST preserve the original bill record and its history when
  bills are cancelled or split.
- **FR-016**: The staff billing experience MUST remain usable across mobile, tablet,
  desktop, portrait, landscape, and native WebView contexts when accessed through a
  web interface.

### POS Risk Requirements *(include when applicable)*

- **Transaction Integrity**: Checkout, cancellation, and split actions MUST preserve
  correct totals and must not create or lose value when bills are closed or divided.
- **Catalog/Inventory Accuracy**: The staff view MUST reflect the same order details
  and references that were captured at submission time.
- **Offline Operation**: If the staff interface cannot confirm a billing action
  immediately, the system MUST avoid duplicate checkout, cancellation, or split
  actions and MUST make the current state clear to staff.
- **Security/Audit**: Cancellation and split actions MUST be logged with staff
  identity, timestamp, reason, and resulting bill references.
- **Data Safety**: Billing actions MUST not delete completed, cancelled, or split bill
  history.
- **Staff Visibility**: Staff MUST be able to locate active orders and historical
  bills without depending on hidden or manual back-office processes.
- **Responsive/WebView**: Staff controls for viewing, checking out, cancelling, and
  splitting bills MUST work across mobile, tablet, desktop, portrait, landscape, and
  native WebView contexts.

### Key Entities *(include if feature involves data)*

- **Active Order**: An order currently waiting for preparation, payment, or other
  staff action.
- **Bill**: The financial record associated with an order, including totals and
  status.
- **Sub-Bill**: A bill created when an original bill is split into multiple payable
  parts.
- **Historical Bill Record**: A completed or cancelled bill kept for review and
  audit.
- **Cancellation Record**: The reasoned log of a bill cancellation, including who
  performed it and when.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: 100% of newly submitted customer orders appear in the staff list without
  manual refresh during testing.
- **SC-002**: Staff can open an incoming order and confirm its details in under 10
  seconds.
- **SC-003**: 100% of completed bills move into historical records with a completed
  status after checkout.
- **SC-004**: 100% of bill cancellations require a reason before completion and are
  recorded in the audit trail.
- **SC-005**: 100% of split bills preserve totals so the combined sub-bill total
  equals the original bill total.
- **SC-006**: Staff can search or filter historical bills by date range and status and
  find matching records in under 10 seconds.
- **SC-007**: Staff can complete the main billing actions on at least one mobile, one
  tablet, and one desktop viewport during usability testing.

## Assumptions

- Staff users are already authorized to manage orders and billing.
- Payment collection is handled within the existing bill checkout process.
- Split bills may include whole items, portions of the bill total, or both depending
  on the selected split method.
- Historical bill records are retained rather than deleted when bills are completed or
  cancelled.
- The order details shown to staff include any table identifier or customer reference
  already associated with the order.
- Date filters use the organization’s standard business timezone.
- The staff interface is web-based and may be accessed on browser, tablet, or native
  WebView surfaces.
