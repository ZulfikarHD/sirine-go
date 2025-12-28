<template>
  <AppLayout>
    <div class="max-w-7xl mx-auto">
      <!-- Header dengan Motion-V animation -->
      <Motion v-bind="entranceAnimations.fadeUp" class="glass-card p-6 rounded-2xl mb-6">
        <div class="flex items-center justify-between flex-wrap gap-4">
          <div>
            <h1 class="text-3xl font-bold text-gray-900">Audit Logs</h1>
            <p class="text-gray-600 mt-1">Monitor aktivitas pengguna dalam sistem</p>
          </div>
        </div>
      </Motion>

      <!-- Filters Card -->
      <Motion
        v-bind="entranceAnimations.fadeUp"
        class="glass-card p-6 rounded-2xl mb-6"
        :style="{ transitionDelay: '0.1s' }"
      >
        <div>
          <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
            <!-- Action Filter -->
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">
                Jenis Aksi
              </label>
              <select
                v-model="filters.action"
                @change="handleFilterChange"
                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-100 focus:border-indigo-500 transition-colors duration-150"
              >
                <option value="">Semua Aksi</option>
                <option value="CREATE">Buat</option>
                <option value="UPDATE">Update</option>
                <option value="DELETE">Hapus</option>
                <option value="LOGIN">Login</option>
                <option value="LOGOUT">Logout</option>
                <option value="PASSWORD_CHANGE">Ganti Password</option>
              </select>
            </div>

            <!-- Entity Type Filter -->
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">
                Tipe Entity
              </label>
              <select
                v-model="filters.entity_type"
                @change="handleFilterChange"
                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-100 focus:border-indigo-500 transition-colors duration-150"
              >
                <option value="">Semua Entity</option>
                <option value="users">Users</option>
                <option value="profile">Profile</option>
              </select>
            </div>

            <!-- Search -->
            <div class="md:col-span-2">
              <label class="block text-sm font-medium text-gray-700 mb-2">
                Cari
              </label>
              <div class="relative">
                <input
                  v-model="filters.search"
                  @input="handleSearchDebounced"
                  type="text"
                  placeholder="Cari berdasarkan entity ID atau tipe..."
                  class="w-full pl-10 pr-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-100 focus:border-indigo-500 transition-colors duration-150"
                />
                <svg class="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
                </svg>
              </div>
            </div>
          </div>

          <!-- Filter Actions -->
        <div class="mt-4 flex items-center justify-between">
          <button
            v-if="hasActiveFilters"
            @click="resetFilters"
            class="text-sm text-indigo-600 hover:text-indigo-700 font-medium active-scale"
          >
            Reset Filter
          </button>
          <div class="text-sm text-gray-600">
            Menampilkan <span class="font-semibold text-gray-900">{{ logs.length }}</span> dari <span class="font-semibold text-gray-900">{{ meta.total }}</span> log
          </div>
        </div>
      </div>
    </Motion>

    <!-- Loading State -->
    <Motion
      v-if="isLoading && logs.length === 0"
      v-bind="entranceAnimations.fadeUp"
      class="glass-card p-12 rounded-2xl text-center"
    >
      <div class="inline-block w-8 h-8 border-4 border-gray-300 border-t-indigo-600 rounded-full animate-spin"></div>
      <p class="mt-4 text-gray-600">Memuat activity logs...</p>
    </Motion>

    <!-- Activity Logs Table -->
    <Motion
      v-else
      v-bind="entranceAnimations.fadeUp"
      :style="{ transitionDelay: '0.2s' }"
    >
      <ActivityLogTable :logs="logs" />
    </Motion>

    <!-- Pagination -->
    <Motion
      v-if="meta.total_pages > 1"
      v-bind="entranceAnimations.fadeUp"
      class="mt-6"
      :style="{ transitionDelay: '0.3s' }"
    >
      <div class="glass-card px-6 py-4 rounded-2xl">
          <div class="flex items-center justify-between flex-wrap gap-4">
            <!-- Page Info -->
            <div class="text-sm text-gray-600">
              Halaman <span class="font-semibold text-gray-900">{{ meta.page }}</span> dari <span class="font-semibold text-gray-900">{{ meta.total_pages }}</span>
            </div>

            <!-- Pagination Controls -->
            <div class="flex items-center gap-2">
              <button
                @click="goToPage(meta.page - 1)"
                :disabled="meta.page === 1"
                class="px-3 py-2 border border-gray-300 rounded-lg text-sm font-medium text-gray-700 hover:bg-gray-50 active-scale transition-colors duration-150 disabled:opacity-50 disabled:cursor-not-allowed"
              >
                <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
                </svg>
              </button>

              <!-- Page Numbers -->
              <div class="hidden sm:flex items-center gap-1">
                <button
                  v-for="page in visiblePages"
                  :key="page"
                  @click="goToPage(page)"
                  class="px-3 py-2 rounded-lg text-sm font-medium transition-all duration-150 active-scale"
                  :class="page === meta.page 
                    ? 'bg-linear-to-r from-indigo-600 to-fuchsia-600 text-white shadow-md' 
                    : 'text-gray-700 hover:bg-gray-50'"
                >
                  {{ page }}
                </button>
              </div>

              <button
                @click="goToPage(meta.page + 1)"
                :disabled="meta.page === meta.total_pages"
                class="px-3 py-2 border border-gray-300 rounded-lg text-sm font-medium text-gray-700 hover:bg-gray-50 active-scale transition-colors duration-150 disabled:opacity-50 disabled:cursor-not-allowed"
              >
                <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
                </svg>
              </button>
            </div>

            <!-- Page Size Selector -->
            <div class="flex items-center gap-2">
              <label class="text-sm text-gray-600">Tampilkan:</label>
              <select
                v-model="pagination.page_size"
                @change="handlePageSizeChange"
                class="px-3 py-2 border border-gray-300 rounded-lg text-sm focus:ring-2 focus:ring-indigo-100 focus:border-indigo-500 transition-colors duration-150"
              >
                <option :value="10">10</option>
                <option :value="20">20</option>
                <option :value="50">50</option>
                <option :value="100">100</option>
              </select>
          </div>
          </div>
        </div>
      </Motion>
    </div>
  </AppLayout>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { Motion } from 'motion-v'
import { entranceAnimations } from '@/composables/useMotion'
import { useApi } from '@/composables/useApi'
import { useAlertDialog } from '@/composables/useModal'
import ActivityLogTable from '@/components/admin/ActivityLogTable.vue'
import AppLayout from '@/components/layout/AppLayout.vue'

const api = useApi()
const alert = useAlertDialog()

// State
const isLoading = ref(false)
const logs = ref([])
const meta = ref({
  total: 0,
  page: 1,
  page_size: 20,
  total_pages: 0,
})

const filters = ref({
  action: '',
  entity_type: '',
  search: '',
})

const pagination = ref({
  page: 1,
  page_size: 20,
})

// Computed
const hasActiveFilters = computed(() => {
  return filters.value.action || filters.value.entity_type || filters.value.search
})

const visiblePages = computed(() => {
  const current = meta.value.page
  const total = meta.value.total_pages
  const delta = 2 // Show 2 pages on each side
  
  const range = []
  const rangeWithDots = []
  
  for (let i = Math.max(2, current - delta); i <= Math.min(total - 1, current + delta); i++) {
    range.push(i)
  }
  
  if (current - delta > 2) {
    rangeWithDots.push(1, '...')
  } else {
    rangeWithDots.push(1)
  }
  
  rangeWithDots.push(...range)
  
  if (current + delta < total - 1) {
    rangeWithDots.push('...', total)
  } else if (total > 1) {
    rangeWithDots.push(total)
  }
  
  return rangeWithDots.filter((item, index, self) => self.indexOf(item) === index)
})

/**
 * fetchActivityLogs mengambil activity logs dari API
 * dengan filters dan pagination
 */
const fetchActivityLogs = async () => {
  try {
    isLoading.value = true
    
    // Build query params
    const params = new URLSearchParams()
    params.append('page', pagination.value.page.toString())
    params.append('page_size', pagination.value.page_size.toString())
    
    if (filters.value.action) {
      params.append('action', filters.value.action)
    }
    if (filters.value.entity_type) {
      params.append('entity_type', filters.value.entity_type)
    }
    if (filters.value.search) {
      params.append('search', filters.value.search)
    }
    
    const response = await api.get(`/admin/activity-logs?${params.toString()}`)
    
    if (response.success) {
      logs.value = response.data || []
      meta.value = response.meta || meta.value
    }
  } catch (error) {
    console.error('Error fetching activity logs:', error)
    await alert.error('Gagal memuat activity logs', { detail: error.message })
  } finally {
    isLoading.value = false
  }
}

/**
 * handleFilterChange handle perubahan filter
 * dengan reset page ke 1
 */
const handleFilterChange = () => {
  pagination.value.page = 1
  fetchActivityLogs()
  
  // Haptic feedback
  if ('vibrate' in navigator) {
    navigator.vibrate(10)
  }
}

/**
 * handleSearchDebounced handle search dengan debounce
 * untuk menghindari terlalu banyak request
 */
let searchTimeout = null
const handleSearchDebounced = () => {
  clearTimeout(searchTimeout)
  searchTimeout = setTimeout(() => {
    pagination.value.page = 1
    fetchActivityLogs()
  }, 300)
}

/**
 * resetFilters reset semua filter ke default
 */
const resetFilters = () => {
  filters.value = {
    action: '',
    entity_type: '',
    search: '',
  }
  pagination.value.page = 1
  fetchActivityLogs()
  
  // Haptic feedback
  if ('vibrate' in navigator) {
    navigator.vibrate(10)
  }
}

/**
 * goToPage navigate ke halaman tertentu
 */
const goToPage = (page) => {
  if (page < 1 || page > meta.value.total_pages) return
  pagination.value.page = page
  fetchActivityLogs()
  
  // Scroll to top
  window.scrollTo({ top: 0, behavior: 'smooth' })
  
  // Haptic feedback
  if ('vibrate' in navigator) {
    navigator.vibrate(10)
  }
}

/**
 * handlePageSizeChange handle perubahan page size
 * dengan reset page ke 1
 */
const handlePageSizeChange = () => {
  pagination.value.page = 1
  fetchActivityLogs()
}

// Lifecycle: Fetch logs saat component mounted
onMounted(() => {
  fetchActivityLogs()
})
</script>
