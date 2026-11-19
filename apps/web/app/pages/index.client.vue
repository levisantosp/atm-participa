<script setup lang="ts">
import { Search } from '@lucide/vue'
import { toTypedSchema } from '@vee-validate/zod'
import { Button, Input } from 'ui'
import { useForm } from 'vee-validate'
import { z } from 'zod'
import CreateIssueDialog from './create-issue-dialog.vue'
import Issues from './issues.vue'

const route = useRoute()
const router = useRouter()
const schema = z.object({
  search: z
    .string()
    .trim()
    .min(2, 'Digite pelo menos 2 caracteres para pesquisar.')
    .max(100, 'A pesquisa deve ter no máximo 100 caracteres.')
})

const { handleSubmit, errors, defineField, setFieldValue } = useForm({
  validationSchema: toTypedSchema(schema)
})

const [search, searchAttrs] = defineField('search')

watch(
  () => route.query.search,
  (v) => {
    setFieldValue('search', v?.toString())
  },
  {
    immediate: true
  }
)

const onSubmit = handleSubmit((data) => {
  void router.push({
    query: {
      ...route.query,
      ...data
    }
  })
})
</script>

<template>
  <div class="mx-auto flex flex-col w-full max-w-5xl gap-5 px-4">
    <div class="flex w-full flex-col gap-5">
      <form class="flex w-full flex-col gap-2" @submit.prevent="onSubmit">
        <div class="flex flex-wrap gap-2">
          <div class="relative min-w-0 flex-1">
            <Search
              class="pointer-events-none absolute left-3 top-1/2 size-4 -translate-y-1/2 text-muted-foreground"
            />
            <Input
              v-model="search"
              v-bind="searchAttrs"
              type="search"
              placeholder="Pesquisar ocorrências"
              class="pl-9"
            />
          </div>

          <Button aria-label="Pesquisar">
            <Search class="size-4" />
            <span class="hidden md:block">Pesquisar</span>
          </Button>
        </div>
        <p
          v-if="!errors.search?.startsWith('Invalid input')"
          class="text-sm text-red-400"
          role="alert"
        >
          {{ errors.search }}
        </p>
      </form>

      <CreateIssueDialog />
      <Issues />
    </div>
  </div>
</template>
