import { useApi } from './useApi'

/**
 * Composable untuk Cutting API operations
 * Menyediakan methods untuk interaksi dengan backend cutting endpoints
 */
export function useCuttingApi() {
  const api = useApi()

  /**
   * Mengambil cutting queue dengan optional filters
   * @param {Object} filters - Filter parameters (priority, date_from, date_to, sort_by, sort_order)
   * @returns {Promise<Object>} Response dengan data dan meta
   */
  const getCuttingQueue = async (filters = {}) => {
    const params = new URLSearchParams()
    
    if (filters.priority) params.append('priority', filters.priority)
    if (filters.date_from) params.append('date_from', filters.date_from)
    if (filters.date_to) params.append('date_to', filters.date_to)
    if (filters.sort_by) params.append('sort_by', filters.sort_by)
    if (filters.sort_order) params.append('sort_order', filters.sort_order)
    
    const queryString = params.toString()
    const url = queryString ? `/khazwal/cutting/queue?${queryString}` : '/khazwal/cutting/queue'
    
    return await api.get(url)
  }

  /**
   * Mengambil detail cutting record berdasarkan ID
   * @param {Number} id - Cutting record ID
   * @returns {Promise<Object>} Cutting detail data
   */
  const getCuttingDetail = async (id) => {
    return await api.get(`/khazwal/cutting/${id}`)
  }

  /**
   * Memulai proses cutting untuk PO tertentu
   * @param {Number} poId - Production Order ID
   * @param {Object} payload - Request payload dengan cutting_machine
   * @returns {Promise<Object>} Start cutting response
   */
  const startCutting = async (poId, payload) => {
    return await api.post(`/khazwal/cutting/po/${poId}/start`, payload)
  }

  /**
   * Mengupdate hasil cutting (sisiran kiri & kanan)
   * @param {Number} id - Cutting record ID
   * @param {Object} payload - Request payload dengan output_sisiran_kiri, output_sisiran_kanan, waste_reason, waste_photo_url
   * @returns {Promise<Object>} Update result response
   */
  const updateCuttingResult = async (id, payload) => {
    return await api.patch(`/khazwal/cutting/${id}/result`, payload)
  }

  /**
   * Finalisasi proses cutting dan generate verification labels
   * @param {Number} id - Cutting record ID
   * @returns {Promise<Object>} Finalize response
   */
  const finalizeCutting = async (id) => {
    return await api.post(`/khazwal/cutting/${id}/finalize`)
  }

  return {
    getCuttingQueue,
    getCuttingDetail,
    startCutting,
    updateCuttingResult,
    finalizeCutting,
  }
}
