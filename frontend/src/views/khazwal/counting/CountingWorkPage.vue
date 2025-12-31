<template>
  <AppLayout>
    <div class="max-w-4xl mx-auto px-4 py-6 space-y-6">
      <!-- Header with Back Button -->
      <Motion v-bind="entranceAnimations.fadeUp">
        <div class="flex items-center gap-4">
          <button
            @click="handleBack"
            class="p-2 hover:bg-gray-100 rounded-lg transition-colors active-scale"
          >
            <svg class="w-6 h-6 text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7"/>
            </svg>
          </button>
          <div class="flex-1">
            <h1 class="text-xl sm:text-2xl font-bold text-gray-900">
              PO {{ currentPO?.po_number || '-' }}
            </h1>
            <p v-if="currentCounting" class="text-sm text-gray-600 mt-1">
              Status: {{ statusLabel }} 
              <span v-if="currentCounting.started_at" class="text-indigo-600 font-medium">
                ({{ countingDuration }})
              </span>
            </p>
          </div>
        </div>
      </Motion>

      <!-- Loading State -->
      <div v-if="isLoadingDetail" class="space-y-4">
        <LoadingSkeleton v-for="i in 3" :key="i" type="card" />
      </div>

      <template v-else-if="currentPO">
        <!-- State 1: Before Start -->
        <template v-if="!currentCounting || currentCounting.status === 'PENDING'">
          <CountingPrintInfo :po-info="currentPO" :print-info="printInfo" />

          <Motion v-bind="entranceAnimations.fadeScale">
            <div class="glass-card p-6 text-center">
              <div class="w-20 h-20 mx-auto mb-4 rounded-full bg-gradient-to-br from-indigo-500 to-fuchsia-500 flex items-center justify-center">
                <svg class="w-10 h-10 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-6 9l2 2 4-4"/>
                </svg>
              </div>
              <h3 class="text-lg font-bold text-gray-900 mb-2">Siap Memulai Penghitungan?</h3>
              <p class="text-gray-600 mb-6">Pastikan material sudah siap untuk dihitung</p>
              
              <button
                @click="handleStartCounting"
                :disabled="isStarting"
                class="inline-flex items-center gap-2 px-6 py-3 bg-gradient-to-r from-indigo-600 to-fuchsia-600 text-white font-bold rounded-xl hover:shadow-lg transition-all active-scale"
                :class="{ 'opacity-50 cursor-not-allowed': isStarting }"
              >
                <svg v-if="!isStarting" class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14.752 11.168l-3.197-2.132A1 1 0 0010 9.87v4.263a1 1 0 001.555.832l3.197-2.132a1 1 0 000-1.664z"/>
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
                </svg>
                <svg v-else class="animate-spin w-5 h-5" fill="none" viewBox="0 0 24 24">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
                <span>{{ isStarting ? 'Memulai...' : 'Mulai Penghitungan' }}</span>
              </button>
            </div>
          </Motion>
        </template>

        <!-- State 2: Input Form (After Start) -->
        <template v-else-if="currentCounting.status === 'IN_PROGRESS'">
          <CountingPrintInfo :po-info="currentPO" :print-info="printInfo" :collapsible="true" />

          <CountingInputForm
            :model-value="countingResult"
            :is-loading="isSaving"
            @update:quantity-good="handleQuantityGoodUpdate"
            @update:quantity-defect="handleQuantityDefectUpdate"
          />

          <CountingCalculation
            v-if="calculations"
            :calculations="calculations"
            :quantity-good="countingResult.quantity_good"
            :quantity-defect="countingResult.quantity_defect"
          />

          <DefectBreakdownForm
            v-if="calculations?.isDefectAboveThreshold"
            v-model="countingResult.defect_breakdown"
            :total-defect="countingResult.quantity_defect"
          />

          <Motion v-bind="entranceAnimations.fadeScale" v-if="calculations?.hasVariance">
            <div class="glass-card p-4">
              <label class="block text-sm font-semibold text-gray-700 mb-2">
                Keterangan Selisih
                <span class="text-red-500">*</span>
              </label>
              <textarea
                v-model="countingResult.variance_reason"
                @input="debouncedSave"
                rows="3"
                class="w-full px-4 py-3 border-2 border-gray-300 rounded-xl focus:border-indigo-500 focus:ring-4 focus:ring-indigo-100 transition-all outline-none resize-none"
                placeholder="Jelaskan alasan terjadinya selisih..."
              ></textarea>
              <p class="text-xs text-gray-500 mt-1">Wajib diisi karena ada selisih dari target</p>
            </div>
          </Motion>

          <!-- Finalize Button -->
          <Motion v-bind="entranceAnimations.fadeScale">
            <div class="flex items-center justify-end gap-3 pt-4">
              <button
                @click="showFinalizeModal = true"
                :disabled="!canFinalize || isSaving"
                class="inline-flex items-center gap-2 px-6 py-3 bg-gradient-to-r from-green-600 to-emerald-600 text-white font-bold rounded-xl hover:shadow-lg transition-all active-scale disabled:opacity-50 disabled:cursor-not-allowed"
              >
                <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/>
                </svg>
                <span>Selesai Penghitungan</span>
              </button>
            </div>
          </Motion>

          <!-- Save Indicator -->
          <div v-if="isSaving" class="fixed bottom-4 right-4 z-50">
            <Motion
              :initial="{ opacity: 0, y: 20 }"
              :animate="{ opacity: 1, y: 0 }"
              :exit="{ opacity: 0, y: 20 }"
            >
              <div class="bg-indigo-600 text-white px-4 py-2 rounded-lg shadow-lg flex items-center gap-2">
                <svg class="animate-spin w-4 h-4" fill="none" viewBox="0 0 24 24">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
                <span class="text-sm font-medium">Menyimpan...</span>
              </div>
            </Motion>
          </div>
        </template>

        <!-- State 3: Completed -->
        <template v-else-if="currentCounting.status === 'COMPLETED'">
          <Motion v-bind="entranceAnimations.fadeScale">
            <div class="glass-card p-6 text-center">
              <div class="w-20 h-20 mx-auto mb-4 rounded-full bg-green-500 flex items-center justify-center">
                <svg class="w-10 h-10 text-white" fill="currentColor" viewBox="0 0 20 20">
                  <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd"/>
                </svg>
              </div>
              <h3 class="text-lg font-bold text-gray-900 mb-2">Penghitungan Sudah Selesai</h3>
              <p class="text-gray-600 mb-6">PO ini sudah diselesaikan dan siap untuk pemotongan</p>
              
              <button
                @click="$router.push({ name: 'counting-queue' })"
                class="inline-flex items-center gap-2 px-6 py-3 bg-gradient-to-r from-indigo-600 to-fuchsia-600 text-white font-bold rounded-xl hover:shadow-lg transition-all active-scale"
              >
                <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7"/>
                </svg>
                <span>Kembali ke Queue</span>
              </button>
            </div>
          </Motion>
        </template>
      </template>

      <!-- Error State -->
      <div v-if="error" class="glass-card p-4 bg-red-50 border-2 border-red-200">
        <div class="flex items-start gap-3">
          <svg class="w-6 h-6 text-red-600 flex-shrink-0" fill="currentColor" viewBox="0 0 20 20">
            <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd"/>
          </svg>
          <div>
            <p class="font-semibold text-red-900">Terjadi Kesalahan</p>
            <p class="text-sm text-red-700">{{ error }}</p>
          </div>
        </div>
      </div>

      <!-- Finalize Modal -->
      <CountingFinalizeModal
        v-model="showFinalizeModal"
        :summary="finalizeSummary"
        :loading="isFinalizing"
        @confirm="handleFinalize"
        @cancel="showFinalizeModal = false"
      />
    </div>
  </AppLayout>
</template>

<script setup>
import { ref, computed, watch, onMounted, onBeforeUnmount } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { storeToRefs } from 'pinia'
import { Motion } from 'motion-v'
import { entranceAnimations } from '@/composables/useMotion'
import { useCountingStore } from '@/stores/counting'
import { useAlertDialog } from '@/composables/useModal'
import AppLayout from '@/components/layout/AppLayout.vue'
import CountingPrintInfo from '@/components/counting/CountingPrintInfo.vue'
import CountingInputForm from '@/components/counting/CountingInputForm.vue'
import CountingCalculation from '@/components/counting/CountingCalculation.vue'
import DefectBreakdownForm from '@/components/counting/DefectBreakdownForm.vue'
import CountingFinalizeModal from '@/components/counting/CountingFinalizeModal.vue'
import LoadingSkeleton from '@/components/common/LoadingSkeleton.vue'

const route = useRoute()
const router = useRouter()
const countingStore = useCountingStore()
const alertDialog = useAlertDialog()

const {
  currentCounting,
  currentPO,
  countingResult,
  calculations,
  canFinalize,
  isLoadingDetail,
  isStarting,
  isSaving,
  isFinalizing,
  error
} = storeToRefs(countingStore)

const showFinalizeModal = ref(false)
const printInfo = computed(() => currentCounting.value?.print_info || null)
let saveTimeout = null

onMounted(async () => {
  const poId = route.params.poId
  
  // Try to load existing counting detail or prepare for new counting
  try {
    // Attempt to fetch counting detail by PO ID
    await countingStore.fetchCountingDetailByPOID(poId)
  } catch (err) {
    // If not found, it means counting hasn't started yet - that's OK
    // Load PO info from queue for display
    await loadPOInfo(poId)
  }
})

onBeforeUnmount(() => {
  if (saveTimeout) {
    clearTimeout(saveTimeout)
  }
  countingStore.clearCurrentCounting()
})

const statusLabel = computed(() => {
  if (!currentCounting.value) return ''
  
  const statusMap = {
    'PENDING': 'Belum Dimulai',
    'IN_PROGRESS': 'Sedang Dihitung',
    'COMPLETED': 'Selesai'
  }
  
  return statusMap[currentCounting.value.status] || currentCounting.value.status
})

const countingDuration = computed(() => {
  if (!currentCounting.value?.started_at) return ''
  
  const start = new Date(currentCounting.value.started_at)
  const now = new Date()
  const diffMinutes = Math.floor((now - start) / 1000 / 60)
  
  if (diffMinutes < 60) {
    return `${diffMinutes} menit`
  }
  
  const hours = Math.floor(diffMinutes / 60)
  const minutes = diffMinutes % 60
  return `${hours} jam ${minutes} menit`
})

const finalizeSummary = computed(() => {
  if (!calculations.value) return {}
  
  return {
    quantityGood: countingResult.value.quantity_good,
    quantityDefect: countingResult.value.quantity_defect,
    totalCounted: calculations.value.totalCounted,
    varianceFromTarget: calculations.value.varianceFromTarget,
    variancePercentage: calculations.value.variancePercentage,
    percentageGood: calculations.value.percentageGood,
    percentageDefect: calculations.value.percentageDefect,
    defectBreakdown: countingResult.value.defect_breakdown,
    varianceReason: countingResult.value.variance_reason
  }
})

const loadPOInfo = async (poId) => {
  // Load from queue to get PO info
  await countingStore.fetchQueue()
  const queueItem = countingStore.queue.find(item => item.po_id == poId)
  
  if (queueItem) {
    countingStore.currentPO = {
      po_number: queueItem.po_number,
      obc_number: queueItem.obc_number,
      target_quantity: queueItem.target_quantity
    }
  }
}

const handleStartCounting = async () => {
  try {
    const poId = route.params.poId
    const poData = {
      po_number: currentPO.value.po_number,
      obc_number: currentPO.value.obc_number,
      target_quantity: currentPO.value.target_quantity
    }
    
    await countingStore.startCounting(poId, poData)
    
    // Haptic feedback
    if ('vibrate' in navigator) {
      navigator.vibrate(10)
    }
    
    alertDialog.success('Penghitungan dimulai!', {
      detail: 'Silakan input hasil penghitungan',
      autoDismiss: true,
      autoDismissDelay: 2000
    })
  } catch (err) {
    alertDialog.error('Gagal memulai penghitungan', {
      detail: error.value || 'Terjadi kesalahan'
    })
  }
}

const handleQuantityGoodUpdate = (value) => {
  countingStore.updateQuantityGood(value)
  debouncedSave()
}

const handleQuantityDefectUpdate = (value) => {
  countingStore.updateQuantityDefect(value)
  debouncedSave()
}

const debouncedSave = () => {
  if (saveTimeout) {
    clearTimeout(saveTimeout)
  }
  
  saveTimeout = setTimeout(async () => {
    try {
      await countingStore.updateResult()
    } catch (err) {
      console.error('Auto-save error:', err)
    }
  }, 1000)
}

const handleFinalize = async () => {
  try {
    await countingStore.finalize()
    
    showFinalizeModal.value = false
    
    // Haptic feedback
    if ('vibrate' in navigator) {
      navigator.vibrate([20, 100, 20])
    }
    
    await alertDialog.success('Penghitungan berhasil diselesaikan!', {
      detail: 'PO siap untuk pemotongan',
      autoDismiss: true,
      autoDismissDelay: 2000
    })
    
    // Redirect to queue after short delay
    setTimeout(() => {
      router.push({ name: 'counting-queue' })
    }, 2000)
  } catch (err) {
    showFinalizeModal.value = false
    
    alertDialog.error('Gagal menyelesaikan penghitungan', {
      detail: error.value || 'Terjadi kesalahan'
    })
  }
}

const handleBack = () => {
  if (currentCounting.value?.status === 'IN_PROGRESS') {
    // Warn user about unsaved changes
    if (confirm('Penghitungan masih berlangsung. Data sudah otomatis tersimpan. Yakin ingin kembali?')) {
      router.push({ name: 'counting-queue' })
    }
  } else {
    router.push({ name: 'counting-queue' })
  }
}

// Watch for defect breakdown changes
watch(() => countingResult.value.defect_breakdown, () => {
  debouncedSave()
}, { deep: true })
</script>
