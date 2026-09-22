import { readdirSync, writeFileSync } from 'node:fs'

let content = ''

for (const folder of readdirSync('./src/components/ui')) {
  if (folder.endsWith('.ts')) continue
  content += `export * from './components/ui/${folder}'\n`
}

writeFileSync('./src/index.ts', content, 'utf-8')
