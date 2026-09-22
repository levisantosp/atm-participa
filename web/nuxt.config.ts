import tailwindcss from '@tailwindcss/vite'

// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  compatibilityDate: '2025-07-15',
  devtools: {
    enabled: true
  },
  pages: {
    pattern: ['**/index.vue', '**/index.client.vue']
  },
  devServer: {
    port: 5174
  },
  vite: {
    plugins: [tailwindcss()]
  },
  css: ['~/assets/css/main.css']
})
