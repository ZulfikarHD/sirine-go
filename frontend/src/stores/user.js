import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import apiClient from '../composables/useApi'

/**
 * User Store merupakan state management untuk user management operations
 * yang mencakup CRUD users, filters, pagination, dan search functionality
 */
export const useUserStore = defineStore('user', () => {
  // State
  const users = ref([])
  const currentEditUser = ref(null)
  const total = ref(0)
  const currentPage = ref(1)
  const perPage = ref(20)
  const totalPages = ref(0)
  const loading = ref(false)
  const error = ref(null)

  // Filters
  const filters = ref({
    role: '',
    department: '',
    status: '',
    search: ''
  })

  /**
   * Fetch users dengan filters dan pagination
   */
  const fetchUsers = async (page = 1) => {
    loading.value = true
    error.value = null

    try {
      const params = {
        page,
        per_page: perPage.value,
        ...filters.value
      }

      // Remove empty filter values
      Object.keys(params).forEach(key => {
        if (params[key] === '' || params[key] === null || params[key] === undefined) {
          delete params[key]
        }
      })

      const response = await apiClient.get('/users', { params })

      if (response.data.success) {
        users.value = response.data.data.users || []
        total.value = response.data.data.total
        currentPage.value = response.data.data.page
        perPage.value = response.data.data.per_page
        totalPages.value = response.data.data.total_pages
      }
    } catch (err) {
      error.value = err.response?.data?.error || 'Gagal mengambil data users'
      throw err
    } finally {
      loading.value = false
    }
  }

  /**
   * Get user by ID
   */
  const getUserById = async (id) => {
    loading.value = true
    error.value = null

    try {
      const response = await apiClient.get(`/users/${id}`)

      if (response.data.success) {
        return response.data.data
      }
    } catch (err) {
      error.value = err.response?.data?.error || 'Gagal mengambil data user'
      throw err
    } finally {
      loading.value = false
    }
  }

  /**
   * Create new user dengan auto-generated password
   */
  const createUser = async (userData) => {
    loading.value = true
    error.value = null

    try {
      const response = await apiClient.post('/users', userData)

      if (response.data.success) {
        // Refresh users list
        await fetchUsers(currentPage.value)
        return response.data.data
      }
    } catch (err) {
      error.value = err.response?.data?.error || 'Gagal membuat user baru'
      throw err
    } finally {
      loading.value = false
    }
  }

  /**
   * Update user data (Admin)
   */
  const updateUser = async (id, userData) => {
    loading.value = true
    error.value = null

    try {
      const response = await apiClient.put(`/users/${id}`, userData)

      if (response.data.success) {
        // Refresh users list
        await fetchUsers(currentPage.value)
        return response.data.data
      }
    } catch (err) {
      error.value = err.response?.data?.error || 'Gagal update user'
      throw err
    } finally {
      loading.value = false
    }
  }

  /**
   * Delete user (soft delete)
   */
  const deleteUser = async (id) => {
    loading.value = true
    error.value = null

    try {
      const response = await apiClient.delete(`/users/${id}`)

      if (response.data.success) {
        // Refresh users list
        await fetchUsers(currentPage.value)
        return true
      }
    } catch (err) {
      error.value = err.response?.data?.error || 'Gagal menghapus user'
      throw err
    } finally {
      loading.value = false
    }
  }

  /**
   * Search users berdasarkan query
   */
  const searchUsers = async (query) => {
    if (!query || query.trim() === '') {
      return []
    }

    loading.value = true
    error.value = null

    try {
      const response = await apiClient.get('/users/search', {
        params: { q: query }
      })

      if (response.data.success) {
        return response.data.data || []
      }
    } catch (err) {
      error.value = err.response?.data?.error || 'Gagal search users'
      throw err
    } finally {
      loading.value = false
    }
  }

  /**
   * Bulk delete users
   */
  const bulkDeleteUsers = async (userIds) => {
    loading.value = true
    error.value = null

    try {
      const response = await apiClient.post('/users/bulk-delete', {
        user_ids: userIds
      })

      if (response.data.success) {
        // Refresh users list
        await fetchUsers(currentPage.value)
        return response.data.data.affected_count
      }
    } catch (err) {
      error.value = err.response?.data?.error || 'Gagal bulk delete users'
      throw err
    } finally {
      loading.value = false
    }
  }

  /**
   * Bulk update status users
   */
  const bulkUpdateStatus = async (userIds, status) => {
    loading.value = true
    error.value = null

    try {
      const response = await apiClient.post('/users/bulk-update-status', {
        user_ids: userIds,
        status
      })

      if (response.data.success) {
        // Refresh users list
        await fetchUsers(currentPage.value)
        return response.data.data.affected_count
      }
    } catch (err) {
      error.value = err.response?.data?.error || 'Gagal bulk update status'
      throw err
    } finally {
      loading.value = false
    }
  }

  /**
   * Set filters dan trigger fetch
   */
  const setFilters = (newFilters) => {
    filters.value = { ...filters.value, ...newFilters }
    fetchUsers(1) // Reset ke page 1 saat filter berubah
  }

  /**
   * Clear filters
   */
  const clearFilters = () => {
    filters.value = {
      role: '',
      department: '',
      status: '',
      search: ''
    }
    fetchUsers(1)
  }

  /**
   * Set current edit user
   */
  const setCurrentEditUser = (user) => {
    currentEditUser.value = user
  }

  /**
   * Clear current edit user
   */
  const clearCurrentEditUser = () => {
    currentEditUser.value = null
  }

  // Computed
  const hasUsers = computed(() => users.value.length > 0)
  const isEmpty = computed(() => !loading.value && users.value.length === 0)
  const hasFilters = computed(() => {
    return filters.value.role !== '' || 
           filters.value.department !== '' || 
           filters.value.status !== '' || 
           filters.value.search !== ''
  })

  return {
    // State
    users,
    currentEditUser,
    total,
    currentPage,
    perPage,
    totalPages,
    loading,
    error,
    filters,

    // Actions
    fetchUsers,
    getUserById,
    createUser,
    updateUser,
    deleteUser,
    searchUsers,
    bulkDeleteUsers,
    bulkUpdateStatus,
    setFilters,
    clearFilters,
    setCurrentEditUser,
    clearCurrentEditUser,

    // Computed
    hasUsers,
    isEmpty,
    hasFilters
  }
})
