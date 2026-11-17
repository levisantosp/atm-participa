<script setup lang="ts">
import 'vue-sonner/style.css'
import { Spinner, Toaster } from 'ui'
import { publicRoutes } from './config'
import { auth } from './lib/auth'
import { isDark } from './lib/theme'

const title = 'Altamira Participa'
const description = 'Reporte problemas na cidade de Altamira'

useSeoMeta({
  title,
  ogTitle: title,
  description,
  ogDescription: description
})

const session = auth.useSession()
const route = useRoute()
const router = useRouter()

watch(
  [() => session.isPending.value, () => session.data.value, () => route.path],
  ([isPending, session, path]) => {
    if (isPending) return

    const publicRoute = publicRoutes.get(path)

    if (!session && publicRoute) return

    if (session && publicRoute?.action === 'redirect') {
      return void router.push('/')
    }

    if (!session && !publicRoute) {
      void auth.signOut()
      void router.push('/sign-in')
    }
  }
)
</script>

<template>
  <Body>
    <div
      v-if="session.isPending.value"
      class="flex gap-1 min-h-screen justify-center items-center"
    >
      <Spinner class="size-5" />
    </div>

    <NuxtLayout v-else>
      <NuxtPage />
    </NuxtLayout>

    <Toaster position="top-center" :theme="isDark ? 'dark' : 'light'" />
  </Body>
</template>
