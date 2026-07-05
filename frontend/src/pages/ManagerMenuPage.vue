<template>
  <ManagerShell>
    <section class="manager-page__hero" aria-labelledby="manager-menu-title">
      <div class="manager-page__hero-copy">
        <p class="manager-page__eyebrow">Catalog control</p>
        <h2 id="manager-menu-title" class="manager-page__title">
          Manage menu items and stock from one place
        </h2>
        <p class="manager-page__description">
          This workspace is ready for menu creation, pricing updates, and stock
          controls with responsive layouts for phone, tablet, desktop, and WebView
          shells.
        </p>
      </div>

      <div class="manager-page__actions" aria-label="Quick actions">
        <button type="button" class="manager-page__action-button">New item</button>
        <button type="button" class="manager-page__action-button">Update stock</button>
      </div>
    </section>

    <section class="manager-page__grid" aria-label="Workspace summary">
      <article class="manager-panel">
        <p class="manager-panel__label">Menu status</p>
        <h3 class="manager-panel__title">Ready for live catalog data</h3>
        <p class="manager-panel__body">
          The manager shell is in place so we can drop in menu tables, edit forms, and
          audit views without reworking the layout.
        </p>
      </article>

      <article class="manager-panel">
        <p class="manager-panel__label">Availability</p>
        <h3 class="manager-panel__title">Zero stock will be blocked automatically</h3>
        <p class="manager-panel__body">
          The foundation already tracks the stock-driven availability rule used by the
          customer menu and future order flows.
        </p>
      </article>

      <article class="manager-panel manager-panel--highlight">
        <p class="manager-panel__label">Next step</p>
        <h3 class="manager-panel__title">Wire in the catalog views</h3>
        <p class="manager-panel__body">
          The next phases can connect entities, repositories, and manager screens
          without changing the shell structure.
        </p>
      </article>
    </section>
  </ManagerShell>
</template>

<script setup lang="ts">
import ManagerShell from '@/components/ManagerShell.vue';
</script>
