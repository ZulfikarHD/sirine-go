<template>
  <div class="relative">
    <!-- Bell Button dengan badge -->
    <button
      @click="toggleDropdown"
      class="relative p-2 rounded-lg hover:bg-gray-100 transition-colors duration-150 active-scale focus:outline-none focus:ring-2 focus:ring-indigo-100"
      aria-label="Notifikasi"
    >
      <!-- Bell Icon -->
      <svg
        class="w-6 h-6 text-gray-600"
        fill="none"
        stroke="currentColor"
        viewBox="0 0 24 24"
      >
        <path
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-width="2"
          d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9"
        />
      </svg>

      <!-- Unread Badge dengan animation -->
      <Motion
        v-if="hasUnread"
        :key="unreadCount"
        :initial="{ scale: 0.8, opacity: 0 }"
        :animate="{ scale: 1, opacity: 1 }"
        :transition="{ type: 'spring', stiffness: 500, damping: 35, mass: 0.6 }"
        class="absolute top-0 right-0 inline-flex items-center justify-center w-5 h-5 text-xs font-bold text-white bg-gradient-to-r from-indigo-600 to-fuchsia-600 rounded-full"
      >
        {{ unreadCount > 99 ? '99+' : unreadCount }}
      </Motion>
    </button>

    <!-- Dropdown Panel dengan Motion-V animation -->
    <Teleport to="body">
      <Motion
        v-if="isOpen"
        :initial="{ opacity: 0 }"
        :animate="{ opacity: 1 }"
        :exit="{ opacity: 0 }"
        :transition="{ duration: 0.2, ease: 'easeOut' }"
        class="fixed inset-0 z-40"
        @click="closeDropdown"
      />
    </Teleport>

    <Motion
      v-if="isOpen"
      :initial="{ opacity: 0, y: -10, scale: 0.95 }"
      :animate="{ opacity: 1, y: 0, scale: 1 }"
      :exit="{ opacity: 0, y: -10, scale: 0.95 }"
      :transition="{ type: 'spring', stiffness: 500, damping: 40, mass: 0.8 }"
      class="absolute right-0 mt-2 w-80 bg-white rounded-xl shadow-xl border border-gray-200/50 z-50 overflow-hidden"
      style="max-height: 480px"
    >
      <!-- Header -->
      <div class="px-4 py-3 border-b border-gray-100 bg-gradient-to-r from-indigo-50 to-fuchsia-50">
        <div class="flex items-center justify-between">
          <h3 class="font-semibold text-gray-900">Notifikasi</h3>
          <button
            v-if="hasUnread"
            @click="handleMarkAllAsRead"
            class="text-xs text-indigo-600 hover:text-indigo-700 font-medium active-scale"
          >
            Tandai Semua Dibaca
          </button>
        </div>
      </div>

      <!-- Notifications List -->
      <div class="overflow-y-auto" style="max-height: 360px">
        <!-- Loading State -->
        <div v-if="isLoading" class="p-4 text-center text-gray-500">
          <div class="inline-block w-5 h-5 border-2 border-gray-300 border-t-indigo-600 rounded-full animate-spin"></div>
          <p class="mt-2 text-sm">Memuat notifikasi...</p>
        </div>

        <!-- Empty State -->
        <div v-else-if="recentNotifications.length === 0" class="p-8 text-center">
          <svg class="w-12 h-12 mx-auto text-gray-300 mb-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9" />
          </svg>
          <p class="text-sm text-gray-500">Tidak ada notifikasi</p>
        </div>

        <!-- Notification Items dengan staggered animation -->
        <Motion
          v-for="(notification, index) in recentNotifications"
          :key="notification.id"
          :initial="{ opacity: 0, x: -10 }"
          :animate="{ opacity: 1, x: 0 }"
          :transition="{ duration: 0.2, delay: index * 0.05, ease: 'easeOut' }"
          @click="handleNotificationClick(notification)"
          class="px-4 py-3 border-b border-gray-100 hover:bg-gray-50 cursor-pointer transition-colors duration-150 active-scale"
          :class="{ 'bg-indigo-50/30': !notification.is_read }"
        >
          <div class="flex items-start gap-3">
            <!-- Icon Badge -->
            <div
              class="flex-shrink-0 w-8 h-8 rounded-full flex items-center justify-center"
              :class="getNotificationBgClass(notification.type)"
            >
              <svg class="w-4 h-4" :class="getNotificationTextClass(notification.type)" fill="currentColor" viewBox="0 0 20 20">
                <path v-if="notification.type === 'SUCCESS'" fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
                <path v-else-if="notification.type === 'ERROR'" fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd" />
                <path v-else-if="notification.type === 'WARNING'" fill-rule="evenodd" d="M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.213 2.98-1.742 2.98H4.42c-1.53 0-2.493-1.646-1.743-2.98l5.58-9.92zM11 13a1 1 0 11-2 0 1 1 0 012 0zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z" clip-rule="evenodd" />
                <path v-else-if="notification.type === 'ACHIEVEMENT'" d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z" />
                <path v-else fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z" clip-rule="evenodd" />
              </svg>
            </div>

            <!-- Content -->
            <div class="flex-1 min-w-0">
              <p class="text-sm font-medium text-gray-900 truncate">
                {{ notification.title }}
              </p>
              <p class="text-xs text-gray-600 line-clamp-2 mt-0.5">
                {{ notification.message }}
              </p>
              <p class="text-xs text-gray-400 mt-1">
                {{ formatTimeAgo(notification.created_at) }}
              </p>
            </div>

            <!-- Unread Indicator -->
            <div v-if="!notification.is_read" class="flex-shrink-0 w-2 h-2 bg-indigo-600 rounded-full"></div>
          </div>
        </Motion>
      </div>

      <!-- Footer -->
      <div class="px-4 py-3 border-t border-gray-100 bg-gray-50">
        <button
          @click="navigateToNotificationCenter"
          class="w-full text-center text-sm text-indigo-600 hover:text-indigo-700 font-medium active-scale"
        >
          Lihat Semua Notifikasi
        </button>
      </div>
    </Motion>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { Motion } from 'motion-v'
import { useNotificationStore } from '@/stores/notification'
import { formatDistanceToNow } from 'date-fns'
import { id as idLocale } from 'date-fns/locale'

const router = useRouter()
const notificationStore = useNotificationStore()

// State
const isOpen = ref(false)
const isLoading = ref(false)
const recentNotifications = ref([])

// Computed
const unreadCount = computed(() => notificationStore.unreadCount)
const hasUnread = computed(() => notificationStore.hasUnread)

/**
 * toggleDropdown membuka/menutup dropdown notification
 * dengan fetch recent notifications saat dibuka
 */
const toggleDropdown = async () => {
  isOpen.value = !isOpen.value
  
  if (isOpen.value) {
    await fetchRecentNotifications()
    
    // Haptic feedback untuk mobile
    if ('vibrate' in navigator) {
      navigator.vibrate(10)
    }
  }
}

/**
 * closeDropdown menutup dropdown
 */
const closeDropdown = () => {
  isOpen.value = false
}

/**
 * fetchRecentNotifications mengambil 5 notifikasi terbaru
 * untuk display di dropdown
 */
const fetchRecentNotifications = async () => {
  try {
    isLoading.value = true
    const notifications = await notificationStore.getRecentNotifications(5)
    recentNotifications.value = notifications
  } catch (error) {
    console.error('Error fetching recent notifications:', error)
  } finally {
    isLoading.value = false
  }
}

/**
 * handleNotificationClick menangani klik pada notifikasi
 * dengan mark as read dan haptic feedback
 */
const handleNotificationClick = async (notification) => {
  if (!notification.is_read) {
    await notificationStore.markAsRead(notification.id)
  }
  
  // Haptic feedback
  if ('vibrate' in navigator) {
    navigator.vibrate(10)
  }
  
  closeDropdown()
  router.push('/notifications')
}

/**
 * handleMarkAllAsRead menandai semua notifikasi sebagai dibaca
 * dengan haptic feedback
 */
const handleMarkAllAsRead = async () => {
  try {
    await notificationStore.markAllAsRead()
    await fetchRecentNotifications()
    
    // Haptic feedback
    if ('vibrate' in navigator) {
      navigator.vibrate([10, 50, 10])
    }
  } catch (error) {
    console.error('Error marking all as read:', error)
  }
}

/**
 * navigateToNotificationCenter redirect ke notification center
 */
const navigateToNotificationCenter = () => {
  closeDropdown()
  router.push('/notifications')
}

/**
 * formatTimeAgo format timestamp menjadi relative time
 * menggunakan date-fns dengan Indonesian locale
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

// Lifecycle: Start polling saat component mounted
onMounted(() => {
  notificationStore.startPolling(30000) // Poll setiap 30 detik
})

// Lifecycle: Stop polling saat component unmounted
onUnmounted(() => {
  notificationStore.stopPolling()
})
</script>
