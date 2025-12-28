import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { useApi } from '@/composables/useApi'

/**
 * Notification Store merupakan Pinia store untuk state management notifikasi
 * yang bertujuan untuk handling in-app notifications dengan real-time updates
 */
export const useNotificationStore = defineStore('notification', () => {
  const api = useApi()

  // State
  const notifications = ref([])
  const unreadCount = ref(0)
  const isLoading = ref(false)
  const pollingInterval = ref(null)

  // Computed
  const unreadNotifications = computed(() => 
    notifications.value.filter(n => !n.is_read)
  )

  const hasUnread = computed(() => unreadCount.value > 0)

  /**
   * fetchNotifications mengambil semua notifikasi user
   * dengan optional filter unread_only untuk performance
   */
  const fetchNotifications = async (unreadOnly = false) => {
    try {
      isLoading.value = true
      const params = unreadOnly ? '?unread_only=true' : ''
      const response = await api.get(`/notifications${params}`)
      
      if (response.success) {
        notifications.value = response.data || []
        return response.data
      }
    } catch (error) {
      console.error('Error fetching notifications:', error)
      throw error
    } finally {
      isLoading.value = false
    }
  }

  /**
   * fetchUnreadCount mengambil jumlah notifikasi belum dibaca
   * untuk badge display di notification bell
   */
  const fetchUnreadCount = async () => {
    try {
      const response = await api.get('/notifications/unread-count')
      if (response.success) {
        unreadCount.value = response.data.count
        return response.data.count
      }
    } catch (error) {
      console.error('Error fetching unread count:', error)
      throw error
    }
  }

  /**
   * getRecentNotifications mengambil N notifikasi terbaru
   * untuk quick preview di notification bell dropdown
   */
  const getRecentNotifications = async (limit = 5) => {
    try {
      const response = await api.get(`/notifications/recent?limit=${limit}`)
      if (response.success) {
        return response.data || []
      }
    } catch (error) {
      console.error('Error fetching recent notifications:', error)
      throw error
    }
  }

  /**
   * markAsRead menandai satu notifikasi sebagai sudah dibaca
   * dengan optimistic update untuk UX yang lebih baik
   */
  const markAsRead = async (notificationId) => {
    try {
      // Optimistic update
      const notification = notifications.value.find(n => n.id === notificationId)
      if (notification && !notification.is_read) {
        notification.is_read = true
        notification.read_at = new Date().toISOString()
        unreadCount.value = Math.max(0, unreadCount.value - 1)
      }

      const response = await api.put(`/notifications/${notificationId}/read`)
      if (response.success) {
        return true
      }
    } catch (error) {
      // Rollback optimistic update on error
      const notification = notifications.value.find(n => n.id === notificationId)
      if (notification) {
        notification.is_read = false
        notification.read_at = null
        unreadCount.value++
      }
      console.error('Error marking notification as read:', error)
      throw error
    }
  }

  /**
   * markAllAsRead menandai semua notifikasi sebagai sudah dibaca
   * untuk bulk action dengan optimistic update
   */
  const markAllAsRead = async () => {
    try {
      // Store original state untuk rollback
      const originalNotifications = JSON.parse(JSON.stringify(notifications.value))
      const originalCount = unreadCount.value

      // Optimistic update
      notifications.value.forEach(notification => {
        if (!notification.is_read) {
          notification.is_read = true
          notification.read_at = new Date().toISOString()
        }
      })
      unreadCount.value = 0

      const response = await api.put('/notifications/read-all')
      if (response.success) {
        return true
      }
    } catch (error) {
      // Rollback on error
      await fetchNotifications()
      await fetchUnreadCount()
      console.error('Error marking all as read:', error)
      throw error
    }
  }

  /**
   * deleteNotification menghapus notifikasi
   * dengan optimistic update
   */
  const deleteNotification = async (notificationId) => {
    try {
      // Store original untuk rollback
      const notificationIndex = notifications.value.findIndex(n => n.id === notificationId)
      const deletedNotification = notifications.value[notificationIndex]

      // Optimistic update
      notifications.value.splice(notificationIndex, 1)
      if (!deletedNotification.is_read) {
        unreadCount.value = Math.max(0, unreadCount.value - 1)
      }

      const response = await api.del(`/notifications/${notificationId}`)
      if (response.success) {
        return true
      }
    } catch (error) {
      // Rollback
      await fetchNotifications()
      await fetchUnreadCount()
      console.error('Error deleting notification:', error)
      throw error
    }
  }

  /**
   * startPolling memulai polling untuk real-time updates
   * dengan interval 30 detik (configurable)
   */
  const startPolling = (intervalMs = 30000) => {
    if (pollingInterval.value) {
      stopPolling()
    }

    // Initial fetch
    fetchUnreadCount()

    // Setup polling interval
    pollingInterval.value = setInterval(() => {
      fetchUnreadCount()
    }, intervalMs)
  }

  /**
   * stopPolling menghentikan polling
   * dipanggil saat user logout atau component unmount
   */
  const stopPolling = () => {
    if (pollingInterval.value) {
      clearInterval(pollingInterval.value)
      pollingInterval.value = null
    }
  }

  /**
   * reset membersihkan state notification
   * dipanggil saat logout
   */
  const reset = () => {
    notifications.value = []
    unreadCount.value = 0
    stopPolling()
  }

  /**
   * getNotificationIcon mengembalikan icon berdasarkan notification type
   * untuk display di UI
   */
  const getNotificationIcon = (type) => {
    const icons = {
      INFO: 'information-circle',
      SUCCESS: 'checkmark-circle',
      WARNING: 'warning',
      ERROR: 'close-circle',
      ACHIEVEMENT: 'trophy',
    }
    return icons[type] || 'information-circle'
  }

  /**
   * getNotificationColor mengembalikan color berdasarkan notification type
   * untuk styling badge dan icon
   */
  const getNotificationColor = (type) => {
    const colors = {
      INFO: 'blue',
      SUCCESS: 'green',
      WARNING: 'yellow',
      ERROR: 'red',
      ACHIEVEMENT: 'purple',
    }
    return colors[type] || 'gray'
  }

  return {
    // State
    notifications,
    unreadCount,
    isLoading,
    
    // Computed
    unreadNotifications,
    hasUnread,
    
    // Actions
    fetchNotifications,
    fetchUnreadCount,
    getRecentNotifications,
    markAsRead,
    markAllAsRead,
    deleteNotification,
    startPolling,
    stopPolling,
    reset,
    
    // Helpers
    getNotificationIcon,
    getNotificationColor,
  }
})
