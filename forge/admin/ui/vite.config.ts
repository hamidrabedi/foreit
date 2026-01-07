import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react-swc'
import path from 'path'

// https://vite.dev/config/
export default defineConfig({
  base: './',
  plugins: [react()],
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'),
      seroval: path.resolve(
        __dirname,
        './node_modules/seroval/dist/esm/production/index.mjs'
      ),
      immer: path.resolve(__dirname, './node_modules/immer/dist/immer.mjs'),
      reselect: path.resolve(
        __dirname,
        './node_modules/reselect/dist/reselect.mjs'
      ),
    },
  },
  server: {
    port: 3000,
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
    },
  },
})
