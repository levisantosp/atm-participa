<script setup lang="ts">
import {
  EllipsisVertical,
  Share,
  SquareExclamationPoint,
  ThumbsUp
} from '@lucide/vue'
import type { InfiniteData } from '@tanstack/vue-query'
import {
  getIssues,
  getIssuesInfiniteQueryKey,
  useGetIssuesInfinite,
  useGetUsersByUserIdUpvotes,
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
  Menubar,
  MenubarContent,
  MenubarItem,
  MenubarMenu,
  MenubarTrigger,
  Spinner
} from 'ui'
import { toast } from 'vue-sonner'
import ReportDialog from '~/components/report-dialog.vue'
import { auth } from '~/lib/auth'
import { queryClient } from '~/lib/query-client'
import { isDark } from '~/lib/theme'

const route = useRoute()
const reportIssueId = ref<Issue['id'] | null>(null)

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

const { data: sessionData, isPending: isSessionPending } = auth.useSession()

const { data: upvotesData, isFetching: isFetchingUpvotes } =
  useGetUsersByUserIdUpvotes(
    {
      path: computed(() => ({
        userId: BigInt(sessionData.value?.userId ?? 0)
      })),
      query: computed(() => ({
        issuesIds: issues.value.slice(0, 100).map((issue) => BigInt(issue.id)),
        limit: 100n
      }))
    },
    {
      query: {
        enabled: computed(
          () =>
            Boolean(data.value) &&
            !isSessionPending.value &&
            Boolean(sessionData.value?.userId) &&
            issues.value.length > 0
        )
      }
    }
  )

const upvotedIssueIDs = computed(
  () => new Set(upvotesData.value?.items?.map((upvote) => upvote.issueId) ?? [])
)

const locallyUpvotedIssueIDs = ref(new Set<string>())

const issuesWithVotes = computed(() =>
  issues.value.map((issue) =>
    Object.assign({}, issue, {
      hasUpvoted:
        upvotedIssueIDs.value.has(issue.id) ||
        locallyUpvotedIssueIDs.value.has(issue.id)
    })
  )
)

const {
  isPending: isUpvoting,
  mutate,
  variables
} = usePostIssuesByIdUpvote({
  mutation: {
    onSuccess(data) {
      locallyUpvotedIssueIDs.value = new Set([
        ...locallyUpvotedIssueIDs.value,
        data.id
      ])

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

function onUpvote(issue: Issue & { hasUpvoted: boolean }) {
  if (!issue.hasUpvoted) {
    mutate({
      path: {
        id: BigInt(issue.id)
      }
    })
  }
}

async function copy(id: string) {
  await navigator.clipboard.writeText(
    `${window.location.origin}/ocorrencias/${id}`
  )
  toast.success('Link copiado para área de transferência')
}
</script>

<template>
  <Spinner v-if="isPending" class="self-center" />

  <span v-else-if="error" class="text-red-400">
    Ocorreu um erro inesperado...
  </span>

  <div v-else v-for="issue in issuesWithVotes" :key="issue.id">
    <Card class="w-full">
      <CardHeader>
        <div class="flex items-center justify-between" v-if="issue.author">
          <div class="flex flex-wrap gap-1">
            <span>{{ issue.author.displayName }}</span>
            <span class="text-muted-foreground">
              @{{ issue.author.username }}
            </span>
          </div>

          <Menubar class="border-none">
            <MenubarMenu>
              <MenubarTrigger class="cursor-pointer">
                <EllipsisVertical />
              </MenubarTrigger>

              <MenubarContent>
                <MenubarItem class="cursor-pointer" @click="copy(issue.id)">
                  <Share />
                  Copiar link
                </MenubarItem>

                <MenubarItem
                  class="cursor-pointer text-destructive hover:bg-destructive!"
                  @select="reportIssueId = issue.id"
                >
                  <SquareExclamationPoint />
                  Reportar
                </MenubarItem>
              </MenubarContent>
            </MenubarMenu>
          </Menubar>
        </div>

        <NuxtLink :href="`/ocorrencias/${issue.id}`">
          <CardTitle
            class="cursor-pointer hover:text-blue-400 transition duration-200"
          >
            {{ issue.title }}
          </CardTitle>
        </NuxtLink>
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
          :disabled="
            isFetchingUpvotes ||
            (isUpvoting && variables?.path.id === BigInt(issue.id))
          "
        >
          <Spinner
            v-if="
              isFetchingUpvotes ||
              (isUpvoting && variables?.path.id === BigInt(issue.id))
            "
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

  <ReportDialog v-model:issue-id="reportIssueId" />

  <Button
    v-if="!isPending"
    :disabled="!hasNextPage || isFetchingNextPage"
    @click="fetchNextPage"
  >
    <Spinner v-if="isFetchingNextPage" />
    <span v-else>Carregar mais</span>
  </Button>
</template>
