# Modal System - Quick Start Guide

## Installation

Modal components sudah tersedia di `/src/components/common/`. Tidak perlu instalasi tambahan.

## Quick Usage

### 1. Import Components

```javascript
// Option 1: Individual imports
import BaseModal from '@/components/common/BaseModal.vue'
import ConfirmDialog from '@/components/common/ConfirmDialog.vue'
import AlertDialog from '@/components/common/AlertDialog.vue'
import PasswordConfirmDialog from '@/components/common/PasswordConfirmDialog.vue'

// Option 2: Bulk import dari index
import { BaseModal, ConfirmDialog, AlertDialog, PasswordConfirmDialog } from '@/components/common'

// Import composables
import { useModal, useConfirmDialog, useAlertDialog, usePasswordConfirmDialog } from '@/composables/useModal'
```

### 2. Basic Modal (CRUD)

```vue
<template>
  <div>
    <button @click="modal.open()" class="btn-primary">
      Open Form
    </button>

    <BaseModal
      v-model="modal.isOpen.value"
      title="Tambah Data"
      :show-footer="true"
      :loading="modal.loading.value"
      @confirm="handleSubmit"
    >
      <!-- Your form here -->
      <form @submit.prevent="handleSubmit">
        <input v-model="name" class="input-field" />
      </form>
    </BaseModal>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { BaseModal } from '@/components/common'
import { useModal } from '@/composables/useModal'

const modal = useModal()
const name = ref('')

const handleSubmit = async () => {
  modal.setLoading(true)
  await saveData()
  modal.setLoading(false)
  modal.close()
}
</script>
```

### 3. Confirmation Dialog

```vue
<template>
  <div>
    <button @click="handleDelete" class="btn-primary">
      Delete
    </button>

    <ConfirmDialog
      v-model="confirm.isOpen.value"
      :title="confirm.config.value.title"
      :message="confirm.config.value.message"
      :variant="confirm.config.value.variant"
      :loading="confirm.loading.value"
      @confirm="confirm.handleConfirm()"
      @cancel="confirm.handleCancel()"
    />
  </div>
</template>

<script setup>
import { ConfirmDialog } from '@/components/common'
import { useConfirmDialog } from '@/composables/useModal'

const confirm = useConfirmDialog()

const handleDelete = async () => {
  const result = await confirm.confirm({
    title: 'Hapus Data',
    message: 'Yakin ingin menghapus?',
    variant: 'danger'
  })

  if (result) {
    confirm.loading.value = true
    await deleteData()
    confirm.close()
  }
}
</script>
```

### 4. Password Confirmation

```vue
<template>
  <div>
    <button @click="handleSecureAction" class="btn-primary">
      Delete Account
    </button>

    <PasswordConfirmDialog
      ref="passwordDialogRef"
      v-model="passwordConfirm.isOpen.value"
      :title="passwordConfirm.config.value.title"
      :message="passwordConfirm.config.value.message"
      :loading="passwordConfirm.loading.value"
      @confirm="passwordConfirm.handleConfirm"
      @cancel="passwordConfirm.handleCancel"
    />
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { PasswordConfirmDialog } from '@/components/common'
import { usePasswordConfirmDialog } from '@/composables/useModal'

const passwordConfirm = usePasswordConfirmDialog()
const passwordDialogRef = ref(null)
passwordConfirm.setDialogRef(passwordDialogRef)

const handleSecureAction = async () => {
  const password = await passwordConfirm.confirm({
    title: 'Konfirmasi Password',
    message: 'Masukkan password untuk menghapus account.'
  })

  if (password) {
    passwordConfirm.setLoading(true)
    try {
      // Verify password with backend
      const verified = await verifyPassword(password)
      if (verified) {
        passwordConfirm.close()
        await deleteAccount()
      } else {
        passwordConfirm.setError('Password salah')
      }
    } catch (error) {
      passwordConfirm.setError(error.message)
    }
  }
}
</script>
```

### 5. Alert Notifications

```vue
<template>
  <AlertDialog
    v-model="alert.isOpen.value"
    :title="alert.config.value.title"
    :message="alert.config.value.message"
    :variant="alert.config.value.variant"
    @close="alert.handleClose()"
  />
</template>

<script setup>
import { AlertDialog } from '@/components/common'
import { useAlertDialog } from '@/composables/useModal'

const alert = useAlertDialog()

// Show alerts
alert.success('Berhasil disimpan!')
alert.error('Terjadi kesalahan!')
alert.warning('Perhatian!')
alert.info('Informasi penting')
</script>
```

## Common Patterns

### Pattern 1: Form dengan Validation

```javascript
const handleSubmit = async () => {
  // Validate
  if (!form.name) {
    await alert.warning('Nama harus diisi!')
    return
  }

  // Submit
  modal.setLoading(true)
  try {
    await saveData(form)
    modal.close()
    await alert.success('Data berhasil disimpan!')
  } catch (error) {
    await alert.error('Gagal menyimpan data', {
      detail: error.message
    })
  } finally {
    modal.setLoading(false)
  }
}
```

### Pattern 2: Delete Confirmation

```javascript
const handleDelete = async (item) => {
  const confirmed = await confirm.confirm({
    title: 'Hapus Data',
    message: `Yakin ingin menghapus "${item.name}"?`,
    variant: 'danger',
    confirmText: 'Ya, Hapus',
    showWarning: true
  })

  if (confirmed) {
    confirm.loading.value = true
    try {
      await deleteItem(item.id)
      confirm.close()
      await alert.success('Data berhasil dihapus!')
    } catch (error) {
      await alert.error('Gagal menghapus data')
    }
  }
}
```

### Pattern 3: Auto-dismiss Toast

```javascript
// Quick success toast
const showToast = () => {
  alert.success('Berhasil!', {
    autoDismiss: true,
    autoDismissDelay: 2000,
    showTimer: true
  })
}
```

## Props Cheatsheet

### BaseModal

```javascript
// Size options
size="xs" | "sm" | "md" | "lg" | "xl" | "full"

// Common props
:show-footer="true"
:show-cancel="true"
:show-confirm="true"
:loading="loading"
confirm-text="Simpan"
cancel-text="Batal"

// Behavior
:dismissible="true"
:close-on-backdrop="true"
:close-on-escape="true"

// Styling
title-gradient  // Gradient effect on title
:no-padding="true"  // Remove content padding
```

### ConfirmDialog

```javascript
// Variants
variant="default" | "danger" | "warning" | "info" | "success"

// Danger variant props
:show-warning="true"
warning-message="Custom warning text"
```

### PasswordConfirmDialog

```javascript
// Common props
title="Konfirmasi Password"
message="Masukkan password Anda"
confirm-text="Verifikasi"

// Security notice
:show-security-notice="true"
security-notice="Custom security message"

// Validation
:validate-password="validateFn"
:auto-clear-password="true"
```

### AlertDialog

```javascript
// Variants
variant="success" | "error" | "warning" | "info"

// Auto-dismiss
:auto-dismiss="true"
:auto-dismiss-delay="3000"
:show-timer="true"
```

## Composable Methods

```javascript
// useModal()
modal.open()
modal.close()
modal.toggle()
modal.setLoading(true)

// useConfirmDialog()
const result = await confirm.confirm({ options })
confirm.handleConfirm()
confirm.handleCancel()
confirm.close()

// useAlertDialog()
await alert.alert({ options })
await alert.success(message, options)
await alert.error(message, options)
await alert.warning(message, options)
await alert.info(message, options)
alert.close()

// usePasswordConfirmDialog()
const password = await passwordConfirm.confirm({ options })
passwordConfirm.setError(message)
passwordConfirm.setLoading(true)
passwordConfirm.close()
passwordConfirm.setDialogRef(ref)
```

## Examples

Lihat `/src/components/common/ModalExamples.vue` untuk comprehensive examples dan interactive demos.

Run development server dan navigate ke route yang menggunakan ModalExamples component.

## Need Help?

- **Full Documentation**: `docs/components/modal-system.md`
- **Examples**: `src/components/common/ModalExamples.vue`
- **Developer**: Zulfikar Hidayatullah (+62 857-1583-8733)
