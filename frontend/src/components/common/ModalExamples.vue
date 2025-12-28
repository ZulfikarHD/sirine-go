<template>
  <div class="max-w-7xl mx-auto p-6 space-y-8">
    <!-- Header -->
    <div class="text-center space-y-3 mb-12">
      <h1 class="text-4xl font-bold text-gradient-indigo-fuchsia">
        Modal Components Examples
      </h1>
      <p class="text-lg text-gray-600">
        iOS-inspired modal system dengan spring animations dan glass morphism
      </p>
    </div>

    <!-- BaseModal Examples -->
    <section class="space-y-6">
      <h2 class="text-2xl font-bold text-gray-900">BaseModal - CRUD Operations</h2>
      
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        <!-- Simple Modal -->
        <button @click="simpleModal.open()" class="btn-primary">
          Simple Modal
        </button>

        <!-- Form Modal -->
        <button @click="formModal.open()" class="btn-primary">
          Form Modal (MD)
        </button>

        <!-- Large Modal -->
        <button @click="largeModal.open()" class="btn-primary">
          Large Modal (LG)
        </button>

        <!-- Full Modal -->
        <button @click="fullModal.open()" class="btn-primary">
          Full Width Modal
        </button>

        <!-- Non-dismissible -->
        <button @click="nonDismissibleModal.open()" class="btn-primary">
          Non-dismissible Modal
        </button>

        <!-- With Loading -->
        <button @click="handleLoadingModal" class="btn-primary">
          Modal with Loading
        </button>
      </div>
    </section>

    <!-- ConfirmDialog Examples -->
    <section class="space-y-6">
      <h2 class="text-2xl font-bold text-gray-900">ConfirmDialog - Confirmations</h2>
      
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        <!-- Default Confirmation -->
        <button @click="showDefaultConfirm" class="btn-primary">
          Default Confirm
        </button>

        <!-- Danger Confirmation -->
        <button @click="showDangerConfirm" class="btn-primary bg-linear-to-r! from-red-600! to-pink-600!">
          Danger Confirm
        </button>

        <!-- Warning Confirmation -->
        <button @click="showWarningConfirm" class="btn-primary bg-linear-to-r! from-amber-600! to-orange-600!">
          Warning Confirm
        </button>

        <!-- Info Confirmation -->
        <button @click="showInfoConfirm" class="btn-primary bg-linear-to-r! from-blue-600! to-cyan-600!">
          Info Confirm
        </button>
      </div>

      <!-- Result Display -->
      <div v-if="confirmResult" class="p-4 bg-gray-100 rounded-xl">
        <p class="text-sm text-gray-700">
          <strong>Result:</strong> {{ confirmResult }}
        </p>
      </div>
    </section>

    <!-- PasswordConfirmDialog Examples -->
    <section class="space-y-6">
      <h2 class="text-2xl font-bold text-gray-900">PasswordConfirmDialog - Security</h2>
      
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        <!-- Basic Password Confirm -->
        <button @click="showPasswordConfirm" class="btn-primary">
          Password Confirm
        </button>

        <!-- Delete Account with Password -->
        <button @click="showDeleteAccountConfirm" class="btn-primary bg-linear-to-r! from-red-600! to-pink-600!">
          Delete Account (Secure)
        </button>

        <!-- Change Settings with Password -->
        <button @click="showChangeSettingsConfirm" class="btn-primary bg-linear-to-r! from-amber-600! to-orange-600!">
          Change Settings
        </button>
      </div>

      <!-- Result Display -->
      <div v-if="passwordResult" class="p-4 bg-gray-100 rounded-xl">
        <p class="text-sm text-gray-700">
          <strong>Result:</strong> {{ passwordResult }}
        </p>
      </div>
    </section>

    <!-- AlertDialog Examples -->
    <section class="space-y-6">
      <h2 class="text-2xl font-bold text-gray-900">AlertDialog - Notifications</h2>
      
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        <!-- Success Alert -->
        <button @click="showSuccessAlert" class="btn-primary bg-linear-to-r! from-emerald-600! to-green-600!">
          Success Alert
        </button>

        <!-- Error Alert -->
        <button @click="showErrorAlert" class="btn-primary bg-linear-to-r! from-red-600! to-pink-600!">
          Error Alert
        </button>

        <!-- Warning Alert -->
        <button @click="showWarningAlert" class="btn-primary bg-linear-to-r! from-amber-600! to-orange-600!">
          Warning Alert
        </button>

        <!-- Info Alert -->
        <button @click="showInfoAlert" class="btn-primary bg-linear-to-r! from-blue-600! to-cyan-600!">
          Info Alert
        </button>

        <!-- Auto Dismiss -->
        <button @click="showAutoDismissAlert" class="btn-secondary">
          Auto Dismiss (3s)
        </button>
      </div>
    </section>

    <!-- Usage Guide -->
    <section class="space-y-6 mt-12 p-6 bg-gray-50 rounded-2xl">
      <h2 class="text-2xl font-bold text-gray-900">Usage Guide</h2>
      
      <div class="space-y-4 text-sm">
        <div class="p-4 bg-white rounded-xl">
          <h3 class="font-semibold text-gray-900 mb-2">1. Import Components</h3>
          <pre class="bg-gray-900 text-gray-100 p-4 rounded-lg overflow-x-auto"><code>import BaseModal from '@/components/common/BaseModal.vue'
import ConfirmDialog from '@/components/common/ConfirmDialog.vue'
import AlertDialog from '@/components/common/AlertDialog.vue'
import PasswordConfirmDialog from '@/components/common/PasswordConfirmDialog.vue'
import { useModal, useConfirmDialog, useAlertDialog, usePasswordConfirmDialog } from '@/composables/useModal'</code></pre>
        </div>

        <div class="p-4 bg-white rounded-xl">
          <h3 class="font-semibold text-gray-900 mb-2">2. Use with Composable</h3>
          <pre class="bg-gray-900 text-gray-100 p-4 rounded-lg overflow-x-auto"><code>const modal = useModal()
const confirmDialog = useConfirmDialog()
const alertDialog = useAlertDialog()
const passwordConfirm = usePasswordConfirmDialog()

// Open modal
modal.open()

// Confirm with promise
const result = await confirmDialog.confirm({
  title: 'Hapus Data',
  message: 'Yakin ingin menghapus data ini?',
  variant: 'danger'
})

// Password verification
const password = await passwordConfirm.confirm({
  title: 'Verifikasi',
  message: 'Masukkan password Anda'
})

// Show alert
await alertDialog.success('Data berhasil disimpan!')</code></pre>
        </div>
      </div>
    </section>

    <!-- Modals -->
    <BaseModal
      v-model="simpleModal.isOpen.value"
      title="Simple Modal"
      subtitle="This is a simple modal example"
      size="sm"
    >
      <p class="text-gray-700">
        This is a simple modal with minimal configuration.
      </p>
    </BaseModal>

    <BaseModal
      v-model="formModal.isOpen.value"
      title="Form Modal"
      subtitle="Tambah data baru"
      title-gradient
      :show-footer="true"
      :loading="formModal.loading.value"
      @confirm="handleFormSubmit"
    >
      <form @submit.prevent="handleFormSubmit" class="space-y-4">
        <div>
          <label class="block text-sm font-semibold text-gray-700 mb-2">
            Nama Lengkap
          </label>
          <input
            v-model="formData.name"
            type="text"
            class="input-field"
            placeholder="Masukkan nama"
          />
        </div>
        <div>
          <label class="block text-sm font-semibold text-gray-700 mb-2">
            Email
          </label>
          <input
            v-model="formData.email"
            type="email"
            class="input-field"
            placeholder="email@example.com"
          />
        </div>
      </form>
    </BaseModal>

    <BaseModal
      v-model="largeModal.isOpen.value"
      title="Large Modal"
      size="lg"
      :show-footer="true"
    >
      <div class="space-y-4">
        <p class="text-gray-700">
          This is a large modal that can contain more content.
        </p>
        <div v-for="i in 10" :key="i" class="p-4 bg-gray-50 rounded-xl">
          <p class="text-sm text-gray-600">Content item {{ i }}</p>
        </div>
      </div>
    </BaseModal>

    <BaseModal
      v-model="fullModal.isOpen.value"
      title="Full Width Modal"
      size="full"
      content-max-height="70vh"
    >
      <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
        <div v-for="i in 12" :key="i" class="p-6 bg-gray-50 rounded-xl">
          <p class="text-sm text-gray-600">Grid item {{ i }}</p>
        </div>
      </div>
    </BaseModal>

    <BaseModal
      v-model="nonDismissibleModal.isOpen.value"
      title="Non-dismissible Modal"
      subtitle="You must click the button to close"
      :dismissible="false"
      :show-footer="true"
      :show-cancel="false"
      @confirm="nonDismissibleModal.close()"
    >
      <p class="text-gray-700">
        This modal cannot be dismissed by clicking outside or pressing Escape.
      </p>
    </BaseModal>

    <BaseModal
      v-model="loadingModal.isOpen.value"
      title="Processing..."
      :loading="loadingModal.loading.value"
      :dismissible="!loadingModal.loading.value"
      :show-footer="true"
      @confirm="handleStopLoading"
    >
      <p class="text-gray-700">
        Modal with loading state demonstration.
      </p>
    </BaseModal>

    <ConfirmDialog
      v-model="confirmDialog.isOpen.value"
      :title="confirmDialog.config.value.title"
      :message="confirmDialog.config.value.message"
      :detail="confirmDialog.config.value.detail"
      :variant="confirmDialog.config.value.variant"
      :confirm-text="confirmDialog.config.value.confirmText"
      :cancel-text="confirmDialog.config.value.cancelText"
      :show-warning="confirmDialog.config.value.showWarning"
      :warning-message="confirmDialog.config.value.warningMessage"
      :loading="confirmDialog.loading.value"
      @confirm="confirmDialog.handleConfirm()"
      @cancel="confirmDialog.handleCancel()"
    />

    <AlertDialog
      v-model="alertDialog.isOpen.value"
      :title="alertDialog.config.value.title"
      :message="alertDialog.config.value.message"
      :detail="alertDialog.config.value.detail"
      :variant="alertDialog.config.value.variant"
      :confirm-text="alertDialog.config.value.confirmText"
      :auto-dismiss="alertDialog.config.value.autoDismiss"
      :auto-dismiss-delay="alertDialog.config.value.autoDismissDelay"
      :show-timer="alertDialog.config.value.showTimer"
      @close="alertDialog.handleClose()"
    />

    <PasswordConfirmDialog
      ref="passwordDialogRef"
      v-model="passwordConfirmDialog.isOpen.value"
      :title="passwordConfirmDialog.config.value.title"
      :subtitle="passwordConfirmDialog.config.value.subtitle"
      :message="passwordConfirmDialog.config.value.message"
      :detail="passwordConfirmDialog.config.value.detail"
      :confirm-text="passwordConfirmDialog.config.value.confirmText"
      :show-security-notice="passwordConfirmDialog.config.value.showSecurityNotice"
      :security-notice="passwordConfirmDialog.config.value.securityNotice"
      :loading="passwordConfirmDialog.loading.value"
      @confirm="passwordConfirmDialog.handleConfirm"
      @cancel="passwordConfirmDialog.handleCancel"
    />
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import BaseModal from './BaseModal.vue'
import ConfirmDialog from './ConfirmDialog.vue'
import AlertDialog from './AlertDialog.vue'
import PasswordConfirmDialog from './PasswordConfirmDialog.vue'
import { useModal, useConfirmDialog, useAlertDialog, usePasswordConfirmDialog } from '../../composables/useModal'

/**
 * ModalExamples - Comprehensive examples untuk modal components
 * yang mendemonstrasikan berbagai use cases dan configurations
 * dengan interactive demos dan code snippets
 */

// Modal instances
const simpleModal = useModal()
const formModal = useModal()
const largeModal = useModal()
const fullModal = useModal()
const nonDismissibleModal = useModal()
const loadingModal = useModal()

// Confirm & Alert
const confirmDialog = useConfirmDialog()
const alertDialog = useAlertDialog()
const passwordConfirmDialog = usePasswordConfirmDialog()

// State
const confirmResult = ref('')
const passwordResult = ref('')
const passwordDialogRef = ref(null)

// Set dialog ref for password confirm
passwordConfirmDialog.setDialogRef(passwordDialogRef)
const formData = reactive({
  name: '',
  email: ''
})

// Form submit handler
const handleFormSubmit = async () => {
  formModal.setLoading(true)
  
  // Simulate API call
  await new Promise(resolve => setTimeout(resolve, 2000))
  
  formModal.setLoading(false)
  formModal.close()
  
  // Show success alert
  alertDialog.success('Data berhasil disimpan!', {
    detail: 'Form data telah berhasil diproses dan disimpan ke database.'
  })
}

// Loading modal handler
const handleLoadingModal = () => {
  loadingModal.open()
  loadingModal.setLoading(true)
  
  // Auto stop after 3s
  setTimeout(() => {
    loadingModal.setLoading(false)
  }, 3000)
}

const handleStopLoading = () => {
  loadingModal.setLoading(false)
  loadingModal.close()
}

// Confirm handlers
const showDefaultConfirm = async () => {
  const result = await confirmDialog.confirm({
    title: 'Konfirmasi Aksi',
    message: 'Apakah Anda yakin ingin melanjutkan aksi ini?',
    detail: 'Aksi ini akan memproses data Anda.',
    variant: 'default'
  })
  confirmResult.value = result ? 'User confirmed' : 'User cancelled'
}

const showDangerConfirm = async () => {
  const result = await confirmDialog.confirm({
    title: 'Hapus Data',
    message: 'Apakah Anda yakin ingin menghapus data ini?',
    detail: 'Data yang dihapus tidak dapat dikembalikan.',
    variant: 'danger',
    confirmText: 'Ya, Hapus',
    showWarning: true,
    warningMessage: 'Tindakan ini bersifat permanen dan tidak dapat dibatalkan.'
  })
  
  if (result) {
    confirmDialog.loading.value = true
    await new Promise(resolve => setTimeout(resolve, 2000))
    confirmDialog.close()
    alertDialog.success('Data berhasil dihapus!')
  }
  
  confirmResult.value = result ? 'Data deleted' : 'Deletion cancelled'
}

const showWarningConfirm = async () => {
  const result = await confirmDialog.confirm({
    title: 'Peringatan',
    message: 'Tindakan ini akan mengubah data penting.',
    detail: 'Pastikan Anda telah membackup data sebelum melanjutkan.',
    variant: 'warning',
    confirmText: 'Lanjutkan',
    showWarning: false
  })
  confirmResult.value = result ? 'User proceeded' : 'User cancelled'
}

const showInfoConfirm = async () => {
  const result = await confirmDialog.confirm({
    title: 'Informasi',
    message: 'Ada update baru yang tersedia.',
    detail: 'Apakah Anda ingin mengupdate sekarang?',
    variant: 'info',
    confirmText: 'Update Sekarang',
    cancelText: 'Nanti Saja',
    showWarning: false
  })
  confirmResult.value = result ? 'Update started' : 'Update postponed'
}

// Alert handlers
const showSuccessAlert = () => {
  alertDialog.success('Operasi Berhasil!', {
    detail: 'Data telah berhasil disimpan ke database.',
    title: 'Sukses'
  })
}

const showErrorAlert = () => {
  alertDialog.error('Terjadi Kesalahan!', {
    detail: 'Gagal menyimpan data. Silakan coba lagi.',
    title: 'Error'
  })
}

const showWarningAlert = () => {
  alertDialog.warning('Perhatian!', {
    detail: 'Beberapa field belum diisi dengan lengkap.',
    title: 'Warning'
  })
}

const showInfoAlert = () => {
  alertDialog.info('Informasi', {
    detail: 'Sistem akan menjalani maintenance pada hari Minggu.',
    title: 'Info'
  })
}

const showAutoDismissAlert = () => {
  alertDialog.success('Berhasil!', {
    detail: 'Dialog ini akan menutup otomatis dalam 3 detik.',
    autoDismiss: true,
    autoDismissDelay: 3000,
    showTimer: true
  })
}

// Password Confirm handlers
const showPasswordConfirm = async () => {
  const password = await passwordConfirmDialog.confirm({
    title: 'Konfirmasi Password',
    message: 'Masukkan password Anda untuk melanjutkan.',
    detail: 'Tindakan ini memerlukan verifikasi identitas.'
  })

  if (password) {
    passwordConfirmDialog.setLoading(true)
    
    // Simulate API verification
    setTimeout(async () => {
      // Simulate password check (in real app, verify with backend)
      if (password === 'demo') {
        passwordConfirmDialog.close()
        passwordResult.value = 'Password verified successfully!'
        await alertDialog.success('Password Benar!', {
          detail: 'Identitas Anda telah diverifikasi.'
        })
      } else {
        passwordConfirmDialog.setError('Password yang Anda masukkan salah. Silakan coba lagi.')
      }
    }, 1500)
  } else {
    passwordResult.value = 'User cancelled verification'
  }
}

const showDeleteAccountConfirm = async () => {
  const password = await passwordConfirmDialog.confirm({
    title: 'Hapus Account',
    message: 'Untuk menghapus account, masukkan password Anda.',
    detail: 'Tindakan ini tidak dapat dibatalkan.',
    confirmText: 'Hapus Account',
    showSecurityNotice: true,
    securityNotice: '⚠️ Semua data Anda akan dihapus secara permanen.'
  })

  if (password) {
    passwordConfirmDialog.setLoading(true)
    
    // Simulate verification
    setTimeout(async () => {
      if (password === 'demo') {
        passwordConfirmDialog.close()
        passwordResult.value = 'Account deletion authorized'
        await alertDialog.success('Account Terhapus', {
          detail: 'Account Anda telah berhasil dihapus.',
          autoDismiss: true
        })
      } else {
        passwordConfirmDialog.setError('Password salah. Account tidak dapat dihapus.')
      }
    }, 1500)
  }
}

const showChangeSettingsConfirm = async () => {
  const password = await passwordConfirmDialog.confirm({
    title: 'Ubah Pengaturan Keamanan',
    message: 'Verifikasi identitas Anda untuk mengubah pengaturan.',
    confirmText: 'Verifikasi & Ubah'
  })

  if (password) {
    passwordConfirmDialog.setLoading(true)
    
    setTimeout(async () => {
      if (password === 'demo') {
        passwordConfirmDialog.close()
        passwordResult.value = 'Settings change authorized'
        await alertDialog.success('Pengaturan Diubah', {
          detail: 'Pengaturan keamanan telah diperbarui.'
        })
      } else {
        passwordConfirmDialog.setError('Password tidak valid.')
      }
    }, 1500)
  }
}
</script>
