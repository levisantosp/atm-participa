<script setup lang="ts">
import { Eye } from '@lucide/vue'
import { useGetAdminIssuesCountInReview } from 'api-client'
import { Card, CardContent, CardHeader, CardTitle, Spinner } from 'ui'

const { isFetching, error, data } = useGetAdminIssuesCountInReview()
</script>

<template>
  <Card class="w-80 h-30">
    <CardHeader class="flex justify-between">
      <CardTitle>Ocorrências em revisão</CardTitle>
      <Eye class="text-yellow-400" />
    </CardHeader>
    <CardContent class="text-muted-foreground">
      <Spinner v-if="isFetching || data === undefined" />
      <span v-else-if="error" class="text-destructive">
        Ocorreu um erro inesperado
      </span>
      <span>{{ data }}</span>
    </CardContent>
  </Card>
</template>
