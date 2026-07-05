package httpadapter

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"pos-project/backend/entity"
	"pos-project/backend/usecase"
)

func TestDecodeJSONBodyAndValidationHelpers(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"name":"ok"}`))
	var payload struct {
		Name string `json:"name"`
	}
	if err := DecodeJSONBody(req, &payload); err != nil {
		t.Fatalf("DecodeJSONBody returned error: %v", err)
	}
	if payload.Name != "ok" {
		t.Fatalf("unexpected payload: %#v", payload)
	}

	if err := ValidateRequiredString("name", "  curry  ", 2, 20); err != nil {
		t.Fatalf("ValidateRequiredString returned error: %v", err)
	}
	if err := ValidatePositiveInt("price", 1); err != nil {
		t.Fatalf("ValidatePositiveInt returned error: %v", err)
	}
	if err := ValidateNonNegativeInt("stock", 0); err != nil {
		t.Fatalf("ValidateNonNegativeInt returned error: %v", err)
	}

	req.SetPathValue("id", "12")
	id, err := ParseInt64PathValue(req, "id")
	if err != nil {
		t.Fatalf("ParseInt64PathValue returned error: %v", err)
	}
	if id != 12 {
		t.Fatalf("unexpected id: %d", id)
	}
}

func TestDecodeJSONBodyRejectsEmptyBody(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodPost, "/", nil)
	var payload struct{}
	if err := DecodeJSONBody(req, &payload); err == nil {
		t.Fatal("expected missing body error")
	}
}

func TestDecodeJSONBodyRejectsUnknownFields(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"name":"ok","unknown":true}`))
	var payload struct {
		Name string `json:"name"`
	}
	if err := DecodeJSONBody(req, &payload); err == nil {
		t.Fatal("expected unknown field error")
	}
}

func TestValidationHelpersRejectInvalidValues(t *testing.T) {
	t.Parallel()

	if err := ValidateRequiredString("name", " ", 1, 5); err == nil {
		t.Fatal("expected required string error")
	}
	if err := ValidateRequiredString("name", "longer-than-five", 1, 5); err == nil {
		t.Fatal("expected max length error")
	}
	if err := ValidatePositiveInt("price", 0); err == nil {
		t.Fatal("expected positive int error")
	}
	if err := ValidateNonNegativeInt("stock", -1); err == nil {
		t.Fatal("expected non-negative int error")
	}
	if _, err := ParseInt64PathValue(httptest.NewRequest(http.MethodGet, "/", nil), "id"); err == nil {
		t.Fatal("expected missing path value error")
	}
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.SetPathValue("id", "not-a-number")
	if _, err := ParseInt64PathValue(req, "id"); err == nil {
		t.Fatal("expected invalid path value error")
	}
}

func TestWriteJSON(t *testing.T) {
	t.Parallel()

	rr := httptest.NewRecorder()
	WriteJSON(rr, http.StatusCreated, errorResponse{Error: "boom"})

	if rr.Code != http.StatusCreated {
		t.Fatalf("unexpected status: %d", rr.Code)
	}
	if got := rr.Header().Get("Content-Type"); got != "application/json" {
		t.Fatalf("unexpected content type: %s", got)
	}

	var decoded errorResponse
	if err := json.NewDecoder(rr.Body).Decode(&decoded); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if decoded.Error != "boom" {
		t.Fatalf("unexpected response: %#v", decoded)
	}
}

func TestMenuHTTPHandlerEndpoints(t *testing.T) {
	t.Parallel()

	now := time.Now().UTC()
	item := entity.MenuItem{
		ID:                1,
		Name:              "Green Curry",
		Description:       "Thai curry",
		PriceCents:        12900,
		StockQuantity:     5,
		LowStockThreshold: 3,
		IsAvailable:       true,
		IsActive:          true,
		CategoryID:        2,
		CreatedAt:         now,
		UpdatedAt:         now,
	}
	menuRepo := &httpMenuRepository{item: item}
	auditRepo := &httpAuditRepository{}
	priceHistoryRepo := &httpPriceHistoryRepository{}
	menuUC := usecase.NewMenuUsecase(menuRepo, auditRepo, priceHistoryRepo)
	inventoryUC := usecase.NewInventoryUsecase(&httpInventoryRepository{
		settings: entity.InventorySettings{
			ID:                1,
			LowStockThreshold: 5,
			CreatedAt:         now,
			UpdatedAt:         now,
		},
		items: []entity.MenuItem{item},
	})

	handler := NewMenuHTTPHandler(menuUC, inventoryUC)

	rr := httptest.NewRecorder()
	handler.GetCustomerMenu(rr, httptest.NewRequest(http.MethodGet, "/api/menu", nil))
	if rr.Code != http.StatusOK {
		t.Fatalf("unexpected customer menu status: %d", rr.Code)
	}

	rr = httptest.NewRecorder()
	handler.ListManagerMenu(rr, httptest.NewRequest(http.MethodGet, "/api/manager/menu", nil))
	if rr.Code != http.StatusOK {
		t.Fatalf("unexpected manager menu status: %d", rr.Code)
	}

	rr = httptest.NewRecorder()
	handler.CreateMenuItem(rr, httptest.NewRequest(http.MethodPost, "/api/manager/menu", strings.NewReader(`{"name":"Pad Thai","description":"Noodles","price_cents":13900,"stock_quantity":4,"low_stock_threshold":2,"category_id":3,"actor":"manager-1"}`)))
	if rr.Code != http.StatusCreated {
		t.Fatalf("unexpected create status: %d", rr.Code)
	}

	rr = httptest.NewRecorder()
	updateReq := httptest.NewRequest(http.MethodPatch, "/api/manager/menu/1", strings.NewReader(`{"name":"Green Curry Deluxe","description":"Thai curry","price_cents":12900,"stock_quantity":5,"low_stock_threshold":3,"category_id":2,"actor":"manager-1"}`))
	updateReq.SetPathValue("id", "1")
	handler.UpdateMenuItem(rr, updateReq)
	if rr.Code != http.StatusOK {
		t.Fatalf("unexpected update status: %d", rr.Code)
	}

	rr = httptest.NewRecorder()
	deactivateReq := httptest.NewRequest(http.MethodPatch, "/api/manager/menu/1/deactivate", strings.NewReader(`{"actor":"manager-1","reason":"seasonal change"}`))
	deactivateReq.SetPathValue("id", "1")
	handler.DeactivateMenuItem(rr, deactivateReq)
	if rr.Code != http.StatusOK {
		t.Fatalf("unexpected deactivate status: %d", rr.Code)
	}

	rr = httptest.NewRecorder()
	priceReq := httptest.NewRequest(http.MethodPatch, "/api/manager/menu/1/price", strings.NewReader(`{"price_cents":14900,"actor":"manager-1"}`))
	priceReq.SetPathValue("id", "1")
	handler.UpdateMenuItemPrice(rr, priceReq)
	if rr.Code != http.StatusOK {
		t.Fatalf("unexpected price status: %d", rr.Code)
	}

	rr = httptest.NewRecorder()
	stockReq := httptest.NewRequest(http.MethodPatch, "/api/manager/menu/1/stock", strings.NewReader(`{"stock_quantity":0,"actor":"manager-1"}`))
	stockReq.SetPathValue("id", "1")
	handler.UpdateMenuItemStock(rr, stockReq)
	if rr.Code != http.StatusOK {
		t.Fatalf("unexpected stock status: %d", rr.Code)
	}
	if len(priceHistoryRepo.entries) != 1 {
		t.Fatalf("expected one price history entry, got %d", len(priceHistoryRepo.entries))
	}
	if priceHistoryRepo.entries[0].PreviousPriceCents != 12900 || priceHistoryRepo.entries[0].NewPriceCents != 14900 {
		t.Fatalf("unexpected price history entry: %#v", priceHistoryRepo.entries[0])
	}
	if len(auditRepo.entries) != 5 {
		t.Fatalf("expected five audit entries, got %d", len(auditRepo.entries))
	}
	if auditRepo.entries[3].ActionType != entity.AuditActionPrice {
		t.Fatalf("expected price audit entry, got %s", auditRepo.entries[3].ActionType)
	}
	if auditRepo.entries[4].ActionType != entity.AuditActionStock {
		t.Fatalf("expected stock audit entry, got %s", auditRepo.entries[4].ActionType)
	}
}

func TestMenuHTTPHandlerErrorBranches(t *testing.T) {
	t.Parallel()

	handler := NewMenuHTTPHandler(usecase.NewMenuUsecase(nil, nil), nil)

	rr := httptest.NewRecorder()
	handler.GetCustomerMenu(rr, httptest.NewRequest(http.MethodGet, "/api/menu", nil))
	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("expected menu repo error status, got %d", rr.Code)
	}

	rr = httptest.NewRecorder()
	handler.ListManagerMenu(rr, httptest.NewRequest(http.MethodGet, "/api/manager/menu", nil))
	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("expected manager menu error status, got %d", rr.Code)
	}

	rr = httptest.NewRecorder()
	handler.CreateMenuItem(rr, httptest.NewRequest(http.MethodPost, "/api/manager/menu", strings.NewReader(`{"unknown":true}`)))
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected create validation status, got %d", rr.Code)
	}

	rr = httptest.NewRecorder()
	handler.CreateMenuItem(rr, httptest.NewRequest(http.MethodPost, "/api/manager/menu", strings.NewReader(`{"name":"Pad Thai","description":"Noodles","price_cents":13900,"stock_quantity":4,"low_stock_threshold":2,"category_id":3,"actor":""}`)))
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected create actor validation status, got %d", rr.Code)
	}

	rr = httptest.NewRecorder()
	updateReq := httptest.NewRequest(http.MethodPatch, "/api/manager/menu/1", strings.NewReader(`{"name":"ok"}`))
	handler.UpdateMenuItem(rr, updateReq)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected path validation status, got %d", rr.Code)
	}

	rr = httptest.NewRecorder()
	updateReq = httptest.NewRequest(http.MethodPatch, "/api/manager/menu/1", strings.NewReader(`{"name":"Green Curry Deluxe","description":"Thai curry","price_cents":12900,"stock_quantity":5,"low_stock_threshold":3,"category_id":2,"actor":""}`))
	updateReq.SetPathValue("id", "1")
	handler.UpdateMenuItem(rr, updateReq)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected update actor validation status, got %d", rr.Code)
	}

	rr = httptest.NewRecorder()
	priceReq := httptest.NewRequest(http.MethodPatch, "/api/manager/menu/1/price", strings.NewReader(`{"price_cents":0,"actor":"manager-1"}`))
	priceReq.SetPathValue("id", "1")
	handler.UpdateMenuItemPrice(rr, priceReq)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected price validation status, got %d", rr.Code)
	}

	rr = httptest.NewRecorder()
	priceReq = httptest.NewRequest(http.MethodPatch, "/api/manager/menu/1/price", strings.NewReader(`{"price_cents":14900,"actor":""}`))
	priceReq.SetPathValue("id", "1")
	handler.UpdateMenuItemPrice(rr, priceReq)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected price actor validation status, got %d", rr.Code)
	}

	rr = httptest.NewRecorder()
	stockReq := httptest.NewRequest(http.MethodPatch, "/api/manager/menu/1/stock", strings.NewReader(`{"stock_quantity":-1,"actor":"manager-1"}`))
	stockReq.SetPathValue("id", "1")
	handler.UpdateMenuItemStock(rr, stockReq)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected stock validation status, got %d", rr.Code)
	}

	rr = httptest.NewRecorder()
	stockReq = httptest.NewRequest(http.MethodPatch, "/api/manager/menu/1/stock", strings.NewReader(`{"stock_quantity":0,"actor":""}`))
	stockReq.SetPathValue("id", "1")
	handler.UpdateMenuItemStock(rr, stockReq)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected stock actor validation status, got %d", rr.Code)
	}
}

func TestSettingsHTTPHandlerEndpoints(t *testing.T) {
	t.Parallel()

	now := time.Now().UTC()
	handler := NewSettingsHTTPHandler(usecase.NewInventoryUsecase(&httpInventoryRepository{
		settings: entity.InventorySettings{
			ID:                1,
			LowStockThreshold: 5,
			CreatedAt:         now,
			UpdatedAt:         now,
		},
		items: []entity.MenuItem{
			{
				ID:            1,
				Name:          "Tea",
				Description:   "Jasmine tea",
				PriceCents:    2500,
				StockQuantity: 3,
				IsAvailable:   true,
				IsActive:      true,
				CategoryID:    1,
				CreatedAt:     now,
				UpdatedAt:     now,
			},
		},
	}))

	rr := httptest.NewRecorder()
	handler.ListLowStock(rr, httptest.NewRequest(http.MethodGet, "/api/manager/menu/low-stock", nil))
	if rr.Code != http.StatusOK {
		t.Fatalf("unexpected low-stock status: %d", rr.Code)
	}
	var lowStockBody lowStockResponse
	if err := json.NewDecoder(rr.Body).Decode(&lowStockBody); err != nil {
		t.Fatalf("failed to decode low-stock response: %v", err)
	}
	if lowStockBody.Threshold != 5 || len(lowStockBody.Items) != 1 || lowStockBody.Items[0].Name != "Tea" {
		t.Fatalf("unexpected low-stock response: %#v", lowStockBody)
	}

	rr = httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPatch, "/api/manager/settings/low-stock-threshold", strings.NewReader(`{"low_stock_threshold":7,"actor":"manager-1"}`))
	handler.UpdateLowStockThreshold(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("unexpected threshold status: %d", rr.Code)
	}
	var thresholdBody lowStockSettingsResponse
	if err := json.NewDecoder(rr.Body).Decode(&thresholdBody); err != nil {
		t.Fatalf("failed to decode threshold response: %v", err)
	}
	if thresholdBody.Threshold != 7 {
		t.Fatalf("unexpected threshold response: %#v", thresholdBody)
	}
}

func TestSettingsHTTPHandlerErrorBranches(t *testing.T) {
	t.Parallel()

	rr := httptest.NewRecorder()
	NewSettingsHTTPHandler(nil).ListLowStock(rr, httptest.NewRequest(http.MethodGet, "/api/manager/menu/low-stock", nil))
	if rr.Code != http.StatusNotImplemented {
		t.Fatalf("expected not configured status, got %d", rr.Code)
	}

	rr = httptest.NewRecorder()
	NewSettingsHTTPHandler(nil).UpdateLowStockThreshold(rr, httptest.NewRequest(http.MethodPatch, "/api/manager/settings/low-stock-threshold", strings.NewReader(`{"low_stock_threshold":1,"actor":"manager-1"}`)))
	if rr.Code != http.StatusNotImplemented {
		t.Fatalf("expected not configured threshold status, got %d", rr.Code)
	}

	rr = httptest.NewRecorder()
	handler := NewSettingsHTTPHandler(usecase.NewInventoryUsecase(nil))
	handler.ListLowStock(rr, httptest.NewRequest(http.MethodGet, "/api/manager/menu/low-stock", nil))
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected inventory error status, got %d", rr.Code)
	}

	rr = httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPatch, "/api/manager/settings/low-stock-threshold", strings.NewReader(`{"low_stock_threshold":7,"actor":""}`))
	handler.UpdateLowStockThreshold(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected threshold actor validation status, got %d", rr.Code)
	}

	rr = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodPatch, "/api/manager/settings/low-stock-threshold", strings.NewReader(`{"low_stock_threshold":-1,"actor":"manager-1"}`))
	handler.UpdateLowStockThreshold(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected threshold validation status, got %d", rr.Code)
	}
}

func TestNewRouter(t *testing.T) {
	t.Parallel()

	mux := NewRouter(NewMenuHTTPHandler(usecase.NewMenuUsecase(nil, nil), nil))
	if mux == nil {
		t.Fatal("expected router")
	}
}

type httpMenuRepository struct {
	item entity.MenuItem
}

func (r *httpMenuRepository) ListCategories(context.Context) ([]entity.Category, error) {
	return nil, nil
}

func (r *httpMenuRepository) ListMenuItems(context.Context) ([]entity.MenuItem, error) {
	return []entity.MenuItem{r.item}, nil
}

func (r *httpMenuRepository) ListAvailableMenuItems(context.Context) ([]entity.MenuItem, error) {
	return []entity.MenuItem{r.item}, nil
}

func (r *httpMenuRepository) GetMenuItemByID(context.Context, int64) (entity.MenuItem, error) {
	return r.item, nil
}

func (r *httpMenuRepository) CreateMenuItem(_ context.Context, item entity.MenuItem) (entity.MenuItem, error) {
	item.ID = 2
	item.CreatedAt = time.Now().UTC()
	item.UpdatedAt = item.CreatedAt
	item.IsAvailable = item.IsActive && item.StockQuantity > 0
	return item, nil
}

func (r *httpMenuRepository) UpdateMenuItem(_ context.Context, item entity.MenuItem) (entity.MenuItem, error) {
	item.UpdatedAt = time.Now().UTC()
	item.IsAvailable = item.IsActive && item.StockQuantity > 0
	return item, nil
}

func (r *httpMenuRepository) UpdateMenuItemPrice(_ context.Context, _ int64, priceCents int64) (entity.MenuItem, error) {
	r.item.PriceCents = priceCents
	r.item.UpdatedAt = time.Now().UTC()
	return r.item, nil
}

func (r *httpMenuRepository) UpdateMenuItemStock(_ context.Context, _ int64, stockQuantity int) (entity.MenuItem, error) {
	r.item.StockQuantity = stockQuantity
	r.item.IsAvailable = r.item.IsActive && stockQuantity > 0
	r.item.UpdatedAt = time.Now().UTC()
	return r.item, nil
}

func (r *httpMenuRepository) DeactivateMenuItem(context.Context, int64) (entity.MenuItem, error) {
	r.item.IsActive = false
	r.item.IsAvailable = false
	r.item.UpdatedAt = time.Now().UTC()
	return r.item, nil
}

type httpAuditRepository struct {
	entries []entity.MenuItemAuditEntry
}

func (r *httpAuditRepository) RecordMenuItemAuditEntry(_ context.Context, entry entity.MenuItemAuditEntry) error {
	r.entries = append(r.entries, entry)
	return nil
}

type httpPriceHistoryRepository struct {
	entries []entity.PriceHistory
}

func (r *httpPriceHistoryRepository) RecordPriceHistory(_ context.Context, history entity.PriceHistory) error {
	r.entries = append(r.entries, history)
	return nil
}

type httpInventoryRepository struct {
	settings entity.InventorySettings
	items    []entity.MenuItem
}

func (r *httpInventoryRepository) GetInventorySettings(context.Context) (entity.InventorySettings, error) {
	return r.settings, nil
}

func (r *httpInventoryRepository) UpdateInventorySettings(_ context.Context, settings entity.InventorySettings) (entity.InventorySettings, error) {
	r.settings = entity.InventorySettings{
		ID:                1,
		LowStockThreshold: settings.LowStockThreshold,
		CreatedAt:         time.Now().UTC(),
		UpdatedAt:         time.Now().UTC(),
	}
	return r.settings, nil
}

func (r *httpInventoryRepository) UpdateMenuItemStock(context.Context, int64, int) (entity.MenuItem, error) {
	return entity.MenuItem{}, nil
}

func (r *httpInventoryRepository) ListLowStockMenuItems(_ context.Context, threshold int) ([]entity.MenuItem, error) {
	var items []entity.MenuItem
	for _, item := range r.items {
		if item.StockQuantity <= threshold {
			items = append(items, item)
		}
	}
	return items, nil
}
