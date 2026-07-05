package entity

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

// MenuItem is the managed menu record used by customers and staff.
type MenuItem struct {
	ID            int64
	Name          string
	Description   string
	PriceCents    int64
	StockQuantity int
	IsAvailable   bool
	IsActive      bool
	CategoryID    int64
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// Validate checks the business rules for a menu item.
func (m MenuItem) Validate() error {
	switch {
	case strings.TrimSpace(m.Name) == "":
		return errors.New("menu item name is required")
	case strings.TrimSpace(m.Description) == "":
		return errors.New("menu item description is required")
	case m.PriceCents <= 0:
		return errors.New("menu item price must be greater than zero")
	case m.StockQuantity < 0:
		return fmt.Errorf("menu item stock quantity must be zero or greater")
	case m.CategoryID <= 0:
		return errors.New("menu item category id is required")
	}

	if m.StockQuantity == 0 && m.IsAvailable {
		return errors.New("menu item must be unavailable when stock is zero")
	}
	if m.StockQuantity > 0 && m.IsActive && !m.IsAvailable {
		return errors.New("active menu items with stock must be available")
	}

	return nil
}

// RecalculateAvailability keeps availability aligned with stock and active state.
func (m *MenuItem) RecalculateAvailability() {
	m.IsAvailable = m.IsActive && m.StockQuantity > 0
}

// CustomerOrderable reports whether customers can order the item.
func (m MenuItem) CustomerOrderable() bool {
	return m.IsActive && m.IsAvailable && m.StockQuantity > 0
}

// ApplyStockQuantity updates the quantity and availability together.
func (m *MenuItem) ApplyStockQuantity(quantity int) {
	m.StockQuantity = quantity
	m.RecalculateAvailability()
}

// ApplyPrice updates the item price while keeping validation responsibility local.
func (m *MenuItem) ApplyPrice(priceCents int64) error {
	if priceCents <= 0 {
		return errors.New("menu item price must be greater than zero")
	}

	m.PriceCents = priceCents
	return nil
}
