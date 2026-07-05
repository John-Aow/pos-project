<template>
  <form class="menu-item-form" @submit.prevent="emitSubmit">
    <div class="menu-item-form__grid">
      <label class="menu-item-form__field">
        <span class="menu-item-form__label">Name</span>
        <input
          v-model="form.name"
          class="menu-item-form__input"
          data-testid="menu-item-name"
          type="text"
          autocomplete="off"
        />
      </label>

      <label class="menu-item-form__field">
        <span class="menu-item-form__label">Description</span>
        <textarea
          v-model="form.description"
          class="menu-item-form__textarea"
          data-testid="menu-item-description"
          rows="3"
        ></textarea>
      </label>

      <div class="menu-item-form__row">
        <label class="menu-item-form__field">
          <span class="menu-item-form__label">Price (cents)</span>
          <input
            v-model.number="form.priceCents"
            class="menu-item-form__input"
            data-testid="menu-item-price"
            type="number"
            min="1"
            step="1"
          />
        </label>

        <label class="menu-item-form__field">
          <span class="menu-item-form__label">Stock</span>
          <input
            v-model.number="form.stockQuantity"
            class="menu-item-form__input"
            data-testid="menu-item-stock"
            type="number"
            min="0"
            step="1"
          />
        </label>

        <label class="menu-item-form__field">
          <span class="menu-item-form__label">Low-stock threshold</span>
          <input
            v-model.number="form.lowStockThreshold"
            class="menu-item-form__input"
            data-testid="menu-item-low-stock-threshold"
            type="number"
            min="0"
            step="1"
          />
        </label>

        <label class="menu-item-form__field">
          <span class="menu-item-form__label">Category</span>
          <input
            v-model.number="form.categoryId"
            class="menu-item-form__input"
            data-testid="menu-item-category"
            type="number"
            min="1"
            step="1"
          />
        </label>
      </div>
    </div>

    <div class="menu-item-form__actions">
      <button type="submit" class="menu-item-form__primary">
        {{ editing ? 'Save changes' : 'Create item' }}
      </button>
      <button
        type="button"
        class="menu-item-form__secondary"
        data-testid="menu-item-form-cancel"
        @click="emitCancel"
      >
        Clear
      </button>
    </div>
  </form>
</template>

<script setup lang="ts">
import { reactive, watch } from 'vue';
import type { MenuItemDraft } from '@/stores/managerMenuStore';

const props = defineProps<{
  modelValue: MenuItemDraft;
  editing: boolean;
}>();

const emit = defineEmits(['submit', 'cancel']);

const form = reactive<MenuItemDraft>({
  name: '',
  description: '',
  priceCents: 0,
  stockQuantity: 0,
  lowStockThreshold: 5,
  categoryId: 1,
});

watch(
  () => props.modelValue,
  (value) => {
    form.name = value.name;
    form.description = value.description;
    form.priceCents = value.priceCents;
    form.stockQuantity = value.stockQuantity;
    form.lowStockThreshold = value.lowStockThreshold;
    form.categoryId = value.categoryId;
  },
  { immediate: true, deep: true },
);

function emitSubmit() {
  emit('submit', {
    name: form.name,
    description: form.description,
    priceCents: form.priceCents,
    stockQuantity: form.stockQuantity,
    lowStockThreshold: form.lowStockThreshold,
    categoryId: form.categoryId,
  });
}

function emitCancel() {
  emit('cancel');
}
</script>
