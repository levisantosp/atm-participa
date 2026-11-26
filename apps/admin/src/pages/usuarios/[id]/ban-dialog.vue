<script setup lang="ts">
import { Ban, CalendarDays } from '@lucide/vue'
import { toTypedSchema } from '@vee-validate/zod'
import {
  getAdminUsersByIdQueryKey,
  usePostAdminUsersByIdBan,
  type GetAdminUsersByIdStatus200,
} from 'api-client'
import dayjs from 'dayjs'
import {
  Button,
  Calendar,
  Dialog,
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
  Label,
  NativeSelect,
  NativeSelectOption,
  Popover,
  PopoverContent,
  PopoverTrigger,
  Spinner,
  Switch
} from 'ui'
import { useForm } from 'vee-validate'
import { computed, ref, shallowRef, watch } from 'vue'
import { toast } from 'vue-sonner'
import { z } from 'zod'
import { auth } from '@/lib/auth'
import { queryClient } from '@/lib/query'

const props = defineProps<{
  user: GetAdminUsersByIdStatus200
}>()

const open = ref(false)
const calendarOpen = ref(false)
const selectedDate = shallowRef<InstanceType<typeof Calendar>['modelValue']>()
const visibleDate = shallowRef<InstanceType<typeof Calendar>['placeholder']>()
const monthFormatter = new Intl.DateTimeFormat('pt-BR', { month: 'short' })

function yearsFor(year: number) {
  const current = dayjs().year()
  const start = Math.min(current, year)
  const end = Math.max(current + 10, year + 1)
  return Array.from({ length: end - start + 1 }, (_, index) => start + index)
}

watch(calendarOpen, (isOpen) => {
  if (isOpen && !Array.isArray(selectedDate.value)) {
    visibleDate.value = selectedDate.value ?? undefined
  }
})
const { data: session } = auth.useSession()
const canBan = computed(
  () => !!session.value && session.value.userId !== props.user.id
)

const schema = z
  .object({
    permanent: z.boolean(),
    banEndsAt: z.iso.date().optional()
  })
  .refine(
    (data) =>
      data.permanent ||
      (data.banEndsAt && dayjs(data.banEndsAt).endOf('day').isAfter(dayjs())),
    {
      path: ['banEndsAt'],
      message: 'Selecione hoje ou uma data futura.'
    }
  )

const { defineField, errors, handleSubmit, resetForm } = useForm({
  validationSchema: toTypedSchema(schema),
  initialValues: {
    permanent: false
  }
})
const [permanent] = defineField('permanent')
const [banEndsAt] = defineField('banEndsAt')

const { isPending, mutate } = usePostAdminUsersByIdBan({
  mutation: {
    onSuccess(data, { path }) {
      queryClient.setQueryData(getAdminUsersByIdQueryKey({ path }), data)
      open.value = false
      toast.success('Usuário banido com sucesso')
    }
  }
})

const onSubmit = handleSubmit((data) => {
  if (!canBan.value || isPending.value) return
  mutate({
    path: {
      id: BigInt(props.user.id)
    },
    body:
      !data.permanent && data.banEndsAt
        ? {
            banEndsAt: dayjs(data.banEndsAt).endOf('day').toISOString()
          }
        : undefined
  })
})

watch(open, (isOpen) => {
  if (!isOpen) {
    resetForm()
    selectedDate.value = undefined
    calendarOpen.value = false
  }
})
</script>

<template>
  <Dialog v-model:open="open">
    <DialogTrigger as-child>
      <Button
        variant="destructive"
        class="relative cursor-pointer"
        aria-label="Banir usuário"
        :disabled="!canBan"
      >
        <Ban />
        <span class="hidden md:block">Banir usuário</span>
      </Button>
    </DialogTrigger>

    <DialogContent>
      <DialogHeader>
        <DialogTitle>Banir usuário</DialogTitle>
        <DialogDescription>
          Defina a duração do banimento de @{{ user.username }}.
        </DialogDescription>
      </DialogHeader>

      <form class="flex flex-col gap-6" @submit.prevent="onSubmit">
        <div class="flex items-center gap-2">
          <Switch
            id="ban-permanent"
            v-model="permanent"
            :disabled="isPending"
          />
          <Label for="ban-permanent">Sem prazo de término</Label>
        </div>
        <div v-if="!permanent" class="flex flex-col gap-2">
          <Label for="ban-ends-at">Último dia do banimento</Label>
          <Popover v-model:open="calendarOpen">
            <PopoverTrigger as-child>
              <Button
                id="ban-ends-at"
                type="button"
                variant="outline"
                class="justify-start"
                :disabled="isPending"
                :aria-invalid="!!errors.banEndsAt"
                :aria-describedby="
                  errors.banEndsAt ? 'ban-date-error' : undefined
                "
              >
                <CalendarDays />
                {{
                  banEndsAt
                    ? dayjs(banEndsAt).format('DD/MM/YYYY')
                    : 'Selecionar data'
                }}
              </Button>
            </PopoverTrigger>
            <PopoverContent align="start" class="w-auto p-0">
              <Calendar
                v-model="selectedDate"
                v-model:placeholder="visibleDate"
                locale="pt-BR"
                initial-focus
                :is-date-disabled="
                  (date) => dayjs(date.toString()).isBefore(dayjs(), 'day')
                "
                @update:model-value="
                  (date) => {
                    banEndsAt = date?.toString()
                    calendarOpen = false
                  }
                "
              >
                <template #calendar-heading="{ date }">
                  <div class="flex items-center gap-1">
                    <NativeSelect
                      size="sm"
                      aria-label="Mês"
                      :model-value="date.month"
                      @update:model-value="
                        visibleDate = date.set({ month: Number($event) })
                      "
                    >
                      <NativeSelectOption
                        v-for="month in 12"
                        :key="month"
                        :value="month"
                      >
                        {{
                          monthFormatter.format(
                            new Date(date.year, month - 1, 1)
                          )
                        }}
                      </NativeSelectOption>
                    </NativeSelect>
                    <NativeSelect
                      size="sm"
                      aria-label="Ano"
                      :model-value="date.year"
                      @update:model-value="
                        visibleDate = date.set({ year: Number($event) })
                      "
                    >
                      <NativeSelectOption
                        v-for="year in yearsFor(date.year)"
                        :key="year"
                        :value="year"
                      >
                        {{ year }}
                      </NativeSelectOption>
                    </NativeSelect>
                  </div>
                </template>
              </Calendar>
            </PopoverContent>
          </Popover>
          <span
            v-if="errors.banEndsAt"
            id="ban-date-error"
            class="text-sm text-red-400"
          >
            {{ errors.banEndsAt }}
          </span>
        </div>

        <DialogFooter>
          <DialogClose as-child>
            <Button type="button" variant="secondary" :disabled="isPending">
              Cancelar
            </Button>
          </DialogClose>
          <Button
            type="submit"
            variant="destructive"
            :disabled="!canBan || isPending"
            class="relative"
          >
            <span :class="{ 'opacity-0': isPending }">Banir</span>
            <Spinner v-if="isPending" class="absolute inset-0 m-auto" />
          </Button>
        </DialogFooter>
      </form>
    </DialogContent>
  </Dialog>
</template>
