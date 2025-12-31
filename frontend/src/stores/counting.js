import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { useCountingApi } from '@/composables/useCountingApi'

/**
 * Counting Store untuk state management penghitungan Epic 2
 * dengan reactive queue, current counting session, dan auto-calculation
 */
export const useCountingStore = defineStore('counting', () => {
  const countingApi = useCountingApi()

  // State
  const queue = ref([])
  const queueMeta = ref({ total: 0, overdue_count: 0 })
  const currentCounting = ref(null)
  const currentPO = ref(null)
  const countingResult = ref({
    quantity_good: 0,
    quantity_defect: 0,
    defect_breakdown: [],
    variance_reason: ''
  })
  
  const isLoadingQueue = ref(false)
  const isLoadingDetail = ref(false)
  const isStarting = ref(false)
  const isSaving = ref(false)
  const isFinalizing = ref(false)
  const error = ref(null)

  // Filters untuk queue
  const queueFilters = ref({
    machine_id: null,
    date_from: null,
    date_to: null
  })

  // Computed - Real-time calculations
  const calculations = computed(() => {
    if (!currentPO.value) {
      return null
    }

    return countingApi.calculateCountingStats(
      countingResult.value.quantity_good,
      countingResult.value.quantity_defect,
      currentPO.value.target_quantity
    )
  })

  // Computed - Validation states
  const isDefectBreakdownRequired = computed(() => {
    return calculations.value?.isDefectAboveThreshold || false
  })

  const isVarianceReasonRequired = computed(() => {
    return calculations.value?.hasVariance || false
  })

  const isDefectBreakdownValid = computed(() => {
    if (!isDefectBreakdownRequired.value) {
      return true
    }

    const validation = countingApi.validateDefectBreakdown(
      countingResult.value.defect_breakdown,
      countingResult.value.quantity_defect
    )
    return validation.isValid
  })

  const canFinalize = computed(() => {
    if (!currentCounting.value || currentCounting.value.status !== 'IN_PROGRESS') {
      return false
    }

    // Check quantities filled
    const hasQuantities = countingResult.value.quantity_good > 0 || 
                          countingResult.value.quantity_defect > 0

    if (!hasQuantities) {
      return false
    }

    // Check defect breakdown requirement
    if (isDefectBreakdownRequired.value && !isDefectBreakdownValid.value) {
      return false
    }

    // Check variance reason requirement
    if (isVarianceReasonRequired.value && !countingResult.value.variance_reason.trim()) {
      return false
    }

    return true
  })

  // Computed - Overdue items in queue
  const overdueItems = computed(() => {
    return queue.value.filter(item => item.is_overdue)
  })

  // Actions

  /**
   * Fetch counting queue dengan optional filters
   */
  const fetchQueue = async (filters = {}) => {
    isLoadingQueue.value = true
    error.value = null

    try {
      const response = await countingApi.getCountingQueue(filters)
      
      if (response.success) {
        queue.value = response.data || []
        queueMeta.value = response.meta || { total: 0, overdue_count: 0 }
      }
    } catch (err) {
      error.value = err.response?.data?.message || 'Gagal mengambil counting queue'
      console.error('Fetch queue error:', err)
    } finally {
      isLoadingQueue.value = false
    }
  }

  /**
   * Fetch counting detail by ID
   */
  const fetchCountingDetail = async (id) => {
    isLoadingDetail.value = true
    error.value = null

    try {
      const response = await countingApi.getCountingDetail(id)
      
      if (response.success) {
        currentCounting.value = response.data
        currentPO.value = response.data.po
        
        // Populate form dengan existing data
        countingResult.value = {
          quantity_good: response.data.quantity_good || 0,
          quantity_defect: response.data.quantity_defect || 0,
          defect_breakdown: response.data.defect_breakdown || [],
          variance_reason: response.data.variance_reason || ''
        }
      }
    } catch (err) {
      error.value = err.response?.data?.message || 'Gagal mengambil counting detail'
      console.error('Fetch detail error:', err)
    } finally {
      isLoadingDetail.value = false
    }
  }

  /**
   * Start counting untuk PO tertentu
   */
  const startCounting = async (poId, poData) => {
    isStarting.value = true
    error.value = null

    try {
      const response = await countingApi.startCounting(poId)
      
      if (response.success) {
        currentCounting.value = response.data
        currentPO.value = poData
        
        // Reset form
        countingResult.value = {
          quantity_good: 0,
          quantity_defect: 0,
          defect_breakdown: [],
          variance_reason: ''
        }

        return response.data
      }
    } catch (err) {
      error.value = err.response?.data?.message || 'Gagal memulai penghitungan'
      console.error('Start counting error:', err)
      throw err
    } finally {
      isStarting.value = false
    }
  }

  /**
   * Update counting result dengan auto-save
   */
  const updateResult = async () => {
    if (!currentCounting.value) {
      return
    }

    isSaving.value = true
    error.value = null

    try {
      const response = await countingApi.updateCountingResult(
        currentCounting.value.id,
        countingResult.value
      )
      
      if (response.success) {
        // Update current counting dengan response data (termasuk calculated fields)
        Object.assign(currentCounting.value, response.data)
        return response.data
      }
    } catch (err) {
      error.value = err.response?.data?.message || 'Gagal menyimpan hasil penghitungan'
      console.error('Update result error:', err)
      throw err
    } finally {
      isSaving.value = false
    }
  }

  /**
   * Finalize counting dengan validation
   */
  const finalize = async () => {
    if (!canFinalize.value) {
      throw new Error('Data belum lengkap untuk finalisasi')
    }

    isFinalizing.value = true
    error.value = null

    try {
      const response = await countingApi.finalizeCounting(currentCounting.value.id)
      
      if (response.success) {
        // Clear current counting session
        currentCounting.value = null
        currentPO.value = null
        countingResult.value = {
          quantity_good: 0,
          quantity_defect: 0,
          defect_breakdown: [],
          variance_reason: ''
        }

        return response.data
      }
    } catch (err) {
      error.value = err.response?.data?.message || 'Gagal menyelesaikan penghitungan'
      console.error('Finalize error:', err)
      throw err
    } finally {
      isFinalizing.value = false
    }
  }

  /**
   * Update quantity good dengan validation
   */
  const updateQuantityGood = (value) => {
    const numValue = parseInt(value) || 0
    countingResult.value.quantity_good = Math.max(0, numValue)
  }

  /**
   * Update quantity defect dengan validation
   */
  const updateQuantityDefect = (value) => {
    const numValue = parseInt(value) || 0
    countingResult.value.quantity_defect = Math.max(0, numValue)
  }

  /**
   * Update defect breakdown item
   */
  const updateDefectBreakdownItem = (type, quantity) => {
    const existingIndex = countingResult.value.defect_breakdown.findIndex(
      item => item.type === type
    )

    const numQuantity = parseInt(quantity) || 0

    if (numQuantity <= 0) {
      // Remove item jika quantity 0
      if (existingIndex !== -1) {
        countingResult.value.defect_breakdown.splice(existingIndex, 1)
      }
    } else {
      // Update atau add item
      if (existingIndex !== -1) {
        countingResult.value.defect_breakdown[existingIndex].quantity = numQuantity
      } else {
        countingResult.value.defect_breakdown.push({ type, quantity: numQuantity })
      }
    }
  }

  /**
   * Clear current counting session
   */
  const clearCurrentCounting = () => {
    currentCounting.value = null
    currentPO.value = null
    countingResult.value = {
      quantity_good: 0,
      quantity_defect: 0,
      defect_breakdown: [],
      variance_reason: ''
    }
    error.value = null
  }

  /**
   * Apply filters to queue
   */
  const applyFilters = async (filters) => {
    queueFilters.value = { ...filters }
    await fetchQueue(filters)
  }

  /**
   * Clear filters
   */
  const clearFilters = async () => {
    queueFilters.value = {
      machine_id: null,
      date_from: null,
      date_to: null
    }
    await fetchQueue()
  }

  return {
    // State
    queue,
    queueMeta,
    currentCounting,
    currentPO,
    countingResult,
    isLoadingQueue,
    isLoadingDetail,
    isStarting,
    isSaving,
    isFinalizing,
    error,
    queueFilters,

    // Computed
    calculations,
    isDefectBreakdownRequired,
    isVarianceReasonRequired,
    isDefectBreakdownValid,
    canFinalize,
    overdueItems,

    // Actions
    fetchQueue,
    fetchCountingDetail,
    startCounting,
    updateResult,
    finalize,
    updateQuantityGood,
    updateQuantityDefect,
    updateDefectBreakdownItem,
    clearCurrentCounting,
    applyFilters,
    clearFilters
  }
})
