<template>
  <Motion v-bind="entranceAnimations.fadeScale">
    <div class="glass-card p-4">
      <h3 class="text-sm font-bold text-gray-700 mb-3">Perhitungan Real-time</h3>

      <div class="space-y-3">
        <!-- Total Counted -->
        <div class="flex items-center justify-between p-4 bg-gradient-to-r from-gray-50 to-gray-100 rounded-xl">
          <div>
            <p class="text-xs text-gray-600 mb-1">Total Dihitung</p>
            <p class="text-2xl font-bold text-gray-900">{{ calculations.totalCounted }} <span class="text-sm font-normal text-gray-500">LB</span></p>
          </div>
          <div class="w-12 h-12 rounded-full bg-gray-200 flex items-center justify-center">
            <svg class="w-6 h-6 text-gray-600" fill="currentColor" viewBox="0 0 20 20">
              <path d="M2 11a1 1 0 011-1h2a1 1 0 011 1v5a1 1 0 01-1 1H3a1 1 0 01-1-1v-5zM8 7a1 1 0 011-1h2a1 1 0 011 1v9a1 1 0 01-1 1H9a1 1 0 01-1-1V7zM14 4a1 1 0 011-1h2a1 1 0 011 1v12a1 1 0 01-1 1h-2a1 1 0 01-1-1V4z"/>
            </svg>
          </div>
        </div>

        <!-- Variance -->
        <div 
          class="flex items-center justify-between p-4 rounded-xl"
          :class="varianceClass"
        >
          <div>
            <p class="text-xs font-semibold mb-1" :class="varianceTextClass">Selisih dari Target</p>
            <p class="text-xl font-bold" :class="varianceTextClass">
              {{ calculations.varianceFromTarget > 0 ? '+' : '' }}{{ calculations.varianceFromTarget }} LB
              <span class="text-sm font-normal">({{ calculations.variancePercentage > 0 ? '+' : '' }}{{ calculations.variancePercentage }}%)</span>
            </p>
          </div>
          <div class="w-10 h-10 rounded-full flex items-center justify-center" :class="varianceIconBg">
            <svg class="w-5 h-5" :class="varianceIconColor" fill="currentColor" viewBox="0 0 20 20">
              <path v-if="calculations.varianceFromTarget === 0" fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd"/>
              <path v-else-if="calculations.varianceFromTarget > 0" fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm1-11a1 1 0 10-2 0v3.586L7.707 9.293a1 1 0 00-1.414 1.414l3 3a1 1 0 001.414 0l3-3a1 1 0 00-1.414-1.414L11 10.586V7z" clip-rule="evenodd"/>
              <path v-else fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm1-11a1 1 0 10-2 0v3.586l-1.293-1.293a1 1 0 00-1.414 1.414l3 3a1 1 0 001.414 0l3-3a1 1 0 00-1.414-1.414L11 10.586V7z" clip-rule="evenodd" transform="rotate(180 10 10)"/>
            </svg>
          </div>
        </div>

        <!-- Tolerance Warning -->
        <div v-if="!calculations.isWithinTolerance && calculations.hasVariance" class="flex items-start gap-2 p-3 bg-yellow-50 border border-yellow-200 rounded-lg">
          <svg class="w-5 h-5 text-yellow-600 flex-shrink-0 mt-0.5" fill="currentColor" viewBox="0 0 20 20">
            <path fill-rule="evenodd" d="M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.213 2.98-1.742 2.98H4.42c-1.53 0-2.493-1.646-1.743-2.98l5.58-9.92zM11 13a1 1 0 11-2 0 1 1 0 012 0zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z" clip-rule="evenodd"/>
          </svg>
          <div class="text-sm">
            <p class="font-semibold text-yellow-800">Selisih Melebihi Toleransi</p>
            <p class="text-yellow-700">Pastikan keterangan selisih diisi dengan jelas</p>
          </div>
        </div>

        <!-- Quality Breakdown -->
        <div class="grid grid-cols-2 gap-3">
          <!-- Good Percentage -->
          <div class="bg-green-50 border border-green-200 rounded-lg p-3">
            <div class="flex items-center gap-2 mb-2">
              <div class="w-8 h-8 rounded-full bg-green-500 flex items-center justify-center">
                <svg class="w-4 h-4 text-white" fill="currentColor" viewBox="0 0 20 20">
                  <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd"/>
                </svg>
              </div>
              <p class="text-xs font-semibold text-green-700">Baik</p>
            </div>
            <p class="text-2xl font-bold text-green-800">{{ calculations.percentageGood }}%</p>
            <p class="text-xs text-green-600 mt-1">{{ quantityGood }} lembar</p>
          </div>

          <!-- Defect Percentage -->
          <div class="bg-red-50 border border-red-200 rounded-lg p-3">
            <div class="flex items-center gap-2 mb-2">
              <div class="w-8 h-8 rounded-full bg-red-500 flex items-center justify-center">
                <svg class="w-4 h-4 text-white" fill="currentColor" viewBox="0 0 20 20">
                  <path fill-rule="evenodd" d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z" clip-rule="evenodd"/>
                </svg>
              </div>
              <p class="text-xs font-semibold text-red-700">Rusak</p>
            </div>
            <p class="text-2xl font-bold text-red-800">{{ calculations.percentageDefect }}%</p>
            <p class="text-xs text-red-600 mt-1">{{ quantityDefect }} lembar</p>
          </div>
        </div>

        <!-- Defect threshold warning -->
        <div v-if="calculations.isDefectAboveThreshold" class="flex items-start gap-2 p-3 bg-red-50 border border-red-200 rounded-lg">
          <svg class="w-5 h-5 text-red-600 flex-shrink-0 mt-0.5" fill="currentColor" viewBox="0 0 20 20">
            <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd"/>
          </svg>
          <div class="text-sm">
            <p class="font-semibold text-red-800">Kerusakan Tinggi</p>
            <p class="text-red-700">Breakdown jenis kerusakan wajib diisi (rusak >5%)</p>
          </div>
        </div>
      </div>
    </div>
  </Motion>
</template>

<script setup>
import { computed } from 'vue'
import { Motion } from 'motion-v'
import { entranceAnimations } from '@/composables/useMotion'

const props = defineProps({
  calculations: {
    type: Object,
    required: true
  },
  quantityGood: {
    type: Number,
    default: 0
  },
  quantityDefect: {
    type: Number,
    default: 0
  }
})

const varianceClass = computed(() => {
  if (props.calculations.varianceFromTarget === 0) {
    return 'bg-green-50 border border-green-200'
  } else if (props.calculations.varianceFromTarget > 0) {
    return 'bg-blue-50 border border-blue-200'
  } else {
    return 'bg-orange-50 border border-orange-200'
  }
})

const varianceTextClass = computed(() => {
  if (props.calculations.varianceFromTarget === 0) {
    return 'text-green-800'
  } else if (props.calculations.varianceFromTarget > 0) {
    return 'text-blue-800'
  } else {
    return 'text-orange-800'
  }
})

const varianceIconBg = computed(() => {
  if (props.calculations.varianceFromTarget === 0) {
    return 'bg-green-500'
  } else if (props.calculations.varianceFromTarget > 0) {
    return 'bg-blue-500'
  } else {
    return 'bg-orange-500'
  }
})

const varianceIconColor = computed(() => {
  return 'text-white'
})
</script>
