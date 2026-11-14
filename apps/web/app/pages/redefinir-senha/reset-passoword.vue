<script setup lang="ts">
import { toTypedSchema } from '@vee-validate/zod'
import { usePostAuthForgotPassword, usePostAuthResetPassword } from 'api-client'
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
import { toast } from 'vue-sonner'
import { z } from 'zod'

const schema = z
  .object({
    password: z
      .string('Campo obrigatório')
      .min(8, 'A senha deve ter no mínimo 8 caracteres')
      .max(72, 'A senha deve ter no máximo 72 caracteres'),
    confirmPassword: z.string('Campo obrigatório').min(1, 'Confirme a senha')
  })
  .refine((data) => data.password === data.confirmPassword, {
    error: 'As senhas não coincidem',
    path: ['confirmPassword']
  })

const { handleSubmit, errors, defineField, resetForm } = useForm({
  validationSchema: toTypedSchema(schema)
})

const [password, passwordAttrs] = defineField('password')
const [confirmPassword, confirmPasswordAttrs] = defineField('confirmPassword')

const route = useRoute()

const { isPending, mutate } = usePostAuthResetPassword({
  mutation: {
    onSuccess() {
      toast.success(
        'Senha redefinida com sucesso! Acesse a página de login e entre com a sua conta.',
        {
          duration: 10_000
        }
      )
      resetForm()
    }
  }
})

const onSubmit = handleSubmit((data) =>
  mutate({
    body: {
      token: route.query.token!.toString(),
      password: data.password
    }
  })
)
</script>

<template>
  <Card class="w-full max-w-md md:max-w-xl">
    <CardHeader>
      <CardTitle>Redefinir Senha</CardTitle>
      <CardDescription> Informe sua nova senha abaixo </CardDescription>
    </CardHeader>

    <CardContent>
      <form id="reset-passoword" @submit="onSubmit">
        <div class="flex flex-col gap-6">
          <div class="grid gap-2">
            <Label for="password">Senha</Label>
            <Input
              id="password"
              type="password"
              v-model="password"
              v-bind="passwordAttrs"
            />

            <span v-if="errors.password" class="text-sm text-destructive">
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

            <span
              v-if="errors.confirmPassword"
              class="text-sm text-destructive"
            >
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
        :disabled="isPending"
        form="reset-passoword"
      >
        <Spinner v-if="isPending" />
        <span v-else>Redefinir senha</span>
      </Button>
    </CardFooter>
  </Card>
</template>
