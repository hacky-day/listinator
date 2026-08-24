import { fileURLToPath } from "node:url";
import { mergeConfig, defineConfig } from "vitest/config";

import viteConfig from "./vite.config";

// Reuses the vite config so tests resolve the "@/" alias the same way the app
// does. No DOM environment is needed, the tests only exercise modules.
export default mergeConfig(
  viteConfig,
  defineConfig({
    test: {
      environment: "node",
      include: ["src/**/__tests__/**/*.test.ts"],
      // Polyfills the indexedDB global (absent in the "node" environment)
      // for src/offline's tests, which exercise real idb transactions
      // rather than mocking storage out.
      setupFiles: ["src/offline/__tests__/setup.ts"],
      root: fileURLToPath(new URL("./", import.meta.url)),
    },
  }),
);
