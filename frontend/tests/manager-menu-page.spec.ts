import { createPinia, setActivePinia } from 'pinia';
import { mount } from '@vue/test-utils';
import { nextTick } from 'vue';
import MenuItemForm from '@/components/MenuItemForm.vue';
import ManagerMenuPage from '@/pages/ManagerMenuPage.vue';
import { useManagerMenuStore } from '@/stores/managerMenuStore';

describe('ManagerMenuPage', () => {
  it('creates, edits, and deactivates a menu item', async () => {
    const pinia = createPinia();
    setActivePinia(pinia);
    const wrapper = mount(ManagerMenuPage, {
      global: {
        plugins: [pinia],
      },
    });
    const store = useManagerMenuStore();

    wrapper.findComponent(MenuItemForm).vm.$emit('submit', {
      name: 'Green Curry',
      description: 'Thai curry',
      priceCents: 12900,
      stockQuantity: 12,
      lowStockThreshold: 4,
      categoryId: 3,
    });
    await nextTick();

    expect(store.items).toHaveLength(1);
    expect(store.items[0].name).toBe('Green Curry');
    expect(store.items[0].lowStockThreshold).toBe(4);
    expect(wrapper.text()).toContain('Green Curry saved.');

    await wrapper.get('[data-testid="edit-item-1"]').trigger('click');
    await nextTick();
    expect(wrapper.text()).toContain('Editing an existing item.');
    wrapper.findComponent(MenuItemForm).vm.$emit('submit', {
      name: 'Green Curry Deluxe',
      description: 'Thai curry',
      priceCents: 12900,
      stockQuantity: 12,
      lowStockThreshold: 5,
      categoryId: 3,
    });
    await nextTick();

    expect(store.items[0].name).toBe('Green Curry Deluxe');
    expect(store.items[0].lowStockThreshold).toBe(5);
    expect(wrapper.text()).toContain('Green Curry Deluxe saved.');

    await wrapper.get('[data-testid="deactivate-item-1"]').trigger('click');

    expect(store.items[0].isActive).toBe(false);
    expect(store.items[0].isAvailable).toBe(false);
    expect(wrapper.text()).toContain('Item deactivated.');
  });

  it('clears the draft when the reset button is pressed', async () => {
    const pinia = createPinia();
    setActivePinia(pinia);
    const wrapper = mount(ManagerMenuPage, {
      global: {
        plugins: [pinia],
      },
    });
    const store = useManagerMenuStore();

    store.updateDraft({
      name: 'Sticky Rice',
      description: 'With mango',
      priceCents: 8000,
      stockQuantity: 5,
      lowStockThreshold: 2,
      categoryId: 2,
    });
    await nextTick();

    await wrapper.get('.manager-menu__reset-button').trigger('click');

    expect(store.draft.name).toBe('');
    expect(wrapper.text()).toContain('Creating a new item.');
  });
});
