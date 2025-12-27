<template>
  <div class="min-h-screen flex items-center justify-center bg-linear-to-br from-indigo-50 via-fuchsia-50 to-indigo-100 p-4">
    <!-- Login Card dengan glass effect -->
    <div 
      class="w-full max-w-md"
      :style="cardAnimation"
    >
      <!-- Glass Card -->
      <div class="glass-card rounded-3xl p-8 shadow-2xl backdrop-blur-xl bg-white/80 border border-white/20">
        <!-- Logo & Title -->
        <div class="text-center mb-8">
          <div class="inline-flex items-center justify-center w-20 h-20 rounded-full bg-linear-to-br from-indigo-500 to-fuchsia-600 mb-4 shadow-lg">
            <svg class="w-10 h-10 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
            </svg>
          </div>
          <h1 class="text-3xl font-bold text-gray-900 mb-2">Sirine Go</h1>
          <p class="text-gray-600">Masuk dengan NIP/Email dan password Anda</p>
        </div>

        <!-- Error Message -->
        <div 
          v-if="errorMessage" 
          class="mb-6 p-4 rounded-2xl bg-red-50 border border-red-200 animate-shake"
        >
          <p class="text-red-600 text-sm text-center font-medium">{{ errorMessage }}</p>
        </div>

        <!-- Login Form -->
        <form @submit.prevent="handleLogin" class="space-y-5">
          <!-- NIP/Email Input -->
          <div>
            <label for="nip" class="block text-sm font-semibold text-gray-700 mb-2">
              NIP atau Email
            </label>
            <input
              id="nip"
              v-model="form.nip"
              type="text"
              placeholder="Masukkan NIP atau Email"
              required
              class="input-field"
              :class="{ 'border-red-300': errors.nip }"
              @focus="clearError('nip')"
            />
            <p v-if="errors.nip" class="mt-1 text-xs text-red-600">{{ errors.nip }}</p>
          </div>

          <!-- Password Input -->
          <div>
            <label for="password" class="block text-sm font-semibold text-gray-700 mb-2">
              Password
            </label>
            <div class="relative">
              <input
                id="password"
                v-model="form.password"
                :type="showPassword ? 'text' : 'password'"
                placeholder="Masukkan password"
                required
                class="input-field pr-12"
                :class="{ 'border-red-300': errors.password }"
                @focus="clearError('password')"
              />
              <button
                type="button"
                @click="showPassword = !showPassword"
                class="absolute right-3 top-1/2 -translate-y-1/2 text-gray-400 hover:text-gray-600 transition-colors active-scale"
              >
                <svg v-if="!showPassword" class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
                </svg>
                <svg v-else class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.88 9.88l-3.29-3.29m7.532 7.532l3.29 3.29M3 3l3.59 3.59m0 0A9.953 9.953 0 0112 5c4.478 0 8.268 2.943 9.543 7a10.025 10.025 0 01-4.132 5.411m0 0L21 21" />
                </svg>
              </button>
            </div>
            <p v-if="errors.password" class="mt-1 text-xs text-red-600">{{ errors.password }}</p>
          </div>

          <!-- Remember Me -->
          <div class="flex items-center">
            <input
              id="remember"
              v-model="form.rememberMe"
              type="checkbox"
              class="w-4 h-4 text-indigo-600 bg-gray-100 border-gray-300 rounded focus:ring-indigo-500 focus:ring-2"
            />
            <label for="remember" class="ml-2 text-sm text-gray-700">
              Ingat saya selama 30 hari
            </label>
          </div>

          <!-- Login Button -->
          <button
            type="submit"
            :disabled="isLoading"
            class="btn-primary w-full"
          >
            <span v-if="!isLoading">Masuk</span>
            <span v-else class="flex items-center justify-center">
              <svg class="animate-spin -ml-1 mr-3 h-5 w-5 text-white" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
              Memproses...
            </span>
          </button>
        </form>

        <!-- Footer -->
        <div class="mt-6 text-center text-sm text-gray-600">
          <p>© 2025 Sirine Go - Sistem Produksi Pita Cukai</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuth } from '../../composables/useAuth'
import { animate, spring } from 'motion-v'

const router = useRouter()
const route = useRoute()
const { login, isLoading, error, getDashboardRoute, triggerHapticFeedback } = useAuth()

// Form state
const form = ref({
  nip: '',
  password: '',
  rememberMe: false,
})

const showPassword = ref(false)
const errors = ref({
  nip: '',
  password: '',
})
const errorMessage = ref('')

// Animations
const cardAnimation = ref({})

onMounted(() => {
  // Spring entrance animation untuk card
  try {
    const card = document.querySelector('.glass-card')
    if (card) {
      animate(
        card,
        { opacity: [0, 1], transform: ['scale(0.95)', 'scale(1)'] },
        { duration: 0.6, easing: spring({ stiffness: 300, damping: 20 }) }
      )
    }
  } catch (error) {
    console.log('Animation not available:', error)
  }
})

/**
 * Clear error untuk field tertentu
 */
const clearError = (field) => {
  errors.value[field] = ''
  errorMessage.value = ''
}

/**
 * Handle login form submission dengan validasi flexible untuk NIP atau Email
 * Memastikan tidak ada full page reload dengan proper error handling
 */
const handleLogin = async (event) => {
  // CRITICAL: Prevent default form submission behavior untuk avoid full reload
  // This is a safeguard even though form has @submit.prevent
  event?.preventDefault()
  event?.stopPropagation()

  // Reset errors
  errors.value = { nip: '', password: '' }
  errorMessage.value = ''

  // Validasi
  if (!form.value.nip || form.value.nip.trim() === '') {
    errors.value.nip = 'NIP atau Email harus diisi'
    return false
  }

  if (!form.value.password) {
    errors.value.password = 'Password harus diisi'
    return false
  }

  try {
    const response = await login(
      form.value.nip.trim(),
      form.value.password,
      form.value.rememberMe
    )

    // Login berhasil - redirect
    const redirectPath = route.query.redirect || getDashboardRoute()
    await router.push(redirectPath)

    // Trigger success haptic
    triggerHapticFeedback('success')

  } catch (err) {
    // Handle error tanpa page reload
    console.error('Login error:', err)
    errorMessage.value = err.response?.data?.message || 'NIP/Email atau password salah'
    
    // Shake animation untuk error
    try {
      const card = document.querySelector('.glass-card')
      if (card) {
        animate(
          card,
          { transform: ['translateX(0)', 'translateX(-10px)', 'translateX(10px)', 'translateX(-10px)', 'translateX(0)'] },
          { duration: 0.4 }
        )
      }
    } catch (animError) {
      console.log('Animation error:', animError)
    }

    // Trigger error haptic
    triggerHapticFeedback('error')
  }
  
  // Ensure we never reload
  return false
}
</script>
