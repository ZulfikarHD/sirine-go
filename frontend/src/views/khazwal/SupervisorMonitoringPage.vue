<template>
  <AppLayout>
    <div class="min-h-screen pb-20">
      <!-- Header Section -->
      <Motion v-bind="entranceAnimations.fadeUp" class="mb-6">
        <div class="flex items-center justify-between">
          <div>
            <h1 class="text-2xl sm:text-3xl font-bold text-gray-900">
              Monitoring
            </h1>
            <p class="text-sm text-gray-500 mt-1">
              Pantau aktivitas tim Khazanah Awal
            </p>
          </div>
          
          <!-- Refresh Button -->
          <button
            @click="refreshData"
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
      </Motion>

      <!-- Stats Cards -->
      <Motion
        :initial="{ opacity: 0, y: 15 }"
        :animate="{ opacity: 1, y: 0 }"
        :transition="{ duration: 0.25, delay: 0.1, ease: 'easeOut' }"
        class="grid grid-cols-2 md:grid-cols-4 gap-4 mb-6"
      >
        <!-- In Queue -->
        <div class="glass-card rounded-2xl p-4">
          <div class="flex items-center gap-3 mb-2">
            <div class="p-2 rounded-xl bg-blue-100">
              <ClipboardList class="w-5 h-5 text-blue-600" />
            </div>
          </div>
          <p class="text-2xl font-bold text-gray-900">
            {{ stats.total_in_queue || 0 }}
          </p>
          <p class="text-sm text-gray-500">Antrian</p>
        </div>

        <!-- In Progress -->
        <div class="glass-card rounded-2xl p-4">
          <div class="flex items-center gap-3 mb-2">
            <div class="p-2 rounded-xl bg-yellow-100">
              <Loader class="w-5 h-5 text-yellow-600" />
            </div>
          </div>
          <p class="text-2xl font-bold text-gray-900">
            {{ stats.total_in_progress || 0 }}
          </p>
          <p class="text-sm text-gray-500">Diproses</p>
        </div>

        <!-- Completed Today -->
        <div class="glass-card rounded-2xl p-4">
          <div class="flex items-center gap-3 mb-2">
            <div class="p-2 rounded-xl bg-emerald-100">
              <CheckCircle class="w-5 h-5 text-emerald-600" />
            </div>
          </div>
          <p class="text-2xl font-bold text-gray-900">
            {{ stats.total_completed_today || 0 }}
          </p>
          <p class="text-sm text-gray-500">Selesai Hari Ini</p>
        </div>

        <!-- Average Duration -->
        <div class="glass-card rounded-2xl p-4">
          <div class="flex items-center gap-3 mb-2">
            <div class="p-2 rounded-xl bg-purple-100">
              <Clock class="w-5 h-5 text-purple-600" />
            </div>
          </div>
          <p class="text-2xl font-bold text-gray-900">
            {{ formatDuration(stats.average_duration_mins || 0) }}
          </p>
          <p class="text-sm text-gray-500">Rata-rata Durasi</p>
        </div>
      </Motion>

      <!-- Staff Activity Section -->
      <Motion
        :initial="{ opacity: 0, y: 15 }"
        :animate="{ opacity: 1, y: 0 }"
        :transition="{ duration: 0.25, delay: 0.2, ease: 'easeOut' }"
        class="mb-6"
      >
        <h2 class="text-lg font-bold text-gray-900 mb-4 flex items-center gap-2">
          <Users class="w-5 h-5 text-indigo-600" />
          Aktivitas Staff
        </h2>

        <div v-if="staffActive.length > 0" class="grid gap-4 md:grid-cols-2">
          <StaffActivityCard
            v-for="(staff, index) in staffActive"
            :key="staff.user_id"
            :staff="staff"
            :index="index"
          />
        </div>

        <div v-else class="glass-card rounded-2xl p-8 text-center">
          <UserX class="w-12 h-12 text-gray-300 mx-auto mb-3" />
          <p class="text-gray-500">Tidak ada staff yang sedang aktif</p>
        </div>
      </Motion>

      <!-- Recent Completions Section -->
      <Motion
        :initial="{ opacity: 0, y: 15 }"
        :animate="{ opacity: 1, y: 0 }"
        :transition="{ duration: 0.25, delay: 0.3, ease: 'easeOut' }"
      >
        <h2 class="text-lg font-bold text-gray-900 mb-4 flex items-center gap-2">
          <History class="w-5 h-5 text-indigo-600" />
          Penyelesaian Terbaru
        </h2>

        <div v-if="recentCompletions.length > 0" class="space-y-3">
          <Motion
            v-for="(item, index) in recentCompletions"
            :key="item.prep_id"
            :initial="{ opacity: 0, x: -15 }"
            :animate="{ opacity: 1, x: 0 }"
            :transition="{ duration: 0.25, delay: index * 0.05, ease: 'easeOut' }"
            class="glass-card rounded-xl p-4 flex items-center gap-4"
          >
            <!-- Timeline Dot -->
            <div class="relative flex-shrink-0">
              <div class="w-3 h-3 rounded-full bg-emerald-500" />
              <div 
                v-if="index < recentCompletions.length - 1"
                class="absolute top-4 left-1/2 -translate-x-1/2 w-0.5 h-12 bg-gray-200"
              />
            </div>

            <!-- Content -->
            <div class="flex-1 min-w-0">
              <div class="flex items-center justify-between gap-2 mb-1">
                <p class="text-sm font-semibold text-gray-900 truncate">
                  PO #{{ item.obc_number }}
                </p>
                <span class="text-xs text-gray-500 flex-shrink-0">
                  {{ formatTimeAgo(item.completed_at) }}
                </span>
              </div>
              <div class="flex items-center justify-between gap-2">
                <p class="text-xs text-gray-500 truncate">
                  {{ item.prepared_by_name }}
                </p>
                <span class="text-xs font-medium text-indigo-600">
                  {{ item.duration_minutes }} menit
                </span>
              </div>
            </div>
          </Motion>
        </div>

        <div v-else class="glass-card rounded-2xl p-8 text-center">
          <Clock class="w-12 h-12 text-gray-300 mx-auto mb-3" />
          <p class="text-gray-500">Belum ada penyelesaian hari ini</p>
        </div>
      </Motion>

      <!-- Auto Refresh Indicator -->
      <div class="fixed bottom-20 left-1/2 -translate-x-1/2 z-30">
        <Motion
          v-if="autoRefreshEnabled"
          :initial="{ opacity: 0, y: 20 }"
          :animate="{ opacity: 1, y: 0 }"
          class="px-4 py-2 rounded-full bg-gray-900/80 text-white text-xs font-medium 
                 flex items-center gap-2 backdrop-blur-sm"
        >
          <div class="w-2 h-2 rounded-full bg-emerald-400 animate-pulse" />
          Auto-refresh: {{ refreshCountdown }}s
        </Motion>
      </div>
    </div>
  </AppLayout>
</template>

<script setup>
/**
 * SupervisorMonitoringPage - Dashboard monitoring untuk Supervisor Khazwal
 * dengan stats, staff activity, dan recent completions
 */
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { Motion } from 'motion-v'
import { entranceAnimations } from '@/composables/useMotion'
import { useKhazwalApi } from '@/composables/useKhazwalApi'
import { useAlertDialog } from '@/composables/useModal'
import AppLayout from '@/components/layout/AppLayout.vue'
import StaffActivityCard from '@/components/khazwal/StaffActivityCard.vue'
import { 
  RefreshCw, 
  ClipboardList,
  Loader,
  CheckCircle,
  Clock,
  Users,
  UserX,
  History
} from 'lucide-vue-next'

const khazwalApi = useKhazwalApi()
const alertDialog = useAlertDialog()

// State
const loading = ref(false)
const stats = ref({})
const staffActive = ref([])
const recentCompletions = ref([])

// Auto Refresh
const autoRefreshEnabled = ref(true)
const refreshCountdown = ref(30)
const AUTO_REFRESH_INTERVAL = 30 // seconds
let refreshIntervalId = null
let countdownIntervalId = null

/**
 * Fetch monitoring data dari API
 */
const fetchMonitoringData = async () => {
  loading.value = true
  try {
    const response = await khazwalApi.getMonitoring()

    if (response.success) {
      stats.value = {
        total_in_queue: response.data.total_in_queue,
        total_in_progress: response.data.total_in_progress,
        total_completed_today: response.data.total_completed_today,
        average_duration_mins: response.data.average_duration_mins
      }
      staffActive.value = response.data.staff_active || []
      recentCompletions.value = response.data.recent_completions || []
    }
  } catch (error) {
    console.error('Error fetching monitoring data:', error)
    // Jangan tampilkan error untuk auto-refresh failure
    if (!autoRefreshEnabled.value) {
      alertDialog.error('Gagal memuat data monitoring', {
        detail: error.response?.data?.message || 'Silakan coba lagi'
      })
    }
  } finally {
    loading.value = false
  }
}

/**
 * Manual refresh data
 */
const refreshData = () => {
  fetchMonitoringData()
  resetCountdown()
  
  if ('vibrate' in navigator) {
    navigator.vibrate(10)
  }
}

/**
 * Reset countdown
 */
const resetCountdown = () => {
  refreshCountdown.value = AUTO_REFRESH_INTERVAL
}

/**
 * Start auto refresh
 */
const startAutoRefresh = () => {
  // Clear existing intervals
  stopAutoRefresh()

  // Start countdown
  countdownIntervalId = setInterval(() => {
    refreshCountdown.value--
    
    if (refreshCountdown.value <= 0) {
      fetchMonitoringData()
      resetCountdown()
    }
  }, 1000)
}

/**
 * Stop auto refresh
 */
const stopAutoRefresh = () => {
  if (refreshIntervalId) {
    clearInterval(refreshIntervalId)
    refreshIntervalId = null
  }
  if (countdownIntervalId) {
    clearInterval(countdownIntervalId)
    countdownIntervalId = null
  }
}

/**
 * Format duration untuk display
 */
const formatDuration = (minutes) => {
  if (!minutes) return '0m'
  
  if (minutes < 60) {
    return `${minutes}m`
  }
  
  const hours = Math.floor(minutes / 60)
  const mins = minutes % 60
  
  if (mins === 0) {
    return `${hours}j`
  }
  
  return `${hours}j ${mins}m`
}

/**
 * Format time ago
 */
const formatTimeAgo = (dateStr) => {
  try {
    const date = new Date(dateStr)
    const now = new Date()
    const diffMs = now - date
    const diffMins = Math.floor(diffMs / 60000)
    const diffHours = Math.floor(diffMs / 3600000)

    if (diffMins < 1) {
      return 'Baru saja'
    }
    if (diffMins < 60) {
      return `${diffMins} menit lalu`
    }
    if (diffHours < 24) {
      return `${diffHours} jam lalu`
    }
    
    return date.toLocaleDateString('id-ID', {
      day: 'numeric',
      month: 'short'
    })
  } catch (e) {
    return dateStr
  }
}

onMounted(() => {
  fetchMonitoringData()
  startAutoRefresh()
})

onUnmounted(() => {
  stopAutoRefresh()
})
</script>
