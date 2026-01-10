<template>
  <AppLayout>
    <div class="max-w-4xl mx-auto px-4 py-6 space-y-6">
      <!-- Header -->
      <Motion v-bind="entranceAnimations.fadeUp">
        <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
          <div>
            <h1 class="text-2xl sm:text-3xl font-bold bg-gradient-to-r from-indigo-600 to-fuchsia-600 bg-clip-text text-transparent">
              Antrian Verifikasi
            </h1>
            <p class="text-gray-600 mt-1">Label-label yang siap untuk diverifikasi</p>
          </div>

          <!-- Refresh Button -->
          <button
            @click="refreshQueue"
            :disabled="isLoadingQueue"
            class="inline-flex items-center gap-2 px-4 py-2 bg-white border-2 border-gray-300 rounded-xl font-semibold text-gray-700 hover:border-indigo-500 hover:text-indigo-600 transition-all active-scale"
            :class="{ 'opacity-50 cursor-not-allowed': isLoadingQueue }"
          >
            <svg 
              class="w-5 h-5"
              :class="{ 'animate-spin': isLoadingQueue }"
              fill="none" 
              stroke="currentColor" 
              viewBox="0 0 24 24"
            >
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"/>
            </svg>
            <span class="hidden sm:inline">Refresh</span>
          </button>
        </div>
      </Motion>

      <!-- Stats Summary -->
      <Motion v-bind="entranceAnimations.fadeScale">
        <div class="grid grid-cols-1 sm:grid-cols-3 gap-4">
          <div class="glass-card p-4">
            <div class="flex items-center gap-3">
              <div class="w-12 h-12 rounded-full bg-gradient-to-br from-indigo-500 to-fuchsia-500 flex items-center justify-center">
                <svg class="w-6 h-6 text-white" fill="currentColor" viewBox="0 0 20 20">
                  <path fill-rule="evenodd" d="M6 2a1 1 0 00-1 1v1H4a2 2 0 00-2 2v10a2 2 0 002 2h12a2 2 0 002-2V6a2 2 0 00-2-2h-1V3a1 1 0 10-2 0v1H7V3a1 1 0 00-1-1zm0 5a1 1 0 000 2h8a1 1 0 100-2H6z" clip-rule="evenodd"/>
                </svg>
              </div>
              <div>
                <p class="text-xs text-gray-500">Total PO</p>
                <p class="text-2xl font-bold text-gray-900">{{ queueMeta.total_po }}</p>
              </div>
            </div>
          </div>

          <div class="glass-card p-4">
            <div class="flex items-center gap-3">
              <div class="w-12 h-12 rounded-full bg-gradient-to-br from-emerald-500 to-cyan-500 flex items-center justify-center">
                <svg class="w-6 h-6 text-white" fill="currentColor" viewBox="0 0 20 20">
                  <path d="M7 3a1 1 0 000 2h6a1 1 0 100-2H7zM4 7a1 1 0 011-1h10a1 1 0 110 2H5a1 1 0 01-1-1zM2 11a2 2 0 012-2h12a2 2 0 012 2v4a2 2 0 01-2 2H4a2 2 0 01-2-2v-4z"/>
                </svg>
              </div>
              <div>
                <p class="text-xs text-gray-500">Total Label</p>
                <p class="text-2xl font-bold text-emerald-600">{{ queueMeta.total_labels }}</p>
              </div>
            </div>
          </div>

          <div class="glass-card p-4">
            <div class="flex items-center gap-3">
              <div class="w-12 h-12 rounded-full bg-gradient-to-br from-amber-500 to-orange-500 flex items-center justify-center">
                <svg class="w-6 h-6 text-white" fill="currentColor" viewBox="0 0 20 20">
                  <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm1-12a1 1 0 10-2 0v4a1 1 0 00.293.707l2.828 2.829a1 1 0 101.415-1.415L11 9.586V6z" clip-rule="evenodd"/>
                </svg>
              </div>
              <div>
                <p class="text-xs text-gray-500">Pending</p>
                <p class="text-2xl font-bold text-amber-600">{{ queueMeta.pending_count }}</p>
              </div>
            </div>
          </div>
        </div>
      </Motion>

      <!-- Filters -->
      <Motion v-bind="entranceAnimations.fadeScale">
        <div class="glass-card p-4">
          <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
            <!-- Priority Filter -->
            <div>
              <label class="block text-sm font-semibold text-gray-700 mb-2">
                <svg class="w-4 h-4 inline mr-1" fill="currentColor" viewBox="0 0 20 20">
                  <path d="M3 3a1 1 0 000 2h11a1 1 0 100-2H3zM3 7a1 1 0 000 2h5a1 1 0 000-2H3zM3 11a1 1 0 100 2h4a1 1 0 100-2H3z"/>
                </svg>
                Prioritas
              </label>
              <select
                v-model="filters.priority"
                @change="applyFilters"
                class="w-full px-4 py-2 border-2 border-gray-300 rounded-xl focus:border-indigo-500 focus:ring-4 focus:ring-indigo-100 transition-all"
              >
                <option value="">Semua Prioritas</option>
                <option value="URGENT">Urgent</option>
                <option value="HIGH">High</option>
                <option value="NORMAL">Normal</option>
                <option value="LOW">Low</option>
              </select>
            </div>

            <!-- PO Search -->
            <div>
              <label class="block text-sm font-semibold text-gray-700 mb-2">
                <svg class="w-4 h-4 inline mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"/>
                </svg>
                Cari PO
              </label>
              <input
                type="text"
                v-model="filters.search"
                @input="applyFilters"
                placeholder="Nomor PO atau OBC"
                class="w-full px-4 py-2 border-2 border-gray-300 rounded-xl focus:border-indigo-500 focus:ring-4 focus:ring-indigo-100 transition-all"
              />
            </div>
          </div>

          <!-- Reset Filters -->
          <button
            v-if="hasActiveFilters"
            @click="resetFilters"
            class="mt-3 text-sm text-indigo-600 hover:text-indigo-700 font-semibold flex items-center gap-1"
          >
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
            </svg>
            Reset Filter
          </button>
        </div>
      </Motion>

      <!-- Loading State -->
      <div v-if="isLoadingQueue" class="space-y-4">
        <LoadingSkeleton v-for="i in 3" :key="i" class="h-48" />
      </div>

      <!-- Empty State -->
      <Motion v-else-if="queueData.length === 0" v-bind="entranceAnimations.fadeScale">
        <div class="glass-card p-12 text-center">
          <div class="w-20 h-20 mx-auto mb-4 rounded-full bg-gradient-to-br from-gray-100 to-gray-200 flex items-center justify-center">
            <svg class="w-10 h-10 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/>
            </svg>
          </div>
          <h3 class="text-xl font-bold text-gray-900 mb-2">
            {{ hasActiveFilters ? 'Tidak Ada Hasil' : 'Belum Ada PO Siap Verifikasi' }}
          </h3>
          <p class="text-gray-600 mb-4">
            {{ hasActiveFilters ? 'Coba ubah filter pencarian Anda' : 'Label verifikasi akan muncul setelah proses pemotongan selesai' }}
          </p>
          <button
            v-if="hasActiveFilters"
            @click="resetFilters"
            class="px-6 py-2 bg-gradient-to-r from-indigo-600 to-fuchsia-600 text-white font-semibold rounded-xl hover:from-indigo-700 hover:to-fuchsia-700 transition-all active-scale"
          >
            Reset Filter
          </button>
        </div>
      </Motion>

      <!-- Queue List - Grouped by PO -->
      <div v-else class="space-y-6">
        <div
          v-for="(po, poIndex) in queueData"
          :key="po.po_id"
        >
          <!-- PO Header -->
          <Motion
            :initial="{ opacity: 0, y: 15 }"
            :animate="{ opacity: 1, y: 0 }"
            :transition="{ duration: 0.25, delay: poIndex * 0.05, ease: 'easeOut' }"
          >
            <div class="glass-card p-4 mb-3">
              <div class="flex items-center justify-between">
                <div class="flex items-center gap-3">
                  <div class="w-10 h-10 rounded-full bg-gradient-to-br from-indigo-500 to-fuchsia-500 flex items-center justify-center">
                    <svg class="w-5 h-5 text-white" fill="currentColor" viewBox="0 0 20 20">
                      <path d="M9 2a1 1 0 000 2h2a1 1 0 100-2H9z"/>
                      <path fill-rule="evenodd" d="M4 5a2 2 0 012-2 3 3 0 003 3h2a3 3 0 003-3 2 2 0 012 2v11a2 2 0 01-2 2H6a2 2 0 01-2-2V5zm3 4a1 1 0 000 2h.01a1 1 0 100-2H7zm3 0a1 1 0 000 2h3a1 1 0 100-2h-3zm-3 4a1 1 0 100 2h.01a1 1 0 100-2H7zm3 0a1 1 0 100 2h3a1 1 0 100-2h-3z" clip-rule="evenodd"/>
                    </svg>
                  </div>
                  <div>
                    <h3 class="text-lg font-bold text-gray-900">PO {{ po.po_number }}</h3>
                    <p class="text-sm text-gray-500">{{ po.obc_number }}</p>
                  </div>
                </div>
                <div class="text-right">
                  <p class="text-sm font-semibold text-gray-700">{{ po.labels.length }} Label</p>
                  <p class="text-xs text-gray-500">{{ formatNumber(po.total_quantity) }} lembar</p>
                </div>
              </div>
            </div>
          </Motion>

          <!-- Labels List -->
          <div class="grid grid-cols-1 md:grid-cols-2 gap-4 ml-4">
            <Motion
              v-for="(label, labelIndex) in po.labels"
              :key="label.label_id"
              :initial="{ opacity: 0, y: 15 }"
              :animate="{ opacity: 1, y: 0 }"
              :transition="{ duration: 0.25, delay: (poIndex * 0.05) + (labelIndex * 0.03), ease: 'easeOut' }"
            >
              <VerificationLabelCard
                :label-id="label.label_id"
                :label-number="label.label_number"
                :total-labels="po.labels.length"
                :target-quantity="label.target_quantity"
                :sisiran="label.sisiran"
                :qc-status="label.qc_status"
                :po-number="po.po_number"
                :obc-number="po.obc_number"
                :cutting-info="label.cutting_info"
                :waste-percentage="label.waste_percentage"
                @start-verification="handleStartVerification(label)"
                @continue-verification="handleContinueVerification(label)"
              />
            </Motion>
          </div>
        </div>
      </div>
    </div>
  </AppLayout>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { Motion } from 'motion-v'
import { entranceAnimations } from '@/composables/useMotion'
import { useVerifikasiApi } from '@/composables/useVerifikasiApi'
import { useHaptic } from '@/composables/useHaptic'
import AppLayout from '@/components/layout/AppLayout.vue'
import VerificationLabelCard from '@/components/verifikasi/VerificationLabelCard.vue'
import { LoadingSkeleton } from '@/components/common'

const router = useRouter()
const verifikasiApi = useVerifikasiApi()
const haptic = useHaptic()

// State
const isLoadingQueue = ref(false)
const queueData = ref([])
const queueMeta = ref({
  total_po: 0,
  total_labels: 0,
  pending_count: 0,
})

// Filters
const filters = ref({
  priority: '',
  search: '',
})

// Computed
const hasActiveFilters = computed(() => {
  return filters.value.priority !== '' || filters.value.search !== ''
})

/**
 * Format number dengan thousand separator
 */
const formatNumber = (num) => {
  return new Intl.NumberFormat('id-ID').format(num)
}

/**
 * Fetch verification queue dari API
 */
const fetchQueue = async () => {
  isLoadingQueue.value = true
  
  try {
    // TODO: Replace with actual API when backend is ready
    // const response = await verifikasiApi.getVerificationQueue(filters.value)
    
    // Mock data for now
    const response = {
      data: [],
      meta: {
        total_po: 0,
        total_labels: 0,
        pending_count: 0,
      }
    }
    
    queueData.value = response.data || []
    queueMeta.value = response.meta || {
      total_po: 0,
      total_labels: 0,
      pending_count: 0,
    }
  } catch (error) {
    console.error('Failed to fetch verification queue:', error)
    // For now, show empty state
    queueData.value = []
  } finally {
    isLoadingQueue.value = false
  }
}

/**
 * Apply filters dan refresh queue
 */
const applyFilters = () => {
  haptic.light()
  fetchQueue()
}

/**
 * Reset semua filters
 */
const resetFilters = () => {
  haptic.light()
  filters.value = {
    priority: '',
    search: '',
  }
  fetchQueue()
}

/**
 * Refresh queue manual
 */
const refreshQueue = () => {
  haptic.light()
  fetchQueue()
}

/**
 * Handle start verification action
 */
const handleStartVerification = (label) => {
  haptic.medium()
  // TODO: Navigate to verification detail page
  console.log('Start verification for label:', label.label_id)
  // router.push({ name: 'verification-detail', params: { labelId: label.label_id } })
}

/**
 * Handle continue verification action
 */
const handleContinueVerification = (label) => {
  haptic.medium()
  // TODO: Navigate to verification detail page
  console.log('Continue verification for label:', label.label_id)
  // router.push({ name: 'verification-detail', params: { labelId: label.label_id } })
}

// Lifecycle
onMounted(() => {
  fetchQueue()
})
</script>
