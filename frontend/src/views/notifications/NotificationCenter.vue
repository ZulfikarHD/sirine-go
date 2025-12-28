<template>
  <AppLayout>
    <!-- Header dengan Motion-V animation -->
    <Motion v-bind="entranceAnimations.fadeUp" class="glass-card p-6 rounded-2xl mb-6">
      <div class="flex items-center justify-between flex-wrap gap-4">
        <div>
          <h1 class="text-3xl font-bold text-gray-900">Pusat Notifikasi</h1>
          <p class="text-gray-600 mt-1">Kelola semua notifikasi Anda</p>
        </div>
        
        <!-- Actions -->
        <div class="flex items-center gap-3">
          <button
            v-if="hasUnread"
            @click="handleMarkAllAsRead"
            :disabled="isLoading"
            class="inline-flex items-center px-4 py-2 bg-indigo-600 hover:bg-indigo-700 text-white rounded-lg text-sm font-medium active-scale transition-colors duration-150 focus:outline-none focus:ring-2 focus:ring-indigo-100"
          >
            <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
            </svg>
            Tandai Semua Dibaca
          </button>
        </div>
      </div>
    </Motion>

    <!-- Tabs -->
    <Motion
      v-bind="entranceAnimations.fadeUp"
      class="mb-6"
      :style="{ transitionDelay: '0.1s' }"
    >
      <div class="glass-card p-1 rounded-xl inline-flex">
          <button
            v-for="tab in tabs"
            :key="tab.value"
            @click="currentTab = tab.value"
            class="px-4 py-2 rounded-lg text-sm font-medium transition-all duration-200 active-scale"
            :class="currentTab === tab.value 
              ? 'bg-linear-to-r from-indigo-600 to-fuchsia-600 text-white shadow-md' 
              : 'text-gray-600 hover:bg-gray-50'"
          >
            {{ tab.label }}
            <span
              v-if="tab.value === 'unread' && unreadCount > 0"
              class="ml-2 inline-flex items-center justify-center px-2 py-0.5 text-xs font-bold rounded-full"
              :class="currentTab === tab.value ? 'bg-white/20 text-white' : 'bg-indigo-100 text-indigo-600'"
            >
              {{ unreadCount }}
            </span>
          </button>
        </div>
      </Motion>

    <!-- Loading State -->
    <Motion
      v-if="isLoading && notifications.length === 0"
      v-bind="entranceAnimations.fadeUp"
      class="glass-card p-12 rounded-2xl text-center"
    >
      <div class="inline-block w-8 h-8 border-4 border-gray-300 border-t-indigo-600 rounded-full animate-spin"></div>
      <p class="mt-4 text-gray-600">Memuat notifikasi...</p>
    </Motion>

    <!-- Notifications List -->
    <div v-else-if="filteredNotifications.length > 0" class="space-y-3">
        <Motion
          v-for="(notification, index) in filteredNotifications"
          :key="notification.id"
          :initial="{ opacity: 0, y: 15 }"
          :animate="{ opacity: 1, y: 0 }"
          :transition="{ duration: 0.25, delay: index * 0.05, ease: 'easeOut' }"
          class="glass-card rounded-2xl overflow-hidden hover:shadow-md transition-shadow duration-200"
          :class="{ 'ring-2 ring-indigo-200': !notification.is_read }"
        >
          <div class="p-4 sm:p-5">
            <div class="flex items-start gap-4">
              <!-- Icon Badge -->
              <div
                class="shrink-0 w-12 h-12 rounded-xl flex items-center justify-center"
                :class="getNotificationBgClass(notification.type)"
              >
                <svg
                  class="w-6 h-6"
                  :class="getNotificationTextClass(notification.type)"
                  fill="currentColor"
                  viewBox="0 0 20 20"
                >
                  <path v-if="notification.type === 'SUCCESS'" fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
                  <path v-else-if="notification.type === 'ERROR'" fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd" />
                  <path v-else-if="notification.type === 'WARNING'" fill-rule="evenodd" d="M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.213 2.98-1.742 2.98H4.42c-1.53 0-2.493-1.646-1.743-2.98l5.58-9.92zM11 13a1 1 0 11-2 0 1 1 0 012 0zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z" clip-rule="evenodd" />
                  <path v-else-if="notification.type === 'ACHIEVEMENT'" d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z" />
                  <path v-else fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z" clip-rule="evenodd" />
                </svg>
              </div>

              <!-- Content -->
              <div class="flex-1 min-w-0">
                <div class="flex items-start justify-between gap-2">
                  <h3 class="text-base font-semibold text-gray-900">
                    {{ notification.title }}
                  </h3>
                  <div v-if="!notification.is_read" class="shrink-0 w-2.5 h-2.5 bg-indigo-600 rounded-full"></div>
                </div>
                
                <p class="mt-1 text-sm text-gray-600">
                  {{ notification.message }}
                </p>
                
                <div class="mt-3 flex items-center gap-4 text-xs text-gray-500">
                  <span class="flex items-center gap-1">
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                    </svg>
                    {{ formatTimeAgo(notification.created_at) }}
                  </span>
                  <span v-if="notification.read_at" class="flex items-center gap-1 text-green-600">
                    <svg class="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
                      <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
                    </svg>
                    Dibaca
                  </span>
                </div>

                <!-- Actions -->
                <div class="mt-4 flex items-center gap-2">
                  <button
                    v-if="!notification.is_read"
                    @click="handleMarkAsRead(notification.id)"
                    class="inline-flex items-center px-3 py-1.5 text-xs font-medium text-indigo-600 hover:bg-indigo-50 rounded-lg transition-colors duration-150 active-scale"
                  >
                    <svg class="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
                    </svg>
                    Tandai Dibaca
                  </button>
                  
                  <button
                    @click="handleDelete(notification.id)"
                    class="inline-flex items-center px-3 py-1.5 text-xs font-medium text-red-600 hover:bg-red-50 rounded-lg transition-colors duration-150 active-scale"
                  >
                    <svg class="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                    </svg>
                    Hapus
                  </button>
                </div>
              </div>
            </div>
          </div>
        </Motion>
      </div>

    <!-- Empty State -->
    <Motion
      v-else
      v-bind="entranceAnimations.fadeScale"
      class="glass-card p-12 rounded-2xl text-center"
    >
        <svg class="w-16 h-16 mx-auto text-gray-300 mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9" />
        </svg>
        <h3 class="text-lg font-semibold text-gray-900 mb-2">
          {{ currentTab === 'unread' ? 'Tidak ada notifikasi baru' : 'Tidak ada notifikasi' }}
        </h3>
      <p class="text-gray-600">
        {{ currentTab === 'unread' ? 'Semua notifikasi sudah dibaca' : 'Belum ada notifikasi untuk ditampilkan' }}
      </p>
    </Motion>
  </AppLayout>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { Motion } from 'motion-v'
import { entranceAnimations } from '@/composables/useMotion'
import { useNotificationStore } from '@/stores/notification'
import { useAlertDialog } from '@/composables/useModal'
import { formatDistanceToNow } from 'date-fns'
import { id as idLocale } from 'date-fns/locale'
import AppLayout from '@/components/layout/AppLayout.vue'

const notificationStore = useNotificationStore()
const alert = useAlertDialog()

// State
const currentTab = ref('all')
const isLoading = ref(false)

const tabs = [
  { label: 'Semua', value: 'all' },
  { label: 'Belum Dibaca', value: 'unread' },
]

// Computed
const notifications = computed(() => notificationStore.notifications)
const unreadCount = computed(() => notificationStore.unreadCount)
const hasUnread = computed(() => notificationStore.hasUnread)

const filteredNotifications = computed(() => {
  if (currentTab.value === 'unread') {
    return notifications.value.filter(n => !n.is_read)
  }
  return notifications.value
})

/**
 * fetchNotifications mengambil semua notifikasi
 * dengan filter berdasarkan tab yang aktif
 */
const fetchNotifications = async () => {
  try {
    isLoading.value = true
    const unreadOnly = currentTab.value === 'unread'
    await notificationStore.fetchNotifications(unreadOnly)
  } catch (error) {
    await alert.error('Gagal memuat notifikasi', { detail: error.message })
  } finally {
    isLoading.value = false
  }
}

/**
 * handleMarkAsRead menandai satu notifikasi sebagai dibaca
 * dengan haptic feedback
 */
const handleMarkAsRead = async (notificationId) => {
  try {
    await notificationStore.markAsRead(notificationId)
    
    // Haptic feedback
    if ('vibrate' in navigator) {
      navigator.vibrate(10)
    }
  } catch (error) {
    await alert.error('Gagal menandai notifikasi', { detail: error.message })
  }
}

/**
 * handleMarkAllAsRead menandai semua notifikasi sebagai dibaca
 * dengan haptic feedback
 */
const handleMarkAllAsRead = async () => {
  try {
    await notificationStore.markAllAsRead()
    await alert.success('Semua notifikasi berhasil ditandai sebagai dibaca')
    
    // Haptic feedback
    if ('vibrate' in navigator) {
      navigator.vibrate([10, 50, 10])
    }
  } catch (error) {
    await alert.error('Gagal menandai semua notifikasi', { detail: error.message })
  }
}

/**
 * handleDelete menghapus notifikasi
 * dengan haptic feedback
 */
const handleDelete = async (notificationId) => {
  try {
    await notificationStore.deleteNotification(notificationId)
    
    // Haptic feedback
    if ('vibrate' in navigator) {
      navigator.vibrate(30)
    }
    
    await alert.success('Notifikasi berhasil dihapus')
  } catch (error) {
    await alert.error('Gagal menghapus notifikasi', { detail: error.message })
  }
}

/**
 * formatTimeAgo format timestamp menjadi relative time
 */
const formatTimeAgo = (timestamp) => {
  return formatDistanceToNow(new Date(timestamp), {
    addSuffix: true,
    locale: idLocale,
  })
}

/**
 * getNotificationBgClass mengembalikan class background berdasarkan type
 */
const getNotificationBgClass = (type) => {
  const classes = {
    INFO: 'bg-blue-100',
    SUCCESS: 'bg-green-100',
    WARNING: 'bg-yellow-100',
    ERROR: 'bg-red-100',
    ACHIEVEMENT: 'bg-purple-100',
  }
  return classes[type] || 'bg-gray-100'
}

/**
 * getNotificationTextClass mengembalikan class text color berdasarkan type
 */
const getNotificationTextClass = (type) => {
  const classes = {
    INFO: 'text-blue-600',
    SUCCESS: 'text-green-600',
    WARNING: 'text-yellow-600',
    ERROR: 'text-red-600',
    ACHIEVEMENT: 'text-purple-600',
  }
  return classes[type] || 'text-gray-600'
}

// Lifecycle: Fetch notifications saat component mounted
onMounted(async () => {
  await fetchNotifications()
  await notificationStore.fetchUnreadCount()
})
</script>
