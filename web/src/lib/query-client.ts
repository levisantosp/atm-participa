import { MutationCache, QueryCache, QueryClient } from '@tanstack/react-query'
import { getAuthMeQueryKey, ResponseError } from 'api-client'
import { toast } from 'sonner'

const ignoredKeys = new Set([JSON.stringify(getAuthMeQueryKey())])

export const queryClient = new QueryClient({
  queryCache: new QueryCache({
    onError(err, query) {
      if (ignoredKeys.has(JSON.stringify(query.queryKey))) return

      if (err instanceof ResponseError) {
        toast.error(err.response.data.detail)
      } else if (err instanceof Error) {
        toast.error(err.message)
      }
    }
  }),
  mutationCache: new MutationCache({
    onError(err) {
      if (err instanceof ResponseError) {
        toast.error(err.response.data.detail)
      } else if (err instanceof Error) {
        toast.error(err.message)
      }
    }
  }),
  defaultOptions: {
    queries: {
      staleTime: 60_000,
      retry: false
    },
    mutations: {
      retry: false
    }
  }
})
