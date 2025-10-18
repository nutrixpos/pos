import { fileURLToPath, URL } from 'node:url'
import path from 'path'; // Import Node.js path module

import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import vueJsx from '@vitejs/plugin-vue-jsx'
import vueDevTools from 'vite-plugin-vue-devtools'
import { viteSingleFile } from "vite-plugin-singlefile"


// https://vite.dev/config/
export default defineConfig({
  css: {
    preprocessorOptions: {
      scss: {
        quietDeps: true, // This will silence deprecation warnings
      }
    }
  },
  build: {
    sourcemap: true, // <-- for production builds
    minify: false,
    terserOptions: {
      compress: false,
      mangle: false,
    },
  },
  plugins: [
    vue(),
    vueJsx(),
    vueDevTools(),
    viteSingleFile(),
  ],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url)),
      'vue-oidc-client/vue3': path.resolve(__dirname, './node_modules/vue-oidc-client/vue3/src/vue-oidc-client.ts'), // Non-minified path
    },
  },
})
