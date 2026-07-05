package httpadapter

import "net/http"

// MenuHandler describes the HTTP handlers used by the manager menu feature.
type MenuHandler interface {
	GetCustomerMenu(http.ResponseWriter, *http.Request)
	ListManagerMenu(http.ResponseWriter, *http.Request)
	CreateMenuItem(http.ResponseWriter, *http.Request)
	UpdateMenuItem(http.ResponseWriter, *http.Request)
	UpdateMenuItemPrice(http.ResponseWriter, *http.Request)
	UpdateMenuItemStock(http.ResponseWriter, *http.Request)
	ListLowStock(http.ResponseWriter, *http.Request)
	UpdateLowStockThreshold(http.ResponseWriter, *http.Request)
}

// NewRouter wires the manager and customer menu routes.
func NewRouter(handler MenuHandler) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/menu", handler.GetCustomerMenu)
	mux.HandleFunc("GET /api/manager/menu", handler.ListManagerMenu)
	mux.HandleFunc("POST /api/manager/menu", handler.CreateMenuItem)
	mux.HandleFunc("PATCH /api/manager/menu/{id}", handler.UpdateMenuItem)
	mux.HandleFunc("PATCH /api/manager/menu/{id}/price", handler.UpdateMenuItemPrice)
	mux.HandleFunc("PATCH /api/manager/menu/{id}/stock", handler.UpdateMenuItemStock)
	mux.HandleFunc("GET /api/manager/menu/low-stock", handler.ListLowStock)
	mux.HandleFunc("PATCH /api/manager/settings/low-stock-threshold", handler.UpdateLowStockThreshold)

	return mux
}
