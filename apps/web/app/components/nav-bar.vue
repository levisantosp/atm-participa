<script setup lang="ts">
import {
  CircleUserIcon,
  History,
  Home,
  LogOut,
  Moon,
  Sun,
  User2
} from '@lucide/vue'
import { getAuthMeQueryKey } from 'api-client'
import {
  Button,
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuGroup,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuTrigger,
  Spinner,
  Switch
} from 'ui'
import { auth } from '~/lib/auth'
import { queryClient } from '~/lib/query-client'
import { isDark } from '~/lib/theme'

const dropItems = {
  account: [
    {
      label: 'Perfil',
      href: '/perfil',
      icon: User2
    },
    {
      label: 'Histórico',
      href: '/historico',
      icon: History
    }
  ],
  nav: [
    {
      label: 'Página principal',
      href: '/',
      icon: Home
    }
  ]
} as const

const isPending = ref(false)

async function signOut() {
  isPending.value = true
  await auth.signOut()
  isPending.value = false
  await queryClient.resetQueries({
    queryKey: getAuthMeQueryKey()
  })
}

function setTheme(v: boolean) {
  isDark.value = v
}
</script>

<template>
  <div class="flex w-full items-center justify-end">
    <DropdownMenu>
      <DropdownMenuTrigger as-child>
        <Button variant="ghost">
          <CircleUserIcon :size="40" />
        </Button>
      </DropdownMenuTrigger>

      <DropdownMenuContent align="end">
        <DropdownMenuGroup>
          <DropdownMenuLabel>Navegação</DropdownMenuLabel>

          <NuxtLink
            v-for="item in dropItems.nav"
            :key="item.href"
            :to="item.href"
          >
            <DropdownMenuItem class="cursor-pointer">
              <component :is="item.icon" />
              {{ item.label }}
            </DropdownMenuItem>
          </NuxtLink>
        </DropdownMenuGroup>

        <DropdownMenuGroup>
          <DropdownMenuLabel>Conta</DropdownMenuLabel>

          <NuxtLink
            v-for="item in dropItems.account"
            :key="item.href"
            :to="item.href"
          >
            <DropdownMenuItem class="cursor-pointer">
              <component :is="item.icon" />
              {{ item.label }}
            </DropdownMenuItem>
          </NuxtLink>

          <DropdownMenuItem
            class="cursor-pointer hover:bg-destructive! text-destructive"
            @click="signOut"
            :disabled="isPending"
            @select.prevent
          >
            <LogOut />
            <span v-if="!isPending">Sair</span>
            <Spinner v-else />
          </DropdownMenuItem>
        </DropdownMenuGroup>

        <DropdownMenuGroup>
          <DropdownMenuLabel>Tema</DropdownMenuLabel>
          <DropdownMenuItem @select.prevent>
            <Sun />
            <Switch
              class="cursor-pointer"
              :model-value="isDark"
              @update:model-value="setTheme"
            />
            <Moon />
          </DropdownMenuItem>
        </DropdownMenuGroup>
      </DropdownMenuContent>
    </DropdownMenu>
  </div>
</template>
