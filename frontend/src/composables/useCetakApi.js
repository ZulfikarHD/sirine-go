import { useApi } from './useApi'

/**
 * useCetakApi composable untuk Unit Cetak API calls
 * yang mencakup queue retrieval dan detail view untuk PO siap cetak
 */
export const useCetakApi = () => {
  const { get } = useApi()

  /**
   * Mengambil queue list PO yang siap untuk dicetak
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

    const queryString = params.toString()
    const url = `/cetak/queue${queryString ? '?' + queryString : ''}`
    
    return await get(url)
  }

  /**
   * Mengambil detail PO untuk cetak termasuk material photos
   * @param {number} id - Production Order ID
   * @returns {Promise<Object>} PO detail dengan material prep info dan photos
   */
  const getDetail = async (id) => {
    return await get(`/cetak/queue/${id}`)
  }

  return {
    getQueue,
    getDetail
  }
}

export default useCetakApi
