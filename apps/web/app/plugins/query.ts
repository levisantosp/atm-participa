import { VueQueryPlugin } from '@tanstack/vue-query'
import { queryClient } from '~/lib/query-client'

export default defineNuxtPlugin((app) => {
  app.vueApp.use(VueQueryPlugin, {
    queryClient
  })
})
