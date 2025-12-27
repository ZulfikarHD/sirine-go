import axios from 'axios'
import { ref } from 'vue'

const api = axios.create({
  baseURL: '/api',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json'
  }
})

// Request interceptor
api.interceptors.request.use(
  (config) => {
    // Tambahkan token jika ada
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// Response interceptor
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      // Handle unauthorized
      localStorage.removeItem('token')
      window.location.href = '/login'
    }
    return Promise.reject(error)
  }
)

export function useApi() {
  const loading = ref(false)
  const error = ref(null)

  const request = async (method, url, data = null) => {
    loading.value = true
    error.value = null

    try {
      const response = await api({
        method,
        url,
        data
      })
      return response.data
    } catch (err) {
      error.value = err.response?.data?.error || 'Terjadi kesalahan'
      throw err
    } finally {
      loading.value = false
    }
  }

  const get = (url) => request('GET', url)
  const post = (url, data) => request('POST', url, data)
  const put = (url, data) => request('PUT', url, data)
  const del = (url) => request('DELETE', url)

  return {
    loading,
    error,
    get,
    post,
    put,
    del
  }
}

export default api
