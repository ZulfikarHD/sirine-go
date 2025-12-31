<template>
  <Motion
    :initial="{ opacity: 0, y: 15 }"
    :animate="{ opacity: 1, y: 0 }"
    :transition="{ duration: 0.25, delay: index * 0.05, ease: 'easeOut' }"
  >
    <div
      class="glass-card p-4 cursor-pointer active-scale transition-all duration-200"
      :class="{ 'border-l-4 border-red-500': item.is_overdue }"
      @click="$emit('click', item)"
    >
      <!-- Header: PO Number + OBC -->
      <div class="flex justify-between items-start mb-3">
        <div>
          <h3 class="text-lg font-bold text-gray-900">
            PO {{ item.po_number }}
          </h3>
          <span class="inline-block mt-1 px-2 py-1 bg-gradient-to-r from-indigo-100 to-fuchsia-100 text-indigo-700 text-xs font-semibold rounded-full">
            {{ item.obc_number }}
          </span>
        </div>
        
        <!-- Overdue badge -->
        <div v-if="item.is_overdue" class="flex items-center gap-1 px-2 py-1 bg-red-100 text-red-700 rounded-lg text-xs font-semibold">
          <svg class="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
            <path fill-rule="evenodd" d="M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.213 2.98-1.742 2.98H4.42c-1.53 0-2.493-1.646-1.743-2.98l5.58-9.92zM11 13a1 1 0 11-2 0 1 1 0 012 0zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z" clip-rule="evenodd"/>
          </svg>
          <span>Overdue</span>
        </div>
      </div>

      <!-- Stats Grid -->
      <div class="grid grid-cols-2 gap-3 mb-3">
        <div class="bg-gray-50 rounded-lg p-3">
          <p class="text-xs text-gray-500 mb-1">Target</p>
          <p class="text-base font-bold text-gray-900">{{ item.target_quantity }} <span class="text-xs font-normal text-gray-500">LB</span></p>
        </div>
        <div class="bg-gray-50 rounded-lg p-3">
          <p class="text-xs text-gray-500 mb-1">Mesin</p>
          <p class="text-base font-bold text-gray-900">{{ item.machine?.name || '-' }}</p>
        </div>
      </div>

      <!-- Operator Info -->
      <div v-if="item.operator" class="flex items-center gap-2 mb-3 p-2 bg-blue-50 rounded-lg">
        <div class="w-8 h-8 rounded-full bg-gradient-to-br from-indigo-500 to-fuchsia-500 flex items-center justify-center text-white font-semibold text-sm">
          {{ getInitials(item.operator.name) }}
        </div>
        <div class="flex-1 min-w-0">
          <p class="text-xs text-gray-500">Operator Cetak</p>
          <p class="text-sm font-semibold text-gray-900 truncate">{{ item.operator.name }}</p>
        </div>
      </div>

      <!-- Waiting Time Alert -->
      <div 
        v-if="item.is_overdue" 
        class="flex items-center gap-2 p-3 bg-red-50 border border-red-200 rounded-lg"
      >
        <svg class="w-5 h-5 text-red-600 flex-shrink-0" fill="currentColor" viewBox="0 0 20 20">
          <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd"/>
        </svg>
        <div class="flex-1">
          <p class="text-sm font-semibold text-red-800">Menunggu {{ formattedWaitingTime }}</p>
          <p class="text-xs text-red-600">Segera lakukan penghitungan</p>
        </div>
      </div>
      <div 
        v-else 
        class="text-sm text-gray-600 flex items-center gap-2"
      >
        <svg class="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
          <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm1-12a1 1 0 10-2 0v4a1 1 0 00.293.707l2.828 2.829a1 1 0 101.415-1.415L11 9.586V6z" clip-rule="evenodd"/>
        </svg>
        <span>Menunggu {{ formattedWaitingTime }}</span>
      </div>

      <!-- Call to Action -->
      <div class="mt-4 pt-3 border-t border-gray-200">
        <div class="flex items-center justify-between">
          <span class="text-xs text-gray-500">Tap untuk mulai</span>
          <svg class="w-5 h-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"/>
          </svg>
        </div>
      </div>
    </div>
  </Motion>
</template>

<script setup>
import { computed } from 'vue'
import { Motion } from 'motion-v'
import { useCountingApi } from '@/composables/useCountingApi'

const props = defineProps({
  item: {
    type: Object,
    required: true
  },
  index: {
    type: Number,
    default: 0
  }
})

defineEmits(['click'])

const { formatWaitingTime } = useCountingApi()

const formattedWaitingTime = computed(() => {
  return formatWaitingTime(props.item.waiting_minutes)
})

const getInitials = (name) => {
  return name
    .split(' ')
    .map(word => word[0])
    .join('')
    .toUpperCase()
    .slice(0, 2)
}
</script>
