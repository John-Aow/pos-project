# Implementation Plan: Manager Menu, Pricing, and Stock Management

**Branch**: `003-menu-pricing-stock` | **Date**: 2026-07-05 | **Spec**: [/Users/johnhu/Desktop/pos-project/specs/003-menu-pricing-stock/spec.md](/Users/johnhu/Desktop/pos-project/specs/003-menu-pricing-stock/spec.md)

**Input**: Feature specification from `/specs/003-menu-pricing-stock/spec.md`

**Note**: This template is filled in by the `/speckit-plan` command. See `.specify/templates/plan-template.md` for the execution workflow.

## Summary

Build manager-facing menu, pricing, and stock controls that keep the customer menu accurate in real time, make out-of-stock items unavailable, and preserve order/bill price history so old bills never change when prices are updated later.

## Technical Context

**Language/Version**: Frontend: Vue 3 + TypeScript; Backend: Go

**Primary Dependencies**: Frontend: Pinia, Vitest, Vue Test Utils; Backend: testify; new dependencies require justification

**Storage**: PostgreSQL

**Testing**: Frontend: Vitest + Vue Test Utils; Backend: Go testing + testify; 80% minimum coverage

**Target Platform**: Web app for desktop browser, tablet, mobile browser, and native iOS/Android WebView

**Project Type**: Monorepo with `frontend/`, `backend/`, and optional `packages/`

**Responsive Targets**: mobile ≥360px, tablet ≥768px, desktop ≥1024px, iOS/Android WebView

**Performance Goals**: Customer-visible menu availability should reflect confirmed stock changes within 5 seconds, and stock-sensitive updates should remain safe under concurrent order activity

**Constraints**: Price changes must not retroactively alter historical bills; stock updates must remain transactionally consistent with concurrent order activity; no existing schema/data deletion; frontend and backend must build/test independently; 80% minimum coverage; WebView-safe responsive UI

**Scale/Scope**: Single-store catalog management with categories, item pricing, stock levels, and low-stock warnings for a moderate menu catalog

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- **Transaction Integrity**: Price changes must apply only to future orders; historical bills must keep their original unit prices and totals.
- **Inventory And Catalog Accuracy**: Stock changes must drive customer-facing availability, and zero stock must make items unavailable automatically.
- **Offline Resilience**: Manager updates must not appear accepted until the system confirms them successfully.
- **Security, Privacy, And Auditability**: Manager menu, price, stock, deactivate, and threshold actions must be traceable.
- **Testable Incremental Delivery**: User stories are independently valuable, with coverage at or above 80% for both frontend and backend.
- **Monorepo Independence**: Frontend and backend must build/test independently, with shared contracts only through documented interfaces.
- **Frontend Standards**: Vue 3 Composition API, TypeScript, Pinia, `<script setup>`, Vitest, Vue Test Utils, ESLint, and Prettier are required.
- **Responsive And WebView Support**: Mobile-first responsive layouts must support mobile, tablet, desktop, portrait, landscape, touch alternatives, 44x44px targets, and WebView-safe viewport behavior.
- **Backend Clean Architecture**: Changes must map cleanly to `entity`, `usecase`, `interface/adapter`, and `infrastructure`.
- **Database Governance**: No existing schema or data may be deleted or edited in place; schema changes must use new migration files only.
- **API Validation And Go Errors**: Backend endpoints must validate input and return explicit Go errors.
- **Dependency Discipline**: Any new dependency must be justified in the plan.

No unresolved violations identified.

## Project Structure

### Documentation (this feature)

```text
specs/003-menu-pricing-stock/
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
│   ├── composables/
│   ├── pages/
│   ├── services/
│   └── stores/
└── tests/

backend/
├── entity/
├── usecase/
├── interface/
│   └── adapter/
│       └── http/
├── infrastructure/
│   ├── postgres/
│   └── audit/
└── tests/

packages/
└── contracts/
```

**Structure Decision**: Use a single monorepo with separate frontend and backend roots. Shared API contracts, if needed, live under `packages/contracts/` so neither side imports the other’s implementation directly.

## Complexity Tracking

No Constitution Check violations were identified, so no complexity justification is required.
