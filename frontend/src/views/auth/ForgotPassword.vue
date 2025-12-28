<template>
  <div class="min-h-screen flex items-center justify-center bg-gradient-to-br from-indigo-50 via-white to-fuchsia-50 p-4">
    <Motion
      v-bind="entranceAnimations.fadeScale"
      class="w-full max-w-md bg-white/95 backdrop-blur-sm rounded-2xl shadow-xl border border-gray-200/30 p-8 space-y-6"
    >
      <!-- Back Button -->
      <button
        class="flex items-center gap-2 text-sm text-gray-600 hover:text-gray-900 transition-colors duration-200 -mt-2 -ml-2 p-2 rounded-lg hover:bg-gray-50 active:scale-[0.97]"
        @click="router.push('/login')"
      >
        <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
        </svg>
        Kembali ke Login
      </button>

      <!-- Icon -->
      <Motion v-bind="iconAnimations.popIn" class="w-16 h-16 mx-auto bg-gradient-to-br from-indigo-100 to-fuchsia-100 rounded-full flex items-center justify-center">
        <svg class="w-10 h-10 text-indigo-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 7a2 2 0 012 2m4 0a6 6 0 01-7.743 5.743L11 17H9v2H7v2H4a1 1 0 01-1-1v-2.586a1 1 0 01.293-.707l5.964-5.964A6 6 0 1121 9z" />
        </svg>
      </Motion>

      <!-- Title -->
      <div class="text-center space-y-2">
        <h1 class="text-2xl font-bold text-gray-900">Lupa Password</h1>
        <p class="text-sm text-gray-600 leading-relaxed">
          Masukkan NIP atau Email Anda, kami akan mengirimkan link untuk reset password
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
          <p class="font-semibold text-green-900">Email Terkirim!</p>
          <p class="text-sm text-green-700 mt-1">
            Jika NIP/Email terdaftar, link reset password telah dikirim ke email Anda. Silakan cek inbox atau spam folder.
          </p>
        </div>
      </Motion>

      <!-- Form -->
      <form v-if="!isSuccess" @submit.prevent="handleSubmit" class="space-y-4">
        <!-- NIP or Email Input -->
        <div class="space-y-2">
          <label for="nipOrEmail" class="block text-sm font-medium text-gray-700">NIP atau Email</label>
          <input
            id="nipOrEmail"
            v-model="form.nipOrEmail"
            type="text"
            placeholder="Masukkan NIP atau Email Anda"
            class="w-full px-4 py-3 rounded-xl border border-gray-300 focus:outline-none focus:ring-4 focus:ring-indigo-100 focus:border-indigo-500 transition-all duration-200 disabled:bg-gray-50 disabled:cursor-not-allowed"
            :class="{ 'border-red-500 focus:ring-red-100 focus:border-red-500': errors.nipOrEmail }"
            :disabled="isLoading"
            required
          />
          <p v-if="errors.nipOrEmail" class="text-sm text-red-600">
            {{ errors.nipOrEmail }}
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
          <span v-else>Kirim Link Reset Password</span>
        </button>
      </form>

      <!-- Back to Login (after success) -->
      <button
        v-if="isSuccess"
        class="w-full px-6 py-3 rounded-xl bg-gray-100 hover:bg-gray-200 text-gray-700 font-semibold focus:outline-none focus:ring-4 focus:ring-gray-100 transition-all duration-200 active:scale-[0.97]"
        @click="router.push('/login')"
      >
        Kembali ke Login
      </button>
    </Motion>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { Motion } from 'motion-v'
import { entranceAnimations, springPresets, iconAnimations } from '@/composables/useMotion'
import { useAuth } from '@/composables/useAuth'

/**
 * ForgotPassword page untuk request password reset link
 * dengan email notification system
 */
const router = useRouter()
const { forgotPassword, isLoading, error } = useAuth()

const form = reactive({
  nipOrEmail: '',
})

const errors = reactive({
  nipOrEmail: '',
})

const errorMessage = ref('')
const isSuccess = ref(false)

/**
 * Validate form sebelum submit
 */
const validateForm = () => {
  errors.nipOrEmail = ''
  errorMessage.value = ''

  if (!form.nipOrEmail.trim()) {
    errors.nipOrEmail = 'NIP atau Email harus diisi'
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
    await forgotPassword(form.nipOrEmail)
    isSuccess.value = true
    
    // Auto redirect ke login setelah 5 detik
    setTimeout(() => {
      router.push('/login')
    }, 5000)
  } catch (err) {
    errorMessage.value = error.value || 'Gagal mengirim email reset password'
  }
}
</script>
