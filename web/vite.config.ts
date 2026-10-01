import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import tailwindcss from '@tailwindcss/vite'

export default defineConfig({
  plugins: [react(), tailwindcss()],
  server: {
    proxy: {
      '/api/ws': {
        target: 'ws://localhost:8080',
        ws: true,
      },
      '/api': 'http://localhost:8080',
    },
  },
  build: {
    // Keep the Vite 7 browser baseline when upgrading the bundler.
    target: ['chrome107', 'edge107', 'firefox104', 'safari16'],
    chunkSizeWarningLimit: 2000, // Increase warning limit to 2MB to avoid warning while maintaining default stable bundling
  },
})
