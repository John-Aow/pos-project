package entity

import (
	"errors"
	"fmt"
	"time"
)

// InventorySettings stores the shared low-stock threshold.
type InventorySettings struct {
	ID                int64
	LowStockThreshold int
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

// Validate ensures the threshold can be saved safely.
func (s InventorySettings) Validate() error {
	if s.LowStockThreshold < 0 {
		return fmt.Errorf("low stock threshold must be zero or greater")
	}
	if s.ID < 0 {
		return errors.New("inventory settings id must not be negative")
	}
	return nil
}
