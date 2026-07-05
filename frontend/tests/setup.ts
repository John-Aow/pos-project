import { afterEach, beforeEach } from 'vitest';

type Viewport = {
  width: number;
  height: number;
};

const defaultViewport: Viewport = {
  width: 1440,
  height: 900,
};

function applyViewport(viewport: Viewport) {
  Object.defineProperty(window, 'innerWidth', {
    configurable: true,
    value: viewport.width,
  });
  Object.defineProperty(window, 'innerHeight', {
    configurable: true,
    value: viewport.height,
  });
  window.dispatchEvent(new Event('resize'));
}

export function setViewport(viewport: Viewport) {
  applyViewport(viewport);
}

beforeEach(() => {
  if (!window.matchMedia) {
    window.matchMedia = ((query: string) =>
      ({
        matches: false,
        media: query,
        onchange: null,
        addEventListener: () => undefined,
        removeEventListener: () => undefined,
        addListener: () => undefined,
        removeListener: () => undefined,
        dispatchEvent: () => false,
      }) as unknown as MediaQueryList);
  }

  applyViewport(defaultViewport);
});

afterEach(() => {
  applyViewport(defaultViewport);
});
