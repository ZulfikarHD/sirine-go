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

// Other Common Components
export { default as Breadcrumbs } from './Breadcrumbs.vue'

// Example Component (untuk development/testing)
export { default as ModalExamples } from './ModalExamples.vue'
