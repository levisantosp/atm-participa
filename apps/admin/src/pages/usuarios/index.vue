<script setup lang="ts">
import { Search } from '@lucide/vue'
import { toTypedSchema } from '@vee-validate/zod'
import { getAdminUsers, useGetAdminUsersInfinite } from 'api-client'
import dayjs from 'dayjs'
import {
  Button,
  Card,
  Input,
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
  search: z
    .string()
    .trim()
    .min(1, 'Mínimo 1 caractere')
    .max(100, 'Máximo 100 caracteres')
    .optional()
})

const { defineField, errors, handleSubmit, setFieldValue } = useForm({
  validationSchema: toTypedSchema(schema)
})

const [search, searchProps] = defineField('search')
const route = useRoute('/usuarios/')
const router = useRouter()

watch(
  () => route.query.search?.toString(),
  (search) => setFieldValue('search', search),
  {
    immediate: true
  }
)

const onSubmit = handleSubmit((data) => {
  void router.push({
    query: {
      ...route.query,
      search: data.search || undefined
    }
  })
})

const {
  isFetching,
  isFetchingNextPage,
  hasNextPage,
  fetchNextPage,
  error,
  data
} = useGetAdminUsersInfinite(
  {
    query: computed(() => ({
      search: route.query.search?.toString(),
      limit: 20n
    }))
  },
  {
    query: {
      queryFn: ({ queryKey, pageParam, signal }) =>
        getAdminUsers({
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

const users = computed(
  () => data.value?.pages.flatMap((page) => page.items ?? []) ?? []
)
</script>

<template>
  <h1 class="font-bold text-xl md:text-2xl">Usuários</h1>

  <form class="mt-10 w-full max-w-md md:max-w-4xl" @submit.prevent="onSubmit">
    <div class="flex min-w-0 gap-2">
      <Input
        id="user-search"
        v-model="search"
        v-bind="searchProps"
        type="search"
        placeholder="Pesquisar por nome de usuário, ID ou e-mail"
      />
      <Button
        type="submit"
        variant="secondary"
        size="icon"
        aria-label="Pesquisar usuários"
      >
        <Search />
      </Button>
    </div>
    <p v-if="errors.search" class="mt-1 text-sm text-red-400">
      {{ errors.search }}
    </p>
  </form>

  <Card class="mt-4 max-w-md md:max-w-4xl w-full">
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead>ID</TableHead>
          <TableHead>Nome</TableHead>
          <TableHead>Nome de usuário</TableHead>
          <TableHead>Banido</TableHead>
          <TableHead class="text-right">Cadastro</TableHead>
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
        <TableRow v-else-if="users.length === 0">
          <TableCell class="text-center text-muted-foreground" colspan="5">
            Nenhum usuário encontrado
          </TableCell>
        </TableRow>
        <template v-else>
          <TableRow
            v-for="user in users"
            :key="user.id"
            class="cursor-pointer"
            @click="router.push(`/usuarios/${user.id}`)"
          >
            <TableCell>{{ user.id }}</TableCell>
            <TableCell class="max-w-80 truncate">
              <RouterLink :to="`/usuarios/${user.id}`" @click.stop>
                {{ user.displayName }}
              </RouterLink>
            </TableCell>
            <TableCell>@{{ user.username }}</TableCell>
            <TableCell>{{ user.banned ? 'Sim' : 'Não' }}</TableCell>
            <TableCell class="text-right">
              {{ dayjs(user.createdAt).format('DD/MM/YYYY HH:mm') }}
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
