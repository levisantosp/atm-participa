<script setup lang="ts">
import { ClipboardList, Home, Users } from '@lucide/vue'
import {
  Sidebar,
  SidebarContent,
  SidebarGroup,
  SidebarGroupContent,
  SidebarGroupLabel,
  SidebarHeader,
  SidebarInset,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarProvider,
  SidebarTrigger
} from 'ui'

const navItems = [
  {
    label: 'Página principal',
    href: '/',
    icon: Home
  },
  {
    label: 'Ocorrências',
    href: '/ocorrencias',
    icon: ClipboardList
  },
  {
    label: 'Usuários',
    href: '/usuarios',
    icon: Users
  }
] as const
</script>

<template>
  <div>
    <SidebarProvider>
      <Sidebar>
        <SidebarHeader>
          <RouterLink to="/">
            <h1 class="text-center font-bold md:text-xl">
              Painel Administrativo
            </h1>
          </RouterLink>
        </SidebarHeader>

        <SidebarContent>
          <SidebarGroup>
            <SidebarGroupLabel>Menu</SidebarGroupLabel>

            <SidebarGroupContent>
              <SidebarMenu>
                <SidebarMenuItem v-for="item in navItems" :key="item.href">
                  <SidebarMenuButton as-child>
                    <RouterLink :to="item.href">
                      <component :is="item.icon" />
                      <span>{{ item.label }}</span>
                    </RouterLink>
                  </SidebarMenuButton>
                </SidebarMenuItem>
              </SidebarMenu>
            </SidebarGroupContent>
          </SidebarGroup>
        </SidebarContent>
      </Sidebar>

      <SidebarInset>
        <header class="flex h-14 items-center gap-2 border-b px-4">
          <SidebarTrigger />
        </header>

        <main
          class="flex flex-col justify-center items-center mt-5 mb-30 ml-30 mr-30"
        >
          <slot />
        </main>
      </SidebarInset>
    </SidebarProvider>
  </div>
</template>
