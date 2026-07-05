<template>
  <section class="low-stock-view" :data-layout-mode="layoutMode">
    <header class="low-stock-view__intro">
      <p class="low-stock-view__eyebrow">Inventory control</p>
      <h2 class="low-stock-view__title">Review stock warnings</h2>
      <p class="low-stock-view__description">
        Keep an eye on items that are close to running out and tune the warning level
        from a touch-friendly workspace.
      </p>
    </header>

    <div class="low-stock-view__status" role="status" aria-live="polite">
      {{ statusMessage }}
    </div>

    <div class="low-stock-view__workspace">
      <LowStockAlert
        :items="store.lowStockItems"
        :loading="store.isLoading"
        :threshold="store.threshold"
        @refresh="handleRefresh"
      />

      <ThresholdControl
        :model-value="store.thresholdDraft"
        @update:model-value="store.updateThresholdDraft"
        @save="handleSaveThreshold"
      />

      <LowStockList :items="store.lowStockItems" :threshold="store.threshold" />
    </div>
  </section>
</template>

<script setup lang="ts">
import { onBeforeUnmount, onMounted, ref } from 'vue';
import LowStockAlert from '@/components/LowStockAlert.vue';
import LowStockList from '@/components/LowStockList.vue';
import ThresholdControl from '@/components/ThresholdControl.vue';
import { useInventoryStore } from '@/stores/inventoryStore';

const store = useInventoryStore();
const statusMessage = ref('Ready to review low-stock items.');
const layoutMode = ref(getLayoutMode());

const demoItems = [
  {
    id: 1,
    name: 'Green Curry',
    description: 'Thai curry with jasmine rice',
    priceCents: 12900,
    stockQuantity: 3,
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
    isAvailable: true,
    isActive: true,
    categoryId: 1,
    createdAt: '2026-07-05T00:00:00Z',
    updatedAt: '2026-07-05T00:00:00Z',
  },
];

function getLayoutMode() {
  const width = window.innerWidth;
  if (width >= 1024) {
    return 'desktop';
  }
  if (width >= 768) {
    return 'tablet';
  }
  return 'mobile';
}

function syncLayoutMode() {
  layoutMode.value = getLayoutMode();
}

function handleSaveThreshold() {
  try {
    const threshold = store.saveThreshold();
    statusMessage.value = `Low-stock threshold set to ${threshold}.`;
  } catch (error) {
    statusMessage.value = error instanceof Error ? error.message : 'Unable to save threshold.';
  }
}

function handleRefresh() {
  store.refreshItems(store.items.length > 0 ? store.items : demoItems);
  statusMessage.value = 'Low-stock warnings refreshed.';
}

onMounted(() => {
  window.addEventListener('resize', syncLayoutMode);
  if (!store.hasItems) {
    store.refreshItems(demoItems);
  }
  store.beginThresholdEdit();
});

onBeforeUnmount(() => {
  window.removeEventListener('resize', syncLayoutMode);
});
</script>
