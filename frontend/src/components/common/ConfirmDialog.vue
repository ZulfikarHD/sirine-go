<template>
  <BaseModal
    v-model="isOpen"
    :title="title"
    :subtitle="subtitle"
    size="sm"
    :show-footer="true"
    :show-cancel="true"
    :show-confirm="true"
    :cancel-text="cancelText"
    :confirm-text="confirmText"
    :loading="loading"
    :loading-text="loadingText"
    :confirm-danger="variant === 'danger'"
    :dismissible="!loading"
    @confirm="handleConfirm"
    @cancel="handleCancel"
    @close="handleClose"
  >
    <!-- Icon dengan Animated Scale -->
    <Motion
      :initial="{ scale: 0, rotate: -180 }"
      :animate="{ scale: 1, rotate: 0 }"
      :transition="{ 
        duration: 0.5,
        type: 'spring',
        stiffness: 200,
        damping: 15
      }"
      class="flex justify-center mb-6"
    >
      <div 
        :class="[
          'w-20 h-20 rounded-full flex items-center justify-center',
          iconClasses
        ]"
      >
        <component :is="iconComponent" :class="['w-10 h-10', iconColorClasses]" />
      </div>
    </Motion>

    <!-- Message Content -->
    <div class="text-center space-y-3">
      <p v-if="message" class="text-base text-gray-700 leading-relaxed">
        {{ message }}
      </p>
      <slot />

      <!-- Detail Message (opsional) -->
      <p v-if="detail" class="text-sm text-gray-500">
        {{ detail }}
      </p>
    </div>

    <!-- Warning Badge (untuk destructive actions) -->
    <div 
      v-if="variant === 'danger' && showWarning"
      class="mt-6 p-4 bg-red-50 border border-red-200 rounded-xl"
    >
      <div class="flex items-start gap-3">
        <AlertTriangle class="w-5 h-5 text-red-600 shrink-0 mt-0.5" />
        <div class="flex-1">
          <p class="text-sm font-semibold text-red-900">Peringatan</p>
          <p class="text-sm text-red-700 mt-1">
            {{ warningMessage || 'Tindakan ini tidak dapat dibatalkan.' }}
          </p>
        </div>
      </div>
    </div>
  </BaseModal>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { Motion } from 'motion-v'
import { 
  AlertCircle, 
  AlertTriangle, 
  HelpCircle, 
  Trash2,
  CheckCircle2
} from 'lucide-vue-next'
import BaseModal from './BaseModal.vue'

/**
 * ConfirmDialog - Komponen confirmation dialog dengan visual feedback
 * untuk mengkonfirmasi user actions seperti delete, submit, atau cancel
 * operations dengan clear visual hierarchy dan haptic feedback
 */

const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false
  },

  // Content props
  title: {
    type: String,
    default: 'Konfirmasi'
  },
  subtitle: {
    type: String,
    default: ''
  },
  message: {
    type: String,
    default: ''
  },
  detail: {
    type: String,
    default: ''
  },

  // Variant untuk different types of confirmation
  variant: {
    type: String,
    default: 'default', // default, danger, warning, info, success
    validator: (value) => ['default', 'danger', 'warning', 'info', 'success'].includes(value)
  },

  // Icon customization
  icon: {
    type: [String, Object],
    default: null
  },

  // Button texts
  confirmText: {
    type: String,
    default: 'Ya, Lanjutkan'
  },
  cancelText: {
    type: String,
    default: 'Batal'
  },
  loadingText: {
    type: String,
    default: 'Memproses...'
  },

  // Warning props (untuk danger variant)
  showWarning: {
    type: Boolean,
    default: true
  },
  warningMessage: {
    type: String,
    default: ''
  },

  // Loading state
  loading: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['update:modelValue', 'confirm', 'cancel', 'close'])

const isOpen = ref(props.modelValue)

// Watch modelValue
watch(() => props.modelValue, (newValue) => {
  isOpen.value = newValue
})

watch(isOpen, (newValue) => {
  emit('update:modelValue', newValue)
})

// Icon configuration berdasarkan variant
const iconComponent = computed(() => {
  if (props.icon) return props.icon

  const icons = {
    default: HelpCircle,
    danger: Trash2,
    warning: AlertTriangle,
    info: AlertCircle,
    success: CheckCircle2
  }
  return icons[props.variant] || icons.default
})

const iconClasses = computed(() => {
  const classes = {
    default: 'bg-indigo-100',
    danger: 'bg-red-100',
    warning: 'bg-amber-100',
    info: 'bg-blue-100',
    success: 'bg-emerald-100'
  }
  return classes[props.variant] || classes.default
})

const iconColorClasses = computed(() => {
  const classes = {
    default: 'text-indigo-600',
    danger: 'text-red-600',
    warning: 'text-amber-600',
    info: 'text-blue-600',
    success: 'text-emerald-600'
  }
  return classes[props.variant] || classes.default
})

// Event handlers
const handleConfirm = async () => {
  emit('confirm')
}

const handleCancel = () => {
  emit('cancel')
  handleClose()
}

const handleClose = () => {
  isOpen.value = false
  emit('close')
}
</script>
