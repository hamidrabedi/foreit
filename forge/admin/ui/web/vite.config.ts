import { defineConfig } from "vite";
import react from "@vitejs/plugin-react-swc";
import { TanStackRouterVite } from "@tanstack/router-plugin/vite";
import path from "path";

// https://vite.dev/config/
export default defineConfig({
  plugins: [TanStackRouterVite(), react()],
  resolve: {
    alias: {
      "@": path.resolve(import.meta.dirname, "./src"),
    },
  },
  server: {
    port: 3000,
    proxy: {
      "/api": {
        target: "http://localhost:8080",
        changeOrigin: true,
      },
    },
  },
  base: "/admin/",
  build: {
    outDir: "../dist",
    emptyOutDir: true,
    chunkSizeWarningLimit: 600,
    // NOTE: `build.rollupOptions.output.manualChunks` as an object map was tried
    // first (per rules/web/performance.md task 7.3), but this project builds with
    // rolldown (see `overrides` in package.json) and rolldown rejects the object
    // form: "Invalid type: Expected Function but received Object". The
    // `codeSplitting.groups` below express the same vendor split in the
    // rolldown-supported shape.
    rolldownOptions: {
      output: {
        codeSplitting: {
          groups: [
            { name: "react", test: /node_modules[\\/](react|react-dom|scheduler)[\\/]/ },
            { name: "router", test: /node_modules[\\/]@tanstack[\\/]react-router[\\/]/ },
            { name: "query", test: /node_modules[\\/]@tanstack[\\/]react-query[\\/]/ },
            { name: "charts", test: /node_modules[\\/]recharts[\\/]/ },
            { name: "motion", test: /node_modules[\\/]framer-motion[\\/]/ },
            {
              name: "radix",
              test: /node_modules[\\/]@radix-ui[\\/]react-(dialog|dropdown-menu|select|tabs|tooltip|alert-dialog|label|slot)[\\/]/,
            },
            {
              name: "forms",
              test: /node_modules[\\/](react-hook-form|zod)[\\/]|node_modules[\\/]@hookform[\\/]resolvers[\\/]/,
            },
          ],
        },
      },
    },
  },
});
