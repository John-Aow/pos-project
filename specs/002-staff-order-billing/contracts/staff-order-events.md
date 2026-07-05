# Contract: Staff Order Realtime Events

## Source

Feature 1 `order_notifier` publishes order events consumed by the staff billing interface.

## Event: `order.created`

**Fields**
- `event_id`
- `order_id`
- `occurred_at`
- `status`

**Rules**
- The frontend treats the event as a refresh hint and loads order details from `GET /api/staff/orders`.
- Duplicate events must not duplicate orders in the staff store.

## Event: `order.updated`

**Fields**
- `event_id`
- `order_id`
- `occurred_at`
- `status`

**Rules**
- The staff store updates or refreshes the matching order.
- If the order is no longer active, it is removed from the active list or moved to the appropriate view.

## Reconnect Behavior

- On websocket connect or reconnect, the staff UI refreshes `GET /api/staff/orders`.
- If websocket events are missed, the API response remains the source of truth.
- UI must show a recoverable connection state without allowing duplicate bill mutations.
