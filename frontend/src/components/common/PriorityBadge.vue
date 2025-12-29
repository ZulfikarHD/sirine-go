<template>
  <div 
    :class="badgeClasses"
    class="inline-flex items-center gap-1.5 font-semibold rounded-lg transition-all duration-150"
  >
    <div :class="dotClasses" class="rounded-full" />
    <span>{{ priorityLabel }}</span>
  </div>
</template>

<script setup>
/**
 * PriorityBadge - Reusable badge component untuk priority display
 * dengan gradient styling sesuai design standard
 */
import { computed } from 'vue'

const props = defineProps({
  priority: {
    type: String,
    required: true,
    validator: (value) => ['URGENT', 'NORMAL', 'LOW'].includes(value)
  },
  size: {
    type: String,
    default: 'md',
    validator: (value) => ['sm', 'md', 'lg'].includes(value)
  }
})

const badgeClasses = computed(() => {
  const sizeClasses = {
    sm: 'px-2 py-0.5 text-xs',
    md: 'px-3 py-1 text-sm',
    lg: 'px-4 py-1.5 text-base'
  }

  const priorityClasses = {
    URGENT: 'bg-gradient-to-r from-red-50 to-orange-50 text-red-700 border border-red-200/50',
    NORMAL: 'bg-gradient-to-r from-indigo-50 to-fuchsia-50 text-indigo-700 border border-indigo-200/50',
    LOW: 'bg-gray-50 text-gray-600 border border-gray-200/50'
  }

  return [
    sizeClasses[props.size],
    priorityClasses[props.priority]
  ]
})

const dotClasses = computed(() => {
  const sizeClasses = {
    sm: 'w-1.5 h-1.5',
    md: 'w-2 h-2',
    lg: 'w-2.5 h-2.5'
  }

  const priorityClasses = {
    URGENT: 'bg-red-500',
    NORMAL: 'bg-indigo-500',
    LOW: 'bg-gray-400'
  }

  return [
    sizeClasses[props.size],
    priorityClasses[props.priority]
  ]
})

const priorityLabel = computed(() => {
  const labels = {
    URGENT: 'Urgent',
    NORMAL: 'Normal',
    LOW: 'Rendah'
  }
  return labels[props.priority] || props.priority
})
</script>
