<script lang="ts" setup>
import { ChevronRightIcon } from '@lucide/vue'
import { reactiveOmit } from '@vueuse/core'
import { cn } from 'cn'
import type { CalendarNextProps } from 'reka-ui'
import { CalendarNext, useForwardProps } from 'reka-ui'
import { buttonVariants } from 'ui/components/ui/button'
import type { HTMLAttributes } from 'vue'

const props = defineProps<
  CalendarNextProps & { class?: HTMLAttributes['class'] }
>()

const delegatedProps = reactiveOmit(props, 'class')

const forwardedProps = useForwardProps(delegatedProps)
</script>

<template>
  <CalendarNext
    data-slot="calendar-next-button"
    :class="
      cn(
        buttonVariants({ variant: 'outline' }),
        'pointer-events-auto size-7 bg-transparent p-0 opacity-50 hover:opacity-100',
        props.class
      )
    "
    v-bind="forwardedProps"
  >
    <slot>
      <ChevronRightIcon class="cn-rtl-flip size-4" />
    </slot>
  </CalendarNext>
</template>
