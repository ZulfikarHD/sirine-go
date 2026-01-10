<template>
  <div class="glass-card p-5 hover:shadow-lg transition-shadow">
    <!-- Header: Label Number & Status -->
    <div class="flex items-start justify-between gap-4 mb-4">
      <div class="flex-1">
        <h3 class="text-lg font-bold text-gray-900">Label {{ labelNumber }} / {{ totalLabels }}</h3>
        <p class="text-sm text-gray-500">{{ sisiran }} - {{ formatNumber(targetQuantity) }} lembar</p>
      </div>
      
      <!-- Status Badge -->
      <span
        class="px-3 py-1 rounded-full text-xs font-semibold"
        :class="statusClasses"
      >
        {{ statusText }}
      </span>
    </div>

    <!-- PO Info -->
    <div class="bg-gradient-to-br from-gray-50 to-gray-100 rounded-lg p-3 mb-4">
      <div class="flex items-center gap-2 mb-1">
        <svg class="w-4 h-4 text-gray-400" fill="currentColor" viewBox="0 0 20 20">
          <path d="M9 2a1 1 0 000 2h2a1 1 0 100-2H9z"/>
          <path fill-rule="evenodd" d="M4 5a2 2 0 012-2 3 3 0 003 3h2a3 3 0 003-3 2 2 0 012 2v11a2 2 0 01-2 2H6a2 2 0 01-2-2V5zm3 4a1 1 0 000 2h.01a1 1 0 100-2H7zm3 0a1 1 0 000 2h3a1 1 0 100-2h-3zm-3 4a1 1 0 100 2h.01a1 1 0 100-2H7zm3 0a1 1 0 100 2h3a1 1 0 100-2h-3z" clip-rule="evenodd"/>
        </svg>
        <span class="text-sm font-semibold text-gray-700">PO {{ poNumber }}</span>
      </div>
      <p class="text-xs text-gray-500">{{ obcNumber }}</p>
    </div>

    <!-- Cutting Info -->
    <div v-if="cuttingInfo" class="flex items-center gap-2 text-sm mb-4">
      <svg class="w-4 h-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"/>
      </svg>
      <span class="text-gray-600">
        Dipotong oleh <span class="font-semibold">{{ cuttingInfo.operator }}</span>
      </span>
    </div>

    <!-- Waste Info (if available) -->
    <div v-if="wastePercentage !== null" class="mb-4">
      <div
        class="rounded-lg p-3 flex items-center gap-2"
        :class="wastePercentage > 2 ? 'bg-red-50' : 'bg-emerald-50'"
      >
        <div
          class="w-8 h-8 rounded-full flex items-center justify-center flex-shrink-0"
          :class="wastePercentage > 2 ? 'bg-red-500' : 'bg-emerald-500'"
        >
          <svg v-if="wastePercentage > 2" class="w-4 h-4 text-white" fill="currentColor" viewBox="0 0 20 20">
            <path fill-rule="evenodd" d="M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.213 2.98-1.742 2.98H4.42c-1.53 0-2.493-1.646-1.743-2.98l5.58-9.92zM11 13a1 1 0 11-2 0 1 1 0 012 0zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z" clip-rule="evenodd"/>
          </svg>
          <svg v-else class="w-4 h-4 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/>
          </svg>
        </div>
        <div class="flex-1">
          <p class="text-xs font-semibold" :class="wastePercentage > 2 ? 'text-red-900' : 'text-emerald-900'">
            Waste: {{ wastePercentage }}%
          </p>
        </div>
      </div>
    </div>

    <!-- Action Button -->
    <button
      v-if="qcStatus === 'PENDING'"
      @click="$emit('start-verification')"
      class="w-full py-3 px-4 bg-gradient-to-r from-indigo-600 to-fuchsia-600 text-white font-semibold rounded-xl hover:from-indigo-700 hover:to-fuchsia-700 transition-all active-scale flex items-center justify-center gap-2"
    >
      <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/>
      </svg>
      <span>Mulai Verifikasi</span>
    </button>
    <button
      v-else-if="qcStatus === 'IN_PROGRESS'"
      @click="$emit('continue-verification')"
      class="w-full py-3 px-4 bg-gradient-to-r from-emerald-600 to-cyan-600 text-white font-semibold rounded-xl hover:from-emerald-700 hover:to-cyan-700 transition-all active-scale flex items-center justify-center gap-2"
    >
      <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"/>
      </svg>
      <span>Lanjutkan Verifikasi</span>
    </button>
    <div v-else-if="qcStatus === 'COMPLETED'" class="text-center py-3 px-4 bg-emerald-50 border-2 border-emerald-200 rounded-xl">
      <div class="flex items-center justify-center gap-2">
        <svg class="w-5 h-5 text-emerald-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/>
        </svg>
        <span class="text-sm font-semibold text-emerald-900">Verifikasi Selesai</span>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'

/**
 * Props untuk VerificationLabelCard component
 */
const props = defineProps({
  labelId: {
    type: Number,
    required: true
  },
  labelNumber: {
    type: Number,
    required: true
  },
  totalLabels: {
    type: Number,
    required: true
  },
  targetQuantity: {
    type: Number,
    required: true
  },
  sisiran: {
    type: String,
    required: true,
    validator: (value) => ['KIRI', 'KANAN'].includes(value)
  },
  qcStatus: {
    type: String,
    required: true,
    validator: (value) => ['PENDING', 'IN_PROGRESS', 'COMPLETED'].includes(value)
  },
  poNumber: {
    type: [Number, String],
    required: true
  },
  obcNumber: {
    type: String,
    required: true
  },
  cuttingInfo: {
    type: Object,
    default: null
  },
  wastePercentage: {
    type: Number,
    default: null
  }
})

/**
 * Events yang di-emit oleh component
 */
defineEmits(['start-verification', 'continue-verification'])

/**
 * Computed status classes untuk badge
 */
const statusClasses = computed(() => {
  switch (props.qcStatus) {
    case 'PENDING':
      return 'bg-gray-200 text-gray-700'
    case 'IN_PROGRESS':
      return 'bg-amber-200 text-amber-800'
    case 'COMPLETED':
      return 'bg-emerald-200 text-emerald-800'
    default:
      return 'bg-gray-200 text-gray-700'
  }
})

/**
 * Computed status text
 */
const statusText = computed(() => {
  switch (props.qcStatus) {
    case 'PENDING':
      return 'Menunggu'
    case 'IN_PROGRESS':
      return 'Sedang Verifikasi'
    case 'COMPLETED':
      return 'Selesai'
    default:
      return 'Unknown'
  }
})

/**
 * Format number dengan thousand separator
 */
const formatNumber = (num) => {
  return new Intl.NumberFormat('id-ID').format(num)
}
</script>
