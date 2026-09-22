import { readdirSync, writeFileSync } from 'node:fs'

let content = `import { client } from './gen/.kubb/client'

client.interceptors.request.use(async (req) => {
  req.withCredentials = true

  if (process.env.NUXT_PUBLIC_ENABLE_API_DELAY || process.env.VITE_ENABLE_API_DELAY) {
    await new Promise((r) => setTimeout(r, Math.round(Math.random() * 4000)))
  }

  return req
})

`

for (const folder of readdirSync('./src/gen')) {
  for (const file of readdirSync(`./src/gen/${folder}`)) {
    if (file === 'index.ts') continue
    content += `export * from './gen/${folder}/${file}'\n`
  }
}

writeFileSync('./src/index.ts', content, 'utf-8')
