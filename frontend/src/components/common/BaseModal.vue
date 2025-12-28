<template>
  <!-- Modal Backdrop - iOS-inspired smooth animation -->
  <Teleport to="body">
    <Motion
      v-if="isOpen"
      :initial="{ opacity: 0 }"
      :animate="{ opacity: 1 }"
      :exit="{ opacity: 0 }"
      :transition="{ duration: 0.2, ease: 'easeOut' }"
      class="modal-backdrop"
      @click.self="handleBackdropClick"
    >
      <!-- Modal Container - iOS-like spring physics -->
      <Motion
        :initial="mobileView 
          ? { opacity: 0, y: '100%' } 
          : { opacity: 0, scale: 0.95, y: 20 }"
        :animate="mobileView 
          ? { opacity: 1, y: 0 } 
          : { opacity: 1, scale: 1, y: 0 }"
        :exit="mobileView 
          ? { opacity: 0, y: '100%' } 
          : { opacity: 0, scale: 0.95, y: 20 }"
        :transition="{ 
          type: 'spring',
          stiffness: 500,
          damping: 40,
          mass: 0.8
        }"
        :class="[
          'modal-container',
          mobileView ? 'modal-mobile' : 'modal-desktop',
          sizeClasses
        ]"
        @click.stop
      >
        <!-- Header dengan Gradient Accent Line -->
        <div v-if="showHeader" class="modal-header">
          <!-- Mobile Drag Handle -->
          <div 
            v-if="mobileView && dismissible" 
            class="modal-drag-handle"
            @click="handleClose"
          />

          <div class="flex items-start justify-between">
            <div class="flex-1 pr-4">
              <!-- Title dengan Optional Gradient -->
              <h2 
                v-if="title"
                :class="[
                  'modal-title',
                  titleGradient ? 'modal-title-gradient' : ''
                ]"
              >
                {{ title }}
              </h2>
              <slot name="title" />

              <!-- Subtitle -->
              <p v-if="subtitle" class="modal-subtitle">
                {{ subtitle }}
              </p>
              <slot name="subtitle" />
            </div>

            <!-- Close Button dengan Press Feedback -->
            <button
              v-if="dismissible"
              @click="handleClose"
              class="modal-close-btn"
              aria-label="Tutup"
            >
              <X />
            </button>
          </div>
        </div>

        <!-- Content Area dengan Custom Scrollbar -->
        <div 
          :class="[
            'modal-content custom-scrollbar',
            scrollable ? 'overflow-y-auto' : '',
            noPadding ? 'p-0!' : ''
          ]"
          :style="{ maxHeight: contentMaxHeight }"
        >
          <slot />
        </div>

        <!-- Footer dengan Glass Effect -->
        <div v-if="showFooter" class="modal-footer">
          <slot name="footer">
            <!-- Default Footer Actions -->
            <div class="modal-actions">
              <button
                v-if="showCancel"
                @click="handleCancel"
                :disabled="loading"
                class="modal-btn-secondary"
              >
                {{ cancelText }}
              </button>
              
              <button
                v-if="showConfirm"
                @click="handleConfirm"
                :disabled="loading || confirmDisabled"
                :class="confirmDanger ? 'modal-btn-danger' : 'modal-btn-primary'"
              >
                <span v-if="loading" class="flex items-center justify-center gap-2">
                  <div class="modal-spinner" />
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
 * yang mengimplementasikan spring physics animations, glass morphism,
 * dan Indigo-Fuchsia gradient theme untuk consistent UX
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
    default: 'md',
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

  // Haptic feedback untuk mobile
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
const mobileView = computed(() => width.value < 640)

const sizeClasses = computed(() => {
  const sizes = {
    xs: 'modal-xs',
    sm: 'modal-sm',
    md: 'modal-md',
    lg: 'modal-lg',
    xl: 'modal-xl',
    full: 'modal-full'
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
  document.body.style.overflow = 'hidden'
  triggerHaptic('light')
  emit('opened')
}

const handleClose = () => {
  if (!props.dismissible || props.loading) return
  
  isOpen.value = false
  triggerHaptic('light')
  
  setTimeout(() => {
    document.body.style.overflow = ''
    emit('update:modelValue', false)
    emit('close')
    emit('closed')
  }, 250)
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
