<script setup lang="ts">
import NavBar from '~/components/nav-bar.vue'

useSeoMeta({
  title: 'Altamira Participa',
  description: 'Reporte problemas na cidade de Altamira'
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
