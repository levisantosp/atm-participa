<script setup lang="ts">
import ResetPassoword from './reset-passoword.vue'
import SendEmail from './send-email.vue'

definePageMeta({
  layout: 'default'
})

const route = useRoute()
</script>

<template>
  <div
    class="flex min-h-[calc(100vh-8.75rem)] w-full items-center justify-center"
  >
    <SendEmail v-if="!route.query.token" />
    <ResetPassoword v-else />
  </div>
</template>
