<template>
  <div class="manager-shell">
    <header class="manager-shell__header">
      <div class="manager-shell__brand">
        <p class="manager-shell__eyebrow">POS Manager</p>
        <h1 class="manager-shell__title">Menu, pricing, and stock</h1>
      </div>

      <nav class="manager-shell__nav" aria-label="Manager navigation">
        <button type="button" class="manager-shell__nav-button" aria-pressed="true">
          Menu
        </button>
        <button type="button" class="manager-shell__nav-button">Pricing</button>
        <button type="button" class="manager-shell__nav-button">Stock</button>
        <button type="button" class="manager-shell__nav-button">History</button>
      </nav>
    </header>

    <main class="manager-shell__main">
      <slot />
    </main>
  </div>
</template>

<script setup lang="ts"></script>
