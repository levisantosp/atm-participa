<script setup lang="ts">
import { cn } from 'ui/lib/utils'
import type { HTMLAttributes } from 'vue'

const props = defineProps<{ class?: HTMLAttributes['class'] }>()
</script>

<template>
  <div
    data-slot="sheet-header"
    :class="cn('gap-1.5 p-6 flex flex-col', props.class)"
  >
    <slot />
  </div>
</template>
