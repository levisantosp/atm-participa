'use client'

import { CircleUser, House } from 'lucide-react'
import Link from 'next/link'
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
    icon: House,
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

export function NavBar() {
  return (
    <div className='flex w-full justify-between items-center'>
      <NavigationMenu>
        <NavigationMenuList>
          {navItems.map((item) => (
            <NavigationMenuItem key={item.href}>
              <NavigationMenuLink render={<Link href='/' />}>
                {item.label}
              </NavigationMenuLink>
            </NavigationMenuItem>
          ))}
        </NavigationMenuList>
      </NavigationMenu>

      <DropdownMenu>
        <DropdownMenuTrigger
          render={
            <Button variant='ghost'>
              <CircleUser size={40} />
            </Button>
          }
        />

        <DropdownMenuContent>
          <DropdownMenuGroup>
            <DropdownMenuLabel>Conta</DropdownMenuLabel>

            {dropItems.account.map((item) => (
              <Link key={item.href} href={item.href}>
                <DropdownMenuItem className='cursor-pointer'>
                  {item.label}
                </DropdownMenuItem>
              </Link>
            ))}
          </DropdownMenuGroup>
        </DropdownMenuContent>
      </DropdownMenu>
    </div>
  )
}
