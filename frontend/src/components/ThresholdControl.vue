<template>
  <section class="threshold-control" aria-labelledby="threshold-control-title">
    <header class="threshold-control__header">
      <p class="manager-panel__label">Stock warnings</p>
      <h3 id="threshold-control-title" class="manager-panel__title">Adjust threshold</h3>
    </header>

    <form class="threshold-control__form" @submit.prevent="handleSubmit">
      <label class="threshold-control__field">
        <span class="threshold-control__label">Low-stock threshold</span>
        <input
          :value="modelValue"
          class="threshold-control__input"
          data-testid="threshold-control-input"
          inputmode="numeric"
          min="0"
          type="number"
          @input="handleInput"
        />
      </label>

      <button type="submit" class="threshold-control__save">Save threshold</button>
    </form>
  </section>
</template>

<script setup lang="ts">
defineProps<{
  modelValue: number;
}>();

const emit = defineEmits(['update:modelValue', 'save']);

function handleInput(event: Event) {
  const value = Number((event.target as HTMLInputElement).value);
  emit('update:modelValue', Number.isFinite(value) ? value : 0);
}

function handleSubmit() {
  emit('save');
}
</script>
