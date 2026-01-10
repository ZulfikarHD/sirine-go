/**
 * Common Components Exports
 * 
 * Centralized exports untuk common/shared components yang mencakup
 * modal system, dialogs, dan reusable UI components untuk
 * consistent imports across application
 */

// Modal Components
export { default as BaseModal } from './BaseModal.vue'
export { default as ConfirmDialog } from './ConfirmDialog.vue'
export { default as AlertDialog } from './AlertDialog.vue'
export { default as PasswordConfirmDialog } from './PasswordConfirmDialog.vue'

// Navigation & Layout Components
export { default as Breadcrumbs } from './Breadcrumbs.vue'

// UI Components
export { default as PriorityBadge } from './PriorityBadge.vue'
export { default as LoadingSkeleton } from './LoadingSkeleton.vue'
export { default as BarcodeScanner } from './BarcodeScanner.vue'
export { default as PhotoUploader } from './PhotoUploader.vue'
export { default as MaterialPhotoViewer } from './MaterialPhotoViewer.vue'

