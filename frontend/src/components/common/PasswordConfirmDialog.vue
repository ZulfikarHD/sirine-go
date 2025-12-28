<template>
  <BaseModal
    v-model="isOpen"
    :title="title"
    :subtitle="subtitle"
    title-gradient
    size="sm"
    :show-footer="true"
    :show-cancel="true"
    :show-confirm="true"
    cancel-text="Batal"
    :confirm-text="confirmText"
    :loading="loading"
    loading-text="Memverifikasi..."
    :confirm-disabled="!password || password.length < 1"
    :dismissible="!loading"
    @confirm="handleConfirm"
    @cancel="handleCancel"
    @close="handleClose"
  >
    <!-- Security Icon -->
    <Motion
      v-bind="iconAnimations.popIn"
      class="flex justify-center mb-6"
    >
      <div class="dialog-icon-container dialog-icon-default">
        <ShieldCheck class="w-10 h-10" :stroke-width="1.5" />
      </div>
    </Motion>

    <!-- Message Content -->
    <Motion
      v-bind="entranceAnimations.fadeUp"
      class="text-center space-y-2 mb-6"
    >
      <p v-if="message" class="text-base text-gray-700 leading-relaxed">
        {{ message }}
      </p>
      <p v-if="detail" class="text-sm text-gray-500">
        {{ detail }}
      </p>
    </Motion>

    <!-- Password Input Field -->
    <Motion
      :initial="{ opacity: 0, y: 10 }"
      :animate="{ opacity: 1, y: 0 }"
      :transition="{ duration: 0.25, delay: 0.1, ease: 'easeOut' }"
    >
      <div class="space-y-2">
        <label class="block text-sm font-semibold text-gray-700">
          Password <span class="text-fuchsia-500">*</span>
        </label>
        <div class="password-input-wrapper">
          <input
            ref="passwordInput"
            v-model="password"
            :type="showPassword ? 'text' : 'password'"
            class="password-input-field"
            placeholder="Masukkan password Anda"
            autocomplete="current-password"
            :disabled="loading"
            @keyup.enter="handleConfirm"
          />
          <button
            type="button"
            @click="togglePasswordVisibility"
            class="password-toggle-btn"
            :disabled="loading"
          >
            <EyeOff v-if="showPassword" />
            <Eye v-else />
          </button>
        </div>
        <p class="text-xs text-gray-400">
          Konfirmasi password Anda untuk melanjutkan
        </p>
      </div>
    </Motion>

    <!-- Error Message -->
    <Motion
      v-if="errorMessage"
      :initial="{ opacity: 0, x: -10 }"
      :animate="{ opacity: 1, x: 0 }"
      :transition="{ ...springPresets.snappy }"
      class="dialog-error-box"
    >
      <div class="flex items-start gap-3">
        <AlertCircle class="w-5 h-5 text-red-500 shrink-0 mt-0.5" />
        <div class="flex-1">
          <p class="text-sm font-semibold text-red-800">Verifikasi Gagal</p>
          <p class="text-sm text-red-600 mt-1">
            {{ errorMessage }}
          </p>
        </div>
      </div>
    </Motion>

    <!-- Security Notice -->
    <Motion
      v-if="showSecurityNotice"
      :initial="{ opacity: 0 }"
      :animate="{ opacity: 1 }"
      :transition="{ duration: 0.2, delay: 0.15 }"
      class="dialog-info-box"
    >
      <div class="flex items-start gap-3">
        <Info class="dialog-info-box-icon" />
        <p class="dialog-info-box-text">
          {{ securityNotice || 'Password Anda tidak akan disimpan dan hanya digunakan untuk verifikasi.' }}
        </p>
      </div>
    </Motion>
  </BaseModal>
</template>

<script setup>
import { ref, watch, nextTick } from 'vue'
import { Motion } from 'motion-v'
import { 
  ShieldCheck, 
  Eye, 
  EyeOff, 
  AlertCircle,
  Info
} from 'lucide-vue-next'
import BaseModal from './BaseModal.vue'
import { entranceAnimations, iconAnimations, springPresets } from '../../composables/useMotion'

/**
 * PasswordConfirmDialog - Komponen untuk password confirmation
 * dengan security-focused UX dan Motion-V animations
 */

const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false
  },
  title: {
    type: String,
    default: 'Konfirmasi Password'
  },
  subtitle: {
    type: String,
    default: ''
  },
  message: {
    type: String,
    default: 'Untuk keamanan, masukkan password Anda untuk melanjutkan.'
  },
  detail: {
    type: String,
    default: ''
  },
  confirmText: {
    type: String,
    default: 'Verifikasi'
  },
  showSecurityNotice: {
    type: Boolean,
    default: true
  },
  securityNotice: {
    type: String,
    default: ''
  },
  validatePassword: {
    type: Function,
    default: null
  },
  autoClearPassword: {
    type: Boolean,
    default: true
  }
})

const emit = defineEmits(['update:modelValue', 'confirm', 'cancel', 'close'])

const isOpen = ref(props.modelValue)
const loading = ref(false)
const password = ref('')
const showPassword = ref(false)
const errorMessage = ref('')
const passwordInput = ref(null)

watch(() => props.modelValue, (newValue) => {
  isOpen.value = newValue
  if (newValue) {
    nextTick(() => {
      passwordInput.value?.focus()
    })
    errorMessage.value = ''
    if (props.autoClearPassword) {
      password.value = ''
      showPassword.value = false
    }
  }
})

watch(isOpen, (newValue) => {
  emit('update:modelValue', newValue)
})

// Toggle password visibility dengan haptic
const togglePasswordVisibility = () => {
  showPassword.value = !showPassword.value
  if ('vibrate' in navigator) {
    navigator.vibrate(10)
  }
}

// Event handlers
const handleConfirm = async () => {
  if (!password.value || loading.value) return

  errorMessage.value = ''
  loading.value = true

  try {
    if (props.validatePassword) {
      const isValid = await props.validatePassword(password.value)
      if (!isValid) {
        errorMessage.value = 'Password yang Anda masukkan tidak valid.'
        loading.value = false
        if ('vibrate' in navigator) {
          navigator.vibrate([30, 100, 30])
        }
        return
      }
    }

    emit('confirm', password.value)
  } catch (error) {
    errorMessage.value = error.message || 'Terjadi kesalahan saat verifikasi.'
    loading.value = false
    if ('vibrate' in navigator) {
      navigator.vibrate([30, 100, 30])
    }
  }
}

const handleCancel = () => {
  if (loading.value) return
  
  errorMessage.value = ''
  if (props.autoClearPassword) {
    password.value = ''
  }
  emit('cancel')
  handleClose()
}

const handleClose = () => {
  if (loading.value) return
  
  isOpen.value = false
  errorMessage.value = ''
  if (props.autoClearPassword) {
    password.value = ''
    showPassword.value = false
  }
  emit('close')
}

// Public methods
const setError = (message) => {
  errorMessage.value = message
  loading.value = false
  if ('vibrate' in navigator) {
    navigator.vibrate([30, 100, 30])
  }
}

const setLoading = (value) => {
  loading.value = value
}

const close = () => {
  loading.value = false
  handleClose()
}

defineExpose({
  setError,
  setLoading,
  close
})
</script>
