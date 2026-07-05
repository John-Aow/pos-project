import { defineStore } from 'pinia';

export interface MenuItem {
  id: number;
  name: string;
  description: string;
  priceCents: number;
  stockQuantity: number;
  lowStockThreshold: number;
  isAvailable: boolean;
  isActive: boolean;
  categoryId: number;
  createdAt: string;
  updatedAt: string;
}

export interface MenuItemDraft {
  name: string;
  description: string;
  priceCents: number;
  stockQuantity: number;
  lowStockThreshold: number;
  categoryId: number;
}

const now = () => new Date().toISOString();

const emptyDraft = (): MenuItemDraft => ({
  name: '',
  description: '',
  priceCents: 0,
  stockQuantity: 0,
  lowStockThreshold: 5,
  categoryId: 1,
});

const normalizeItem = (item: MenuItem): MenuItem => ({
  ...item,
  lowStockThreshold: item.lowStockThreshold ?? 5,
  isAvailable: item.isActive && item.stockQuantity > 0,
});

export const useManagerMenuStore = defineStore('managerMenu', {
  state: () => ({
    items: [] as MenuItem[],
    draft: emptyDraft(),
    editingId: null as number | null,
    selectedItemId: null as number | null,
    editorMode: 'menu' as 'menu' | 'price' | 'stock',
    priceDraft: 0,
    stockDraft: 0,
    nextId: 1,
  }),
  getters: {
    hasItems: (state) => state.items.length > 0,
    activeItems: (state) => state.items.filter((item) => item.isActive),
    customerItems: (state) =>
      state.items.filter((item) => item.isActive && item.stockQuantity > 0),
    isEditing: (state) => state.editingId !== null,
    selectedItem: (state) =>
      state.items.find((item) => item.id === state.selectedItemId) ?? null,
  },
  actions: {
    seedItems(items: MenuItem[]) {
      this.items = items.map(normalizeItem);
      const largestId = this.items.reduce((max, item) => Math.max(max, item.id), 0);
      this.nextId = largestId + 1;
    },
    beginCreate() {
      this.editingId = null;
      this.draft = emptyDraft();
    },
    clearSelection() {
      this.selectedItemId = null;
      this.editorMode = 'menu';
      this.priceDraft = 0;
      this.stockDraft = 0;
    },
    beginEdit(id: number) {
      const item = this.items.find((candidate) => candidate.id === id);
      if (!item) {
        return;
      }

      this.editingId = id;
      this.draft = {
        name: item.name,
        description: item.description,
        priceCents: item.priceCents,
        stockQuantity: item.stockQuantity,
        lowStockThreshold: item.lowStockThreshold,
        categoryId: item.categoryId,
      };
    },
    beginPriceEdit(id: number) {
      const item = this.items.find((candidate) => candidate.id === id);
      if (!item) {
        return;
      }

      this.selectedItemId = id;
      this.editorMode = 'price';
      this.priceDraft = item.priceCents;
    },
    beginStockEdit(id: number) {
      const item = this.items.find((candidate) => candidate.id === id);
      if (!item) {
        return;
      }

      this.selectedItemId = id;
      this.editorMode = 'stock';
      this.stockDraft = item.stockQuantity;
    },
    updateDraft(partial: Partial<MenuItemDraft>) {
      this.draft = { ...this.draft, ...partial };
    },
    updatePriceDraft(value: number) {
      this.priceDraft = value;
    },
    updateStockDraft(value: number) {
      this.stockDraft = value;
    },
    saveDraft(): MenuItem {
      if (this.draft.name.trim() === '') {
        throw new Error('Menu item name is required');
      }
      if (this.draft.description.trim() === '') {
        throw new Error('Menu item description is required');
      }
      if (this.draft.priceCents <= 0) {
        throw new Error('Menu item price must be greater than zero');
      }
      if (this.draft.stockQuantity < 0) {
        throw new Error('Menu item stock must be zero or greater');
      }
      if (this.draft.lowStockThreshold < 0) {
        throw new Error('Menu item low-stock threshold must be zero or greater');
      }
      if (this.draft.categoryId <= 0) {
        throw new Error('Menu item category is required');
      }

      const timestamp = now();
      const existing = this.items.find((item) => item.id === this.editingId);
      const nextItem: MenuItem = {
        id: this.editingId ?? this.nextId,
        name: this.draft.name.trim(),
        description: this.draft.description.trim(),
        priceCents: this.draft.priceCents,
        stockQuantity: this.draft.stockQuantity,
        lowStockThreshold: this.draft.lowStockThreshold,
        isActive: existing?.isActive ?? true,
        isAvailable: (existing?.isActive ?? true) && this.draft.stockQuantity > 0,
        categoryId: this.draft.categoryId,
        createdAt: timestamp,
        updatedAt: timestamp,
      };

      if (this.editingId === null) {
        this.items = [...this.items, nextItem];
        this.nextId += 1;
      } else {
        this.items = this.items.map((item) =>
          item.id === this.editingId
            ? {
                ...item,
                ...nextItem,
                id: item.id,
                createdAt: item.createdAt,
              }
            : item,
        );
      }

      this.beginCreate();
      return nextItem;
    },
    deactivateItem(id: number) {
      this.items = this.items.map((item) =>
        item.id === id
          ? {
              ...item,
              isActive: false,
              isAvailable: false,
              updatedAt: now(),
            }
          : item,
      );
      if (this.editingId === id) {
        this.beginCreate();
      }
      if (this.selectedItemId === id) {
        this.clearSelection();
      }
    },
    savePriceDraft(): MenuItem {
      if (this.selectedItemId === null) {
        throw new Error('Select a menu item before updating price');
      }
      if (this.priceDraft <= 0) {
        throw new Error('Menu item price must be greater than zero');
      }

      const timestamp = now();
      let updatedItem: MenuItem | null = null;
      this.items = this.items.map((item) =>
        item.id === this.selectedItemId
          ? (() => {
              updatedItem = {
                ...item,
                priceCents: this.priceDraft,
                updatedAt: timestamp,
              };
              return updatedItem;
            })()
          : item,
      );

      if (!updatedItem) {
        throw new Error('Selected menu item no longer exists');
      }

      return updatedItem;
    },
    saveStockDraft(): MenuItem {
      if (this.selectedItemId === null) {
        throw new Error('Select a menu item before updating stock');
      }
      if (this.stockDraft < 0) {
        throw new Error('Menu item stock must be zero or greater');
      }

      const timestamp = now();
      let updatedItem: MenuItem | null = null;
      this.items = this.items.map((item) =>
        item.id === this.selectedItemId
          ? (() => {
              updatedItem = {
                ...item,
                stockQuantity: this.stockDraft,
                isAvailable: item.isActive && this.stockDraft > 0,
                updatedAt: timestamp,
              };
              return updatedItem;
            })()
          : item,
      );

      if (!updatedItem) {
        throw new Error('Selected menu item no longer exists');
      }

      return updatedItem;
    },
  },
});
