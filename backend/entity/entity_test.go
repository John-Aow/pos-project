package entity

import (
	"testing"
	"time"
)

func TestCategoryValidate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		category Category
		wantErr  bool
	}{
		{
			name: "valid",
			category: Category{
				Name:      "Mains",
				SortOrder: 1,
			},
		},
		{
			name: "missing name",
			category: Category{
				SortOrder: 1,
			},
			wantErr: true,
		},
		{
			name: "negative sort",
			category: Category{
				Name:      "Drinks",
				SortOrder: -1,
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			err := tc.category.Validate()
			if tc.wantErr && err == nil {
				t.Fatal("expected error")
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestMenuItemRules(t *testing.T) {
	t.Parallel()

	item := MenuItem{
		Name:              "Green Curry",
		Description:       "Thai curry",
		PriceCents:        12900,
		StockQuantity:     5,
		LowStockThreshold: 3,
		IsActive:          true,
		CategoryID:        1,
	}
	item.RecalculateAvailability()

	if err := item.Validate(); err != nil {
		t.Fatalf("expected valid item, got %v", err)
	}
	if !item.CustomerOrderable() {
		t.Fatal("expected item to be customer orderable")
	}
	if item.LowStock() {
		t.Fatal("expected item above low-stock threshold")
	}

	item.ApplyStockQuantity(0)
	if item.IsAvailable {
		t.Fatal("expected zero stock to disable availability")
	}
	if !item.LowStock() {
		t.Fatal("expected zero stock to be low stock")
	}

	if err := item.ApplyPrice(15000); err != nil {
		t.Fatalf("expected price update to succeed, got %v", err)
	}
	if item.PriceCents != 15000 {
		t.Fatalf("expected price 15000, got %d", item.PriceCents)
	}
	if err := item.ApplyPrice(0); err == nil {
		t.Fatal("expected invalid price error")
	}
}

func TestMenuItemValidationFailures(t *testing.T) {
	t.Parallel()

	tests := []MenuItem{
		{},
		{
			Name:        "No description",
			Description: "",
			PriceCents:  1,
			CategoryID:  1,
		},
		{
			Name:        "Bad price",
			Description: "x",
			PriceCents:  0,
			CategoryID:  1,
		},
		{
			Name:          "Bad stock",
			Description:   "x",
			PriceCents:    1,
			StockQuantity: -1,
			CategoryID:    1,
		},
		{
			Name:              "Bad threshold",
			Description:       "x",
			PriceCents:        1,
			StockQuantity:     1,
			LowStockThreshold: -1,
			CategoryID:        1,
		},
		{
			Name:          "Bad category",
			Description:   "x",
			PriceCents:    1,
			StockQuantity: 1,
		},
	}

	for _, item := range tests {
		item := item
		t.Run(item.Name, func(t *testing.T) {
			t.Parallel()

			if err := item.Validate(); err == nil {
				t.Fatal("expected validation error")
			}
		})
	}
}

func TestPriceHistoryValidate(t *testing.T) {
	t.Parallel()

	history := PriceHistory{
		MenuItemID:         1,
		PreviousPriceCents: 12900,
		NewPriceCents:      13900,
		Actor:              "manager-1",
		CreatedAt:          time.Now(),
	}
	if err := history.Validate(); err != nil {
		t.Fatalf("expected valid price history, got %v", err)
	}

	tests := []PriceHistory{
		{},
		{
			MenuItemID:         1,
			PreviousPriceCents: 0,
			NewPriceCents:      13900,
			Actor:              "manager-1",
			CreatedAt:          time.Now(),
		},
		{
			MenuItemID:         1,
			PreviousPriceCents: 13900,
			NewPriceCents:      13900,
			Actor:              "manager-1",
			CreatedAt:          time.Now(),
		},
		{
			MenuItemID:         1,
			PreviousPriceCents: 12900,
			NewPriceCents:      13900,
			CreatedAt:          time.Now(),
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run("invalid price history", func(t *testing.T) {
			t.Parallel()

			if err := tc.Validate(); err == nil {
				t.Fatal("expected validation error")
			}
		})
	}
}

func TestMenuAuditLogValidate(t *testing.T) {
	t.Parallel()

	log := MenuAuditLog{
		MenuItemID: 1,
		ActionType: AuditActionStock,
		Actor:      "manager-1",
		CreatedAt:  time.Now(),
	}
	if err := log.Validate(); err != nil {
		t.Fatalf("expected valid menu audit log, got %v", err)
	}

	tests := []MenuAuditLog{
		{},
		{MenuItemID: 1, Actor: "manager-1", CreatedAt: time.Now()},
		{MenuItemID: 1, ActionType: AuditActionStock, CreatedAt: time.Now()},
		{MenuItemID: 1, ActionType: AuditActionStock, Actor: "manager-1"},
	}
	for _, tc := range tests {
		tc := tc
		t.Run("invalid menu audit log", func(t *testing.T) {
			t.Parallel()

			if err := tc.Validate(); err == nil {
				t.Fatal("expected validation error")
			}
		})
	}
}

func TestInventorySettingsValidate(t *testing.T) {
	t.Parallel()

	if err := (InventorySettings{LowStockThreshold: 2}).Validate(); err != nil {
		t.Fatalf("expected valid settings, got %v", err)
	}
	if err := (InventorySettings{LowStockThreshold: -1}).Validate(); err == nil {
		t.Fatal("expected threshold validation error")
	}
}

func TestMenuItemAuditEntryValidate(t *testing.T) {
	t.Parallel()

	entry := MenuItemAuditEntry{
		MenuItemID: 1,
		ActionType: AuditActionCreate,
		Actor:      "manager-1",
		CreatedAt:  time.Now(),
	}
	if err := entry.Validate(); err != nil {
		t.Fatalf("expected valid audit entry, got %v", err)
	}

	entry.Actor = ""
	if err := entry.Validate(); err == nil {
		t.Fatal("expected validation error for missing actor")
	}

	if err := (MenuItemAuditEntry{}).Validate(); err == nil {
		t.Fatal("expected validation error for missing audit fields")
	}
}
