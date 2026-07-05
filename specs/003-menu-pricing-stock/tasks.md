---

description: "Task list template for feature implementation"
---

# Tasks: Manager Menu, Pricing, and Stock Management

**Input**: Design documents from `/specs/003-menu-pricing-stock/`

**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md, contracts/

**Tests**: Tests are required for this feature because the spec includes acceptance criteria for menu visibility, stock behavior, price snapshots, audit logging, and responsive WebView usability.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Path Conventions

- **Frontend**: `frontend/src/`, `frontend/tests/`
- **Backend**: `backend/entity/`, `backend/usecase/`, `backend/interface/adapter/`, `backend/infrastructure/`, `backend/tests/`
- **Shared**: `packages/contracts/` for shared API types or schemas if needed
- Frontend and backend tasks must build/test independently. Do not create direct implementation imports between `frontend/` and `backend/`.

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

**Purpose**: Project initialization and shared scaffolding

- [X] T001 Create the feature-ready project structure for frontend, backend, and shared contracts in `frontend/`, `backend/`, and `packages/contracts/`
- [X] T002 Configure the frontend toolchain, test setup, and responsive viewport test helpers in `frontend/vitest.config.ts`, `frontend/src/main.ts`, and `frontend/tests/setup.ts`
- [X] T003 Configure the backend toolchain, test setup, and migration workflow in `backend/go.mod`, `backend/tests/setup_test.go`, and `backend/infrastructure/postgres/migrations/`

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

Examples of foundational tasks (adjust based on your project):

- [X] T004 [P] Create the menu catalog migration for `categories`, `menu_items`, `inventory_settings`, and `menu_item_audit_entries` in `backend/infrastructure/postgres/migrations/20260705_create_menu_catalog.sql`
- [X] T005 [P] Implement core entity structs and validation rules in `backend/entity/category.go`, `backend/entity/menu_item.go`, `backend/entity/inventory_settings.go`, and `backend/entity/menu_item_audit_entry.go`
- [X] T006 [P] Define repository, inventory, and audit interfaces in `backend/usecase/menu_repository.go`, `backend/usecase/inventory_repository.go`, and `backend/usecase/audit_repository.go`
- [X] T007 [P] Add HTTP routing and request validation helpers for manager APIs in `backend/interface/adapter/http/router.go` and `backend/interface/adapter/http/validation.go`
- [X] T008 [P] Add the manager shell layout and responsive app scaffolding in `frontend/src/pages/ManagerMenuPage.vue`, `frontend/src/components/ManagerShell.vue`, and `frontend/src/styles/manager.css`

**Checkpoint**: Foundation ready - user story implementation can now begin in parallel

---

## Phase 3: User Story 1 - Manage Menu Items (Priority: P1)

**Goal**: Managers can create, edit, and deactivate menu items so the customer menu stays accurate.

**Independent Test**: Create a new item, edit it, deactivate it, and confirm the customer-facing menu reflects the active state without deleting history.

### Tests for User Story 1

> **NOTE: Write required tests first and ensure they fail before implementation**

- [ ] T009 [P] [US1] Add backend unit tests for menu item create/edit/deactivate use cases in `backend/tests/menu_item_usecase_test.go`
- [ ] T010 [P] [US1] Add frontend component tests for menu item creation, editing, and deactivation in `frontend/tests/manager-menu-page.spec.ts`
- [ ] T011 [P] [US1] Add responsive breakpoint tests for the manager menu layout in `frontend/tests/manager-menu-responsive.spec.ts`

### Implementation for User Story 1

- [ ] T012 [P] [US1] Implement menu item create/edit/deactivate use case methods in `backend/usecase/menu_usecase.go`
- [ ] T013 [US1] Implement menu item HTTP handlers and presenters in `backend/interface/adapter/http/menu_handler.go` and `backend/interface/adapter/http/menu_presenter.go`
- [ ] T014 [P] [US1] Implement menu repository persistence for create/read/update/deactivate in `backend/infrastructure/postgres/menu_repository.go`
- [ ] T015 [US1] Implement the manager menu page, form, table, and active-state UI in `frontend/src/pages/ManagerMenuPage.vue`, `frontend/src/components/MenuItemForm.vue`, `frontend/src/components/MenuItemTable.vue`, and `frontend/src/stores/managerMenuStore.ts`

**Checkpoint**: User Story 1 is complete when menu item maintenance is testable independently.

---

## Phase 4: User Story 2 - Manage Price And Stock (Priority: P1)

**Goal**: Managers can update prices and stock levels while preserving historical bill prices.

**Independent Test**: Update a price, update stock, set stock to zero, and confirm the customer menu changes immediately while historical bills keep the old price snapshot.

### Tests for User Story 2

- [ ] T016 [P] [US2] Add backend unit tests for price updates, stock updates, zero-stock availability, and historical price snapshot behavior in `backend/tests/price_stock_usecase_test.go`
- [ ] T017 [P] [US2] Add frontend component tests for price and stock editing flows in `frontend/tests/price-stock-editor.spec.ts`

### Implementation for User Story 2

- [ ] T018 [P] [US2] Implement price and stock use case methods, including availability recalculation and historical price snapshot rules, in `backend/usecase/menu_usecase.go`
- [ ] T019 [US2] Implement price and stock HTTP handlers plus audit logging in `backend/interface/adapter/http/menu_handler.go` and `backend/interface/adapter/http/menu_presenter.go`
- [ ] T020 [P] [US2] Implement price and stock persistence updates with transactional safety in `backend/infrastructure/postgres/menu_repository.go`
- [ ] T021 [US2] Implement the manager price/stock editor UI and availability badges in `frontend/src/pages/ManagerMenuPage.vue`, `frontend/src/components/PriceEditor.vue`, `frontend/src/components/StockEditor.vue`, and `frontend/src/stores/managerMenuStore.ts`

**Checkpoint**: User Story 2 is complete when new orders use the new price, zero stock makes items unavailable, and old bills keep their original price snapshot.

---

## Phase 5: User Story 3 - Review Stock Warnings (Priority: P2)

**Goal**: Managers can identify low-stock items and tune the warning threshold before items run out.

**Independent Test**: Set items below the threshold, verify they appear in the low-stock view, then adjust the threshold and confirm the list updates.

### Tests for User Story 3

- [ ] T022 [P] [US3] Add backend unit tests for low-stock queries and threshold updates in `backend/tests/low_stock_usecase_test.go`
- [ ] T023 [P] [US3] Add frontend component tests for the low-stock warning view and threshold control in `frontend/tests/low-stock-warning.spec.ts`

### Implementation for User Story 3

- [ ] T024 [P] [US3] Implement low-stock query and threshold update use case methods in `backend/usecase/inventory_usecase.go`
- [ ] T025 [US3] Implement low-stock HTTP endpoints and settings handlers in `backend/interface/adapter/http/settings_handler.go` and `backend/interface/adapter/http/menu_handler.go`
- [ ] T026 [P] [US3] Implement threshold persistence and low-stock query support in `backend/infrastructure/postgres/inventory_repository.go`
- [ ] T027 [US3] Implement the low-stock warning view and threshold controls in `frontend/src/pages/LowStockView.vue`, `frontend/src/components/LowStockList.vue`, `frontend/src/components/ThresholdControl.vue`, and `frontend/src/stores/inventoryStore.ts`

**Checkpoint**: User Story 3 is complete when the low-stock list is accurate and configurable.

---

## Phase N: Polish & Cross-Cutting Concerns

**Purpose**: Improvements that affect multiple user stories

- [ ] T028 [P] Run backend formatter, linter, tests, and coverage in `backend/`
- [ ] T029 [P] Run frontend formatter, linter, tests, coverage, and viewport checks in `frontend/`
- [ ] T030 Update the quickstart validation guide with final verification notes in `specs/003-menu-pricing-stock/quickstart.md`
- [ ] T031 Update the API contract documentation if endpoint shapes changed in `specs/003-menu-pricing-stock/contracts/menu-management-api.md`

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion - BLOCKS all user stories
- **User Stories (Phase 3+)**: All depend on Foundational phase completion
  - User stories can then proceed in parallel if staffed
  - Or sequentially in priority order (P1 → P2)
- **Polish (Final Phase)**: Depends on all desired user stories being complete

### User Story Dependencies

- **User Story 1 (P1)**: Can start after Foundational (Phase 2) - No dependencies on other stories
- **User Story 2 (P1)**: Can start after Foundational (Phase 2) - Shares backend menu items and repository work with US1 but remains independently testable
- **User Story 3 (P2)**: Can start after Foundational (Phase 2) - Depends on low-stock threshold support and menu/stock data structures

### Within Each User Story

- Tests MUST be written before implementation for the story tasks above
- Backend entities and interfaces before usecases
- Usecases before HTTP adapters and repositories
- Repositories before integration wiring
- Frontend stores/composables before dependent components when shared state is needed
- Responsive layout, touch alternatives, and WebView-safe navigation before story completion
- Story complete before moving to next priority

### Parallel Opportunities

- All Setup tasks marked [P] can run in parallel
- All Foundational tasks marked [P] can run in parallel
- Backend test tasks and frontend test tasks within a story can run in parallel
- Frontend and backend tasks for a story can run in parallel when they touch separate files
- Story 1 and Story 2 can proceed in parallel after foundational work if staffing allows

---

## Parallel Example: User Story 1

```bash
# Launch tests for User Story 1 together:
Task: "Add backend unit tests for menu item create/edit/deactivate use cases in backend/tests/menu_item_usecase_test.go"
Task: "Add frontend component tests for menu item creation, editing, and deactivation in frontend/tests/manager-menu-page.spec.ts"
Task: "Add responsive breakpoint tests for the manager menu layout in frontend/tests/manager-menu-responsive.spec.ts"

# Launch independent implementation work together:
Task: "Implement menu item create/edit/deactivate use case methods in backend/usecase/menu_usecase.go"
Task: "Implement the manager menu page, form, table, and active-state UI in frontend/src/pages/ManagerMenuPage.vue"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup
2. Complete Phase 2: Foundational
3. Complete Phase 3: User Story 1
4. **STOP and VALIDATE**: Verify menu item maintenance independently
5. Deploy/demo if ready

### Incremental Delivery

1. Complete Setup + Foundational → foundation ready
2. Add User Story 1 → menu item maintenance works
3. Add User Story 2 → price and stock controls work without changing historical bills
4. Add User Story 3 → low-stock warnings work
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
- Never edit applied migrations or delete database schema/data
- Stop and ask for human confirmation if a task requires schema modification or data deletion
- Commit after each task or logical group
- Stop at any checkpoint to validate story independently
- Avoid vague tasks, same file conflicts, and cross-story dependencies that break independence
