<script setup lang="ts">
import type { InfiniteData } from '@tanstack/vue-query'
import { toTypedSchema } from '@vee-validate/zod'
import { getIssuesInfiniteQueryKey, useCreateIssue } from 'api-client'
import type { GetIssuesStatus200 } from 'api-client'
import {
  Button,
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
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

const schema = z.object({
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
    validationSchema: toTypedSchema(schema)
  }
)

const [title, titleAttrs] = defineField('title')
const [description, descriptionAttrs] = defineField('description')
const imageInput = ref<HTMLInputElement | null>(null)
const hasCreatedDemand = ref(false)

const createDemandMutation = useCreateIssue({
  mutation: {
    onSuccess(issue) {
      queryClient.setQueryData<InfiniteData<GetIssuesStatus200>>(
        getIssuesInfiniteQueryKey({
          query: {
            limit: 100n
          }
        }),
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

      hasCreatedDemand.value = true
      toast.success('Demanda cadastrada com sucesso')
    },
    onError() {
      toast.error('Não foi possível cadastrar a demanda', {
        description: 'Tente novamente em instantes.'
      })
    }
  }
})
const { isPending: isCreatingDemand } = createDemandMutation

const createDemand = handleSubmit((formData) => {
  createDemandMutation.mutate({
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

function createAnotherDemand() {
  resetForm()

  if (imageInput.value) {
    imageInput.value.value = ''
  }

  hasCreatedDemand.value = false
}
</script>

<template>
  <div class="mx-auto w-full max-w-3xl px-4 pb-10">
    <div class="mb-6">
      <h1 class="text-2xl font-semibold tracking-tight sm:text-3xl">
        Cadastrar demanda
      </h1>
      <p class="mt-2 text-muted-foreground">
        Descreva o problema e, se quiser, anexe uma fotografia.
      </p>
    </div>

    <Card>
      <CardHeader>
        <CardTitle>Dados da ocorrência</CardTitle>
        <CardDescription>
          Preencha o título e a descrição da demanda.
        </CardDescription>
      </CardHeader>

      <template v-if="hasCreatedDemand">
        <CardContent>
          <div
            role="status"
            class="rounded-xl border border-green-200 bg-green-50 p-5 text-green-950 dark:border-green-900 dark:bg-green-950/30 dark:text-green-100"
          >
            <div class="flex items-start gap-3">
              <span
                aria-hidden="true"
                class="flex size-7 shrink-0 items-center justify-center rounded-full bg-green-600 font-semibold text-white"
              >
                ✓
              </span>
              <div>
                <h2 class="font-semibold">Demanda cadastrada com sucesso!</h2>
                <p class="mt-1 text-sm opacity-85">
                  Seu relato foi enviado e já pode ser acompanhado pela
                  comunidade.
                </p>
              </div>
            </div>
          </div>
        </CardContent>

        <CardFooter class="flex-col-reverse justify-end gap-3 sm:flex-row">
          <Button type="button" variant="outline" @click="createAnotherDemand">
            Cadastrar outra demanda
          </Button>
          <NuxtLink to="/">
            <Button type="button">Voltar para a página inicial</Button>
          </NuxtLink>
        </CardFooter>
      </template>

      <form v-else id="create-demand" @submit="createDemand">
        <CardContent>
          <div class="grid gap-5">
            <div class="grid gap-2">
              <Label for="demand-title">Título</Label>
              <Input
                id="demand-title"
                v-model="title"
                v-bind="titleAttrs"
                placeholder="Ex.: Poste sem iluminação na praça"
                maxlength="72"
                required
                :disabled="isCreatingDemand"
                :aria-invalid="!!errors.title"
              />
              <span v-if="errors.title" class="text-sm text-red-400">
                {{ errors.title }}
              </span>
            </div>

            <div class="grid gap-2">
              <Label for="demand-description">Descrição</Label>
              <Textarea
                id="demand-description"
                v-model="description"
                v-bind="descriptionAttrs"
                placeholder="Explique o problema e como ele afeta a região."
                maxlength="65000"
                required
                :disabled="isCreatingDemand"
                :aria-invalid="!!errors.description"
                class="min-h-32"
              />
              <span v-if="errors.description" class="text-sm text-red-400">
                {{ errors.description }}
              </span>
            </div>

            <div class="grid gap-2">
              <Label for="demand-photo"
                >Fotografia do problema (opcional)</Label
              >
              <input
                id="demand-photo"
                ref="imageInput"
                type="file"
                accept="image/png,image/jpeg,image/webp,image/gif"
                :disabled="isCreatingDemand"
                :aria-invalid="!!errors.file"
                @change="setImage"
                class="bg-input/50 border-transparent focus-visible:border-ring focus-visible:ring-ring/30 aria-invalid:ring-destructive/20 dark:aria-invalid:ring-destructive/40 aria-invalid:border-destructive dark:aria-invalid:border-destructive/50 h-9 w-full min-w-0 rounded-3xl border px-3 py-1 text-base outline-none transition-[color,box-shadow,background-color] file:inline-flex file:h-7 file:border-0 file:bg-transparent file:text-sm file:font-medium file:text-foreground focus-visible:ring-3 aria-invalid:ring-3 disabled:pointer-events-none disabled:cursor-not-allowed disabled:opacity-50 md:text-sm"
              />
              <span class="text-muted-foreground text-xs">
                PNG, JPEG, WebP ou GIF; até 15 MB.
              </span>
              <span v-if="errors.file" class="text-sm text-red-400">
                {{ errors.file }}
              </span>
            </div>
          </div>
        </CardContent>

        <CardFooter class="flex-col-reverse justify-end gap-3 sm:flex-row">
          <Button
            type="button"
            variant="outline"
            :disabled="isCreatingDemand"
            @click="navigateTo('/')"
          >
            Cancelar
          </Button>
          <Button
            type="submit"
            form="create-demand"
            class="w-full sm:w-44"
            :disabled="isCreatingDemand"
          >
            <Spinner v-if="isCreatingDemand" />
            <span v-else>Enviar cadastro</span>
          </Button>
        </CardFooter>
      </form>
    </Card>
  </div>
</template>
