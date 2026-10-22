<script setup lang="ts">
import 'vue-sonner/style.css'
import { Spinner, Toaster } from 'ui'
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
  async ([isFetching, session, path]) => {
    if (isFetching) return

    const publicRoute = publicRoutes.has(path)

    if (!session && publicRoute) return

    if (session?.isAdmin && publicRoute) {
      return await router.push('/')
    }

    if (!session?.isAdmin && !publicRoute) {
      if (session) {
        await auth.signOut()
      }
      await router.push('/sign-in')
    }
  },
  {
    immediate: true
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

  <Toaster position="top-center" theme="dark" />
</template>
