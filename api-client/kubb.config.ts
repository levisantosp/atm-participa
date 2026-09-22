import { pluginAxios } from '@kubb/plugin-axios'
import { pluginTs } from '@kubb/plugin-ts'
import { pluginVueQuery } from '@kubb/plugin-vue-query'
import { pluginZod } from '@kubb/plugin-zod'
import { defineConfig } from 'kubb/config'

process.loadEnvFile()

export default defineConfig({
  input: `${process.env.API_URL}/openapi.json`,
  output: {
    path: './src/gen',
    clean: true
  },
  plugins: [
    pluginTs(),
    pluginAxios({
      baseURL: process.env.API_URL
    }),
    pluginVueQuery({
      client: 'axios'
    }),
    pluginZod()
  ]
})
