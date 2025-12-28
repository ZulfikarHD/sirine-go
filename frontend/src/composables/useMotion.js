/**
 * useMotion - Composable untuk standardized Motion-V animation presets
 * yang mengimplementasikan iOS-inspired spring physics dengan optimal
 * performance settings untuk consistent UX across the application
 */

/**
 * iOS-like Spring Physics Presets
 * Optimized untuk 60fps smooth animations
 */
export const springPresets = {
  // Default spring - natural, balanced
  default: {
    type: 'spring',
    stiffness: 400,
    damping: 30,
    mass: 0.8
  },
  // Snappy - quick responsive actions
  snappy: {
    type: 'spring',
    stiffness: 500,
    damping: 35,
    mass: 0.6
  },
  // Gentle - subtle, soft movements
  gentle: {
    type: 'spring',
    stiffness: 300,
    damping: 25,
    mass: 1
  },
  // Bouncy - playful, more overshoot
  bouncy: {
    type: 'spring',
    stiffness: 350,
    damping: 20,
    mass: 0.8
  }
}

/**
 * Easing Presets untuk non-spring animations
 */
export const easingPresets = {
  easeOut: 'easeOut',
  easeIn: 'easeIn',
  easeInOut: 'easeInOut',
  linear: 'linear'
}

/**
 * Page/Section Entrance Animations
 * untuk konsisten page load experience
 */
export const entranceAnimations = {
  // Fade in from bottom - default page entrance
  fadeUp: {
    initial: { opacity: 0, y: 15 },
    animate: { opacity: 1, y: 0 },
    transition: { duration: 0.25, ease: 'easeOut' }
  },
  // Fade in with scale - modal/dialog entrance
  fadeScale: {
    initial: { opacity: 0, scale: 0.95 },
    animate: { opacity: 1, scale: 1 },
    transition: { ...springPresets.default }
  },
  // Simple fade - subtle elements
  fade: {
    initial: { opacity: 0 },
    animate: { opacity: 1 },
    transition: { duration: 0.2, ease: 'easeOut' }
  },
  // Slide from left
  slideLeft: {
    initial: { opacity: 0, x: -20 },
    animate: { opacity: 1, x: 0 },
    transition: { duration: 0.25, ease: 'easeOut' }
  },
  // Slide from right
  slideRight: {
    initial: { opacity: 0, x: 20 },
    animate: { opacity: 1, x: 0 },
    transition: { duration: 0.25, ease: 'easeOut' }
  }
}

/**
 * Stagger Animation Helper
 * untuk list items dengan sequential entrance
 * @param {number} index - Item index dalam list
 * @param {number} baseDelay - Base delay sebelum animasi dimulai
 * @param {number} staggerDelay - Delay antar item (default 0.05s)
 */
export const getStaggerDelay = (index, baseDelay = 0, staggerDelay = 0.05) => {
  return baseDelay + (index * staggerDelay)
}

/**
 * List Item Animation dengan Stagger
 * @param {number} index - Item index
 * @param {number} baseDelay - Base delay
 */
export const listItemAnimation = (index, baseDelay = 0.1) => ({
  initial: { opacity: 0, y: 15 },
  animate: { opacity: 1, y: 0 },
  transition: {
    duration: 0.25,
    delay: getStaggerDelay(index, baseDelay),
    ease: 'easeOut'
  }
})

/**
 * Modal/Dialog Animations
 */
export const modalAnimations = {
  // Desktop modal
  desktop: {
    initial: { opacity: 0, scale: 0.95, y: 20 },
    animate: { opacity: 1, scale: 1, y: 0 },
    exit: { opacity: 0, scale: 0.95, y: 20 },
    transition: { ...springPresets.snappy }
  },
  // Mobile bottom sheet
  mobile: {
    initial: { opacity: 0, y: '100%' },
    animate: { opacity: 1, y: 0 },
    exit: { opacity: 0, y: '100%' },
    transition: { ...springPresets.snappy }
  },
  // Backdrop fade
  backdrop: {
    initial: { opacity: 0 },
    animate: { opacity: 1 },
    exit: { opacity: 0 },
    transition: { duration: 0.2, ease: 'easeOut' }
  }
}

/**
 * Icon Animations - untuk dialog icons, notifications
 */
export const iconAnimations = {
  // Pop in with spring
  popIn: {
    initial: { opacity: 0, scale: 0.8 },
    animate: { opacity: 1, scale: 1 },
    transition: { ...springPresets.bouncy }
  },
  // Rotate in - untuk refresh, loading states
  rotateIn: {
    initial: { opacity: 0, rotate: -90, scale: 0.8 },
    animate: { opacity: 1, rotate: 0, scale: 1 },
    transition: { ...springPresets.default }
  }
}

/**
 * Button/Interactive Element Animations
 */
export const interactiveAnimations = {
  // Press feedback scale
  pressScale: 0.97,
  // Hover lift
  hoverLift: { y: -2 },
  // Tap duration
  tapDuration: 0.12
}

/**
 * Exit Animations
 */
export const exitAnimations = {
  fadeOut: {
    exit: { opacity: 0 },
    transition: { duration: 0.15, ease: 'easeIn' }
  },
  fadeDown: {
    exit: { opacity: 0, y: 10 },
    transition: { duration: 0.15, ease: 'easeIn' }
  },
  scaleOut: {
    exit: { opacity: 0, scale: 0.95 },
    transition: { duration: 0.15, ease: 'easeIn' }
  }
}

/**
 * Shake Animation untuk Error States
 */
export const shakeAnimation = {
  animate: {
    x: [0, -8, 8, -8, 0]
  },
  transition: {
    duration: 0.4,
    ease: 'easeInOut'
  }
}

/**
 * Composable Hook untuk Vue components
 */
export function useMotion() {
  return {
    // Presets
    springPresets,
    easingPresets,
    
    // Entrance animations
    entranceAnimations,
    listItemAnimation,
    getStaggerDelay,
    
    // Modal animations
    modalAnimations,
    
    // Icon animations
    iconAnimations,
    
    // Interactive
    interactiveAnimations,
    
    // Exit
    exitAnimations,
    
    // Error
    shakeAnimation
  }
}

export default useMotion
