<script setup lang="ts">
import { UserCheck } from '@lucide/vue'
import {
  getAdminUsersByIdQueryKey,
  useGetAdminUsersById,
  usePostAdminUsersByIdUnban
} from 'api-client'
import dayjs from 'dayjs'
import {
  Button,
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
  Spinner
} from 'ui'
import { computed, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { toast } from 'vue-sonner'
import { auth } from '@/lib/auth'
import { queryClient } from '@/lib/query'
import BanDialog from './ban-dialog.vue'

const route = useRoute('/usuarios/[id]/')
const router = useRouter()

const { data: session } = auth.useSession()

const { isFetching, data, error } = useGetAdminUsersById(
  {
    path: computed(() => ({
      id: BigInt(route.params.id)
    }))
  },
  {
    query: {
      meta: {
        handledErrorStatus: 404
      }
    }
  }
)

const { isPending, mutate } = usePostAdminUsersByIdUnban({
  mutation: {
    onSuccess(data, { path }) {
      queryClient.setQueryData(getAdminUsersByIdQueryKey({ path }), data)
      toast.success('Usuário desbanido com sucesso')
    }
  }
})

watch(error, (err) => {
  if (err?.status === 404) {
    toast.error('Usuário não encontrado')
    router.back()
  }
})
</script>

<template>
  <div v-if="isFetching && !data" class="flex justify-center items-center">
    <Spinner />
  </div>
  <div v-else-if="!data" class="mt-10 text-center">
    <p v-if="error && error.status !== 404">Ocorreu um erro inesperado</p>
    <RouterLink to="/usuarios" class="underline">
      Voltar para usuários
    </RouterLink>
  </div>
  <div v-else class="mt-10 max-w-md md:max-w-4xl w-full flex flex-col gap-10">
    <div class="flex justify-end gap-1">
      <Button
        v-if="data.banned"
        aria-label="Desbanir usuário"
        :disabled="isPending || !session || session.userId === data.id"
        @click="mutate({ path: { id: BigInt(data.id) } })"
      >
        <Spinner v-if="isPending" />
        <UserCheck />
        <span class="hidden md:block">Desbanir usuário</span>
      </Button>
      <BanDialog v-else :key="data.id" :user="data" />
    </div>

    <Card>
      <CardHeader>
        <CardTitle>{{ data.displayName }}</CardTitle>
        <CardDescription>@{{ data.username }}</CardDescription>
      </CardHeader>
      <CardContent>
        <dl class="flex flex-wrap gap-10">
          <div>
            <dt class="text-sm text-muted-foreground">Permissão</dt>
            <dd>{{ data.isAdmin ? 'Administrador' : 'Usuário' }}</dd>
          </div>
          <div>
            <dt class="text-sm text-muted-foreground">Cadastro</dt>
            <dd>{{ dayjs(data.createdAt).format('DD/MM/YYYY HH:mm') }}</dd>
          </div>
          <div>
            <dt class="text-sm text-muted-foreground">Última atualização</dt>
            <dd>{{ dayjs(data.updatedAt).format('DD/MM/YYYY HH:mm') }}</dd>
          </div>
          <div>
            <dt class="text-sm text-muted-foreground">Banido</dt>
            <dd>{{ data.banned ? 'Sim' : 'Não' }}</dd>
          </div>
          <div v-if="data.banned && data.bannedAt">
            <dt class="text-sm text-muted-foreground">Banido em</dt>
            <dd>{{ dayjs(data.bannedAt).format('DD/MM/YYYY HH:mm') }}</dd>
          </div>
          <div v-if="data.banned">
            <dt class="text-sm text-muted-foreground">Término do banimento</dt>
            <dd>
              {{
                data.banEndsAt
                  ? dayjs(data.banEndsAt).format('DD/MM/YYYY HH:mm')
                  : 'Sem prazo'
              }}
            </dd>
          </div>
        </dl>
      </CardContent>
    </Card>
  </div>
</template>
