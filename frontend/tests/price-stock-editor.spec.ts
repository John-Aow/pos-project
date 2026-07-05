import { createPinia, setActivePinia } from 'pinia';
import { mount } from '@vue/test-utils';
import { nextTick } from 'vue';
import PriceEditor from '@/components/PriceEditor.vue';
import StockEditor from '@/components/StockEditor.vue';
import ManagerMenuPage from '@/pages/ManagerMenuPage.vue';
import { useManagerMenuStore } from '@/stores/managerMenuStore';

const seededItem = {
  id: 1,
  name: 'Green Curry',
  description: 'Thai curry with rice',
  priceCents: 12900,
  stockQuantity: 6,
  lowStockThreshold: 3,
  isAvailable: true,
  isActive: true,
  categoryId: 3,
  createdAt: '2026-07-05T00:00:00Z',
  updatedAt: '2026-07-05T00:00:00Z',
};

describe('PriceEditor', () => {
  it('emits the updated price when saved', async () => {
    const wrapper = mount(PriceEditor, {
      props: {
        item: seededItem,
        modelValue: seededItem.priceCents,
      },
    });

    await wrapper.get('[data-testid="price-editor-input"]').setValue('13900');
    await wrapper.get('form').trigger('submit');

    expect(wrapper.emitted('update:modelValue')?.[0]).toEqual([13900]);
    expect(wrapper.emitted('save')).toHaveLength(1);
  });
});

describe('StockEditor', () => {
  it('emits the updated stock quantity when saved', async () => {
    const wrapper = mount(StockEditor, {
      props: {
        item: seededItem,
        modelValue: seededItem.stockQuantity,
      },
    });

    await wrapper.get('[data-testid="stock-editor-input"]').setValue('0');
    await wrapper.get('form').trigger('submit');

    expect(wrapper.emitted('update:modelValue')?.[0]).toEqual([0]);
    expect(wrapper.emitted('save')).toHaveLength(1);
  });
});

describe('ManagerMenuPage price and stock flows', () => {
  it('updates prices immediately and blocks out-of-stock items', async () => {
    const pinia = createPinia();
    setActivePinia(pinia);
    const store = useManagerMenuStore();
    store.seedItems([seededItem]);

    const wrapper = mount(ManagerMenuPage, {
      global: {
        plugins: [pinia],
      },
    });

    const historicalSnapshot = { ...store.items[0] };

    await wrapper.get('[data-testid="price-item-1"]').trigger('click');
    await nextTick();
    wrapper.findComponent(PriceEditor).vm.$emit('update:modelValue', 14900);
    wrapper.findComponent(PriceEditor).vm.$emit('save');
    await nextTick();

    expect(store.items[0].priceCents).toBe(14900);
    expect(historicalSnapshot.priceCents).toBe(12900);
    expect(wrapper.text()).toContain('price updated');

    await wrapper.get('[data-testid="stock-item-1"]').trigger('click');
    await nextTick();
    wrapper.findComponent(StockEditor).vm.$emit('update:modelValue', 0);
    wrapper.findComponent(StockEditor).vm.$emit('save');
    await nextTick();

    expect(store.items[0].stockQuantity).toBe(0);
    expect(store.items[0].isAvailable).toBe(false);
    expect(wrapper.text()).toContain('marked unavailable');
  });
});
