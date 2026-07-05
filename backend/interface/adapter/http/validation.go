package httpadapter

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"
)

const maxRequestBodyBytes int64 = 1 << 20

// ErrInvalidJSON indicates that the request payload could not be decoded.
var ErrInvalidJSON = errors.New("invalid json payload")

// DecodeJSONBody decodes a JSON request body with a strict field check.
func DecodeJSONBody(r *http.Request, dst any) error {
	if r.Body == nil {
		return errors.New("request body is required")
	}

	decoder := json.NewDecoder(io.LimitReader(r.Body, maxRequestBodyBytes))
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(dst); err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidJSON, err)
	}

	return nil
}

// ValidateRequiredString checks that a string is present and within bounds.
func ValidateRequiredString(field, value string, minLen, maxLen int) error {
	trimmed := strings.TrimSpace(value)
	length := utf8.RuneCountInString(trimmed)
	if trimmed == "" {
		return fmt.Errorf("%s is required", field)
	}
	if minLen > 0 && length < minLen {
		return fmt.Errorf("%s must be at least %d characters", field, minLen)
	}
	if maxLen > 0 && length > maxLen {
		return fmt.Errorf("%s must be at most %d characters", field, maxLen)
	}
	return nil
}

// ValidatePositiveInt ensures a numeric value is greater than zero.
func ValidatePositiveInt(field string, value int) error {
	if value <= 0 {
		return fmt.Errorf("%s must be greater than zero", field)
	}
	return nil
}

// ValidateNonNegativeInt ensures a numeric value is zero or greater.
func ValidateNonNegativeInt(field string, value int) error {
	if value < 0 {
		return fmt.Errorf("%s must be zero or greater", field)
	}
	return nil
}

// ValidatePositiveInt64 ensures a numeric value is greater than zero.
func ValidatePositiveInt64(field string, value int64) error {
	if value <= 0 {
		return fmt.Errorf("%s must be greater than zero", field)
	}
	return nil
}

// ValidateMenuItemInput checks manager create/edit payload fields.
func ValidateMenuItemInput(name, description string, priceCents int64, stockQuantity int, lowStockThreshold int, categoryID int64) error {
	if err := ValidateRequiredString("name", name, 1, 120); err != nil {
		return err
	}
	if err := ValidateRequiredString("description", description, 1, 500); err != nil {
		return err
	}
	if err := ValidatePositiveInt64("price_cents", priceCents); err != nil {
		return err
	}
	if err := ValidateNonNegativeInt("stock_quantity", stockQuantity); err != nil {
		return err
	}
	if err := ValidateNonNegativeInt("low_stock_threshold", lowStockThreshold); err != nil {
		return err
	}
	if categoryID <= 0 {
		return fmt.Errorf("category_id must be greater than zero")
	}
	return nil
}

// ValidatePriceUpdateInput checks manager price update payload fields.
func ValidatePriceUpdateInput(priceCents int64, actor string) error {
	if err := ValidatePositiveInt64("price_cents", priceCents); err != nil {
		return err
	}
	return ValidateRequiredString("actor", actor, 1, 120)
}

// ValidateStockUpdateInput checks manager stock update payload fields.
func ValidateStockUpdateInput(stockQuantity int, actor string) error {
	if err := ValidateNonNegativeInt("stock_quantity", stockQuantity); err != nil {
		return err
	}
	return ValidateRequiredString("actor", actor, 1, 120)
}

// ValidateLowStockThresholdInput checks manager low-stock threshold payload fields.
func ValidateLowStockThresholdInput(lowStockThreshold int, actor string) error {
	if err := ValidateNonNegativeInt("low_stock_threshold", lowStockThreshold); err != nil {
		return err
	}
	return ValidateRequiredString("actor", actor, 1, 120)
}

// ParseInt64PathValue reads a route parameter from the new ServeMux path API.
func ParseInt64PathValue(r *http.Request, key string) (int64, error) {
	raw := strings.TrimSpace(r.PathValue(key))
	if raw == "" {
		return 0, fmt.Errorf("%s is required", key)
	}

	value, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("%s must be a valid integer", key)
	}

	return value, nil
}

// WriteJSON writes a JSON response with the supplied status code.
func WriteJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
