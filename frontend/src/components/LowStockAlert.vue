<template>
  <section class="low-stock-alert" aria-labelledby="low-stock-alert-title">
    <header class="low-stock-alert__header">
      <div>
        <p class="manager-panel__label">Stock warnings</p>
        <h3 id="low-stock-alert-title" class="manager-panel__title">Low-stock alert</h3>
      </div>
      <button
        type="button"
        class="low-stock-alert__refresh"
        data-testid="low-stock-alert-refresh"
        @click="emit('refresh')"
      >
        Refresh
      </button>
    </header>

    <p class="low-stock-alert__summary">
      {{ loading ? 'Refreshing low-stock items.' : summaryText }}
    </p>

    <p v-if="loading" class="low-stock-alert__empty">Loading warning data.</p>
    <p v-else-if="items.length === 0" class="low-stock-alert__empty">
      No items are currently below the warning level.
    </p>
    <ul v-else class="low-stock-alert__list">
      <li v-for="item in visibleItems" :key="item.id" class="low-stock-alert__item">
        <strong>{{ item.name }}</strong>
        <span>{{ item.stockQuantity }} in stock</span>
      </li>
    </ul>
  </section>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import type { InventoryItem } from '@/stores/inventoryStore';

const props = defineProps<{
  items: InventoryItem[];
  threshold: number;
  loading?: boolean;
}>();

const emit = defineEmits(['refresh']);

const visibleItems = computed(() => props.items.slice(0, 4));

const summaryText = computed(() => {
  const count = props.items.length;
  const noun = count === 1 ? 'item' : 'items';
  return `${count} ${noun} at or below ${props.threshold}.`;
});
</script>
