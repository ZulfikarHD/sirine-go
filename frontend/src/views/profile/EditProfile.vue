<template>
  <AppLayout>
    <div class="max-w-2xl mx-auto">
      <!-- Header -->
      <Motion
        v-bind="entranceAnimations.fadeUp"
        class="mb-6"
      >
        <button
          @click="goBack"
          class="flex items-center space-x-2 text-gray-600 hover:text-gray-900 mb-4"
        >
          <ArrowLeft class="w-5 h-5" />
          <span>Kembali ke Profile</span>
        </button>
        <h1 class="text-2xl font-bold text-gray-900">Edit Profile</h1>
        <p class="text-sm text-gray-600 mt-1">Update informasi profile Anda</p>
      </Motion>

      <!-- Form -->
      <Motion
        :initial="{ opacity: 0, y: 15 }"
        :animate="{ opacity: 1, y: 0 }"
        :transition="{ duration: 0.25, delay: 0.1, ease: 'easeOut' }"
      >
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
              class="input-field"
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
              class="input-field"
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
              class="input-field"
              placeholder="08xxxxxxxxxx"
            />
            <p v-if="errors.phone" class="mt-1 text-sm text-red-600">{{ errors.phone }}</p>
          </div>

          <!-- Read-only fields -->
          <div class="pt-4 border-t border-gray-200">
            <p class="text-sm font-semibold text-gray-700 mb-3">Informasi yang tidak dapat diubah:</p>
            
            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <label class="block text-sm text-gray-600 mb-1">NIP</label>
                <div class="px-4 py-2.5 rounded-xl bg-gray-100 text-gray-500 border border-gray-200">
                  {{ user?.nip }}
                </div>
              </div>

              <div>
                <label class="block text-sm text-gray-600 mb-1">Role</label>
                <div class="px-4 py-2.5 rounded-xl bg-gray-100 text-gray-500 border border-gray-200">
                  {{ user?.role }}
                </div>
              </div>

              <div>
                <label class="block text-sm text-gray-600 mb-1">Department</label>
                <div class="px-4 py-2.5 rounded-xl bg-gray-100 text-gray-500 border border-gray-200">
                  {{ user?.department }}
                </div>
              </div>

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
          <Motion
            v-if="errorMessage"
            :initial="{ opacity: 0, x: -10 }"
            :animate="{ opacity: 1, x: 0 }"
            :transition="{ ...springPresets.snappy }"
            class="p-4 bg-red-50 border border-red-200 rounded-xl"
          >
            <p class="text-sm text-red-700">{{ errorMessage }}</p>
          </Motion>

          <!-- Actions -->
          <div class="flex items-center space-x-3 pt-4">
            <button
              type="submit"
              :disabled="loading"
              class="flex-1 btn-primary disabled:opacity-50 disabled:cursor-not-allowed"
            >
              <span v-if="loading" class="flex items-center justify-center">
                <div class="modal-spinner mr-2"></div>
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
              class="px-6 py-2.5 rounded-xl border border-gray-200 text-gray-700 font-semibold hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              Batal
            </button>
          </div>
        </form>
      </Motion>
    </div>

    <!-- Alert Dialog untuk Success/Error -->
    <AlertDialog
      v-model="alertDialog.isOpen.value"
      :title="alertDialog.config.value.title"
      :message="alertDialog.config.value.message"
      :variant="alertDialog.config.value.variant"
      :detail="alertDialog.config.value.detail"
      @close="alertDialog.handleClose"
    />
  </AppLayout>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { Motion } from 'motion-v'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../../stores/auth'
import apiClient from '../../composables/useApi'
import AppLayout from '../../components/layout/AppLayout.vue'
import { AlertDialog } from '../../components/common'
import { useAlertDialog } from '../../composables/useModal'
import { entranceAnimations, springPresets } from '../../composables/useMotion'
import { ArrowLeft, Save } from 'lucide-vue-next'

const router = useRouter()
const authStore = useAuthStore()
const alertDialog = useAlertDialog()

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
onMounted(() => {
  if (user.value) {
    form.full_name = user.value.full_name
    form.email = user.value.email
    form.phone = user.value.phone || ''
  }
})

// Validate form
const validateForm = () => {
  let isValid = true
  
  errors.full_name = ''
  errors.email = ''
  errors.phone = ''
  
  if (!form.full_name || form.full_name.trim().length < 3) {
    errors.full_name = 'Nama lengkap minimal 3 karakter'
    isValid = false
  }
  
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
  if (!form.email || !emailRegex.test(form.email)) {
    errors.email = 'Format email tidak valid'
    isValid = false
  }
  
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
      await authStore.fetchCurrentUser()
      
      await alertDialog.success('Profile berhasil diupdate!', {
        detail: 'Perubahan profile Anda telah disimpan.',
        autoDismiss: true,
        autoDismissDelay: 2000
      })
      
      setTimeout(() => {
        router.push('/profile')
      }, 2000)
    }
  } catch (error) {
    errorMessage.value = error.response?.data?.error || 'Gagal update profile'
    
    await alertDialog.error('Gagal update profile', {
      detail: error.response?.data?.error || 'Terjadi kesalahan saat menyimpan perubahan.'
    })
  } finally {
    loading.value = false
  }
}

const goBack = () => {
  router.push('/profile')
}
</script>
