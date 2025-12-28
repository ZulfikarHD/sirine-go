/**
 * useHaptic - Composable untuk haptic feedback dengan vibration API
 * yang mengimplementasikan tactile response patterns untuk enhance UX
 */

/**
 * Haptic Feedback Patterns
 * Optimized untuk berbagai interaction types
 */
export const hapticPatterns = {
  // Light tap - untuk button press, toggle
  light: 10,
  
  // Medium tap - untuk selection, confirmation
  medium: 20,
  
  // Heavy tap - untuk important actions
  heavy: 30,
  
  // Success - 2 quick pulses
  success: [10, 50, 10],
  
  // Error - single long pulse dengan pause
  error: [30, 100, 30],
  
  // Warning - 3 short pulses
  warning: [15, 50, 15, 50, 15],
  
  // Achievement unlock - celebratory pattern
  achievement: [10, 50, 10, 100, 20, 50, 10],
  
  // Notification - subtle double tap
  notification: [10, 100, 10],
}

/**
 * useHaptic composable untuk trigger vibration dengan pattern support
 */
export function useHaptic() {
  /**
   * Check apakah device support vibration API
   */
  const isSupported = () => {
    return 'vibrate' in navigator
  }

  /**
   * Trigger vibration dengan pattern atau duration
   * @param {number|array} pattern - Duration (ms) atau pattern array
   */
  const vibrate = (pattern) => {
    if (!isSupported()) {
      return false
    }

    try {
      navigator.vibrate(pattern)
      return true
    } catch (error) {
      console.warn('Haptic feedback failed:', error)
      return false
    }
  }

  /**
   * Stop ongoing vibration
   */
  const stop = () => {
    if (!isSupported()) {
      return false
    }

    try {
      navigator.vibrate(0)
      return true
    } catch (error) {
      console.warn('Stop haptic failed:', error)
      return false
    }
  }

  /**
   * Predefined haptic methods untuk common interactions
   */
  const light = () => vibrate(hapticPatterns.light)
  const medium = () => vibrate(hapticPatterns.medium)
  const heavy = () => vibrate(hapticPatterns.heavy)
  const success = () => vibrate(hapticPatterns.success)
  const error = () => vibrate(hapticPatterns.error)
  const warning = () => vibrate(hapticPatterns.warning)
  const achievement = () => vibrate(hapticPatterns.achievement)
  const notification = () => vibrate(hapticPatterns.notification)

  /**
   * Trigger haptic untuk button press dengan iOS-like feedback
   */
  const buttonPress = () => light()

  /**
   * Trigger haptic untuk login success
   */
  const loginSuccess = () => success()

  /**
   * Trigger haptic untuk form submission
   */
  const formSubmit = () => medium()

  /**
   * Trigger haptic untuk toggle switch
   */
  const toggle = () => light()

  /**
   * Trigger haptic untuk delete/destructive action
   */
  const destructive = () => heavy()

  return {
    // Core methods
    isSupported,
    vibrate,
    stop,
    
    // Pattern methods
    light,
    medium,
    heavy,
    success,
    error,
    warning,
    achievement,
    notification,
    
    // Context methods
    buttonPress,
    loginSuccess,
    formSubmit,
    toggle,
    destructive,
    
    // Patterns reference
    patterns: hapticPatterns,
  }
}

export default useHaptic
