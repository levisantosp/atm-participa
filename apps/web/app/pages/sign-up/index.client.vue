<script setup lang="ts">
import { toTypedSchema } from '@vee-validate/zod'
import { ResponseError } from 'api-client'
import {
  Button,
  Card,
  CardAction,
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
import { toast } from 'vue-sonner'
import { z } from 'zod'
import { auth } from '~/lib/auth'

definePageMeta({
  layout: 'default'
})

const schema = z
  .object({
    displayName: z
      .string()
      .min(1, 'Informe o nome completo')
      .max(100, 'O nome deve ter no máximo 100 caracteres'),
    username: z
      .string()
      .min(3, 'O nome de usuário deve ter no mínimo 3 caracteres')
      .max(32, 'O nome de usuário deve ter no máximo 32 caracteres')
      .regex(
        /^[a-zA-Z0-9_]+$/,
        'O nome de usuário só pode conter letras, números e underscore'
      ),
    email: z.email('Informe um e-mail válido'),
    password: z
      .string()
      .min(8, 'A senha deve ter no mínimo 8 caracteres')
      .max(72, 'A senha deve ter no máximo 72 caracteres'),
    confirmPassword: z.string().min(1, 'Confirme a senha')
  })
  .refine((data) => data.password === data.confirmPassword, {
    error: 'As senhas não coincidem',
    path: ['confirmPassword']
  })

const { handleSubmit, errors, defineField, isSubmitting } = useForm({
  validationSchema: toTypedSchema(schema)
})

const [displayName, displayNameAttrs] = defineField('displayName')
const [username, usernameAttrs] = defineField('username')
const [email, emailAttrs] = defineField('email')
const [password, passwordAttrs] = defineField('password')
const [confirmPassword, confirmPasswordAttrs] = defineField('confirmPassword')

const router = useRouter()

const signUp = handleSubmit(async (data) => {
  try {
    await auth.signUp.email({
      email: data.email,
      password: data.password,
      username: data.username,
      displayName: data.displayName
    })

    router.push('/')
  } catch (err) {
    console.error(err)
    if (err instanceof ResponseError) {
      toast.error(err.data.detail)
    } else if (err instanceof Error) {
      toast.error('Ocorreu um erro inesperado', {
        description: err.message
      })
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
        <CardTitle>Crie sua conta</CardTitle>
        <CardDescription>
          Preencha os dados abaixo para se cadastrar
        </CardDescription>

        <CardAction>
          <NuxtLink to="/sign-in">
            <Button variant="link">Entrar</Button>
          </NuxtLink>
        </CardAction>
      </CardHeader>

      <CardContent>
        <form id="sign-up" @submit="signUp">
          <div class="flex flex-col gap-6">
            <div class="grid gap-2">
              <Label for="displayName">Nome completo</Label>
              <Input
                id="displayName"
                type="text"
                placeholder="Seu nome"
                v-model="displayName"
                v-bind="displayNameAttrs"
              />

              <span v-if="errors.displayName" class="text-sm text-red-400">
                {{ errors.displayName }}
              </span>
            </div>
            <div class="grid gap-2">
              <Label for="username">Nome de usuário</Label>
              <Input
                id="username"
                type="text"
                v-model="username"
                v-bind="usernameAttrs"
              />

              <span v-if="errors.username" class="text-sm text-red-400">
                {{ errors.username }}
              </span>
            </div>
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
              <Label for="password">Senha</Label>
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
            <div class="grid gap-2">
              <Label for="confirmPassword">Confirmar senha</Label>
              <Input
                id="confirmPassword"
                type="password"
                v-model="confirmPassword"
                v-bind="confirmPasswordAttrs"
              />

              <span v-if="errors.confirmPassword" class="text-sm text-red-400">
                {{ errors.confirmPassword }}
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
          form="sign-up"
        >
          <Spinner v-if="isSubmitting" />
          <span v-else>Criar conta</span>
        </Button>
      </CardFooter>
    </Card>
  </div>
</template>
