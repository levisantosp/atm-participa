<script setup lang="ts">
import { cn } from 'ui/lib/utils'
import type { HTMLAttributes } from 'vue'

const props = defineProps<{
  class?: HTMLAttributes['class']
}>()
</script>

<template>
  <div
    data-slot="card-title"
    :class="cn('text-base font-medium cn-font-heading', props.class)"
  >
    <slot />
  </div>
</template>
