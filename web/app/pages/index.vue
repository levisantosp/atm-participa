<script setup lang="ts">
import { SearchIcon, ThumbsUpIcon } from '@lucide/vue'
import { Card, CardContent, CardHeader, CardTitle, Input, Spinner } from 'ui'

definePageMeta({
  layout: 'default'
})

type Issue = {
  id: number
  title: string
  description: string
  status: 'open' | 'closed' | 'in_review'
  createdAt: string
  updatedAt: string
}

type IssuesResponse = {
  items: Issue[]
  hasNextPage: boolean
}

const search = ref('')

const { data, error, status } = await useFetch<IssuesResponse>(
  'http://localhost:3333/issues',
  {
    query: { limit: 100 },
    credentials: 'include',
    default: () => ({ items: [], hasNextPage: false })
  }
)

const issues = computed(() => data.value?.items ?? [])

const filteredIssues = computed(() => {
  const term = search.value.trim().toLocaleLowerCase()

  if (!term) return issues.value

  return issues.value.filter((issue) =>
    [issue.title, issue.description].some((value) =>
      value.toLocaleLowerCase().includes(term)
    )
  )
})

const statusLabels: Record<Issue['status'], string> = {
  open: 'Aberta',
  closed: 'Encerrada',
  in_review: 'Em análise'
}

function formatDate(date: string) {
  return new Intl.DateTimeFormat('pt-BR', {
    dateStyle: 'short'
  }).format(new Date(date))
}
</script>

<template>
  <div class="flex w-full max-w-5xl flex-col gap-8 pt-10">
    <header class="flex flex-col gap-2">
      <p class="text-sm text-muted-foreground">Altamira Participa</p>
      <h1 class="text-3xl font-semibold tracking-tight">Ocorrências</h1>
      <p class="text-muted-foreground">
        Consulte os problemas registrados pela comunidade.
      </p>
    </header>

    <div class="relative w-full">
      <SearchIcon
        class="pointer-events-none absolute left-3 top-1/2 size-4 -translate-y-1/2 text-muted-foreground"
      />
      <Input
        id="issue-search"
        v-model="search"
        class="pl-9"
        placeholder="Pesquisar por título ou descrição"
        type="search"
      />
    </div>

    <div v-if="status === 'pending'" class="flex justify-center py-16">
      <Spinner class="size-6" />
    </div>

    <Card v-else-if="error" class="p-6">
      <CardContent class="p-0 text-center text-muted-foreground">
        Não foi possível carregar as ocorrências.
      </CardContent>
    </Card>

    <div v-else-if="filteredIssues.length" class="grid gap-4">
      <div v-for="issue in filteredIssues" :key="issue.id" class="block">
        <Card class="transition-colors hover:bg-accent/50">
          <CardHeader class="gap-2">
            <div
              class="flex flex-col justify-between gap-2 sm:flex-row sm:items-start"
            >
              <CardTitle class="text-xl">{{ issue.title }}</CardTitle>
              <span class="shrink-0 text-sm text-muted-foreground">
                {{ statusLabels[issue.status] }}
              </span>
            </div>
          </CardHeader>
          <CardContent class="grid gap-4 pt-0">
            <p class="line-clamp-2 text-sm text-muted-foreground">
              {{ issue.description }}
            </p>
            <div
              class="flex flex-wrap items-center gap-x-6 gap-y-2 text-sm text-muted-foreground"
            >
              <span class="inline-flex items-center gap-2">
                <ThumbsUpIcon class="size-4" />
                Apoios não informados
              </span>
              <span>Data: {{ formatDate(issue.createdAt) }}</span>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>

    <Card v-else class="p-6">
      <CardContent class="p-0 text-center text-muted-foreground">
        {{
          search
            ? 'Nenhuma ocorrência encontrada.'
            : 'Nenhuma ocorrência cadastrada.'
        }}
      </CardContent>
    </Card>
  </div>
</template>
