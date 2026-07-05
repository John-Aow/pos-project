package httpadapter

import "pos-project/backend/entity"

type menuItemResponse struct {
	ID                int64  `json:"id"`
	Name              string `json:"name"`
	Description       string `json:"description"`
	PriceCents        int64  `json:"price_cents"`
	StockQuantity     int    `json:"stock_quantity"`
	LowStockThreshold int    `json:"low_stock_threshold"`
	IsAvailable       bool   `json:"is_available"`
	IsActive          bool   `json:"is_active"`
	CategoryID        int64  `json:"category_id"`
	CreatedAt         string `json:"created_at"`
	UpdatedAt         string `json:"updated_at"`
}

type menuListResponse struct {
	Items []menuItemResponse `json:"items"`
}

type lowStockResponse struct {
	Threshold int                `json:"threshold"`
	Items     []menuItemResponse `json:"items"`
}

type lowStockSettingsResponse struct {
	Threshold int `json:"threshold"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func presentMenuItem(item entity.MenuItem) menuItemResponse {
	return menuItemResponse{
		ID:                item.ID,
		Name:              item.Name,
		Description:       item.Description,
		PriceCents:        item.PriceCents,
		StockQuantity:     item.StockQuantity,
		LowStockThreshold: item.LowStockThreshold,
		IsAvailable:       item.IsAvailable,
		IsActive:          item.IsActive,
		CategoryID:        item.CategoryID,
		CreatedAt:         item.CreatedAt.Format(timeLayout),
		UpdatedAt:         item.UpdatedAt.Format(timeLayout),
	}
}

func presentMenuItems(items []entity.MenuItem) []menuItemResponse {
	responses := make([]menuItemResponse, 0, len(items))
	for _, item := range items {
		responses = append(responses, presentMenuItem(item))
	}
	return responses
}

const timeLayout = "2006-01-02T15:04:05Z07:00"
