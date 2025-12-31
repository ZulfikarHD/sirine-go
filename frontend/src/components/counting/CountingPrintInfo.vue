<template>
  <Motion v-bind="entranceAnimations.fadeScale">
    <div class="glass-card p-4">
      <div class="flex items-center justify-between mb-3">
        <h3 class="text-sm font-bold text-gray-700">Info Hasil Cetak</h3>
        <button 
          v-if="collapsible"
          @click="isExpanded = !isExpanded"
          class="p-1 hover:bg-gray-100 rounded-lg transition-colors active-scale"
        >
          <svg 
            class="w-5 h-5 text-gray-500 transition-transform duration-200"
            :class="{ 'rotate-180': !isExpanded }"
            fill="none" 
            stroke="currentColor" 
            viewBox="0 0 24 24"
          >
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"/>
          </svg>
        </button>
      </div>

      <div v-show="isExpanded" class="space-y-3">
        <!-- Target Quantity -->
        <div class="flex items-center justify-between p-3 bg-gradient-to-r from-indigo-50 to-fuchsia-50 rounded-lg">
          <div class="flex items-center gap-3">
            <div class="w-10 h-10 rounded-full bg-gradient-to-br from-indigo-500 to-fuchsia-500 flex items-center justify-center">
              <svg class="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"/>
              </svg>
            </div>
            <div>
              <p class="text-xs text-gray-600">Target Quantity</p>
              <p class="text-xl font-bold bg-gradient-to-r from-indigo-600 to-fuchsia-600 bg-clip-text text-transparent">
                {{ poInfo.target_quantity }} LB
              </p>
            </div>
          </div>
        </div>

        <!-- Print Details Grid -->
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
          <!-- Machine -->
          <div class="bg-gray-50 rounded-lg p-3">
            <div class="flex items-center gap-2 mb-1">
              <svg class="w-4 h-4 text-gray-500" fill="currentColor" viewBox="0 0 20 20">
                <path fill-rule="evenodd" d="M6 2a2 2 0 00-2 2v12a2 2 0 002 2h8a2 2 0 002-2V4a2 2 0 00-2-2H6zm1 2a1 1 0 000 2h6a1 1 0 100-2H7zm6 7a1 1 0 011 1v3a1 1 0 11-2 0v-3a1 1 0 011-1zm-3 3a1 1 0 100 2h.01a1 1 0 100-2H10zm-4 1a1 1 0 011-1h.01a1 1 0 110 2H7a1 1 0 01-1-1zm1-4a1 1 0 100 2h.01a1 1 0 100-2H7zm2 1a1 1 0 011-1h.01a1 1 0 110 2H10a1 1 0 01-1-1zm4-4a1 1 0 100 2h.01a1 1 0 100-2H13zM9 9a1 1 0 011-1h.01a1 1 0 110 2H10a1 1 0 01-1-1zM7 8a1 1 0 000 2h.01a1 1 0 000-2H7z" clip-rule="evenodd"/>
              </svg>
              <p class="text-xs text-gray-500">Mesin</p>
            </div>
            <p class="text-sm font-bold text-gray-900">{{ printInfo?.machine_name || '-' }}</p>
          </div>

          <!-- Operator -->
          <div class="bg-gray-50 rounded-lg p-3">
            <div class="flex items-center gap-2 mb-1">
              <svg class="w-4 h-4 text-gray-500" fill="currentColor" viewBox="0 0 20 20">
                <path fill-rule="evenodd" d="M10 9a3 3 0 100-6 3 3 0 000 6zm-7 9a7 7 0 1114 0H3z" clip-rule="evenodd"/>
              </svg>
              <p class="text-xs text-gray-500">Operator</p>
            </div>
            <p class="text-sm font-bold text-gray-900">{{ printInfo?.operator_name || '-' }}</p>
          </div>
        </div>

        <!-- Print Completed Time -->
        <div v-if="printInfo?.finalized_at" class="flex items-center gap-2 p-3 bg-green-50 border border-green-200 rounded-lg">
          <svg class="w-5 h-5 text-green-600" fill="currentColor" viewBox="0 0 20 20">
            <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd"/>
          </svg>
          <div>
            <p class="text-xs text-green-600">Selesai Cetak</p>
            <p class="text-sm font-semibold text-green-800">{{ formatDateTime(printInfo.finalized_at) }}</p>
          </div>
        </div>
      </div>
    </div>
  </Motion>
</template>

<script setup>
import { ref } from 'vue'
import { Motion } from 'motion-v'
import { entranceAnimations } from '@/composables/useMotion'

const props = defineProps({
  poInfo: {
    type: Object,
    required: true
  },
  printInfo: {
    type: Object,
    default: null
  },
  collapsible: {
    type: Boolean,
    default: false
  }
})

const isExpanded = ref(true)

const formatDateTime = (dateString) => {
  const date = new Date(dateString)
  return date.toLocaleString('id-ID', {
    day: 'numeric',
    month: 'short',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
    timeZone: 'Asia/Jakarta'
  })
}
</script>
