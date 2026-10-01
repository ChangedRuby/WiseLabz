/// <reference types="vite/client" />

// This font package exposes CSS without TypeScript declarations.
declare module '@fontsource-variable/big-shoulders-text';

interface ImportMetaEnv {
  /** Force-enable/disable the mock layer (MSW + WS). Default: on in dev. */
  readonly VITE_USE_MOCKS?: 'true' | 'false';
}
