<script setup lang="ts">
import { Search } from '@lucide/vue'
import { toTypedSchema } from '@vee-validate/zod'
import { getAdminIssues, useGetAdminIssuesInfinite } from 'api-client'
import dayjs from 'dayjs'
import {
  Button,
  Card,
  Input,
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
  Spinner,
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow
} from 'ui'
import { useForm } from 'vee-validate'
import { computed, toValue, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { z } from 'zod'

const schema = z.object({
  search: z.string().min(2, 'Mínimo 2 caracteres').optional(),
  status: z.enum(['open', 'in_review', 'closed'])
})

const { defineField, errors, handleSubmit, setFieldValue } = useForm({
  validationSchema: toTypedSchema(schema)
})

const [search, searchProps] = defineField('search')
const [status] = defineField('status')

const route = useRoute('/ocorrencias/')
const router = useRouter()

watch(
  [() => route.query.search?.toString(), () => route.query.status?.toString()],
  ([search, status]) => {
    setFieldValue('search', search)
    setFieldValue('status', (status ?? 'open') as any)
  },
  {
    immediate: true
  }
)

const onSubmit = handleSubmit((data) => {
  void router.push({
    query: {
      ...route.query,
      ...data
    }
  })
})

const { isFetching, isFetchingNextPage, fetchNextPage, error, data } =
  useGetAdminIssuesInfinite(
    {
      query: computed(() => ({
        search: route.query.search?.toString(),
        status: route.query.status?.toString() ?? 'open',
        limit: 20n
      }))
    },
    {
      query: {
        queryFn: ({ queryKey, pageParam, signal }) =>
          getAdminIssues({
            query: {
              ...toValue(queryKey[1]),
              cursor: typeof pageParam === 'number' ? pageParam : undefined
            },
            signal
          }).unwrap()
      }
    }
  )

const issues = computed(
  () => data.value?.pages.flatMap((page) => page.items ?? []) ?? []
)
const hasNextCursor = computed(() => {
  const pages = data.value?.pages ?? []

  return !!pages[pages.length - 1]?.nextCursor
})
</script>

<template>
  <h1 class="font-bold text-xl md:text-2xl">Ocorrências</h1>

  <form
    class="mt-10 grid w-full gap-3 sm:grid-cols-[minmax(0,1fr)_12rem] md:max-w-4xl"
    @submit.prevent="onSubmit"
  >
    <div class="min-w-0">
      <div class="flex min-w-0 gap-2">
        <Input
          id="issue-search"
          v-model="search"
          v-bind="searchProps"
          type="search"
          placeholder="Pesquisar ocorrências por título e descrição"
        />
        <Button
          type="submit"
          variant="secondary"
          size="icon"
          aria-label="Pesquisar ocorrências"
        >
          <Search />
        </Button>
      </div>
      <p
        v-if="!errors.search?.includes('Invalid input')"
        class="mt-1 text-sm text-red-400"
      >
        {{ errors.search }}
      </p>
    </div>

    <div>
      <Select v-model="status">
        <SelectTrigger class="w-full">
          <SelectValue placeholder="Filtrar por status" />
        </SelectTrigger>
        <SelectContent>
          <SelectItem value="open">Aberta</SelectItem>
          <SelectItem value="in_review">Em análise</SelectItem>
          <SelectItem value="closed">Encerrada</SelectItem>
        </SelectContent>
      </Select>
    </div>
  </form>

  <Card class="mt-4 max-w-md md:max-w-4xl w-full">
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead>ID</TableHead>
          <TableHead>Título</TableHead>
          <TableHead class="text-right">Criada em</TableHead>
          <TableHead class="text-right">Atualizada em</TableHead>
        </TableRow>
      </TableHeader>

      <TableBody>
        <TableRow v-if="isFetching && !data">
          <TableCell class="h-20 text-center" colspan="4">
            <Spinner class="mx-auto" />
          </TableCell>
        </TableRow>
        <TableRow v-else-if="!data && error">
          <TableCell class="text-center text-destructive" colspan="4">
            Ocorreu um erro inesperado
          </TableCell>
        </TableRow>
        <TableRow v-else-if="issues.length === 0">
          <TableCell class="text-center text-muted-foreground" colspan="4">
            Nenhuma ocorrência encontrada
          </TableCell>
        </TableRow>
        <template v-else>
          <TableRow
            v-for="issue in issues"
            :key="issue.id"
            class="cursor-pointer"
            @click="router.push(`/ocorrencias/${issue.id}`)"
          >
            <TableCell>{{ issue.id }}</TableCell>
            <TableCell class="max-w-80 truncate">
              {{ issue.title }}
            </TableCell>
            <TableCell class="text-right">
              {{ dayjs(issue.createdAt).format('DD/MM/YYYY HH:mm') }}
            </TableCell>
            <TableCell class="text-right">
              {{ dayjs(issue.updatedAt).format('DD/MM/YYYY HH:mm') }}
            </TableCell>
          </TableRow>
          <TableRow v-if="error">
            <TableCell class="text-center text-destructive" colspan="4">
              Ocorreu um erro inesperado
            </TableCell>
          </TableRow>
        </template>
      </TableBody>
    </Table>
  </Card>

  <div class="mt-4 flex w-full max-w-md justify-center md:max-w-4xl">
    <Button
      type="button"
      :disabled="!hasNextCursor || isFetchingNextPage"
      @click="fetchNextPage()"
    >
      <Spinner v-if="isFetchingNextPage" />
      Carregar mais
    </Button>
  </div>
</template>
