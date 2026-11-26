<script setup lang="ts">
import {
  ArrowLeftIcon,
  CalendarDaysIcon,
  CheckCircle2Icon,
  CircleDotIcon,
  Clock3Icon,
  ImageIcon,
  MapPinIcon,
  ThumbsUpIcon,
  UserRoundIcon
} from '@lucide/vue'
import {
  useGetIssuesById,
  useGetUsersByUserIdUpvotes,
  usePostIssuesByIdUpvote
} from 'api-client'
import dayjs from 'dayjs'
import { Button, Card, CardContent, CardHeader, CardTitle, Spinner } from 'ui'
import { auth } from '~/lib/auth'

definePageMeta({ layout: 'default' })

const route = useRoute('/ocorrencias/[id]')
const router = useRouter()
const config = useRuntimeConfig()
const session = auth.useSession()
const issueId = computed(() => {
  const id = route.params.id
  return /^\d+$/.test(id) ? BigInt(id) : null
})

const {
  data: issue,
  isPending,
  isError,
  error,
  refetch
} = useGetIssuesById(
  { path: computed(() => ({ id: issueId.value ?? 0n })) },
  { query: { enabled: computed(() => issueId.value !== null) } }
)

const { data: upvotes, isFetching: isFetchingUpvotes } =
  useGetUsersByUserIdUpvotes(
    {
      path: computed(() => ({
        userId: BigInt(session.data.value?.userId ?? 0)
      })),
      query: computed(() => ({
        issuesIds: issue.value ? [BigInt(issue.value.id)] : []
      }))
    },
    {
      query: {
        enabled: computed(() =>
          Boolean(session.data.value?.userId && issue.value)
        )
      }
    }
  )

const hasUpvoted = computed(
  () =>
    upvotes.value?.items?.some(
      (upvote) => upvote.issueId === issue.value?.id
    ) ?? false
)
const { isPending: isUpvoting, mutate } = usePostIssuesByIdUpvote()

const statusInfo = computed(() => {
  const statuses: Record<
    string,
    { label: string; icon: typeof CircleDotIcon; class: string }
  > = {
    open: {
      label: 'Aberta',
      icon: CircleDotIcon,
      class: 'bg-sky-500/10 text-sky-300 ring-sky-400/20'
    },
    in_review: {
      label: 'Em análise',
      icon: Clock3Icon,
      class: 'bg-amber-500/10 text-amber-300 ring-amber-400/20'
    },
    closed: {
      label: 'Concluída',
      icon: CheckCircle2Icon,
      class: 'bg-emerald-500/10 text-emerald-300 ring-emerald-400/20'
    }
  }
  return issue.value
    ? (statuses[issue.value.status] ?? {
        label: issue.value.status,
        icon: CircleDotIcon,
        class: 'bg-muted text-muted-foreground ring-border'
      })
    : null
})

const isNotFound = computed(() => {
  const requestError = error.value as {
    status?: number
    statusCode?: number
  } | null
  return (
    issueId.value === null ||
    requestError?.status === 404 ||
    requestError?.statusCode === 404
  )
})

function supportIssue() {
  if (issue.value && session.data.value?.userId && !hasUpvoted.value) {
    mutate({ path: { id: BigInt(issue.value.id) } })
  }
}
</script>

<template>
  <div class="mx-auto flex w-full max-w-4xl flex-col gap-5 px-4 pb-12">
    <header class="flex items-center gap-3">
      <Button
        variant="ghost"
        size="icon"
        aria-label="Voltar"
        @click="router.back()"
      >
        <ArrowLeftIcon :size="20" />
      </Button>
      <h1 class="text-xl font-bold tracking-tight sm:text-2xl">
        Detalhes da ocorrência
      </h1>
    </header>

    <div
      v-if="isPending"
      class="flex justify-center py-16"
      role="status"
      aria-label="Carregando ocorrência"
    >
      <Spinner />
    </div>

    <Card v-else-if="isNotFound" class="border-dashed">
      <CardContent class="flex flex-col items-center gap-3 py-12 text-center">
        <p class="font-medium">Ocorrência não encontrada</p>
        <p class="text-sm text-muted-foreground">
          Este registro não existe ou não está disponível.
        </p>
        <Button variant="outline" @click="router.push('/')"
          >Voltar para ocorrências</Button
        >
      </CardContent>
    </Card>

    <Card v-else-if="isError" class="border-dashed">
      <CardContent class="flex flex-col items-center gap-3 py-12 text-center">
        <p class="text-sm text-destructive">
          Não foi possível carregar esta ocorrência.
        </p>
        <Button variant="outline" @click="refetch()">Tentar novamente</Button>
      </CardContent>
    </Card>

    <template v-else-if="issue">
      <Card>
        <CardHeader class="gap-4">
          <div class="flex flex-wrap items-center gap-2">
            <span
              v-if="statusInfo"
              class="inline-flex items-center gap-1 rounded-full px-2.5 py-1 text-xs font-medium ring-1 ring-inset"
              :class="statusInfo.class"
            >
              <component :is="statusInfo.icon" :size="13" />{{
                statusInfo.label
              }}
            </span>
          </div>
          <CardTitle class="text-2xl sm:text-3xl">{{ issue.title }}</CardTitle>
          <div
            class="flex flex-wrap gap-x-5 gap-y-2 text-sm text-muted-foreground"
          >
            <span class="inline-flex items-center gap-1.5">
              <UserRoundIcon :size="16" />
              <template v-if="issue.author">
                {{ issue.author.displayName }}
                <span>@{{ issue.author.username }}</span>
              </template>
              <span v-else>Autor não informado</span>
            </span>
            <span class="inline-flex items-center gap-1.5"
              ><CalendarDaysIcon :size="16" />{{
                dayjs(issue.createdAt).format('DD/MM/YYYY [às] HH:mm')
              }}</span
            >
            <span class="inline-flex items-center gap-1.5"
              ><MapPinIcon :size="16" />Bairro e localização não
              informados</span
            >
          </div>
        </CardHeader>
        <CardContent class="space-y-6">
          <section aria-labelledby="description-heading">
            <h2 id="description-heading" class="mb-2 font-semibold">
              Descrição
            </h2>
            <p
              class="whitespace-pre-wrap leading-relaxed text-muted-foreground"
            >
              {{ issue.description }}
            </p>
          </section>

          <section aria-labelledby="evidence-heading">
            <h2 id="evidence-heading" class="mb-3 font-semibold">Evidências</h2>
            <div v-if="issue.images?.length" class="grid gap-3 sm:grid-cols-2">
              <NuxtImg
                v-for="(image, index) in issue.images"
                :key="image"
                :src="`${config.public.cdnUrl}/issues/${issue.id}/${image}`"
                :alt="`Evidência ${index + 1} da ocorrência: ${issue.title}`"
                class="aspect-video w-full rounded-xl object-cover"
              />
            </div>
            <div
              v-else
              class="flex items-center gap-3 rounded-xl border border-dashed p-5 text-sm text-muted-foreground"
            >
              <ImageIcon :size="20" />Nenhuma fotografia foi anexada.
            </div>
          </section>

          <section
            aria-labelledby="support-heading"
            class="flex flex-wrap items-center justify-between gap-4 border-t pt-5"
          >
            <div>
              <h2 id="support-heading" class="font-semibold">
                Apoie esta ocorrência
              </h2>
              <p class="text-sm text-muted-foreground">
                {{ issue.upvotes.toString() }}
                {{ issue.upvotes === 1n ? 'apoio' : 'apoios' }}
              </p>
            </div>
            <Button
              v-if="session.data.value?.userId"
              :disabled="hasUpvoted || isFetchingUpvotes || isUpvoting"
              @click="supportIssue"
            >
              <Spinner v-if="isFetchingUpvotes || isUpvoting" />
              <template v-else
                ><ThumbsUpIcon :size="16" />{{
                  hasUpvoted ? 'Você já apoiou' : 'Apoiar'
                }}</template
              >
            </Button>
            <Button v-else variant="outline" @click="router.push('/sign-in')"
              >Entre para apoiar</Button
            >
          </section>
          <p
            v-if="!session.data.value?.userId && !session.isPending.value"
            class="text-xs text-muted-foreground"
          >
            O apoio fica associado à sua conta. A API atual não oferece
            cancelamento.
          </p>
        </CardContent>
      </Card>
    </template>
  </div>
</template>
