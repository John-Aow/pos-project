package httpadapter

import (
	"net/http"

	"pos-project/backend/usecase"
)

// SettingsHTTPHandler serves the low-stock warning and threshold endpoints.
type SettingsHTTPHandler struct {
	Usecase *usecase.InventoryUsecase
}

// NewSettingsHTTPHandler constructs a settings handler.
func NewSettingsHTTPHandler(uc *usecase.InventoryUsecase) *SettingsHTTPHandler {
	return &SettingsHTTPHandler{Usecase: uc}
}

func (h *SettingsHTTPHandler) ListLowStock(w http.ResponseWriter, r *http.Request) {
	if h == nil || h.Usecase == nil {
		writeHandlerError(w, http.StatusNotImplemented, errHandlerNotConfigured)
		return
	}

	overview, err := h.Usecase.GetLowStockItems(r.Context())
	if err != nil {
		writeHandlerError(w, http.StatusBadRequest, err)
		return
	}

	WriteJSON(w, http.StatusOK, lowStockResponse{
		Threshold: overview.Settings.LowStockThreshold,
		Items:     presentMenuItems(overview.Items),
	})
}

func (h *SettingsHTTPHandler) UpdateLowStockThreshold(w http.ResponseWriter, r *http.Request) {
	if h == nil || h.Usecase == nil {
		writeHandlerError(w, http.StatusNotImplemented, errHandlerNotConfigured)
		return
	}

	var input struct {
		LowStockThreshold int    `json:"low_stock_threshold"`
		Actor             string `json:"actor"`
		Reason            string `json:"reason"`
	}
	if err := DecodeJSONBody(r, &input); err != nil {
		writeHandlerError(w, http.StatusBadRequest, err)
		return
	}
	if err := ValidateNonNegativeInt("low_stock_threshold", input.LowStockThreshold); err != nil {
		writeHandlerError(w, http.StatusBadRequest, err)
		return
	}
	if err := ValidateRequiredString("actor", input.Actor, 1, 120); err != nil {
		writeHandlerError(w, http.StatusBadRequest, err)
		return
	}

	settings, err := h.Usecase.UpdateLowStockThreshold(r.Context(), usecase.UpdateLowStockThresholdInput{
		LowStockThreshold: input.LowStockThreshold,
		Actor:             input.Actor,
		Reason:            input.Reason,
	})
	if err != nil {
		writeHandlerError(w, http.StatusBadRequest, err)
		return
	}

	WriteJSON(w, http.StatusOK, lowStockSettingsResponse{
		Threshold: settings.LowStockThreshold,
	})
}

var errHandlerNotConfigured = &handlerError{"handler not configured"}

type handlerError struct {
	message string
}

func (e *handlerError) Error() string {
	return e.message
}
