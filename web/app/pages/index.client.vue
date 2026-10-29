<script setup lang="ts">
import { PlusIcon } from '@lucide/vue'
import { useGetIssuesInfinite } from 'api-client'
import {
  Button,
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
  Spinner
} from 'ui'

const { isFetching, error, data } = useGetIssuesInfinite({
  query: {
    limit: 100n
  }
})

const config = useRuntimeConfig()

const issues = computed(
  () => data.value?.pages.flatMap((page) => page.items ?? []) ?? []
)
</script>

<template>
  <div class="mx-auto flex flex-col w-full max-w-5xl gap-5 px-4">
    <div class="flex w-full justify-end">
      <Button as-child>
        <NuxtLink to="/criar-demanda">
          <PlusIcon />
          <span>Cadastrar demanda</span>
        </NuxtLink>
      </Button>
    </div>
    <Spinner v-if="isFetching || !data" class="self-center" />

    <span v-else-if="error" class="text-red-400">
      Ocorreu um erro inesperado...
    </span>

    <div v-else class="flex w-full flex-col gap-5">
      <Card v-for="issue in issues" :key="issue.id" class="w-full">
        <CardHeader>
          <CardTitle>{{ issue.title }}</CardTitle>
          <CardDescription>{{ issue.description }}</CardDescription>
        </CardHeader>

        <CardContent v-if="issue.images?.length">
          <NuxtImg
            :src="`${config.public.cdnUrl}/issues/${issue.id}/${issue.images[0]}`"
            :alt="issue.title"
            class="aspect-video w-full rounded-xl object-cover"
          />
        </CardContent>
      </Card>
    </div>
  </div>
</template>
