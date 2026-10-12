<script setup lang="ts">
import { toTypedSchema } from '@vee-validate/zod'
import { ResponseError } from 'api-client'
import {
  Button,
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
  Input,
  Label,
  Spinner
} from 'ui'
import { useForm } from 'vee-validate'
import { useRouter } from 'vue-router'
import { toast } from 'vue-sonner'
import { z } from 'zod'
import { auth } from '@/lib/auth'

definePage({
  meta: {
    layout: 'public'
  }
})

const schema = z.object({
  email: z.email('Informe um e-mail válido.'),
  password: z.string('Informe a senha').trim()
})

const { defineField, errors, handleSubmit, isSubmitting } = useForm({
  validationSchema: toTypedSchema(schema)
})

const [email, emailAttrs] = defineField('email')
const [password, passwordAttrs] = defineField('password')

const router = useRouter()

const signIn = handleSubmit(async (data) => {
  try {
    await auth.signIn.email(data.email, data.password)
    router.push('/')
  } catch (err) {
    console.error(err)

    if (err instanceof ResponseError) {
      toast.error(err.response.data.detail)
    }
  }
})
</script>

<template>
  <div
    class="flex min-h-[calc(100vh-8.75rem)] w-full items-center justify-center"
  >
    <Card class="w-full max-w-md md:max-w-xl">
      <CardHeader>
        <CardTitle>Bem vindo(a) de volta!</CardTitle>
        <CardDescription>
          Informe seu e-mail e senha abaixo para entrar na sua conta
        </CardDescription>
      </CardHeader>

      <CardContent>
        <form id="sign-in" @submit="signIn">
          <div class="flex flex-col gap-6">
            <div class="grid gap-2">
              <Label for="email">E-mail</Label>
              <Input
                id="email"
                type="email"
                placeholder="usuario@exemplo.com"
                v-model="email"
                v-bind="emailAttrs"
              />

              <span v-if="errors.email" class="text-sm text-red-400">
                {{ errors.email }}
              </span>
            </div>
            <div class="grid gap-2">
              <div class="flex items-center">
                <Label for="password">Senha</Label>
                <a
                  href="#"
                  class="ml-auto inline-block text-sm underline-offset-4 hover:underline"
                >
                  Esqueci minha senha
                </a>
              </div>
              <Input
                id="password"
                type="password"
                v-model="password"
                v-bind="passwordAttrs"
              />

              <span v-if="errors.password" class="text-sm text-red-400">
                {{ errors.password }}
              </span>
            </div>
          </div>
        </form>
      </CardContent>

      <CardFooter class="flex-col gap-2">
        <Button
          type="submit"
          class="w-full"
          :disabled="isSubmitting"
          form="sign-in"
        >
          <Spinner v-if="isSubmitting" />
          <span v-else>Entrar</span>
        </Button>
        <Button variant="outline" class="w-full" disabled>
          Entrar com Google
        </Button>
      </CardFooter>
    </Card>
  </div>
</template>
