<script setup lang="ts">
import InReviewIssues from './in-review-issues.vue'
import OpenIssues from './open-issues.vue'
</script>

<template>
  <h1 class="font-bold text-xl md:text-2xl">Visão geral</h1>

  <div class="mt-10 flex flex-wrap justify-center gap-5">
    <OpenIssues />
    <InReviewIssues />
  </div>
</template>
