<template>
  <section class="menu-management-view">
    <ManagerMenuPage />

    <LowStockAlert
      :items="store.lowStockItems"
      :loading="store.isLoading"
      :threshold="store.threshold"
      @refresh="handleRefresh"
    />
  </section>
</template>

<script setup lang="ts">
import { onMounted } from 'vue';
import LowStockAlert from '@/components/LowStockAlert.vue';
import ManagerMenuPage from './ManagerMenuPage.vue';
import { useInventoryStore } from '@/stores/inventoryStore';

const store = useInventoryStore();

const demoItems = [
  {
    id: 1,
    name: 'Green Curry',
    description: 'Thai curry with jasmine rice',
    priceCents: 12900,
    stockQuantity: 3,
    lowStockThreshold: 5,
    isAvailable: true,
    isActive: true,
    categoryId: 1,
    createdAt: '2026-07-05T00:00:00Z',
    updatedAt: '2026-07-05T00:00:00Z',
  },
  {
    id: 2,
    name: 'Pad Thai',
    description: 'Rice noodles with prawns',
    priceCents: 13900,
    stockQuantity: 8,
    lowStockThreshold: 5,
    isAvailable: true,
    isActive: true,
    categoryId: 1,
    createdAt: '2026-07-05T00:00:00Z',
    updatedAt: '2026-07-05T00:00:00Z',
  },
];

function handleRefresh() {
  store.refreshItems(store.items.length > 0 ? store.items : demoItems);
}

onMounted(() => {
  if (!store.hasItems) {
    store.refreshItems(demoItems);
  }
});
</script>
