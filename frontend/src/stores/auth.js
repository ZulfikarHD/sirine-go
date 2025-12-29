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
   * Validate JWT token format dan struktur
   * untuk memastikan token valid sebelum digunakan
   * Note: Tidak check expiry untuk access token, biarkan interceptor handle dengan refresh token
   */
  const isValidTokenFormat = (token, checkExpiry = false) => {
    if (!token || typeof token !== 'string') return false
    
    // JWT format: header.payload.signature
    const parts = token.split('.')
    if (parts.length !== 3) return false
    
    try {
      // Decode payload untuk validate structure
      const payload = JSON.parse(atob(parts[1]))
      
      // Check expiry hanya jika diminta (untuk refresh token validation)
      if (checkExpiry && payload.exp && payload.exp * 1000 < Date.now()) {
        console.warn('Token sudah expired')
        return false
      }
      
      return true
    } catch (e) {
      console.error('Token format tidak valid:', e)
      return false
    }
  }

  /**
   * Restore authentication dari localStorage
   * untuk persistent login setelah page reload
   * 
   * Strategy: Restore semua tokens jika format valid, tanpa check expiry access token.
   * Biarkan API interceptor handle expired access token dengan refresh token automatically.
   */
  const restoreAuth = () => {
    const storedToken = localStorage.getItem('auth_token')
    const storedRefreshToken = localStorage.getItem('refresh_token')
    const storedUser = localStorage.getItem('user_data')

    // Butuh minimal refresh token dan user data untuk restore session
    if (storedRefreshToken && storedUser) {
      // Validate refresh token format dan expiry
      if (!isValidTokenFormat(storedRefreshToken, true)) {
        console.warn('Refresh token tidak valid atau expired, clearing auth')
        clearAuth()
        return
      }

      // Restore semua data (access token boleh expired, akan di-refresh oleh interceptor)
      token.value = storedToken
      refreshToken.value = storedRefreshToken
      
      try {
        user.value = JSON.parse(storedUser)
      } catch (e) {
        console.error('Error parsing stored user data:', e)
        clearAuth()
      }
    } else if (storedToken || storedRefreshToken || storedUser) {
      // Jika ada incomplete auth data (missing refresh token atau user), clear semua
      console.warn('Incomplete auth data detected (missing refresh token or user), clearing all')
      clearAuth()
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

  /**
   * Update specific user field tanpa full replace
   * untuk partial updates seperti photo upload
   */
  const updateUserField = (field, value) => {
    if (!user.value) return
    
    user.value = {
      ...user.value,
      [field]: value
    }
    
    localStorage.setItem('user_data', JSON.stringify(user.value))
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
    updateUserField,
    isValidTokenFormat,
  }
})
