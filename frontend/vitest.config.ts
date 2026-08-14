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
      root: fileURLToPath(new URL("./", import.meta.url)),
    },
  }),
);
