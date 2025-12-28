<template>
  <BaseModal
    v-model="isOpen"
    :title="isEditMode ? 'Edit User' : 'Tambah User Baru'"
    :subtitle="isEditMode ? 'Update informasi user' : 'Buat user baru dengan kredensial auto-generated'"
    size="lg"
    :show-footer="!generatedPassword"
    :show-cancel="!loading"
    :show-confirm="true"
    :confirm-text="isEditMode ? 'Simpan Perubahan' : 'Buat User'"
    :loading="loading"
    :loading-text="isEditMode ? 'Menyimpan...' : 'Membuat...'"
    :dismissible="!loading"
    @confirm="handleSubmit"
    @cancel="handleClose"
    @close="handleClose"
  >
    <!-- Form Content -->

    <form @submit.prevent="handleSubmit" class="space-y-5">
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
    </form>

    <!-- Custom Footer untuk Success State -->
    <template v-if="generatedPassword" #footer>
      <!-- Success Message dengan Generated Password -->
      <div class="space-y-4">
        <div class="p-4 bg-emerald-50 border border-emerald-200 rounded-xl space-y-3">
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
              class="p-2 hover:bg-emerald-50 rounded-lg transition-colors active-scale"
              title="Copy password"
            >
              <Copy class="w-4 h-4 text-emerald-600" />
            </button>
          </div>
          <p class="text-xs text-emerald-600">
            ⚠️ Password ini hanya ditampilkan sekali. User harus mengubah password saat login pertama kali.
          </p>
        </div>
        
        <!-- Close Button -->
        <button
          type="button"
          @click="handleClose"
          class="w-full btn-primary"
        >
          Selesai
        </button>
      </div>
    </template>
  </BaseModal>

  <!-- Alert Dialog untuk Copy Success -->
  <AlertDialog
    v-model="alertDialog.isOpen.value"
    :message="alertDialog.config.value.message"
    :variant="alertDialog.config.value.variant"
    @close="alertDialog.handleClose"
  />
</template>

<script setup>
import { ref, reactive, computed, watch } from 'vue'
import { useUserStore } from '../../stores/user'
import { BaseModal, AlertDialog } from '../common'
import { useAlertDialog } from '../../composables/useModal'
import { CheckCircle, Copy } from 'lucide-vue-next'

const props = defineProps({
  user: {
    type: Object,
    default: null
  }
})

const emit = defineEmits(['close', 'success'])

const userStore = useUserStore()
const alertDialog = useAlertDialog()

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

// Copy password dengan AlertDialog
const copyPassword = async () => {
  try {
    await navigator.clipboard.writeText(generatedPassword.value)
    await alertDialog.success('Password berhasil di-copy!', {
      detail: 'Password telah disalin ke clipboard.',
      autoDismiss: true,
      autoDismissDelay: 2000
    })
  } catch (error) {
    await alertDialog.error('Gagal menyalin password', {
      detail: 'Silakan copy secara manual.'
    })
  }
}

// Close modal
const handleClose = () => {
  isOpen.value = false
  setTimeout(() => {
    emit('close')
  }, 300)
}
</script>
