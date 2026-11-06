<script setup lang="ts">
import { Search, ThumbsUp } from '@lucide/vue'
import type { InfiniteData } from '@tanstack/vue-query'
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
import { queryClient } from '~/lib/query-client'
import CreateIssueDialog from './create-issue-dialog.vue'

const search = ref('')

const { isFetching, error, data } = useGetIssuesInfinite({
  query: computed(() => ({
    limit: 100n,
    search: search.value
  }))
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
            search: search.value
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
      <div class="relative w-full">
        <Search
          class="pointer-events-none absolute left-3 top-1/2 size-4 -translate-y-1/2 text-muted-foreground"
        />
        <Input
          v-model="search"
          type="search"
          placeholder="Pesquisar ocorrências"
          class="pl-9"
        />
      </div>

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
