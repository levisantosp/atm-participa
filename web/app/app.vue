<script setup lang="ts">
import 'vue-sonner/style.css'
import { Spinner, Toaster } from 'ui'
import { publicRoutes } from './config'
import { auth } from './lib/auth'

const session = auth.useSession()
const route = useRoute()
const router = useRouter()

watch(
  [() => session.isPending.value, () => session.data.value],
  ([isPending, session]) => {
    if (isPending) return

    const publicRoute = publicRoutes.get(route.path)

    if (!session && publicRoute) return

    if (session && publicRoute?.action === 'redirect') {
      return router.push('/')
    }

    if (!session && !publicRoute) {
      router.push('/sign-in')
    }
  }
)
</script>

<template>
  <Html class="dark">
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

      <Toaster position="top-center" theme="dark" />
    </Body>
  </Html>
</template>
