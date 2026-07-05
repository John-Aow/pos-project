<template>
  <div class="menu-item-table">
    <table>
      <thead>
        <tr>
          <th>Name</th>
          <th>Price</th>
          <th>Stock</th>
          <th>Threshold</th>
          <th>Status</th>
          <th>Actions</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="item in items" :key="item.id">
          <td data-label="Item">
            <strong>{{ item.name }}</strong>
            <p class="menu-item-table__description">{{ item.description }}</p>
          </td>
          <td data-label="Price">{{ formatPrice(item.priceCents) }}</td>
          <td data-label="Stock">{{ item.stockQuantity }}</td>
          <td data-label="Threshold">{{ item.lowStockThreshold }}</td>
          <td data-label="Status">
            <span
              class="menu-item-table__status"
              :class="item.isActive && item.isAvailable ? 'is-available' : 'is-inactive'"
            >
              {{ item.isActive && item.isAvailable ? 'Available' : 'Inactive' }}
            </span>
          </td>
          <td class="menu-item-table__actions" data-label="Actions">
            <button
              type="button"
              :data-testid="`edit-item-${item.id}`"
              @click="emit('edit', item.id)"
            >
              Edit
            </button>
            <button
              type="button"
              :data-testid="`price-item-${item.id}`"
              @click="emit('price', item.id)"
            >
              Price
            </button>
            <button
              type="button"
              :data-testid="`stock-item-${item.id}`"
              @click="emit('stock', item.id)"
            >
              Stock
            </button>
            <button
              type="button"
              :data-testid="`deactivate-item-${item.id}`"
              :disabled="!item.isActive"
              @click="emit('deactivate', item.id)"
            >
              Deactivate
            </button>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script setup lang="ts">
import type { MenuItem } from '@/stores/managerMenuStore';

const emit = defineEmits(['edit', 'price', 'stock', 'deactivate']);

defineProps<{
  items: MenuItem[];
}>();

function formatPrice(priceCents: number) {
  return `$${(priceCents / 100).toFixed(2)}`;
}
</script>
