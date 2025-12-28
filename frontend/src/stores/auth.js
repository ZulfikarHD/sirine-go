import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

/**
 * Auth Store merupakan state management untuk authentication
 * yang mencakup user data, token management, dan authentication status
 */
export const useAuthStore = defineStore('auth', () => {
  // State
  const user = ref(null)
  const token = ref(localStorage.getItem('auth_token') || null)
  const refreshToken = ref(localStorage.getItem('refresh_token') || null)

  // Computed
  const isAuthenticated = computed(() => !!token.value && !!user.value)
  const userRole = computed(() => user.value?.role || null)
  const userDepartment = computed(() => user.value?.department || null)
  const requirePasswordChange = computed(() => user.value?.must_change_password || false)

  /**
   * Set authentication data dengan token dan user info
   * untuk menyimpan session setelah login berhasil
   */
  const setAuth = (authData) => {
    token.value = authData.token
    refreshToken.value = authData.refresh_token
    user.value = authData.user

    // Persist ke localStorage
    localStorage.setItem('auth_token', authData.token)
    localStorage.setItem('refresh_token', authData.refresh_token)
    localStorage.setItem('user_data', JSON.stringify(authData.user))
  }

  /**
   * Update user data tanpa mengubah token
   * untuk sync profile changes
   */
  const setUser = (userData) => {
    user.value = userData
    localStorage.setItem('user_data', JSON.stringify(userData))
  }

  /**
   * Clear authentication data dan logout user
   * untuk menghapus session setelah logout
   */
  const clearAuth = () => {
    user.value = null
    token.value = null
    refreshToken.value = null

    // Clear localStorage
    localStorage.removeItem('auth_token')
    localStorage.removeItem('refresh_token')
    localStorage.removeItem('user_data')
  }

  /**
   * Restore authentication dari localStorage
   * untuk persistent login setelah page reload
   */
  const restoreAuth = () => {
    const storedToken = localStorage.getItem('auth_token')
    const storedRefreshToken = localStorage.getItem('refresh_token')
    const storedUser = localStorage.getItem('user_data')

    if (storedToken && storedUser) {
      token.value = storedToken
      refreshToken.value = storedRefreshToken
      try {
        user.value = JSON.parse(storedUser)
      } catch (e) {
        console.error('Error parsing stored user data:', e)
        clearAuth()
      }
    }
  }

  /**
   * Check apakah user memiliki role tertentu
   * untuk authorization checks
   */
  const hasRole = (...roles) => {
    if (!user.value) return false
    return roles.includes(user.value.role)
  }

  /**
   * Check apakah user adalah admin atau manager
   */
  const isAdmin = () => {
    return hasRole('ADMIN', 'MANAGER')
  }

  /**
   * Check apakah user memiliki department tertentu
   */
  const hasDepartment = (...departments) => {
    if (!user.value) return false
    return departments.includes(user.value.department)
  }

  /**
   * Fetch current user data dari server
   * untuk sync latest user information
   */
  const fetchCurrentUser = async () => {
    if (!token.value) return

    try {
      const API_BASE = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api'
      const response = await fetch(`${API_BASE}/auth/me`, {
        headers: {
          'Authorization': `Bearer ${token.value}`
        }
      })

      if (response.ok) {
        const data = await response.json()
        if (data.success) {
          setUser(data.data)
        }
      }
    } catch (error) {
      console.error('Error fetching current user:', error)
    }
  }

  return {
    // State
    user,
    token,
    refreshToken,
    
    // Computed
    isAuthenticated,
    userRole,
    userDepartment,
    requirePasswordChange,
    
    // Actions
    setAuth,
    setUser,
    clearAuth,
    restoreAuth,
    hasRole,
    isAdmin,
    hasDepartment,
    fetchCurrentUser,
  }
})
