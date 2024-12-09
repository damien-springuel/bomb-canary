import { defineConfig } from 'vite'
import { svelte } from '@sveltejs/vite-plugin-svelte'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [svelte()],
  server: {
    proxy: {
      '/events': {
        ws: true,
        target: 'ws://localhost:44324'
      },
      '^/actions/.*': {
        target: 'http://localhost:44324'
      },
      '^/party/.*': {
        target: 'http://localhost:44324'
      }
    }
  }
})
