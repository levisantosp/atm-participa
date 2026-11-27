<script setup lang="ts">
import { SearchCheck } from '@lucide/vue'
import {
  getAdminReportsByIdQueryKey,
  getAdminReportsInfiniteQueryKey,
  useGetAdminReportsById,
  type GetAdminReportsByIdStatus200,
  usePatchAdminReportsById
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
import { queryClient } from '@/lib/query'

const route = useRoute('/denuncias/[id]/')
const router = useRouter()

const { isFetching, data, error } = useGetAdminReportsById({
  path: computed(() => ({
    id: BigInt(route.params.id)
  }))
})

const { isPending, mutate } = usePatchAdminReportsById({
  mutation: {
    onSuccess(data, { path }) {
      queryClient.setQueryData<GetAdminReportsByIdStatus200>(
        getAdminReportsByIdQueryKey({ path }),
        (current) => ({
          ...current,
          ...data
        })
      )

      queryClient.invalidateQueries({
        queryKey: getAdminReportsInfiniteQueryKey()
      })
    }
  }
})

watch(error, (err) => {
  if (err?.status === 404) {
    toast.error('Denúncia não encontrada')
    void router.push('/denuncias')
  }
})
</script>

<template>
  <div v-if="isFetching && !data" class="flex justify-center items-center">
    <Spinner />
  </div>
  <div v-else-if="!data" class="mt-10 text-center">
    <p v-if="error && error.status !== 404">Ocorreu um erro inesperado</p>
    <RouterLink to="/denuncias" class="underline">
      Voltar para denúncias
    </RouterLink>
  </div>
  <div v-else class="mt-10 max-w-md md:max-w-4xl w-full flex flex-col gap-10">
    <div class="flex justify-end">
      <Button
        :disabled="isPending"
        @click="
          mutate({
            path: { id: BigInt(data.id) },
            body: { status: data.status === 'open' ? 'closed' : 'open' }
          })
        "
      >
        <Spinner v-if="isPending" />
        <SearchCheck v-else />
        <span class="sr-only md:not-sr-only">
          {{ data.status === 'open' ? 'Encerrar' : 'Reabrir' }}
        </span>
      </Button>
    </div>

    <Card>
      <CardHeader>
        <CardTitle>Denúncia #{{ data.id }}</CardTitle>
        <CardDescription class="whitespace-pre-wrap wrap-anywhere">{{
          data.reason
        }}</CardDescription>
      </CardHeader>
      <CardContent>
        <dl class="flex flex-wrap gap-10">
          <div>
            <dt class="text-sm text-muted-foreground">Alvo da denúncia</dt>
            <dd>
              <RouterLink
                v-if="data.targetId"
                :to="
                  data.entity === 'issue'
                    ? `/ocorrencias/${data.targetId}`
                    : `/usuarios/${data.targetId}`
                "
                class="underline"
              >
                {{ data.entity === 'issue' ? 'Ocorrência' : 'Usuário' }} #{{
                  data.targetId
                }}
              </RouterLink>
              <span v-else>{{
                data.entity === 'issue'
                  ? 'Ocorrência excluída'
                  : 'Usuário excluído'
              }}</span>
            </dd>
          </div>
          <div>
            <dt class="text-sm text-muted-foreground">Criada em</dt>
            <dd>{{ dayjs(data.createdAt).format('DD/MM/YYYY HH:mm') }}</dd>
          </div>
          <div>
            <dt class="text-sm text-muted-foreground">Última atualização</dt>
            <dd>{{ dayjs(data.updatedAt).format('DD/MM/YYYY HH:mm') }}</dd>
          </div>
          <div>
            <dt class="text-sm text-muted-foreground">Autor</dt>
            <dd>
              <RouterLink
                v-if="data.author"
                :to="`/usuarios/${data.author.id}`"
                class="underline"
              >
                @{{ data.author.username }}
              </RouterLink>
              <span v-else>Usuário excluído</span>
            </dd>
          </div>
        </dl>
      </CardContent>
    </Card>
  </div>
</template>
