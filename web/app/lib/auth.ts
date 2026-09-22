import { queryOptions, useQuery } from '@tanstack/vue-query'
import {
  getAuthMeQueryKey,
  getAuthMeQueryOptions,
  postAuthSignInEmail,
  postAuthSignOut,
  postAuthSignUpEmail
} from 'api-client'
import { queryClient } from './query-client'

const query = queryOptions({
  ...getAuthMeQueryOptions(),
  staleTime: 5 * 60_000
})

export const auth = {
  signIn: {
    async email(email: string, password: string) {
      await postAuthSignInEmail({
        body: {
          email,
          password
        }
      })
      await queryClient.invalidateQueries({ queryKey: getAuthMeQueryKey() })
    }
  },
  signUp: {
    async email(body: {
      email: string
      password: string
      username: string
      displayName: string
    }) {
      await postAuthSignUpEmail({ body })
      await queryClient.invalidateQueries({ queryKey: getAuthMeQueryKey() })
    }
  },
  async fetchSession() {
    return await queryClient.query(query).catch(() => null)
  },
  async signOut() {
    await postAuthSignOut()
    await queryClient.invalidateQueries({ queryKey: getAuthMeQueryKey() })
  },
  useSession() {
    return useQuery({
      ...query,
      staleTime: 5 * 60_000
    })
  }
} as const
