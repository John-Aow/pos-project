<template>
  <section class="manager-menu" :data-layout-mode="layoutMode">
    <header class="manager-menu__intro">
      <p class="manager-menu__eyebrow">Catalog control</p>
      <h2 class="manager-menu__title">Manage menu items and stock</h2>
      <p class="manager-menu__description">
        Create, edit, and deactivate menu items from a responsive workspace that
        works on mobile, tablet, desktop, and WebView shells.
      </p>
    </header>

    <div class="manager-menu__status" role="status" aria-live="polite">
      {{ statusMessage }}
    </div>

    <div class="manager-menu__workspace">
      <MenuItemForm
        :model-value="store.draft"
        :editing="store.isEditing"
        @submit="handleSubmit"
        @cancel="handleCancel"
      />

      <section class="manager-menu__table-panel" aria-labelledby="menu-items-title">
        <div class="manager-menu__table-header">
          <h3 id="menu-items-title" class="manager-menu__section-title">Menu items</h3>
          <button type="button" class="manager-menu__reset-button" @click="handleStartCreate">
            New item
          </button>
        </div>

        <p v-if="!store.hasItems" class="manager-menu__empty-state">
          No menu items yet.
        </p>

        <MenuItemTable
          v-else
          :items="store.items"
          @edit="handleEdit"
          @price="handlePriceSelect"
          @stock="handleStockSelect"
          @deactivate="handleDeactivate"
        />
      </section>
    </div>

    <section class="manager-menu__editors" aria-label="Price and stock editors">
      <PriceEditor
        :item="store.selectedItem"
        :model-value="store.priceDraft"
        @update:model-value="store.updatePriceDraft"
        @save="handleSavePrice"
        @cancel="handleClearSelection"
      />

      <StockEditor
        :item="store.selectedItem"
        :model-value="store.stockDraft"
        @update:model-value="store.updateStockDraft"
        @save="handleSaveStock"
        @cancel="handleClearSelection"
      />
    </section>
  </section>
</template>

<script setup lang="ts">
import { onBeforeUnmount, onMounted, ref } from 'vue';
import MenuItemForm from '@/components/MenuItemForm.vue';
import MenuItemTable from '@/components/MenuItemTable.vue';
import PriceEditor from '@/components/PriceEditor.vue';
import StockEditor from '@/components/StockEditor.vue';
import { useManagerMenuStore, type MenuItemDraft } from '@/stores/managerMenuStore';

const store = useManagerMenuStore();
const statusMessage = ref('Ready to manage menu items.');
const layoutMode = ref(getLayoutMode());

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

function handleSubmit(draft: MenuItemDraft) {
  store.updateDraft(draft);
  try {
    const item = store.saveDraft();
    statusMessage.value = `${item.name} saved.`;
  } catch (error) {
    statusMessage.value = error instanceof Error ? error.message : 'Unable to save item.';
  }
}

function handleCancel() {
  store.beginCreate();
  store.clearSelection();
  statusMessage.value = 'Draft cleared.';
}

function handleStartCreate() {
  store.beginCreate();
  store.clearSelection();
  statusMessage.value = 'Creating a new item.';
}

function handleEdit(id: number) {
  store.clearSelection();
  store.beginEdit(id);
  statusMessage.value = 'Editing an existing item.';
}

function handlePriceSelect(id: number) {
  store.beginPriceEdit(id);
  statusMessage.value = 'Editing price.';
}

function handleStockSelect(id: number) {
  store.beginStockEdit(id);
  statusMessage.value = 'Editing stock.';
}

function handleDeactivate(id: number) {
  store.deactivateItem(id);
  statusMessage.value = 'Item deactivated.';
}

function handleSavePrice() {
  try {
    const item = store.savePriceDraft();
    statusMessage.value = `${item.name} price updated.`;
  } catch (error) {
    statusMessage.value = error instanceof Error ? error.message : 'Unable to save price.';
  }
}

function handleSaveStock() {
  try {
    const item = store.saveStockDraft();
    statusMessage.value =
      item.stockQuantity === 0
        ? `${item.name} stock set to zero and marked unavailable.`
        : `${item.name} stock updated.`;
  } catch (error) {
    statusMessage.value = error instanceof Error ? error.message : 'Unable to save stock.';
  }
}

function handleClearSelection() {
  store.clearSelection();
  statusMessage.value = 'Selection cleared.';
}

onMounted(() => {
  window.addEventListener('resize', syncLayoutMode);
  if (!store.hasItems && !store.isEditing) {
    store.beginCreate();
  }
});

onBeforeUnmount(() => {
  window.removeEventListener('resize', syncLayoutMode);
});
</script>
