# Research: Staff Order and Billing Management

## Decision 1: Add billing tables with reviewed, additive migrations only

- Decision: Create new migrations for `bills` and `bill_audit_logs`, including `parent_bill_id`, `cancelled_reason`, and `cancelled_by`, with foreign keys to Feature 1 `orders`.
- Rationale: Billing state and audit history are new Feature 2 concepts, and the constitution requires schema changes to be additive and reviewed before apply. Referencing `orders` preserves traceability from order intake to financial record.
- Alternatives considered: Reusing only the `orders` table was rejected because checkout/cancel/split state would blur operational order state with financial bill state. Editing existing Feature 1 migrations was rejected by database governance.

## Decision 2: Keep monetary rules in entity and usecase layers

- Decision: Model `Bill` and `BillAuditLog` entities with validation for required cancellation reason and exact split total preservation, then enforce transaction orchestration in usecases.
- Rationale: Clean Architecture keeps business rules testable without HTTP or PostgreSQL. Entity tests can cover rule failures quickly, while usecase tests can verify repository calls and transaction behavior.
- Alternatives considered: Validating only in HTTP handlers was rejected because repository or background callers could bypass rules. Validating only in SQL constraints was rejected because split strategy and audit behavior need domain-aware errors.

## Decision 3: Represent split bills as child bills linked by `parent_bill_id`

- Decision: Preserve the original bill as the parent, create sub-bills with `parent_bill_id`, and update the original status to indicate it has been split after all child bills are inserted successfully.
- Rationale: This preserves the original financial record and makes sub-bills auditable without deleting or mutating historical totals. A transaction ensures no partial split is committed.
- Alternatives considered: Replacing the original bill with multiple rows was rejected because it loses the original record. Storing split details only as JSON was rejected because child bills need independent checkout/history behavior.

## Decision 4: Use deterministic cent allocation for even splits

- Decision: Store money in integer minor units and allocate any remainder cents deterministically, such as assigning earlier sub-bills one extra cent until the remainder is exhausted.
- Rationale: Even splits may not divide cleanly, but financial totals must remain exact. Deterministic allocation makes tests, receipts, and audit trails reproducible.
- Alternatives considered: Floating point division was rejected because it risks rounding drift. Rejecting all non-divisible splits was rejected because common totals cannot always divide evenly.

## Decision 5: Make cancellation audit mandatory and reasoned

- Decision: `CancelBill(reason)` rejects empty reasons and always writes a `bill_audit_logs` entry with actor, reason, action, timestamp, and bill reference when cancellation succeeds.
- Rationale: Cancellation changes financial history and must be reviewable. The feature and constitution both require reasoned audit trails for sensitive staff actions.
- Alternatives considered: Optional cancellation comments were rejected because acceptance requires a reason. UI-only confirmation was rejected because backend callers must follow the same rule.

## Decision 6: Expose staff billing through validated HTTP contracts

- Decision: Provide `GET /api/staff/orders`, bill checkout/cancel/split endpoints, and `GET /api/staff/bills/history` with date/status filters.
- Rationale: Explicit staff contracts let frontend and backend remain independent and give integration tests stable request/response expectations.
- Alternatives considered: A generic bill mutation endpoint was rejected because checkout, cancellation, and split have different validation and audit rules. Direct frontend database access is outside project architecture.

## Decision 7: Use order notifier events as realtime hints, with API refresh as source of truth

- Decision: The websocket listener consumes Feature 1 `order_notifier` events to update the staff UI, while `GET /api/staff/orders` remains the authoritative recovery path after reconnect or missed events.
- Rationale: Staff need realtime order visibility, but network or websocket loss must not corrupt billing state. Treating events as hints keeps the UI fresh and recoverable.
- Alternatives considered: Polling only was rejected because acceptance requires realtime visibility without manual refresh. Websocket-only state was rejected because reconnect and missed-event recovery need a durable read API.

## Decision 8: Use responsive staff-specific layouts

- Decision: Active orders use a dense table on desktop and scan-friendly cards on mobile/tablet; split is a full-screen modal on mobile and dialog on wider viewports.
- Rationale: Staff may work on phones, tablets, desktops, or WebViews. The controls must be efficient for repeated use while preserving touch target and viewport requirements.
- Alternatives considered: A single desktop table was rejected because it fails mobile/WebView ergonomics. A card-only layout was rejected because desktop staff workflows benefit from dense scanning.
