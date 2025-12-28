# Common Components - Modal System

## Overview

Modal system dengan iOS-inspired design yang mengimplementasikan:
- **Spring Physics**: Natural, bouncy animations dengan motion-v
- **Glass Morphism**: Frosted glass effect dengan backdrop blur
- **Indigo-Fuchsia Theme**: Gradient theme yang modern dan techy
- **Press Feedback**: Scale-down effect saat element di-tap (0.97 scale)
- **Haptic Feedback**: Vibration API untuk tactile response
- **Mobile-First**: Bottom sheet style pada mobile dengan drag handle

## Components

### 1. BaseModal
Modal dasar yang reusable untuk CRUD operations dan custom content.

```vue
<BaseModal
  v-model="isOpen"
  title="Tambah Data"
  subtitle="Isi form berikut"
  title-gradient
  size="md"
  :show-footer="true"
  :loading="loading"
  @confirm="handleSubmit"
  @cancel="handleCancel"
>
  <!-- Content -->
</BaseModal>
```

**Props:**
- `size`: 'xs' | 'sm' | 'md' | 'lg' | 'xl' | 'full'
- `title-gradient`: Boolean - Menggunakan gradient Indigo-Fuchsia
- `confirm-danger`: Boolean - Menggunakan danger button style

### 2. ConfirmDialog
Dialog konfirmasi dengan animated icons dan variant-based styling.

```vue
<ConfirmDialog
  v-model="isOpen"
  title="Hapus Data"
  message="Yakin ingin menghapus?"
  variant="danger"
  :show-warning="true"
  @confirm="handleDelete"
/>
```

**Variants:**
- `default`: Indigo theme dengan HelpCircle icon
- `danger`: Red-Pink gradient dengan Trash icon
- `warning`: Amber theme dengan AlertTriangle icon
- `success`: Emerald theme dengan CheckCircle icon
- `info`: Blue theme dengan AlertCircle icon

### 3. AlertDialog
Alert/notification dialog dengan auto-dismiss dan visual timer.

```vue
<AlertDialog
  v-model="isOpen"
  message="Data berhasil disimpan!"
  variant="success"
  :auto-dismiss="true"
  :auto-dismiss-delay="3000"
/>
```

**Variants:**
- `success`: Emerald dengan pulse animation
- `error`: Red dengan XCircle icon
- `warning`: Amber dengan AlertTriangle icon
- `info`: Blue dengan Info icon

### 4. PasswordConfirmDialog
Password verification dialog untuk sensitive actions.

```vue
<PasswordConfirmDialog
  v-model="isOpen"
  title="Konfirmasi Password"
  message="Masukkan password untuk melanjutkan"
  :show-security-notice="true"
  @confirm="handleVerify"
/>
```

## Design Principles

### iOS-Inspired Animations
- Spring physics dengan stiffness: 380, damping: 32
- Staggered entrance animations dengan delay
- Pulse rings untuk success states
- Shake animation untuk error states

### Glass Morphism
- Backdrop blur: 8-20px
- Semi-transparent backgrounds
- Subtle borders dengan opacity
- Gradient shadows

### Theme Colors
- Primary: Indigo (#6366f1) to Fuchsia (#d946ef)
- Success: Emerald (#10b981)
- Danger: Red (#ef4444) to Pink (#ec4899)
- Warning: Amber (#f59e0b)
- Info: Blue (#3b82f6)

## Composables

Gunakan composables untuk programmatic control:

```javascript
import { useModal, useConfirmDialog, useAlertDialog } from '@/composables/useModal'

const modal = useModal()
const confirm = useConfirmDialog()
const alert = useAlertDialog()

// Confirm dialog
const result = await confirm.confirm({
  title: 'Hapus Data',
  message: 'Yakin?',
  variant: 'danger'
})

// Alert dialog
await alert.success('Berhasil!')
await alert.error('Gagal!')
```

## CSS Classes

Semua styles didefinisikan di `src/style.css`:
- `.modal-backdrop` - Frosted glass backdrop
- `.modal-container` - Glass card container
- `.modal-header` - Header dengan gradient accent
- `.modal-btn-primary` - Gradient primary button
- `.modal-btn-danger` - Gradient danger button
- `.dialog-icon-*` - Icon container variants

## Developer

Zulfikar Hidayatullah (+62 857-1583-8733)
