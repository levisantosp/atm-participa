<script setup lang="ts">
import { CircleUserIcon, History, HouseIcon, LogOut, User2 } from '@lucide/vue'
import { getAuthMeQueryKey } from 'api-client'
import {
  Button,
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuGroup,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuTrigger,
  NavigationMenu,
  NavigationMenuItem,
  NavigationMenuLink,
  NavigationMenuList,
  Spinner
} from 'ui'
import { auth } from '~/lib/auth'
import { queryClient } from '~/lib/query-client'

const navItems = [
  {
    label: 'Página principal',
    icon: HouseIcon,
    href: '/'
  }
] as const

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
</script>

<template>
  <div class="flex w-full items-center justify-between">
    <NavigationMenu>
      <NavigationMenuList>
        <NavigationMenuItem v-for="item in navItems" :key="item.href">
          <NavigationMenuLink as-child>
            <NuxtLink :to="item.href">{{ item.label }}</NuxtLink>
          </NavigationMenuLink>
        </NavigationMenuItem>
      </NavigationMenuList>
    </NavigationMenu>

    <DropdownMenu>
      <DropdownMenuTrigger as-child>
        <Button variant="ghost">
          <CircleUserIcon :size="40" />
        </Button>
      </DropdownMenuTrigger>

      <DropdownMenuContent align="end">
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
      </DropdownMenuContent>
    </DropdownMenu>
  </div>
</template>
