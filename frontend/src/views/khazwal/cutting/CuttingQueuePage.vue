<template>
  <AppLayout>
    <div class="max-w-4xl mx-auto px-4 py-6 space-y-6">
      <!-- Header -->
      <Motion v-bind="entranceAnimations.fadeUp">
        <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
          <div>
            <h1 class="text-2xl sm:text-3xl font-bold bg-gradient-to-r from-indigo-600 to-fuchsia-600 bg-clip-text text-transparent">
              Pemotongan
            </h1>
            <p class="text-gray-600 mt-1">Daftar PO menunggu pemotongan</p>
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
                  <path d="M9 2a1 1 0 000 2h2a1 1 0 100-2H9z"/>
                  <path fill-rule="evenodd" d="M4 5a2 2 0 012-2 3 3 0 003 3h2a3 3 0 003-3 2 2 0 012 2v11a2 2 0 01-2 2H6a2 2 0 01-2-2V5zm3 4a1 1 0 000 2h.01a1 1 0 100-2H7zm3 0a1 1 0 000 2h3a1 1 0 100-2h-3zm-3 4a1 1 0 100 2h.01a1 1 0 100-2H7zm3 0a1 1 0 100 2h3a1 1 0 100-2h-3z" clip-rule="evenodd"/>
                </svg>
              </div>
              <div>
                <p class="text-xs text-gray-500">Total Queue</p>
                <p class="text-2xl font-bold text-gray-900">{{ queueMeta.total }}</p>
              </div>
            </div>
          </div>

          <div class="glass-card p-4">
            <div class="flex items-center gap-3">
              <div class="w-12 h-12 rounded-full bg-gradient-to-br from-red-500 to-orange-500 flex items-center justify-center">
                <svg class="w-6 h-6 text-white" fill="currentColor" viewBox="0 0 20 20">
                  <path fill-rule="evenodd" d="M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.213 2.98-1.742 2.98H4.42c-1.53 0-2.493-1.646-1.743-2.98l5.58-9.92zM11 13a1 1 0 11-2 0 1 1 0 012 0zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z" clip-rule="evenodd"/>
                </svg>
              </div>
              <div>
                <p class="text-xs text-gray-500">Urgent</p>
                <p class="text-2xl font-bold text-red-600">{{ queueMeta.urgent_count }}</p>
              </div>
            </div>
          </div>

          <div class="glass-card p-4">
            <div class="flex items-center gap-3">
              <div class="w-12 h-12 rounded-full bg-gradient-to-br from-emerald-500 to-cyan-500 flex items-center justify-center">
                <svg class="w-6 h-6 text-white" fill="currentColor" viewBox="0 0 20 20">
                  <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm1-12a1 1 0 10-2 0v4a1 1 0 00.293.707l2.828 2.829a1 1 0 101.415-1.415L11 9.586V6z" clip-rule="evenodd"/>
                </svg>
              </div>
              <div>
                <p class="text-xs text-gray-500">Normal</p>
                <p class="text-2xl font-bold text-emerald-600">{{ queueMeta.normal_count }}</p>
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

            <!-- Sort Filter -->
            <div>
              <label class="block text-sm font-semibold text-gray-700 mb-2">
                <svg class="w-4 h-4 inline mr-1" fill="currentColor" viewBox="0 0 20 20">
                  <path d="M3 3a1 1 0 000 2h11a1 1 0 100-2H3zM3 7a1 1 0 000 2h7a1 1 0 100-2H3zM3 11a1 1 0 100 2h4a1 1 0 100-2H3z"/>
                </svg>
                Urutan
              </label>
              <select
                v-model="filters.sort_by"
                @change="applyFilters"
                class="w-full px-4 py-2 border-2 border-gray-300 rounded-xl focus:border-indigo-500 focus:ring-4 focus:ring-indigo-100 transition-all"
              >
                <option value="">Default (Prioritas + FIFO)</option>
                <option value="date">Tanggal Selesai Counting</option>
              </select>
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
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14.121 14.121L19 19m-7-7l7-7m-7 7l-2.879 2.879M12 12L9.121 9.121m0 5.758a3 3 0 10-4.243 4.243 3 3 0 004.243-4.243zm0-5.758a3 3 0 10-4.243-4.243 3 3 0 004.243 4.243z"/>
            </svg>
          </div>
          <h3 class="text-xl font-bold text-gray-900 mb-2">
            {{ hasActiveFilters ? 'Tidak Ada Hasil' : 'Tidak Ada PO' }}
          </h3>
          <p class="text-gray-600 mb-4">
            {{ hasActiveFilters ? 'Coba ubah filter pencarian Anda' : 'Semua PO sudah selesai dipotong' }}
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

      <!-- Queue List -->
      <div v-else class="space-y-4">
        <Motion
          v-for="(item, index) in queueData"
          :key="item.po_id"
          :initial="{ opacity: 0, y: 15 }"
          :animate="{ opacity: 1, y: 0 }"
          :transition="{ duration: 0.25, delay: index * 0.05, ease: 'easeOut' }"
        >
          <CuttingQueueCard
            :po-id="item.po_id"
            :po-number="item.po_number"
            :obc-number="item.obc_number"
            :priority="item.priority"
            :input-lembar-besar="item.input_lembar_besar"
            :estimated-output="item.estimated_output"
            :counting-completed-at="item.counting_completed_at"
            :waiting-minutes="item.waiting_minutes"
            :is-overdue="item.is_overdue"
            @start-cutting="handleStartCutting(item)"
          />
        </Motion>
      </div>
    </div>
  </AppLayout>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { Motion } from 'motion-v'
import { entranceAnimations } from '@/composables/useMotion'
import { useCuttingApi } from '@/composables/useCuttingApi'
import { useHaptic } from '@/composables/useHaptic'
import AppLayout from '@/components/layout/AppLayout.vue'
import CuttingQueueCard from '@/components/cutting/CuttingQueueCard.vue'
import { LoadingSkeleton } from '@/components/common'

const router = useRouter()
const cuttingApi = useCuttingApi()
const haptic = useHaptic()

// State
const isLoadingQueue = ref(false)
const queueData = ref([])
const queueMeta = ref({
  total: 0,
  urgent_count: 0,
  normal_count: 0,
})

// Filters
const filters = ref({
  priority: '',
  sort_by: '',
})

// Computed
const hasActiveFilters = computed(() => {
  return filters.value.priority !== '' || filters.value.sort_by !== ''
})

/**
 * Fetch cutting queue dari API
 */
const fetchQueue = async () => {
  isLoadingQueue.value = true
  
  try {
    const response = await cuttingApi.getCuttingQueue(filters.value)
    queueData.value = response.data || []
    queueMeta.value = response.meta || {
      total: 0,
      urgent_count: 0,
      normal_count: 0,
    }
  } catch (error) {
    console.error('Failed to fetch cutting queue:', error)
    // TODO: Show error toast notification
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
    sort_by: '',
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
 * Handle start cutting action
 */
const handleStartCutting = (item) => {
  haptic.medium()
  router.push({
    name: 'cutting-start',
    params: { poId: item.po_id }
  })
}

// Lifecycle
onMounted(() => {
  fetchQueue()
})
</script>
