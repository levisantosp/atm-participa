type PublicRoute = {
  action: 'redirect' | 'next'
}

export const publicRoutes = new Map<string, PublicRoute>([
  [
    '/sign-in',
    {
      action: 'redirect'
    }
  ],
  [
    '/sign-up',
    {
      action: 'redirect'
    }
  ]
])
