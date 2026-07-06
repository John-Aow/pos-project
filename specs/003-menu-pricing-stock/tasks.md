---

description: "Task list for Manager Menu, Pricing, and Stock Management"
---

# Tasks: Manager Menu, Pricing, and Stock Management

**Input**: Design documents from `/specs/003-menu-pricing-stock/`

**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/, quickstart.md

**Tests**: Required by the constitution and by the T3 acceptance criteria. Backend and frontend coverage must remain at or above 80%. Tests must cover soft-delete behavior, price history creation, old `OrderItem.unit_price` snapshots, stock zero auto-disable, low-stock warnings, API status/schema, and responsive frontend layouts.

**Organization**: Tasks are grouped by user story so each story can be implemented and tested independently after the foundational migrations/entities/interfaces are complete.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel because it touches different files and does not depend on incomplete tasks
- **[Story]**: Maps to user stories from `specs/003-menu-pricing-stock/spec.md`
- Every task includes exact file paths

## Path Conventions

- **Frontend**: `frontend/src/`, `frontend/tests/`
- **Backend**: `backend/entity/`, `backend/usecase/`, `backend/interface/adapter/`, `backend/infrastructure/`, `backend/tests/`
- **Migrations**: `backend/infrastructure/postgres/migrations/`
- **Contracts**: `specs/003-menu-pricing-stock/contracts/`

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Align the feature workspace, contracts, and validation targets before implementation.

- [X] T001 Update the manager menu API contract for `low_stock_threshold`, `menu_item_price_history`, and `menu_audit_logs` behavior in `specs/003-menu-pricing-stock/contracts/menu-management-api.md`
- [X] T002 Update the quickstart validation guide with reviewed migration, price-history, soft-delete, and coverage checks in `specs/003-menu-pricing-stock/quickstart.md`
- [X] T003 [P] Verify frontend test helpers include mobile, tablet, and desktop viewport support in `frontend/tests/setup.ts`
- [X] T004 [P] Verify backend test setup can run migration/repository integration tests in `backend/tests/setup_test.go`

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Database, entities, and interfaces that all user stories need.

**CRITICAL**: No user story work can begin until this phase is complete.

- [X] T005 Create reviewed additive migration adding `low_stock_threshold` to `menu_items` in `backend/infrastructure/postgres/migrations/20260705_add_low_stock_threshold_to_menu_items.sql`
- [X] T006 Create rollback migration for `low_stock_threshold` column in `backend/infrastructure/postgres/migrations/20260705_add_low_stock_threshold_to_menu_items.down.sql`
- [X] T007 Create reviewed additive migration for `menu_item_price_history` and `menu_audit_logs` with FKs to `menu_items` in `backend/infrastructure/postgres/migrations/20260705_create_menu_price_history_and_audit_logs.sql`
- [X] T008 Create rollback migration for `menu_item_price_history` and `menu_audit_logs` in `backend/infrastructure/postgres/migrations/20260705_create_menu_price_history_and_audit_logs.down.sql`
- [X] T009 [P] Add migration verification tests for additive schema, preserved menu data, and `menu_items` foreign keys in `backend/infrastructure/postgres/menu_repository_test.go`
- [X] T010 [P] Extend `MenuItem` validation for price > 0, stock >= 0, active/available rules, and per-item low-stock threshold in `backend/entity/menu_item.go`
- [X] T011 [P] Add `PriceHistory` entity validation in `backend/entity/price_history.go`
- [X] T012 [P] Add `MenuAuditLog` entity validation in `backend/entity/menu_audit_log.go`
- [X] T013 [P] Add entity unit tests for menu item, price history, and audit log validation in `backend/entity/entity_test.go`
- [X] T014 Define `MenuManagementRepository`, `PriceHistoryRepository`, and `AuditLogRepository` interfaces in `backend/usecase/menu_repository.go`
- [X] T015 [P] Add HTTP request validation helpers for manager menu, price, stock, and low-stock inputs in `backend/interface/adapter/http/validation.go`

**Checkpoint**: Database migrations, entities, and repository interfaces are ready for story implementation.

---

## Phase 3: User Story 1 - Manage Menu Items (Priority: P1) MVP

**Goal**: Managers can create, edit, and deactivate menu items without deleting historical records.

**Independent Test**: Create a menu item, update it, deactivate it, and verify the record remains persisted while `is_available=false` removes it from customer ordering.

### Tests for User Story 1

- [X] T016 [P] [US1] Add usecase tests for `CreateMenuItem`, `UpdateMenuItem`, and `DeactivateMenuItem` soft-delete behavior in `backend/tests/menu_usecase_test.go`
- [X] T017 [P] [US1] Add repository integration test proving deactivate keeps the record in `backend/infrastructure/postgres/menu_repository_test.go`
- [X] T018 [P] [US1] Add HTTP integration tests for menu item CRUD status/schema and validation failures in `backend/interface/adapter/http/menu_handler_test.go`
- [X] T019 [P] [US1] Add frontend store and page tests for create/edit/deactivate flows in `frontend/tests/manager-menu-page.spec.ts`
- [X] T020 [P] [US1] Add responsive tests for desktop data table and mobile accordion-card layout in `frontend/tests/manager-menu-responsive.spec.ts`

### Implementation for User Story 1

- [X] T021 [US1] Implement `CreateMenuItem`, `UpdateMenuItem`, and `DeactivateMenuItem` usecases with soft-delete semantics in `backend/usecase/menu_usecase.go`
- [X] T022 [US1] Implement PostgreSQL create/update/deactivate persistence without physical deletes in `backend/infrastructure/postgres/menu_repository.go`
- [X] T023 [US1] Implement menu item CRUD HTTP handlers with boundary validation in `backend/interface/adapter/http/menu_handler.go`
- [X] T024 [P] [US1] Implement menu response mapping for active and manager-visible records in `backend/interface/adapter/http/menu_presenter.go`
- [X] T025 [US1] Implement `managerMenuStore.ts` actions for listing, creating, updating, and deactivating menu items in `frontend/src/stores/managerMenuStore.ts`
- [X] T026 [US1] Implement manager menu desktop table and mobile accordion-card layout in `frontend/src/pages/MenuManagementView.vue`
- [X] T027 [US1] Implement create/edit form with price and stock validation, full-screen mobile form, and tablet/desktop modal behavior in `frontend/src/components/MenuItemForm.vue`

**Checkpoint**: User Story 1 is independently shippable when menu item maintenance works and deactivate never deletes records.

---

## Phase 4: User Story 2 - Manage Price And Stock (Priority: P1)

**Goal**: Managers can update item prices and stock while preserving historical order item price snapshots and updating customer availability immediately.

**Independent Test**: Update a price, confirm old `OrderItem.unit_price` remains unchanged, set stock to zero, and confirm the customer menu marks the item unavailable.

### Tests for User Story 2

- [X] T028 [P] [US2] Add usecase tests for `UpdatePrice`, `UpdateStock`, price history creation, and stock=0 auto-disable in `backend/tests/price_stock_usecase_test.go`
- [X] T029 [P] [US2] Add integration test proving old `OrderItem.unit_price` remains unchanged after price update in `backend/infrastructure/postgres/menu_repository_test.go`
- [X] T030 [P] [US2] Add HTTP integration tests for price update, stock update, audit log creation, and validation errors in `backend/interface/adapter/http/menu_handler_test.go`
- [X] T031 [P] [US2] Add frontend tests for price and stock editing flows in `frontend/tests/price-stock-editor.spec.ts`

### Implementation for User Story 2

- [X] T032 [US2] Implement `UpdatePrice` with required `PriceHistory` recording and no mutation of historical order item snapshots in `backend/usecase/menu_usecase.go`
- [X] T033 [US2] Implement `UpdateStock` with stock >= 0 validation and automatic `is_available=false` when stock is zero in `backend/usecase/menu_usecase.go`
- [X] T034 [US2] Implement price history and audit log persistence for price and stock changes in `backend/infrastructure/postgres/menu_repository.go`
- [X] T035 [P] [US2] Implement price history repository methods in `backend/usecase/price_history_repository.go`
- [X] T036 [P] [US2] Implement menu audit log repository methods in `backend/usecase/audit_log_repository.go`
- [X] T037 [US2] Implement price and stock HTTP handlers that sync with `GET /api/menu` availability in `backend/interface/adapter/http/menu_handler.go`
- [X] T038 [US2] Implement frontend price and stock actions in `frontend/src/stores/managerMenuStore.ts`
- [X] T039 [P] [US2] Implement price editor UI in `frontend/src/components/PriceEditor.vue`
- [X] T040 [P] [US2] Implement stock editor UI and zero-stock availability state in `frontend/src/components/StockEditor.vue`

**Checkpoint**: User Story 2 is independently shippable when new orders use updated prices, old order snapshots remain unchanged, and zero stock disables ordering.

---

## Phase 5: User Story 3 - Review Stock Warnings (Priority: P2)

**Goal**: Managers can see low-stock items and use the per-item low-stock threshold for warnings.

**Independent Test**: Set item stock below its threshold, confirm it appears in the low-stock list, then change the threshold and confirm the warning list updates.

### Tests for User Story 3

- [X] T041 [P] [US3] Add usecase tests for `GetLowStockItems` and low-stock threshold behavior in `backend/usecase/inventory_usecase_test.go`
- [X] T042 [P] [US3] Add repository tests for low-stock sorting and threshold filtering in `backend/infrastructure/postgres/menu_repository_test.go`
- [X] T043 [P] [US3] Add HTTP integration tests for low-stock query status/schema in `backend/interface/adapter/http/menu_handler_test.go`
- [X] T044 [P] [US3] Add frontend tests for low-stock alert rendering and refresh behavior in `frontend/tests/low-stock-warning.spec.ts`

### Implementation for User Story 3

- [X] T045 [US3] Implement `GetLowStockItems` usecase using per-item `low_stock_threshold` in `backend/usecase/inventory_usecase.go`
- [X] T046 [US3] Implement PostgreSQL low-stock query support in `backend/infrastructure/postgres/inventory_repository.go`
- [X] T047 [US3] Implement low-stock query HTTP handler and presenter response fields in `backend/interface/adapter/http/menu_handler.go`
- [X] T048 [US3] Implement inventory store low-stock loading state and actions in `frontend/src/stores/inventoryStore.ts`
- [X] T049 [US3] Implement `LowStockAlert.vue` for manager warnings in `frontend/src/components/LowStockAlert.vue`
- [X] T050 [US3] Wire low-stock alerts into the manager menu view in `frontend/src/pages/MenuManagementView.vue`

**Checkpoint**: User Story 3 is independently shippable when low-stock warnings are accurate and visible.

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Coverage, quality gates, docs, and final validation across all stories.

- [X] T051 [P] Run backend formatting and tests with coverage >=80% and record results in `backend/coverage.out`
- [X] T052 [P] Run frontend typecheck, lint, tests, and coverage >=80% and record results in `frontend/coverage/coverage-final.json`
- [X] T053 [P] Verify no migration edits or destructive schema/data operations were introduced in `backend/infrastructure/postgres/migrations/`
- [X] T054 [P] Verify manager UI touch targets, mobile/tablet/desktop layouts, and WebView-safe behavior in `frontend/tests/manager-menu-responsive.spec.ts`
- [X] T055 Update quickstart results and any final API shape changes in `specs/003-menu-pricing-stock/quickstart.md`

---

## Dependencies & Execution Order

### Phase Dependencies

- **Phase 1 Setup**: No dependencies.
- **Phase 2 Foundational**: Depends on Phase 1 and blocks all user stories.
- **Phase 3 US1**: Depends on Phase 2; suggested MVP.
- **Phase 4 US2**: Depends on Phase 2 and can run in parallel with US1 after shared repository interfaces exist, but price/stock UI should integrate with the manager store from US1.
- **Phase 5 US3**: Depends on Phase 2 and can run in parallel with US1/US2 after stock fields and repository interfaces exist.
- **Phase 6 Polish**: Depends on selected user stories being complete.

### User Story Dependencies

- **US1 Manage Menu Items (P1)**: Starts after foundational migrations/entities/interfaces.
- **US2 Manage Price And Stock (P1)**: Starts after foundational migrations/entities/interfaces; depends on price history/audit tables.
- **US3 Review Stock Warnings (P2)**: Starts after foundational low-stock threshold migration and stock repository interface.

### Within Each User Story

- Write tests before implementation and confirm they fail.
- Entity validation precedes usecase implementation.
- Usecases precede HTTP handlers and PostgreSQL repositories.
- HTTP handlers validate input before calling usecases.
- Frontend store actions precede dependent Vue components.
- Responsive behavior and WebView constraints must pass before story completion.

### Parallel Opportunities

- T003 and T004 can run in parallel.
- T009 through T013 and T015 can run in parallel after migration filenames are chosen.
- Test tasks inside each user story can run in parallel.
- Backend handler/presenter work and frontend component work can run in parallel after contracts are stable.
- US2 and US3 can run in parallel with US1 after Phase 2 when staffed separately.

---

## Parallel Example: User Story 1

```bash
# Launch tests for US1 together:
Task: "T016 Add usecase tests for CreateMenuItem, UpdateMenuItem, and DeactivateMenuItem soft-delete behavior in backend/tests/menu_usecase_test.go"
Task: "T018 Add HTTP integration tests for menu item CRUD status/schema and validation failures in backend/interface/adapter/http/menu_handler_test.go"
Task: "T019 Add frontend store and page tests for create/edit/deactivate flows in frontend/tests/manager-menu-page.spec.ts"
Task: "T020 Add responsive tests for desktop data table and mobile accordion-card layout in frontend/tests/manager-menu-responsive.spec.ts"

# Launch independent implementation after tests are in place:
Task: "T021 Implement CreateMenuItem, UpdateMenuItem, and DeactivateMenuItem usecases in backend/usecase/menu_usecase.go"
Task: "T025 Implement managerMenuStore.ts actions in frontend/src/stores/managerMenuStore.ts"
Task: "T027 Implement create/edit form in frontend/src/components/MenuItemForm.vue"
```

## Parallel Example: User Story 2

```bash
Task: "T032 Implement UpdatePrice with PriceHistory recording in backend/usecase/menu_usecase.go"
Task: "T039 Implement price editor UI in frontend/src/components/PriceEditor.vue"
Task: "T040 Implement stock editor UI in frontend/src/components/StockEditor.vue"
```

## Parallel Example: User Story 3

```bash
Task: "T041 Add usecase tests for GetLowStockItems in backend/usecase/inventory_usecase_test.go"
Task: "T044 Add frontend tests for low-stock alert rendering in frontend/tests/low-stock-warning.spec.ts"
Task: "T049 Implement LowStockAlert.vue in frontend/src/components/LowStockAlert.vue"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1 setup.
2. Complete Phase 2 foundational migration/entity/interface work.
3. Complete Phase 3 US1 for menu item create/edit/deactivate.
4. Stop and validate US1 independently, including soft-delete and responsive layout tests.

### Incremental Delivery

1. Deliver US1 menu item management.
2. Deliver US2 price/stock updates with price history and stock auto-disable.
3. Deliver US3 low-stock alert visibility.
4. Run Phase 6 quality gates and quickstart validation.

### Quality Gates

- Migration review completed before apply.
- No existing schema/data deletion.
- Backend coverage >=80%.
- Frontend coverage >=80%.
- Endpoint integration tests pass for status/schema.
- Responsive tests cover mobile, tablet, desktop, and WebView-safe behavior.

## Notes

- `[P]` means the task is parallelizable.
- `[US1]`, `[US2]`, and `[US3]` map directly to spec user stories.
- Keep frontend and backend independent; communicate through documented API contracts only.
- Stop for human confirmation if any implementation requires editing or deleting existing schema/data.

## Phase 7: Convergence

- [ ] T056 Add inventory-settings audit logging for low-stock threshold changes in `backend/usecase/inventory_usecase.go` and `backend/interface/adapter/http/settings_handler.go` so threshold updates write an auditable record per `contracts/menu-management-api.md` audit contract (partial)
- [ ] T057 Replace demo-seeded inventory data with backend-backed loading in `frontend/src/pages/MenuManagementView.vue`, `frontend/src/pages/LowStockView.vue`, and `frontend/src/stores/inventoryStore.ts` so menu edits and low-stock warnings share persisted state per FR-008/FR-011/FR-012 (partial)
- [ ] T058 Add a real historical order/bill snapshot integration check against Feature 1 order persistence in `backend/infrastructure/postgres/menu_repository_test.go` and `backend/tests/price_stock_usecase_test.go` so price updates are proven not to mutate historical `OrderItem.unit_price` snapshots per FR-010/SC-004/US2/AC1 (partial)
