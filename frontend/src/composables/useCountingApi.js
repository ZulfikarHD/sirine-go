import { useApi } from './useApi'

/**
 * useCountingApi composable untuk counting API operations
 * yang mencakup queue management, start counting, update results, dan finalize
 */
export const useCountingApi = () => {
  const { get, post, patch } = useApi()

  /**
   * Get counting queue - mengambil list PO yang menunggu penghitungan
   * dengan sorting FIFO berdasarkan print completion time
   * @param {Object} filters - Optional filters (machine_id, date_from, date_to)
   * @returns {Promise} Queue response dengan data dan metadata
   */
  const getCountingQueue = async (filters = {}) => {
    const params = new URLSearchParams()
    
    if (filters.machine_id) {
      params.append('machine_id', filters.machine_id)
    }
    if (filters.date_from) {
      params.append('date_from', filters.date_from)
    }
    if (filters.date_to) {
      params.append('date_to', filters.date_to)
    }

    const queryString = params.toString()
    const url = `/khazwal/counting/queue${queryString ? `?${queryString}` : ''}`
    
    return await get(url)
  }

  /**
   * Get counting detail - mengambil detail counting record dengan relasi
   * @param {number} id - Counting ID
   * @returns {Promise} Counting detail dengan PO info, print info, dll
   */
  const getCountingDetail = async (id) => {
    return await get(`/khazwal/counting/${id}`)
  }

  /**
   * Start counting - memulai proses penghitungan untuk PO tertentu
   * dengan create counting record dan update PO status ke SEDANG_DIHITUNG
   * @param {number} poId - Production Order ID
   * @returns {Promise} Start counting response dengan counting ID dan timestamps
   */
  const startCounting = async (poId) => {
    return await post(`/khazwal/counting/${poId}/start`, {})
  }

  /**
   * Update counting result - update hasil penghitungan
   * dapat dipanggil multiple times sebelum finalize untuk edit
   * @param {number} countingId - Counting ID
   * @param {Object} data - Result data
   * @param {number} data.quantity_good - Jumlah lembar besar yang baik
   * @param {number} data.quantity_defect - Jumlah lembar besar yang rusak
   * @param {Array} data.defect_breakdown - Breakdown jenis kerusakan (required jika defect > 5%)
   * @param {string} data.variance_reason - Alasan selisih (required jika variance != 0)
   * @returns {Promise} Updated result dengan calculations
   */
  const updateCountingResult = async (countingId, data) => {
    return await patch(`/khazwal/counting/${countingId}/result`, data)
  }

  /**
   * Finalize counting - menyelesaikan penghitungan dengan lock data
   * dan advance PO ke stage KHAZWAL_CUTTING dengan status SIAP_POTONG
   * @param {number} countingId - Counting ID
   * @returns {Promise} Finalize response dengan completion timestamp dan duration
   */
  const finalizeCounting = async (countingId) => {
    return await post(`/khazwal/counting/${countingId}/finalize`, {})
  }

  /**
   * Calculate real-time counting statistics untuk display
   * @param {number} quantityGood - Jumlah baik
   * @param {number} quantityDefect - Jumlah rusak
   * @param {number} targetQuantity - Target quantity dari PO
   * @returns {Object} Calculated statistics
   */
  const calculateCountingStats = (quantityGood, quantityDefect, targetQuantity) => {
    const totalCounted = quantityGood + quantityDefect
    const varianceFromTarget = totalCounted - targetQuantity
    const variancePercentage = targetQuantity > 0 
      ? ((varianceFromTarget / targetQuantity) * 100).toFixed(2)
      : 0

    const percentageGood = totalCounted > 0 
      ? ((quantityGood / totalCounted) * 100).toFixed(2)
      : 0
    const percentageDefect = totalCounted > 0 
      ? ((quantityDefect / totalCounted) * 100).toFixed(2)
      : 0

    return {
      totalCounted,
      varianceFromTarget,
      variancePercentage,
      percentageGood,
      percentageDefect,
      isDefectAboveThreshold: parseFloat(percentageDefect) > 5,
      hasVariance: varianceFromTarget !== 0,
      isWithinTolerance: Math.abs(parseFloat(variancePercentage)) <= 2
    }
  }

  /**
   * Validate defect breakdown sum untuk memastikan match dengan quantity_defect
   * @param {Array} breakdown - Defect breakdown array
   * @param {number} quantityDefect - Total quantity defect
   * @returns {Object} Validation result dengan isValid dan sum
   */
  const validateDefectBreakdown = (breakdown, quantityDefect) => {
    const sum = breakdown.reduce((total, item) => total + (item.quantity || 0), 0)
    return {
      isValid: sum === quantityDefect,
      sum,
      difference: sum - quantityDefect
    }
  }

  /**
   * Format waiting time untuk display user-friendly
   * @param {number} minutes - Waiting time dalam menit
   * @returns {string} Formatted waiting time
   */
  const formatWaitingTime = (minutes) => {
    if (minutes < 60) {
      return `${minutes} menit`
    }
    
    const hours = Math.floor(minutes / 60)
    const remainingMinutes = minutes % 60
    
    if (remainingMinutes === 0) {
      return `${hours} jam`
    }
    
    return `${hours} jam ${remainingMinutes} menit`
  }

  /**
   * Predefined defect types untuk dropdown selection
   */
  const defectTypes = [
    { type: 'Warna pudar', label: 'Warna Pudar' },
    { type: 'Tinta blobor', label: 'Tinta Blobor' },
    { type: 'Kertas sobek', label: 'Kertas Sobek' },
    { type: 'Register tidak pas', label: 'Register Tidak Pas' },
    { type: 'Lainnya', label: 'Lainnya' }
  ]

  return {
    // API calls
    getCountingQueue,
    getCountingDetail,
    startCounting,
    updateCountingResult,
    finalizeCounting,
    
    // Helper functions
    calculateCountingStats,
    validateDefectBreakdown,
    formatWaitingTime,
    
    // Constants
    defectTypes,
    
    // Thresholds
    DEFECT_BREAKDOWN_THRESHOLD: 5, // 5% - require breakdown
    TOLERANCE_THRESHOLD: 2, // 2% - warning threshold
    OVERDUE_THRESHOLD_MINUTES: 120 // 2 hours - overdue warning
  }
}
