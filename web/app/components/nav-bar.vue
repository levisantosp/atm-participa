<script setup lang="ts">
import { CircleUserIcon, HouseIcon } from '@lucide/vue'
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
  NavigationMenuList
} from 'ui'

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
      href: '/perfil'
    }
  ]
} as const
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

      <DropdownMenuContent>
        <DropdownMenuGroup>
          <DropdownMenuLabel>Conta</DropdownMenuLabel>

          <NuxtLink
            v-for="item in dropItems.account"
            :key="item.href"
            :to="item.href"
          >
            <DropdownMenuItem class="cursor-pointer">
              {{ item.label }}
            </DropdownMenuItem>
          </NuxtLink>
        </DropdownMenuGroup>
      </DropdownMenuContent>
    </DropdownMenu>
  </div>
</template>
