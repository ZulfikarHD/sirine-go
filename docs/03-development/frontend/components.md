# 🧩 Component Development

Panduan untuk membuat Vue components yang reusable, maintainable, dan mengikuti best practices dalam Sirine Go App.

---

## 📋 Overview

Components merupakan building blocks aplikasi Vue, yaitu reusable pieces yang combine template, logic, dan styling untuk create UI elements.

### Component Types

**Base Components** - Generic reusable UI elements:
- Buttons, inputs, modals, cards
- No business logic
- Highly customizable via props

**Feature Components** - Specific to features:
- User cards, notification bells, achievement badges
- Contains business logic
- May call API atau use composables

**Layout Components** - Application structure:
- Navbar, sidebar, footer
- Define app layout dan navigation
- Usually persistent across pages

**Page Components** - Top-level Inertia pages:
- Dashboard, user list, profile pages
- Compose smaller components
- Receive props dari backend

---

## 🏗️ Component Structure

### Basic Component Template

```vue
<script setup>
/**
 * ComponentName untuk [purpose] dengan [approach]
 * digunakan di [where] untuk [benefit]
 */
import { ref, computed, onMounted } from 'vue'
import { Motion } from 'motion-v'

// Props - Input dari parent component
const props = defineProps({
  title: {
    type: String,
    required: true,
    validator: (value) => value.length > 0
  },
  variant: {
    type: String,
    default: 'primary',
    validator: (value) => ['primary', 'secondary', 'danger'].includes(value)
  },
  disabled: {
    type: Boolean,
    default: false
  }
})

// Emits - Events ke parent component
const emit = defineEmits(['click', 'update:modelValue'])

// State - Internal reactive state
const isActive = ref(false)
const localValue = ref('')

// Computed - Derived state
const buttonClasses = computed(() => ({
  'btn-primary': props.variant === 'primary',
  'btn-secondary': props.variant === 'secondary',
  'btn-danger': props.variant === 'danger',
  'opacity-50 cursor-not-allowed': props.disabled
}))

// Methods - Component logic
const handleClick = () => {
  if (props.disabled) return
  emit('click')
}

// Lifecycle - Component lifecycle hooks
onMounted(() => {
  console.log('Component mounted')
})
</script>

<template>
  <Motion
    :initial="{ opacity: 0, y: 10 }"
    :animate="{ opacity: 1, y: 0 }"
    :transition="{ duration: 0.2, ease: 'easeOut' }"
  >
    <button
      :class="buttonClasses"
      :disabled="disabled"
      @click="handleClick"
      class="active-scale px-4 py-2 rounded-lg transition-colors"
    >
      <slot>{{ title }}</slot>
    </button>
  </Motion>
</template>
```

---

## 🎨 Base Components Examples

### BaseButton Component

```vue
<script setup>
/**
 * BaseButton merupakan reusable button component dengan variants
 * untuk consistent styling across aplikasi
 */
import { computed } from 'vue'
import { Motion } from 'motion-v'

const props = defineProps({
  variant: {
    type: String,
    default: 'primary',
    validator: (v) => ['primary', 'secondary', 'danger', 'ghost'].includes(v)
  },
  size: {
    type: String,
    default: 'md',
    validator: (v) => ['sm', 'md', 'lg'].includes(v)
  },
  disabled: Boolean,
  loading: Boolean,
  block: Boolean  // Full width
})

const emit = defineEmits(['click'])

const buttonClasses = computed(() => {
  const classes = ['active-scale transition-all duration-200']
  
  // Variant styles
  const variants = {
    primary: 'bg-gradient-to-r from-indigo-600 to-fuchsia-600 text-white hover:from-indigo-700 hover:to-fuchsia-700',
    secondary: 'bg-gray-100 text-gray-900 hover:bg-gray-200',
    danger: 'bg-red-600 text-white hover:bg-red-700',
    ghost: 'bg-transparent text-indigo-600 hover:bg-indigo-50'
  }
  classes.push(variants[props.variant])
  
  // Size styles
  const sizes = {
    sm: 'px-3 py-1.5 text-sm',
    md: 'px-4 py-2 text-base',
    lg: 'px-6 py-3 text-lg'
  }
  classes.push(sizes[props.size])
  
  // States
  if (props.disabled || props.loading) {
    classes.push('opacity-50 cursor-not-allowed')
  }
  
  if (props.block) {
    classes.push('w-full')
  }
  
  classes.push('rounded-lg font-medium')
  
  return classes.join(' ')
})

const handleClick = (e) => {
  if (!props.disabled && !props.loading) {
    emit('click', e)
  }
}
</script>

<template>
  <button
    :class="buttonClasses"
    :disabled="disabled || loading"
    @click="handleClick"
    type="button"
  >
    <span v-if="loading" class="inline-block mr-2">
      <svg class="animate-spin h-4 w-4" viewBox="0 0 24 24">
        <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" fill="none"/>
        <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"/>
      </svg>
    </span>
    <slot />
  </button>
</template>
```

**Usage:**
```vue
<BaseButton variant="primary" size="md" @click="handleSubmit">
  Submit
</BaseButton>

<BaseButton variant="danger" :loading="isDeleting" @click="handleDelete">
  Delete User
</BaseButton>
```

### BaseInput Component

```vue
<script setup>
/**
 * BaseInput untuk form input dengan validation styling
 * dan error message display
 */
import { computed } from 'vue'

const props = defineProps({
  modelValue: [String, Number],
  type: {
    type: String,
    default: 'text'
  },
  label: String,
  placeholder: String,
  error: String,
  disabled: Boolean,
  required: Boolean
})

const emit = defineEmits(['update:modelValue'])

const inputClasses = computed(() => {
  const classes = [
    'w-full px-4 py-2 rounded-lg border transition-all duration-200',
    'focus:outline-none focus:ring-4'
  ]
  
  if (props.error) {
    classes.push('border-red-500 focus:border-red-500 focus:ring-red-100')
  } else {
    classes.push('border-gray-300 focus:border-indigo-500 focus:ring-indigo-100')
  }
  
  if (props.disabled) {
    classes.push('bg-gray-100 cursor-not-allowed')
  } else {
    classes.push('bg-white')
  }
  
  return classes.join(' ')
})

const handleInput = (event) => {
  emit('update:modelValue', event.target.value)
}
</script>

<template>
  <div class="flex flex-col gap-1">
    <!-- Label -->
    <label v-if="label" class="text-sm font-medium text-gray-700">
      {{ label }}
      <span v-if="required" class="text-red-500">*</span>
    </label>
    
    <!-- Input -->
    <input
      :type="type"
      :value="modelValue"
      :placeholder="placeholder"
      :disabled="disabled"
      :class="inputClasses"
      @input="handleInput"
    />
    
    <!-- Error message -->
    <p v-if="error" class="text-sm text-red-600">
      {{ error }}
    </p>
  </div>
</template>
```

**Usage:**
```vue
<script setup>
import { ref } from 'vue'

const email = ref('')
const emailError = ref('')

const validateEmail = () => {
  if (!email.value.includes('@')) {
    emailError.value = 'Email tidak valid'
  } else {
    emailError.value = ''
  }
}
</script>

<template>
  <BaseInput
    v-model="email"
    type="email"
    label="Email Address"
    placeholder="you@example.com"
    :error="emailError"
    required
    @blur="validateEmail"
  />
</template>
```

### BaseModal Component

```vue
<script setup>
/**
 * BaseModal untuk reusable modal/dialog dengan backdrop
 * dan Motion-V animations
 */
import { Motion } from 'motion-v'
import { onMounted, onUnmounted } from 'vue'

const props = defineProps({
  show: {
    type: Boolean,
    required: true
  },
  title: String,
  size: {
    type: String,
    default: 'md',
    validator: (v) => ['sm', 'md', 'lg', 'xl'].includes(v)
  },
  closeOnBackdrop: {
    type: Boolean,
    default: true
  }
})

const emit = defineEmits(['close'])

const sizeClasses = {
  sm: 'max-w-md',
  md: 'max-w-lg',
  lg: 'max-w-2xl',
  xl: 'max-w-4xl'
}

const handleBackdropClick = () => {
  if (props.closeOnBackdrop) {
    emit('close')
  }
}

const handleEscape = (e) => {
  if (e.key === 'Escape' && props.show) {
    emit('close')
  }
}

onMounted(() => {
  document.addEventListener('keydown', handleEscape)
})

onUnmounted(() => {
  document.removeEventListener('keydown', handleEscape)
})
</script>

<template>
  <Teleport to="body">
    <Transition>
      <div v-if="show" class="fixed inset-0 z-50 flex items-center justify-center p-4">
        <!-- Backdrop -->
        <Motion
          :initial="{ opacity: 0 }"
          :animate="{ opacity: 1 }"
          :exit="{ opacity: 0 }"
          :transition="{ duration: 0.2, ease: 'easeOut' }"
          class="fixed inset-0 bg-black/50 backdrop-blur-sm"
          @click="handleBackdropClick"
        />
        
        <!-- Modal Container -->
        <Motion
          :initial="{ opacity: 0, scale: 0.95, y: 20 }"
          :animate="{ opacity: 1, scale: 1, y: 0 }"
          :exit="{ opacity: 0, scale: 0.95, y: 20 }"
          :transition="{ type: 'spring', stiffness: 500, damping: 40, mass: 0.8 }"
          :class="['relative bg-white rounded-2xl shadow-2xl w-full', sizeClasses[size]]"
        >
          <!-- Header -->
          <div class="flex items-center justify-between p-6 border-b border-gray-200">
            <h3 class="text-xl font-semibold text-gray-900">
              {{ title }}
            </h3>
            <button
              @click="emit('close')"
              class="p-1 hover:bg-gray-100 rounded-lg transition-colors active-scale"
            >
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
              </svg>
            </button>
          </div>
          
          <!-- Body -->
          <div class="p-6">
            <slot />
          </div>
          
          <!-- Footer -->
          <div v-if="$slots.footer" class="p-6 border-t border-gray-200">
            <slot name="footer" />
          </div>
        </Motion>
      </div>
    </Transition>
  </Teleport>
</template>
```

**Usage:**
```vue
<script setup>
import { ref } from 'vue'

const showModal = ref(false)
</script>

<template>
  <BaseButton @click="showModal = true">
    Open Modal
  </BaseButton>
  
  <BaseModal
    :show="showModal"
    title="Confirm Delete"
    size="md"
    @close="showModal = false"
  >
    <p>Are you sure you want to delete this user?</p>
    
    <template #footer>
      <div class="flex gap-3 justify-end">
        <BaseButton variant="ghost" @click="showModal = false">
          Cancel
        </BaseButton>
        <BaseButton variant="danger" @click="handleDelete">
          Delete
        </BaseButton>
      </div>
    </template>
  </BaseModal>
</template>
```

---

## 🔧 Composables Pattern

Composables merupakan reusable stateful logic yang dapat digunakan across multiple components.

### useToast Composable

```javascript
// src/composables/useToast.js
import { ref } from 'vue'

/**
 * useToast untuk global toast notification system
 * yang manage toast messages dengan auto-dismiss
 */

const toasts = ref([])
let nextId = 0

export function useToast() {
  /**
   * add untuk menambah toast notification dengan auto-dismiss
   * @param {String} message - Message text
   * @param {String} type - Toast type (success, error, info, warning)
   * @param {Number} duration - Duration dalam ms sebelum auto-dismiss
   */
  const add = (message, type = 'info', duration = 3000) => {
    const id = nextId++
    
    toasts.value.push({
      id,
      message,
      type,
      timestamp: Date.now()
    })
    
    // Auto-dismiss setelah duration
    if (duration > 0) {
      setTimeout(() => {
        remove(id)
      }, duration)
    }
    
    return id
  }
  
  const remove = (id) => {
    const index = toasts.value.findIndex(t => t.id === id)
    if (index > -1) {
      toasts.value.splice(index, 1)
    }
  }
  
  // Shorthand methods
  const success = (msg, duration) => add(msg, 'success', duration)
  const error = (msg, duration) => add(msg, 'error', duration)
  const info = (msg, duration) => add(msg, 'info', duration)
  const warning = (msg, duration) => add(msg, 'warning', duration)
  
  return {
    toasts,
    add,
    remove,
    success,
    error,
    info,
    warning
  }
}
```

**Usage:**
```vue
<script setup>
import { useToast } from '@/composables/useToast'

const toast = useToast()

const handleSubmit = async () => {
  try {
    await api.post('/users', formData.value)
    toast.success('User berhasil dibuat')
  } catch (error) {
    toast.error('Gagal membuat user')
  }
}
</script>
```

### useModal Composable

```javascript
// src/composables/useModal.js
import { ref } from 'vue'

/**
 * useModal untuk manage modal state dengan multiple modals support
 */
export function useModal() {
  const modals = ref({})
  
  const open = (name) => {
    modals.value[name] = true
  }
  
  const close = (name) => {
    modals.value[name] = false
  }
  
  const toggle = (name) => {
    modals.value[name] = !modals.value[name]
  }
  
  const isOpen = (name) => {
    return !!modals.value[name]
  }
  
  return {
    modals,
    open,
    close,
    toggle,
    isOpen
  }
}
```

**Usage:**
```vue
<script setup>
import { useModal } from '@/composables/useModal'

const modal = useModal()
</script>

<template>
  <BaseButton @click="modal.open('deleteUser')">
    Delete
  </BaseButton>
  
  <BaseModal
    :show="modal.isOpen('deleteUser')"
    @close="modal.close('deleteUser')"
  >
    <!-- Modal content -->
  </BaseModal>
</template>
```

---

## ✅ Best Practices

### 1. Props Validation

```vue
<script setup>
// Do: Validate props dengan detailed validators
const props = defineProps({
  status: {
    type: String,
    required: true,
    validator: (value) => ['active', 'inactive', 'pending'].includes(value)
  }
})

// Don't: Skip validation
const props = defineProps({
  status: String  // No validation
})
</script>
```

### 2. Event Naming

```vue
<script setup>
// Do: Use kebab-case untuk custom events
emit('update:modelValue', newValue)
emit('user-selected', user)

// Don't: Use camelCase
emit('updateModelValue', newValue)  // Tidak konsisten dengan HTML
</script>
```

### 3. Component Size

- Keep components < 250 lines
- Extract complex logic ke composables
- Split large components ke smaller ones

### 4. Slots for Flexibility

```vue
<!-- Do: Use slots untuk flexible content -->
<BaseCard>
  <template #header>
    <h2>Custom Header</h2>
  </template>
  
  <p>Card content</p>
  
  <template #footer>
    <BaseButton>Action</BaseButton>
  </template>
</BaseCard>

<!-- Component definition -->
<template>
  <div class="card">
    <div v-if="$slots.header" class="card-header">
      <slot name="header" />
    </div>
    <div class="card-body">
      <slot />
    </div>
    <div v-if="$slots.footer" class="card-footer">
      <slot name="footer" />
    </div>
  </div>
</template>
```

### 5. Provide/Inject for Deep Props

```vue
<!-- Parent component -->
<script setup>
import { provide } from 'vue'

const theme = ref('light')
provide('theme', theme)
</script>

<!-- Deep child component -->
<script setup>
import { inject } from 'vue'

const theme = inject('theme')
</script>
```

---

## 📚 Related Documentation

- [Getting Started Guide](./getting-started.md) - Setup environment
- [Animation Guide](./animations.md) - Motion-V animations
- [Styling Guide](./styling.md) - Design system
- [API Reference](../../04-api-reference/README.md) - Backend endpoints

---

**Last Updated:** 28 Desember 2025  
**Status:** ✅ Production Ready
