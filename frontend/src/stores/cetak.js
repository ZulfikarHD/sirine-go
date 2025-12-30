/**
 * Cetak Store - Pinia store untuk Cetak Queue Management
 * yang mencakup queue retrieval, detail view, dan filtering
 */
import { defineStore } from 'pinia'
import apiClient from '@/composables/useApi'

export const useCetakStore = defineStore('cetak', {
  state: () => ({
    // Queue state
    queue: [],
    queueLoading: false,
    queueError: null,
    queuePagination: {
      total: 0,
      page: 1,
      per_page: 20,
      total_pages: 0
    },

    // Current detail state
    currentDetail: null,
    detailLoading: false,
    detailError: null
  }),

  getters: {
    /**
     * Get queue items yang urgent (priority URGENT)
     */
    urgentQueue: (state) => {
      return state.queue.filter(item => item.priority === 'URGENT')
    },

    /**
     * Get queue items yang past due
     */
    pastDueQueue: (state) => {
      return state.queue.filter(item => item.is_past_due === true)
    },

    /**
     * Get queue items yang due within 3 days
     */
    soonDueQueue: (state) => {
      return state.queue.filter(item => 
        !item.is_past_due && item.days_until_due <= 3 && item.days_until_due >= 0
      )
    },

    /**
     * Check apakah ada urgent items di queue
     */
    hasUrgentItems: (state) => {
      return state.queue.some(item => item.priority === 'URGENT')
    },

    /**
     * Check apakah ada past due items
     */
    hasPastDueItems: (state) => {
      return state.queue.some(item => item.is_past_due === true)
    },

    /**
     * Get total items count
     */
    totalItems: (state) => {
      return state.queuePagination.total
    }
  },

  actions: {
    /**
     * Get Cetak Queue dengan filter dan pagination
     * @param {Object} filters - Filter options (search, priority, page, per_page)
     * @returns {Promise<Object>} Queue response dengan pagination
     */
    async getCetakQueue(filters = {}) {
      this.queueLoading = true
      this.queueError = null

      try {
        const response = await apiClient.get('/cetak/queue', {
          params: {
            search: filters.search || '',
            priority: filters.priority || '',
            page: filters.page || 1,
            per_page: filters.per_page || 20
          }
        })

        if (response.data.success) {
          this.queue = response.data.data.items || []
          this.queuePagination = {
            total: response.data.data.total || 0,
            page: response.data.data.page || 1,
            per_page: response.data.data.per_page || 20,
            total_pages: response.data.data.total_pages || 0
          }
          return response.data.data
        } else {
          throw new Error(response.data.message || 'Gagal mengambil queue cetak')
        }
      } catch (error) {
        this.queueError = error.response?.data?.message || error.message || 'Gagal mengambil queue cetak'
        throw error
      } finally {
        this.queueLoading = false
      }
    },

    /**
     * Get Cetak Detail untuk single PO
     * @param {Number} poId - Production Order ID
     * @returns {Promise<Object>} PO detail dengan OBC Master dan material prep data
     */
    async getCetakDetail(poId) {
      this.detailLoading = true
      this.detailError = null

      try {
        const response = await apiClient.get(`/cetak/queue/${poId}`)

        if (response.data.success) {
          this.currentDetail = response.data.data
          return response.data.data
        } else {
          throw new Error(response.data.message || 'Gagal mengambil detail PO')
        }
      } catch (error) {
        this.detailError = error.response?.data?.message || error.message || 'Gagal mengambil detail PO'
        throw error
      } finally {
        this.detailLoading = false
      }
    },

    /**
     * Clear queue data
     */
    clearQueue() {
      this.queue = []
      this.queuePagination = {
        total: 0,
        page: 1,
        per_page: 20,
        total_pages: 0
      }
      this.queueError = null
    },

    /**
     * Clear current detail
     */
    clearDetail() {
      this.currentDetail = null
      this.detailError = null
    },

    /**
     * Clear all state (untuk logout atau reset)
     */
    clearAll() {
      this.queue = []
      this.currentDetail = null
      this.queuePagination = {
        total: 0,
        page: 1,
        per_page: 20,
        total_pages: 0
      }
      this.queueError = null
      this.detailError = null
    }
  }
})
