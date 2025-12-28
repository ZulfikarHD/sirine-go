<template>
  <!-- Modal Backdrop dengan Glass Effect -->
  <Teleport to="body">
    <Motion
      v-if="isOpen"
      :initial="{ opacity: 0 }"
      :animate="{ opacity: 1 }"
      :exit="{ opacity: 0 }"
      :transition="{ duration: 0.3 }"
      class="fixed inset-0 bg-black/60 backdrop-blur-md z-50 flex items-end sm:items-center justify-center"
      @click.self="handleBackdropClick"
    >
      <!-- Modal Content dengan Spring Animation -->
      <Motion
        :initial="mobileView ? { opacity: 0, y: 100 } : { opacity: 0, scale: 0.95, y: 20 }"
        :animate="mobileView ? { opacity: 1, y: 0 } : { opacity: 1, scale: 1, y: 0 }"
        :exit="mobileView ? { opacity: 0, y: 100 } : { opacity: 0, scale: 0.95, y: 20 }"
        :transition="{ 
          duration: 0.4, 
          easing: 'spring',
          type: 'spring',
          stiffness: 300,
          damping: 30
        }"
        :class="[
          'glass-card shadow-2xl overflow-hidden',
          mobileView 
            ? 'w-full max-h-[90vh] rounded-t-3xl' 
            : 'w-full mx-4 rounded-2xl max-h-[90vh]',
          sizeClasses
        ]"
        @click.stop
      >
        <!-- Header dengan Gradient Accent -->
        <div 
          v-if="showHeader"
          :class="[
            'sticky top-0 z-10 px-6 py-4 border-b border-gray-200/50',
            mobileView ? 'glass-navbar' : 'bg-white/95 backdrop-blur-xl'
          ]"
        >
          <!-- Mobile Drag Handle -->
          <div v-if="mobileView && dismissible" class="flex justify-center mb-3">
            <div class="w-12 h-1.5 bg-gray-300 rounded-full active-scale cursor-grab" />
          </div>

          <div class="flex items-start justify-between">
            <div class="flex-1 pr-4">
              <!-- Title dengan Gradient -->
              <h2 
                v-if="title"
                :class="[
                  'font-bold leading-tight',
                  mobileView ? 'text-xl' : 'text-2xl',
                  titleGradient ? 'text-gradient-indigo-fuchsia' : 'text-gray-900'
                ]"
              >
                {{ title }}
              </h2>
              <slot name="title" />

              <!-- Subtitle -->
              <p 
                v-if="subtitle"
                class="text-sm text-gray-600 mt-1.5"
              >
                {{ subtitle }}
              </p>
              <slot name="subtitle" />
            </div>

            <!-- Close Button dengan Press Feedback -->
            <button
              v-if="dismissible"
              @click="handleClose"
              class="p-2.5 hover:bg-gray-100 rounded-xl transition-all duration-200 active:scale-95 shrink-0"
              :aria-label="$t?.('close') || 'Tutup'"
            >
              <X class="w-5 h-5 text-gray-500" />
            </button>
          </div>
        </div>

        <!-- Content dengan Custom Scrollbar -->
        <div 
          :class="[
            'custom-scrollbar',
            scrollable ? 'overflow-y-auto' : '',
            noPadding ? '' : 'p-6'
          ]"
          :style="{ maxHeight: contentMaxHeight }"
        >
          <slot />
        </div>

        <!-- Footer Sticky (untuk actions) -->
        <div 
          v-if="showFooter"
          :class="[
            'sticky bottom-0 px-6 py-4 border-t border-gray-200/50',
            mobileView ? 'glass-navbar' : 'bg-white/95 backdrop-blur-xl'
          ]"
        >
          <slot name="footer">
            <!-- Default Footer Actions -->
            <div class="flex items-center gap-3">
              <button
                v-if="showCancel"
                @click="handleCancel"
                :disabled="loading"
                class="flex-1 sm:flex-none btn-secondary disabled:opacity-50 disabled:cursor-not-allowed"
              >
                {{ cancelText }}
              </button>
              
              <button
                v-if="showConfirm"
                @click="handleConfirm"
                :disabled="loading || confirmDisabled"
                :class="[
                  'flex-1 btn-primary disabled:opacity-50 disabled:cursor-not-allowed',
                  confirmDanger ? 'bg-linear-to-r! from-red-600! to-pink-600! hover:from-red-700! hover:to-pink-700!' : ''
                ]"
              >
                <span v-if="loading" class="flex items-center justify-center gap-2">
                  <div class="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin" />
                  {{ loadingText }}
                </span>
                <span v-else>{{ confirmText }}</span>
              </button>
            </div>
          </slot>
        </div>
      </Motion>
    </Motion>
  </Teleport>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'
import { Motion } from 'motion-v'
import { X } from 'lucide-vue-next'
import { useWindowSize } from '@vueuse/core'

/**
 * BaseModal - Komponen modal reusable dengan iOS-inspired design
 * untuk CRUD operations, confirmations, dan alerts dengan
 * spring physics animations serta glass morphism effect
 */

const props = defineProps({
  // Visibility control
  modelValue: {
    type: Boolean,
    default: false
  },

  // Header props
  title: {
    type: String,
    default: ''
  },
  subtitle: {
    type: String,
    default: ''
  },
  titleGradient: {
    type: Boolean,
    default: false
  },
  showHeader: {
    type: Boolean,
    default: true
  },

  // Size variants untuk different use cases
  size: {
    type: String,
    default: 'md', // xs, sm, md, lg, xl, full
    validator: (value) => ['xs', 'sm', 'md', 'lg', 'xl', 'full'].includes(value)
  },

  // Content props
  scrollable: {
    type: Boolean,
    default: true
  },
  noPadding: {
    type: Boolean,
    default: false
  },
  contentMaxHeight: {
    type: String,
    default: '60vh'
  },

  // Footer props
  showFooter: {
    type: Boolean,
    default: false
  },
  showCancel: {
    type: Boolean,
    default: true
  },
  showConfirm: {
    type: Boolean,
    default: true
  },
  cancelText: {
    type: String,
    default: 'Batal'
  },
  confirmText: {
    type: String,
    default: 'Simpan'
  },
  loadingText: {
    type: String,
    default: 'Memproses...'
  },
  confirmDanger: {
    type: Boolean,
    default: false
  },
  confirmDisabled: {
    type: Boolean,
    default: false
  },

  // Behavior props
  dismissible: {
    type: Boolean,
    default: true
  },
  closeOnBackdrop: {
    type: Boolean,
    default: true
  },
  closeOnEscape: {
    type: Boolean,
    default: true
  },
  loading: {
    type: Boolean,
    default: false
  },

  // Haptic feedback (untuk mobile)
  enableHaptics: {
    type: Boolean,
    default: true
  }
})

const emit = defineEmits(['update:modelValue', 'close', 'cancel', 'confirm', 'opened', 'closed'])

// Reactive state
const isOpen = ref(props.modelValue)
const { width } = useWindowSize()

// Computed properties
const mobileView = computed(() => width.value < 640) // sm breakpoint

const sizeClasses = computed(() => {
  const sizes = {
    xs: 'sm:max-w-xs',
    sm: 'sm:max-w-sm',
    md: 'sm:max-w-md',
    lg: 'sm:max-w-lg',
    xl: 'sm:max-w-xl',
    full: 'sm:max-w-7xl'
  }
  return sizes[props.size] || sizes.md
})

// Watch modelValue untuk sync
watch(() => props.modelValue, (newValue) => {
  isOpen.value = newValue
  if (newValue) {
    handleOpen()
  }
})

// Methods dengan haptic feedback support
const triggerHaptic = (type = 'light') => {
  if (props.enableHaptics && 'vibrate' in navigator) {
    const patterns = {
      light: [10],
      medium: [20],
      heavy: [30],
      success: [10, 50, 10],
      warning: [20, 100, 20],
      error: [30, 100, 30, 100, 30]
    }
    navigator.vibrate(patterns[type] || patterns.light)
  }
}

const handleOpen = () => {
  // Prevent body scroll saat modal open
  document.body.style.overflow = 'hidden'
  triggerHaptic('light')
  emit('opened')
}

const handleClose = () => {
  if (!props.dismissible || props.loading) return
  
  isOpen.value = false
  triggerHaptic('light')
  
  // Restore body scroll
  setTimeout(() => {
    document.body.style.overflow = ''
    emit('update:modelValue', false)
    emit('close')
    emit('closed')
  }, 300)
}

const handleBackdropClick = () => {
  if (props.closeOnBackdrop && props.dismissible) {
    triggerHaptic('medium')
    handleClose()
  }
}

const handleCancel = () => {
  triggerHaptic('light')
  emit('cancel')
  handleClose()
}

const handleConfirm = () => {
  triggerHaptic('medium')
  emit('confirm')
}

const handleEscape = (e) => {
  if (e.key === 'Escape' && props.closeOnEscape && props.dismissible && !props.loading) {
    handleClose()
  }
}

// Lifecycle hooks
onMounted(() => {
  document.addEventListener('keydown', handleEscape)
  if (props.modelValue) {
    handleOpen()
  }
})

onUnmounted(() => {
  document.removeEventListener('keydown', handleEscape)
  document.body.style.overflow = ''
})

// Expose methods untuk parent component
defineExpose({
  close: handleClose
})
</script>
