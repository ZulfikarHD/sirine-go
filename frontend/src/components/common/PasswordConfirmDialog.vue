<template>
  <BaseModal
    v-model="isOpen"
    :title="title"
    :subtitle="subtitle"
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
      <div class="w-20 h-20 rounded-full bg-indigo-100 flex items-center justify-center">
        <ShieldCheck class="w-10 h-10 text-indigo-600" />
      </div>
    </Motion>

    <!-- Message Content -->
    <div class="text-center space-y-3 mb-6">
      <p v-if="message" class="text-base text-gray-700 leading-relaxed">
        {{ message }}
      </p>
      <p v-if="detail" class="text-sm text-gray-500">
        {{ detail }}
      </p>
    </div>

    <!-- Password Input -->
    <Motion
      :initial="{ opacity: 0, y: 10 }"
      :animate="{ opacity: 1, y: 0 }"
      :transition="{ duration: 0.4, delay: 0.2 }"
    >
      <div class="space-y-2">
        <label class="block text-sm font-semibold text-gray-700">
          Password <span class="text-red-500">*</span>
        </label>
        <div class="relative">
          <input
            ref="passwordInput"
            v-model="password"
            :type="showPassword ? 'text' : 'password'"
            class="input-field pr-12"
            placeholder="Masukkan password Anda"
            autocomplete="current-password"
            :disabled="loading"
            @keyup.enter="handleConfirm"
          />
          <button
            type="button"
            @click="togglePasswordVisibility"
            class="absolute right-3 top-1/2 -translate-y-1/2 p-2 hover:bg-gray-100 rounded-lg transition-colors active:scale-95"
            :disabled="loading"
          >
            <EyeOff v-if="showPassword" class="w-5 h-5 text-gray-500" />
            <Eye v-else class="w-5 h-5 text-gray-500" />
          </button>
        </div>
        <p class="text-xs text-gray-500">
          Konfirmasi password Anda untuk melanjutkan tindakan ini
        </p>
      </div>
    </Motion>

    <!-- Error Message -->
    <Motion
      v-if="errorMessage"
      :initial="{ opacity: 0, x: -10 }"
      :animate="{ opacity: 1, x: 0 }"
      :transition="{ duration: 0.3 }"
      class="mt-4 p-4 bg-red-50 border border-red-200 rounded-xl"
    >
      <div class="flex items-start gap-3">
        <AlertCircle class="w-5 h-5 text-red-600 shrink-0 mt-0.5" />
        <div class="flex-1">
          <p class="text-sm font-semibold text-red-900">Password Salah</p>
          <p class="text-sm text-red-700 mt-1">
            {{ errorMessage }}
          </p>
        </div>
      </div>
    </Motion>

    <!-- Security Notice (opsional) -->
    <Motion
      v-if="showSecurityNotice"
      :initial="{ opacity: 0, y: 10 }"
      :animate="{ opacity: 1, y: 0 }"
      :transition="{ duration: 0.4, delay: 0.3 }"
      class="mt-4 p-4 bg-blue-50 border border-blue-200 rounded-xl"
    >
      <div class="flex items-start gap-3">
        <Info class="w-5 h-5 text-blue-600 shrink-0 mt-0.5" />
        <p class="text-sm text-blue-700">
          {{ securityNotice || 'Password Anda tidak akan disimpan dan hanya digunakan untuk verifikasi tindakan ini.' }}
        </p>
      </div>
    </Motion>
  </BaseModal>
</template>

<script setup>
import { ref, computed, watch, nextTick } from 'vue'
import { Motion } from 'motion-v'
import { 
  ShieldCheck, 
  Eye, 
  EyeOff, 
  AlertCircle,
  Info
} from 'lucide-vue-next'
import BaseModal from './BaseModal.vue'

/**
 * PasswordConfirmDialog - Komponen untuk password confirmation yang digunakan
 * untuk memverifikasi identitas user sebelum melakukan sensitive actions
 * seperti delete account, change settings, atau approve critical operations
 * dengan security-focused design dan clear error handling
 */

const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false
  },

  // Content props
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

  // Button text
  confirmText: {
    type: String,
    default: 'Verifikasi'
  },

  // Security notice
  showSecurityNotice: {
    type: Boolean,
    default: true
  },
  securityNotice: {
    type: String,
    default: ''
  },

  // Validation function (opsional - untuk client-side check)
  validatePassword: {
    type: Function,
    default: null
  },

  // Auto-clear password on close
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

// Watch modelValue untuk sync
watch(() => props.modelValue, (newValue) => {
  isOpen.value = newValue
  if (newValue) {
    // Focus input saat modal dibuka
    nextTick(() => {
      passwordInput.value?.focus()
    })
    // Clear previous state
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

// Toggle password visibility
const togglePasswordVisibility = () => {
  showPassword.value = !showPassword.value
  // Trigger light haptic
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
    // Client-side validation jika ada
    if (props.validatePassword) {
      const isValid = await props.validatePassword(password.value)
      if (!isValid) {
        errorMessage.value = 'Password yang Anda masukkan tidak valid.'
        loading.value = false
        return
      }
    }

    // Emit password untuk parent component verify
    emit('confirm', password.value)
    
    // Note: Parent component harus handle verification dan close modal
    // Jangan auto-close disini karena parent perlu verify dulu
  } catch (error) {
    errorMessage.value = error.message || 'Terjadi kesalahan saat verifikasi.'
    loading.value = false
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

// Public methods untuk parent component
const setError = (message) => {
  errorMessage.value = message
  loading.value = false
  // Shake animation via haptic
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

// Expose methods
defineExpose({
  setError,
  setLoading,
  close
})
</script>
