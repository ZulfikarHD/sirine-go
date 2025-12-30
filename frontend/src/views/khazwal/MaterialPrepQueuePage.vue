<template>
  <AppLayout>
    <!-- Header Section -->
    <Motion v-bind="entranceAnimations.fadeUp">
      <div class="mb-6">
        <h1 class="text-3xl font-bold text-gray-900 mb-2">
          Persiapan Material
        </h1>
        <p class="text-gray-600">
          Queue PO untuk persiapan plat, kertas, dan tinta
        </p>
      </div>
    </Motion>

    <!-- Filter Bar -->
    <Motion 
      :initial="{ opacity: 0, y: 15 }"
      :animate="{ opacity: 1, y: 0 }"
      :transition="{ duration: 0.25, delay: 0.1, ease: 'easeOut' }"
      class="glass-card rounded-2xl p-4 mb-6"
    >
      <div class="flex flex-col sm:flex-row gap-3">
        <!-- Search Input -->
        <div class="flex-1">
          <div class="relative">
            <Search class="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-gray-400" />
            <input
              v-model="filters.search"
              type="text"
              placeholder="Cari PO number, OBC, atau produk..."
              class="w-full pl-10 pr-4 py-2.5 bg-white border border-gray-300 rounded-xl text-sm focus:ring-2 focus:ring-indigo-500 focus:border-transparent transition-all"
              @input="debouncedSearch"
            />
          </div>
        </div>

        <!-- Priority Filter -->
        <div class="sm:w-48">
          <select
            v-model="filters.priority"
            class="w-full px-4 py-2.5 bg-white border border-gray-300 rounded-xl text-sm focus:ring-2 focus:ring-indigo-500 focus:border-transparent transition-all"
            @change="fetchQueue"
          >
            <option value="">Semua Priority</option>
            <option value="URGENT">Urgent</option>
            <option value="NORMAL">Normal</option>
            <option value="LOW">Rendah</option>
          </select>
        </div>

        <!-- Refresh Button -->
        <button
          @click="fetchQueue"
          :disabled="loading"
          class="px-4 py-2.5 bg-white border border-gray-300 rounded-xl text-sm font-semibold text-gray-700 hover:bg-gray-50 active-scale transition-all"
        >
          <RefreshCw :class="{ 'animate-spin': loading }" class="w-4 h-4 inline-block" />
        </button>
      </div>
    </Motion>

    <!-- Loading State -->
    <div v-if="loading && !queueItems.length" class="space-y-4">
      <div v-for="i in 3" :key="i" class="glass-card rounded-2xl p-5 animate-pulse">
        <div class="flex items-start justify-between mb-3">
          <div class="flex-1">
            <div class="h-6 bg-gray-200 rounded w-32 mb-2"></div>
            <div class="h-4 bg-gray-200 rounded w-24"></div>
          </div>
          <div class="h-6 bg-gray-200 rounded w-16"></div>
        </div>
        <div class="h-4 bg-gray-200 rounded w-full mb-3"></div>
        <div class="grid grid-cols-2 gap-3">
          <div class="h-12 bg-gray-200 rounded"></div>
          <div class="h-12 bg-gray-200 rounded"></div>
        </div>
      </div>
    </div>

    <!-- Queue List -->
    <div v-else-if="queueItems.length > 0" class="space-y-4">
      <POQueueCard
        v-for="(item, index) in queueItems"
        :key="item.id"
        :po-item="item"
        :index="index"
        @click="handleCardClick"
      />

      <!-- Pagination -->
      <div v-if="pagination.total_pages > 1" class="flex justify-center gap-2 pt-4">
        <button
          v-for="page in pagination.total_pages"
          :key="page"
          @click="changePage(page)"
          :class="[
            'px-4 py-2 rounded-xl font-semibold text-sm transition-all active-scale',
            page === pagination.page
              ? 'bg-gradient-to-r from-indigo-600 to-fuchsia-600 text-white'
              : 'bg-white border border-gray-300 text-gray-700 hover:bg-gray-50'
          ]"
        >
          {{ page }}
        </button>
      </div>
    </div>

    <!-- Empty State -->
    <Motion
      v-else
      v-bind="entranceAnimations.fadeScale"
      class="glass-card rounded-2xl p-12 text-center"
    >
      <Motion v-bind="iconAnimations.popIn">
        <div class="inline-flex items-center justify-center w-20 h-20 rounded-full bg-gradient-to-br from-indigo-100 to-fuchsia-100 mb-6">
          <Package class="w-10 h-10 text-indigo-600" />
        </div>
      </Motion>

      <h3 class="text-xl font-bold text-gray-900 mb-2">
        {{ emptyStateTitle }}
      </h3>
      <p class="text-gray-500 mb-6 max-w-md mx-auto">
        {{ emptyStateMessage }}
      </p>

      <button
        v-if="hasActiveFilters"
        @click="clearFilters"
        class="px-6 py-3 bg-gradient-to-r from-indigo-600 to-fuchsia-600 text-white font-semibold rounded-xl hover:shadow-lg active-scale"
        style="transition: all 0.15s ease-out"
      >
        <X class="w-4 h-4 inline-block mr-2" />
        Reset Filter
      </button>
    </Motion>

    <!-- Error State -->
    <Motion
      v-if="error"
      v-bind="entranceAnimations.fadeScale"
      class="glass-card rounded-2xl p-8 text-center bg-red-50 border-red-200"
    >
      <AlertTriangle class="w-12 h-12 text-red-600 mx-auto mb-4" />
      <h3 class="text-lg font-bold text-red-900 mb-2">
        Terjadi Kesalahan
      </h3>
      <p class="text-red-700 mb-4">
        {{ error }}
      </p>
      <button
        @click="fetchQueue"
        class="px-6 py-3 bg-red-600 text-white font-semibold rounded-xl hover:bg-red-700 active-scale"
      >
        <RefreshCw class="w-4 h-4 inline-block mr-2" />
        Coba Lagi
      </button>
    </Motion>
  </AppLayout>
</template>

<script setup>
/**
 * Material Prep Queue Page (Sprint 2: Full Implementation)
 * Halaman queue PO untuk persiapan material dengan filtering,
 * search, dan pagination functionality
 */
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { Motion } from 'motion-v'
import { entranceAnimations, iconAnimations } from '@/composables/useMotion'
import { useKhazwalStore } from '@/stores/khazwal'
import AppLayout from '@/components/layout/AppLayout.vue'
import POQueueCard from '@/components/khazwal/POQueueCard.vue'
import { 
  Package, 
  RefreshCw,
  Search,
  X,
  AlertTriangle
} from 'lucide-vue-next'

const router = useRouter()
const khazwalStore = useKhazwalStore()

// State from store
const loading = computed(() => khazwalStore.queueLoading)
const error = computed(() => khazwalStore.queueError)
const queueItems = computed(() => khazwalStore.queue)
const pagination = computed(() => khazwalStore.queuePagination)

// Filter state
const filters = ref({
  search: '',
  priority: '',
  page: 1,
  per_page: 20
})

/**
 * Fetch queue data dari store dengan current filters
 */
const fetchQueue = async () => {
  try {
    await khazwalStore.getMaterialPrepQueue({
      ...filters.value,
      page: filters.value.page || 1
    })
  } catch (err) {
    console.error('Error fetching queue:', err)
  }
}

/**
 * Debounced search untuk avoid excessive API calls
 */
let searchTimeout
const debouncedSearch = () => {
  clearTimeout(searchTimeout)
  searchTimeout = setTimeout(() => {
    filters.value.page = 1
    fetchQueue()
  }, 500)
}

/**
 * Handle page change untuk pagination
 */
const changePage = (page) => {
  filters.value.page = page
  fetchQueue()
  
  // Haptic feedback
  if ('vibrate' in navigator) {
    navigator.vibrate(10)
  }
}

/**
 * Handle card click untuk navigate ke detail page
 */
const handleCardClick = (poItem) => {
  router.push(`/khazwal/material-prep/${poItem.id}`)
}

/**
 * Clear all active filters
 */
const clearFilters = () => {
  filters.value = {
    search: '',
    priority: '',
    page: 1,
    per_page: 20
  }
  fetchQueue()
}

/**
 * Check if any filters are active
 */
const hasActiveFilters = computed(() => {
  return filters.value.search !== '' || filters.value.priority !== ''
})

/**
 * Empty state title berdasarkan filter status
 */
const emptyStateTitle = computed(() => {
  if (hasActiveFilters.value) {
    return 'Tidak Ada Hasil'
  }
  return 'Belum Ada PO'
})

/**
 * Empty state message berdasarkan filter status
 */
const emptyStateMessage = computed(() => {
  if (hasActiveFilters.value) {
    return 'Tidak ada PO yang sesuai dengan filter Anda. Coba ubah kriteria pencarian.'
  }
  return 'Belum ada Production Order yang menunggu persiapan material.'
})

// Fetch queue on mount
onMounted(() => {
  fetchQueue()
})
</script>
