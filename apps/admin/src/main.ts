import { VueQueryPlugin } from '@tanstack/vue-query'
import './style.css'
import { createApp } from 'vue'
import { createRouter, createWebHistory } from 'vue-router'
import { routes } from 'vue-router/auto-routes'
import App from './app.vue'
import { queryClient } from './lib/query'

createApp(App)
  .use(
    createRouter({
      history: createWebHistory(),
      routes
    })
  )
  .use(VueQueryPlugin, {
    queryClient
  })
  .mount('#app')
