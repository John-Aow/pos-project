# Implementation Plan: Customer Food Ordering

**Feature Branch**: `001-customer-food-ordering`
**Spec**: `specs/001-customer-food-ordering/spec.md`
**Tasks**: `specs/001-customer-food-ordering/tasks.md`
**Date**: 2026-07-05

## Summary

Build a customer-facing web ordering flow that lets customers browse available menu
items, build a cart with item notes, submit an order, receive an order number, and
view the current order status. The implementation will extend the existing
Vue 3 frontend and Go backend while preserving frontend/backend independence.

## Technical Context

**Frontend**: Vue 3, TypeScript, Composition API, Pinia, Vite, Vitest, Vue Test Utils

**Backend**: Go, standard `net/http`, Clean Architecture layers:
`entity`, `usecase`, `interface/adapter`, `infrastructure`

**Database**: PostgreSQL migrations under
`backend/infrastructure/postgres/migrations/`

**Shared Contracts**: Human-readable and OpenAPI contracts under
`packages/contracts/` and `specs/001-customer-food-ordering/contracts/`

**Testing**: Frontend and backend tests must remain independently runnable and keep
coverage at or above 80%.

## Constitution Check

- **Transaction Integrity**: Order totals are calculated from selected available
  items using submission-time prices and persisted with the created bill/order.
- **Inventory And Catalog Accuracy**: Availability is checked when browsing, before
  cart add, and again during order submission.
- **Offline-Resilient Checkout**: The frontend uses an idempotency key for submit
  attempts and surfaces uncertain submission state instead of silently retrying.
- **Security, Privacy, And Auditability**: APIs validate input before usecases.
  Submitted orders store order number, submitted time, item selections, quantities,
  notes, total, and status changes. No payment card data is stored.
- **Spec-First, Tested Incremental Delivery**: Tasks are split by user story with
  tests before implementation.
- **Monorepo Independence**: Frontend consumes backend APIs through documented
  contracts only. No direct implementation imports cross app boundaries.
- **Clean Backend Architecture**: Entities contain validation, usecases depend on
  repository interfaces, adapters handle HTTP validation, and infrastructure handles
  PostgreSQL persistence.
- **Responsive & Multi-Device Support**: Customer UI is validated at mobile 360px,
  tablet 768px, desktop 1024px, portrait and landscape, with 44px touch targets.
- **Database Governance**: New schema work is introduced only by new migration files.
  Existing migrations are not edited.

## Project Structure

```text
frontend/
  src/
    components/
    pages/
    services/
    stores/
    styles/
  tests/
backend/
  entity/
  usecase/
  interface/adapter/http/
  infrastructure/postgres/
  tests/
packages/
  contracts/
specs/001-customer-food-ordering/
  contracts/
```

## API Surface

- `GET /api/menu`: list customer-visible menu items
- `POST /api/orders`: submit a customer order
- `GET /api/orders/{orderNumber}/status`: read current customer-visible order status
- `GET /api/staff/orders`: list incoming staff-visible orders

## Implementation Phases

1. Setup design artifacts and contracts.
2. Add shared domain entities, persistence, routes, and frontend shell.
3. Implement customer menu browsing.
4. Implement cart building and item notes.
5. Implement order submission, idempotency, confirmation, and staff visibility.
6. Implement current order status.
7. Run quality gates, quickstart, responsive checks, and coverage.

## Risks And Mitigations

| Risk | Mitigation |
|------|------------|
| Item price or stock changes during cart building | Recheck catalog during submission and return item-level errors |
| Duplicate order submission after network loss | Require frontend-generated idempotency key on `POST /api/orders` |
| Staff list misses submitted order | Persist order status as `pending` and expose `GET /api/staff/orders` sorted by submission time |
| WebView hides browser navigation | Provide in-app navigation between menu, cart, confirmation, and status |
| Existing database history is altered | Use only new migrations; do not edit or delete existing migration files |

## Validation Gates

- Frontend: lint, tests, responsive/WebView component tests, coverage >= 80%
- Backend: gofmt, tests, usecase/API/migration coverage >= 80%
- Quickstart: customer can browse, add item, submit order, receive order number,
  staff can see order, customer can view status
