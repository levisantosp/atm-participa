<script setup lang="ts">
import { toTypedSchema } from '@vee-validate/zod'
import { usePostAuthForgotPassword } from 'api-client'
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

const schema = z.object({
  email: z.email('Informe um e-mail válido')
})

const { handleSubmit, errors, defineField } = useForm({
  validationSchema: toTypedSchema(schema)
})

const [email, emailAttrs] = defineField('email')

const route = useRoute()
const router = useRouter()

const { isPending, mutate } = usePostAuthForgotPassword({
  mutation: {
    onSuccess() {
      toast.success(
        'Enviamos um e-mail com as instruções para redefinir sua senha. Verifique sua caixa de entrada.',
        {
          duration: 10_000
        }
      )
      void router.push('/sign-in')
    }
  }
})

const onSubmit = handleSubmit((data) =>
  mutate({
    body: data
  })
)
</script>

<template>
  <Card v-if="!route.query.token" class="w-full max-w-md md:max-w-xl">
    <CardHeader>
      <CardTitle>Redefinir Senha</CardTitle>
      <CardDescription>
        Informe seu e-mail abaixo para redefinir sua senha
      </CardDescription>
    </CardHeader>

    <CardContent>
      <form id="send-email" @submit="onSubmit">
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
        </div>
      </form>
    </CardContent>

    <CardFooter class="flex-col gap-2">
      <Button class="w-full" :disabled="isPending" form="send-email">
        <Spinner v-if="isPending" />
        <span v-else>Enviar e-mail</span>
      </Button>
    </CardFooter>
  </Card>
</template>
