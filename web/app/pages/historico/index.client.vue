<script setup lang="ts">
import {
  CalendarDaysIcon,
  CheckCircle2Icon,
  ChevronRightIcon,
  CircleDotIcon,
  Clock3Icon,
  FileTextIcon,
  HeartIcon,
  MapPinIcon,
  SearchIcon,
  XIcon
} from '@lucide/vue'
import { useGetUsersByUserIdIssuesInfinite } from 'api-client'
import dayjs from 'dayjs'
import {
  Button,
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
  Input,
  Spinner
} from 'ui'
import { auth } from '~/lib/auth'

definePageMeta({ layout: 'default' })

type HistoryTab = 'created' | 'supported'

const session = auth.useSession()
const config = useRuntimeConfig()
const activeTab = ref<HistoryTab>('created')
const search = ref('')
const selectedIssue = ref<{
  id: string
  title: string
  description: string
  status: string
  images?: Array<string | null> | null
  createdAt: string | Date
} | null>(null)

const issuePath = computed(() => ({
  userId: BigInt(session.data.value?.userId ?? '0')
}))

const {
  data,
  isPending,
  isError,
  hasNextPage,
  isFetchingNextPage,
  fetchNextPage,
  refetch
} = useGetUsersByUserIdIssuesInfinite(
  { path: issuePath, query: { limit: 20n } },
  { query: { enabled: computed(() => !!session.data.value) } }
)

const issues = computed(
  () => data.value?.pages.flatMap((page) => page.items ?? []) ?? []
)
const filteredIssues = computed(() => {
  const term = search.value.trim().toLocaleLowerCase('pt-BR')
  if (!term) return issues.value
  return issues.value.filter((issue) =>
    `${issue.title} ${issue.description}`
      .toLocaleLowerCase('pt-BR')
      .includes(term)
  )
})

const statusInfo = (status: string) => {
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
  return (
    statuses[status] ?? {
      label: status,
      icon: CircleDotIcon,
      class: 'bg-muted text-muted-foreground ring-border'
    }
  )
}
</script>

<template>
  <div class="mx-auto flex w-full max-w-5xl flex-col gap-6 px-4 pb-12">
    <header class="space-y-2">
      <div class="flex items-center gap-2">
        <CalendarDaysIcon :size="23" class="text-primary" />
        <h1 class="text-2xl font-bold tracking-tight">Meu histórico</h1>
      </div>
      <p class="max-w-2xl text-sm text-muted-foreground">
        Acompanhe as ocorrências que você registrou e suas contribuições na
        plataforma.
      </p>
    </header>

    <div class="grid gap-3 sm:grid-cols-2">
      <button
        type="button"
        class="flex items-center gap-3 rounded-xl border p-4 text-left transition-colors"
        :class="
          activeTab === 'created'
            ? 'border-primary/50 bg-primary/5 ring-1 ring-primary/20'
            : 'bg-card hover:bg-muted/50'
        "
        @click="activeTab = 'created'"
      >
        <span
          class="flex size-10 shrink-0 items-center justify-center rounded-lg bg-primary/10 text-primary"
          ><FileTextIcon :size="19"
        /></span>
        <span class="min-w-0 flex-1">
          <span class="block font-medium">Cadastradas</span>
          <span class="text-xs text-muted-foreground"
            >Suas ocorrências enviadas</span
          >
        </span>
        <span class="rounded-full bg-muted px-2.5 py-1 text-xs tabular-nums">{{
          issues.length
        }}</span>
      </button>
      <button
        type="button"
        class="flex items-center gap-3 rounded-xl border p-4 text-left transition-colors"
        :class="
          activeTab === 'supported'
            ? 'border-primary/50 bg-primary/5 ring-1 ring-primary/20'
            : 'bg-card hover:bg-muted/50'
        "
        @click="activeTab = 'supported'"
      >
        <span
          class="flex size-10 shrink-0 items-center justify-center rounded-lg bg-rose-500/10 text-rose-300"
          ><HeartIcon :size="19"
        /></span>
        <span class="min-w-0 flex-1">
          <span class="block font-medium">Apoiadas</span>
          <span class="text-xs text-muted-foreground"
            >Ocorrências que você apoiou</span
          >
        </span>
        <ChevronRightIcon :size="17" class="text-muted-foreground" />
      </button>
    </div>

    <Card v-if="activeTab === 'created'">
      <CardHeader class="gap-4 sm:flex-row sm:items-center sm:justify-between">
        <div>
          <CardTitle>Ocorrências cadastradas</CardTitle>
          <CardDescription class="mt-1"
            >{{ issues.length }}
            {{
              issues.length === 1
                ? 'ocorrência encontrada'
                : 'ocorrências encontradas'
            }}</CardDescription
          >
        </div>
        <div class="relative w-full sm:max-w-xs">
          <SearchIcon
            :size="16"
            class="absolute left-3 top-1/2 -translate-y-1/2 text-muted-foreground"
          />
          <Input
            v-model="search"
            class="pl-9"
            placeholder="Buscar no histórico"
            aria-label="Buscar ocorrências"
          />
        </div>
      </CardHeader>
      <CardContent>
        <div
          v-if="isPending"
          class="flex justify-center py-12"
          role="status"
          aria-label="Carregando histórico"
        >
          <Spinner />
        </div>
        <div
          v-else-if="isError"
          class="rounded-xl border border-dashed p-8 text-center"
        >
          <p class="text-sm text-destructive">
            Não foi possível carregar seu histórico agora.
          </p>
          <Button variant="outline" class="mt-4" @click="refetch()"
            >Tentar novamente</Button
          >
        </div>
        <div
          v-else-if="filteredIssues.length === 0"
          class="flex flex-col items-center rounded-xl border border-dashed px-6 py-12 text-center"
        >
          <span
            class="mb-3 flex size-12 items-center justify-center rounded-full bg-muted text-muted-foreground"
            ><SearchIcon v-if="search" :size="21" /><FileTextIcon
              v-else
              :size="21"
          /></span>
          <p class="font-medium">
            {{
              search
                ? 'Nenhuma ocorrência encontrada'
                : 'Seu histórico ainda está vazio'
            }}
          </p>
          <p class="mt-1 max-w-sm text-sm text-muted-foreground">
            {{
              search
                ? 'Tente buscar por outro título ou descrição.'
                : 'Quando você cadastrar uma ocorrência, ela aparecerá aqui para acompanhamento.'
            }}
          </p>
          <NuxtLink v-if="!search" to="/" class="mt-4"
            ><Button variant="outline" size="sm"
              >Ver ocorrências</Button
            ></NuxtLink
          >
        </div>
        <div v-else class="flex flex-col divide-y">
          <article
            v-for="issue in filteredIssues"
            :key="issue.id"
            class="flex flex-col gap-4 py-5 first:pt-0 last:pb-0 sm:flex-row sm:items-start"
          >
            <div class="min-w-0 flex-1">
              <div class="flex flex-wrap items-center gap-2">
                <h3 class="font-semibold">{{ issue.title }}</h3>
                <span
                  class="inline-flex items-center gap-1 rounded-full px-2.5 py-1 text-xs font-medium ring-1 ring-inset"
                  :class="statusInfo(issue.status).class"
                >
                  <component :is="statusInfo(issue.status).icon" :size="13" />{{
                    statusInfo(issue.status).label
                  }}
                </span>
              </div>
              <p class="mt-1 line-clamp-2 text-sm text-muted-foreground">
                {{ issue.description }}
              </p>
              <div
                class="mt-3 flex flex-wrap items-center gap-x-4 gap-y-2 text-xs text-muted-foreground"
              >
                <span class="inline-flex items-center gap-1.5"
                  ><CalendarDaysIcon :size="14" />Cadastrada em
                  {{ dayjs(issue.createdAt).format('DD/MM/YYYY') }}</span
                >
                <span class="inline-flex items-center gap-1.5"
                  ><MapPinIcon :size="14" />Bairro não informado</span
                >
              </div>
            </div>
            <Button
              variant="outline"
              size="sm"
              class="shrink-0"
              @click="selectedIssue = issue"
              >Ver detalhes<ChevronRightIcon :size="15" class="ml-1"
            /></Button>
          </article>
        </div>
        <Button
          v-if="hasNextPage && !isError"
          variant="outline"
          class="mt-5 w-full"
          :disabled="isFetchingNextPage"
          @click="fetchNextPage()"
        >
          <Spinner v-if="isFetchingNextPage" class="mr-2" />Carregar mais
        </Button>
      </CardContent>
    </Card>

    <Card v-else>
      <CardContent class="flex flex-col items-center px-6 py-12 text-center">
        <span
          class="mb-3 flex size-12 items-center justify-center rounded-full bg-rose-500/10 text-rose-300"
          ><HeartIcon :size="21"
        /></span>
        <p class="font-medium">Histórico de apoios indisponível</p>
        <p class="mt-1 max-w-md text-sm text-muted-foreground">
          A plataforma ainda não disponibiliza a consulta das ocorrências
          apoiadas por você. Seus apoios continuam registrados, mas não podem
          ser listados nesta tela.
        </p>
      </CardContent>
    </Card>

    <div
      v-if="selectedIssue"
      class="fixed inset-0 z-50 flex items-center justify-center bg-background/80 p-4 backdrop-blur-sm"
      @click.self="selectedIssue = null"
      @keydown.esc="selectedIssue = null"
    >
      <Card
        class="max-h-[90vh] w-full max-w-2xl overflow-y-auto shadow-xl"
        role="dialog"
        aria-modal="true"
        aria-labelledby="history-detail-title"
      >
        <CardHeader class="flex-row items-start justify-between gap-4">
          <div>
            <CardTitle id="history-detail-title">{{
              selectedIssue.title
            }}</CardTitle
            ><CardDescription class="mt-2"
              >Cadastrada em
              {{
                dayjs(selectedIssue.createdAt).format('DD/MM/YYYY [às] HH:mm')
              }}</CardDescription
            >
          </div>
          <Button
            variant="ghost"
            size="icon"
            aria-label="Fechar detalhes"
            @click="selectedIssue = null"
            ><XIcon :size="18"
          /></Button>
        </CardHeader>
        <CardContent>
          <span
            class="inline-flex items-center gap-1 rounded-full px-2.5 py-1 text-xs font-medium ring-1 ring-inset"
            :class="statusInfo(selectedIssue.status).class"
            ><component
              :is="statusInfo(selectedIssue.status).icon"
              :size="13"
            />{{ statusInfo(selectedIssue.status).label }}</span
          >
          <p
            class="mt-4 whitespace-pre-line text-sm leading-relaxed text-muted-foreground"
          >
            {{ selectedIssue.description }}
          </p>
          <p class="mt-4 flex items-center gap-2 text-xs text-muted-foreground">
            <MapPinIcon :size="14" />Bairro não informado no cadastro
          </p>
          <NuxtImg
            v-if="selectedIssue.images?.[0]"
            :src="`${config.public.cdnUrl}/issues/${selectedIssue.id}/${selectedIssue.images[0]}`"
            :alt="selectedIssue.title"
            class="mt-5 max-h-80 w-full rounded-xl object-cover"
          />
        </CardContent>
      </Card>
    </div>
  </div>
</template>
