---

description: "Task list template for feature implementation"
---

# Tasks: [FEATURE NAME]

**Input**: Design documents from `/specs/[###-feature-name]/`

**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md, contracts/

**Tests**: Unit tests are REQUIRED for every feature's acceptance criteria. Tests are
also REQUIRED for transaction, inventory, catalog, offline, permission, audit, payment
integration, API validation, Clean Architecture usecases, database migration, and
reconciliation behavior. Frontend interactive components MUST be tested at one mobile
breakpoint, one tablet breakpoint, and one desktop breakpoint. Frontend and backend
coverage MUST remain at or above 80%.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Path Conventions

- **Frontend**: `frontend/src/`, `frontend/tests/`
- **Backend**: `backend/entity/`, `backend/usecase/`, `backend/interface/adapter/`,
  `backend/infrastructure/`, `backend/tests/`
- **Shared**: `packages/` for contracts, generated types, proto, scripts, or shared
  config only
- Frontend and backend tasks must build/test independently. Do not create direct
  implementation imports between `frontend/` and `backend/`.

<!--
  ============================================================================
  IMPORTANT: The tasks below are SAMPLE TASKS for illustration purposes only.

  The /speckit-tasks command MUST replace these with actual tasks based on:
  - User stories from spec.md (with their priorities P1, P2, P3...)
  - Feature requirements from plan.md
  - Entities from data-model.md
  - Endpoints from contracts/

  Tasks MUST be organized by user story so each story can be:
  - Implemented independently
  - Tested independently
  - Delivered as an MVP increment

  DO NOT keep these sample tasks in the generated tasks.md file.
  ============================================================================
-->

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and basic structure

- [ ] T001 Create project structure per implementation plan
- [ ] T002 Initialize frontend Vue 3 + TypeScript project in frontend/
- [ ] T003 Initialize backend Go project in backend/
- [ ] T004 [P] Configure frontend ESLint, Prettier, Vitest, and Vue Test Utils
- [ ] T005 [P] Configure backend gofmt, golangci-lint, testing, and testify
- [ ] T006 [P] Configure frontend responsive/WebView test targets for mobile ≥360px, tablet ≥768px, desktop ≥1024px, portrait, and landscape

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

Examples of foundational tasks (adjust based on your project):

- [ ] T007 Setup database schema and migrations framework using new migrations only
- [ ] T008 [P] Implement authentication/authorization framework
- [ ] T009 [P] Setup backend API routing and middleware structure with input validation
- [ ] T010 Create backend entity models that all stories depend on
- [ ] T011 Define backend usecase interfaces and repository contracts
- [ ] T012 Configure explicit Go error handling and logging infrastructure
- [ ] T013 Setup environment configuration management
- [ ] T014 Define audit event capture for sensitive POS operations
- [ ] T015 Define offline queue/reconciliation foundation if checkout-critical flows are affected

**Checkpoint**: Foundation ready - user story implementation can now begin in parallel

---

## Phase 3: User Story 1 - [Title] (Priority: P1) 🎯 MVP

**Goal**: [Brief description of what this story delivers]

**Independent Test**: [How to verify this story works on its own]

### Tests for User Story 1 (REQUIRED when constitution risks apply) ⚠️

> **NOTE: Write required constitution tests FIRST and ensure they FAIL before implementation**

- [ ] T016 [P] [US1] Frontend unit/component test for [acceptance criterion] in frontend/tests/[name].test.ts
- [ ] T017 [P] [US1] Responsive/WebView test for mobile, tablet, desktop, touch target, orientation, and viewport meta behavior in [path]
- [ ] T018 [P] [US1] Backend usecase unit test with mocked repositories in backend/tests/[name]_test.go
- [ ] T019 [P] [US1] Contract test for [endpoint/API contract] in [path]
- [ ] T020 [P] [US1] POS risk test for [transaction/inventory/offline/permission/audit/database behavior] in [path]

### Implementation for User Story 1

- [ ] T021 [P] [US1] Create backend entity in backend/entity/[name].go
- [ ] T022 [US1] Implement backend usecase in backend/usecase/[name].go
- [ ] T023 [US1] Implement backend adapter/controller in backend/interface/adapter/[name].go
- [ ] T024 [US1] Implement backend infrastructure integration in backend/infrastructure/[name].go
- [ ] T025 [US1] Implement frontend store/composable in frontend/src/stores/[name].ts or frontend/src/composables/[name].ts
- [ ] T026 [US1] Implement responsive Vue component/page with `<script setup>` in frontend/src/[location]/[name].vue
- [ ] T027 [US1] Add touch-friendly alternatives, 44x44px touch targets, fluid layout, in-app navigation, and WebView-safe viewport behavior
- [ ] T028 [US1] Add API input validation, explicit Go error returns, and audit logging

**Checkpoint**: At this point, User Story 1 should be fully functional and testable independently

---

## Phase 4: User Story 2 - [Title] (Priority: P2)

**Goal**: [Brief description of what this story delivers]

**Independent Test**: [How to verify this story works on its own]

### Tests for User Story 2 (REQUIRED when constitution risks apply) ⚠️

- [ ] T029 [P] [US2] Frontend unit/component test for [acceptance criterion] in frontend/tests/[name].test.ts
- [ ] T030 [P] [US2] Responsive/WebView test for mobile, tablet, desktop, touch target, orientation, and viewport meta behavior in [path]
- [ ] T031 [P] [US2] Backend usecase unit test with mocked repositories in backend/tests/[name]_test.go
- [ ] T032 [P] [US2] POS risk test for [transaction/inventory/offline/permission/audit/database behavior] in [path]

### Implementation for User Story 2

- [ ] T033 [P] [US2] Update backend entity/usecase files in backend/
- [ ] T034 [US2] Implement backend adapter/infrastructure files in backend/
- [ ] T035 [US2] Implement responsive frontend Pinia/component changes in frontend/
- [ ] T036 [US2] Integrate with User Story 1 components through documented API contracts only

**Checkpoint**: At this point, User Stories 1 AND 2 should both work independently

---

## Phase 5: User Story 3 - [Title] (Priority: P3)

**Goal**: [Brief description of what this story delivers]

**Independent Test**: [How to verify this story works on its own]

### Tests for User Story 3 (REQUIRED when constitution risks apply) ⚠️

- [ ] T037 [P] [US3] Frontend unit/component test for [acceptance criterion] in frontend/tests/[name].test.ts
- [ ] T038 [P] [US3] Responsive/WebView test for mobile, tablet, desktop, touch target, orientation, and viewport meta behavior in [path]
- [ ] T039 [P] [US3] Backend usecase unit test with mocked repositories in backend/tests/[name]_test.go
- [ ] T040 [P] [US3] POS risk test for [transaction/inventory/offline/permission/audit/database behavior] in [path]

### Implementation for User Story 3

- [ ] T041 [P] [US3] Update backend entity/usecase files in backend/
- [ ] T042 [US3] Implement backend adapter/infrastructure files in backend/
- [ ] T043 [US3] Implement responsive frontend Pinia/component changes in frontend/

**Checkpoint**: All user stories should now be independently functional

---

[Add more user story phases as needed, following the same pattern]

---

## Phase N: Polish & Cross-Cutting Concerns

**Purpose**: Improvements that affect multiple user stories

- [ ] TXXX [P] Documentation updates in docs/
- [ ] TXXX Code cleanup and refactoring
- [ ] TXXX Performance optimization across all stories
- [ ] TXXX [P] Additional frontend unit/component tests in frontend/tests/
- [ ] TXXX [P] Additional backend usecase/unit tests in backend/tests/
- [ ] TXXX Security hardening
- [ ] TXXX Audit log review for sensitive POS operations
- [ ] TXXX Offline/reconciliation validation for checkout-critical flows
- [ ] TXXX Responsive review for mobile ≥360px, tablet ≥768px, desktop ≥1024px, portrait, and landscape
- [ ] TXXX WebView review for touch targets, no hover-only dependency, in-app navigation, and viewport meta behavior
- [ ] TXXX Run frontend lint, format check, tests, and coverage
- [ ] TXXX Run backend gofmt, golangci-lint, tests, and coverage
- [ ] TXXX Run quickstart.md validation

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion - BLOCKS all user stories
- **User Stories (Phase 3+)**: All depend on Foundational phase completion
  - User stories can then proceed in parallel (if staffed)
  - Or sequentially in priority order (P1 → P2 → P3)
- **Polish (Final Phase)**: Depends on all desired user stories being complete

### User Story Dependencies

- **User Story 1 (P1)**: Can start after Foundational (Phase 2) - No dependencies on other stories
- **User Story 2 (P2)**: Can start after Foundational (Phase 2) - May integrate with US1 but should be independently testable
- **User Story 3 (P3)**: Can start after Foundational (Phase 2) - May integrate with US1/US2 but should be independently testable

### Within Each User Story

- Constitution-required tests MUST be written and FAIL before implementation
- Backend entities before usecases
- Backend usecases before adapters/controllers
- Backend adapters before infrastructure integration
- Frontend stores/composables before dependent components when shared state is needed
- Responsive layout, touch alternatives, and in-app navigation before UI story completion
- API input validation before calling usecases
- Story complete before moving to next priority

### Parallel Opportunities

- All Setup tasks marked [P] can run in parallel
- All Foundational tasks marked [P] can run in parallel (within Phase 2)
- Once Foundational phase completes, all user stories can start in parallel (if team capacity allows)
- All tests for a user story marked [P] can run in parallel
- Frontend and backend tasks for a story can run in parallel when they touch separate
  files and communicate through documented contracts
- Different user stories can be worked on in parallel by different team members

---

## Parallel Example: User Story 1

```bash
# Launch all tests for User Story 1 together:
Task: "Frontend unit/component test for [acceptance criterion] in frontend/tests/[name].test.ts"
Task: "Responsive/WebView test for mobile, tablet, desktop, touch target, orientation, and viewport meta behavior in [path]"
Task: "Backend usecase unit test with mocked repositories in backend/tests/[name]_test.go"
Task: "Contract test for [endpoint/API contract] in [path]"

# Launch independent frontend/backend implementation together:
Task: "Implement backend usecase in backend/usecase/[name].go"
Task: "Implement frontend store/composable in frontend/src/stores/[name].ts"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup
2. Complete Phase 2: Foundational (CRITICAL - blocks all stories)
3. Complete Phase 3: User Story 1
4. **STOP and VALIDATE**: Test User Story 1 independently
5. Deploy/demo if ready

### Incremental Delivery

1. Complete Setup + Foundational → Foundation ready
2. Add User Story 1 → Test independently → Deploy/Demo (MVP!)
3. Add User Story 2 → Test independently → Deploy/Demo
4. Add User Story 3 → Test independently → Deploy/Demo
5. Each story adds value without breaking previous stories

### Parallel Team Strategy

With multiple developers:

1. Team completes Setup + Foundational together
2. Once Foundational is done:
   - Developer A: User Story 1
   - Developer B: User Story 2
   - Developer C: User Story 3
3. Stories complete and integrate independently

---

## Notes

- [P] tasks = different files, no dependencies
- [Story] label maps task to specific user story for traceability
- Each user story should be independently completable and testable
- Verify tests fail before implementing
- Maintain 80% minimum coverage for frontend and backend
- Test interactive frontend components at mobile, tablet, and desktop breakpoints
- Preserve WebView usability without relying on browser chrome or hover-only interactions
- Never edit applied migrations or delete database schema/data
- Stop and ask for human confirmation if a task requires schema modification or data deletion
- Commit after each task or logical group
- Stop at any checkpoint to validate story independently
- Avoid: vague tasks, same file conflicts, cross-story dependencies that break independence
