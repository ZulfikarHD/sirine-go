<template>
  <AppLayout>
    <div class="max-w-2xl mx-auto px-4 py-6 space-y-6">
      <!-- Breadcrumbs -->
      <Motion v-bind="entranceAnimations.fadeUp">
        <Breadcrumbs :items="breadcrumbs" />
      </Motion>

      <!-- Header -->
      <Motion v-bind="entranceAnimations.fadeUp">
        <div>
          <h1 class="text-2xl sm:text-3xl font-bold bg-gradient-to-r from-indigo-600 to-fuchsia-600 bg-clip-text text-transparent">
            Input Hasil Pemotongan
          </h1>
          <p class="text-gray-600 mt-1">Masukkan hasil sisiran kiri dan kanan</p>
        </div>
      </Motion>

      <!-- Loading State -->
      <div v-if="isLoading" class="space-y-4">
        <LoadingSkeleton class="h-48" />
        <LoadingSkeleton class="h-64" />
      </div>

      <!-- Content -->
      <template v-else-if="cuttingData">
        <!-- PO Info Card -->
        <Motion v-bind="entranceAnimations.fadeScale">
          <div class="glass-card p-6">
            <div class="flex items-center gap-3 mb-4">
              <div class="w-12 h-12 rounded-full bg-gradient-to-br from-indigo-500 to-fuchsia-500 flex items-center justify-center">
                <svg class="w-6 h-6 text-white" fill="currentColor" viewBox="0 0 20 20">
                  <path d="M9 2a1 1 0 000 2h2a1 1 0 100-2H9z"/>
                  <path fill-rule="evenodd" d="M4 5a2 2 0 012-2 3 3 0 003 3h2a3 3 0 003-3 2 2 0 012 2v11a2 2 0 01-2 2H6a2 2 0 01-2-2V5zm3 4a1 1 0 000 2h.01a1 1 0 100-2H7zm3 0a1 1 0 000 2h3a1 1 0 100-2h-3zm-3 4a1 1 0 100 2h.01a1 1 0 100-2H7zm3 0a1 1 0 100 2h3a1 1 0 100-2h-3z" clip-rule="evenodd"/>
                </svg>
              </div>
              <div class="flex-1">
                <h3 class="text-lg font-bold text-gray-900">PO {{ cuttingData.po?.po_number }}</h3>
                <p class="text-sm text-gray-500">{{ cuttingData.po?.obc_number }}</p>
              </div>
              <PriorityBadge :priority="cuttingData.po?.priority" />
            </div>

            <!-- Input & Output Info -->
            <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
              <div class="bg-gradient-to-br from-indigo-50 to-fuchsia-50 rounded-lg p-4">
                <p class="text-sm text-gray-600 mb-1">Input (Lembar Besar)</p>
                <p class="text-3xl font-bold text-indigo-600">{{ formatNumber(cuttingData.input_lembar_besar) }}</p>
                <p class="text-xs text-gray-500 mt-1">Dari hasil counting</p>
              </div>

              <div class="bg-gradient-to-br from-emerald-50 to-cyan-50 rounded-lg p-4">
                <p class="text-sm text-gray-600 mb-1">Estimasi Output</p>
                <p class="text-3xl font-bold text-emerald-600">{{ formatNumber(cuttingData.expected_output) }}</p>
                <p class="text-xs text-gray-500 mt-1">Input × 2 (sisiran kiri + kanan)</p>
              </div>
            </div>
          </div>
        </Motion>

        <!-- Result Form Card -->
        <Motion v-bind="entranceAnimations.fadeScale">
          <div class="glass-card p-6 space-y-5">
            <h3 class="text-lg font-bold text-gray-900 mb-4">Hasil Pemotongan</h3>

            <!-- Sisiran Kiri Input -->
            <div>
              <label class="block text-sm font-semibold text-gray-700 mb-2">
                <svg class="w-4 h-4 inline mr-1" fill="currentColor" viewBox="0 0 20 20">
                  <path d="M10 18a8 8 0 100-16 8 8 0 000 16zm1-12a1 1 0 10-2 0v4a1 1 0 00.293.707l2.828 2.829a1 1 0 101.415-1.415L11 9.586V6z"/>
                </svg>
                Sisiran Kiri <span class="text-red-500">*</span>
              </label>
              <input
                type="number"
                v-model.number="formData.output_sisiran_kiri"
                :disabled="isSubmitting"
                min="0"
                @input="validateInput"
                class="w-full px-4 py-3 border-2 border-gray-300 rounded-xl focus:border-indigo-500 focus:ring-4 focus:ring-indigo-100 transition-all"
                :class="{ 'opacity-50 cursor-not-allowed': isSubmitting }"
                placeholder="Masukkan jumlah sisiran kiri"
              />
            </div>

            <!-- Sisiran Kanan Input -->
            <div>
              <label class="block text-sm font-semibold text-gray-700 mb-2">
                <svg class="w-4 h-4 inline mr-1" fill="currentColor" viewBox="0 0 20 20">
                  <path d="M10 18a8 8 0 100-16 8 8 0 000 16zm1-12a1 1 0 10-2 0v4a1 1 0 00.293.707l2.828 2.829a1 1 0 101.415-1.415L11 9.586V6z"/>
                </svg>
                Sisiran Kanan <span class="text-red-500">*</span>
              </label>
              <input
                type="number"
                v-model.number="formData.output_sisiran_kanan"
                :disabled="isSubmitting"
                min="0"
                @input="validateInput"
                class="w-full px-4 py-3 border-2 border-gray-300 rounded-xl focus:border-indigo-500 focus:ring-4 focus:ring-indigo-100 transition-all"
                :class="{ 'opacity-50 cursor-not-allowed': isSubmitting }"
                placeholder="Masukkan jumlah sisiran kanan"
              />
            </div>

            <!-- Calculation Results -->
            <div v-if="hasInputValues" class="space-y-3 pt-2">
              <div class="bg-gradient-to-br from-gray-50 to-gray-100 rounded-lg p-4 space-y-2">
                <div class="flex justify-between items-center">
                  <span class="text-sm font-semibold text-gray-700">Total Output:</span>
                  <span class="text-lg font-bold text-gray-900">{{ formatNumber(calculatedTotalOutput) }} lembar</span>
                </div>
                <div class="flex justify-between items-center">
                  <span class="text-sm font-semibold text-gray-700">Waste:</span>
                  <span class="text-lg font-bold text-gray-900">{{ formatNumber(calculatedWaste) }} lembar</span>
                </div>
              </div>

              <!-- Waste Percentage Indicator -->
              <div
                class="rounded-lg p-4 flex items-center gap-3"
                :class="wasteExceedsThreshold ? 'bg-red-50 border-2 border-red-200' : 'bg-emerald-50 border-2 border-emerald-200'"
              >
                <div
                  class="w-12 h-12 rounded-full flex items-center justify-center flex-shrink-0"
                  :class="wasteExceedsThreshold ? 'bg-red-500' : 'bg-emerald-500'"
                >
                  <svg v-if="wasteExceedsThreshold" class="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"/>
                  </svg>
                  <svg v-else class="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/>
                  </svg>
                </div>
                <div class="flex-1">
                  <p class="text-sm font-semibold" :class="wasteExceedsThreshold ? 'text-red-900' : 'text-emerald-900'">
                    Waste Percentage: {{ calculatedWastePercentage }}%
                  </p>
                  <p class="text-xs mt-1" :class="wasteExceedsThreshold ? 'text-red-700' : 'text-emerald-700'">
                    {{ wasteExceedsThreshold ? '⚠️ Waste melebihi toleransi 2%' : '✓ Waste dalam toleransi' }}
                  </p>
                </div>
              </div>
            </div>

            <!-- Waste Documentation (conditional) -->
            <div v-if="wasteExceedsThreshold" class="space-y-4 pt-2">
              <div class="bg-amber-50 border-2 border-amber-200 rounded-lg p-4">
                <div class="flex gap-2 mb-2">
                  <svg class="w-5 h-5 text-amber-600 flex-shrink-0" fill="currentColor" viewBox="0 0 20 20">
                    <path fill-rule="evenodd" d="M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.213 2.98-1.742 2.98H4.42c-1.53 0-2.493-1.646-1.743-2.98l5.58-9.92zM11 13a1 1 0 11-2 0 1 1 0 012 0zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z" clip-rule="evenodd"/>
                  </svg>
                  <div>
                    <p class="text-sm font-semibold text-amber-900">Dokumentasi Waste Diperlukan</p>
                    <p class="text-xs text-amber-700 mt-1">Waste melebihi 2%, wajib mengisi alasan dan foto bukti</p>
                  </div>
                </div>
              </div>

              <!-- Waste Reason -->
              <div>
                <label class="block text-sm font-semibold text-gray-700 mb-2">
                  Alasan Waste <span class="text-red-500">*</span>
                </label>
                <textarea
                  v-model="formData.waste_reason"
                  :disabled="isSubmitting"
                  rows="3"
                  class="w-full px-4 py-3 border-2 border-gray-300 rounded-xl focus:border-indigo-500 focus:ring-4 focus:ring-indigo-100 transition-all resize-none"
                  :class="{ 'opacity-50 cursor-not-allowed': isSubmitting }"
                  placeholder="Jelaskan penyebab waste (minimal 10 karakter)"
                ></textarea>
                <p v-if="formData.waste_reason" class="text-xs text-gray-500 mt-1">
                  {{ formData.waste_reason.length }} karakter
                </p>
              </div>

              <!-- Waste Photo Upload -->
              <WastePhotoUpload
                v-model="formData.waste_photo_url"
                :required="true"
                :disabled="isSubmitting"
              />
            </div>

            <!-- Info Alert -->
            <div class="bg-blue-50 border-2 border-blue-200 rounded-xl p-4 flex gap-3">
              <svg class="w-5 h-5 text-blue-600 flex-shrink-0 mt-0.5" fill="currentColor" viewBox="0 0 20 20">
                <path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z" clip-rule="evenodd"/>
              </svg>
              <div class="text-sm text-blue-900">
                <p class="font-semibold mb-1">Catatan:</p>
                <ul class="list-disc list-inside space-y-1 text-blue-800">
                  <li>Pastikan data yang diinput sudah benar</li>
                  <li>Waste akan dihitung otomatis dari estimasi output</li>
                  <li>Data dapat disimpan untuk dilanjutkan kemudian</li>
                </ul>
              </div>
            </div>
          </div>
        </Motion>

        <!-- Action Buttons -->
        <Motion v-bind="entranceAnimations.fadeUp">
          <div class="flex flex-col sm:flex-row gap-3">
            <!-- Cancel Button -->
            <button
              @click="handleCancel"
              :disabled="isSubmitting"
              class="flex-1 py-3 px-6 bg-white border-2 border-gray-300 text-gray-700 font-semibold rounded-xl hover:bg-gray-50 transition-all active-scale"
              :class="{ 'opacity-50 cursor-not-allowed': isSubmitting }"
            >
              Batal
            </button>

            <!-- Submit Button -->
            <button
              @click="handleSubmit"
              :disabled="!canSubmit || isSubmitting"
              class="flex-1 py-3 px-6 bg-gradient-to-r from-indigo-600 to-fuchsia-600 text-white font-semibold rounded-xl hover:from-indigo-700 hover:to-fuchsia-700 transition-all active-scale flex items-center justify-center gap-2"
              :class="{ 'opacity-50 cursor-not-allowed': !canSubmit || isSubmitting }"
            >
              <svg v-if="isSubmitting" class="animate-spin w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
              <svg v-else class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/>
              </svg>
              <span>{{ isSubmitting ? 'Menyimpan...' : 'Simpan Hasil' }}</span>
            </button>
          </div>
        </Motion>
      </template>

      <!-- Error State -->
      <Motion v-else v-bind="entranceAnimations.fadeScale">
        <div class="glass-card p-12 text-center">
          <div class="w-20 h-20 mx-auto mb-4 rounded-full bg-gradient-to-br from-red-100 to-red-200 flex items-center justify-center">
            <svg class="w-10 h-10 text-red-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
            </svg>
          </div>
          <h3 class="text-xl font-bold text-gray-900 mb-2">Data Tidak Ditemukan</h3>
          <p class="text-gray-600 mb-4">Cutting record tidak tersedia</p>
          <button
            @click="router.push({ name: 'cutting-queue' })"
            class="px-6 py-2 bg-gradient-to-r from-indigo-600 to-fuchsia-600 text-white font-semibold rounded-xl hover:from-indigo-700 hover:to-fuchsia-700 transition-all active-scale"
          >
            Kembali ke Queue
          </button>
        </div>
      </Motion>
    </div>
  </AppLayout>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { Motion } from 'motion-v'
import { entranceAnimations } from '@/composables/useMotion'
import { useCuttingApi } from '@/composables/useCuttingApi'
import { useHaptic } from '@/composables/useHaptic'
import AppLayout from '@/components/layout/AppLayout.vue'
import { Breadcrumbs, LoadingSkeleton, PriorityBadge } from '@/components/common'
import WastePhotoUpload from '@/components/cutting/WastePhotoUpload.vue'

const router = useRouter()
const route = useRoute()
const cuttingApi = useCuttingApi()
const haptic = useHaptic()

// State
const isLoading = ref(false)
const isSubmitting = ref(false)
const cuttingData = ref(null)
const formData = ref({
  output_sisiran_kiri: null,
  output_sisiran_kanan: null,
  waste_reason: '',
  waste_photo_url: '',
})

// Computed - Calculations
const hasInputValues = computed(() => {
  return formData.value.output_sisiran_kiri !== null && 
         formData.value.output_sisiran_kanan !== null &&
         formData.value.output_sisiran_kiri >= 0 &&
         formData.value.output_sisiran_kanan >= 0
})

const calculatedTotalOutput = computed(() => {
  if (!hasInputValues.value) return 0
  return formData.value.output_sisiran_kiri + formData.value.output_sisiran_kanan
})

const calculatedWaste = computed(() => {
  if (!hasInputValues.value || !cuttingData.value) return 0
  return cuttingData.value.expected_output - calculatedTotalOutput.value
})

const calculatedWastePercentage = computed(() => {
  if (!hasInputValues.value || !cuttingData.value || cuttingData.value.expected_output === 0) return '0.00'
  const percentage = (calculatedWaste.value / cuttingData.value.expected_output) * 100
  return percentage.toFixed(2)
})

const wasteExceedsThreshold = computed(() => {
  return parseFloat(calculatedWastePercentage.value) > 2.0
})

// Computed - Validation
const canSubmit = computed(() => {
  // Basic validation
  if (!hasInputValues.value || isSubmitting.value) return false
  
  // If waste exceeds threshold, require reason and photo
  if (wasteExceedsThreshold.value) {
    const hasReason = formData.value.waste_reason.trim().length >= 10
    const hasPhoto = formData.value.waste_photo_url.trim() !== ''
    return hasReason && hasPhoto
  }
  
  return true
})

const breadcrumbs = computed(() => [
  { label: 'Pemotongan', path: '/khazwal/cutting' },
  { label: `Input Hasil`, path: '' },
])

/**
 * Format number dengan thousand separator
 */
const formatNumber = (num) => {
  return new Intl.NumberFormat('id-ID').format(num)
}

/**
 * Validate input untuk mencegah nilai negatif
 */
const validateInput = () => {
  if (formData.value.output_sisiran_kiri < 0) {
    formData.value.output_sisiran_kiri = 0
  }
  if (formData.value.output_sisiran_kanan < 0) {
    formData.value.output_sisiran_kanan = 0
  }
}

/**
 * Fetch cutting data dari API
 */
const fetchCuttingData = async () => {
  isLoading.value = true
  
  try {
    const response = await cuttingApi.getCuttingDetail(route.params.id)
    cuttingData.value = response
    
    // Pre-fill existing data if available
    if (response.output_sisiran_kiri !== null) {
      formData.value.output_sisiran_kiri = response.output_sisiran_kiri
    }
    if (response.output_sisiran_kanan !== null) {
      formData.value.output_sisiran_kanan = response.output_sisiran_kanan
    }
    if (response.waste_reason) {
      formData.value.waste_reason = response.waste_reason
    }
    if (response.waste_photo_url) {
      formData.value.waste_photo_url = response.waste_photo_url
    }
  } catch (error) {
    console.error('Failed to fetch cutting data:', error)
    cuttingData.value = null
  } finally {
    isLoading.value = false
  }
}

/**
 * Handle submit - save cutting results
 */
const handleSubmit = async () => {
  if (!canSubmit.value) return
  
  haptic.medium()
  isSubmitting.value = true
  
  try {
    const payload = {
      output_sisiran_kiri: formData.value.output_sisiran_kiri,
      output_sisiran_kanan: formData.value.output_sisiran_kanan,
    }
    
    // Add waste documentation if required
    if (wasteExceedsThreshold.value) {
      payload.waste_reason = formData.value.waste_reason
      payload.waste_photo_url = formData.value.waste_photo_url
    }
    
    const response = await cuttingApi.updateCuttingResult(route.params.id, payload)
    
    // Success haptic
    haptic.success()
    
    // TODO: Show success toast
    console.log('Cutting result saved successfully:', response)
    
    // Redirect back to queue
    router.push({ name: 'cutting-queue' })
  } catch (error) {
    console.error('Failed to save cutting result:', error)
    haptic.error()
    
    // TODO: Show error toast with specific message
    if (error.response?.status === 400) {
      alert('Data tidak valid: ' + (error.response?.data?.error || error.message))
    } else {
      alert('Gagal menyimpan hasil pemotongan: ' + (error.response?.data?.error || error.message))
    }
  } finally {
    isSubmitting.value = false
  }
}

/**
 * Handle cancel - back to queue
 */
const handleCancel = () => {
  haptic.light()
  router.push({ name: 'cutting-queue' })
}

// Lifecycle
onMounted(() => {
  fetchCuttingData()
})
</script>
