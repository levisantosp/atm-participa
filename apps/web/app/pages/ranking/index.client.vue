<script setup lang="ts">
import { useGetIssuesInfinite } from 'api-client'
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

const {
  isPending,
  error,
  data,
  fetchNextPage,
  hasNextPage,
  isFetchingNextPage
} = useGetIssuesInfinite({
  query: {
    orderBy: 'upvotes',
    limit: 10n
  }
})

const issues = computed(
  () => data.value?.pages.flatMap((page) => page.items ?? []) ?? []
)

const router = useRouter()
</script>

<template>
  <div class="flex flex-col gap-5 w-full max-w-5xl">
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
                <TableHead class="text-right">Autor</TableHead>
              </TableRow>
            </TableHeader>

            <TableBody>
              <TableRow v-if="!issues.length">
                <TableCell
                  class="text-center text-muted-foreground"
                  colspan="4"
                >
                  Nenhuma ocorrência encontrada
                </TableCell>
              </TableRow>

              <TableRow
                v-else
                v-for="(issue, i) in issues"
                :key="issue.id"
                class="text-muted-foreground cursor-pointer"
                @click="router.push(`/ocorrencias/${issue.id}`)"
              >
                <TableCell>{{ i + 1 }}</TableCell>
                <TableCell>{{ issue.upvotes }}</TableCell>
                <TableCell class="whitespace-normal wrap-break-word">
                  {{ issue.title }}
                </TableCell>
                <TableCell class="text-right">
                  {{ issue.author ? `@${issue.author.username}` : '' }}
                </TableCell>
              </TableRow>
            </TableBody>
          </Table>
        </div>
      </CardContent>
    </Card>

    <Button
      v-if="!isPending"
      :disabled="!hasNextPage || isFetchingNextPage"
      @click="fetchNextPage"
    >
      <Spinner v-if="isFetchingNextPage" />
      <span v-else>Carregar mais</span>
    </Button>
  </div>
</template>
