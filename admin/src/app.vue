<script setup lang="ts">
import { Spinner } from 'ui'
import { computed, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { publicRoutes } from './config'
import Private from './layouts/private.vue'
import Public from './layouts/public.vue'
import { auth } from './lib/auth'

const route = useRoute()
const router = useRouter()
const session = auth.useSession()

watch(
  [() => session.isFetching.value, () => session.data.value, () => route.path],
  ([isFetching, session, path]) => {
    if (isFetching) return

    const publicRoute = publicRoutes.has(path)

    if (!session && publicRoute) return

    if (session && publicRoute) {
      return router.push('/')
    }

    if (!session && !publicRoute) {
      router.push('/sign-in')
    }
  }
)

const layouts = {
  public: Public,
  private: Private
} as const

const layout = computed(() => {
  const layout = route.meta.layout as keyof typeof layouts
  return layouts[layout] ?? Private
})
</script>

<template>
  <div
    v-if="session.isFetching.value"
    class="flex gap-1 min-h-screen justify-center items-center"
  >
    <Spinner class="size-5" />
  </div>

  <RouterView v-else v-slot="{ Component }">
    <component :is="layout">
      <component :is="Component" />
    </component>
  </RouterView>
</template>
