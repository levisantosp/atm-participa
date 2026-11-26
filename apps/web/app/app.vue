<script setup lang="ts">
import 'vue-sonner/style.css'
import dayjs from 'dayjs'
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
  },
  {
    immediate: true
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

    <div
      v-else-if="session.data.value?.banned"
      class="flex min-h-screen flex-col items-center justify-center gap-4 px-6 text-center"
    >
      <h1 class="text-2xl font-semibold">Sua conta foi banida</h1>
      <dl class="flex flex-col gap-3">
        <div>
          <dt class="text-sm text-muted-foreground">Início do banimento</dt>
          <dd>
            {{
              session.data.value.bannedAt
                ? dayjs(session.data.value.bannedAt).format(
                    'DD/MM/YYYY [às] HH:mm'
                  )
                : 'Data não informada'
            }}
          </dd>
        </div>
        <div>
          <dt class="text-sm text-muted-foreground">Término do banimento</dt>
          <dd>
            {{
              session.data.value.banEndsAt
                ? dayjs(session.data.value.banEndsAt).format(
                    'DD/MM/YYYY [às] HH:mm'
                  )
                : 'Sem prazo de término'
            }}
          </dd>
        </div>
      </dl>
    </div>

    <NuxtLayout v-else>
      <NuxtPage />
    </NuxtLayout>

    <Toaster position="top-center" :theme="isDark ? 'dark' : 'light'" />
  </Body>
</template>
