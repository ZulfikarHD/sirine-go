<template>
  <!-- Modal Backdrop -->
  <Teleport to="body">
    <Motion
      v-if="isOpen"
      :initial="{ opacity: 0 }"
      :animate="{ opacity: 1 }"
      :exit="{ opacity: 0 }"
      :transition="{ duration: 0.3 }"
      class="fixed inset-0 bg-black bg-opacity-50 backdrop-blur-sm z-50 flex items-center justify-center p-4"
      @click.self="handleClose"
    >
      <!-- Modal Content -->
      <Motion
        :initial="{ opacity: 0, scale: 0.95, y: 20 }"
        :animate="{ opacity: 1, scale: 1, y: 0 }"
        :exit="{ opacity: 0, scale: 0.95, y: 20 }"
        :transition="{ duration: 0.3, easing: 'spring' }"
        class="bg-white rounded-2xl shadow-2xl max-w-2xl w-full max-h-[90vh] overflow-y-auto"
        @click.stop
      >
            <!-- Header -->
            <div class="sticky top-0 bg-white border-b border-gray-200 px-6 py-4 flex items-center justify-between rounded-t-2xl">
              <div>
                <h2 class="text-xl font-bold text-gray-900">
                  {{ isEditMode ? 'Edit User' : 'Tambah User Baru' }}
                </h2>
                <p class="text-sm text-gray-600 mt-1">
                  {{ isEditMode ? 'Update informasi user' : 'Buat user baru dengan kredensial auto-generated' }}
                </p>
              </div>
              <button
                @click="handleClose"
                class="p-2 hover:bg-gray-100 rounded-lg transition-colors"
              >
                <X class="w-5 h-5 text-gray-500" />
              </button>
            </div>

            <!-- Form -->
            <form @submit.prevent="handleSubmit" class="p-6 space-y-5">
              <!-- NIP (only for create) -->
              <div v-if="!isEditMode">
                <label class="block text-sm font-semibold text-gray-700 mb-2">
                  NIP <span class="text-red-500">*</span>
                </label>
                <input
                  v-model="form.nip"
                  type="text"
                  maxlength="5"
                  required
                  class="w-full px-4 py-2.5 rounded-xl border border-gray-200 focus:border-indigo-500 focus:ring-2 focus:ring-indigo-200 transition-all"
                  placeholder="Maksimal 5 digit"
                />
                <p class="mt-1 text-xs text-gray-500">NIP harus unik dan maksimal 5 karakter</p>
              </div>

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
              </div>

              <!-- Role & Department Grid -->
              <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                <!-- Role -->
                <div>
                  <label class="block text-sm font-semibold text-gray-700 mb-2">
                    Role <span class="text-red-500">*</span>
                  </label>
                  <select
                    v-model="form.role"
                    required
                    class="w-full px-4 py-2.5 rounded-xl border border-gray-200 focus:border-indigo-500 focus:ring-2 focus:ring-indigo-200 transition-all"
                  >
                    <option value="">Pilih Role</option>
                    <option value="ADMIN">Admin</option>
                    <option value="MANAGER">Manager</option>
                    <option value="STAFF_KHAZWAL">Staff Khazwal</option>
                    <option value="OPERATOR_CETAK">Operator Cetak</option>
                    <option value="QC_INSPECTOR">QC Inspector</option>
                    <option value="VERIFIKATOR">Verifikator</option>
                    <option value="STAFF_KHAZKHIR">Staff Khazkhir</option>
                  </select>
                </div>

                <!-- Department -->
                <div>
                  <label class="block text-sm font-semibold text-gray-700 mb-2">
                    Department <span class="text-red-500">*</span>
                  </label>
                  <select
                    v-model="form.department"
                    required
                    class="w-full px-4 py-2.5 rounded-xl border border-gray-200 focus:border-indigo-500 focus:ring-2 focus:ring-indigo-200 transition-all"
                  >
                    <option value="">Pilih Department</option>
                    <option value="KHAZWAL">Khazwal</option>
                    <option value="CETAK">Cetak</option>
                    <option value="VERIFIKASI">Verifikasi</option>
                    <option value="KHAZKHIR">Khazkhir</option>
                  </select>
                </div>
              </div>

              <!-- Shift & Status Grid -->
              <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                <!-- Shift -->
                <div>
                  <label class="block text-sm font-semibold text-gray-700 mb-2">
                    Shift
                  </label>
                  <select
                    v-model="form.shift"
                    class="w-full px-4 py-2.5 rounded-xl border border-gray-200 focus:border-indigo-500 focus:ring-2 focus:ring-indigo-200 transition-all"
                  >
                    <option value="PAGI">Pagi</option>
                    <option value="SIANG">Siang</option>
                    <option value="MALAM">Malam</option>
                  </select>
                </div>

                <!-- Status (only for edit) -->
                <div v-if="isEditMode">
                  <label class="block text-sm font-semibold text-gray-700 mb-2">
                    Status
                  </label>
                  <select
                    v-model="form.status"
                    class="w-full px-4 py-2.5 rounded-xl border border-gray-200 focus:border-indigo-500 focus:ring-2 focus:ring-indigo-200 transition-all"
                  >
                    <option value="ACTIVE">Active</option>
                    <option value="INACTIVE">Inactive</option>
                    <option value="SUSPENDED">Suspended</option>
                  </select>
                </div>
              </div>

              <!-- Error Message -->
              <div v-if="errorMessage" class="p-4 bg-red-50 border border-red-200 rounded-xl">
                <p class="text-sm text-red-700">{{ errorMessage }}</p>
              </div>

              <!-- Success Message (untuk create dengan generated password) -->
              <div v-if="generatedPassword" class="p-4 bg-emerald-50 border border-emerald-200 rounded-xl space-y-3">
                <div class="flex items-start space-x-3">
                  <CheckCircle class="w-5 h-5 text-emerald-600 mt-0.5" />
                  <div class="flex-1">
                    <p class="text-sm font-semibold text-emerald-900">User berhasil dibuat!</p>
                    <p class="text-sm text-emerald-700 mt-1">Password telah di-generate otomatis. Pastikan untuk menyimpan password ini:</p>
                  </div>
                </div>
                <div class="flex items-center space-x-2 bg-white p-3 rounded-lg border border-emerald-200">
                  <code class="flex-1 text-sm font-mono text-gray-900">{{ generatedPassword }}</code>
                  <button
                    type="button"
                    @click="copyPassword"
                    class="p-2 hover:bg-emerald-50 rounded-lg transition-colors"
                    title="Copy password"
                  >
                    <Copy class="w-4 h-4 text-emerald-600" />
                  </button>
                </div>
                <p class="text-xs text-emerald-600">
                  ⚠️ Password ini hanya ditampilkan sekali. User harus mengubah password saat login pertama kali.
                </p>
              </div>

              <!-- Actions -->
              <div class="flex items-center space-x-3 pt-4">
                <button
                  v-if="!generatedPassword"
                  type="submit"
                  :disabled="loading"
                  class="flex-1 btn-primary disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  <span v-if="loading" class="flex items-center justify-center">
                    <div class="w-5 h-5 border-2 border-white border-t-transparent rounded-full animate-spin mr-2"></div>
                    {{ isEditMode ? 'Menyimpan...' : 'Membuat...' }}
                  </span>
                  <span v-else>
                    {{ isEditMode ? 'Simpan Perubahan' : 'Buat User' }}
                  </span>
                </button>
                <button
                  v-if="generatedPassword"
                  type="button"
                  @click="handleClose"
                  class="flex-1 btn-primary"
                >
                  Selesai
                </button>
                <button
                  v-if="!generatedPassword"
                  type="button"
                  @click="handleClose"
                  :disabled="loading"
                  class="px-6 py-2.5 rounded-xl border border-gray-200 text-gray-700 font-semibold hover:bg-gray-50 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  Batal
                </button>
              </div>
            </form>
          </Motion>
        </Motion>
  </Teleport>
</template>

<script setup>
import { ref, reactive, computed, watch } from 'vue'
import { useUserStore } from '../../stores/user'
import { Motion } from 'motion-v'
import { X, CheckCircle, Copy } from 'lucide-vue-next'

const props = defineProps({
  user: {
    type: Object,
    default: null
  }
})

const emit = defineEmits(['close', 'success'])

const userStore = useUserStore()
const isOpen = ref(true)
const loading = ref(false)
const errorMessage = ref('')
const generatedPassword = ref('')

const isEditMode = computed(() => !!props.user)

const form = reactive({
  nip: '',
  full_name: '',
  email: '',
  phone: '',
  role: '',
  department: '',
  shift: 'PAGI',
  status: 'ACTIVE'
})

// Load user data untuk edit mode
watch(() => props.user, (newUser) => {
  if (newUser) {
    form.full_name = newUser.full_name
    form.email = newUser.email
    form.phone = newUser.phone
    form.role = newUser.role
    form.department = newUser.department
    form.shift = newUser.shift
    form.status = newUser.status
  }
}, { immediate: true })

// Submit form
const handleSubmit = async () => {
  errorMessage.value = ''
  loading.value = true

  try {
    if (isEditMode.value) {
      // Update user
      await userStore.updateUser(props.user.id, {
        full_name: form.full_name,
        email: form.email,
        phone: form.phone,
        role: form.role,
        department: form.department,
        shift: form.shift,
        status: form.status
      })
      emit('success')
      handleClose()
    } else {
      // Create user
      const response = await userStore.createUser({
        nip: form.nip,
        full_name: form.full_name,
        email: form.email,
        phone: form.phone,
        role: form.role,
        department: form.department,
        shift: form.shift
      })
      
      // Show generated password
      generatedPassword.value = response.generated_password
      emit('success')
    }
  } catch (error) {
    errorMessage.value = error.response?.data?.error || userStore.error || 'Terjadi kesalahan'
  } finally {
    loading.value = false
  }
}

// Copy password
const copyPassword = () => {
  navigator.clipboard.writeText(generatedPassword.value)
  alert('Password berhasil di-copy!')
}

// Close modal
const handleClose = () => {
  isOpen.value = false
  setTimeout(() => {
    emit('close')
  }, 300)
}
</script>
