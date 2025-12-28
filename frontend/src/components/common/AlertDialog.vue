<template>
  <BaseModal
    v-model="isOpen"
    :title="title"
    size="sm"
    :show-footer="true"
    :show-cancel="false"
    :show-confirm="true"
    :confirm-text="confirmText"
    :dismissible="dismissible"
    @confirm="handleClose"
    @close="handleClose"
  >
    <!-- Animated Icon dengan Spring Effect -->
    <Motion
      :initial="{ scale: 0, rotate: -180 }"
      :animate="{ scale: 1, rotate: 0 }"
      :transition="{ 
        duration: 0.6,
        type: 'spring',
        stiffness: 200,
        damping: 15
      }"
      class="flex justify-center mb-6"
    >
      <div 
        :class="[
          'relative w-24 h-24 rounded-full flex items-center justify-center',
          backgroundClasses
        ]"
      >
        <!-- Pulse Ring untuk Success -->
        <Motion
          v-if="variant === 'success'"
          :animate="{ 
            scale: [1, 1.3, 1.3],
            opacity: [0.5, 0, 0]
          }"
          :transition="{ 
            duration: 2,
            repeat: Infinity,
            ease: 'easeOut'
          }"
          :class="['absolute inset-0 rounded-full', backgroundClasses]"
        />

        <!-- Icon -->
        <component 
          :is="iconComponent" 
          :class="['w-12 h-12 relative z-10', iconColorClasses]" 
          :stroke-width="2.5"
        />
      </div>
    </Motion>

    <!-- Message Content dengan Staggered Animation -->
    <div class="text-center space-y-3">
      <Motion
        :initial="{ opacity: 0, y: 10 }"
        :animate="{ opacity: 1, y: 0 }"
        :transition="{ duration: 0.4, delay: 0.2 }"
      >
        <p class="text-lg font-semibold text-gray-900">
          {{ message }}
        </p>
      </Motion>

      <Motion
        v-if="detail"
        :initial="{ opacity: 0, y: 10 }"
        :animate="{ opacity: 1, y: 0 }"
        :transition="{ duration: 0.4, delay: 0.3 }"
      >
        <p class="text-sm text-gray-600 leading-relaxed">
          {{ detail }}
        </p>
      </Motion>

      <slot />
    </div>

    <!-- Additional Content Card (opsional) -->
    <Motion
      v-if="showContent"
      :initial="{ opacity: 0, y: 10 }"
      :animate="{ opacity: 1, y: 0 }"
      :transition="{ duration: 0.4, delay: 0.4 }"
      :class="[
        'mt-6 p-4 rounded-xl border',
        contentClasses
      ]"
    >
      <slot name="content" />
    </Motion>

    <!-- Auto-dismiss Timer (visual indicator) -->
    <Motion
      v-if="autoDismiss && showTimer"
      :initial="{ opacity: 0, y: 10 }"
      :animate="{ opacity: 1, y: 0 }"
      :transition="{ duration: 0.4, delay: 0.5 }"
      class="mt-6"
    >
      <div class="flex items-center justify-center gap-2 text-sm text-gray-500">
        <div class="relative w-4 h-4">
          <svg class="w-4 h-4 transform -rotate-90">
            <circle
              cx="8"
              cy="8"
              r="7"
              stroke="currentColor"
              stroke-width="2"
              fill="none"
              class="opacity-20"
            />
            <circle
              cx="8"
              cy="8"
              r="7"
              stroke="currentColor"
              stroke-width="2"
              fill="none"
              :style="{
                strokeDasharray: `${2 * Math.PI * 7}`,
                strokeDashoffset: `${2 * Math.PI * 7 * (1 - timerProgress)}`,
                transition: 'stroke-dashoffset 0.1s linear'
              }"
            />
          </svg>
        </div>
        <span>Menutup otomatis dalam {{ remainingTime }}s</span>
      </div>
    </Motion>
  </BaseModal>
</template>

<script setup>
import { ref, computed, watch, onMounted, onUnmounted } from 'vue'
import { Motion } from 'motion-v'
import { 
  CheckCircle2, 
  AlertCircle, 
  AlertTriangle, 
  Info,
  XCircle
} from 'lucide-vue-next'
import BaseModal from './BaseModal.vue'

/**
 * AlertDialog - Komponen alert/notification dialog untuk menampilkan
 * success, error, warning, atau info messages dengan animated feedback
 * dan auto-dismiss capability untuk better user experience
 */

const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false
  },

  // Content props
  title: {
    type: String,
    default: ''
  },
  message: {
    type: String,
    required: true
  },
  detail: {
    type: String,
    default: ''
  },

  // Alert variant
  variant: {
    type: String,
    default: 'success', // success, error, warning, info
    validator: (value) => ['success', 'error', 'warning', 'info'].includes(value)
  },

  // Icon customization
  icon: {
    type: [String, Object],
    default: null
  },

  // Button text
  confirmText: {
    type: String,
    default: 'OK'
  },

  // Auto dismiss props
  autoDismiss: {
    type: Boolean,
    default: false
  },
  autoDismissDelay: {
    type: Number,
    default: 3000 // 3 seconds
  },
  showTimer: {
    type: Boolean,
    default: true
  },

  // Additional content
  showContent: {
    type: Boolean,
    default: false
  },

  // Dismissible
  dismissible: {
    type: Boolean,
    default: true
  }
})

const emit = defineEmits(['update:modelValue', 'close', 'dismissed'])

const isOpen = ref(props.modelValue)
const remainingTime = ref(0)
const timerProgress = ref(0)
let dismissTimer = null
let countdownInterval = null

// Watch modelValue
watch(() => props.modelValue, (newValue) => {
  isOpen.value = newValue
  if (newValue && props.autoDismiss) {
    startAutoDismiss()
  }
})

watch(isOpen, (newValue) => {
  emit('update:modelValue', newValue)
  if (!newValue) {
    clearTimers()
  }
})

// Icon configuration berdasarkan variant
const iconComponent = computed(() => {
  if (props.icon) return props.icon

  const icons = {
    success: CheckCircle2,
    error: XCircle,
    warning: AlertTriangle,
    info: Info
  }
  return icons[props.variant] || icons.success
})

const backgroundClasses = computed(() => {
  const classes = {
    success: 'bg-emerald-100',
    error: 'bg-red-100',
    warning: 'bg-amber-100',
    info: 'bg-blue-100'
  }
  return classes[props.variant] || classes.success
})

const iconColorClasses = computed(() => {
  const classes = {
    success: 'text-emerald-600',
    error: 'text-red-600',
    warning: 'text-amber-600',
    info: 'text-blue-600'
  }
  return classes[props.variant] || classes.success
})

const contentClasses = computed(() => {
  const classes = {
    success: 'bg-emerald-50 border-emerald-200',
    error: 'bg-red-50 border-red-200',
    warning: 'bg-amber-50 border-amber-200',
    info: 'bg-blue-50 border-blue-200'
  }
  return classes[props.variant] || classes.success
})

// Auto dismiss functionality
const startAutoDismiss = () => {
  clearTimers()
  
  remainingTime.value = Math.ceil(props.autoDismissDelay / 1000)
  timerProgress.value = 0

  // Countdown timer
  countdownInterval = setInterval(() => {
    remainingTime.value -= 1
    timerProgress.value = 1 - (remainingTime.value / Math.ceil(props.autoDismissDelay / 1000))
    
    if (remainingTime.value <= 0) {
      clearInterval(countdownInterval)
    }
  }, 1000)

  // Dismiss timer
  dismissTimer = setTimeout(() => {
    handleClose()
    emit('dismissed')
  }, props.autoDismissDelay)
}

const clearTimers = () => {
  if (dismissTimer) {
    clearTimeout(dismissTimer)
    dismissTimer = null
  }
  if (countdownInterval) {
    clearInterval(countdownInterval)
    countdownInterval = null
  }
}

// Event handlers
const handleClose = () => {
  clearTimers()
  isOpen.value = false
  emit('close')
}

// Lifecycle
onMounted(() => {
  if (props.modelValue && props.autoDismiss) {
    startAutoDismiss()
  }
})

onUnmounted(() => {
  clearTimers()
})
</script>
