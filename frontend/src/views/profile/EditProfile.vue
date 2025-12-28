<template>
  <AppLayout>
    <div class="max-w-2xl mx-auto">
      <!-- Header -->
      <div class="mb-6">
        <button
          @click="goBack"
          class="flex items-center space-x-2 text-gray-600 hover:text-gray-900 transition-colors mb-4"
        >
          <ArrowLeft class="w-5 h-5" />
          <span>Kembali ke Profile</span>
        </button>
        <h1 class="text-2xl font-bold text-gray-900">Edit Profile</h1>
        <p class="text-sm text-gray-600 mt-1">Update informasi profile Anda</p>
      </div>

      <!-- Form -->
      <form @submit.prevent="handleSubmit" class="glass-card p-6 rounded-2xl space-y-6">
        <!-- Full Name -->
        <div>
          <label class="block text-sm font-semibold text-gray-700 mb-2">
            Nama Lengkap <span class="text-red-500">*</span>
          </label>
          <input
            v-model="form.full_name"
            type="text"
            required
            class="w-full px-4 py-2.5 rounded-xl border border-gray-200 focus:border-indigo-500 focus:ring-2 focus:ring-indigo-200 transition-all"
            placeholder="Masukkan nama lengkap"
          />
          <p v-if="errors.full_name" class="mt-1 text-sm text-red-600">{{ errors.full_name }}</p>
        </div>

        <!-- Email -->
        <div>
          <label class="block text-sm font-semibold text-gray-700 mb-2">
            Email <span class="text-red-500">*</span>
          </label>
          <input
            v-model="form.email"
            type="email"
            required
            class="w-full px-4 py-2.5 rounded-xl border border-gray-200 focus:border-indigo-500 focus:ring-2 focus:ring-indigo-200 transition-all"
            placeholder="email@example.com"
          />
          <p v-if="errors.email" class="mt-1 text-sm text-red-600">{{ errors.email }}</p>
        </div>

        <!-- Phone -->
        <div>
          <label class="block text-sm font-semibold text-gray-700 mb-2">
            Nomor Telepon <span class="text-red-500">*</span>
          </label>
          <input
            v-model="form.phone"
            type="tel"
            required
            class="w-full px-4 py-2.5 rounded-xl border border-gray-200 focus:border-indigo-500 focus:ring-2 focus:ring-indigo-200 transition-all"
            placeholder="08xxxxxxxxxx"
          />
          <p v-if="errors.phone" class="mt-1 text-sm text-red-600">{{ errors.phone }}</p>
        </div>

        <!-- Read-only fields -->
        <div class="pt-4 border-t border-gray-200">
          <p class="text-sm font-semibold text-gray-700 mb-3">Informasi yang tidak dapat diubah:</p>
          
          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <!-- NIP -->
            <div>
              <label class="block text-sm text-gray-600 mb-1">NIP</label>
              <div class="px-4 py-2.5 rounded-xl bg-gray-100 text-gray-500 border border-gray-200">
                {{ user?.nip }}
              </div>
            </div>

            <!-- Role -->
            <div>
              <label class="block text-sm text-gray-600 mb-1">Role</label>
              <div class="px-4 py-2.5 rounded-xl bg-gray-100 text-gray-500 border border-gray-200">
                {{ user?.role }}
              </div>
            </div>

            <!-- Department -->
            <div>
              <label class="block text-sm text-gray-600 mb-1">Department</label>
              <div class="px-4 py-2.5 rounded-xl bg-gray-100 text-gray-500 border border-gray-200">
                {{ user?.department }}
              </div>
            </div>

            <!-- Shift -->
            <div>
              <label class="block text-sm text-gray-600 mb-1">Shift</label>
              <div class="px-4 py-2.5 rounded-xl bg-gray-100 text-gray-500 border border-gray-200">
                {{ user?.shift }}
              </div>
            </div>
          </div>

          <p class="mt-3 text-xs text-gray-500">
            * Untuk mengubah NIP, Role, Department, atau Shift, hubungi Administrator
          </p>
        </div>

        <!-- Error Alert -->
        <div v-if="errorMessage" class="p-4 bg-red-50 border border-red-200 rounded-xl">
          <p class="text-sm text-red-700">{{ errorMessage }}</p>
        </div>

        <!-- Actions -->
        <div class="flex items-center space-x-3 pt-4">
          <button
            type="submit"
            :disabled="loading"
            class="flex-1 btn-primary disabled:opacity-50 disabled:cursor-not-allowed"
          >
            <span v-if="loading" class="flex items-center justify-center">
              <div class="w-5 h-5 border-2 border-white border-t-transparent rounded-full animate-spin mr-2"></div>
              Menyimpan...
            </span>
            <span v-else class="flex items-center justify-center">
              <Save class="w-5 h-5 mr-2" />
              Simpan Perubahan
            </span>
          </button>
          <button
            type="button"
            @click="goBack"
            :disabled="loading"
            class="px-6 py-2.5 rounded-xl border border-gray-200 text-gray-700 font-semibold hover:bg-gray-50 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
          >
            Batal
          </button>
        </div>
      </form>
    </div>
  </AppLayout>
</template>

<script setup>
import { ref, reactive, computed, onMounted, nextTick } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../../stores/auth'
import apiClient from '../../composables/useApi'
import AppLayout from '../../components/layout/AppLayout.vue'
import { ArrowLeft, Save } from 'lucide-vue-next'
import { animate } from 'motion-v'

const router = useRouter()
const authStore = useAuthStore()

const user = computed(() => authStore.user)
const loading = ref(false)
const errorMessage = ref('')
const errors = reactive({
  full_name: '',
  email: '',
  phone: ''
})

const form = reactive({
  full_name: '',
  email: '',
  phone: ''
})

// Load current user data
onMounted(async () => {
  if (user.value) {
    form.full_name = user.value.full_name
    form.email = user.value.email
    form.phone = user.value.phone || ''
  }
  
  // Spring entrance animation untuk form
  await nextTick()
  animate(
    '.glass-card',
    { opacity: [0, 1], transform: ['translateY(20px)', 'translateY(0)'] },
    { duration: 0.6, easing: 'spring' }
  )
})

// Validate form
const validateForm = () => {
  let isValid = true
  
  // Reset errors
  errors.full_name = ''
  errors.email = ''
  errors.phone = ''
  
  // Validate full name
  if (!form.full_name || form.full_name.trim().length < 3) {
    errors.full_name = 'Nama lengkap minimal 3 karakter'
    isValid = false
  }
  
  // Validate email
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
  if (!form.email || !emailRegex.test(form.email)) {
    errors.email = 'Format email tidak valid'
    isValid = false
  }
  
  // Validate phone
  const phoneRegex = /^08[0-9]{8,13}$/
  if (!form.phone || !phoneRegex.test(form.phone)) {
    errors.phone = 'Format nomor telepon tidak valid (08xxxxxxxxxx)'
    isValid = false
  }
  
  return isValid
}

// Submit form
const handleSubmit = async () => {
  errorMessage.value = ''
  
  if (!validateForm()) {
    return
  }
  
  loading.value = true
  
  try {
    const response = await apiClient.put('/profile', form)
    
    if (response.data.success) {
      // Update auth store dengan data terbaru
      await authStore.fetchCurrentUser()
      
      // Show success message
      alert('Profile berhasil diupdate')
      
      // Navigate back to profile
      router.push('/profile')
    }
  } catch (error) {
    errorMessage.value = error.response?.data?.error || 'Gagal update profile'
  } finally {
    loading.value = false
  }
}

// Go back
const goBack = () => {
  router.push('/profile')
}
</script>
