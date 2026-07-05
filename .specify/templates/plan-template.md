# Implementation Plan: [FEATURE]

**Branch**: `[###-feature-name]` | **Date**: [DATE] | **Spec**: [link]

**Input**: Feature specification from `/specs/[###-feature-name]/spec.md`

**Note**: This template is filled in by the `/speckit-plan` command. See `.specify/templates/plan-template.md` for the execution workflow.

## Summary

[Extract from feature spec: primary requirement + technical approach from research]

## Technical Context

<!--
  ACTION REQUIRED: Replace the content in this section with the technical details
  for the project. The structure here is presented in advisory capacity to guide
  the iteration process.
-->

**Language/Version**: [Frontend: Vue 3 + TypeScript; Backend: Go; or NEEDS CLARIFICATION]

**Primary Dependencies**: [Frontend: Pinia, Vitest, Vue Test Utils; Backend: testify; new dependencies require justification]

**Storage**: [if applicable, e.g., PostgreSQL, CoreData, files or N/A]

**Testing**: [Frontend: Vitest + Vue Test Utils; Backend: Go testing + testify; 80% minimum coverage]

**Target Platform**: [e.g., Linux server, iOS 15+, WASM or NEEDS CLARIFICATION]

**Project Type**: Monorepo with `frontend/`, `backend/`, and optional `packages/`

**Responsive Targets**: [mobile ≥360px, tablet ≥768px, desktop ≥1024px, iOS/Android WebView, or N/A]

**Performance Goals**: [domain-specific, e.g., 1000 req/s, 10k lines/sec, 60 fps or NEEDS CLARIFICATION]

**Constraints**: [domain-specific, e.g., <200ms p95, <100MB memory, offline-capable or NEEDS CLARIFICATION]

**Scale/Scope**: [domain-specific, e.g., 10k users, 1M LOC, 50 screens or NEEDS CLARIFICATION]

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- **Transaction Integrity**: If money movement is affected, document deterministic
  totals, rounding, tax, discount, tender, failure, retry, and idempotency behavior.
- **Inventory And Catalog Accuracy**: If product, price, stock, modifier, tax, or
  availability data is affected, document validation, traceability, and catalog
  snapshot behavior.
- **Offline Resilience**: If checkout-critical work depends on networked services,
  document network loss, latency, restart, duplicate submission, queueing, and sync
  reconciliation behavior.
- **Security, Privacy, And Auditability**: Document staff/customer/payment data
  exposure, permission boundaries, secret handling, and audit events.
- **Testable Incremental Delivery**: Confirm each user story is independently
  valuable and list required automated tests for transaction, inventory, offline,
  permission, integration behavior, and acceptance criteria. Confirm coverage remains
  at or above 80% for both frontend and backend.
- **Monorepo Independence**: Identify affected side(s). Confirm frontend and backend
  can build/test independently and no direct cross-side implementation imports are
  introduced.
- **Frontend Standards**: For frontend work, confirm Vue 3 Composition API,
  TypeScript, Pinia, `<script setup>`, Vitest, Vue Test Utils, ESLint, and Prettier
  are used.
- **Responsive And WebView Support**: For frontend work, confirm mobile-first layouts
  support mobile ≥360px, tablet ≥768px, desktop ≥1024px, portrait, landscape, and
  iOS/Android WebView. Confirm touch alternatives for hover-only behavior, 44x44px
  minimum touch targets, no fixed-width layout dependency, in-app navigation that
  does not rely on browser chrome, and viewport meta validation.
- **Backend Clean Architecture**: For backend work, map changes to `entity`,
  `usecase`, `interface/adapter`, and `infrastructure`. Confirm dependencies point
  inward and usecases are testable without infrastructure.
- **Database Governance**: If database work is involved, confirm no existing schema
  or data is modified/deleted without prior spec approval. Confirm schema changes use
  new migration files only.
- **API Validation And Go Errors**: For backend API work, confirm endpoint input
  validation and explicit Go error returns with no panics for normal control flow.
- **Dependency Discipline**: List any new dependency and explain why it is necessary.

Document any unresolved violation in Complexity Tracking before continuing.

## Project Structure

### Documentation (this feature)

```text
specs/[###-feature]/
├── plan.md              # This file (/speckit-plan command output)
├── research.md          # Phase 0 output (/speckit-plan command)
├── data-model.md        # Phase 1 output (/speckit-plan command)
├── quickstart.md        # Phase 1 output (/speckit-plan command)
├── contracts/           # Phase 1 output (/speckit-plan command)
└── tasks.md             # Phase 2 output (/speckit-tasks command - NOT created by /speckit-plan)
```

### Source Code (repository root)
<!--
  ACTION REQUIRED: Replace the placeholder tree below with the concrete layout
  for this feature. Delete unused options and expand the chosen structure with
  real paths (e.g., apps/admin, packages/something). The delivered plan must
  not include Option labels.
-->

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
├── infrastructure/
└── tests/

packages/
└── [shared contracts, generated types, proto, scripts, or shared config only]
```

**Structure Decision**: [Document the selected structure and reference the real
directories captured above]

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| [e.g., 4th project] | [current need] | [why 3 projects insufficient] |
| [e.g., Repository pattern] | [specific problem] | [why direct DB access insufficient] |
