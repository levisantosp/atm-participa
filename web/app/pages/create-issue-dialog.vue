<script setup lang="ts">
import { PlusIcon } from '@lucide/vue'
import type { InfiniteData } from '@tanstack/vue-query'
import { toTypedSchema } from '@vee-validate/zod'
import { getIssuesInfiniteQueryKey, useCreateIssue } from 'api-client'
import type { GetIssuesStatus200 } from 'api-client'
import {
  Button,
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
  Input,
  Label,
  Spinner,
  Textarea
} from 'ui'
import { useForm } from 'vee-validate'
import { toast } from 'vue-sonner'
import { z } from 'zod'
import { queryClient } from '~/lib/query-client'

const allowedImageTypes = new Set([
  'image/png',
  'image/jpeg',
  'image/webp',
  'image/gif'
])

const createIssueSchema = z.object({
  title: z
    .string()
    .min(3, 'O título deve ter pelo menos 3 caracteres')
    .max(72, 'O título deve ter no máximo 72 caracteres'),
  description: z
    .string()
    .min(10, 'A descrição deve ter pelo menos 10 caracteres')
    .max(65_000, 'A descrição deve ter no máximo 65.000 caracteres'),
  file: z
    .instanceof(File)
    .refine(
      (file) => allowedImageTypes.has(file.type),
      'A imagem deve ser PNG, JPEG, WebP ou GIF'
    )
    .refine(
      (file) => file.size <= 15_000_000,
      'A imagem deve ter no máximo 15 MB'
    )
    .optional()
})

const { defineField, errors, handleSubmit, resetForm, setFieldValue } = useForm(
  {
    validationSchema: toTypedSchema(createIssueSchema)
  }
)

const [title, titleAttrs] = defineField('title')
const [description, descriptionAttrs] = defineField('description')
const imageInput = ref<HTMLInputElement | null>(null)
const isCreateIssueDialogOpen = ref(false)

const createIssueMutation = useCreateIssue({
  mutation: {
    onSuccess(issue) {
      queryClient.setQueryData<InfiniteData<GetIssuesStatus200>>(
        getIssuesInfiniteQueryKey({ query: { limit: 100 } }),
        (currentData) => {
          if (!currentData?.pages.length) {
            return currentData
          }

          return {
            ...currentData,
            pages: currentData.pages.map((page, index) => {
              const items = page.items?.filter((item) => item.id !== issue.id)

              if (index === 0) {
                return {
                  ...page,
                  items: [issue, ...(items ?? [])]
                }
              }

              return items ? { ...page, items } : page
            })
          }
        }
      )

      toast.success('Ocorrência criada com sucesso')
      isCreateIssueDialogOpen.value = false
    }
  }
})
const { isPending: isCreatingIssue } = createIssueMutation

const createIssue = handleSubmit((formData) => {
  createIssueMutation.mutate({
    body: {
      title: formData.title,
      description: formData.description,
      ...(formData.file ? { file: formData.file } : {})
    }
  })
})

function setImage(event: Event) {
  const input = event.currentTarget

  if (input instanceof HTMLInputElement) {
    setFieldValue('file', input.files?.[0])
  }
}

watch(isCreateIssueDialogOpen, (isOpen) => {
  if (!isOpen) {
    resetForm()

    if (imageInput.value) {
      imageInput.value.value = ''
    }
  }
})
</script>

<template>
  <Dialog v-model:open="isCreateIssueDialogOpen">
    <div class="flex justify-end">
      <DialogTrigger as-child>
        <Button type="button">
          <PlusIcon />
          <span class="sr-only sm:not-sr-only">Nova ocorrência</span>
        </Button>
      </DialogTrigger>
    </div>

    <DialogContent class="max-h-[calc(100vh-2rem)] overflow-y-auto sm:max-w-lg">
      <DialogHeader>
        <DialogTitle>Nova ocorrência</DialogTitle>
        <DialogDescription>
          Descreva o problema que você encontrou. A imagem é opcional.
        </DialogDescription>
      </DialogHeader>

      <form id="create-issue" @submit="createIssue">
        <div class="flex flex-col gap-5">
          <div class="grid gap-2">
            <Label for="issue-title">Título</Label>
            <Input
              id="issue-title"
              v-model="title"
              v-bind="titleAttrs"
              maxlength="72"
              placeholder="Ex.: Buraco na rua"
              :aria-invalid="!!errors.title"
            />
            <span v-if="errors.title" class="text-sm text-red-400">
              {{ errors.title }}
            </span>
          </div>

          <div class="grid gap-2">
            <Label for="issue-description">Descrição</Label>
            <Textarea
              id="issue-description"
              v-model="description"
              v-bind="descriptionAttrs"
              maxlength="65000"
              placeholder="Conte mais detalhes sobre o problema"
              :aria-invalid="!!errors.description"
            />
            <span v-if="errors.description" class="text-sm text-red-400">
              {{ errors.description }}
            </span>
          </div>

          <div class="grid gap-2">
            <Label for="issue-image">Imagem (opcional)</Label>
            <Input
              id="issue-image"
              ref="imageInput"
              type="file"
              accept="image/png,image/jpeg,image/webp,image/gif"
              :aria-invalid="!!errors.file"
              @change="setImage"
            />
            <span class="text-muted-foreground text-xs">
              PNG, JPEG, WebP ou GIF; até 15 MB.
            </span>
            <span v-if="errors.file" class="text-sm text-red-400">
              {{ errors.file }}
            </span>
          </div>
        </div>
      </form>

      <DialogFooter>
        <Button
          type="button"
          variant="outline"
          :disabled="isCreatingIssue"
          @click="isCreateIssueDialogOpen = false"
        >
          Cancelar
        </Button>
        <Button
          type="submit"
          form="create-issue"
          class="w-40"
          :disabled="isCreatingIssue"
        >
          <Spinner v-if="isCreatingIssue" />
          <span v-else>Enviar ocorrência</span>
        </Button>
      </DialogFooter>
    </DialogContent>
  </Dialog>
</template>
