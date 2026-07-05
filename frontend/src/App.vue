<template>
  <ManagerShell :active-view="activeView" @navigate="activeView = $event">
    <MenuManagementView v-if="activeView === 'menu'" />
    <LowStockView v-else />
  </ManagerShell>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import ManagerShell from './components/ManagerShell.vue';
import MenuManagementView from './pages/MenuManagementView.vue';
import LowStockView from './pages/LowStockView.vue';
import './styles/manager.css';

const activeView = ref<'menu' | 'low-stock'>('menu');
</script>
