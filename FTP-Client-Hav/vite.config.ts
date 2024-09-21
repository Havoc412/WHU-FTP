import { fileURLToPath, URL } from 'node:url'

import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import electronRenderer from "vite-plugin-electron-renderer";

// https://vitejs.dev/config/
export default defineConfig({
  base: "./",
  plugins: [vue(), electronRenderer()],
  resolve: {
    alias: {
      "@": fileURLToPath(new URL("./src", import.meta.url)),
    },
  },
  server: {
    port: 3000,
    open: false,
  },
});
