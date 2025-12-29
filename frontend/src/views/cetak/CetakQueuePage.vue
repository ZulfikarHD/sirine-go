<template>
  <AppLayout>
    <div class="min-h-screen pb-20">
      <!-- Header Section -->
      <Motion v-bind="entranceAnimations.fadeUp" class="mb-6">
        <div class="flex items-center justify-between mb-4">
          <div>
            <h1 class="text-2xl sm:text-3xl font-bold text-gray-900">
              Antrian Cetak
            </h1>
            <p class="text-sm text-gray-500 mt-1">
              {{ totalItems }} PO siap untuk dicetak
            </p>
          </div>
          
          <!-- Refresh Button -->
          <button
            @click="refreshQueue"
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

        <!-- Search & Filter Bar -->
        <div class="flex gap-3">
          <!-- Search Input -->
          <div class="flex-1 relative">
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

          <!-- Priority Filter -->
          <div class="relative">
            <button
              @click="showPriorityFilter = !showPriorityFilter"
              class="flex items-center gap-2 px-4 py-3 rounded-xl bg-white/80 border 
                     border-gray-200/50 hover:bg-gray-50 active-scale shadow-sm"
              :class="{ 'ring-2 ring-indigo-500/20 border-indigo-400': selectedPriority }"
            >
              <Filter class="w-5 h-5 text-gray-500" />
              <span v-if="selectedPriority" class="text-sm font-medium text-indigo-600">
                {{ priorityLabel }}
              </span>
              <ChevronDown class="w-4 h-4 text-gray-400" />
            </button>

            <!-- Priority Dropdown -->
            <Motion
              v-if="showPriorityFilter"
              :initial="{ opacity: 0, y: -10, scale: 0.95 }"
              :animate="{ opacity: 1, y: 0, scale: 1 }"
              :transition="{ duration: 0.15, ease: 'easeOut' }"
              class="absolute right-0 top-full mt-2 w-44 bg-white rounded-xl shadow-xl 
                     border border-gray-200/50 overflow-hidden z-50"
            >
              <button
                v-for="option in priorityOptions"
                :key="option.value"
                @click="selectPriority(option.value)"
                class="w-full px-4 py-3 text-left text-sm hover:bg-gray-50 
                       flex items-center justify-between"
                :class="{ 'bg-indigo-50 text-indigo-700': selectedPriority === option.value }"
              >
                {{ option.label }}
                <Check v-if="selectedPriority === option.value" class="w-4 h-4" />
              </button>
            </Motion>
          </div>
        </div>
      </Motion>

      <!-- Queue Grid -->
      <div v-if="!loading && queueItems.length > 0" class="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
        <CetakQueueCard
          v-for="(item, index) in queueItems"
          :key="item.po_id"
          :item="item"
          :index="index"
          @click="handleCardClick"
        />
      </div>

      <!-- Empty State -->
      <Motion
        v-if="!loading && queueItems.length === 0"
        v-bind="entranceAnimations.fadeUp"
        class="flex flex-col items-center justify-center py-20"
      >
        <div class="w-20 h-20 rounded-2xl bg-gradient-to-br from-indigo-100 to-fuchsia-100 
                    flex items-center justify-center mb-4">
          <ClipboardList class="w-10 h-10 text-indigo-500" />
        </div>
        <h3 class="text-lg font-semibold text-gray-900 mb-2">
          Tidak Ada Antrian
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
        size="lg"
        :show-footer="false"
        scrollable
      >
        <div v-if="detailLoading" class="flex items-center justify-center py-12">
          <div class="w-8 h-8 border-3 border-indigo-200 border-t-indigo-600 rounded-full animate-spin" />
        </div>

        <div v-else-if="detailData" class="space-y-6">
          <!-- Product Info -->
          <div class="glass-card rounded-xl p-4">
            <h4 class="text-sm font-semibold text-gray-700 mb-3 flex items-center gap-2">
              <Package class="w-4 h-4" />
              Informasi Produk
            </h4>
            <p class="text-gray-900 font-medium mb-2">
              {{ detailData.product_name }}
            </p>
            <div class="grid grid-cols-2 gap-3 text-sm">
              <div>
                <p class="text-gray-500">Quantity</p>
                <p class="font-semibold">{{ detailData.quantity_ordered.toLocaleString() }}</p>
              </div>
              <div>
                <p class="text-gray-500">Est. Rims</p>
                <p class="font-semibold">{{ detailData.estimated_rims }}</p>
              </div>
              <div>
                <p class="text-gray-500">Target Lembar</p>
                <p class="font-semibold">{{ detailData.quantity_target_lembar_besar?.toLocaleString() || '-' }}</p>
              </div>
              <div>
                <p class="text-gray-500">Due Date</p>
                <p class="font-semibold" :class="detailData.is_past_due ? 'text-red-600' : ''">
                  {{ formatDate(detailData.due_date) }}
                </p>
              </div>
            </div>
          </div>

          <!-- Material Prep Info -->
          <div v-if="detailData.material_prep" class="glass-card rounded-xl p-4">
            <h4 class="text-sm font-semibold text-gray-700 mb-3 flex items-center gap-2">
              <ClipboardCheck class="w-4 h-4" />
              Info Persiapan Material
            </h4>
            <div class="grid grid-cols-2 gap-3 text-sm mb-4">
              <div>
                <p class="text-gray-500">Prepared By</p>
                <p class="font-semibold">{{ detailData.material_prep.prepared_by_name }}</p>
              </div>
              <div>
                <p class="text-gray-500">Durasi</p>
                <p class="font-semibold">{{ detailData.material_prep.duration_minutes }} menit</p>
              </div>
              <div>
                <p class="text-gray-500">Plat Code</p>
                <p class="font-semibold">{{ detailData.material_prep.sap_plat_code }}</p>
              </div>
              <div>
                <p class="text-gray-500">Kertas Blanko</p>
                <p class="font-semibold">
                  {{ detailData.material_prep.kertas_blanko_actual || detailData.material_prep.kertas_blanko_quantity }}
                  <span v-if="detailData.material_prep.kertas_blanko_variance !== 0" 
                        :class="detailData.material_prep.kertas_blanko_variance > 0 ? 'text-green-600' : 'text-red-600'"
                  >
                    ({{ detailData.material_prep.kertas_blanko_variance > 0 ? '+' : '' }}{{ detailData.material_prep.kertas_blanko_variance }})
                  </span>
                </p>
              </div>
            </div>

            <!-- Notes jika ada -->
            <div v-if="detailData.material_prep.notes" class="p-3 rounded-lg bg-gray-50">
              <p class="text-xs text-gray-500 mb-1">Catatan</p>
              <p class="text-sm text-gray-700">{{ detailData.material_prep.notes }}</p>
            </div>
          </div>

          <!-- Material Photos -->
          <div v-if="detailData.material_prep?.material_photos?.length > 0" class="glass-card rounded-xl p-4">
            <h4 class="text-sm font-semibold text-gray-700 mb-3 flex items-center gap-2">
              <Image class="w-4 h-4" />
              Foto Material ({{ detailData.material_prep.material_photos.length }})
            </h4>
            <div class="grid grid-cols-3 gap-2">
              <button
                v-for="(photo, idx) in detailData.material_prep.material_photos"
                :key="idx"
                @click="openPhotoViewer(idx)"
                class="aspect-square rounded-lg overflow-hidden hover:opacity-90 active-scale"
              >
                <img 
                  :src="photo" 
                  :alt="`Material photo ${idx + 1}`"
                  class="w-full h-full object-cover"
                />
              </button>
            </div>
          </div>
        </div>
      </BaseModal>

      <!-- Photo Viewer -->
      <MaterialPhotoViewer
        v-model="showPhotoViewer"
        :photos="detailPhotos"
        :initial-index="currentPhotoIndex"
        :title="`Foto Material - PO #${selectedItem?.obc_number || ''}`"
      />
    </div>
  </AppLayout>
</template>

<script setup>
/**
 * CetakQueuePage - Halaman antrian cetak untuk Unit Cetak
 * yang menampilkan list PO siap cetak dengan filter dan search
 */
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { Motion } from 'motion-v'
import { entranceAnimations } from '@/composables/useMotion'
import { useCetakApi } from '@/composables/useCetakApi'
import { useAlertDialog } from '@/composables/useModal'
import AppLayout from '@/components/layout/AppLayout.vue'
import BaseModal from '@/components/common/BaseModal.vue'
import MaterialPhotoViewer from '@/components/common/MaterialPhotoViewer.vue'
import CetakQueueCard from '@/components/cetak/CetakQueueCard.vue'
import { 
  Search, 
  Filter, 
  RefreshCw, 
  ChevronDown,
  ChevronLeft,
  ChevronRight,
  Check,
  ClipboardList,
  Package,
  ClipboardCheck,
  Image
} from 'lucide-vue-next'

const cetakApi = useCetakApi()
const alertDialog = useAlertDialog()

// State
const loading = ref(false)
const queueItems = ref([])
const totalItems = ref(0)
const currentPage = ref(1)
const totalPages = ref(0)
const perPage = ref(12)

// Search & Filter
const searchQuery = ref('')
const selectedPriority = ref('')
const showPriorityFilter = ref(false)
let searchTimeout = null

// Detail Modal
const showDetailModal = ref(false)
const selectedItem = ref(null)
const detailLoading = ref(false)
const detailData = ref(null)

// Photo Viewer
const showPhotoViewer = ref(false)
const currentPhotoIndex = ref(0)

// Priority options
const priorityOptions = [
  { value: '', label: 'Semua Prioritas' },
  { value: 'URGENT', label: 'Urgent' },
  { value: 'NORMAL', label: 'Normal' },
  { value: 'LOW', label: 'Rendah' }
]

/**
 * Computed untuk priority label
 */
const priorityLabel = computed(() => {
  const option = priorityOptions.find(o => o.value === selectedPriority.value)
  return option?.label || 'Filter'
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
  if (searchQuery.value || selectedPriority.value) {
    return 'Tidak ada PO yang sesuai dengan filter. Coba ubah kriteria pencarian.'
  }
  return 'Belum ada PO yang siap untuk dicetak. PO akan muncul setelah material preparation selesai.'
})

/**
 * Computed untuk detail photos array
 */
const detailPhotos = computed(() => {
  return detailData.value?.material_prep?.material_photos || []
})

/**
 * Fetch queue data dari API
 */
const fetchQueue = async () => {
  loading.value = true
  try {
    const response = await cetakApi.getQueue({
      search: searchQuery.value,
      priority: selectedPriority.value,
      page: currentPage.value,
      per_page: perPage.value
    })

    if (response.success) {
      queueItems.value = response.data.items
      totalItems.value = response.data.total
      totalPages.value = response.data.total_pages
    }
  } catch (error) {
    console.error('Error fetching queue:', error)
    alertDialog.error('Gagal memuat antrian cetak', {
      detail: error.response?.data?.message || 'Silakan coba lagi'
    })
  } finally {
    loading.value = false
  }
}

/**
 * Refresh queue data
 */
const refreshQueue = () => {
  fetchQueue()
  
  // Haptic feedback
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
    fetchQueue()
  }, 300)
}

/**
 * Select priority filter
 */
const selectPriority = (value) => {
  selectedPriority.value = value
  showPriorityFilter.value = false
  currentPage.value = 1
  fetchQueue()
  
  // Haptic feedback
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
  fetchQueue()
  
  // Scroll to top
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

/**
 * Handle card click untuk show detail
 */
const handleCardClick = async (item) => {
  selectedItem.value = item
  showDetailModal.value = true
  detailLoading.value = true
  detailData.value = null

  try {
    const response = await cetakApi.getDetail(item.po_id)
    if (response.success) {
      detailData.value = response.data
    }
  } catch (error) {
    console.error('Error fetching detail:', error)
    alertDialog.error('Gagal memuat detail', {
      detail: error.response?.data?.message || 'Silakan coba lagi'
    })
  } finally {
    detailLoading.value = false
  }
}

/**
 * Open photo viewer
 */
const openPhotoViewer = (index) => {
  currentPhotoIndex.value = index
  showPhotoViewer.value = true
}

/**
 * Format date untuk display
 */
const formatDate = (dateStr) => {
  try {
    const date = new Date(dateStr)
    return date.toLocaleDateString('id-ID', {
      day: 'numeric',
      month: 'short',
      year: 'numeric'
    })
  } catch (e) {
    return dateStr
  }
}

/**
 * Close priority filter on click outside
 */
const handleClickOutside = (e) => {
  if (showPriorityFilter.value && !e.target.closest('.relative')) {
    showPriorityFilter.value = false
  }
}

onMounted(() => {
  fetchQueue()
  document.addEventListener('click', handleClickOutside)
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
  if (searchTimeout) clearTimeout(searchTimeout)
})
</script>
