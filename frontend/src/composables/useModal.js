import { ref, markRaw } from 'vue'

/**
 * useModal - Composable untuk managing modal state dan actions
 * dengan support untuk multiple modal instances dan programmatic control
 * yang memudahkan handling modal lifecycle dan user interactions
 */

export function useModal() {
  const isOpen = ref(false)
  const loading = ref(false)

  /**
   * Open modal dengan optional callback
   */
  const open = (callback) => {
    isOpen.value = true
    if (callback && typeof callback === 'function') {
      callback()
    }
  }

  /**
   * Close modal dengan optional callback
   */
  const close = (callback) => {
    isOpen.value = false
    if (callback && typeof callback === 'function') {
      // Delay callback untuk animation
      setTimeout(callback, 300)
    }
  }

  /**
   * Toggle modal state
   */
  const toggle = () => {
    isOpen.value = !isOpen.value
  }

  /**
   * Set loading state untuk async operations
   */
  const setLoading = (value) => {
    loading.value = value
  }

  return {
    isOpen,
    loading,
    open,
    close,
    toggle,
    setLoading
  }
}

/**
 * useConfirmDialog - Composable khusus untuk confirmation dialogs
 * dengan promise-based API untuk better async handling
 */
export function useConfirmDialog() {
  const isOpen = ref(false)
  const loading = ref(false)
  const config = ref({
    title: 'Konfirmasi',
    message: '',
    detail: '',
    variant: 'default',
    confirmText: 'Ya, Lanjutkan',
    cancelText: 'Batal',
    showWarning: true,
    warningMessage: ''
  })

  let resolvePromise = null
  let rejectPromise = null

  /**
   * Show confirmation dialog dan return promise
   * @param {Object} options - Configuration options
   * @returns {Promise<boolean>} - Resolves dengan true jika confirmed, false jika cancelled
   */
  const confirm = (options = {}) => {
    config.value = { ...config.value, ...options }
    isOpen.value = true

    return new Promise((resolve, reject) => {
      resolvePromise = resolve
      rejectPromise = reject
    })
  }

  /**
   * Handle confirm action
   */
  const handleConfirm = () => {
    if (resolvePromise) {
      resolvePromise(true)
      resolvePromise = null
    }
    close()
  }

  /**
   * Handle cancel action
   */
  const handleCancel = () => {
    if (resolvePromise) {
      resolvePromise(false)
      resolvePromise = null
    }
    close()
  }

  /**
   * Close dialog
   */
  const close = () => {
    isOpen.value = false
    loading.value = false
  }

  return {
    isOpen,
    loading,
    config,
    confirm,
    handleConfirm,
    handleCancel,
    close
  }
}

/**
 * useAlertDialog - Composable untuk alert/notification dialogs
 * dengan support untuk auto-dismiss dan different alert types
 */
export function useAlertDialog() {
  const isOpen = ref(false)
  const config = ref({
    title: '',
    message: '',
    detail: '',
    variant: 'success',
    confirmText: 'OK',
    autoDismiss: false,
    autoDismissDelay: 3000,
    showTimer: true
  })

  let resolvePromise = null

  /**
   * Show alert dialog
   * @param {Object} options - Configuration options
   * @returns {Promise<void>}
   */
  const alert = (options = {}) => {
    config.value = { ...config.value, ...options }
    isOpen.value = true

    return new Promise((resolve) => {
      resolvePromise = resolve
    })
  }

  /**
   * Shorthand methods untuk different alert types
   */
  const success = (message, options = {}) => {
    return alert({
      message,
      variant: 'success',
      ...options
    })
  }

  const error = (message, options = {}) => {
    return alert({
      message,
      variant: 'error',
      ...options
    })
  }

  const warning = (message, options = {}) => {
    return alert({
      message,
      variant: 'warning',
      ...options
    })
  }

  const info = (message, options = {}) => {
    return alert({
      message,
      variant: 'info',
      ...options
    })
  }

  /**
   * Handle close action
   */
  const handleClose = () => {
    if (resolvePromise) {
      resolvePromise()
      resolvePromise = null
    }
    close()
  }

  /**
   * Close dialog
   */
  const close = () => {
    isOpen.value = false
  }

  return {
    isOpen,
    config,
    alert,
    success,
    error,
    warning,
    info,
    handleClose,
    close
  }
}

/**
 * usePasswordConfirmDialog - Composable untuk password confirmation dialogs
 * dengan promise-based API dan error handling untuk secure verification
 */
export function usePasswordConfirmDialog() {
  const isOpen = ref(false)
  const loading = ref(false)
  const config = ref({
    title: 'Konfirmasi Password',
    subtitle: '',
    message: 'Untuk keamanan, masukkan password Anda untuk melanjutkan.',
    detail: '',
    confirmText: 'Verifikasi',
    showSecurityNotice: true,
    securityNotice: '',
    validatePassword: null
  })

  let resolvePromise = null
  let rejectPromise = null
  let dialogRef = null

  /**
   * Show password confirmation dialog dan return promise dengan password
   * @param {Object} options - Configuration options
   * @returns {Promise<string|null>} - Resolves dengan password jika confirmed, null jika cancelled
   */
  const confirm = (options = {}) => {
    config.value = { ...config.value, ...options }
    isOpen.value = true

    return new Promise((resolve, reject) => {
      resolvePromise = resolve
      rejectPromise = reject
    })
  }

  /**
   * Handle confirm action dengan password
   * @param {string} password - User's password
   */
  const handleConfirm = (password) => {
    if (resolvePromise) {
      resolvePromise(password)
      // Don't clear promise yet - let parent verify and close
    }
  }

  /**
   * Handle cancel action
   */
  const handleCancel = () => {
    if (resolvePromise) {
      resolvePromise(null)
      resolvePromise = null
    }
    close()
  }

  /**
   * Set error message
   * @param {string} message - Error message to display
   */
  const setError = (message) => {
    if (dialogRef) {
      dialogRef.setError(message)
    }
  }

  /**
   * Set loading state
   * @param {boolean} value - Loading state
   */
  const setLoading = (value) => {
    loading.value = value
    if (dialogRef) {
      dialogRef.setLoading(value)
    }
  }

  /**
   * Close dialog dan reset state
   */
  const close = () => {
    isOpen.value = false
    loading.value = false
    if (dialogRef) {
      dialogRef.close()
    }
    if (resolvePromise) {
      resolvePromise = null
      rejectPromise = null
    }
  }

  /**
   * Set dialog reference untuk direct control
   */
  const setDialogRef = (ref) => {
    dialogRef = ref
  }

  return {
    isOpen,
    loading,
    config,
    confirm,
    handleConfirm,
    handleCancel,
    setError,
    setLoading,
    close,
    setDialogRef
  }
}

/**
 * Global modal manager untuk managing multiple modals
 * dengan z-index stacking dan keyboard navigation
 */
class ModalManager {
  constructor() {
    this.modals = []
    this.zIndexBase = 1000
  }

  /**
   * Register modal instance
   */
  register(modal) {
    this.modals.push(modal)
    return this.modals.length - 1
  }

  /**
   * Unregister modal instance
   */
  unregister(modalId) {
    this.modals = this.modals.filter((_, index) => index !== modalId)
  }

  /**
   * Get z-index untuk modal
   */
  getZIndex(modalId) {
    return this.zIndexBase + modalId * 10
  }

  /**
   * Get active modal (top-most)
   */
  getActiveModal() {
    return this.modals[this.modals.length - 1]
  }

  /**
   * Check if any modal is open
   */
  hasOpenModal() {
    return this.modals.length > 0
  }
}

export const modalManager = new ModalManager()

export default useModal
