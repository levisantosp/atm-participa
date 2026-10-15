<script setup lang="ts">
import { Home } from '@lucide/vue'
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
    label: 'Home',
    href: '/',
    icon: Home
  }
] as const
</script>

<template>
  <div>
    <SidebarProvider>
      <Sidebar>
        <SidebarHeader>
          <RouterLink to="/">
            <h1 class="text-center font-bold md:text-xl">Sabine Dashboard</h1>
          </RouterLink>
        </SidebarHeader>

        <SidebarContent>
          <SidebarGroup>
            <SidebarGroupLabel>Menu</SidebarGroupLabel>

            <SidebarGroupContent>
              <SidebarMenu>
                <SidebarMenuItem>
                  <SidebarMenuButton as-child v-for="item in navItems">
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

        <main class="flex flex-col mt-5 mb-30 ml-30 mr-30">
          <slot />
        </main>
      </SidebarInset>
    </SidebarProvider>
  </div>
</template>
