<template>
  <AppLayout>
    <div class="min-h-screen pb-20">
      <!-- Header Section -->
      <Motion v-bind="entranceAnimations.fadeUp" class="mb-6">
        <div class="flex items-center justify-between mb-4">
          <div>
            <h1 class="text-2xl sm:text-3xl font-bold text-gray-900">
              Riwayat Persiapan
            </h1>
            <p class="text-sm text-gray-500 mt-1">
              {{ totalItems }} persiapan selesai
            </p>
          </div>
          
          <!-- Refresh Button -->
          <button
            @click="refreshHistory"
            class="p-2.5 rounded-xl bg-white/80 border border-gray-200/50 
                   hover:bg-gray-50 active-scale shadow-sm"
            :disabled="loading"
          >
            <RefreshCw 
              class="w-5 h-5 text-gray-600" 
              :class="{ 'animate-spin': loading }"
            />
          </button>
        </div>

        <!-- Filter Bar -->
        <div class="space-y-3">
          <!-- Search Input -->
          <div class="relative">
            <Search class="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-gray-400" />
            <input
              v-model="searchQuery"
              type="text"
              placeholder="Cari PO, OBC, atau produk..."
              class="w-full pl-10 pr-4 py-3 rounded-xl bg-white/80 border border-gray-200/50 
                     text-gray-900 placeholder-gray-400 focus:outline-none focus:ring-2 
                     focus:ring-indigo-500/20 focus:border-indigo-400 shadow-sm"
              @input="handleSearch"
            />
          </div>

          <!-- Date Range Filter -->
          <div class="flex gap-3">
            <div class="flex-1 relative">
              <Calendar class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400" />
              <input
                v-model="dateFrom"
                type="date"
                class="w-full pl-10 pr-3 py-2.5 rounded-xl bg-white/80 border border-gray-200/50 
                       text-gray-900 text-sm focus:outline-none focus:ring-2 
                       focus:ring-indigo-500/20 focus:border-indigo-400 shadow-sm"
                @change="applyDateFilter"
              />
            </div>
            <div class="flex items-center text-gray-400 text-sm">s/d</div>
            <div class="flex-1 relative">
              <Calendar class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400" />
              <input
                v-model="dateTo"
                type="date"
                class="w-full pl-10 pr-3 py-2.5 rounded-xl bg-white/80 border border-gray-200/50 
                       text-gray-900 text-sm focus:outline-none focus:ring-2 
                       focus:ring-indigo-500/20 focus:border-indigo-400 shadow-sm"
                @change="applyDateFilter"
              />
            </div>
            
            <!-- Clear Filter Button -->
            <button
              v-if="hasActiveFilter"
              @click="clearFilters"
              class="p-2.5 rounded-xl bg-red-50 border border-red-200/50 
                     hover:bg-red-100 active-scale"
            >
              <X class="w-5 h-5 text-red-600" />
            </button>
          </div>
        </div>
      </Motion>

      <!-- History Grid -->
      <div v-if="!loading && historyItems && historyItems.length > 0" class="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
        <PrepHistoryCard
          v-for="(item, index) in historyItems"
          :key="item.prep_id"
          :item="item"
          :index="index"
          @click="handleCardClick"
        />
      </div>

      <!-- Empty State -->
      <Motion
        v-if="!loading && historyItems && historyItems.length === 0"
        v-bind="entranceAnimations.fadeUp"
        class="flex flex-col items-center justify-center py-20"
      >
        <div class="w-20 h-20 rounded-2xl bg-gradient-to-br from-indigo-100 to-fuchsia-100 
                    flex items-center justify-center mb-4">
          <History class="w-10 h-10 text-indigo-500" />
        </div>
        <h3 class="text-lg font-semibold text-gray-900 mb-2">
          Belum Ada Riwayat
        </h3>
        <p class="text-gray-500 text-center max-w-xs">
          {{ emptyMessage }}
        </p>
      </Motion>

      <!-- Loading Skeleton -->
      <div v-if="loading" class="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
        <div
          v-for="i in 6"
          :key="i"
          class="glass-card rounded-2xl p-5 animate-pulse"
        >
          <div class="flex items-start justify-between mb-3">
            <div class="space-y-2 flex-1">
              <div class="h-5 bg-gray-200 rounded w-32" />
              <div class="h-4 bg-gray-200 rounded w-24" />
            </div>
            <div class="h-6 w-16 bg-gray-200 rounded-lg" />
          </div>
          <div class="h-4 bg-gray-200 rounded w-full mb-3" />
          <div class="grid grid-cols-2 gap-3 mb-3">
            <div class="h-12 bg-gray-200 rounded-lg" />
            <div class="h-12 bg-gray-200 rounded-lg" />
          </div>
          <div class="h-10 bg-gray-200 rounded-lg" />
        </div>
      </div>

      <!-- Pagination -->
      <Motion
        v-if="!loading && totalPages > 1"
        v-bind="entranceAnimations.fadeUp"
        class="flex items-center justify-center gap-2 mt-8"
      >
        <button
          @click="goToPage(currentPage - 1)"
          :disabled="currentPage <= 1"
          class="p-2 rounded-lg hover:bg-gray-100 disabled:opacity-40 
                 disabled:cursor-not-allowed active-scale"
        >
          <ChevronLeft class="w-5 h-5 text-gray-600" />
        </button>

        <div class="flex items-center gap-1">
          <button
            v-for="page in visiblePages"
            :key="page"
            @click="goToPage(page)"
            class="min-w-[36px] h-9 px-3 rounded-lg text-sm font-medium active-scale"
            :class="currentPage === page 
              ? 'bg-gradient-to-r from-indigo-500 to-fuchsia-500 text-white' 
              : 'hover:bg-gray-100 text-gray-700'"
          >
            {{ page }}
          </button>
        </div>

        <button
          @click="goToPage(currentPage + 1)"
          :disabled="currentPage >= totalPages"
          class="p-2 rounded-lg hover:bg-gray-100 disabled:opacity-40 
                 disabled:cursor-not-allowed active-scale"
        >
          <ChevronRight class="w-5 h-5 text-gray-600" />
        </button>
      </Motion>

      <!-- Detail Modal -->
      <BaseModal
        v-model="showDetailModal"
        :title="`Detail PO #${selectedItem?.po_number || ''}`"
        :subtitle="selectedItem?.obc_number"
        title-gradient
        size="md"
        :show-footer="false"
        scrollable
      >
        <div v-if="selectedItem" class="space-y-4">
          <!-- Product Info -->
          <div class="glass-card rounded-xl p-4">
            <h4 class="text-sm font-semibold text-gray-700 mb-2">Informasi Produk</h4>
            <p class="text-gray-900 font-medium mb-2">{{ selectedItem.product_name }}</p>
            <div class="flex items-center gap-4 text-sm">
              <div>
                <p class="text-gray-500">Quantity</p>
                <p class="font-semibold">{{ selectedItem.quantity.toLocaleString() }}</p>
              </div>
              <div>
                <p class="text-gray-500">Prioritas</p>
                <PriorityBadge :priority="selectedItem.priority" size="sm" />
              </div>
            </div>
          </div>

          <!-- Completion Info -->
          <div class="glass-card rounded-xl p-4">
            <h4 class="text-sm font-semibold text-gray-700 mb-2">Info Penyelesaian</h4>
            <div class="grid grid-cols-2 gap-4 text-sm">
              <div>
                <p class="text-gray-500">Selesai Pada</p>
                <p class="font-semibold">{{ formatDateTime(selectedItem.completed_at) }}</p>
              </div>
              <div>
                <p class="text-gray-500">Durasi</p>
                <p class="font-semibold">{{ formatDuration(selectedItem.duration_minutes) }}</p>
              </div>
              <div>
                <p class="text-gray-500">Dikerjakan Oleh</p>
                <p class="font-semibold">{{ selectedItem.prepared_by_name }}</p>
              </div>
              <div>
                <p class="text-gray-500">Foto Material</p>
                <p class="font-semibold">{{ selectedItem.photos_count }} foto</p>
              </div>
            </div>
          </div>
        </div>
      </BaseModal>
    </div>
  </AppLayout>
</template>

<script setup>
/**
 * MaterialPrepHistoryPage - Halaman riwayat material preparation
 * yang menampilkan list history dengan filter date range
 */
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { Motion } from 'motion-v'
import { entranceAnimations } from '@/composables/useMotion'
import { useKhazwalApi } from '@/composables/useKhazwalApi'
import { useAlertDialog } from '@/composables/useModal'
import AppLayout from '@/components/layout/AppLayout.vue'
import BaseModal from '@/components/common/BaseModal.vue'
import PrepHistoryCard from '@/components/khazwal/PrepHistoryCard.vue'
import PriorityBadge from '@/components/common/PriorityBadge.vue'
import { 
  Search, 
  RefreshCw, 
  ChevronLeft,
  ChevronRight,
  Calendar,
  History,
  X
} from 'lucide-vue-next'

const khazwalApi = useKhazwalApi()
const alertDialog = useAlertDialog()

// State
const loading = ref(false)
const historyItems = ref([])
const totalItems = ref(0)
const currentPage = ref(1)
const totalPages = ref(0)
const perPage = ref(12)

// Filters
const searchQuery = ref('')
const dateFrom = ref('')
const dateTo = ref('')
let searchTimeout = null

// Detail Modal
const showDetailModal = ref(false)
const selectedItem = ref(null)

/**
 * Computed untuk check active filter
 */
const hasActiveFilter = computed(() => {
  return searchQuery.value || dateFrom.value || dateTo.value
})

/**
 * Computed untuk visible pagination pages
 */
const visiblePages = computed(() => {
  const pages = []
  const total = totalPages.value
  const current = currentPage.value
  
  if (total <= 5) {
    for (let i = 1; i <= total; i++) {
      pages.push(i)
    }
  } else {
    if (current <= 3) {
      pages.push(1, 2, 3, 4, 5)
    } else if (current >= total - 2) {
      pages.push(total - 4, total - 3, total - 2, total - 1, total)
    } else {
      pages.push(current - 2, current - 1, current, current + 1, current + 2)
    }
  }
  
  return pages
})

/**
 * Computed untuk empty state message
 */
const emptyMessage = computed(() => {
  if (hasActiveFilter.value) {
    return 'Tidak ada riwayat yang sesuai dengan filter. Coba ubah kriteria pencarian.'
  }
  return 'Belum ada riwayat persiapan material. Riwayat akan muncul setelah proses persiapan selesai.'
})

/**
 * Fetch history data dari API
 */
const fetchHistory = async () => {
  loading.value = true
  try {
    const response = await khazwalApi.getHistory({
      search: searchQuery.value,
      date_from: dateFrom.value,
      date_to: dateTo.value,
      page: currentPage.value,
      per_page: perPage.value
    })

    if (response.success && response.data) {
      historyItems.value = response.data.items || []
      totalItems.value = response.data.total || 0
      totalPages.value = response.data.total_pages || 0
    } else {
      // Fallback jika response tidak sesuai expected structure
      historyItems.value = []
      totalItems.value = 0
      totalPages.value = 0
    }
  } catch (error) {
    console.error('Error fetching history:', error)
    
    // Ensure historyItems tetap array meskipun error
    historyItems.value = []
    totalItems.value = 0
    totalPages.value = 0
    
    alertDialog.error('Gagal memuat riwayat', {
      detail: error.response?.data?.message || 'Silakan coba lagi'
    })
  } finally {
    loading.value = false
  }
}

/**
 * Refresh history data
 */
const refreshHistory = () => {
  fetchHistory()
  
  if ('vibrate' in navigator) {
    navigator.vibrate(10)
  }
}

/**
 * Handle search input dengan debounce
 */
const handleSearch = () => {
  if (searchTimeout) clearTimeout(searchTimeout)
  searchTimeout = setTimeout(() => {
    currentPage.value = 1
    fetchHistory()
  }, 300)
}

/**
 * Apply date filter
 */
const applyDateFilter = () => {
  currentPage.value = 1
  fetchHistory()
}

/**
 * Clear all filters
 */
const clearFilters = () => {
  searchQuery.value = ''
  dateFrom.value = ''
  dateTo.value = ''
  currentPage.value = 1
  fetchHistory()
  
  if ('vibrate' in navigator) {
    navigator.vibrate(10)
  }
}

/**
 * Go to specific page
 */
const goToPage = (page) => {
  if (page < 1 || page > totalPages.value) return
  currentPage.value = page
  fetchHistory()
  
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

/**
 * Handle card click
 */
const handleCardClick = (item) => {
  selectedItem.value = item
  showDetailModal.value = true
}

/**
 * Format datetime untuk display
 */
const formatDateTime = (dateStr) => {
  try {
    const date = new Date(dateStr)
    return date.toLocaleDateString('id-ID', {
      day: 'numeric',
      month: 'long',
      year: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    })
  } catch (e) {
    return dateStr
  }
}

/**
 * Format duration
 */
const formatDuration = (minutes) => {
  if (!minutes) return '-'
  
  if (minutes < 60) {
    return `${minutes} menit`
  }
  
  const hours = Math.floor(minutes / 60)
  const mins = minutes % 60
  
  if (mins === 0) {
    return `${hours} jam`
  }
  
  return `${hours} jam ${mins} menit`
}

onMounted(() => {
  fetchHistory()
})

onUnmounted(() => {
  if (searchTimeout) clearTimeout(searchTimeout)
})
</script>
