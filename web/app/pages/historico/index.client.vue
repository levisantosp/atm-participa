<script setup lang="ts">
import {
  CalendarDaysIcon,
  CheckCircle2Icon,
  ChevronRightIcon,
  CircleDotIcon,
  Clock3Icon,
  FileTextIcon,
  HeartIcon,
  MapPinIcon
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
  Spinner,
  Tabs,
  TabsContent,
  TabsList,
  TabsTrigger
} from 'ui'
import { auth } from '~/lib/auth'

definePageMeta({ layout: 'default' })

const session = auth.useSession()

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

    <Tabs default-value="created" class="w-full">
      <TabsList>
        <TabsTrigger value="created">Cadastradas</TabsTrigger>
        <TabsTrigger value="supported">Apoiadas</TabsTrigger>
      </TabsList>

      <TabsContent value="created">
        <Card>
          <CardHeader
            class="gap-4 sm:flex-row sm:items-center sm:justify-between"
          >
            <div>
              <CardTitle>Ocorrências cadastradas</CardTitle>
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
              v-else-if="issues.length === 0"
              class="flex flex-col items-center rounded-xl border border-dashed px-6 py-12 text-center"
            >
              <span
                class="mb-3 flex size-12 items-center justify-center rounded-full bg-muted text-muted-foreground"
                ><FileTextIcon :size="21"
              /></span>
              <p class="font-medium">Seu histórico ainda está vazio</p>
              <p class="mt-1 max-w-sm text-sm text-muted-foreground">
                Quando você cadastrar uma ocorrência, ela aparecerá aqui para
                acompanhamento.
              </p>
            </div>
            <div v-else class="flex flex-col divide-y">
              <article
                v-for="issue in issues"
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
                      <component
                        :is="statusInfo(issue.status).icon"
                        :size="13"
                      />{{ statusInfo(issue.status).label }}
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
      </TabsContent>

      <TabsContent value="supported">
        <Card>
          <CardContent
            class="flex flex-col items-center px-6 py-12 text-center"
          >
            <span
              class="mb-3 flex size-12 items-center justify-center rounded-full bg-rose-500/10 text-rose-300"
              ><HeartIcon :size="21"
            /></span>
            <p class="font-medium">Histórico de apoios indisponível</p>
            <p class="mt-1 max-w-md text-sm text-muted-foreground">
              A plataforma ainda não disponibiliza a consulta das ocorrências
              apoiadas por você. Seus apoios continuam registrados, mas não
              podem ser listados nesta tela.
            </p>
          </CardContent>
        </Card>
      </TabsContent>
    </Tabs>
  </div>
</template>
