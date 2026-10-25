<script setup lang="ts">
import {
  BadgeCheckIcon,
  CalendarDaysIcon,
  ChevronRightIcon,
  LoaderCircleIcon,
  LogOutIcon,
  LockKeyholeIcon,
  MegaphoneIcon,
  PencilLineIcon,
  ShieldCheckIcon,
  UserRoundIcon
} from '@lucide/vue'
import { toTypedSchema } from '@vee-validate/zod'
import {
  getAuthMeQueryKey,
  useGetUsersByUserIdIssuesInfinite,
  usePutUsersMe
} from 'api-client'
import dayjs from 'dayjs'
import {
  Button,
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
  Input,
  Label,
  Spinner
} from 'ui'
import { useForm } from 'vee-validate'
import { toast } from 'vue-sonner'
import { z } from 'zod'
import { auth } from '~/lib/auth'
import { queryClient } from '~/lib/query-client'

definePageMeta({
  layout: 'default'
})

const session = auth.useSession()
const config = useRuntimeConfig()
const isEditing = ref(false)

const { mutate: updateProfile, isPending: isUpdatingProfile } = usePutUsersMe({
  mutation: {
    onSuccess(data) {
      isEditing.value = false
      queryClient.setQueryData(getAuthMeQueryKey(), data)
      toast.success('Perfil atualizado com sucesso')
    }
  }
})

const schema = z.object({
  displayName: z
    .string()
    .trim()
    .min(1, 'Informe seu nome')
    .max(100, 'O nome deve ter no máximo 100 caracteres'),
  username: z
    .string()
    .trim()
    .min(3, 'Use pelo menos 3 caracteres')
    .max(32, 'Use no máximo 32 caracteres')
    .regex(/^[a-zA-Z0-9_]+$/, 'Use apenas letras, números e underscore')
})

const { handleSubmit, errors, defineField, setValues } = useForm({
  validationSchema: toTypedSchema(schema)
})

const [displayName, displayNameAttrs] = defineField('displayName')
const [username, usernameAttrs] = defineField('username')

watch(
  () => session.data.value,
  (user) => {
    if (user) {
      setValues({
        displayName: user.displayName,
        username: user.username
      })
    }
  },
  {
    immediate: true
  }
)

watch(
  () => session.isPending.value,
  (isPending) => {
    if (!isPending && !session.data.value) {
      navigateTo('/sign-in')
    }
  },
  {
    immediate: true
  }
)

const issuePath = computed(() => ({
  userId: BigInt(session.data.value?.userId ?? '0')
}))

const {
  data: issuesData,
  isPending: isLoadingIssues,
  isError: hasIssuesError,
  hasNextPage: hasMoreIssues,
  isFetchingNextPage: isLoadingMoreIssues,
  fetchNextPage: fetchNextIssues
} = useGetUsersByUserIdIssuesInfinite(
  {
    path: issuePath,
    query: {
      limit: 5n
    }
  },
  {
    query: {
      enabled: computed(() => !!session.data.value)
    }
  }
)

const issues = computed(
  () => issuesData.value?.pages.flatMap((page) => page.items ?? []) ?? []
)

const initials = computed(() => {
  const name = session.data.value?.displayName ?? ''
  return name
    .split(/\s+/)
    .filter(Boolean)
    .slice(0, 2)
    .map((part) => part[0]?.toLocaleUpperCase())
    .join('')
})

const statusLabel = (status: string) => {
  const labels: Record<string, string> = {
    open: 'Aberta',
    in_review: 'Em análise',
    closed: 'Concluída'
  }
  return labels[status] ?? status
}

const onSubmit = handleSubmit((values) =>
  updateProfile({
    body: {
      displayName: values.displayName,
      username: values.username
    }
  })
)

const cancelEditing = () => {
  setValues({
    displayName: session.data.value?.displayName ?? '',
    username: session.data.value?.username ?? ''
  })
  isEditing.value = false
}
</script>

<template>
  <div class="mx-auto flex w-full max-w-xl flex-col gap-6 px-4 pb-12">
    <header class="space-y-1">
      <div class="flex items-center gap-2">
        <UserRoundIcon :size="23" class="text-primary" />
        <h1 class="text-2xl font-bold tracking-tight">Meu perfil</h1>
      </div>
      <p class="text-sm text-muted-foreground">
        Visualize e gerencie suas informações na plataforma.
      </p>
    </header>

    <Card class="overflow-hidden">
      <CardContent class="p-5 sm:p-6">
        <div class="flex flex-col items-center text-center">
          <div
            class="flex size-24 items-center justify-center rounded-full border-4 border-background bg-primary/10 text-2xl font-semibold text-primary ring-1 ring-border sm:size-28"
            aria-label="Avatar do perfil"
          >
            <span v-if="initials">{{ initials }}</span>
            <UserRoundIcon v-else :size="42" />
          </div>
          <h2 class="mt-4 max-w-full truncate text-xl font-bold tracking-tight">
            {{ session.data.value?.displayName || 'Seu perfil' }}
          </h2>
          <p class="mt-1 text-sm text-muted-foreground">
            @{{ session.data.value?.username }}
          </p>
          <span
            v-if="session.data.value?.isAdmin"
            class="mt-3 inline-flex items-center gap-1.5 rounded-full border border-primary/20 bg-primary/10 px-3 py-1 text-xs font-medium text-primary"
          >
            <BadgeCheckIcon :size="14" />
            Administrador
          </span>
        </div>

        <form class="mt-6 space-y-4 border-t pt-5" @submit="onSubmit">
          <div class="grid gap-2">
            <Label for="displayName">Nome de exibição</Label>
            <Input
              id="displayName"
              v-model="displayName"
              v-bind="displayNameAttrs"
              autocomplete="name"
              placeholder="Seu nome"
              maxlength="100"
              :readonly="!isEditing || isUpdatingProfile"
              class="read-only:border-transparent read-only:bg-muted/50 read-only:shadow-none"
            />
            <span v-if="errors.displayName" class="text-sm text-destructive">
              {{ errors.displayName }}
            </span>
          </div>

          <div class="grid gap-2">
            <Label for="username">Nome de usuário</Label>
            <div class="relative">
              <span
                class="pointer-events-none absolute left-3 top-1/2 -translate-y-1/2 text-muted-foreground"
              >
                @
              </span>
              <Input
                id="username"
                v-model="username"
                v-bind="usernameAttrs"
                class="pl-8 read-only:border-transparent read-only:bg-muted/50 read-only:shadow-none"
                autocomplete="username"
                placeholder="seu_usuario"
                maxlength="32"
                :readonly="!isEditing || isUpdatingProfile"
              />
            </div>
            <span v-if="isEditing" class="text-xs text-muted-foreground">
              Letras, números e underscore; mínimo de 3 caracteres.
            </span>
            <span v-if="errors.username" class="text-sm text-destructive">
              {{ errors.username }}
            </span>
          </div>

          <div class="flex flex-wrap gap-2 pt-1">
            <Button
              v-if="!isEditing"
              type="button"
              variant="outline"
              class="w-full gap-2"
              @click="isEditing = true"
            >
              <PencilLineIcon :size="16" />
              Editar perfil
            </Button>
            <Button
              v-if="isEditing"
              type="submit"
              class="w-full gap-2 sm:w-auto"
              :disabled="isUpdatingProfile"
            >
              <Spinner v-if="isUpdatingProfile" />
              <span>Salvar alterações</span>
            </Button>
            <Button
              v-if="isEditing"
              type="button"
              variant="outline"
              class="w-full sm:w-auto"
              :disabled="isUpdatingProfile"
              @click="cancelEditing"
            >
              Cancelar
            </Button>
          </div>
        </form>
      </CardContent>
    </Card>

    <div
      class="flex items-center gap-3 rounded-2xl border border-dashed bg-muted/30 p-4 text-xs text-muted-foreground"
    >
      <span
        class="flex size-10 shrink-0 items-center justify-center rounded-xl border bg-card text-primary shadow-sm"
      >
        <ShieldCheckIcon :size="20" />
      </span>
      <p class="leading-relaxed">
        Seus dados são protegidos e utilizados para apoiar sua participação
        cidadã.
      </p>
    </div>

    <Card id="historico" class="scroll-mt-6">
      <CardHeader>
        <CardTitle class="flex items-center gap-2">
          <MegaphoneIcon :size="20" class="text-primary" />
          Histórico de ocorrências
        </CardTitle>
        <CardDescription>
          Acompanhe as contribuições que você enviou.
        </CardDescription>
      </CardHeader>

      <CardContent>
        <div v-if="isLoadingIssues" class="flex justify-center py-10">
          <Spinner />
        </div>

        <div
          v-else-if="hasIssuesError"
          class="rounded-xl border border-dashed p-6 text-center"
        >
          <p class="text-sm text-destructive">
            Não foi possível carregar suas ocorrências.
          </p>
        </div>

        <div
          v-else-if="issues.length === 0"
          class="flex flex-col items-center rounded-xl border border-dashed px-6 py-10 text-center"
        >
          <div
            class="mb-3 flex size-12 items-center justify-center rounded-full bg-muted text-muted-foreground"
          >
            <MegaphoneIcon :size="21" />
          </div>
          <p class="font-medium">Você ainda não enviou ocorrências</p>
          <p class="mt-1 max-w-xs text-sm text-muted-foreground">
            Participe das discussões e ajude a construir melhorias para a
            cidade.
          </p>
          <NuxtLink to="/" class="mt-4">
            <Button variant="outline" size="sm">Explorar propostas</Button>
          </NuxtLink>
        </div>

        <div v-else class="flex flex-col divide-y">
          <article
            v-for="issue in issues"
            :key="issue.id"
            class="py-4 first:pt-0 last:pb-0"
          >
            <div class="flex items-start justify-between gap-3">
              <div class="min-w-0">
                <h3 class="line-clamp-2 font-medium">{{ issue.title }}</h3>
                <p class="mt-1 line-clamp-2 text-sm text-muted-foreground">
                  {{ issue.description }}
                </p>
              </div>
              <span
                class="shrink-0 rounded-full bg-primary/10 px-2.5 py-1 text-xs font-medium text-primary"
              >
                {{ statusLabel(issue.status) }}
              </span>
            </div>
            <div
              class="mt-3 flex items-center gap-2 text-xs text-muted-foreground"
            >
              <CalendarDaysIcon :size="14" />
              Enviada em {{ dayjs(issue.createdAt).format('DD/MM/YYYY HH:mm') }}
            </div>
            <NuxtImg
              v-if="issue.images?.length"
              :src="`${config.public.cdnUrl}/issues/${issue.id}/${issue.images[0]}`"
              :alt="issue.title"
              class="mt-3 max-h-48 w-full rounded-xl object-cover"
            />
          </article>
        </div>

        <Button
          v-if="hasMoreIssues"
          variant="outline"
          class="mt-5 w-full"
          :disabled="isLoadingMoreIssues"
          @click="fetchNextIssues()"
        >
          <LoaderCircleIcon
            v-if="isLoadingMoreIssues"
            class="mr-2 animate-spin"
            :size="16"
          />
          Carregar mais
        </Button>
      </CardContent>
    </Card>
  </div>
</template>
