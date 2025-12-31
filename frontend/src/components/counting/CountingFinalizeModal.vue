<template>
  <ConfirmDialog
    v-model="isOpen"
    title="Selesaikan Penghitungan"
    message="Pastikan semua data sudah benar sebelum finalisasi."
    variant="warning"
    confirm-text="Ya, Selesaikan"
    cancel-text="Cek Lagi"
    :show-warning="true"
    warning-message="Data tidak dapat diubah setelah finalisasi"
    :loading="loading"
    @confirm="handleConfirm"
    @cancel="handleCancel"
  >
    <!-- Custom Content Slot - Ringkasan -->
    <div class="space-y-4 mt-4">
      <div class="bg-gray-50 rounded-xl p-4">
        <h4 class="text-sm font-bold text-gray-700 mb-3">Ringkasan Hasil</h4>
        
        <div class="space-y-3">
          <!-- Good & Defect -->
          <div class="grid grid-cols-2 gap-3">
            <div class="bg-green-50 border border-green-200 rounded-lg p-3">
              <p class="text-xs text-green-600 mb-1">Baik</p>
              <p class="text-xl font-bold text-green-800">
                {{ summary.quantityGood }}
                <span class="text-xs font-normal text-green-600">({{ summary.percentageGood }}%)</span>
              </p>
            </div>
            <div class="bg-red-50 border border-red-200 rounded-lg p-3">
              <p class="text-xs text-red-600 mb-1">Rusak</p>
              <p class="text-xl font-bold text-red-800">
                {{ summary.quantityDefect }}
                <span class="text-xs font-normal text-red-600">({{ summary.percentageDefect }}%)</span>
              </p>
            </div>
          </div>

          <!-- Total & Variance -->
          <div class="flex items-center justify-between p-3 bg-white border border-gray-200 rounded-lg">
            <span class="text-sm font-semibold text-gray-700">Total Dihitung</span>
            <span class="text-lg font-bold text-gray-900">{{ summary.totalCounted }} LB</span>
          </div>

          <div 
            class="flex items-center justify-between p-3 rounded-lg"
            :class="varianceClass"
          >
            <span class="text-sm font-semibold">Selisih</span>
            <span class="text-lg font-bold">
              {{ summary.varianceFromTarget > 0 ? '+' : '' }}{{ summary.varianceFromTarget }} LB
              <span class="text-sm font-normal">({{ summary.variancePercentage > 0 ? '+' : '' }}{{ summary.variancePercentage }}%)</span>
            </span>
          </div>

          <!-- Defect Breakdown if exists -->
          <div v-if="summary.defectBreakdown && summary.defectBreakdown.length > 0" class="bg-red-50 border border-red-200 rounded-lg p-3">
            <p class="text-xs font-semibold text-red-700 mb-2">Breakdown Kerusakan</p>
            <ul class="space-y-1">
              <li 
                v-for="item in summary.defectBreakdown" 
                :key="item.type"
                class="flex items-center justify-between text-sm"
              >
                <span class="text-red-800">{{ item.type }}</span>
                <span class="font-semibold text-red-900">{{ item.quantity }} lembar</span>
              </li>
            </ul>
          </div>

          <!-- Variance Reason if exists -->
          <div v-if="summary.varianceReason" class="bg-yellow-50 border border-yellow-200 rounded-lg p-3">
            <p class="text-xs font-semibold text-yellow-700 mb-1">Alasan Selisih</p>
            <p class="text-sm text-yellow-900">{{ summary.varianceReason }}</p>
          </div>
        </div>
      </div>

      <!-- Warning Box -->
      <div class="flex items-start gap-3 p-4 bg-red-50 border-2 border-red-200 rounded-xl">
        <svg class="w-6 h-6 text-red-600 flex-shrink-0 mt-0.5" fill="currentColor" viewBox="0 0 20 20">
          <path fill-rule="evenodd" d="M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.213 2.98-1.742 2.98H4.42c-1.53 0-2.493-1.646-1.743-2.98l5.58-9.92zM11 13a1 1 0 11-2 0 1 1 0 012 0zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z" clip-rule="evenodd"/>
        </svg>
        <div>
          <p class="text-sm font-bold text-red-900 mb-1">Perhatian!</p>
          <ul class="text-sm text-red-800 space-y-1">
            <li>• Data tidak dapat diubah setelah finalisasi</li>
            <li>• PO akan otomatis masuk ke antrian pemotongan</li>
            <li>• Pastikan semua angka sudah dicek dengan teliti</li>
          </ul>
        </div>
      </div>
    </div>
  </ConfirmDialog>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import ConfirmDialog from '@/components/common/ConfirmDialog.vue'

const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false
  },
  summary: {
    type: Object,
    required: true
  },
  loading: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['update:modelValue', 'confirm', 'cancel'])

const isOpen = ref(props.modelValue)

watch(() => props.modelValue, (value) => {
  isOpen.value = value
})

watch(isOpen, (value) => {
  emit('update:modelValue', value)
})

const varianceClass = computed(() => {
  if (props.summary.varianceFromTarget === 0) {
    return 'bg-green-50 border border-green-200 text-green-900'
  } else if (props.summary.varianceFromTarget > 0) {
    return 'bg-blue-50 border border-blue-200 text-blue-900'
  } else {
    return 'bg-orange-50 border border-orange-200 text-orange-900'
  }
})

const handleConfirm = () => {
  emit('confirm')
}

const handleCancel = () => {
  emit('cancel')
  isOpen.value = false
}
</script>
