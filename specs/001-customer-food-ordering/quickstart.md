# Quickstart: Customer Food Ordering

## Prerequisites

- Frontend dependencies installed in `frontend/`
- Backend dependencies available in `backend/`
- PostgreSQL migrations applied for catalog and customer orders

## Run Locally

Frontend:

```bash
cd frontend
npm run dev
```

Backend:

```bash
cd backend
go test ./...
```

When the backend server is wired, run it from:

```bash
cd backend
go run ./cmd/server
```

## Scenario 1: Browse Customer Menu

1. Open the customer ordering route in the frontend.
2. Verify available menu items show name, description, price, and orderable action.
3. Verify inactive or out-of-stock items are disabled or absent from add-to-cart actions.
4. Verify attempting to select an unavailable item shows a clear message.

Expected result: Customer can understand what is available and cannot order unavailable items.

## Scenario 2: Build Cart With Notes

1. Add an available item to the cart.
2. Increase and decrease the quantity.
3. Add an item-specific note.
4. Remove the item.
5. Re-add one or more items and confirm line totals and order total update.

Expected result: Cart state remains accurate and note text is preserved for submission.

## Scenario 3: Submit Order

1. Submit a non-empty cart.
2. Confirm the backend returns an order number and persisted total.
3. Confirm the frontend shows the same order number and total.
4. Confirm `GET /api/staff/orders` includes the submitted order within 5 seconds.

Expected result: Customer receives confirmation and staff can see the incoming order.

## Scenario 4: Submission Conflict

1. Add an available item to the cart.
2. Make that item unavailable before submission.
3. Submit the cart.

Expected result: Submission is blocked with a customer-readable message identifying the unavailable item.

## Scenario 5: Current Order Status

1. Submit an order and keep the order number.
2. Open the current order status view.
3. Update the order status through backend or staff workflow.
4. Refresh the customer status view.

Expected result: Customer sees pending, preparing, or served as the latest status.

## Responsive And WebView Validation Targets

Validate all interactive customer ordering screens at:

- Mobile: 360px and wider
- Tablet: 768px and wider
- Desktop: 1024px and wider
- Portrait orientation
- Landscape orientation
- Native iOS and Android WebView or equivalent viewport simulation

Required checks:

- Touch targets are at least 44px by 44px.
- The flow does not depend on hover-only behavior.
- In-app navigation exists between menu, cart, confirmation, and status.
- Layouts avoid fixed pixel widths that break small screens.
- The viewport meta configuration prevents unwanted zooming or scaling.
- Text does not overlap controls or other content.

## Quality Gates

Frontend:

```bash
cd frontend
npm run lint
npm run test
npm run build
```

Backend:

```bash
cd backend
gofmt -w .
go test ./...
```
