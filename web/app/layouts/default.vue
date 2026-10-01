<script setup lang="ts">
import NavBar from '~/components/nav-bar.vue'

const title = 'Altamira Participa'
const description = 'Reporte problemas na cidade de Altamira'

useSeoMeta({
  title,
  ogTitle: title,
  description,
  ogDescription: description
})
</script>

<template>
  <main
    class="flex flex-col justify-center items-center mt-5 mb-30 md:ml-30 md:mr-30"
  >
    <nav-bar />
    <slot />
  </main>
</template>
