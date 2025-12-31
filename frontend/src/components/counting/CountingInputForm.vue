<template>
  <Motion v-bind="entranceAnimations.fadeScale">
    <div class="glass-card p-4">
      <h3 class="text-lg font-bold text-gray-900 mb-4">Input Hasil Penghitungan</h3>

      <div class="space-y-4">
        <!-- Quantity Good -->
        <div>
          <label class="block text-sm font-semibold text-gray-700 mb-2">
            Jumlah Baik
            <span class="text-red-500">*</span>
          </label>
          <div class="relative">
            <input
              :value="modelValue.quantity_good"
              @input="$emit('update:quantity-good', $event.target.value)"
              type="number"
              inputmode="numeric"
              min="0"
              class="w-full px-4 py-4 text-lg font-semibold border-2 border-gray-300 rounded-xl focus:border-indigo-500 focus:ring-4 focus:ring-indigo-100 transition-all outline-none"
              :class="{ 'border-red-300 bg-red-50': error.quantity_good }"
              placeholder="0"
              :disabled="isLoading"
            />
            <div class="absolute right-4 top-1/2 -translate-y-1/2 text-gray-400 font-medium">
              LB
            </div>
          </div>
          <p v-if="error.quantity_good" class="mt-1 text-sm text-red-600">{{ error.quantity_good }}</p>
        </div>

        <!-- Quantity Defect -->
        <div>
          <label class="block text-sm font-semibold text-gray-700 mb-2">
            Jumlah Rusak
            <span class="text-red-500">*</span>
          </label>
          <div class="relative">
            <input
              :value="modelValue.quantity_defect"
              @input="$emit('update:quantity-defect', $event.target.value)"
              type="number"
              inputmode="numeric"
              min="0"
              class="w-full px-4 py-4 text-lg font-semibold border-2 border-gray-300 rounded-xl focus:border-indigo-500 focus:ring-4 focus:ring-indigo-100 transition-all outline-none"
              :class="{ 'border-red-300 bg-red-50': error.quantity_defect }"
              placeholder="0"
              :disabled="isLoading"
            />
            <div class="absolute right-4 top-1/2 -translate-y-1/2 text-gray-400 font-medium">
              LB
            </div>
          </div>
          <p v-if="error.quantity_defect" class="mt-1 text-sm text-red-600">{{ error.quantity_defect }}</p>
        </div>

        <!-- Help Text -->
        <div class="flex items-start gap-2 p-3 bg-blue-50 border border-blue-200 rounded-lg">
          <svg class="w-5 h-5 text-blue-600 flex-shrink-0 mt-0.5" fill="currentColor" viewBox="0 0 20 20">
            <path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z" clip-rule="evenodd"/>
          </svg>
          <div class="text-sm text-blue-800">
            <p class="font-semibold mb-1">Petunjuk Penghitungan</p>
            <ul class="list-disc list-inside space-y-1 text-blue-700">
              <li>Hitung jumlah lembar besar yang baik (tidak ada cacat)</li>
              <li>Hitung jumlah lembar besar yang rusak (ada cacat apapun)</li>
              <li>Hasil akan otomatis disimpan</li>
            </ul>
          </div>
        </div>

        <!-- Loading indicator -->
        <div v-if="isLoading" class="flex items-center justify-center gap-2 p-3 bg-indigo-50 rounded-lg">
          <svg class="animate-spin h-5 w-5 text-indigo-600" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
          </svg>
          <span class="text-sm font-medium text-indigo-600">Menyimpan...</span>
        </div>
      </div>
    </div>
  </Motion>
</template>

<script setup>
import { Motion } from 'motion-v'
import { entranceAnimations } from '@/composables/useMotion'

defineProps({
  modelValue: {
    type: Object,
    required: true
  },
  isLoading: {
    type: Boolean,
    default: false
  },
  error: {
    type: Object,
    default: () => ({})
  }
})

defineEmits(['update:quantity-good', 'update:quantity-defect'])
</script>
