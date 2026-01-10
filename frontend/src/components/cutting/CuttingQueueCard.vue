<template>
  <div class="glass-card p-5 hover:shadow-lg transition-shadow">
    <!-- Header: PO Number & Priority -->
    <div class="flex items-start justify-between gap-4 mb-4">
      <div class="flex-1">
        <h3 class="text-lg font-bold text-gray-900">PO {{ poNumber }}</h3>
        <p class="text-sm text-gray-500">{{ obcNumber }}</p>
      </div>
      
      <!-- Priority Badge -->
      <PriorityBadge :priority="priority" />
    </div>

    <!-- Info Grid -->
    <div class="grid grid-cols-2 gap-3 mb-4">
      <!-- Input Lembar Besar -->
      <div class="bg-gradient-to-br from-indigo-50 to-fuchsia-50 rounded-lg p-3">
        <p class="text-xs text-gray-600 mb-1">Input (Lembar Besar)</p>
        <p class="text-xl font-bold text-indigo-600">{{ formatNumber(inputLembarBesar) }}</p>
      </div>

      <!-- Estimasi Output -->
      <div class="bg-gradient-to-br from-emerald-50 to-cyan-50 rounded-lg p-3">
        <p class="text-xs text-gray-600 mb-1">Estimasi Output</p>
        <p class="text-xl font-bold text-emerald-600">{{ formatNumber(estimatedOutput) }}</p>
      </div>
    </div>

    <!-- Timing Info -->
    <div class="flex items-center gap-2 text-sm mb-4">
      <svg class="w-4 h-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"/>
      </svg>
      <span :class="isOverdue ? 'text-red-600 font-semibold' : 'text-gray-600'">
        Menunggu {{ waitingMinutes }} menit
        <span v-if="isOverdue" class="text-red-600">⚠️</span>
      </span>
    </div>

    <!-- Action Button -->
    <button
      @click="$emit('start-cutting')"
      class="w-full py-3 px-4 bg-gradient-to-r from-indigo-600 to-fuchsia-600 text-white font-semibold rounded-xl hover:from-indigo-700 hover:to-fuchsia-700 transition-all active-scale flex items-center justify-center gap-2"
    >
      <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14.121 14.121L19 19m-7-7l7-7m-7 7l-2.879 2.879M12 12L9.121 9.121m0 5.758a3 3 0 10-4.243 4.243 3 3 0 004.243-4.243zm0-5.758a3 3 0 10-4.243-4.243 3 3 0 004.243 4.243z"/>
      </svg>
      <span>Mulai Pemotongan</span>
    </button>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { PriorityBadge } from '../common'

/**
 * Props untuk CuttingQueueCard component
 */
const props = defineProps({
  poId: {
    type: Number,
    required: true
  },
  poNumber: {
    type: [Number, String],
    required: true
  },
  obcNumber: {
    type: String,
    required: true
  },
  priority: {
    type: String,
    required: true
  },
  inputLembarBesar: {
    type: Number,
    required: true
  },
  estimatedOutput: {
    type: Number,
    required: true
  },
  countingCompletedAt: {
    type: String,
    required: true
  },
  waitingMinutes: {
    type: Number,
    required: true
  },
  isOverdue: {
    type: Boolean,
    default: false
  }
})

/**
 * Events yang di-emit oleh component
 */
defineEmits(['start-cutting'])

/**
 * Format number dengan thousand separator
 */
const formatNumber = (num) => {
  return new Intl.NumberFormat('id-ID').format(num)
}
</script>
