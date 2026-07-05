package httpadapter

import (
	"net/http"

	"pos-project/backend/usecase"
)

// MenuHTTPHandler implements the router contract for menu maintenance APIs.
type MenuHTTPHandler struct {
	Usecase *usecase.MenuUsecase
	*SettingsHTTPHandler
}

// NewMenuHTTPHandler constructs a menu HTTP handler.
func NewMenuHTTPHandler(menuUsecase *usecase.MenuUsecase, inventoryUsecase *usecase.InventoryUsecase) *MenuHTTPHandler {
	return &MenuHTTPHandler{
		Usecase:             menuUsecase,
		SettingsHTTPHandler: NewSettingsHTTPHandler(inventoryUsecase),
	}
}

func (h *MenuHTTPHandler) GetCustomerMenu(w http.ResponseWriter, r *http.Request) {
	items, err := h.Usecase.ListCustomerMenu(r.Context())
	if err != nil {
		writeHandlerError(w, http.StatusInternalServerError, err)
		return
	}
	WriteJSON(w, http.StatusOK, menuListResponse{Items: presentMenuItems(items)})
}

func (h *MenuHTTPHandler) ListManagerMenu(w http.ResponseWriter, r *http.Request) {
	items, err := h.Usecase.ListManagerMenu(r.Context())
	if err != nil {
		writeHandlerError(w, http.StatusInternalServerError, err)
		return
	}
	WriteJSON(w, http.StatusOK, menuListResponse{Items: presentMenuItems(items)})
}

func (h *MenuHTTPHandler) CreateMenuItem(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name              string `json:"name"`
		Description       string `json:"description"`
		PriceCents        int64  `json:"price_cents"`
		StockQuantity     int    `json:"stock_quantity"`
		LowStockThreshold int    `json:"low_stock_threshold"`
		CategoryID        int64  `json:"category_id"`
		Actor             string `json:"actor"`
		Reason            string `json:"reason"`
	}
	if err := DecodeJSONBody(r, &input); err != nil {
		writeHandlerError(w, http.StatusBadRequest, err)
		return
	}
	if err := ValidateMenuItemInput(input.Name, input.Description, input.PriceCents, input.StockQuantity, input.LowStockThreshold, input.CategoryID); err != nil {
		writeHandlerError(w, http.StatusBadRequest, err)
		return
	}
	if err := ValidateRequiredString("actor", input.Actor, 1, 120); err != nil {
		writeHandlerError(w, http.StatusBadRequest, err)
		return
	}

	item, err := h.Usecase.CreateMenuItem(r.Context(), usecase.CreateMenuItemInput{
		Name:              input.Name,
		Description:       input.Description,
		PriceCents:        input.PriceCents,
		StockQuantity:     input.StockQuantity,
		LowStockThreshold: input.LowStockThreshold,
		CategoryID:        input.CategoryID,
		Actor:             input.Actor,
		Reason:            input.Reason,
	})
	if err != nil {
		writeHandlerError(w, http.StatusBadRequest, err)
		return
	}

	WriteJSON(w, http.StatusCreated, presentMenuItem(item))
}

func (h *MenuHTTPHandler) UpdateMenuItem(w http.ResponseWriter, r *http.Request) {
	id, err := ParseInt64PathValue(r, "id")
	if err != nil {
		writeHandlerError(w, http.StatusBadRequest, err)
		return
	}

	var input struct {
		Name              string `json:"name"`
		Description       string `json:"description"`
		PriceCents        int64  `json:"price_cents"`
		StockQuantity     int    `json:"stock_quantity"`
		LowStockThreshold int    `json:"low_stock_threshold"`
		CategoryID        int64  `json:"category_id"`
		IsActive          bool   `json:"is_active"`
		Actor             string `json:"actor"`
		Reason            string `json:"reason"`
	}
	if err := DecodeJSONBody(r, &input); err != nil {
		writeHandlerError(w, http.StatusBadRequest, err)
		return
	}

	if err := ValidateMenuItemInput(input.Name, input.Description, input.PriceCents, input.StockQuantity, input.LowStockThreshold, input.CategoryID); err != nil {
		writeHandlerError(w, http.StatusBadRequest, err)
		return
	}
	if err := ValidateRequiredString("actor", input.Actor, 1, 120); err != nil {
		writeHandlerError(w, http.StatusBadRequest, err)
		return
	}

	item, err := h.Usecase.UpdateMenuItem(r.Context(), usecase.UpdateMenuItemInput{
		ID:                id,
		Name:              input.Name,
		Description:       input.Description,
		PriceCents:        input.PriceCents,
		StockQuantity:     input.StockQuantity,
		LowStockThreshold: input.LowStockThreshold,
		CategoryID:        input.CategoryID,
		IsActive:          input.IsActive,
		Actor:             input.Actor,
		Reason:            input.Reason,
	})
	if err != nil {
		writeHandlerError(w, http.StatusBadRequest, err)
		return
	}

	WriteJSON(w, http.StatusOK, presentMenuItem(item))
}

func (h *MenuHTTPHandler) DeactivateMenuItem(w http.ResponseWriter, r *http.Request) {
	id, err := ParseInt64PathValue(r, "id")
	if err != nil {
		writeHandlerError(w, http.StatusBadRequest, err)
		return
	}

	var input struct {
		Actor  string `json:"actor"`
		Reason string `json:"reason"`
	}
	if err := DecodeJSONBody(r, &input); err != nil {
		writeHandlerError(w, http.StatusBadRequest, err)
		return
	}
	if err := ValidateRequiredString("actor", input.Actor, 1, 120); err != nil {
		writeHandlerError(w, http.StatusBadRequest, err)
		return
	}
	if err := ValidateRequiredString("reason", input.Reason, 1, 500); err != nil {
		writeHandlerError(w, http.StatusBadRequest, err)
		return
	}

	item, err := h.Usecase.DeactivateMenuItem(r.Context(), id, input.Actor, input.Reason)
	if err != nil {
		writeHandlerError(w, http.StatusBadRequest, err)
		return
	}

	WriteJSON(w, http.StatusOK, presentMenuItem(item))
}

func (h *MenuHTTPHandler) UpdateMenuItemPrice(w http.ResponseWriter, r *http.Request) {
	id, err := ParseInt64PathValue(r, "id")
	if err != nil {
		writeHandlerError(w, http.StatusBadRequest, err)
		return
	}

	var input struct {
		PriceCents int64  `json:"price_cents"`
		Actor      string `json:"actor"`
		Reason     string `json:"reason"`
	}
	if err := DecodeJSONBody(r, &input); err != nil {
		writeHandlerError(w, http.StatusBadRequest, err)
		return
	}
	if err := ValidatePriceUpdateInput(input.PriceCents, input.Actor); err != nil {
		writeHandlerError(w, http.StatusBadRequest, err)
		return
	}

	item, err := h.Usecase.UpdateMenuItemPrice(r.Context(), usecase.UpdateMenuItemPriceInput{
		ID:         id,
		PriceCents: input.PriceCents,
		Actor:      input.Actor,
		Reason:     input.Reason,
	})
	if err != nil {
		writeHandlerError(w, http.StatusBadRequest, err)
		return
	}

	WriteJSON(w, http.StatusOK, presentMenuItem(item))
}

func (h *MenuHTTPHandler) UpdateMenuItemStock(w http.ResponseWriter, r *http.Request) {
	id, err := ParseInt64PathValue(r, "id")
	if err != nil {
		writeHandlerError(w, http.StatusBadRequest, err)
		return
	}

	var input struct {
		StockQuantity int    `json:"stock_quantity"`
		Actor         string `json:"actor"`
		Reason        string `json:"reason"`
	}
	if err := DecodeJSONBody(r, &input); err != nil {
		writeHandlerError(w, http.StatusBadRequest, err)
		return
	}
	if err := ValidateStockUpdateInput(input.StockQuantity, input.Actor); err != nil {
		writeHandlerError(w, http.StatusBadRequest, err)
		return
	}

	item, err := h.Usecase.UpdateMenuItemStock(r.Context(), usecase.UpdateMenuItemStockInput{
		ID:            id,
		StockQuantity: input.StockQuantity,
		Actor:         input.Actor,
		Reason:        input.Reason,
	})
	if err != nil {
		writeHandlerError(w, http.StatusBadRequest, err)
		return
	}

	WriteJSON(w, http.StatusOK, presentMenuItem(item))
}

func writeHandlerError(w http.ResponseWriter, status int, err error) {
	WriteJSON(w, status, errorResponse{Error: err.Error()})
}
