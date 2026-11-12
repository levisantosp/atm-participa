<script setup lang="ts">
import { Search, ThumbsUp } from '@lucide/vue'
import type { InfiniteData } from '@tanstack/vue-query'
import { toTypedSchema } from '@vee-validate/zod'
import {
  getIssuesInfiniteQueryKey,
  useGetIssuesInfinite,
  usePostIssuesByIdUpvote,
  type GetIssuesStatus200,
  type Issue
} from 'api-client'
import {
  Button,
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
  Input,
  Spinner
} from 'ui'
import { useForm } from 'vee-validate'
import { z } from 'zod'
import { queryClient } from '~/lib/query-client'
import CreateIssueDialog from './create-issue-dialog.vue'

definePageMeta({
  key: (route) => route.query.search?.toString() ?? 'search-empty'
})

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

const { isFetching, error, data } = useGetIssuesInfinite({
  query: {
    limit: 100n,
    search: route.query.search?.toString()
  }
})

const config = useRuntimeConfig()

const issues = computed(
  () => data.value?.pages.flatMap((page) => page.items ?? []) ?? []
)

const { isPending, mutate, variables } = usePostIssuesByIdUpvote({
  mutation: {
    onSuccess(data) {
      queryClient.setQueryData<InfiniteData<GetIssuesStatus200>>(
        getIssuesInfiniteQueryKey({
          query: {
            limit: 100n,
            search: route.query.search?.toString()
          }
        }),
        (old) => {
          if (!old) return old

          return {
            ...old,
            pages: old.pages.map((page) => ({
              ...page,
              items: page.items
                ? page.items.map((i) =>
                    i.id === data.id
                      ? {
                          ...i,
                          hasUpvoted: true,
                          upvotes: data.upvotes
                        }
                      : i
                  )
                : page.items
            }))
          }
        }
      )
    }
  }
})

function onUpvote(issue: Issue) {
  if (!issue.hasUpvoted) {
    mutate({
      path: {
        id: BigInt(issue.id)
      }
    })
  }
}
</script>

<template>
  <div class="mx-auto flex flex-col w-full max-w-5xl gap-5 px-4">
    <Spinner v-if="isFetching || !data" class="self-center" />

    <span v-else-if="error" class="text-red-400">
      Ocorreu um erro inesperado...
    </span>

    <div v-else class="flex w-full flex-col gap-5">
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

      <Card v-for="issue in issues" :key="issue.id" class="w-full">
        <CardHeader>
          <CardTitle>{{ issue.title }}</CardTitle>
          <CardDescription>{{ issue.description }}</CardDescription>
        </CardHeader>

        <CardContent class="flex flex-col gap-4">
          <NuxtImg
            v-if="issue.images?.length"
            :src="`${config.public.cdnUrl}/issues/${issue.id}/${issue.images[0]}`"
            :alt="issue.title"
            class="aspect-video w-full rounded-xl object-cover"
          />

          <Button
            variant="ghost"
            class="self-end"
            @click="onUpvote(issue)"
            :disabled="isPending && variables?.path.id === BigInt(issue.id)"
          >
            <Spinner
              v-if="isPending && variables?.path.id === BigInt(issue.id)"
            />
            <div v-else class="flex gap-1 items-center">
              <ThumbsUp :class="issue.hasUpvoted ? `fill-white` : ''" />
              <span>
                {{ issue.upvotes }}
              </span>
            </div>
          </Button>
        </CardContent>
      </Card>

      <Card v-if="!issues.length" class="w-full">
        <CardContent class="py-6 text-center text-muted-foreground">
          Nenhuma ocorrência encontrada para a pesquisa.
        </CardContent>
      </Card>
    </div>
  </div>
</template>
