<script setup lang="ts">
import { ThumbsUp } from '@lucide/vue'
import type { InfiniteData } from '@tanstack/vue-query'
import {
  getIssues,
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
  Spinner
} from 'ui'
import { queryClient } from '~/lib/query-client'
import { isDark } from '~/lib/theme'

const route = useRoute()

const {
  isPending,
  error,
  data,
  fetchNextPage,
  hasNextPage,
  isFetchingNextPage
} = useGetIssuesInfinite(
  {
    query: computed(() => ({
      limit: 100n,
      search: route.query.search?.toString()
    }))
  },
  {
    query: {
      queryFn: ({ queryKey, pageParam, signal }) =>
        getIssues({
          query: {
            ...toValue(queryKey[1]),
            cursor: typeof pageParam === 'number' ? pageParam : undefined
          },
          signal
        }).unwrap()
    }
  }
)

const config = useRuntimeConfig()

const issues = computed(
  () => data.value?.pages.flatMap((page) => page.items ?? []) ?? []
)

const {
  isPending: isUpvoting,
  mutate,
  variables
} = usePostIssuesByIdUpvote({
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
  <Spinner v-if="isPending || !data" class="self-center" />

  <span v-else-if="error" class="text-red-400">
    Ocorreu um erro inesperado...
  </span>

  <div v-else v-for="issue in issues" :key="issue.id">
    <Card class="w-full">
      <CardHeader>
        <div class="flex flex-wrap gap-1" v-if="issue.author">
          <span>{{ issue.author.displayName }}</span>
          <span class="text-muted-foreground">
            @{{ issue.author.username }}
          </span>
        </div>
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
          :disabled="isUpvoting && variables?.path.id === BigInt(issue.id)"
        >
          <Spinner
            v-if="isUpvoting && variables?.path.id === BigInt(issue.id)"
          />
          <div v-else class="flex gap-1 items-center">
            <ThumbsUp
              :class="
                issue.hasUpvoted
                  ? `${isDark ? 'fill-white' : 'fill-black'}`
                  : ''
              "
            />
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

  <Button
    v-if="!isPending"
    :disabled="!hasNextPage || isFetchingNextPage"
    @click="fetchNextPage"
  >
    <Spinner v-if="isFetchingNextPage" />
    <span v-else>Carregar mais</span>
  </Button>
</template>
