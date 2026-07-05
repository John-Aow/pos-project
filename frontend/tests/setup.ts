import { afterEach, beforeEach } from 'vitest';

type Viewport = {
  width: number;
  height: number;
};

export const mobileViewport: Viewport = {
  width: 375,
  height: 812,
};

export const tabletViewport: Viewport = {
  width: 768,
  height: 1024,
};

export const desktopViewport: Viewport = {
  width: 1280,
  height: 900,
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

export function setMobileViewport() {
  applyViewport(mobileViewport);
}

export function setTabletViewport() {
  applyViewport(tabletViewport);
}

export function setDesktopViewport() {
  applyViewport(desktopViewport);
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
