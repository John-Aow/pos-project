import { createPinia, setActivePinia } from 'pinia';
import { mount } from '@vue/test-utils';
import ManagerMenuPage from '@/pages/ManagerMenuPage.vue';
import { setViewport } from './setup';

describe('ManagerMenuPage responsive layout', () => {
  it.each([
    [{ width: 375, height: 812 }, 'mobile'],
    [{ width: 768, height: 1024 }, 'tablet'],
    [{ width: 1280, height: 900 }, 'desktop'],
  ])('sets the layout mode for viewport %j', async (viewport, expectedMode) => {
    const pinia = createPinia();
    setActivePinia(pinia);
    setViewport(viewport);

    const wrapper = mount(ManagerMenuPage, {
      global: {
        plugins: [pinia],
      },
    });

    expect(wrapper.get('[data-layout-mode]').attributes('data-layout-mode')).toBe(
      expectedMode,
    );
  });
});
