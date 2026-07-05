<template>
  <section class="low-stock-list" aria-labelledby="low-stock-list-title">
    <header class="low-stock-list__header">
      <p class="manager-panel__label">Low-stock items</p>
      <h3 id="low-stock-list-title" class="manager-panel__title">Items at or below threshold</h3>
      <p class="low-stock-list__summary">Current warning level: {{ threshold }}</p>
    </header>

    <p v-if="items.length === 0" class="low-stock-list__empty">
      No items are currently below the warning level.
    </p>

    <table v-else class="low-stock-list__table">
      <thead>
        <tr>
          <th scope="col">Item</th>
          <th scope="col">Stock</th>
          <th scope="col">Status</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="item in items" :key="item.id">
          <td>
            <strong>{{ item.name }}</strong>
            <p class="low-stock-list__description">{{ item.description }}</p>
          </td>
          <td>{{ item.stockQuantity }}</td>
          <td>
            <span
              class="low-stock-list__badge"
              :class="{ 'is-empty': item.stockQuantity === 0 }"
            >
              {{ item.stockQuantity === 0 ? 'Out of stock' : 'Low stock' }}
            </span>
          </td>
        </tr>
      </tbody>
    </table>
  </section>
</template>

<script setup lang="ts">
import type { InventoryItem } from '@/stores/inventoryStore';

defineProps<{
  items: InventoryItem[];
  threshold: number;
}>();
</script>
