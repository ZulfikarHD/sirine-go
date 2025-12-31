<template>
  <Motion v-bind="entranceAnimations.fadeScale">
    <div class="bg-red-50 border-2 border-red-200 rounded-xl p-4">
      <div class="flex items-start gap-3 mb-4">
        <div class="w-10 h-10 rounded-full bg-red-500 flex items-center justify-center flex-shrink-0">
          <svg class="w-5 h-5 text-white" fill="currentColor" viewBox="0 0 20 20">
            <path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7 4a1 1 0 11-2 0 1 1 0 012 0zm-1-9a1 1 0 00-1 1v4a1 1 0 102 0V6a1 1 0 00-1-1z" clip-rule="evenodd"/>
          </svg>
        </div>
        <div class="flex-1">
          <h3 class="text-base font-bold text-red-900 mb-1">Breakdown Kerusakan</h3>
          <p class="text-sm text-red-700">Wajib diisi karena persentase rusak >5%</p>
        </div>
      </div>

      <div class="space-y-3">
        <!-- Defect Type Inputs -->
        <div 
          v-for="defectType in defectTypes" 
          :key="defectType.type"
          class="bg-white rounded-lg p-3"
        >
          <label class="block text-sm font-semibold text-gray-700 mb-2">
            {{ defectType.label }}
          </label>
          <div class="relative">
            <input
              :value="getBreakdownQuantity(defectType.type)"
              @input="handleBreakdownInput(defectType.type, $event.target.value)"
              type="number"
              inputmode="numeric"
              min="0"
              class="w-full px-4 py-3 text-base font-semibold border-2 border-gray-300 rounded-lg focus:border-red-500 focus:ring-4 focus:ring-red-100 transition-all outline-none"
              placeholder="0"
            />
            <div class="absolute right-4 top-1/2 -translate-y-1/2 text-gray-400 font-medium text-sm">
              lembar
            </div>
          </div>
        </div>

        <!-- Summary & Validation -->
        <div class="mt-4 p-4 rounded-lg" :class="validationClass">
          <div class="flex items-center justify-between mb-2">
            <span class="text-sm font-semibold" :class="validationTextClass">Total Breakdown</span>
            <span class="text-lg font-bold" :class="validationTextClass">{{ breakdownSum }} / {{ totalDefect }}</span>
          </div>
          
          <div class="h-2 bg-white/50 rounded-full overflow-hidden">
            <div 
              class="h-full transition-all duration-300"
              :class="progressBarClass"
              :style="{ width: progressPercentage + '%' }"
            ></div>
          </div>

          <!-- Validation Message -->
          <div v-if="!isValid" class="flex items-center gap-2 mt-3 p-2 bg-white rounded-lg">
            <svg class="w-5 h-5 text-red-600 flex-shrink-0" fill="currentColor" viewBox="0 0 20 20">
              <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd"/>
            </svg>
            <p class="text-sm font-medium text-red-900">
              {{ validationMessage }}
            </p>
          </div>

          <div v-else class="flex items-center gap-2 mt-3 p-2 bg-white rounded-lg">
            <svg class="w-5 h-5 text-green-600 flex-shrink-0" fill="currentColor" viewBox="0 0 20 20">
              <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd"/>
            </svg>
            <p class="text-sm font-medium text-green-900">
              Breakdown valid ✓
            </p>
          </div>
        </div>

        <!-- Help Text -->
        <div class="flex items-start gap-2 p-3 bg-white rounded-lg">
          <svg class="w-5 h-5 text-blue-600 flex-shrink-0 mt-0.5" fill="currentColor" viewBox="0 0 20 20">
            <path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z" clip-rule="evenodd"/>
          </svg>
          <div class="text-sm text-gray-700">
            <p class="font-semibold mb-1">Petunjuk</p>
            <p>Total breakdown harus sama dengan jumlah rusak ({{ totalDefect }} lembar)</p>
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
import { useCountingApi } from '@/composables/useCountingApi'

const props = defineProps({
  modelValue: {
    type: Array,
    default: () => []
  },
  totalDefect: {
    type: Number,
    required: true
  }
})

const emit = defineEmits(['update:modelValue'])

const { defectTypes } = useCountingApi()

const breakdownSum = computed(() => {
  return props.modelValue.reduce((sum, item) => sum + (item.quantity || 0), 0)
})

const isValid = computed(() => {
  return breakdownSum.value === props.totalDefect
})

const progressPercentage = computed(() => {
  if (props.totalDefect === 0) return 0
  return Math.min((breakdownSum.value / props.totalDefect) * 100, 100)
})

const validationClass = computed(() => {
  if (isValid.value) {
    return 'bg-green-100 border-2 border-green-300'
  } else if (breakdownSum.value > props.totalDefect) {
    return 'bg-red-100 border-2 border-red-300'
  } else {
    return 'bg-yellow-100 border-2 border-yellow-300'
  }
})

const validationTextClass = computed(() => {
  if (isValid.value) return 'text-green-900'
  if (breakdownSum.value > props.totalDefect) return 'text-red-900'
  return 'text-yellow-900'
})

const progressBarClass = computed(() => {
  if (isValid.value) return 'bg-green-500'
  if (breakdownSum.value > props.totalDefect) return 'bg-red-500'
  return 'bg-yellow-500'
})

const validationMessage = computed(() => {
  const diff = breakdownSum.value - props.totalDefect
  if (diff > 0) {
    return `Kelebihan ${diff} lembar`
  } else {
    return `Kurang ${Math.abs(diff)} lembar lagi`
  }
})

const getBreakdownQuantity = (type) => {
  const item = props.modelValue.find(i => i.type === type)
  return item?.quantity || ''
}

const handleBreakdownInput = (type, value) => {
  const numValue = parseInt(value) || 0
  const newBreakdown = [...props.modelValue]
  
  const existingIndex = newBreakdown.findIndex(item => item.type === type)
  
  if (numValue <= 0) {
    if (existingIndex !== -1) {
      newBreakdown.splice(existingIndex, 1)
    }
  } else {
    if (existingIndex !== -1) {
      newBreakdown[existingIndex].quantity = numValue
    } else {
      newBreakdown.push({ type, quantity: numValue })
    }
  }
  
  emit('update:modelValue', newBreakdown)
}
</script>
