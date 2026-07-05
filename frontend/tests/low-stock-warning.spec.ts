import { createPinia, setActivePinia } from 'pinia';
import { mount } from '@vue/test-utils';
import { nextTick } from 'vue';
import { readFileSync } from 'node:fs';
import App from '@/App.vue';
import LowStockAlert from '@/components/LowStockAlert.vue';
import LowStockView from '@/pages/LowStockView.vue';
import { useInventoryStore } from '@/stores/inventoryStore';
import { setViewport } from './setup';

const demoItems = [
  {
    id: 1,
    name: 'Green Curry',
    description: 'Thai curry with jasmine rice',
    priceCents: 12900,
    stockQuantity: 3,
    lowStockThreshold: 5,
    isAvailable: true,
    isActive: true,
    categoryId: 1,
    createdAt: '2026-07-05T00:00:00Z',
    updatedAt: '2026-07-05T00:00:00Z',
  },
  {
    id: 2,
    name: 'Pad Thai',
    description: 'Rice noodles with prawns',
    priceCents: 13900,
    stockQuantity: 8,
    lowStockThreshold: 5,
    isAvailable: true,
    isActive: true,
    categoryId: 1,
    createdAt: '2026-07-05T00:00:00Z',
    updatedAt: '2026-07-05T00:00:00Z',
  },
];

describe('LowStockView', () => {
  it('shows low-stock items and updates the warning threshold', async () => {
    const pinia = createPinia();
    setActivePinia(pinia);
    const store = useInventoryStore();
    store.seedItems(demoItems);
    store.beginThresholdEdit();

    const wrapper = mount(LowStockView, {
      global: {
        plugins: [pinia],
      },
    });

    expect(wrapper.text()).toContain('Green Curry');
    expect(wrapper.text()).not.toContain('Pad Thai');

    await wrapper.get('[data-testid="threshold-control-input"]').setValue('2');
    await wrapper.get('form').trigger('submit');
    await nextTick();

    expect(store.threshold).toBe(2);
    expect(wrapper.text()).toContain('Low-stock threshold set to 2.');
    expect(wrapper.text()).not.toContain('Green Curry');
  });

  it('shows a refreshable alert summary', async () => {
    const wrapper = mount(LowStockAlert, {
      props: {
        items: demoItems,
        threshold: 8,
        loading: false,
      },
    });

    expect(wrapper.text()).toContain('2 items at or below 8.');
    await wrapper.get('[data-testid="low-stock-alert-refresh"]').trigger('click');
    expect(wrapper.emitted('refresh')).toHaveLength(1);
  });

  it('renders a loading state while refreshing', () => {
    const wrapper = mount(LowStockAlert, {
      props: {
        items: [],
        threshold: 5,
        loading: true,
      },
    });

    expect(wrapper.text()).toContain('Loading warning data.');
  });

  it.each([
    [{ width: 375, height: 812 }, 'mobile'],
    [{ width: 768, height: 1024 }, 'tablet'],
    [{ width: 1280, height: 900 }, 'desktop'],
  ])('uses the layout mode for viewport %j', async (viewport, expectedMode) => {
    const pinia = createPinia();
    setActivePinia(pinia);
    setViewport(viewport);

    const wrapper = mount(LowStockView, {
      global: {
        plugins: [pinia],
      },
    });

    expect(wrapper.get('[data-layout-mode]').attributes('data-layout-mode')).toBe(
      expectedMode,
    );
  });
});

describe('App navigation', () => {
  it('switches between menu and low-stock views inside the shell', async () => {
    const pinia = createPinia();
    setActivePinia(pinia);

    const wrapper = mount(App, {
      global: {
        plugins: [pinia],
      },
    });

    expect(wrapper.text()).toContain('Manage menu items and stock');
    expect(wrapper.text()).toContain('Low-stock alert');

    await wrapper.get('button[aria-pressed="false"]').trigger('click');
    await nextTick();

    expect(wrapper.text()).toContain('Review stock warnings');

    await wrapper.get('button[aria-pressed="false"]').trigger('click');
    await nextTick();

    expect(wrapper.text()).toContain('Manage menu items and stock');
  });
});

describe('Viewport meta', () => {
  it('keeps the WebView-safe viewport configuration', () => {
    const html = readFileSync('index.html', 'utf8');
    expect(html).toContain('width=device-width, initial-scale=1.0, viewport-fit=cover');
  });
});
