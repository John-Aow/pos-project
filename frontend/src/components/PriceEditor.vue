<template>
  <section class="menu-editor" aria-labelledby="price-editor-title">
    <header class="menu-editor__header">
      <div>
        <p class="menu-editor__eyebrow">Pricing</p>
        <h3 id="price-editor-title" class="menu-editor__title">Update item price</h3>
      </div>
      <span class="menu-editor__badge" :class="item ? 'is-ready' : 'is-empty'">
        {{ item ? 'Selected' : 'Choose an item' }}
      </span>
    </header>

    <p class="menu-editor__summary">
      {{ item ? `${item.name} is ready for a new price.` : 'Pick an item from the menu.' }}
    </p>

    <form class="menu-editor__form" @submit.prevent="submit">
      <label class="menu-editor__field">
        <span class="menu-editor__label">Price (cents)</span>
        <input
          :value="draft"
          class="menu-editor__input"
          data-testid="price-editor-input"
          min="1"
          step="1"
          type="number"
          @input="updateDraft"
        />
      </label>

      <div class="menu-editor__actions">
        <button type="submit" class="menu-editor__primary" :disabled="!item">
          Save price
        </button>
        <button type="button" class="menu-editor__secondary" @click="cancel">
          Clear
        </button>
      </div>
    </form>
  </section>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue';
import type { MenuItem } from '@/stores/managerMenuStore';

const props = defineProps<{
  item: MenuItem | null;
  modelValue: number;
}>();

const emit = defineEmits(['update:modelValue', 'save', 'cancel']);

const draft = ref(props.modelValue);

watch(
  () => props.modelValue,
  (value) => {
    draft.value = value;
  },
  { immediate: true },
);

function updateDraft(event: Event) {
  const value = Number((event.target as HTMLInputElement).value);
  draft.value = value;
  emit('update:modelValue', value);
}

function submit() {
  emit('save');
}

function cancel() {
  emit('cancel');
}
</script>
