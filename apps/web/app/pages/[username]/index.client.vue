<script setup lang="ts">
import {
  BadgeCheckIcon,
  EllipsisVertical,
  Share,
  SquareExclamationPoint,
  UserRoundIcon
} from '@lucide/vue'
import { useGetUsersByUsername, type User } from 'api-client'
import {
  Card,
  CardContent,
  Input,
  Label,
  Menubar,
  MenubarContent,
  MenubarItem,
  MenubarMenu,
  MenubarTrigger,
  Spinner
} from 'ui'
import { toast } from 'vue-sonner'
import { auth } from '~/lib/auth'
import ReportDialog from './report-dialog.vue'

definePageMeta({
  layout: 'default'
})

const route = useRoute('username')
const router = useRouter()
const reportUserId = ref<User['id'] | null>(null)
const session = auth.useSession()

watch(
  () => route.params.username,
  (username) => {
    if (username === session.data.value?.username) {
      navigateTo('/perfil', {
        replace: true
      })
    }
  },
  {
    immediate: true
  }
)

const { isPending, data, error } = useGetUsersByUsername({
  path: computed(() => ({
    username: route.params.username
  }))
})

watch(error, (err) => {
  if (err?.status === 404) {
    toast.error('Usuário não encontrado')
    router.back()
  }
})

async function copy() {
  await navigator.clipboard.writeText(
    new URL(route.path, window.location.origin).href
  )
  toast.success('Link copiado para área de transferência')
}
</script>

<template>
  <div class="mx-auto flex w-full max-w-xl flex-col gap-6 px-4 pb-12">
    <header class="flex items-center gap-2">
      <UserRoundIcon :size="23" class="text-primary" />
      <h1 class="text-2xl font-bold tracking-tight">Perfil de usuário</h1>
    </header>

    <div v-if="isPending" class="flex justify-center py-16">
      <Spinner />
    </div>

    <Card v-else-if="data" class="relative overflow-hidden">
      <CardContent class="p-5 sm:p-6">
        <Menubar class="absolute right-3 top-3 border-none">
          <MenubarMenu>
            <MenubarTrigger aria-label="Ações do perfil" class="cursor-pointer">
              <EllipsisVertical />
            </MenubarTrigger>
            <MenubarContent align="end">
              <MenubarItem class="cursor-pointer" @select="copy">
                <Share />
                Copiar link
              </MenubarItem>
              <MenubarItem
                class="cursor-pointer text-destructive hover:bg-destructive!"
                @select="reportUserId = data.id"
              >
                <SquareExclamationPoint />
                Reportar
              </MenubarItem>
            </MenubarContent>
          </MenubarMenu>
        </Menubar>

        <div class="flex flex-col items-center text-center">
          <div
            class="flex size-24 items-center justify-center rounded-full border-4 border-background bg-primary/10 text-2xl font-semibold text-primary ring-1 ring-border sm:size-28"
            aria-label="Avatar do perfil"
          >
            <UserRoundIcon :size="42" />
          </div>
          <h2 class="mt-4 max-w-full truncate text-xl font-bold tracking-tight">
            {{ data.displayName }}
          </h2>
          <p class="mt-1 text-sm text-muted-foreground">@{{ data.username }}</p>
          <span
            v-if="data.isAdmin"
            class="mt-3 inline-flex items-center gap-1.5 rounded-full border border-primary/20 bg-primary/10 px-3 py-1 text-xs font-medium text-primary"
          >
            <BadgeCheckIcon :size="14" />
            Administrador
          </span>
        </div>

        <div class="mt-6 space-y-4 border-t pt-5">
          <div class="grid gap-2">
            <Label for="displayName">Nome de exibição</Label>
            <Input
              id="displayName"
              :model-value="data.displayName"
              disabled
              class="disabled:border-transparent disabled:bg-muted/50 disabled:shadow-none"
            />
          </div>
          <div class="grid gap-2">
            <Label for="username">Nome de usuário</Label>
            <Input
              id="username"
              :model-value="`@${data.username}`"
              disabled
              class="disabled:border-transparent disabled:bg-muted/50 disabled:shadow-none"
            />
          </div>
        </div>
      </CardContent>
    </Card>
    <p v-else-if="error" class="text-sm text-muted-foreground">
      {{
        error.status === 404
          ? 'Usuário não encontrado'
          : 'Não foi possível carregar o perfil'
      }}
    </p>

    <ReportDialog v-model:user-id="reportUserId" />
  </div>
</template>
