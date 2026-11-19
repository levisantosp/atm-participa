<script setup lang="ts">
import { Trophy } from '@lucide/vue'
import { useGetIssues } from 'api-client'
import {
  Button,
  Card,
  CardContent,
  CardHeader,
  CardTitle,
  Spinner,
  Table,
  TableBody,
  TableCaption,
  TableCell,
  TableHead,
  TableHeader,
  TableRow
} from 'ui'

const { isPending, data, error } = useGetIssues({
  query: {
    orderBy: 'upvotes'
  }
})
</script>

<template>
  <Card>
    <CardHeader>
      <CardTitle class="text-center">Ocorrências mais votadas</CardTitle>
    </CardHeader>

    <CardContent>
      <Spinner v-if="isPending || !data" class="mx-auto" />

      <div v-else-if="error" class="text-center text-destructive">
        <span> Ocorreu um erro inesperado... </span>
      </div>

      <div v-else class="flex flex-col gap-10 items-center">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Posição</TableHead>
              <TableHead>Votos</TableHead>
              <TableHead>Título</TableHead>
            </TableRow>
          </TableHeader>

          <TableBody>
            <TableRow v-if="!data.items?.length">
              <TableCell class="text-center text-muted-foreground" colspan="3">
                Nenhuma ocorrência encontrada
              </TableCell>
            </TableRow>

            <TableRow
              v-else
              v-for="(issue, i) in data.items"
              :key="issue.id"
              class="text-muted-foreground"
            >
              <TableCell>{{ i + 1 }}</TableCell>
              <TableCell>{{ issue.upvotes }}</TableCell>
              <TableCell class="whitespace-normal wrap-break-word">
                {{ issue.title }}
              </TableCell>
            </TableRow>
          </TableBody>
        </Table>

        <NuxtLink href="/ranking">
          <Button>
            <Trophy />
            Ver ranking completo
          </Button>
        </NuxtLink>
      </div>
    </CardContent>
  </Card>
</template>
