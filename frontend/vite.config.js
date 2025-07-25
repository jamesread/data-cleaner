import { defineConfig } from 'vite'

export default defineConfig({
  server: {
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
        secure: false,
      },
      '/lang': {
        target: 'http://localhost:8080',
        changeOrigin: true,
        secure: false,
      }
    },
  },
  plugins: [],
})
