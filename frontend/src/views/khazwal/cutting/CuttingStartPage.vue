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
            Mulai Pemotongan
          </h1>
          <p class="text-gray-600 mt-1">Konfirmasi detail sebelum memulai proses pemotongan</p>
        </div>
      </Motion>

      <!-- Loading State -->
      <div v-if="isLoading" class="space-y-4">
        <LoadingSkeleton class="h-32" />
        <LoadingSkeleton class="h-48" />
      </div>

      <!-- Content -->
      <template v-else-if="poData">
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
                <h3 class="text-lg font-bold text-gray-900">PO {{ poData.po_number }}</h3>
                <p class="text-sm text-gray-500">{{ poData.obc_number }}</p>
              </div>
              <PriorityBadge :priority="poData.priority" />
            </div>

            <!-- Input & Output Info -->
            <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
              <div class="bg-gradient-to-br from-indigo-50 to-fuchsia-50 rounded-lg p-4">
                <p class="text-sm text-gray-600 mb-1">Input (Lembar Besar)</p>
                <p class="text-3xl font-bold text-indigo-600">{{ formatNumber(poData.input_lembar_besar) }}</p>
                <p class="text-xs text-gray-500 mt-1">Dari hasil counting</p>
              </div>

              <div class="bg-gradient-to-br from-emerald-50 to-cyan-50 rounded-lg p-4">
                <p class="text-sm text-gray-600 mb-1">Estimasi Output</p>
                <p class="text-3xl font-bold text-emerald-600">{{ formatNumber(poData.estimated_output) }}</p>
                <p class="text-xs text-gray-500 mt-1">Input × 2 (sisiran kiri + kanan)</p>
              </div>
            </div>
          </div>
        </Motion>

        <!-- Form Card -->
        <Motion v-bind="entranceAnimations.fadeScale">
          <div class="glass-card p-6 space-y-5">
            <h3 class="text-lg font-bold text-gray-900 mb-4">Detail Pemotongan</h3>

            <!-- Operator (Auto-filled) -->
            <div>
              <label class="block text-sm font-semibold text-gray-700 mb-2">
                <svg class="w-4 h-4 inline mr-1" fill="currentColor" viewBox="0 0 20 20">
                  <path fill-rule="evenodd" d="M10 9a3 3 0 100-6 3 3 0 000 6zm-7 9a7 7 0 1114 0H3z" clip-rule="evenodd"/>
                </svg>
                Operator
              </label>
              <div class="px-4 py-3 bg-gray-100 border-2 border-gray-300 rounded-xl text-gray-700 font-semibold">
                {{ operatorName }}
              </div>
              <p class="text-xs text-gray-500 mt-1">Terisi otomatis dari akun login</p>
            </div>

            <!-- Cutting Machine Selector -->
            <div>
              <label class="block text-sm font-semibold text-gray-700 mb-2">
                <svg class="w-4 h-4 inline mr-1" fill="currentColor" viewBox="0 0 20 20">
                  <path fill-rule="evenodd" d="M12.316 3.051a1 1 0 01.633 1.265l-4 12a1 1 0 11-1.898-.632l4-12a1 1 0 011.265-.633zM5.707 6.293a1 1 0 010 1.414L3.414 10l2.293 2.293a1 1 0 11-1.414 1.414l-3-3a1 1 0 010-1.414l3-3a1 1 0 011.414 0zm8.586 0a1 1 0 011.414 0l3 3a1 1 0 010 1.414l-3 3a1 1 0 11-1.414-1.414L16.586 10l-2.293-2.293a1 1 0 010-1.414z" clip-rule="evenodd"/>
                </svg>
                Mesin Potong <span class="text-red-500">*</span>
              </label>
              <select
                v-model="formData.cutting_machine"
                :disabled="isSubmitting"
                class="w-full px-4 py-3 border-2 border-gray-300 rounded-xl focus:border-indigo-500 focus:ring-4 focus:ring-indigo-100 transition-all"
                :class="{ 'opacity-50 cursor-not-allowed': isSubmitting }"
              >
                <option value="">-- Pilih Mesin Potong --</option>
                <option value="Mesin A">Mesin A</option>
                <option value="Mesin B">Mesin B</option>
                <option value="Mesin C">Mesin C</option>
              </select>
              <p v-if="!formData.cutting_machine" class="text-xs text-gray-500 mt-1">
                Pilih mesin yang akan digunakan untuk pemotongan
              </p>
            </div>

            <!-- Info Alert -->
            <div class="bg-blue-50 border-2 border-blue-200 rounded-xl p-4 flex gap-3">
              <svg class="w-5 h-5 text-blue-600 flex-shrink-0 mt-0.5" fill="currentColor" viewBox="0 0 20 20">
                <path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z" clip-rule="evenodd"/>
              </svg>
              <div class="text-sm text-blue-900">
                <p class="font-semibold mb-1">Catatan Penting:</p>
                <ul class="list-disc list-inside space-y-1 text-blue-800">
                  <li>Proses pemotongan akan dimulai setelah konfirmasi</li>
                  <li>Status PO akan berubah menjadi "Sedang Dipotong"</li>
                  <li>Anda akan melanjutkan ke input hasil setelah pemotongan selesai</li>
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
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14.121 14.121L19 19m-7-7l7-7m-7 7l-2.879 2.879M12 12L9.121 9.121m0 5.758a3 3 0 10-4.243 4.243 3 3 0 004.243-4.243zm0-5.758a3 3 0 10-4.243-4.243 3 3 0 004.243 4.243z"/>
              </svg>
              <span>{{ isSubmitting ? 'Memulai...' : 'Mulai Pemotongan' }}</span>
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
          <h3 class="text-xl font-bold text-gray-900 mb-2">PO Tidak Ditemukan</h3>
          <p class="text-gray-600 mb-4">PO tidak tersedia atau sudah diproses</p>
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
import { useAuthStore } from '@/stores/auth'
import AppLayout from '@/components/layout/AppLayout.vue'
import { Breadcrumbs, LoadingSkeleton, PriorityBadge } from '@/components/common'

const router = useRouter()
const route = useRoute()
const cuttingApi = useCuttingApi()
const haptic = useHaptic()
const authStore = useAuthStore()

// State
const isLoading = ref(false)
const isSubmitting = ref(false)
const poData = ref(null)
const formData = ref({
  cutting_machine: '',
})

// Computed
const operatorName = computed(() => {
  return authStore.user?.full_name || 'Unknown'
})

const canSubmit = computed(() => {
  return formData.value.cutting_machine !== '' && !isSubmitting.value
})

const breadcrumbs = computed(() => [
  { label: 'Pemotongan', path: '/khazwal/cutting' },
  { label: `PO ${poData.value?.po_number || route.params.poId}`, path: '' },
])

/**
 * Format number dengan thousand separator
 */
const formatNumber = (num) => {
  return new Intl.NumberFormat('id-ID').format(num)
}

/**
 * Fetch PO data dari queue
 */
const fetchPOData = async () => {
  isLoading.value = true
  
  try {
    const response = await cuttingApi.getCuttingQueue()
    const queue = response.data || []
    
    // Find PO by ID
    const po = queue.find(item => item.po_id === parseInt(route.params.poId))
    
    if (po) {
      poData.value = po
    } else {
      poData.value = null
    }
  } catch (error) {
    console.error('Failed to fetch PO data:', error)
    poData.value = null
  } finally {
    isLoading.value = false
  }
}

/**
 * Handle submit - start cutting
 */
const handleSubmit = async () => {
  if (!canSubmit.value) return
  
  haptic.medium()
  isSubmitting.value = true
  
  try {
    const response = await cuttingApi.startCutting(route.params.poId, {
      cutting_machine: formData.value.cutting_machine
    })
    
    // Success haptic
    haptic.success()
    
    // TODO: Show success toast
    console.log('Cutting started successfully:', response)
    
    // Redirect back to queue
    router.push({ name: 'cutting-queue' })
  } catch (error) {
    console.error('Failed to start cutting:', error)
    haptic.error()
    
    // TODO: Show error toast with specific message
    if (error.response?.status === 409) {
      alert('PO ini sudah dimulai oleh user lain')
    } else if (error.response?.status === 400) {
      alert('PO belum siap untuk dipotong')
    } else {
      alert('Gagal memulai pemotongan: ' + (error.response?.data?.error || error.message))
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
  fetchPOData()
})
</script>
