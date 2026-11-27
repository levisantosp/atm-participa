<script setup lang="ts">
import { toTypedSchema } from '@vee-validate/zod'
import { useReportIssue, type Issue } from 'api-client'
import {
  Button,
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogTitle,
  Label,
  RadioGroup,
  RadioGroupItem,
  Spinner,
  Textarea
} from 'ui'
import { useForm } from 'vee-validate'
import { toast } from 'vue-sonner'
import { z } from 'zod'

const issueId = defineModel<Issue['id'] | null>('issueId', {
  required: true
})

const isDialogOpen = computed({
  get: () => issueId.value !== null,
  set(open) {
    if (!open) {
      issueId.value = null
    }
  }
})

const reasons = [
  'Assédio ou ameaça',
  'Spam ou propaganda',
  'Conteúdo impróprio',
  'Conteúdo ofensivo',
  'Informação falsa ou enganosa'
] as const

const schema = z.object({
  reason: z.string('Campo obrigatório').min(5, 'Mínimo 5 caracteres')
})

const { defineField, errors, handleSubmit, resetForm, setFieldValue } = useForm(
  {
    validationSchema: toTypedSchema(schema)
  }
)

const [reason, reasonProps] = defineField('reason')
const selected = ref<string>()

const { isPending, mutate } = useReportIssue({
  mutation: {
    onSuccess() {
      toast.success(
        'Denúncia enviada com sucesso! Obrigado por nos ajudar a manter o ambiente agradável.'
      )
      resetForm()
      issueId.value = null
      isDialogOpen.value = false
    }
  }
})

watch(isDialogOpen, (isOpen) => {
  if (!isOpen) {
    resetForm()
  }
})

const onSubmit = handleSubmit((data) =>
  mutate({
    path: {
      id: BigInt(issueId.value ?? -1)
    },
    body: data
  })
)
</script>

<template>
  <Dialog v-model:open="isDialogOpen">
    <DialogContent>
      <DialogTitle>Denunciar ocorrência</DialogTitle>
      <DialogDescription>Ocorrência #{{ issueId }}</DialogDescription>

      <form id="report-issue" @submit.prevent="onSubmit">
        <div class="flex flex-col gap-5">
          <div class="grid gap-2">
            <RadioGroup
              v-model="selected"
              @update:model-value="
                (v) =>
                  setFieldValue(
                    'reason',
                    v === 'other' ? undefined : v?.toString()
                  )
              "
            >
              <div v-for="(r, i) in reasons" class="flex items-center gap-1">
                <RadioGroupItem :id="i.toString()" :value="r" />
                <Label :for="i.toString()">{{ r }}</Label>
              </div>

              <div class="flex items-center gap-1">
                <RadioGroupItem id="other" value="other" />
                <Label for="other">Outro</Label>
              </div>
            </RadioGroup>

            <div v-if="selected === 'other'" class="grid gap-2 mt-2">
              <Label>Descreva o motivo</Label>
              <Textarea v-model="reason" v-bind="reasonProps" />
            </div>

            <span v-if="errors.reason" class="text-sm text-red-400">
              {{ errors.reason }}
            </span>
          </div>
        </div>
      </form>

      <DialogFooter>
        <Button
          type="button"
          variant="outline"
          :disabled="isPending"
          @click="isDialogOpen = false"
        >
          Cancelar
        </Button>

        <Button
          variant="destructive"
          form="report-issue"
          class="w-40"
          :disabled="isPending"
        >
          <Spinner v-if="isPending" />
          <span v-else>Enviar denúncia</span>
        </Button>
      </DialogFooter>
    </DialogContent>
  </Dialog>
</template>
