package entity

import (
	"errors"
	"strings"
	"time"
)

// PriceHistory records every accepted menu item price change.
type PriceHistory struct {
	ID                 int64
	MenuItemID         int64
	PreviousPriceCents int64
	NewPriceCents      int64
	Actor              string
	Reason             string
	CreatedAt          time.Time
}

// Validate ensures price history is auditable and financially meaningful.
func (h PriceHistory) Validate() error {
	if h.MenuItemID <= 0 {
		return errors.New("price history menu item id is required")
	}
	if h.PreviousPriceCents <= 0 {
		return errors.New("price history previous price must be greater than zero")
	}
	if h.NewPriceCents <= 0 {
		return errors.New("price history new price must be greater than zero")
	}
	if h.PreviousPriceCents == h.NewPriceCents {
		return errors.New("price history requires a changed price")
	}
	if strings.TrimSpace(h.Actor) == "" {
		return errors.New("price history actor is required")
	}
	if h.CreatedAt.IsZero() {
		return errors.New("price history created_at is required")
	}
	return nil
}
