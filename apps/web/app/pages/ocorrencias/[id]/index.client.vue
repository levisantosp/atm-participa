<script setup lang="ts">
import {
  CalendarDaysIcon,
  CheckCircle2Icon,
  CircleDotIcon,
  Clock3Icon,
  ThumbsUp
} from '@lucide/vue'
import {
  getIssuesByIdQueryKey,
  getIssuesQueryKey,
  getUsersByUserIdUpvotesQueryKey,
  type Issue,
  useGetIssuesById,
  useGetUsersByUserIdUpvotes,
  usePostIssuesByIdUpvote
} from 'api-client'
import dayjs from 'dayjs'
import {
  Button,
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
  Spinner
} from 'ui'
import { toast } from 'vue-sonner'
import { auth } from '~/lib/auth'
import { queryClient } from '~/lib/query-client'

definePageMeta({
  layout: 'default'
})

const route = useRoute('ocorrencias-id')
const router = useRouter()
const config = useRuntimeConfig()
const { data: session } = auth.useSession()

const {
  data: issue,
  isPending,
  error
} = useGetIssuesById({
  path: computed(() => ({
    id: BigInt(route.params.id)
  }))
})

watch(error, (err) => {
  if (err?.status === 404) {
    toast.error('Ocorrência não encontrada')
    router.back()
  }
})

const { data: upvotes, isFetching: isFetchingUpvotes } =
  useGetUsersByUserIdUpvotes(
    {
      path: computed(() => ({
        userId: BigInt(session.value?.userId ?? 0)
      })),
      query: computed(() => ({
        issuesIds: issue.value ? [BigInt(issue.value.id)] : []
      }))
    },
    {
      query: {
        enabled: computed(() => Boolean(session.value?.userId && issue.value)),
        staleTime: 0
      }
    }
  )

const hasUpvoted = computed(
  () =>
    upvotes.value?.items?.some(
      (upvote) => upvote.issueId === issue.value?.id
    ) ?? false
)
const { isPending: isUpvoting, mutate } = usePostIssuesByIdUpvote({
  mutation: {
    onMutate: () => session.value?.userId,
    async onSuccess(updatedIssue, _, userId) {
      const id = BigInt(updatedIssue.id)
      queryClient.setQueryData<Issue>(
        getIssuesByIdQueryKey({ path: { id } }),
        (current) => current && { ...current, upvotes: updatedIssue.upvotes }
      )
      void queryClient.invalidateQueries({ queryKey: getIssuesQueryKey() })
      if (userId) {
        await queryClient.invalidateQueries({
          queryKey: getUsersByUserIdUpvotesQueryKey({
            path: {
              userId: BigInt(userId)
            }
          })
        })
      }
    }
  }
})

const statusInfo = computed(() => {
  const statuses: Record<
    string,
    { label: string; icon: typeof CircleDotIcon; class: string }
  > = {
    open: {
      label: 'Aberta',
      icon: CircleDotIcon,
      class: 'bg-sky-500/10 text-sky-700 dark:text-sky-300 ring-sky-400/20'
    },
    in_review: {
      label: 'Em análise',
      icon: Clock3Icon,
      class:
        'bg-amber-500/10 text-amber-700 dark:text-amber-300 ring-amber-400/20'
    },
    closed: {
      label: 'Concluída',
      icon: CheckCircle2Icon,
      class:
        'bg-emerald-500/10 text-emerald-700 dark:text-emerald-300 ring-emerald-400/20'
    }
  }
  return issue.value
    ? (statuses[issue.value.status] ?? {
        label: issue.value.status,
        icon: CircleDotIcon,
        class: 'bg-muted text-muted-foreground ring-border'
      })
    : null
})

function upvote() {
  if (issue.value && !hasUpvoted.value && !isUpvoting.value) {
    mutate({
      path: {
        id: BigInt(issue.value.id)
      }
    })
  }
}
</script>

<template>
  <div class="mx-auto flex w-full max-w-5xl flex-col gap-5 px-4 pb-12">
    <header>
      <h1 class="text-xl font-bold tracking-tight sm:text-2xl">
        Detalhes da ocorrência
      </h1>
    </header>

    <div v-if="isPending" class="flex justify-center py-16">
      <Spinner />
    </div>

    <Card v-else-if="issue" class="w-full">
      <CardHeader>
        <div class="flex flex-wrap items-center justify-between gap-2">
          <div v-if="issue.author" class="flex flex-wrap gap-1">
            <span>{{ issue.author.displayName }}</span>
            <span class="text-muted-foreground"
              >@{{ issue.author.username }}</span
            >
          </div>
          <span
            v-if="statusInfo"
            class="inline-flex items-center gap-1 rounded-full px-2.5 py-1 text-xs font-medium ring-1 ring-inset"
            :class="statusInfo.class"
          >
            <component :is="statusInfo.icon" :size="13" />
            {{ statusInfo.label }}
          </span>
        </div>
        <CardTitle>{{ issue.title }}</CardTitle>
        <CardDescription class="whitespace-pre-wrap">{{
          issue.description
        }}</CardDescription>
      </CardHeader>

      <CardContent class="flex flex-col gap-4">
        <div
          v-if="issue.images?.length"
          class="grid gap-3"
          :class="{ 'sm:grid-cols-2': issue.images.length > 1 }"
        >
          <NuxtImg
            v-for="(image, index) in issue.images"
            :key="image"
            :src="`${config.public.cdnUrl}/issues/${issue.id}/${image}`"
            :alt="`Foto ${index + 1} da ocorrência: ${issue.title}`"
            class="aspect-video w-full rounded-xl object-cover"
          />
        </div>
        <div class="flex items-center justify-between gap-2">
          <time
            :datetime="issue.createdAt"
            class="inline-flex items-center gap-1.5 text-xs text-muted-foreground"
          >
            <CalendarDaysIcon :size="14" />
            {{ dayjs(issue.createdAt).format('DD/MM/YYYY [às] HH:mm') }}
          </time>
          <Button
            variant="ghost"
            :disabled="isFetchingUpvotes || isUpvoting"
            :aria-disabled="isFetchingUpvotes || isUpvoting || hasUpvoted"
            :aria-pressed="hasUpvoted"
            :aria-label="`Apoiar ocorrência, total de apoios: ${issue.upvotes}`"
            @click="upvote"
          >
            <Spinner v-if="isFetchingUpvotes || isUpvoting" />
            <span v-else class="flex items-center gap-1">
              <ThumbsUp :class="{ 'fill-current': hasUpvoted }" />
              {{ issue.upvotes }}
            </span>
          </Button>
        </div>
      </CardContent>
    </Card>
  </div>
</template>
