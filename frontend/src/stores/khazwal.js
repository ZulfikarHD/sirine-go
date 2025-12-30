/**
 * Khazwal Store - Pinia store untuk Khazwal Material Preparation
 * yang mencakup queue management, material prep workflow, dan history
 */
import { defineStore } from 'pinia'
import apiClient from '@/composables/useApi'

export const useKhazwalStore = defineStore('khazwal', {
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

    // Current material prep state
    currentPrep: null,
    currentPrepLoading: false,
    currentPrepError: null,

    // History state
    history: [],
    historyLoading: false,
    historyError: null,
    historyPagination: {
      total: 0,
      page: 1,
      per_page: 20,
      total_pages: 0
    },

    // Monitoring state (for supervisor)
    monitoringStats: null,
    monitoringLoading: false,
    monitoringError: null
  }),

  getters: {
    /**
     * Get queue items yang status PENDING
     */
    pendingQueue: (state) => {
      return state.queue.filter(item => 
        item.current_status === 'WAITING_MATERIAL_PREP'
      )
    },

    /**
     * Get queue items yang status IN_PROGRESS
     */
    inProgressQueue: (state) => {
      return state.queue.filter(item => 
        item.current_status === 'MATERIAL_PREP_IN_PROGRESS'
      )
    },

    /**
     * Check apakah current prep sudah siap untuk finalize
     * yaitu: plat, kertas, dan tinta sudah dikonfirmasi
     */
    canFinalize: (state) => {
      if (!state.currentPrep || !state.currentPrep.khazwal_material_prep) {
        return false
      }

      const prep = state.currentPrep.khazwal_material_prep
      return (
        prep.plat_retrieved_at !== null &&
        prep.kertas_blanko_actual !== null &&
        prep.tinta_actual !== null
      )
    },

    /**
     * Check apakah plat sudah dikonfirmasi dan match
     */
    isPlatConfirmed: (state) => {
      if (!state.currentPrep || !state.currentPrep.khazwal_material_prep) {
        return false
      }
      const prep = state.currentPrep.khazwal_material_prep
      return prep.plat_match === true && prep.plat_retrieved_at !== null
    }
  },

  actions: {
    /**
     * Get Material Prep Queue dengan filter dan pagination
     * @param {Object} filters - Filter options (search, priority, page, per_page)
     * @returns {Promise<Object>} Queue response dengan pagination
     */
    async getMaterialPrepQueue(filters = {}) {
      this.queueLoading = true
      this.queueError = null

      try {
        const response = await apiClient.get('/khazwal/material-prep/queue', {
          params: {
            search: filters.search || '',
            priority: filters.priority || '',
            page: filters.page || 1,
            per_page: filters.per_page || 20,
            sort_by: filters.sort_by || 'priority_score',
            sort_dir: filters.sort_dir || 'DESC'
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
          throw new Error(response.data.message || 'Gagal mengambil queue')
        }
      } catch (error) {
        this.queueError = error.response?.data?.message || error.message || 'Gagal mengambil queue'
        throw error
      } finally {
        this.queueLoading = false
      }
    },

    /**
     * Get Material Prep Detail dengan full OBC Master data
     * @param {Number} poId - Production Order ID
     * @returns {Promise<Object>} PO detail dengan OBC Master relationship
     */
    async getMaterialPrepDetail(poId) {
      this.currentPrepLoading = true
      this.currentPrepError = null

      try {
        const response = await apiClient.get(`/khazwal/material-prep/${poId}`)

        if (response.data.success) {
          this.currentPrep = response.data.data
          return response.data.data
        } else {
          throw new Error(response.data.message || 'Gagal mengambil detail')
        }
      } catch (error) {
        this.currentPrepError = error.response?.data?.message || error.message || 'Gagal mengambil detail'
        throw error
      } finally {
        this.currentPrepLoading = false
      }
    },

    /**
     * Start Material Prep untuk PO tertentu
     * @param {Number} poId - Production Order ID
     * @returns {Promise<Object>} Updated PO dengan material prep initialized
     */
    async startMaterialPrep(poId) {
      this.currentPrepLoading = true
      this.currentPrepError = null

      try {
        const response = await apiClient.post(`/khazwal/material-prep/${poId}/start`)

        if (response.data.success) {
          this.currentPrep = response.data.data
          
          // Remove dari queue jika ada
          this.queue = this.queue.filter(item => item.id !== poId)
          
          return response.data.data
        } else {
          throw new Error(response.data.message || 'Gagal memulai persiapan')
        }
      } catch (error) {
        this.currentPrepError = error.response?.data?.message || error.message || 'Gagal memulai persiapan'
        throw error
      } finally {
        this.currentPrepLoading = false
      }
    },

    /**
     * Confirm Plat dengan scanned barcode
     * @param {Number} prepId - Material Prep ID
     * @param {String} scannedCode - Barcode yang di-scan
     * @returns {Promise<Object>} Response dari server
     */
    async confirmPlat(prepId, scannedCode) {
      try {
        const response = await apiClient.post(`/khazwal/material-prep/${prepId}/confirm-plat`, {
          plat_code: scannedCode
        })

        if (response.data.success) {
          // Update current prep jika ada
          if (this.currentPrep && this.currentPrep.khazwal_material_prep?.id === prepId) {
            this.currentPrep.khazwal_material_prep.plat_scanned_code = scannedCode
            this.currentPrep.khazwal_material_prep.plat_match = true
            this.currentPrep.khazwal_material_prep.plat_retrieved_at = new Date().toISOString()
          }
          
          return response.data
        } else {
          throw new Error(response.data.message || 'Konfirmasi plat gagal')
        }
      } catch (error) {
        throw error
      }
    },

    /**
     * Update Kertas Blanko actual quantity
     * @param {Number} prepId - Material Prep ID
     * @param {Number} actualQuantity - Jumlah kertas actual
     * @param {String} varianceReason - Alasan variance (jika > 5%)
     * @returns {Promise<Object>} Response dari server
     */
    async updateKertas(prepId, actualQuantity, varianceReason = '') {
      try {
        const response = await apiClient.patch(`/khazwal/material-prep/${prepId}/kertas`, {
          actual_quantity: actualQuantity,
          variance_reason: varianceReason
        })

        if (response.data.success) {
          // Update current prep jika ada
          if (this.currentPrep && this.currentPrep.khazwal_material_prep?.id === prepId) {
            const prep = this.currentPrep.khazwal_material_prep
            prep.kertas_blanko_actual = actualQuantity
            const variance = actualQuantity - prep.kertas_blanko_quantity
            prep.kertas_blanko_variance = variance
            prep.kertas_blanko_variance_percentage = 
              (variance / prep.kertas_blanko_quantity) * 100
            if (varianceReason) {
              prep.kertas_blanko_variance_reason = varianceReason
            }
          }
          
          return response.data
        } else {
          throw new Error(response.data.message || 'Update kertas gagal')
        }
      } catch (error) {
        throw error
      }
    },

    /**
     * Update Tinta actual per color
     * @param {Number} prepId - Material Prep ID
     * @param {Object} tintaActual - Tinta actual data
     * @returns {Promise<Object>} Response dari server
     */
    async updateTinta(prepId, tintaActual) {
      try {
        const response = await apiClient.patch(`/khazwal/material-prep/${prepId}/tinta`, tintaActual)

        if (response.data.success) {
          // Update current prep jika ada
          if (this.currentPrep && this.currentPrep.khazwal_material_prep?.id === prepId) {
            this.currentPrep.khazwal_material_prep.tinta_actual = tintaActual
            // Update low stock flags jika ada dari server
            if (response.data.data?.tinta_low_stock_flags) {
              this.currentPrep.khazwal_material_prep.tinta_low_stock_flags = 
                response.data.data.tinta_low_stock_flags
            }
          }
          
          return response.data
        } else {
          throw new Error(response.data.message || 'Update tinta gagal')
        }
      } catch (error) {
        throw error
      }
    },

    /**
     * Finalize Material Prep dan kirim ke Unit Cetak
     * @param {Number} prepId - Material Prep ID
     * @param {Array<String>} photos - Array of photo URLs (optional)
     * @param {String} notes - Additional notes (optional)
     * @returns {Promise<Object>} Finalize result dengan duration
     */
    async finalizeMaterialPrep(prepId, photos = [], notes = '') {
      try {
        const response = await apiClient.post(`/khazwal/material-prep/${prepId}/finalize`, {
          photos: photos,
          notes: notes
        })

        if (response.data.success) {
          // Update current prep status
          if (this.currentPrep && this.currentPrep.khazwal_material_prep?.id === prepId) {
            this.currentPrep.khazwal_material_prep.status = 'COMPLETED'
            this.currentPrep.khazwal_material_prep.completed_at = new Date().toISOString()
            this.currentPrep.current_stage = 'CETAK'
            this.currentPrep.current_status = 'READY_FOR_CETAK'
          }
          
          return response.data
        } else {
          throw new Error(response.data.message || 'Finalisasi gagal')
        }
      } catch (error) {
        throw error
      }
    },

    /**
     * Get Material Prep History dengan filter
     * @param {Object} filters - Filter options (search, staff_id, date_from, date_to)
     * @returns {Promise<Object>} History response dengan pagination
     */
    async getMaterialPrepHistory(filters = {}) {
      this.historyLoading = true
      this.historyError = null

      try {
        const response = await apiClient.get('/khazwal/material-prep/history', {
          params: {
            search: filters.search || '',
            staff_id: filters.staff_id || '',
            date_from: filters.date_from || '',
            date_to: filters.date_to || '',
            page: filters.page || 1,
            per_page: filters.per_page || 20
          }
        })

        if (response.data.success) {
          this.history = response.data.data.items || []
          this.historyPagination = {
            total: response.data.data.total || 0,
            page: response.data.data.page || 1,
            per_page: response.data.data.per_page || 20,
            total_pages: response.data.data.total_pages || 0
          }
          return response.data.data
        } else {
          throw new Error(response.data.message || 'Gagal mengambil history')
        }
      } catch (error) {
        this.historyError = error.response?.data?.message || error.message || 'Gagal mengambil history'
        throw error
      } finally {
        this.historyLoading = false
      }
    },

    /**
     * Get Monitoring Stats untuk Supervisor Dashboard
     * @returns {Promise<Object>} Monitoring statistics
     */
    async getMonitoringStats() {
      this.monitoringLoading = true
      this.monitoringError = null

      try {
        const response = await apiClient.get('/khazwal/monitoring')

        if (response.data.success) {
          this.monitoringStats = response.data.data
          return response.data.data
        } else {
          throw new Error(response.data.message || 'Gagal mengambil monitoring stats')
        }
      } catch (error) {
        this.monitoringError = error.response?.data?.message || error.message || 'Gagal mengambil monitoring stats'
        throw error
      } finally {
        this.monitoringLoading = false
      }
    },

    /**
     * Clear current prep state
     */
    clearCurrentPrep() {
      this.currentPrep = null
      this.currentPrepError = null
    },

    /**
     * Clear all state (untuk logout atau reset)
     */
    clearAll() {
      this.queue = []
      this.currentPrep = null
      this.history = []
      this.monitoringStats = null
      this.queueError = null
      this.currentPrepError = null
      this.historyError = null
      this.monitoringError = null
    }
  }
})
