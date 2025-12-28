<template>
  <Teleport to="body">
    <div class="fixed inset-0 z-[9999] bg-gradient-to-br from-indigo-50 via-white to-fuchsia-50 flex items-center justify-center p-4 overflow-y-auto">
      <Motion
        v-bind="entranceAnimations.fadeScale"
        class="w-full max-w-lg bg-white rounded-2xl shadow-2xl border border-gray-200 p-8 space-y-6 my-8"
      >
        <!-- Icon -->
        <Motion v-bind="iconAnimations.popIn" class="w-20 h-20 mx-auto bg-gradient-to-br from-yellow-100 to-orange-100 rounded-full flex items-center justify-center">
          <svg class="w-12 h-12 text-orange-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
          </svg>
        </Motion>

        <!-- Header -->
        <div class="text-center space-y-2">
          <h1 class="text-2xl font-bold text-gray-900">Ubah Password Anda</h1>
          <p class="text-sm text-gray-600 leading-relaxed">
            Untuk keamanan akun Anda, Anda harus mengubah password sebelum melanjutkan
          </p>
        </div>

        <!-- Info Alert -->
        <div class="flex gap-3 p-4 rounded-xl bg-blue-50 border border-blue-200">
          <svg class="w-6 h-6 text-blue-600 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
          <p class="text-sm text-blue-700 leading-relaxed">
            Ini adalah login pertama Anda atau password Anda telah direset oleh administrator. Silakan buat password baru yang kuat.
          </p>
        </div>

        <!-- Form -->
        <form @submit.prevent="handleSubmit" class="space-y-5">
          <!-- Current Password (from admin) -->
          <div class="space-y-2">
            <label for="currentPassword" class="block text-sm font-medium text-gray-700">Password Sementara</label>
            <div class="relative">
              <input
                id="currentPassword"
                v-model="form.currentPassword"
                :type="showCurrentPassword ? 'text' : 'password'"
                placeholder="Password yang diberikan admin"
                class="w-full px-4 py-3 pr-12 rounded-xl border border-gray-300 focus:outline-none focus:ring-4 focus:ring-indigo-100 focus:border-indigo-500 transition-all duration-200 disabled:bg-gray-50 disabled:cursor-not-allowed"
                :class="errors.currentPassword ? 'border-red-500 focus:ring-red-100 focus:border-red-500' : ''"
                :disabled="isLoading"
                required
              />
              <button
                type="button"
                class="absolute right-3 top-1/2 -translate-y-1/2 text-gray-400 hover:text-gray-600 transition-colors duration-200"
                @click="showCurrentPassword = !showCurrentPassword"
              >
                <svg v-if="showCurrentPassword" class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
                </svg>
                <svg v-else class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.88 9.88l-3.29-3.29m7.532 7.532l3.29 3.29M3 3l3.59 3.59m0 0A9.953 9.953 0 0112 5c4.478 0 8.268 2.943 9.543 7a10.025 10.025 0 01-4.132 5.411m0 0L21 21" />
                </svg>
              </button>
            </div>
            <p v-if="errors.currentPassword" class="text-sm text-red-600">
              {{ errors.currentPassword }}
            </p>
          </div>

          <!-- New Password -->
          <div class="space-y-2">
            <label for="newPassword" class="block text-sm font-medium text-gray-700">Password Baru</label>
            <div class="relative">
              <input
                id="newPassword"
                v-model="form.newPassword"
                :type="showNewPassword ? 'text' : 'password'"
                placeholder="Buat password baru yang kuat"
                class="w-full px-4 py-3 pr-12 rounded-xl border border-gray-300 focus:outline-none focus:ring-4 focus:ring-indigo-100 focus:border-indigo-500 transition-all duration-200 disabled:bg-gray-50 disabled:cursor-not-allowed"
                :class="errors.newPassword ? 'border-red-500 focus:ring-red-100 focus:border-red-500' : ''"
                :disabled="isLoading"
                required
              />
              <button
                type="button"
                class="absolute right-3 top-1/2 -translate-y-1/2 text-gray-400 hover:text-gray-600 transition-colors duration-200"
                @click="showNewPassword = !showNewPassword"
              >
                <svg v-if="showNewPassword" class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
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

          <!-- Confirm Password -->
          <div class="space-y-2">
            <label for="confirmPassword" class="block text-sm font-medium text-gray-700">Konfirmasi Password Baru</label>
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
            <span v-else>Ubah Password & Lanjutkan</span>
          </button>
        </form>
      </Motion>
    </div>
  </Teleport>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { Motion } from 'motion-v'
import { entranceAnimations, springPresets, iconAnimations } from '@/composables/useMotion'
import { useAuth } from '@/composables/useAuth'
import PasswordStrength from '@/components/auth/PasswordStrength.vue'

/**
 * ForceChangePassword modal untuk first-time login atau admin reset
 * Fullscreen blocking modal yang tidak bisa di-dismiss
 */
const router = useRouter()
const { changePassword, isLoading, error } = useAuth()

const form = reactive({
  currentPassword: '',
  newPassword: '',
  confirmPassword: '',
})

const errors = reactive({
  currentPassword: '',
  newPassword: '',
  confirmPassword: '',
})

const errorMessage = ref('')
const showCurrentPassword = ref(false)
const showNewPassword = ref(false)
const showConfirmPassword = ref(false)

/**
 * Validate form sebelum submit
 */
const validateForm = () => {
  errors.currentPassword = ''
  errors.newPassword = ''
  errors.confirmPassword = ''
  errorMessage.value = ''

  if (!form.currentPassword) {
    errors.currentPassword = 'Password sementara harus diisi'
    return false
  }

  if (!form.newPassword) {
    errors.newPassword = 'Password baru harus diisi'
    return false
  }

  if (form.newPassword.length < 8) {
    errors.newPassword = 'Password minimal 8 karakter'
    return false
  }

  if (form.currentPassword === form.newPassword) {
    errors.newPassword = 'Password baru harus berbeda dari password sementara'
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

  try {
    await changePassword(form.currentPassword, form.newPassword)
    // Success - will auto logout and redirect to login
  } catch (err) {
    errorMessage.value = error.value || 'Gagal mengubah password'
  }
}
</script>
