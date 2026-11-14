<script setup lang="ts">
import { ClipboardList } from '@lucide/vue'
import { useGetAdminIssuesCountOpen } from 'api-client'
import { Card, CardContent, CardHeader, CardTitle, Spinner } from 'ui'

const { isFetching, error, data } = useGetAdminIssuesCountOpen()
</script>

<template>
  <Card class="w-80 h-30">
    <CardHeader class="flex justify-between">
      <CardTitle>Ocorrências abertas</CardTitle>
      <ClipboardList class="text-green-700" />
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
