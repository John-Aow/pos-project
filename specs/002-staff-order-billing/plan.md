# Implementation Plan: Staff Order and Billing Management

**Branch**: `002-staff-order-billing` | **Date**: 2026-07-05 | **Spec**: [/Users/johnhu/Desktop/pos-project/specs/002-staff-order-billing/spec.md](/Users/johnhu/Desktop/pos-project/specs/002-staff-order-billing/spec.md)

**Input**: Feature specification from `/specs/002-staff-order-billing/spec.md` plus T2 task scope for bills, audit logs, usecases, repositories, HTTP endpoints, websocket order updates, staff billing UI, and coverage gates.

**Note**: This template is filled in by the `/speckit-plan` command. See `.specify/templates/plan-template.md` for the execution workflow.

## Summary

Build staff-facing active order and billing management so staff can see incoming orders in real time, check out bills, cancel bills with required reasons and audit logs, split bills without changing the original total, and review completed/cancelled bill history. The backend will add new billing tables through reviewed migrations, keep billing rules in Clean Architecture entity/usecase layers, expose validated staff APIs, and persist split/cancel/checkout changes transactionally. The frontend will add a Pinia-backed responsive staff billing experience with desktop tables, mobile/tablet cards, cancellation and split modals, websocket updates, and history filters.

## Technical Context

**Language/Version**: Frontend: Vue 3 + TypeScript; Backend: Go 1.22

**Primary Dependencies**: Frontend: Pinia, Vitest, Vue Test Utils; Backend: testify; no new dependencies planned

**Storage**: PostgreSQL with new migration files for `bills` and `bill_audit_logs`; migration must be reviewed before apply and must not modify existing Feature 1 schemas

**Testing**: Frontend: Vitest + Vue Test Utils; Backend: Go testing + testify; endpoint integration tests; backend and frontend coverage each at or above 80%

**Target Platform**: Web app for desktop browser, tablet, mobile browser, and native iOS/Android WebView; Go backend on the existing server runtime

**Project Type**: Monorepo with `frontend/`, `backend/`, and optional `packages/`

**Responsive Targets**: mobile >=360px, tablet >=768px, desktop >=1024px, iOS/Android WebView, portrait, and landscape

**Performance Goals**: New orders should appear in staff views without manual refresh during testing; staff can open incoming order details and filter bill history in under 10 seconds; checkout/cancel/split operations complete atomically or return a clear failure

**Constraints**: Billing totals must be deterministic; cancellation requires a non-empty reason; split sub-bill totals must equal the original bill total before commit; split operations must use a DB transaction for sub-bill inserts and original status update; duplicate checkout/cancel/split submissions must not create duplicate financial state; APIs validate input before usecases; no schema edits/deletions beyond new migration files

**Scale/Scope**: Single-store staff order and billing workflow covering active orders, bill checkout, cancellation audit, even/by-item/custom split, bill history filters, and realtime order notifications

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- **Transaction Integrity**: Checkout marks one active bill completed once; cancellation records reason/actor without deleting the bill; split creates sub-bills only when combined totals equal the original total. Rounding and indivisible cents are assigned deterministically.
- **Inventory And Catalog Accuracy**: Staff order details must reflect the order snapshot created by Feature 1, including items, quantities, notes, and table/customer reference.
- **Offline Resilience**: Staff billing actions must not be shown as accepted until confirmed. Network failures, duplicate submissions, latency, and websocket disconnects must leave the current bill/order state clear and refreshable without duplicate checkout, cancellation, or split.
- **Security, Privacy, And Auditability**: Staff identity is required for cancellation and split audit records. Cancellation reason is required. Endpoints expose only staff billing/order data and validate all inputs before usecase calls.
- **Testable Incremental Delivery**: User stories are independently valuable: active order list, checkout/cancel, split, and history. Required tests include entity validation, usecase mocks/fakes, repository integration, HTTP integration, websocket update behavior, responsive UI, and coverage >=80% for frontend/backend.
- **Monorepo Independence**: Frontend and backend build/test independently. Shared behavior is documented through contracts; neither side imports the other's implementation code.
- **Frontend Standards**: Vue 3 Composition API, TypeScript, Pinia, `<script setup>`, Vitest, Vue Test Utils, ESLint, and Prettier are required.
- **Responsive And WebView Support**: Active orders use desktop table layout and mobile/tablet cards. Split modal is full-screen on mobile and dialog on tablet/desktop. Touch targets must be at least 44x44px, with in-app navigation independent of browser chrome.
- **Backend Clean Architecture**: Entities hold bill validation rules; usecases coordinate checkout/cancel/split/list history; interface/adapters define repositories and HTTP/websocket delivery; infrastructure implements PostgreSQL persistence.
- **Database Governance**: New migrations add `bills` and `bill_audit_logs` only. They must be reviewed before apply, must reference Feature 1 `orders`, and must not edit applied migrations or existing schema objects.
- **API Validation And Go Errors**: HTTP handlers validate route/body/query fields and return explicit Go errors; no panics for normal validation or business-rule failure.
- **Dependency Discipline**: No new runtime dependency is planned. If websocket or database migration tooling requires a dependency later, it must be justified before implementation.

No unresolved violations identified.

## Project Structure

### Documentation (this feature)

```text
specs/002-staff-order-billing/
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/
└── tasks.md
```

### Source Code (repository root)

```text
frontend/
├── src/
│   ├── components/
│   │   ├── CancelBillModal.vue
│   │   └── SplitBillModal.vue
│   ├── pages/
│   │   ├── ActiveOrdersView.vue
│   │   └── HistoryView.vue
│   ├── services/
│   └── stores/
│       └── staffBillStore.ts
└── tests/

backend/
├── entity/
│   ├── bill.go
│   └── bill_audit_log.go
├── usecase/
│   ├── bill_repository.go
│   ├── bill_usecase.go
│   └── audit_log_repository.go
├── interface/
│   └── adapter/
│       ├── http/
│       └── websocket/
├── infrastructure/
│   └── postgres/
└── tests/

packages/
└── contracts/
```

**Structure Decision**: Keep the existing monorepo and Clean Architecture boundaries. Backend entity/usecase code owns billing rules and repository interfaces; HTTP/websocket adapters translate external requests/events; PostgreSQL implementations live in infrastructure. Frontend staff billing UI lives under existing Vue pages/components/stores and consumes documented contracts only.

## Complexity Tracking

No Constitution Check violations were identified, so no complexity justification is required.
