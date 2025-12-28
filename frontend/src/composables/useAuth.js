import { ref } from 'vue'
import { useAuthStore } from '../stores/auth'
import { useApi } from './useApi'
import { useRouter } from 'vue-router'

/**
 * useAuth composable untuk authentication operations
 * yang mencakup login, logout, dan user management
 */
export const useAuth = () => {
  const authStore = useAuthStore()
  const api = useApi()
  const router = useRouter()

  const isLoading = ref(false)
  const error = ref(null)

  /**
   * Login dengan NIP dan password
   * @param {string} nip - Nomor Induk Pegawai
   * @param {string} password - Password user
   * @param {boolean} rememberMe - Remember me option
   * @returns {Promise} Login response
   */
  const login = async (nip, password, rememberMe = false) => {
    isLoading.value = true
    error.value = null

    try {
      const response = await api.post('/auth/login', {
        nip,
        password,
        remember_me: rememberMe,
      })

      if (response.success) {
        authStore.setAuth(response.data)

        // Trigger haptic feedback untuk success (jika available)
        triggerHapticFeedback('success')

        return response.data
      } else {
        throw new Error(response.message || 'Login gagal')
      }
    } catch (err) {
      error.value = err.response?.data?.message || err.message || 'Login gagal'
      
      // Trigger haptic feedback untuk error
      triggerHapticFeedback('error')
      
      throw err
    } finally {
      isLoading.value = false
    }
  }

  /**
   * Logout dan clear session
   */
  const logout = async () => {
    isLoading.value = true
    error.value = null

    try {
      await api.post('/auth/logout')
      authStore.clearAuth()
      router.push('/login')
      
      // Trigger haptic feedback
      triggerHapticFeedback('success')
    } catch (err) {
      // Tetap clear auth meski API call gagal
      authStore.clearAuth()
      router.push('/login')
      console.error('Logout error:', err)
    } finally {
      isLoading.value = false
    }
  }

  /**
   * Fetch current user data untuk refresh user info
   */
  const fetchCurrentUser = async () => {
    isLoading.value = true
    error.value = null

    try {
      const response = await api.get('/auth/me')

      if (response.success) {
        authStore.setUser(response.data)
        return response.data
      } else {
        throw new Error(response.message || 'Gagal mengambil data user')
      }
    } catch (err) {
      error.value = err.response?.data?.message || err.message
      throw err
    } finally {
      isLoading.value = false
    }
  }

  /**
   * Check authentication status dan validate token
   */
  const checkAuth = async () => {
    if (!authStore.token) {
      return false
    }

    try {
      await fetchCurrentUser()
      return true
    } catch (err) {
      authStore.clearAuth()
      return false
    }
  }

  /**
   * Refresh authentication token
   */
  const refreshAuth = async () => {
    if (!authStore.refreshToken) {
      throw new Error('No refresh token available')
    }

    try {
      const response = await api.post('/auth/refresh', {
        refresh_token: authStore.refreshToken,
      })

      if (response.success) {
        authStore.setAuth(response.data)
        return response.data
      } else {
        throw new Error(response.message || 'Token refresh gagal')
      }
    } catch (err) {
      authStore.clearAuth()
      throw err
    }
  }

  /**
   * Get role-based dashboard route
   */
  const getDashboardRoute = () => {
    if (!authStore.user) return '/login'

    const role = authStore.user.role

    // Route mapping berdasarkan role
    const dashboardMap = {
      ADMIN: '/dashboard/admin',
      MANAGER: '/dashboard/admin',
      STAFF_KHAZWAL: '/dashboard/staff',
      OPERATOR_CETAK: '/dashboard/staff',
      QC_INSPECTOR: '/dashboard/staff',
      VERIFIKATOR: '/dashboard/staff',
      STAFF_KHAZKHIR: '/dashboard/staff',
    }

    return dashboardMap[role] || '/dashboard/staff'
  }

  /**
   * Forgot password - request reset link via email
   * @param {string} nipOrEmail - NIP atau Email user
   */
  const forgotPassword = async (nipOrEmail) => {
    isLoading.value = true
    error.value = null

    try {
      const response = await api.post('/auth/forgot-password', {
        nip_or_email: nipOrEmail,
      })

      if (response.success) {
        return response
      } else {
        throw new Error(response.message || 'Gagal mengirim email reset password')
      }
    } catch (err) {
      error.value = err.response?.data?.message || err.message
      throw err
    } finally {
      isLoading.value = false
    }
  }

  /**
   * Reset password menggunakan token dari email
   * @param {string} token - Reset token dari email
   * @param {string} newPassword - Password baru
   */
  const resetPassword = async (token, newPassword) => {
    isLoading.value = true
    error.value = null

    try {
      const response = await api.post('/auth/reset-password', {
        token,
        new_password: newPassword,
      })

      if (response.success) {
        triggerHapticFeedback('success')
        return response
      } else {
        throw new Error(response.message || 'Gagal reset password')
      }
    } catch (err) {
      error.value = err.response?.data?.message || err.message
      triggerHapticFeedback('error')
      throw err
    } finally {
      isLoading.value = false
    }
  }

  /**
   * Change password untuk user yang sedang login
   * @param {string} currentPassword - Password saat ini
   * @param {string} newPassword - Password baru
   */
  const changePassword = async (currentPassword, newPassword) => {
    isLoading.value = true
    error.value = null

    try {
      const response = await api.put('/profile/password', {
        current_password: currentPassword,
        new_password: newPassword,
      })

      if (response.success) {
        triggerHapticFeedback('success')
        
        // Auto logout setelah change password
        setTimeout(() => {
          authStore.clearAuth()
          router.push('/login')
        }, 2000)
        
        return response
      } else {
        throw new Error(response.message || 'Gagal mengubah password')
      }
    } catch (err) {
      error.value = err.response?.data?.message || err.message
      triggerHapticFeedback('error')
      throw err
    } finally {
      isLoading.value = false
    }
  }

  /**
   * Trigger haptic feedback (vibration) jika available
   * @param {string} type - Type of feedback: 'success', 'error', 'warning'
   */
  const triggerHapticFeedback = (type = 'success') => {
    if ('vibrate' in navigator) {
      switch (type) {
        case 'success':
          navigator.vibrate(200) // Single pulse untuk success
          break
        case 'error':
          navigator.vibrate([100, 50, 100]) // Double pulse untuk error
          break
        case 'warning':
          navigator.vibrate([50, 50, 50]) // Triple short pulse
          break
        case 'achievement':
          navigator.vibrate([100, 50, 100, 50, 100]) // Multiple pulses
          break
        default:
          navigator.vibrate(100)
      }
    }
  }

  return {
    // State
    isLoading,
    error,

    // Methods
    login,
    logout,
    fetchCurrentUser,
    checkAuth,
    refreshAuth,
    getDashboardRoute,
    triggerHapticFeedback,
    forgotPassword,
    resetPassword,
    changePassword,

    // Store access
    isAuthenticated: authStore.isAuthenticated,
    user: authStore.user,
    hasRole: authStore.hasRole,
    isAdmin: authStore.isAdmin,
  }
}
