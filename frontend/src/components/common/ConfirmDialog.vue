<template>
  <BaseModal
    v-model="isOpen"
    :title="title"
    :subtitle="subtitle"
    :title-gradient="variant === 'default'"
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
    <!-- Icon dengan Motion-V animation -->
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
      <p v-if="message" class="text-lg text-gray-800 leading-relaxed font-medium">
        {{ message }}
      </p>
      
      <slot />

      <p v-if="detail" class="text-sm text-gray-500 leading-relaxed">
        {{ detail }}
      </p>
    </Motion>

    <!-- Warning Box untuk Destructive Actions -->
    <Motion
      v-if="variant === 'danger' && showWarning"
      :initial="{ opacity: 0, y: 10 }"
      :animate="{ opacity: 1, y: 0 }"
      :transition="{ duration: 0.25, delay: 0.15, ease: 'easeOut' }"
      class="dialog-warning-box"
    >
      <div class="flex items-start gap-3">
        <AlertTriangle class="dialog-warning-box-icon" />
        <div class="flex-1">
          <p class="dialog-warning-box-title">Peringatan</p>
          <p class="dialog-warning-box-text">
            {{ warningMessage || 'Tindakan ini tidak dapat dibatalkan.' }}
          </p>
        </div>
      </div>
    </Motion>
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
import { entranceAnimations, iconAnimations } from '../../composables/useMotion'

/**
 * ConfirmDialog - Komponen confirmation dialog dengan iOS-inspired design
 * yang mengimplementasikan Motion-V animations untuk smooth UX
 */

const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false
  },
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
  variant: {
    type: String,
    default: 'default',
    validator: (value) => ['default', 'danger', 'warning', 'info', 'success'].includes(value)
  },
  icon: {
    type: [String, Object],
    default: null
  },
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
  showWarning: {
    type: Boolean,
    default: true
  },
  warningMessage: {
    type: String,
    default: ''
  },
  loading: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['update:modelValue', 'confirm', 'cancel', 'close'])

const isOpen = ref(props.modelValue)

watch(() => props.modelValue, (newValue) => {
  isOpen.value = newValue
})

watch(isOpen, (newValue) => {
  emit('update:modelValue', newValue)
})

// Icon configuration
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

const iconContainerClass = computed(() => {
  const classes = {
    default: 'dialog-icon-default',
    danger: 'dialog-icon-danger',
    warning: 'dialog-icon-warning',
    info: 'dialog-icon-info',
    success: 'dialog-icon-success'
  }
  return classes[props.variant] || classes.default
})

// Event handlers
const handleConfirm = () => {
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
