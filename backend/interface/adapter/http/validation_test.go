package httpadapter

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestValidationHelpersExerciseAllBranches(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	if _, err := ParseInt64PathValue(req, "id"); err == nil {
		t.Fatal("expected missing path value error")
	}

	req = httptest.NewRequest(http.MethodGet, "/", nil)
	req.SetPathValue("id", "not-a-number")
	if _, err := ParseInt64PathValue(req, "id"); err == nil {
		t.Fatal("expected invalid path value error")
	}

	if err := ValidateMenuItemInput("", "desc", 1, 0, 0, 1); err == nil {
		t.Fatal("expected menu item name validation error")
	}
	if err := ValidateMenuItemInput("name", "", 1, 0, 0, 1); err == nil {
		t.Fatal("expected menu item description validation error")
	}
	if err := ValidateMenuItemInput("name", "desc", 0, 0, 0, 1); err == nil {
		t.Fatal("expected menu item price validation error")
	}
	if err := ValidateMenuItemInput("name", "desc", 1, -1, 0, 1); err == nil {
		t.Fatal("expected menu item stock validation error")
	}
	if err := ValidateMenuItemInput("name", "desc", 1, 0, -1, 1); err == nil {
		t.Fatal("expected menu item threshold validation error")
	}
	if err := ValidateMenuItemInput("name", "desc", 1, 0, 0, 0); err == nil {
		t.Fatal("expected menu item category validation error")
	}
	if err := ValidateMenuItemInput("name", "desc", 1, 0, 0, 1); err != nil {
		t.Fatalf("expected valid menu item input, got %v", err)
	}

	if err := ValidatePriceUpdateInput(0, "manager-1"); err == nil {
		t.Fatal("expected price validation error")
	}
	if err := ValidatePriceUpdateInput(1, ""); err == nil {
		t.Fatal("expected price actor validation error")
	}
	if err := ValidatePriceUpdateInput(1, "manager-1"); err != nil {
		t.Fatalf("expected valid price input, got %v", err)
	}

	if err := ValidateStockUpdateInput(-1, "manager-1"); err == nil {
		t.Fatal("expected stock validation error")
	}
	if err := ValidateStockUpdateInput(0, ""); err == nil {
		t.Fatal("expected stock actor validation error")
	}
	if err := ValidateStockUpdateInput(0, "manager-1"); err != nil {
		t.Fatalf("expected valid stock input, got %v", err)
	}

	if err := ValidateLowStockThresholdInput(-1, "manager-1"); err == nil {
		t.Fatal("expected low-stock validation error")
	}
	if err := ValidateLowStockThresholdInput(1, ""); err == nil {
		t.Fatal("expected low-stock actor validation error")
	}
	if err := ValidateLowStockThresholdInput(1, "manager-1"); err != nil {
		t.Fatalf("expected valid low-stock input, got %v", err)
	}

	req = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"name":"ok"}`))
	var payload struct {
		Name string `json:"name"`
	}
	if err := DecodeJSONBody(req, &payload); err != nil {
		t.Fatalf("DecodeJSONBody returned error: %v", err)
	}
}
