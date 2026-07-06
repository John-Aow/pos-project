# Tasks: Customer Food Ordering

**Input**: Design documents from `/specs/001-customer-food-ordering/`

**Available Docs**: `spec.md`, `checklists/requirements.md`

**Missing Docs Noted**: `plan.md`, `research.md`, `data-model.md`, `contracts/`, and `quickstart.md` are not present for this feature. Tasks are generated from the feature specification, project constitution, and current Vue/Go monorepo structure.

**Tests**: Required by the project constitution for each acceptance criterion, POS risk requirement, Clean Architecture usecase, API validation, database migration, responsive frontend behavior, and WebView behavior. Frontend and backend coverage must remain at or above 80%.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel because it touches different files and has no dependency on incomplete tasks
- **[Story]**: Which user story this task belongs to
- Every task includes an exact file path

## Path Conventions

- **Frontend**: `frontend/src/`, `frontend/tests/`
- **Backend**: `backend/entity/`, `backend/usecase/`, `backend/interface/adapter/`, `backend/infrastructure/`, `backend/tests/`
- **Shared**: `packages/contracts/`
- Frontend and backend must build and test independently. Do not create direct implementation imports between `frontend/` and `backend/`.

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Establish the missing feature design and shared contract surfaces needed before implementation.

- [X] T001 Create customer ordering implementation plan in specs/001-customer-food-ordering/plan.md
- [X] T002 Create customer ordering data model in specs/001-customer-food-ordering/data-model.md
- [X] T003 Create customer ordering quickstart scenarios in specs/001-customer-food-ordering/quickstart.md
- [X] T004 [P] Create customer ordering API contract directory in specs/001-customer-food-ordering/contracts/customer-ordering.openapi.yaml
- [X] T005 [P] Create shared customer ordering contract notes in packages/contracts/customer-ordering.md
- [X] T006 [P] Document responsive and WebView validation targets for customer ordering in specs/001-customer-food-ordering/quickstart.md

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core entities, persistence, routing, and frontend shell that all customer-ordering stories depend on.

**Critical**: No user story work can begin until this phase is complete.

- [ ] T007 Create order and bill database migration in backend/infrastructure/postgres/migrations/20260705_create_customer_orders.sql
- [ ] T008 Create order and bill rollback migration in backend/infrastructure/postgres/migrations/20260705_create_customer_orders.down.sql
- [ ] T009 [P] Create order domain entities with validation in backend/entity/order.go
- [ ] T010 [P] Create bill domain entity with validation in backend/entity/bill.go
- [ ] T011 [P] Create order status domain entity in backend/entity/order_status.go
- [ ] T012 Create order repository interface in backend/usecase/order_repository.go
- [ ] T013 Create order usecase skeleton and shared input types in backend/usecase/order_usecase.go
- [ ] T014 Create Postgres order repository skeleton in backend/infrastructure/postgres/order_repository.go
- [ ] T015 Create customer order HTTP handler skeleton in backend/interface/adapter/http/order_handler.go
- [ ] T016 Update backend router with customer order routes in backend/interface/adapter/http/router.go
- [ ] T017 [P] Create customer ordering API service client in frontend/src/services/customerOrderingApi.ts
- [ ] T018 [P] Create customer ordering Pinia store skeleton in frontend/src/stores/customerOrderStore.ts
- [ ] T019 [P] Create customer ordering route/view switch in frontend/src/App.vue
- [ ] T020 [P] Create customer ordering stylesheet in frontend/src/styles/customer-ordering.css

**Checkpoint**: Foundation ready - user story implementation can now begin in priority order or in parallel by separate contributors.

---

## Phase 3: User Story 1 - Browse Available Menu (Priority: P1) MVP

**Goal**: Customers can open the ordering web app, view available menu items with name, price, description, and availability status, and unavailable items cannot be ordered.

**Independent Test**: Open the customer menu and verify orderable items show complete details while unavailable or out-of-stock items are clearly disabled and cannot be selected.

### Tests for User Story 1

- [ ] T021 [P] [US1] Add backend usecase tests for customer menu availability filtering in backend/tests/customer_menu_usecase_test.go
- [ ] T022 [P] [US1] Add backend HTTP contract tests for GET /api/menu in backend/interface/adapter/http/customer_menu_handler_test.go
- [ ] T023 [P] [US1] Add frontend component tests for menu item display and disabled unavailable items in frontend/tests/customer-menu.spec.ts
- [ ] T024 [P] [US1] Add responsive and WebView tests for customer menu mobile tablet desktop touch targets in frontend/tests/customer-menu-responsive.spec.ts
- [ ] T025 [P] [US1] Add Postgres repository tests for available menu filtering in backend/infrastructure/postgres/customer_menu_repository_test.go

### Implementation for User Story 1

- [ ] T026 [US1] Extend customer menu presenter with availability fields in backend/interface/adapter/http/menu_presenter.go
- [ ] T027 [US1] Ensure ListCustomerMenu rejects inactive and out-of-stock items in backend/usecase/menu_usecase.go
- [ ] T028 [US1] Ensure Postgres ListAvailableMenuItems filters active in-stock items in backend/infrastructure/postgres/menu_repository.go
- [ ] T029 [US1] Implement customer menu loading action in frontend/src/stores/customerOrderStore.ts
- [ ] T030 [US1] Implement customer menu page in frontend/src/pages/CustomerMenuPage.vue
- [ ] T031 [US1] Implement customer menu item component with disabled unavailable state in frontend/src/components/CustomerMenuItem.vue
- [ ] T032 [US1] Add customer menu route navigation state in frontend/src/App.vue
- [ ] T033 [US1] Add mobile-first customer menu styles and 44px tap targets in frontend/src/styles/customer-ordering.css

**Checkpoint**: User Story 1 is fully functional and testable independently.

---

## Phase 4: User Story 2 - Build Cart With Notes (Priority: P1)

**Goal**: Customers can add available items to a cart, adjust quantities, remove items, add item notes, and see line totals and order total update correctly.

**Independent Test**: Add available items, change quantities, remove an item, add a note, and confirm cart totals update correctly after each change.

### Tests for User Story 2

- [ ] T034 [P] [US2] Add frontend store tests for cart add update remove note and totals in frontend/tests/customer-cart-store.spec.ts
- [ ] T035 [P] [US2] Add frontend component tests for cart controls and note preservation in frontend/tests/customer-cart.spec.ts
- [ ] T036 [P] [US2] Add frontend validation tests for empty cart zero quantity negative quantity and long notes in frontend/tests/customer-cart-validation.spec.ts
- [ ] T037 [P] [US2] Add responsive and WebView tests for cart controls on mobile tablet desktop in frontend/tests/customer-cart-responsive.spec.ts

### Implementation for User Story 2

- [ ] T038 [US2] Add cart state line total order total and note actions in frontend/src/stores/customerOrderStore.ts
- [ ] T039 [US2] Implement customer cart page in frontend/src/pages/CustomerCartPage.vue
- [ ] T040 [US2] Implement cart item quantity and remove controls in frontend/src/components/CustomerCartItem.vue
- [ ] T041 [US2] Implement item note editor with length validation in frontend/src/components/CustomerItemNoteEditor.vue
- [ ] T042 [US2] Add cart summary and order total component in frontend/src/components/CustomerCartSummary.vue
- [ ] T043 [US2] Add in-app navigation between menu and cart in frontend/src/App.vue
- [ ] T044 [US2] Add mobile-first cart styles and touch-friendly controls in frontend/src/styles/customer-ordering.css

**Checkpoint**: User Stories 1 and 2 are fully functional and testable independently.

---

## Phase 5: User Story 3 - Submit Order And Receive Confirmation (Priority: P1)

**Goal**: Customers can submit a valid cart, receive an order number, and staff can see the new order with ordered items, quantities, notes, and total.

**Independent Test**: Submit a cart with available items, confirm an order number is shown, and verify the staff incoming order list can read the created order immediately.

### Tests for User Story 3

- [ ] T045 [P] [US3] Add backend usecase tests for submit order totals availability recheck idempotency and audit data in backend/tests/customer_order_usecase_test.go
- [ ] T046 [P] [US3] Add backend HTTP validation tests for POST /api/orders empty cart invalid quantity unavailable item and duplicate key in backend/interface/adapter/http/order_handler_test.go
- [ ] T047 [P] [US3] Add Postgres repository tests for order bill line item persistence in backend/infrastructure/postgres/order_repository_test.go
- [ ] T048 [P] [US3] Add frontend submit flow tests for confirmation unavailable item error and interrupted confirmation recovery in frontend/tests/customer-order-submit.spec.ts
- [ ] T049 [P] [US3] Add staff visibility contract tests for incoming orders API in backend/interface/adapter/http/staff_order_handler_test.go
- [ ] T050 [P] [US3] Add offline and duplicate submission tests for customer ordering store in frontend/tests/customer-order-offline.spec.ts

### Implementation for User Story 3

- [ ] T051 [US3] Implement order submission input validation and total calculation in backend/usecase/order_usecase.go
- [ ] T052 [US3] Implement order number generation and idempotency key handling in backend/usecase/order_usecase.go
- [ ] T053 [US3] Implement order audit data capture for submitted time items quantities notes totals and status in backend/usecase/order_usecase.go
- [ ] T054 [US3] Implement Postgres order bill and line item persistence in backend/infrastructure/postgres/order_repository.go
- [ ] T055 [US3] Implement POST /api/orders request decoding validation and response in backend/interface/adapter/http/order_handler.go
- [ ] T056 [US3] Implement GET /api/staff/orders incoming order list for staff visibility in backend/interface/adapter/http/order_handler.go
- [ ] T057 [US3] Update backend router with POST /api/orders and GET /api/staff/orders in backend/interface/adapter/http/router.go
- [ ] T058 [US3] Implement submit order API call and idempotency key persistence in frontend/src/services/customerOrderingApi.ts
- [ ] T059 [US3] Implement submit order action offline state and confirmation recovery in frontend/src/stores/customerOrderStore.ts
- [ ] T060 [US3] Implement order confirmation page with order number and total in frontend/src/pages/CustomerConfirmationPage.vue
- [ ] T061 [US3] Implement unavailable item error display in frontend/src/components/CustomerOrderError.vue
- [ ] T062 [US3] Add confirmation and submission styles in frontend/src/styles/customer-ordering.css

**Checkpoint**: User Stories 1, 2, and 3 provide the complete P1 customer ordering flow.

---

## Phase 6: User Story 4 - View Current Order Status (Priority: P2)

**Goal**: Customers can view the current status of their confirmed order, including pending, preparing, and served.

**Independent Test**: Place an order, update its status through backend or staff workflow, and confirm the customer status view reflects the latest state.

### Tests for User Story 4

- [ ] T063 [P] [US4] Add backend usecase tests for order status lookup and status transition visibility in backend/tests/order_status_usecase_test.go
- [ ] T064 [P] [US4] Add backend HTTP contract tests for GET /api/orders/{orderNumber}/status in backend/interface/adapter/http/order_status_handler_test.go
- [ ] T065 [P] [US4] Add frontend status view tests for pending preparing served and missing order in frontend/tests/customer-order-status.spec.ts
- [ ] T066 [P] [US4] Add responsive and WebView tests for status view navigation in frontend/tests/customer-order-status-responsive.spec.ts

### Implementation for User Story 4

- [ ] T067 [US4] Add order status repository methods in backend/usecase/order_repository.go
- [ ] T068 [US4] Implement order status lookup in backend/usecase/order_usecase.go
- [ ] T069 [US4] Implement Postgres order status lookup in backend/infrastructure/postgres/order_repository.go
- [ ] T070 [US4] Implement GET /api/orders/{orderNumber}/status handler in backend/interface/adapter/http/order_handler.go
- [ ] T071 [US4] Update backend router with GET /api/orders/{orderNumber}/status in backend/interface/adapter/http/router.go
- [ ] T072 [US4] Add order status API call in frontend/src/services/customerOrderingApi.ts
- [ ] T073 [US4] Add status loading state and refresh action in frontend/src/stores/customerOrderStore.ts
- [ ] T074 [US4] Implement current order status page in frontend/src/pages/CustomerOrderStatusPage.vue
- [ ] T075 [US4] Add status page navigation from confirmation in frontend/src/App.vue
- [ ] T076 [US4] Add status page styles and touch-friendly refresh control in frontend/src/styles/customer-ordering.css

**Checkpoint**: All customer ordering user stories are independently functional.

---

## Phase 7: Polish & Cross-Cutting Concerns

**Purpose**: Quality gates and refinements that affect multiple stories.

- [ ] T077 [P] Update customer ordering quickstart validation steps in specs/001-customer-food-ordering/quickstart.md
- [ ] T078 [P] Update shared API contract documentation in packages/contracts/customer-ordering.md
- [ ] T079 [P] Add customer ordering feature notes to frontend README in frontend/README.md
- [ ] T080 [P] Add backend customer ordering API notes to backend/README.md
- [ ] T081 Run frontend lint command and fix customer ordering issues configured in frontend/package.json
- [ ] T082 Run frontend tests with coverage and keep coverage at or above 80% using frontend/vitest.config.ts
- [ ] T083 Run backend gofmt and fix formatting for customer ordering files tracked from backend/go.mod
- [ ] T084 Run backend tests with coverage and keep coverage at or above 80% using backend/go.mod
- [ ] T085 Validate customer order migrations apply and roll back cleanly using backend/infrastructure/postgres/migrations/20260705_create_customer_orders.sql
- [ ] T086 Validate customer ordering quickstart end-to-end in specs/001-customer-food-ordering/quickstart.md
- [ ] T087 Review frontend for mobile 360px tablet 768px desktop 1024px portrait landscape and WebView behavior in frontend/src/styles/customer-ordering.css
- [ ] T088 Review order submission audit and idempotency behavior against constitution in backend/usecase/order_usecase.go

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately.
- **Foundational (Phase 2)**: Depends on Setup completion - blocks all user stories.
- **User Stories (Phase 3+)**: Depend on Foundational completion.
- **Polish (Phase 7)**: Depends on all desired user stories being complete.

### User Story Dependencies

- **User Story 1 (P1)**: Can start after Foundational - no dependency on other stories.
- **User Story 2 (P1)**: Depends on User Story 1 menu data being available in `frontend/src/stores/customerOrderStore.ts`.
- **User Story 3 (P1)**: Depends on User Stories 1 and 2 because it submits selected cart items.
- **User Story 4 (P2)**: Depends on User Story 3 because it displays a confirmed order status.

### Within Each User Story

- Tests must be written first and fail before implementation.
- Backend entities before usecases.
- Backend usecases before HTTP handlers.
- HTTP handlers before frontend API integration.
- Frontend stores before dependent components.
- Responsive layout, touch targets, and in-app navigation before story completion.
- API input validation before calling usecases.

---

## Parallel Opportunities

- T004, T005, and T006 can run in parallel after T001-T003 are started.
- T009, T010, T011, T017, T018, T019, and T020 can run in parallel during Foundation.
- All test tasks within each story can run in parallel because they touch separate frontend, backend, and repository test files.
- Backend implementation and frontend component work can run in parallel once contracts are stable.
- Polish documentation tasks T077-T080 can run in parallel.

---

## Parallel Example: User Story 1

```bash
Task: "T021 [P] [US1] Add backend usecase tests for customer menu availability filtering in backend/tests/customer_menu_usecase_test.go"
Task: "T022 [P] [US1] Add backend HTTP contract tests for GET /api/menu in backend/interface/adapter/http/customer_menu_handler_test.go"
Task: "T023 [P] [US1] Add frontend component tests for menu item display and disabled unavailable items in frontend/tests/customer-menu.spec.ts"
Task: "T024 [P] [US1] Add responsive and WebView tests for customer menu mobile tablet desktop touch targets in frontend/tests/customer-menu-responsive.spec.ts"
```

---

## Parallel Example: User Story 2

```bash
Task: "T034 [P] [US2] Add frontend store tests for cart add update remove note and totals in frontend/tests/customer-cart-store.spec.ts"
Task: "T035 [P] [US2] Add frontend component tests for cart controls and note preservation in frontend/tests/customer-cart.spec.ts"
Task: "T036 [P] [US2] Add frontend validation tests for empty cart zero quantity negative quantity and long notes in frontend/tests/customer-cart-validation.spec.ts"
Task: "T037 [P] [US2] Add responsive and WebView tests for cart controls on mobile tablet desktop in frontend/tests/customer-cart-responsive.spec.ts"
```

---

## Parallel Example: User Story 3

```bash
Task: "T045 [P] [US3] Add backend usecase tests for submit order totals availability recheck idempotency and audit data in backend/tests/customer_order_usecase_test.go"
Task: "T046 [P] [US3] Add backend HTTP validation tests for POST /api/orders empty cart invalid quantity unavailable item and duplicate key in backend/interface/adapter/http/order_handler_test.go"
Task: "T048 [P] [US3] Add frontend submit flow tests for confirmation unavailable item error and interrupted confirmation recovery in frontend/tests/customer-order-submit.spec.ts"
Task: "T050 [P] [US3] Add offline and duplicate submission tests for customer ordering store in frontend/tests/customer-order-offline.spec.ts"
```

---

## Parallel Example: User Story 4

```bash
Task: "T063 [P] [US4] Add backend usecase tests for order status lookup and status transition visibility in backend/tests/order_status_usecase_test.go"
Task: "T064 [P] [US4] Add backend HTTP contract tests for GET /api/orders/{orderNumber}/status in backend/interface/adapter/http/order_status_handler_test.go"
Task: "T065 [P] [US4] Add frontend status view tests for pending preparing served and missing order in frontend/tests/customer-order-status.spec.ts"
Task: "T066 [P] [US4] Add responsive and WebView tests for status view navigation in frontend/tests/customer-order-status-responsive.spec.ts"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup.
2. Complete Phase 2: Foundational.
3. Complete Phase 3: User Story 1.
4. Stop and validate customer menu browsing independently.
5. Demo the ordering menu at the customer frontend route before adding cart and submission.

### P1 Ordering Flow

1. Complete User Story 2 for cart building.
2. Complete User Story 3 for order submission and confirmation.
3. Validate the complete customer ordering flow from menu to confirmation.

### Incremental Delivery

1. Add User Story 4 for order status after the P1 ordering flow is stable.
2. Run all polish and quality gate tasks.
3. Validate quickstart, coverage, migrations, responsive behavior, and audit/idempotency behavior before marking the feature complete.
