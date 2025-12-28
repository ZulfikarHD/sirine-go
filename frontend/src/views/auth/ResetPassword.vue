<template>
  <div class="min-h-screen flex items-center justify-center bg-gradient-to-br from-indigo-50 via-white to-fuchsia-50 p-4">
    <Motion
      v-bind="entranceAnimations.fadeScale"
      class="w-full max-w-md bg-white/95 backdrop-blur-sm rounded-2xl shadow-xl border border-gray-200/30 p-8 space-y-6"
    >
      <!-- Icon -->
      <Motion v-bind="iconAnimations.popIn" class="w-16 h-16 mx-auto bg-gradient-to-br from-indigo-100 to-fuchsia-100 rounded-full flex items-center justify-center">
        <svg class="w-10 h-10 text-indigo-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
        </svg>
      </Motion>

      <!-- Title -->
      <div class="text-center space-y-2">
        <h1 class="text-2xl font-bold text-gray-900">Reset Password</h1>
        <p class="text-sm text-gray-600">
          Buat password baru untuk akun Anda
        </p>
      </div>

      <!-- Success Message -->
      <Motion
        v-if="isSuccess"
        :initial="{ opacity: 0, y: -10 }"
        :animate="{ opacity: 1, y: 0 }"
        :transition="springPresets.snappy"
        class="flex gap-3 p-4 rounded-xl bg-green-50 border border-green-200"
      >
        <svg class="w-6 h-6 text-green-600 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
        <div>
          <p class="font-semibold text-green-900">Password Berhasil Direset!</p>
          <p class="text-sm text-green-700 mt-1">
            Anda akan diarahkan ke halaman login dalam beberapa detik...
          </p>
        </div>
      </Motion>

      <!-- Form -->
      <form v-if="!isSuccess" @submit.prevent="handleSubmit" class="space-y-4">
        <!-- New Password Input -->
        <div class="space-y-2">
          <label for="newPassword" class="block text-sm font-medium text-gray-700">Password Baru</label>
          <div class="relative">
            <input
              id="newPassword"
              v-model="form.newPassword"
              :type="showPassword ? 'text' : 'password'"
              placeholder="Masukkan password baru"
              class="w-full px-4 py-3 pr-12 rounded-xl border border-gray-300 focus:outline-none focus:ring-4 focus:ring-indigo-100 focus:border-indigo-500 transition-all duration-200 disabled:bg-gray-50 disabled:cursor-not-allowed"
              :class="errors.newPassword ? 'border-red-500 focus:ring-red-100 focus:border-red-500' : ''"
              :disabled="isLoading"
              required
            />
            <button
              type="button"
              class="absolute right-3 top-1/2 -translate-y-1/2 text-gray-400 hover:text-gray-600 transition-colors duration-200"
              @click="showPassword = !showPassword"
            >
              <svg v-if="showPassword" class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
              </svg>
              <svg v-else class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.88 9.88l-3.29-3.29m7.532 7.532l3.29 3.29M3 3l3.59 3.59m0 0A9.953 9.953 0 0112 5c4.478 0 8.268 2.943 9.543 7a10.025 10.025 0 01-4.132 5.411m0 0L21 21" />
              </svg>
            </button>
          </div>
          <p v-if="errors.newPassword" class="text-sm text-red-600">
            {{ errors.newPassword }}
          </p>
        </div>

        <!-- Password Strength Indicator -->
        <PasswordStrength
          :password="form.newPassword"
          :show-requirements="true"
        />

        <!-- Confirm Password Input -->
        <div class="space-y-2">
          <label for="confirmPassword" class="block text-sm font-medium text-gray-700">Konfirmasi Password</label>
          <div class="relative">
            <input
              id="confirmPassword"
              v-model="form.confirmPassword"
              :type="showConfirmPassword ? 'text' : 'password'"
              placeholder="Masukkan ulang password baru"
              class="w-full px-4 py-3 pr-12 rounded-xl border border-gray-300 focus:outline-none focus:ring-4 focus:ring-indigo-100 focus:border-indigo-500 transition-all duration-200 disabled:bg-gray-50 disabled:cursor-not-allowed"
              :class="errors.confirmPassword ? 'border-red-500 focus:ring-red-100 focus:border-red-500' : ''"
              :disabled="isLoading"
              required
            />
            <button
              type="button"
              class="absolute right-3 top-1/2 -translate-y-1/2 text-gray-400 hover:text-gray-600 transition-colors duration-200"
              @click="showConfirmPassword = !showConfirmPassword"
            >
              <svg v-if="showConfirmPassword" class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
              </svg>
              <svg v-else class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.88 9.88l-3.29-3.29m7.532 7.532l3.29 3.29M3 3l3.59 3.59m0 0A9.953 9.953 0 0112 5c4.478 0 8.268 2.943 9.543 7a10.025 10.025 0 01-4.132 5.411m0 0L21 21" />
              </svg>
            </button>
          </div>
          <p v-if="errors.confirmPassword" class="text-sm text-red-600">
            {{ errors.confirmPassword }}
          </p>
        </div>

        <!-- Error Message -->
        <Motion
          v-if="errorMessage"
          :initial="{ opacity: 0, y: -10 }"
          :animate="{ opacity: 1, y: 0 }"
          :transition="springPresets.snappy"
          class="flex items-center gap-2 p-3 rounded-xl bg-red-50 border border-red-200 text-sm text-red-700"
        >
          <svg class="w-5 h-5 text-red-600 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
          {{ errorMessage }}
        </Motion>

        <!-- Submit Button -->
        <button
          type="submit"
          class="w-full px-6 py-3 rounded-xl bg-gradient-to-r from-indigo-600 to-fuchsia-600 text-white font-semibold hover:from-indigo-700 hover:to-fuchsia-700 focus:outline-none focus:ring-4 focus:ring-indigo-100 transition-all duration-200 disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-2 active:scale-[0.97]"
          :disabled="isLoading"
        >
          <svg
            v-if="isLoading"
            class="w-5 h-5 animate-spin"
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 24 24"
          >
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
          </svg>
          <span v-else>Reset Password</span>
        </button>
      </form>

      <!-- Back to Login Link -->
      <div class="text-center pt-4 border-t border-gray-200">
        <button
          class="text-sm text-indigo-600 hover:text-indigo-700 font-medium transition-colors duration-200 active:scale-[0.97]"
          @click="router.push('/login')"
        >
          Kembali ke Login
        </button>
      </div>
    </Motion>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { Motion } from 'motion-v'
import { entranceAnimations, springPresets, iconAnimations } from '@/composables/useMotion'
import { useAuth } from '@/composables/useAuth'
import PasswordStrength from '@/components/auth/PasswordStrength.vue'

/**
 * ResetPassword page untuk reset password menggunakan token dari email
 * dengan password strength indicator dan validation
 */
const router = useRouter()
const route = useRoute()
const { resetPassword, isLoading, error } = useAuth()

const form = reactive({
  newPassword: '',
  confirmPassword: '',
})

const errors = reactive({
  newPassword: '',
  confirmPassword: '',
})

const errorMessage = ref('')
const isSuccess = ref(false)
const showPassword = ref(false)
const showConfirmPassword = ref(false)
const token = ref('')

/**
 * Get token dari URL query params
 */
onMounted(() => {
  token.value = route.query.token || ''
  
  if (!token.value) {
    errorMessage.value = 'Token reset tidak valid. Silakan request ulang link reset password.'
  }
})

/**
 * Validate form sebelum submit
 */
const validateForm = () => {
  errors.newPassword = ''
  errors.confirmPassword = ''
  errorMessage.value = ''

  if (!form.newPassword) {
    errors.newPassword = 'Password baru harus diisi'
    return false
  }

  if (form.newPassword.length < 8) {
    errors.newPassword = 'Password minimal 8 karakter'
    return false
  }

  if (!form.confirmPassword) {
    errors.confirmPassword = 'Konfirmasi password harus diisi'
    return false
  }

  if (form.newPassword !== form.confirmPassword) {
    errors.confirmPassword = 'Password tidak cocok'
    return false
  }

  return true
}

/**
 * Handle form submit
 */
const handleSubmit = async () => {
  if (!validateForm()) return

  if (!token.value) {
    errorMessage.value = 'Token reset tidak valid'
    return
  }

  try {
    await resetPassword(token.value, form.newPassword)
    isSuccess.value = true
    
    // Auto redirect ke login setelah 3 detik
    setTimeout(() => {
      router.push('/login')
    }, 3000)
  } catch (err) {
    errorMessage.value = error.value || 'Gagal reset password'
  }
}
</script>
