<script setup lang="ts">
import { reactiveOmit } from '@vueuse/core'
import type { SelectSeparatorProps } from 'reka-ui'
import { SelectSeparator } from 'reka-ui'
import { cn } from 'ui/lib/utils'
import type { HTMLAttributes } from 'vue'

const props = defineProps<
  SelectSeparatorProps & { class?: HTMLAttributes['class'] }
>()

const delegatedProps = reactiveOmit(props, 'class')
</script>

<template>
  <SelectSeparator
    data-slot="select-separator"
    v-bind="delegatedProps"
    :class="
      cn('bg-border -mx-1.5 my-1.5 h-px pointer-events-none', props.class)
    "
  />
</template>
