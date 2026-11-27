<script setup lang="ts">
import { getAdminReports, useGetAdminReportsInfinite } from 'api-client'
import dayjs from 'dayjs'
import {
  Button,
  Card,
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
import { computed, toValue } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute('/denuncias/')
const router = useRouter()

const {
  isFetching,
  isFetchingNextPage,
  hasNextPage,
  fetchNextPage,
  error,
  data
} = useGetAdminReportsInfinite(
  {
    query: computed(() => ({
      entity:
        route.query.entity === 'issue' || route.query.entity === 'user'
          ? route.query.entity
          : undefined,
      limit: 20n
    }))
  },
  {
    query: {
      queryFn: ({ queryKey, pageParam, signal }) =>
        getAdminReports({
          query: {
            ...toValue(queryKey[1]),
            cursor:
              typeof pageParam === 'number' || typeof pageParam === 'bigint'
                ? BigInt(pageParam)
                : undefined
          },
          signal
        }).unwrap()
    }
  }
)

const reports = computed(
  () => data.value?.pages.flatMap((page) => page.items ?? []) ?? []
)
</script>

<template>
  <h1 class="font-bold text-xl md:text-2xl">Denúncias</h1>

  <div class="mt-10 w-full max-w-md md:max-w-4xl">
    <Select
      :model-value="route.query.entity?.toString() ?? 'all'"
      @update:model-value="
        router.push({
          query: {
            ...route.query,
            entity: $event === 'issue' || $event === 'user' ? $event : undefined
          }
        })
      "
    >
      <SelectTrigger
        class="w-full sm:w-48"
        aria-label="Filtrar por tipo de denúncia"
      >
        <SelectValue />
      </SelectTrigger>
      <SelectContent>
        <SelectItem value="all">Todos os tipos</SelectItem>
        <SelectItem value="issue">Ocorrências</SelectItem>
        <SelectItem value="user">Usuários</SelectItem>
      </SelectContent>
    </Select>
  </div>

  <Card class="mt-4 max-w-md md:max-w-4xl w-full">
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead>ID</TableHead>
          <TableHead>Motivo</TableHead>
          <TableHead>Tipo</TableHead>
          <TableHead>Autor</TableHead>
          <TableHead class="text-right">Criada em</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        <TableRow v-if="isFetching && !data">
          <TableCell class="h-20 text-center" colspan="5">
            <Spinner class="mx-auto" />
          </TableCell>
        </TableRow>
        <TableRow v-else-if="!data && error">
          <TableCell class="text-center text-destructive" colspan="5">
            Ocorreu um erro inesperado
          </TableCell>
        </TableRow>
        <TableRow v-else-if="reports.length === 0">
          <TableCell class="text-center text-muted-foreground" colspan="5">
            Nenhuma denúncia encontrada
          </TableCell>
        </TableRow>
        <template v-else>
          <TableRow
            v-for="report in reports"
            :key="report.id"
            class="cursor-pointer"
            @click="router.push(`/denuncias/${report.id}`)"
          >
            <TableCell>{{ report.id }}</TableCell>
            <TableCell class="max-w-80 truncate">
              <RouterLink :to="`/denuncias/${report.id}`" @click.stop>
                {{ report.reason }}
              </RouterLink>
            </TableCell>
            <TableCell>{{
              report.entity === 'issue' ? 'Ocorrência' : 'Usuário'
            }}</TableCell>
            <TableCell>{{
              report.author ? '@' + report.author.username : 'Usuário excluído'
            }}</TableCell>
            <TableCell class="text-right">
              {{ dayjs(report.createdAt).format('DD/MM/YYYY HH:mm') }}
            </TableCell>
          </TableRow>
          <TableRow v-if="error">
            <TableCell class="text-center text-destructive" colspan="5">
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
      :disabled="!hasNextPage || isFetchingNextPage"
      @click="fetchNextPage()"
    >
      <Spinner v-if="isFetchingNextPage" />
      Carregar mais
    </Button>
  </div>
</template>
