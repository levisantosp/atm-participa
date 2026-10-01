<script setup lang="ts">
import { PlusIcon } from '@lucide/vue'
import type { InfiniteData } from '@tanstack/vue-query'
import { useQueryClient } from '@tanstack/vue-query'
import { getIssuesInfiniteQueryKey, useGetIssuesInfinite } from 'api-client'
import type { GetIssuesStatus200, Issue } from 'api-client'
import {
  Button,
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
  Spinner
} from 'ui'
import CreateIssueDialog from './create-issue-dialog.vue'

const { isFetching, error, data } = useGetIssuesInfinite({
  query: {
    limit: 100
  }
})

const queryClient = useQueryClient()

function addIssueToFirstPage(issue: Issue) {
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
}

const issues = computed(
  () => data.value?.pages.flatMap((page) => page.items ?? []) ?? []
)
</script>

<template>
  <div class="mx-auto flex flex-col w-full max-w-5xl gap-5 px-4">
    <Spinner v-if="isFetching || !data" class="self-center" />

    <span v-else-if="error" class="text-red-400">
      Ocorreu um erro inesperado...
    </span>

    <div v-else class="flex w-full flex-col gap-5">
      <CreateIssueDialog @created="addIssueToFirstPage" />

      <Card v-for="issue in issues" :key="issue.id" class="w-full">
        <CardHeader>
          <CardTitle>{{ issue.title }}</CardTitle>
          <CardDescription>{{ issue.description }}</CardDescription>
        </CardHeader>

        <CardContent v-if="issue.imageUrl">
          <NuxtImg
            :src="issue.imageUrl"
            :alt="issue.title"
            class="aspect-video w-full rounded-xl object-cover"
          />
        </CardContent>
      </Card>
    </div>
  </div>
</template>
