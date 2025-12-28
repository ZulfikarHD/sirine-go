# Common Components

## Overview

Common components merupakan kumpulan reusable UI components yang digunakan di seluruh aplikasi, yaitu: modal system, dialogs, breadcrumbs, dan shared UI elements dengan consistent design language dan behavior patterns.

## Components

### Modal System

#### BaseModal
Core modal component untuk CRUD operations, form submissions, dan general-purpose dialogs.

```javascript
import { BaseModal } from '@/components/common'
```

**Use Cases:**
- Form modals (create/edit)
- Detail views
- Multi-step wizards
- Content previews

#### ConfirmDialog
Specialized confirmation dialog untuk user action confirmations.

```javascript
import { ConfirmDialog } from '@/components/common'
```

**Use Cases:**
- Delete confirmations
- Destructive action warnings
- Navigation confirmations
- Save/discard prompts

#### AlertDialog
Notification dialog untuk success, error, warning, dan info messages.

```javascript
import { AlertDialog } from '@/components/common'
```

**Use Cases:**
- Success notifications
- Error messages
- Warning alerts
- Info announcements

### Navigation

#### Breadcrumbs
Breadcrumb navigation component untuk hierarchical navigation.

```javascript
import { Breadcrumbs } from '@/components/common'
```

## Quick Import

```javascript
// Import individual components
import { 
  BaseModal, 
  ConfirmDialog, 
  AlertDialog,
  Breadcrumbs 
} from '@/components/common'

// Import composables
import { 
  useModal, 
  useConfirmDialog, 
  useAlertDialog 
} from '@/composables/useModal'
```

## Documentation

- **Modal System Full Docs:** `/docs/components/modal-system.md`
- **Quick Start Guide:** `/docs/components/QUICK_START_MODAL.md`
- **Interactive Examples:** Navigate to `/dev/modal-examples`

## Design Principles

Semua components mengikuti design standards:
- ✅ iOS-inspired animations
- ✅ Indigo & Fuchsia color theme
- ✅ Mobile-first responsive
- ✅ Glass morphism effects
- ✅ Haptic feedback
- ✅ Press feedback
- ✅ Spring physics animations

## Development

### Adding New Components

1. Create component file di `/components/common/`
2. Follow naming convention: PascalCase
3. Add comprehensive JSDoc comments
4. Export dari `index.js`
5. Add documentation
6. Create usage examples

### Testing Components

Use ModalExamples.vue sebagai template untuk creating interactive component demos.

## Support

Developer: Zulfikar Hidayatullah (+62 857-1583-8733)
