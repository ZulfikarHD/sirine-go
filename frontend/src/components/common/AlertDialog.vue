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
    <!-- Icon dengan Motion-V Animation -->
    <Motion
      v-bind="iconAnimations.popIn"
      class="flex justify-center mb-6"
    >
      <div :class="['dialog-icon-container', iconContainerClass]">
        <component 
          :is="iconComponent" 
          class="w-10 h-10" 
          :stroke-width="1.5"
        />
      </div>
    </Motion>

    <!-- Message Content -->
    <Motion
      v-bind="entranceAnimations.fadeUp"
      class="text-center space-y-2"
    >
      <p class="text-xl font-bold text-gray-900">
        {{ message }}
      </p>

      <p v-if="detail" class="text-sm text-gray-500 leading-relaxed">
        {{ detail }}
      </p>

      <slot />
    </Motion>

    <!-- Additional Content Card -->
    <Motion
      v-if="showContent"
      :initial="{ opacity: 0, y: 10 }"
      :animate="{ opacity: 1, y: 0 }"
      :transition="{ duration: 0.25, delay: 0.15, ease: 'easeOut' }"
      :class="['mt-6 p-4 rounded-2xl border', contentBoxClass]"
    >
      <slot name="content" />
    </Motion>

    <!-- Auto-dismiss Timer -->
    <Motion
      v-if="autoDismiss && showTimer"
      :initial="{ opacity: 0 }"
      :animate="{ opacity: 1 }"
      :transition="{ duration: 0.2, delay: 0.2 }"
      class="dialog-timer"
    >
      <div class="dialog-timer-ring">
        <svg viewBox="0 0 20 20">
          <circle
            cx="10"
            cy="10"
            r="8"
            stroke="currentColor"
            stroke-width="2"
            fill="none"
            class="opacity-20"
          />
          <circle
            cx="10"
            cy="10"
            r="8"
            :stroke="timerStrokeColor"
            stroke-width="2"
            fill="none"
            stroke-linecap="round"
            :style="{
              strokeDasharray: `${2 * Math.PI * 8}`,
              strokeDashoffset: `${2 * Math.PI * 8 * (1 - timerProgress)}`,
              transition: 'stroke-dashoffset 0.1s linear'
            }"
          />
        </svg>
      </div>
      <span>Menutup dalam {{ remainingTime }}s</span>
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
import { entranceAnimations, iconAnimations } from '../../composables/useMotion'

/**
 * AlertDialog - Komponen alert/notification dialog dengan Motion-V animations
 * untuk success, error, warning, dan info states
 */

const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false
  },
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
  variant: {
    type: String,
    default: 'success',
    validator: (value) => ['success', 'error', 'warning', 'info'].includes(value)
  },
  icon: {
    type: [String, Object],
    default: null
  },
  confirmText: {
    type: String,
    default: 'OK'
  },
  autoDismiss: {
    type: Boolean,
    default: false
  },
  autoDismissDelay: {
    type: Number,
    default: 3000
  },
  showTimer: {
    type: Boolean,
    default: true
  },
  showContent: {
    type: Boolean,
    default: false
  },
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

// Icon configuration
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

const iconContainerClass = computed(() => {
  const classes = {
    success: 'dialog-icon-success',
    error: 'dialog-icon-danger',
    warning: 'dialog-icon-warning',
    info: 'dialog-icon-info'
  }
  return classes[props.variant] || classes.success
})

const contentBoxClass = computed(() => {
  const classes = {
    success: 'bg-emerald-50/80 border-emerald-200/50',
    error: 'bg-red-50/80 border-red-200/50',
    warning: 'bg-amber-50/80 border-amber-200/50',
    info: 'bg-blue-50/80 border-blue-200/50'
  }
  return classes[props.variant] || classes.success
})

const timerStrokeColor = computed(() => {
  const colors = {
    success: '#10b981',
    error: '#ef4444',
    warning: '#f59e0b',
    info: '#3b82f6'
  }
  return colors[props.variant] || colors.success
})

// Auto dismiss
const startAutoDismiss = () => {
  clearTimers()
  
  remainingTime.value = Math.ceil(props.autoDismissDelay / 1000)
  timerProgress.value = 0

  const totalSeconds = Math.ceil(props.autoDismissDelay / 1000)
  
  countdownInterval = setInterval(() => {
    remainingTime.value -= 1
    timerProgress.value = 1 - (remainingTime.value / totalSeconds)
    
    if (remainingTime.value <= 0) {
      clearInterval(countdownInterval)
    }
  }, 1000)

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

const handleClose = () => {
  clearTimers()
  isOpen.value = false
  emit('close')
}

onMounted(() => {
  if (props.modelValue && props.autoDismiss) {
    startAutoDismiss()
  }
})

onUnmounted(() => {
  clearTimers()
})
</script>
