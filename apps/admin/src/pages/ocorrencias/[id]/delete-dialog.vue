<script setup lang="ts">
import { Trash } from '@lucide/vue'
import {
  getAdminIssuesByIdQueryKey,
  useDeleteAdminIssuesById,
  type GetAdminIssuesByIdStatus200
} from 'api-client'
import {
  Button,
  Dialog,
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
  Spinner
} from 'ui'
import { useRouter } from 'vue-router'
import { toast } from 'vue-sonner'
import { queryClient } from '@/lib/query'

const props = defineProps<{ update: GetAdminIssuesByIdStatus200 }>()
const router = useRouter()

const { isPending, mutate } = useDeleteAdminIssuesById({
  mutation: {
    onSuccess() {
      toast.success('Ocorrência deletada com sucesso')
      queryClient.removeQueries({
        queryKey: getAdminIssuesByIdQueryKey({
          path: {
            id: BigInt(props.update.id)
          }
        })
      })
      void router.push('/ocorrencias')
    }
  }
})
</script>

<template>
  <Dialog>
    <DialogTrigger as-child>
      <Button variant="destructive" class="relative cursor-pointer">
        <Trash />
        <span class="hidden md:block">Deletar</span>
      </Button>
    </DialogTrigger>

    <DialogContent>
      <DialogHeader>
        <DialogTitle>Deletar ocorrência</DialogTitle>

        <DialogDescription>
          Tem certeza que deseja deletar esta ocorrência? Esta ação é
          irreverssível
        </DialogDescription>

        <DialogFooter>
          <DialogClose>
            <Button variant="secondary" :disabled="isPending">Cancelar</Button>
          </DialogClose>

          <Button
            :disabled="isPending"
            @click="
              mutate({
                path: {
                  id: BigInt(props.update.id)
                }
              })
            "
            class="relative"
            variant="destructive"
          >
            <span :class="{ 'opacity-0': isPending }">Deletar</span>
            <Spinner v-if="isPending" class="absolute inset-0 m-auto" />
          </Button>
        </DialogFooter>
      </DialogHeader>
    </DialogContent>
  </Dialog>
</template>
