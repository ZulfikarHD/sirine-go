import { useApi } from './useApi'

/**
 * Composable untuk Verifikasi API operations
 * Menyediakan methods untuk interaksi dengan backend verification endpoints
 */
export function useVerifikasiApi() {
  const api = useApi()

  /**
   * Mengambil verification queue (PO yang siap untuk diverifikasi)
   * @param {Object} filters - Filter parameters (priority, date_from, date_to)
   * @returns {Promise<Object>} Response dengan data dan meta
   */
  const getVerificationQueue = async (filters = {}) => {
    const params = new URLSearchParams()
    
    if (filters.priority) params.append('priority', filters.priority)
    if (filters.date_from) params.append('date_from', filters.date_from)
    if (filters.date_to) params.append('date_to', filters.date_to)
    
    const queryString = params.toString()
    const url = queryString ? `/verifikasi/queue?${queryString}` : '/verifikasi/queue'
    
    return await api.get(url)
  }

  /**
   * Mengambil detail verification labels untuk PO tertentu
   * @param {Number} poId - Production Order ID
   * @returns {Promise<Object>} Verification detail dengan list labels
   */
  const getVerificationDetail = async (poId) => {
    return await api.get(`/verifikasi/${poId}`)
  }

  /**
   * Mengambil detail label spesifik
   * @param {Number} labelId - Verification label ID
   * @returns {Promise<Object>} Label detail data
   */
  const getLabelDetail = async (labelId) => {
    return await api.get(`/verifikasi/label/${labelId}`)
  }

  /**
   * Memulai proses verifikasi untuk label tertentu
   * @param {Number} labelId - Verification label ID
   * @param {Object} payload - Request payload dengan inspector info
   * @returns {Promise<Object>} Start verification response
   */
  const startVerification = async (labelId, payload) => {
    return await api.post(`/verifikasi/label/${labelId}/start`, payload)
  }

  /**
   * Mengupdate hasil verifikasi (HCS/HCTS)
   * @param {Number} labelId - Verification label ID
   * @param {Object} payload - Request payload dengan verification result
   * @returns {Promise<Object>} Update verification response
   */
  const updateVerificationResult = async (labelId, payload) => {
    return await api.patch(`/verifikasi/label/${labelId}/result`, payload)
  }

  /**
   * Finalisasi verifikasi label
   * @param {Number} labelId - Verification label ID
   * @returns {Promise<Object>} Finalize response
   */
  const finalizeVerification = async (labelId) => {
    return await api.post(`/verifikasi/label/${labelId}/finalize`)
  }

  return {
    getVerificationQueue,
    getVerificationDetail,
    getLabelDetail,
    startVerification,
    updateVerificationResult,
    finalizeVerification,
  }
}
