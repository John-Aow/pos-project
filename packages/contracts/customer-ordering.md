# Customer Ordering Contract

This document mirrors the source contract at
`specs/001-customer-food-ordering/contracts/customer-ordering.openapi.yaml` for
developers working from the shared `packages/contracts/` directory.

## Endpoints

- `GET /api/menu`: returns customer-visible menu items with price, description,
  stock quantity, active state, and availability state.
- `POST /api/orders`: submits a customer cart using an idempotency key and returns
  an order number, status, persisted total, submitted time, and item snapshots.
- `GET /api/orders/{orderNumber}/status`: returns the current customer-visible order
  status.
- `GET /api/staff/orders`: returns incoming and active staff-visible orders.

## Required Client Behavior

- Do not submit an empty cart.
- Generate and reuse an idempotency key for a submit attempt until the result is
  known.
- Treat `400` responses from `POST /api/orders` as customer-correctable validation
  or availability errors.
- Treat `409` responses as duplicate or idempotency conflicts and surface a clear
  recovery path.
- Do not depend on backend implementation details outside the documented API shape.

## Required Server Behavior

- Validate all request input at the HTTP boundary before invoking usecases.
- Recheck item availability and price during submission.
- Persist order totals that match the confirmation response.
- Preserve submitted item names, prices, quantities, and notes as snapshots.
- Keep staff-visible submitted orders discoverable without requiring frontend imports
  from backend implementation code.
