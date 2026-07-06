package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"sort"
	"sync"
	"time"

	"pos-project/backend/entity"
	httpadapter "pos-project/backend/interface/adapter/http"
	"pos-project/backend/usecase"
)

func main() {
	port := getEnv("PORT", "8080")

	store := newMemoryStore()
	menuRepo := &memoryMenuRepository{store: store}
	inventoryRepo := &memoryInventoryRepository{store: store}
	auditRepo := &memoryAuditRepository{store: store}
	priceHistoryRepo := &memoryPriceHistoryRepository{store: store}
	auditLogRepo := &memoryAuditLogRepository{store: store}

	menuUsecase := usecase.NewMenuUsecase(menuRepo, auditRepo, priceHistoryRepo, auditLogRepo)
	inventoryUsecase := usecase.NewInventoryUsecase(inventoryRepo)

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	mux.Handle("/", httpadapter.NewRouter(httpadapter.NewMenuHTTPHandler(menuUsecase, inventoryUsecase)))

	server := &http.Server{
		Addr:              ":" + port,
		Handler:           withCORS(mux),
		ReadHeaderTimeout: 5 * time.Second,
	}

	fmt.Printf("backend listening on :%s\n", port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		panic(err)
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PATCH,OPTIONS")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

type memoryStore struct {
	mu sync.RWMutex

	categories map[int64]entity.Category
	items      map[int64]entity.MenuItem
	settings   entity.InventorySettings

	menuItemAudits   []entity.MenuItemAuditEntry
	priceHistory     []entity.PriceHistory
	menuAuditEntries []entity.MenuAuditLog

	nextItemID int64
}

func newMemoryStore() *memoryStore {
	now := time.Now().UTC()
	store := &memoryStore{
		categories: map[int64]entity.Category{
			1: {ID: 1, Name: "Mains", SortOrder: 1, CreatedAt: now, UpdatedAt: now},
			2: {ID: 2, Name: "Drinks", SortOrder: 2, CreatedAt: now, UpdatedAt: now},
		},
		items: map[int64]entity.MenuItem{},
		settings: entity.InventorySettings{
			ID:                1,
			LowStockThreshold: 5,
			CreatedAt:         now,
			UpdatedAt:         now,
		},
		nextItemID: 1,
	}

	store.seedItems([]entity.MenuItem{
		{
			Name:              "Green Curry",
			Description:       "Thai curry with jasmine rice",
			PriceCents:        12900,
			StockQuantity:     3,
			LowStockThreshold: 5,
			IsActive:          true,
			CategoryID:        1,
		},
		{
			Name:              "Pad Thai",
			Description:       "Rice noodles with prawns",
			PriceCents:        13900,
			StockQuantity:     8,
			LowStockThreshold: 5,
			IsActive:          true,
			CategoryID:        1,
		},
		{
			Name:              "Thai Iced Tea",
			Description:       "Sweet tea with milk",
			PriceCents:        4500,
			StockQuantity:     0,
			LowStockThreshold: 5,
			IsActive:          true,
			CategoryID:        2,
		},
	})

	return store
}

func (s *memoryStore) seedItems(items []entity.MenuItem) {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now().UTC()
	for _, item := range items {
		item.ID = s.nextItemID
		s.nextItemID++
		item.CreatedAt = now
		item.UpdatedAt = now
		item.RecalculateAvailability()
		s.items[item.ID] = item
	}
}

func (s *memoryStore) listItems() []entity.MenuItem {
	items := make([]entity.MenuItem, 0, len(s.items))
	for _, item := range s.items {
		items = append(items, item)
	}
	sort.Slice(items, func(i, j int) bool {
		if !items[i].CreatedAt.Equal(items[j].CreatedAt) {
			return items[i].CreatedAt.After(items[j].CreatedAt)
		}
		return items[i].ID > items[j].ID
	})
	return items
}

type memoryMenuRepository struct {
	store *memoryStore
}

func (r *memoryMenuRepository) ListCategories(context.Context) ([]entity.Category, error) {
	r.store.mu.RLock()
	defer r.store.mu.RUnlock()

	categories := make([]entity.Category, 0, len(r.store.categories))
	for _, category := range r.store.categories {
		categories = append(categories, category)
	}
	sort.Slice(categories, func(i, j int) bool {
		if categories[i].SortOrder != categories[j].SortOrder {
			return categories[i].SortOrder < categories[j].SortOrder
		}
		return categories[i].Name < categories[j].Name
	})
	return categories, nil
}

func (r *memoryMenuRepository) ListMenuItems(context.Context) ([]entity.MenuItem, error) {
	r.store.mu.RLock()
	defer r.store.mu.RUnlock()

	return r.store.listItems(), nil
}

func (r *memoryMenuRepository) ListAvailableMenuItems(context.Context) ([]entity.MenuItem, error) {
	r.store.mu.RLock()
	defer r.store.mu.RUnlock()

	items := make([]entity.MenuItem, 0, len(r.store.items))
	for _, item := range r.store.items {
		if item.CustomerOrderable() {
			items = append(items, item)
		}
	}
	sort.Slice(items, func(i, j int) bool {
		return items[i].Name < items[j].Name
	})
	return items, nil
}

func (r *memoryMenuRepository) GetMenuItemByID(_ context.Context, id int64) (entity.MenuItem, error) {
	r.store.mu.RLock()
	defer r.store.mu.RUnlock()

	item, ok := r.store.items[id]
	if !ok {
		return entity.MenuItem{}, fmt.Errorf("menu item not found")
	}
	return item, nil
}

func (r *memoryMenuRepository) CreateMenuItem(_ context.Context, item entity.MenuItem) (entity.MenuItem, error) {
	r.store.mu.Lock()
	defer r.store.mu.Unlock()

	item.ID = r.store.nextItemID
	r.store.nextItemID++
	item.CreatedAt = time.Now().UTC()
	item.UpdatedAt = item.CreatedAt
	item.RecalculateAvailability()
	r.store.items[item.ID] = item
	return item, nil
}

func (r *memoryMenuRepository) UpdateMenuItem(_ context.Context, item entity.MenuItem) (entity.MenuItem, error) {
	r.store.mu.Lock()
	defer r.store.mu.Unlock()

	existing, ok := r.store.items[item.ID]
	if !ok {
		return entity.MenuItem{}, fmt.Errorf("menu item not found")
	}

	item.CreatedAt = existing.CreatedAt
	item.UpdatedAt = time.Now().UTC()
	item.RecalculateAvailability()
	r.store.items[item.ID] = item
	return item, nil
}

func (r *memoryMenuRepository) UpdateMenuItemPrice(_ context.Context, id int64, priceCents int64) (entity.MenuItem, error) {
	r.store.mu.Lock()
	defer r.store.mu.Unlock()

	item, ok := r.store.items[id]
	if !ok {
		return entity.MenuItem{}, fmt.Errorf("menu item not found")
	}
	item.PriceCents = priceCents
	item.UpdatedAt = time.Now().UTC()
	r.store.items[id] = item
	return item, nil
}

func (r *memoryMenuRepository) UpdateMenuItemStock(_ context.Context, id int64, stockQuantity int) (entity.MenuItem, error) {
	r.store.mu.Lock()
	defer r.store.mu.Unlock()

	item, ok := r.store.items[id]
	if !ok {
		return entity.MenuItem{}, fmt.Errorf("menu item not found")
	}
	item.StockQuantity = stockQuantity
	item.RecalculateAvailability()
	item.UpdatedAt = time.Now().UTC()
	r.store.items[id] = item
	return item, nil
}

func (r *memoryMenuRepository) DeactivateMenuItem(_ context.Context, id int64) (entity.MenuItem, error) {
	r.store.mu.Lock()
	defer r.store.mu.Unlock()

	item, ok := r.store.items[id]
	if !ok {
		return entity.MenuItem{}, fmt.Errorf("menu item not found")
	}
	item.IsActive = false
	item.RecalculateAvailability()
	item.UpdatedAt = time.Now().UTC()
	r.store.items[id] = item
	return item, nil
}

type memoryInventoryRepository struct {
	store *memoryStore
}

func (r *memoryInventoryRepository) GetInventorySettings(context.Context) (entity.InventorySettings, error) {
	r.store.mu.RLock()
	defer r.store.mu.RUnlock()

	return r.store.settings, nil
}

func (r *memoryInventoryRepository) UpdateInventorySettings(_ context.Context, settings entity.InventorySettings) (entity.InventorySettings, error) {
	r.store.mu.Lock()
	defer r.store.mu.Unlock()

	r.store.settings = entity.InventorySettings{
		ID:                1,
		LowStockThreshold: settings.LowStockThreshold,
		CreatedAt:         r.store.settings.CreatedAt,
		UpdatedAt:         time.Now().UTC(),
	}
	return r.store.settings, nil
}

func (r *memoryInventoryRepository) UpdateMenuItemStock(ctx context.Context, id int64, stockQuantity int) (entity.MenuItem, error) {
	return (&memoryMenuRepository{store: r.store}).UpdateMenuItemStock(ctx, id, stockQuantity)
}

func (r *memoryInventoryRepository) ListLowStockMenuItems(_ context.Context, threshold int) ([]entity.MenuItem, error) {
	r.store.mu.RLock()
	defer r.store.mu.RUnlock()

	items := make([]entity.MenuItem, 0)
	for _, item := range r.store.items {
		if item.IsActive && item.StockQuantity <= threshold {
			items = append(items, item)
		}
	}
	sort.Slice(items, func(i, j int) bool {
		if items[i].StockQuantity != items[j].StockQuantity {
			return items[i].StockQuantity < items[j].StockQuantity
		}
		return items[i].Name < items[j].Name
	})
	return items, nil
}

type memoryAuditRepository struct {
	store *memoryStore
}

func (r *memoryAuditRepository) RecordMenuItemAuditEntry(_ context.Context, entry entity.MenuItemAuditEntry) error {
	r.store.mu.Lock()
	defer r.store.mu.Unlock()

	entry.ID = int64(len(r.store.menuItemAudits) + 1)
	if entry.CreatedAt.IsZero() {
		entry.CreatedAt = time.Now().UTC()
	}
	r.store.menuItemAudits = append(r.store.menuItemAudits, entry)
	return nil
}

type memoryPriceHistoryRepository struct {
	store *memoryStore
}

func (r *memoryPriceHistoryRepository) RecordPriceHistory(_ context.Context, history entity.PriceHistory) error {
	r.store.mu.Lock()
	defer r.store.mu.Unlock()

	history.ID = int64(len(r.store.priceHistory) + 1)
	if history.CreatedAt.IsZero() {
		history.CreatedAt = time.Now().UTC()
	}
	r.store.priceHistory = append(r.store.priceHistory, history)
	return nil
}

type memoryAuditLogRepository struct {
	store *memoryStore
}

func (r *memoryAuditLogRepository) RecordMenuAuditLog(_ context.Context, log entity.MenuAuditLog) error {
	r.store.mu.Lock()
	defer r.store.mu.Unlock()

	log.ID = int64(len(r.store.menuAuditEntries) + 1)
	if log.CreatedAt.IsZero() {
		log.CreatedAt = time.Now().UTC()
	}
	r.store.menuAuditEntries = append(r.store.menuAuditEntries, log)
	return nil
}
