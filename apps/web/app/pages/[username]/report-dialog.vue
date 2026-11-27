<script setup lang="ts">
import { toTypedSchema } from '@vee-validate/zod'
import { useReportUser, type User } from 'api-client'
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

const userId = defineModel<User['id'] | null>('userId', {
  required: true
})

const isDialogOpen = computed({
  get: () => userId.value !== null,
  set(open) {
    if (!open) {
      userId.value = null
    }
  }
})

const reasons = [
  'Assédio ou ameaça',
  'Perfil se passando por outra pessoa',
  'Spam ou golpes',
  'Nome de usuário ofensivo',
  'Discurso de ódio ou discriminação'
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

const { isPending, mutate } = useReportUser({
  mutation: {
    onSuccess() {
      toast.success(
        'Denúncia enviada com sucesso! Obrigado por nos ajudar a manter o ambiente agradável.'
      )
      isDialogOpen.value = false
    }
  }
})

watch(isDialogOpen, (isOpen) => {
  if (!isOpen) {
    selected.value = undefined
    resetForm()
  }
})

const onSubmit = handleSubmit((data) => {
  if (!userId.value) return

  mutate({
    path: {
      userId: BigInt(userId.value)
    },
    body: data
  })
})
</script>

<template>
  <Dialog v-model:open="isDialogOpen">
    <DialogContent>
      <DialogTitle>Denunciar usuário</DialogTitle>

      <form id="report-user" @submit.prevent="onSubmit">
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
              <div
                v-for="(r, i) in reasons"
                :key="r"
                class="flex items-center gap-1"
              >
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
          form="report-user"
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
