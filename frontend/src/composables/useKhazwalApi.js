import { useApi } from './useApi'

/**
 * useKhazwalApi composable untuk Khazwal Material Preparation API calls
 * yang mencakup queue management, detail retrieval, dan workflow actions
 */
export const useKhazwalApi = () => {
  const { get, post, put, apiClient } = useApi()

  /**
   * Mengambil queue list PO untuk material preparation
   * dengan support untuk filter, search, dan pagination
   * @param {Object} filters - Filter parameters (priority, search, page, per_page)
   * @returns {Promise<Object>} Response dengan items array dan pagination metadata
   */
  const getQueue = async (filters = {}) => {
    const params = new URLSearchParams()
    
    // Add filters ke query params jika ada
    if (filters.search) params.append('search', filters.search)
    if (filters.priority) params.append('priority', filters.priority)
    if (filters.page) params.append('page', filters.page)
    if (filters.per_page) params.append('per_page', filters.per_page)
    if (filters.sort_by) params.append('sort_by', filters.sort_by)
    if (filters.sort_dir) params.append('sort_dir', filters.sort_dir)

    const queryString = params.toString()
    const url = `/khazwal/material-prep/queue${queryString ? '?' + queryString : ''}`
    
    return await get(url)
  }

  /**
   * Mengambil detail PO dengan full relations dan computed fields
   * untuk display di detail page
   * @param {number} id - Production Order ID
   * @returns {Promise<Object>} PO detail dengan material prep info
   */
  const getDetail = async (id) => {
    return await get(`/khazwal/material-prep/${id}`)
  }

  /**
   * Memulai proses material preparation untuk PO
   * dengan transaction untuk update status dan tracking
   * @param {number} id - Production Order ID
   * @returns {Promise<Object>} Updated PO data
   */
  const startPrep = async (id) => {
    return await post(`/khazwal/material-prep/${id}/start`)
  }

  /**
   * Mengkonfirmasi pengambilan plat dengan barcode scan
   * @param {number} id - Material Prep ID (bukan PO ID)
   * @param {string} platCode - SAP Plat Code dari barcode scan
   * @returns {Promise<Object>} Updated material prep data
   */
  const confirmPlat = async (id, platCode) => {
    return await post(`/khazwal/material-prep/${id}/confirm-plat`, {
      plat_code: platCode
    })
  }

  /**
   * Update kertas blanko actual quantity dengan variance calculation
   * @param {number} id - Material Prep ID (bukan PO ID)
   * @param {number} actualQty - Actual quantity yang diambil
   * @param {string} varianceReason - Alasan jika ada variance
   * @returns {Promise<Object>} Updated material prep data
   */
  const updateKertas = async (id, actualQty, varianceReason = '') => {
    const response = await apiClient.patch(`/khazwal/material-prep/${id}/kertas`, {
      actual_qty: actualQty,
      variance_reason: varianceReason
    })
    return response.data
  }

  /**
   * Update informasi tinta yang digunakan dengan checklist
   * @param {number} id - Material Prep ID (bukan PO ID)
   * @param {Array} tintaActual - Array tinta yang sudah dicek
   * @returns {Promise<Object>} Updated material prep data
   */
  const updateTinta = async (id, tintaActual) => {
    const response = await apiClient.patch(`/khazwal/material-prep/${id}/tinta`, {
      tinta_actual: tintaActual
    })
    return response.data
  }

  /**
   * Finalize material preparation dengan upload photos dan notes
   * yang mengirim ke Unit Cetak untuk proses selanjutnya
   * @param {number} id - Material Prep ID (bukan PO ID)
   * @param {Object} data - Data finalize (photos, notes)
   * @param {Array} data.photos - Array foto material dalam base64 format
   * @param {string} data.notes - Catatan tambahan (opsional)
   * @returns {Promise<Object>} Completion summary dengan duration dan status
   */
  const finalize = async (id, { photos = [], notes = '' } = {}) => {
    return await post(`/khazwal/material-prep/${id}/finalize`, {
      photos,
      notes
    })
  }

  /**
   * Mengambil riwayat material preparation yang sudah selesai
   * dengan support untuk filter by date range dan staff
   * @param {Object} filters - Filter parameters (search, staff_id, date_from, date_to, page, per_page)
   * @returns {Promise<Object>} Response dengan items array dan pagination metadata
   */
  const getHistory = async (filters = {}) => {
    const params = new URLSearchParams()
    
    // Add filters ke query params jika ada
    if (filters.search) params.append('search', filters.search)
    if (filters.staff_id) params.append('staff_id', filters.staff_id)
    if (filters.date_from) params.append('date_from', filters.date_from)
    if (filters.date_to) params.append('date_to', filters.date_to)
    if (filters.page) params.append('page', filters.page)
    if (filters.per_page) params.append('per_page', filters.per_page)

    const queryString = params.toString()
    const url = `/khazwal/material-prep/history${queryString ? '?' + queryString : ''}`
    
    return await get(url)
  }

  /**
   * Mengambil statistik monitoring untuk supervisor dashboard
   * @returns {Promise<Object>} Response dengan stats dan staff activity
   */
  const getMonitoring = async () => {
    return await get('/khazwal/monitoring')
  }

  return {
    // Queue & Detail
    getQueue,
    getDetail,
    
    // Workflow actions
    startPrep,
    confirmPlat,
    updateKertas,
    updateTinta,
    finalize,

    // History & Monitoring (Sprint 5)
    getHistory,
    getMonitoring
  }
}

export default useKhazwalApi
