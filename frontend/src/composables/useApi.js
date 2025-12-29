import axios from 'axios'
import { useAuthStore } from '../stores/auth'
import router from '../router'

/**
 * API base URL dari environment variable
 * dengan fallback ke localhost untuk development
 */
const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api'

/**
 * Flag untuk prevent multiple simultaneous refresh token attempts
 */
let isRefreshing = false
let failedQueue = []

/**
 * Process queue untuk pending requests setelah token refresh
 */
const processQueue = (error, token = null) => {
  failedQueue.forEach(prom => {
    if (error) {
      prom.reject(error)
    } else {
      prom.resolve(token)
    }
  })
  failedQueue = []
}

/**
 * Axios instance dengan auto token injection dan interceptors
 * untuk handling authentication dan error responses
 */
const apiClient = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
})

// Request interceptor untuk inject token
apiClient.interceptors.request.use(
  (config) => {
    const authStore = useAuthStore()
    if (authStore.token) {
      config.headers.Authorization = `Bearer ${authStore.token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// Response interceptor untuk handle errors
apiClient.interceptors.response.use(
  (response) => {
    return response
  },
  async (error) => {
    const authStore = useAuthStore()
    const originalRequest = error.config

    // Handle 401 Unauthorized (hanya untuk authenticated requests)
    if (error.response?.status === 401 && !originalRequest._retry) {
      // Skip auto-redirect jika ini adalah login/refresh request yang gagal
      const isLoginRequest = originalRequest.url?.includes('/auth/login')
      const isRefreshRequest = originalRequest.url?.includes('/auth/refresh')
      
      if (isLoginRequest) {
        // Untuk login request, reject error tanpa redirect
        return Promise.reject(error)
      }

      if (isRefreshRequest) {
        // Refresh token sendiri yang gagal, logout immediately
        console.error('API Error: refresh token tidak valid atau expired')
        isRefreshing = false
        processQueue(error, null)
        authStore.clearAuth()
        if (window.location.pathname !== '/login') {
          router.push('/login').catch(() => {})
        }
        return Promise.reject(error)
      }

      originalRequest._retry = true

      // Jika sedang refresh, queue request ini
      if (isRefreshing) {
        return new Promise((resolve, reject) => {
          failedQueue.push({ resolve, reject })
        }).then(token => {
          originalRequest.headers.Authorization = `Bearer ${token}`
          return apiClient(originalRequest)
        }).catch(err => {
          return Promise.reject(err)
        })
      }

      // Coba refresh token untuk request lainnya
      if (authStore.refreshToken) {
        isRefreshing = true

        try {
          const response = await axios.post(`${API_BASE_URL}/auth/refresh`, {
            refresh_token: authStore.refreshToken,
          })

          if (response.data.success) {
            const newToken = response.data.data.token
            authStore.setAuth(response.data.data)
            originalRequest.headers.Authorization = `Bearer ${newToken}`
            
            processQueue(null, newToken)
            isRefreshing = false
            
            return apiClient(originalRequest)
          }
        } catch (refreshError) {
          // Refresh token gagal, logout user
          console.error('API Error: gagal refresh token -', refreshError.response?.data?.message || 'Token expired')
          isRefreshing = false
          processQueue(refreshError, null)
          authStore.clearAuth()
          
          // Use router untuk avoid full page reload
          if (window.location.pathname !== '/login') {
            router.push('/login').catch(() => {})
          }
          return Promise.reject(refreshError)
        }
      } else {
        // Tidak ada refresh token, redirect ke login
        console.warn('API Error: tidak ada refresh token, redirect ke login')
        authStore.clearAuth()
        if (window.location.pathname !== '/login') {
          router.push('/login').catch(() => {})
        }
      }
    }

    return Promise.reject(error)
  }
)

/**
 * useApi composable untuk API calls dengan built-in error handling
 */
export const useApi = () => {
  const authStore = useAuthStore()

  /**
   * Generic GET request
   */
  const get = async (url, config = {}) => {
    try {
      const response = await apiClient.get(url, config)
      return response.data
    } catch (error) {
      handleApiError(error)
      throw error
    }
  }

  /**
   * Generic POST request
   */
  const post = async (url, data = {}, config = {}) => {
    try {
      const response = await apiClient.post(url, data, config)
      return response.data
    } catch (error) {
      handleApiError(error)
      throw error
    }
  }

  /**
   * Generic PUT request
   */
  const put = async (url, data = {}, config = {}) => {
    try {
      const response = await apiClient.put(url, data, config)
      return response.data
    } catch (error) {
      handleApiError(error)
      throw error
    }
  }

  /**
   * Generic DELETE request
   */
  const del = async (url, config = {}) => {
    try {
      const response = await apiClient.delete(url, config)
      return response.data
    } catch (error) {
      handleApiError(error)
      throw error
    }
  }

  /**
   * Handle API errors dengan user-friendly messages
   */
  const handleApiError = (error) => {
    if (error.response) {
      // Server responded dengan error status
      const message = error.response.data?.message || 'Terjadi kesalahan pada server'
      console.error('API Error:', message)
    } else if (error.request) {
      // Request dibuat tapi tidak ada response
      console.error('Network Error: Tidak dapat terhubung ke server')
    } else {
      // Error lain
      console.error('Error:', error.message)
    }
  }

  return {
    get,
    post,
    put,
    del,
    apiClient,
  }
}

export default apiClient
