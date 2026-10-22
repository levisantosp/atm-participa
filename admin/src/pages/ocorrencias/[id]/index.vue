<script setup lang="ts">
import { SearchAlert, SearchCheck } from '@lucide/vue'
import {
  getAdminIssuesByIdQueryKey,
  useGetAdminIssuesById,
  usePatchAdminIssuesById
} from 'api-client'
import {
  Button,
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
  Spinner
} from 'ui'
import { watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { toast } from 'vue-sonner'
import { env } from '@/env'
import { queryClient } from '@/lib/query'
import DeleteDialog from './delete-dialog.vue'

const route = useRoute('/ocorrencias/[id]/')
const router = useRouter()

const { isFetching, data, error } = useGetAdminIssuesById({
  path: {
    id: BigInt(route.params.id)
  }
})

const { isPending, mutate } = usePatchAdminIssuesById({
  mutation: {
    onSuccess(data) {
      queryClient.setQueryData(
        getAdminIssuesByIdQueryKey({
          path: {
            id: BigInt(route.params.id)
          }
        }),
        data
      )
    }
  }
})

watch(error, (err) => {
  if (err?.status === 404) {
    toast.error('Ocorrência não encontrada')
    void router.push('/ocorrencias')
  }
})
</script>

<template>
  <div v-if="isFetching || !data" class="flex justify-center items-center">
    <Spinner />
  </div>

  <div v-else class="mt-10 max-w-md md:max-w-4xl w-full flex flex-col gap-10">
    <div class="flex justify-end gap-1">
      <Button
        v-if="data.status === 'open'"
        @click="
          mutate({
            path: {
              id: BigInt(data.id)
            },
            body: {
              status: 'in_review'
            }
          })
        "
        :disabled="isPending"
      >
        <Spinner v-if="isPending" />
        <SearchAlert />
        <span class="hidden md:block">Colocar em revisão</span>
      </Button>

      <Button
        v-else-if="data.status === 'in_review'"
        @click="
          mutate({
            path: {
              id: BigInt(data.id)
            },
            body: {
              status: 'closed'
            }
          })
        "
        :disabled="isPending"
      >
        <Spinner v-if="isPending" />
        <SearchCheck />
        <span class="hidden md:block">Encerrar</span>
      </Button>

      <Button
        v-else
        @click="
          mutate({
            path: {
              id: BigInt(data.id)
            },
            body: {
              status: 'open'
            }
          })
        "
        :disabled="isPending"
      >
        <Spinner v-if="isPending" />
        <SearchCheck />
        <span class="hidden md:block">Reabrir</span>
      </Button>

      <DeleteDialog :update="data" />
    </div>

    <Card>
      <CardHeader>
        <CardTitle>{{ data.title }}</CardTitle>
        <CardDescription>{{ data.description }}</CardDescription>
      </CardHeader>

      <CardContent v-if="data.images?.length">
        <img
          :src="`${env.VITE_CDN_URL}/issues/${data.id}/${data.images[0]}`"
          :alt="data.title"
          class="w-full rounded-xl"
        />
      </CardContent>
    </Card>
  </div>
</template>
