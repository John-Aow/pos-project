import { defineStore } from 'pinia';

export interface InventoryItem {
  id: number;
  name: string;
  description: string;
  priceCents: number;
  stockQuantity: number;
  isAvailable: boolean;
  isActive: boolean;
  categoryId: number;
  createdAt: string;
  updatedAt: string;
}

const now = () => new Date().toISOString();

const normalizeItem = (item: InventoryItem): InventoryItem => ({
  ...item,
  isAvailable: item.isActive && item.stockQuantity > 0,
});

export const useInventoryStore = defineStore('inventory', {
  state: () => ({
    items: [] as InventoryItem[],
    threshold: 5,
    thresholdDraft: 5,
    isLoading: false,
    lastRefreshedAt: '' as string,
  }),
  getters: {
    hasItems: (state) => state.items.length > 0,
    lowStockItems: (state) =>
      [...state.items]
        .filter((item) => item.isActive && item.stockQuantity <= state.threshold)
        .sort((left, right) => {
          if (left.stockQuantity !== right.stockQuantity) {
            return left.stockQuantity - right.stockQuantity;
          }
          return left.name.localeCompare(right.name);
        }),
  },
  actions: {
    seedItems(items: InventoryItem[]) {
      this.items = items.map(normalizeItem);
      this.thresholdDraft = this.threshold;
      this.isLoading = false;
      this.lastRefreshedAt = now();
    },
    beginRefresh() {
      this.isLoading = true;
    },
    finishRefresh() {
      this.isLoading = false;
      this.lastRefreshedAt = now();
    },
    refreshItems(items: InventoryItem[]) {
      this.beginRefresh();
      this.items = items.map(normalizeItem);
      this.thresholdDraft = this.threshold;
      this.finishRefresh();
    },
    beginThresholdEdit() {
      this.thresholdDraft = this.threshold;
    },
    updateThresholdDraft(value: number) {
      this.thresholdDraft = value;
    },
    saveThreshold() {
      if (this.thresholdDraft < 0) {
        throw new Error('Low-stock threshold must be zero or greater');
      }

      this.threshold = this.thresholdDraft;
      return this.threshold;
    },
    updateItemStock(id: number, stockQuantity: number) {
      this.items = this.items.map((item) =>
        item.id === id
          ? {
              ...item,
              stockQuantity,
              isAvailable: item.isActive && stockQuantity > 0,
              updatedAt: now(),
            }
          : item,
      );
    },
  },
});
