# Modal System Documentation

## Overview

Modal System merupakan sekumpulan komponen dialog yang reusable dengan iOS-inspired design, yaitu: BaseModal untuk CRUD operations, ConfirmDialog untuk confirmations, dan AlertDialog untuk notifications. Sistem ini dibangun dengan Vue 3, motion-vue animations, dan Tailwind CSS yang mengikuti design principles dari Apple ecosystem.

## Features

Sistem modal ini mencakup fitur-fitur, antara lain:

- **iOS Design Principles** - Spring physics animations, press feedback, glass morphism effect
- **Mobile-First Approach** - Responsive design dengan bottom sheet untuk mobile devices
- **Haptic Feedback** - Vibration API untuk tactile response pada mobile
- **Accessibility** - Keyboard navigation, focus trap, ARIA attributes
- **Flexible Sizing** - Multiple size options (xs, sm, md, lg, xl, full)
- **Auto-dismiss** - Optional auto-close untuk alert notifications
- **Programmatic API** - Composable-based API untuk easy usage
- **Stacked Modals** - Support untuk multiple modal instances
- **Animation** - Spring-based animations dengan motion-vue

## Components

### 1. BaseModal

BaseModal merupakan core modal component yang bertujuan untuk CRUD operations, form submissions, dan general-purpose dialogs dengan customizable header, content, dan footer sections.

#### Props

| Prop | Type | Default | Description |
|------|------|---------|-------------|
| `modelValue` | Boolean | `false` | Two-way binding untuk modal visibility |
| `title` | String | `''` | Modal title text |
| `subtitle` | String | `''` | Subtitle text di bawah title |
| `titleGradient` | Boolean | `false` | Apply gradient effect pada title |
| `showHeader` | Boolean | `true` | Tampilkan header section |
| `size` | String | `'md'` | Modal size: xs, sm, md, lg, xl, full |
| `scrollable` | Boolean | `true` | Enable scrolling untuk content |
| `noPadding` | Boolean | `false` | Remove padding dari content area |
| `contentMaxHeight` | String | `'60vh'` | Max height untuk scrollable content |
| `showFooter` | Boolean | `false` | Tampilkan footer section |
| `showCancel` | Boolean | `true` | Tampilkan cancel button |
| `showConfirm` | Boolean | `true` | Tampilkan confirm button |
| `cancelText` | String | `'Batal'` | Cancel button text |
| `confirmText` | String | `'Simpan'` | Confirm button text |
| `loadingText` | String | `'Memproses...'` | Loading state text |
| `confirmDanger` | Boolean | `false` | Style confirm button sebagai danger (red) |
| `confirmDisabled` | Boolean | `false` | Disable confirm button |
| `dismissible` | Boolean | `true` | Allow modal to be dismissed |
| `closeOnBackdrop` | Boolean | `true` | Close modal saat click backdrop |
| `closeOnEscape` | Boolean | `true` | Close modal saat press Escape key |
| `loading` | Boolean | `false` | Loading state untuk async operations |
| `enableHaptics` | Boolean | `true` | Enable haptic feedback (vibration) |

#### Events

- `update:modelValue` - Emitted saat modal visibility berubah
- `close` - Emitted saat modal ditutup
- `cancel` - Emitted saat cancel button clicked
- `confirm` - Emitted saat confirm button clicked
- `opened` - Emitted saat modal selesai dibuka
- `closed` - Emitted saat modal selesai ditutup

#### Slots

- `default` - Main content area
- `title` - Custom title content (alternative ke title prop)
- `subtitle` - Custom subtitle content
- `footer` - Custom footer content

#### Usage Example

```vue
<template>
  <div>
    <button @click="modal.open()" class="btn-primary">
      Open Modal
    </button>

    <BaseModal
      v-model="modal.isOpen.value"
      title="Tambah Data User"
      subtitle="Isi form di bawah untuk menambah user baru"
      title-gradient
      size="md"
      :show-footer="true"
      :loading="modal.loading.value"
      @confirm="handleSubmit"
      @cancel="modal.close()"
    >
      <form @submit.prevent="handleSubmit" class="space-y-4">
        <div>
          <label class="block text-sm font-semibold text-gray-700 mb-2">
            Nama Lengkap
          </label>
          <input
            v-model="form.name"
            type="text"
            class="input-field"
            required
          />
        </div>
        <!-- More form fields -->
      </form>
    </BaseModal>
  </div>
</template>

<script setup>
import { reactive } from 'vue'
import BaseModal from '@/components/common/BaseModal.vue'
import { useModal } from '@/composables/useModal'

const modal = useModal()
const form = reactive({
  name: '',
  email: ''
})

const handleSubmit = async () => {
  modal.setLoading(true)
  
  try {
    // API call untuk save data
    await saveUser(form)
    modal.close()
  } catch (error) {
    console.error(error)
  } finally {
    modal.setLoading(false)
  }
}
</script>
```

### 2. ConfirmDialog

ConfirmDialog merupakan specialized modal untuk confirmation actions dengan visual indicators dan warning messages untuk destructive operations seperti delete atau irreversible changes.

#### Props

| Prop | Type | Default | Description |
|------|------|---------|-------------|
| `modelValue` | Boolean | `false` | Two-way binding untuk dialog visibility |
| `title` | String | `'Konfirmasi'` | Dialog title |
| `subtitle` | String | `''` | Subtitle text |
| `message` | String | `''` | Main confirmation message |
| `detail` | String | `''` | Additional detail text |
| `variant` | String | `'default'` | Visual variant: default, danger, warning, info, success |
| `icon` | String/Object | `null` | Custom icon component |
| `confirmText` | String | `'Ya, Lanjutkan'` | Confirm button text |
| `cancelText` | String | `'Batal'` | Cancel button text |
| `loadingText` | String | `'Memproses...'` | Loading state text |
| `showWarning` | Boolean | `true` | Show warning box untuk danger variant |
| `warningMessage` | String | `''` | Custom warning message |
| `loading` | Boolean | `false` | Loading state |

#### Events

- `update:modelValue` - Visibility change
- `confirm` - User confirmed action
- `cancel` - User cancelled action
- `close` - Dialog closed

#### Usage Example

```vue
<template>
  <div>
    <!-- Trigger button -->
    <button @click="handleDelete" class="btn-primary">
      Hapus Data
    </button>

    <!-- Confirm Dialog -->
    <ConfirmDialog
      v-model="confirmDialog.isOpen.value"
      :title="confirmDialog.config.value.title"
      :message="confirmDialog.config.value.message"
      :detail="confirmDialog.config.value.detail"
      :variant="confirmDialog.config.value.variant"
      :confirm-text="confirmDialog.config.value.confirmText"
      :loading="confirmDialog.loading.value"
      @confirm="confirmDialog.handleConfirm()"
      @cancel="confirmDialog.handleCancel()"
    />
  </div>
</template>

<script setup>
import ConfirmDialog from '@/components/common/ConfirmDialog.vue'
import { useConfirmDialog } from '@/composables/useModal'

const confirmDialog = useConfirmDialog()

const handleDelete = async () => {
  const confirmed = await confirmDialog.confirm({
    title: 'Hapus User',
    message: 'Apakah Anda yakin ingin menghapus user ini?',
    detail: 'User akan dihapus secara permanen dari sistem.',
    variant: 'danger',
    confirmText: 'Ya, Hapus',
    showWarning: true,
    warningMessage: 'Data yang dihapus tidak dapat dikembalikan.'
  })

  if (confirmed) {
    confirmDialog.loading.value = true
    
    try {
      await deleteUser(userId)
      // Show success alert
    } catch (error) {
      // Show error alert
    } finally {
      confirmDialog.close()
    }
  }
}
</script>
```

### 3. PasswordConfirmDialog

PasswordConfirmDialog merupakan specialized modal untuk password verification yang digunakan sebelum melakukan sensitive actions seperti delete account, change critical settings, atau approve important operations dengan security-focused design.

#### Props

| Prop | Type | Default | Description |
|------|------|---------|-------------|
| `modelValue` | Boolean | `false` | Dialog visibility |
| `title` | String | `'Konfirmasi Password'` | Dialog title |
| `subtitle` | String | `''` | Subtitle text |
| `message` | String | `'Untuk keamanan...'` | Main message |
| `detail` | String | `''` | Additional detail text |
| `confirmText` | String | `'Verifikasi'` | Confirm button text |
| `showSecurityNotice` | Boolean | `true` | Show security notice box |
| `securityNotice` | String | `''` | Custom security notice text |
| `validatePassword` | Function | `null` | Client-side validation function |
| `autoClearPassword` | Boolean | `true` | Clear password on close |

#### Events

- `update:modelValue` - Visibility change
- `confirm(password)` - User submitted password
- `cancel` - User cancelled
- `close` - Dialog closed

#### Exposed Methods

- `setError(message)` - Display error message
- `setLoading(value)` - Set loading state
- `close()` - Close dialog

#### Usage Example

```vue
<template>
  <div>
    <!-- Trigger button -->
    <button @click="handleDeleteAccount" class="btn-primary">
      Delete Account
    </button>

    <!-- Password Confirm Dialog -->
    <PasswordConfirmDialog
      ref="passwordDialogRef"
      v-model="passwordConfirm.isOpen.value"
      :title="passwordConfirm.config.value.title"
      :message="passwordConfirm.config.value.message"
      :detail="passwordConfirm.config.value.detail"
      :confirm-text="passwordConfirm.config.value.confirmText"
      :show-security-notice="passwordConfirm.config.value.showSecurityNotice"
      :security-notice="passwordConfirm.config.value.securityNotice"
      :loading="passwordConfirm.loading.value"
      @confirm="passwordConfirm.handleConfirm"
      @cancel="passwordConfirm.handleCancel"
    />
  </div>
</template>

<script setup>
import { ref } from 'vue'
import PasswordConfirmDialog from '@/components/common/PasswordConfirmDialog.vue'
import { usePasswordConfirmDialog } from '@/composables/useModal'

const passwordConfirm = usePasswordConfirmDialog()
const passwordDialogRef = ref(null)

// Set dialog reference untuk direct control
passwordConfirm.setDialogRef(passwordDialogRef)

const handleDeleteAccount = async () => {
  const password = await passwordConfirm.confirm({
    title: 'Hapus Account',
    message: 'Untuk menghapus account, masukkan password Anda.',
    detail: 'Semua data akan dihapus secara permanen.',
    confirmText: 'Hapus Account',
    showSecurityNotice: true,
    securityNotice: '⚠️ Tindakan ini tidak dapat dibatalkan.'
  })

  if (password) {
    passwordConfirm.setLoading(true)
    
    try {
      // Verify password dengan backend
      const response = await api.verifyPassword({ password })
      
      if (response.verified) {
        // Password benar, lanjutkan delete
        await api.deleteAccount()
        passwordConfirm.close()
        
        // Show success
        alert.success('Account berhasil dihapus')
      } else {
        // Password salah
        passwordConfirm.setError('Password yang Anda masukkan salah.')
      }
    } catch (error) {
      passwordConfirm.setError('Terjadi kesalahan saat verifikasi.')
    }
  } else {
    // User cancelled
    console.log('User cancelled password confirmation')
  }
}
</script>
```

#### Security Best Practices

**1. Always Verify on Backend**
```javascript
// ❌ JANGAN verify password di frontend
if (password === userPassword) { ... }

// ✅ SELALU verify di backend
const response = await api.verifyPassword({ password })
if (response.verified) { ... }
```

**2. Use HTTPS**
Pastikan semua password verification dilakukan melalui HTTPS untuk encrypt data in transit.

**3. Rate Limiting**
Implement rate limiting di backend untuk prevent brute force attacks.

**4. Clear Sensitive Data**
Password akan otomatis di-clear saat dialog ditutup (via `autoClearPassword` prop).

**5. Show Clear Error Messages**
```javascript
if (!verified) {
  passwordConfirm.setError('Password salah. Silakan coba lagi.')
}
```

### 4. AlertDialog

AlertDialog merupakan notification dialog untuk menampilkan success, error, warning, atau info messages dengan animated icons dan optional auto-dismiss functionality untuk quick notifications.

#### Props

| Prop | Type | Default | Description |
|------|------|---------|-------------|
| `modelValue` | Boolean | `false` | Dialog visibility |
| `title` | String | `''` | Dialog title |
| `message` | String | `''` (required) | Main message text |
| `detail` | String | `''` | Additional detail text |
| `variant` | String | `'success'` | Alert type: success, error, warning, info |
| `icon` | String/Object | `null` | Custom icon |
| `confirmText` | String | `'OK'` | Button text |
| `autoDismiss` | Boolean | `false` | Enable auto-dismiss |
| `autoDismissDelay` | Number | `3000` | Auto-dismiss delay (ms) |
| `showTimer` | Boolean | `true` | Show countdown timer |
| `showContent` | Boolean | `false` | Enable content slot |
| `dismissible` | Boolean | `true` | Allow manual dismiss |

#### Events

- `update:modelValue` - Visibility change
- `close` - Dialog closed
- `dismissed` - Auto-dismissed

#### Usage Example

```vue
<template>
  <div>
    <!-- Alert Dialog -->
    <AlertDialog
      v-model="alertDialog.isOpen.value"
      :title="alertDialog.config.value.title"
      :message="alertDialog.config.value.message"
      :detail="alertDialog.config.value.detail"
      :variant="alertDialog.config.value.variant"
      :auto-dismiss="alertDialog.config.value.autoDismiss"
      @close="alertDialog.handleClose()"
    />
  </div>
</template>

<script setup>
import AlertDialog from '@/components/common/AlertDialog.vue'
import { useAlertDialog } from '@/composables/useModal'

const alertDialog = useAlertDialog()

// Success alert
const showSuccess = () => {
  alertDialog.success('Data Berhasil Disimpan!', {
    detail: 'User telah berhasil ditambahkan ke sistem.',
    title: 'Sukses'
  })
}

// Error alert
const showError = () => {
  alertDialog.error('Terjadi Kesalahan!', {
    detail: 'Gagal menyimpan data. Silakan coba lagi.',
    title: 'Error'
  })
}

// Auto-dismiss alert
const showAutoAlert = () => {
  alertDialog.success('Berhasil!', {
    detail: 'Dialog akan menutup otomatis.',
    autoDismiss: true,
    autoDismissDelay: 3000
  })
}
</script>
```

## Composables

### useModal()

Composable untuk managing basic modal state dan actions.

```javascript
const modal = useModal()

// Properties
modal.isOpen    // ref<boolean>
modal.loading   // ref<boolean>

// Methods
modal.open()              // Open modal
modal.close()             // Close modal
modal.toggle()            // Toggle modal
modal.setLoading(true)    // Set loading state
```

### useConfirmDialog()

Composable dengan promise-based API untuk confirmation dialogs.

```javascript
const confirmDialog = useConfirmDialog()

// Properties
confirmDialog.isOpen      // ref<boolean>
confirmDialog.loading     // ref<boolean>
confirmDialog.config      // ref<object>

// Methods
const result = await confirmDialog.confirm({
  title: 'Konfirmasi',
  message: 'Lanjutkan?',
  variant: 'danger'
})
// result: true (confirmed) atau false (cancelled)

confirmDialog.handleConfirm()  // Resolve promise dengan true
confirmDialog.handleCancel()   // Resolve promise dengan false
confirmDialog.close()          // Close dialog
```

### useAlertDialog()

Composable untuk alert notifications dengan shorthand methods.

```javascript
const alertDialog = useAlertDialog()

// Properties
alertDialog.isOpen     // ref<boolean>
alertDialog.config     // ref<object>

// Methods
await alertDialog.alert({ message, variant })  // Generic alert
await alertDialog.success('Success!')          // Success variant
await alertDialog.error('Error!')              // Error variant
await alertDialog.warning('Warning!')          // Warning variant
await alertDialog.info('Info!')                // Info variant
alertDialog.handleClose()                      // Close dialog
```

### usePasswordConfirmDialog()

Composable untuk password confirmation dengan promise-based API dan error handling.

```javascript
const passwordConfirm = usePasswordConfirmDialog()
const passwordDialogRef = ref(null)

// Set dialog reference
passwordConfirm.setDialogRef(passwordDialogRef)

// Properties
passwordConfirm.isOpen      // ref<boolean>
passwordConfirm.loading     // ref<boolean>
passwordConfirm.config      // ref<object>

// Methods
const password = await passwordConfirm.confirm({
  title: 'Konfirmasi Password',
  message: 'Masukkan password Anda',
  confirmText: 'Verifikasi'
})
// Returns: password string (confirmed) atau null (cancelled)

passwordConfirm.handleConfirm(password)  // Handle confirm dengan password
passwordConfirm.handleCancel()           // Handle cancel
passwordConfirm.setError(message)        // Display error message
passwordConfirm.setLoading(true)         // Set loading state
passwordConfirm.close()                  // Close dialog
passwordConfirm.setDialogRef(ref)        // Set dialog reference
```

## Design Principles

### iOS-Inspired Animations

Modal system menggunakan spring physics animations yang memberikan natural, bouncy feel seperti native iOS apps dengan transition parameters:

```javascript
{
  duration: 0.4,
  type: 'spring',
  stiffness: 300,
  damping: 30
}
```

### Mobile-First Responsive

Modal behavior berbeda antara mobile dan desktop views:

**Mobile (< 640px):**
- Bottom sheet style dengan rounded top corners
- Slide up animation dari bottom
- Drag handle untuk visual affordance
- Full width dengan max-height 90vh

**Desktop (≥ 640px):**
- Centered modal dengan shadow
- Scale + fade animation
- Fixed max-width berdasarkan size prop
- Glass morphism backdrop

### Glass Morphism Effect

Menggunakan backdrop-filter untuk frosted glass appearance:

```css
backdrop-filter: blur(16px) saturate(180%);
-webkit-backdrop-filter: blur(16px) saturate(180%);
background-color: rgba(255, 255, 255, 0.9);
```

### Haptic Feedback

Vibration API digunakan untuk tactile response pada mobile devices dengan patterns:

- **Light** (10ms) - Open/close actions
- **Medium** (20ms) - Confirm actions
- **Heavy** (30ms) - Danger/destructive actions

### Press Feedback

Interactive elements menggunakan scale transform saat active:

```css
active:scale-[0.97]
```

## Best Practices

### 1. Loading States

Selalu set loading state untuk async operations:

```javascript
const handleSubmit = async () => {
  modal.setLoading(true)
  try {
    await saveData()
  } finally {
    modal.setLoading(false)
  }
}
```

### 2. Confirmation untuk Destructive Actions

Gunakan ConfirmDialog dengan variant='danger' untuk delete/destructive operations:

```javascript
const confirmed = await confirmDialog.confirm({
  variant: 'danger',
  showWarning: true,
  warningMessage: 'Tindakan tidak dapat dibatalkan'
})
```

### 3. Success Feedback

Tampilkan alert setelah successful operations:

```javascript
await alertDialog.success('Data berhasil disimpan!', {
  autoDismiss: true,
  autoDismissDelay: 2000
})
```

### 4. Error Handling

Tangkap errors dan tampilkan error alerts:

```javascript
try {
  await saveData()
} catch (error) {
  alertDialog.error('Gagal menyimpan data', {
    detail: error.message
  })
}
```

### 5. Mobile UX

Perhatikan mobile user experience:
- Gunakan size='sm' atau 'md' untuk mobile
- Enable haptic feedback
- Gunakan clear, concise text
- Avoid complex forms di modal untuk mobile

### 6. Password Confirmation untuk Sensitive Actions

Gunakan PasswordConfirmDialog untuk actions yang memerlukan extra security:

```javascript
// Delete account, change email, transfer money, etc.
const password = await passwordConfirm.confirm({
  title: 'Konfirmasi Identitas',
  message: 'Masukkan password untuk melanjutkan.'
})

if (password) {
  // Verify dengan backend
  const verified = await verifyPassword(password)
  if (verified) {
    // Lanjutkan action
  } else {
    passwordConfirm.setError('Password salah')
  }
}
```

## Troubleshooting

### Modal tidak menutup

**Problem:** Modal tetap terbuka setelah action

**Solution:** 
```javascript
// Pastikan loading state di-reset
modal.setLoading(false)
// Manual close jika perlu
modal.close()
```

### Backdrop blur tidak bekerja

**Problem:** Glass effect tidak muncul di beberapa browser

**Solution:** Sudah include webkit prefix:
```css
-webkit-backdrop-filter: blur(16px);
backdrop-filter: blur(16px);
```

### Haptic feedback tidak bekerja

**Problem:** Vibration API tidak tersedia

**Solution:** Feature detection sudah built-in:
```javascript
if ('vibrate' in navigator) {
  navigator.vibrate(pattern)
}
```

### Multiple modals z-index issue

**Problem:** Modal stacking tidak bekerja dengan benar

**Solution:** Gunakan modalManager dari composable:
```javascript
import { modalManager } from '@/composables/useModal'
```

## Testing

Lihat `ModalExamples.vue` untuk interactive testing dan examples dari semua modal variants dan configurations.

**Run test:**
```bash
yarn dev
# Navigate ke /modal-examples route
```

## Changelog

### v1.0.0 (Initial Release)
- BaseModal component dengan iOS-inspired design
- ConfirmDialog untuk confirmations
- AlertDialog untuk notifications
- useModal, useConfirmDialog, useAlertDialog composables
- Mobile-first responsive design
- Haptic feedback support
- Spring physics animations
- Glass morphism effects
- Auto-dismiss functionality
- Comprehensive documentation

## Contributing

Untuk improvements atau bug fixes:

1. Test perubahan dengan ModalExamples.vue
2. Update documentation jika ada API changes
3. Pastikan mobile responsiveness
4. Test di iOS dan Android devices
5. Verify haptic feedback
6. Check accessibility (keyboard navigation, screen readers)

## Support

Developer: Zulfikar Hidayatullah (+62 857-1583-8733)
