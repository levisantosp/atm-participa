import { z } from 'zod'

export const env = z
  .object({
    VITE_CDN_URL: z.url()
  })
  .parse(import.meta.env)
