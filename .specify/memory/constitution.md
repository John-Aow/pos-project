<!--
Sync Impact Report
Version change: 1.1.0 -> 1.2.0
Modified principles:
- I. Transaction Integrity Is Non-Negotiable -> I. Transaction Integrity Is Non-Negotiable
- II. Inventory And Catalog Accuracy -> II. Inventory And Catalog Accuracy
- III. Offline-Resilient Checkout -> III. Offline-Resilient Checkout
- IV. Security, Privacy, And Auditability -> IV. Security, Privacy, And Auditability
- V. Testable Incremental Delivery -> V. Spec-First, Tested Incremental Delivery
Added principles:
- None
Added sections:
- Responsive & Multi-Device Support
Removed sections:
- None
Templates requiring updates:
- Updated .specify/templates/plan-template.md
- Updated .specify/templates/spec-template.md
- Updated .specify/templates/tasks-template.md
- Updated .specify/templates/checklist-template.md
- Not present .specify/templates/commands/*.md
Follow-up TODOs:
- None
-->
# POS Project Constitution

## Core Principles

### I. Transaction Integrity Is Non-Negotiable
Every checkout, refund, void, discount, tax calculation, payment handoff, and receipt
MUST be recorded as an auditable transaction with deterministic totals. Features that
change money movement MUST define rounding, tax, discount, tender, failure, retry, and
idempotency behavior before implementation begins. Rationale: a POS system is trusted
only when financial records can be reconciled exactly.

### II. Inventory And Catalog Accuracy
Product, price, modifier, tax, stock, and availability changes MUST be modeled as
explicit domain updates with validation and traceability. Checkout flows MUST use a
consistent catalog snapshot for each transaction and MUST not silently sell unknown,
inactive, or incorrectly priced items. Rationale: operators rely on the system for
both customer-facing prices and back-office stock decisions.

### III. Offline-Resilient Checkout
Checkout-critical workflows MUST specify behavior for network loss, service latency,
device restart, duplicate submission, and later synchronization. Local actions that
cannot be completed remotely MUST be queued, visibly marked, and reconciled without
data loss or duplicate charges. Rationale: store operations cannot stop whenever a
connection or upstream service is unreliable.

### IV. Security, Privacy, And Auditability
Features MUST protect customer, staff, payment, and store data with least-privilege
access, secret-safe configuration, and audit logs for sensitive actions. Payment card
data MUST not be stored unless a later, explicit compliance decision defines the
approved scope and controls. Every backend API endpoint MUST validate input before
calling usecase logic. Rationale: POS systems sit at a high-risk boundary between
customer trust, staff access, and financial systems.

### V. Spec-First, Tested Incremental Delivery
No implementation may begin without a specification and acceptance criteria. Each
feature MUST be split into independently valuable user stories with measurable
acceptance criteria. Unit tests MUST cover each feature's acceptance criteria before
the feature is considered done, and frontend and backend code coverage MUST remain at
or above 80%. Rationale: small verified increments reduce operational risk and keep
releases demonstrable to store operators.

### VI. Monorepo Independence
The repository MUST keep frontend and backend applications in one monorepo with
`frontend/`, `backend/`, and optional `packages/` directories. Frontend and backend
MUST build and test independently, and neither side may import implementation code
directly from the other. Shared code or configuration MUST live under `packages/` or
a documented contract such as an API schema, generated type, proto, or script.
Rationale: one repository improves coordination, but independent build and test
boundaries keep releases and ownership clear.

### VII. Clean Backend Architecture
The Go backend MUST follow Clean Architecture with `entity`, `usecase`,
`interface/adapter`, and `infrastructure` layers. Dependencies MUST point inward
toward `entity`; inner layers MUST NOT import outer layers. Usecases MUST depend on
repository or service interfaces and MUST be tested independently from infrastructure
with mocks or fakes. Go errors MUST be explicit return values, and panics MUST NOT be
used for normal control flow. Rationale: POS business rules must remain testable and
independent from delivery frameworks, database drivers, and external services.

## Operational Constraints

The project MUST treat checkout availability, data correctness, and staff efficiency
as first-class constraints. Plans and specifications MUST document performance targets
for cashier-facing actions, accessibility expectations for daily-use screens, and
failure modes for any external dependency involved in checkout or reconciliation.

Implementation choices MUST keep domain behavior understandable in code. Shared
models, schemas, and APIs MUST name business concepts consistently across features.
Code MUST be readable, directly named, and not over-engineered. New dependencies or
libraries MUST NOT be added unless necessary and explicitly justified in the plan.

## Monorepo Technology Standards

Frontend work MUST target Vue 3 with the Composition API, TypeScript, Pinia, and
`<script setup>` for components. Frontend tests MUST use Vitest and Vue Test Utils.
Frontend formatting and static checks MUST follow ESLint and Prettier.

Backend work MUST target Go. Backend formatting and static checks MUST follow
`gofmt` and `golangci-lint`. Backend tests MUST use Go's standard `testing` package
and `testify`. Backend API endpoints MUST keep input validation at the delivery
boundary and pass validated input to usecases.

## Responsive & Multi-Device Support

Frontend work MUST be fully responsive across mobile, tablet, and desktop viewports:
mobile at 360px and wider, tablet at 768px and wider, and desktop at 1024px and
wider. Vue 3 interfaces MUST use a mobile-first CSS approach through Tailwind CSS
breakpoints, CSS Grid/Flexbox, or media queries.

The web app MUST be usable inside native iOS and Android WebViews. Interactive
behavior MUST NOT depend on desktop-only browser features; hover-only interactions
MUST have touch-friendly alternatives, touch targets MUST be at least 44x44px,
layouts MUST avoid fixed pixel widths, and in-app navigation controls MUST exist
where browser chrome may be hidden. Layouts MUST support portrait and landscape
orientation without breaking, and viewport meta configuration MUST be tested to
prevent unwanted zooming or scaling in WebView.

All interactive components MUST be tested at one mobile, one tablet, and one desktop
breakpoint. Rationale: POS workflows may run on phones, tablets, desktop browsers,
and native WebView shells in real store environments.

## Database Governance

Existing database schemas MUST NOT be modified without prior approval through a
specification. Tables, columns, and data MUST NEVER be deleted during implementation,
refactoring, cleanup, or migration work. Schema changes MUST be introduced only by a
new migration file; applied migrations MUST NOT be edited.

If a task appears to require modifying an existing schema object, deleting schema, or
deleting data, implementation MUST stop and human confirmation MUST be requested
before any change is made. Rationale: preserving production data and migration
history takes precedence over implementation convenience.

## Delivery Workflow And Quality Gates

Specifications MUST describe prioritized user journeys, edge cases, measurable
success criteria, and POS-specific risks before planning begins. Implementation plans
MUST pass the Constitution Check before Phase 0 research and again after Phase 1
design. Tasks MUST preserve traceability from each user story to its implementation,
tests, and validation steps.

No feature that touches transactions, catalog data, inventory, offline sync,
permissions, audit logs, payment integration, backend API input handling, Clean
Architecture boundaries, responsive frontend behavior, WebView behavior, or database
migrations may be marked complete until its automated tests, coverage checks,
lint/format checks, viewport checks, and quickstart validation pass. Any exception
MUST be documented in the plan's Complexity Tracking table with risk, owner, and
follow-up mitigation.

## Governance

This constitution supersedes conflicting project practices and templates. Amendments
MUST be proposed as a documented change to this file, include the rationale and
migration impact, and update dependent templates in the same change. Compliance is
reviewed during specification, planning, task generation, implementation review, and
release validation.

Versioning follows semantic versioning. MAJOR increments remove or redefine a
principle in a backward-incompatible way. MINOR increments add a principle, section,
or materially expanded governance requirement. PATCH increments clarify wording
without changing compliance obligations.

**Version**: 1.2.0 | **Ratified**: 2026-07-05 | **Last Amended**: 2026-07-05
