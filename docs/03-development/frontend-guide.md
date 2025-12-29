# 🎨 Frontend Development Guide

Panduan comprehensive untuk frontend development dalam Sirine Go App menggunakan Vue 3, Vite, dan Motion-V.

**Tech Stack:**
- Vue 3 (Composition API)
- Vite (Build tool)
- Wayfinder (Routing untuk Inertia)
- Motion-V (Animations)
- Tailwind CSS (Styling)
- Axios (HTTP client)

**Last Updated:** 28 Desember 2025

---

## 📋 Overview

Frontend Sirine Go App dibangun dengan Vue 3 Composition API untuk reactive dan maintainable code, dengan focus pada:
- Apple-inspired design dengan Indigo & Fuchsia gradient
- Mobile-first UX approach
- Motion-V untuk smooth iOS-like animations
- Minimal visual weight, maksimal UX

---

## 🚀 Quick Start

### Prerequisites

**Required:**
- Node.js v18+ dan Yarn
- Code editor (VS Code recommended)
- Basic knowledge Vue 3 Composition API

### Setup Development Environment

```bash
# Navigate to frontend folder
cd frontend

# Install dependencies
yarn install

# Start dev server
yarn dev

# Server akan running di http://localhost:5173
```

### Development Commands

```bash
# Start dev server
yarn dev

# Build for production
yarn build

# Preview production build
yarn preview

# Run linter
yarn lint

# Run tests
yarn test
```

---

## 📁 Project Structure

```
frontend/
├── src/
│   ├── assets/          # Static assets (images, fonts)
│   ├── components/      # Reusable Vue components
│   │   ├── base/        # Base components (Button, Input, Modal)
│   │   ├── layout/      # Layout components (Navbar, Sidebar)
│   │   └── features/    # Feature-specific components
│   ├── composables/     # Composition API composables
│   │   ├── useAuth.js   # Authentication composable
│   │   ├── useToast.js  # Toast notifications
│   │   ├── useMotion.js # Animation presets
│   │   └── useModal.js  # Modal management
│   ├── pages/           # Page components (Inertia pages)
│   │   ├── Auth/        # Auth pages (Login, ForgotPassword)
│   │   ├── Dashboard/   # Dashboard pages
│   │   ├── Users/       # User management pages
│   │   └── Profile/     # Profile pages
│   ├── utils/           # Utility functions
│   │   ├── axios.js     # Axios configuration
│   │   ├── helpers.js   # Helper functions
│   │   └── validators.js # Validation utilities
│   ├── style.css        # Global styles (minimal)
│   └── app.js           # Main app entry
├── public/              # Public static files
├── tailwind.config.js   # Tailwind configuration
├── vite.config.js       # Vite configuration
└── package.json         # Dependencies
```

---

## 🎨 Design System

### Theme Colors

```javascript
// tailwind.config.js
colors: {
  primary: {
    DEFAULT: '#6366f1',  // Indigo
    light: '#818cf8',
    dark: '#4f46e5',
  },
  secondary: {
    DEFAULT: '#d946ef',  // Fuchsia
    light: '#e879f9',
    dark: '#c026d3',
  },
  // ... other colors
}
```

### Typography

```css
/* Base font: Inter */
font-family: 'Inter', system-ui, sans-serif;

/* Font sizes (Tailwind) */
text-xs    /* 12px */
text-sm    /* 14px */
text-base  /* 16px */
text-lg    /* 18px */
text-xl    /* 20px */
```

### Spacing System

Menggunakan Tailwind spacing scale (4px base):
- `p-1` = 4px
- `p-2` = 8px
- `p-4` = 16px
- `p-6` = 24px
- `p-8` = 32px

---

## 🎬 Animation System (Motion-V)

### Core Principle

**SEMUA animasi menggunakan Motion-V**, bukan CSS animations.

### Import Pattern

```javascript
import { Motion } from 'motion-v'
import { entranceAnimations, springPresets, iconAnimations } from '@/composables/useMotion'
```

### Spring Presets

```javascript
// src/composables/useMotion.js
export const springPresets = {
  // Default - natural, balanced
  default: { type: 'spring', stiffness: 400, damping: 30, mass: 0.8 },
  
  // Snappy - quick responsive
  snappy: { type: 'spring', stiffness: 500, damping: 35, mass: 0.6 },
  
  // Gentle - subtle, soft
  gentle: { type: 'spring', stiffness: 300, damping: 25, mass: 1 },
}
```

### Page Entrance Animations

```vue
<script setup>
import { Motion } from 'motion-v'
import { entranceAnimations } from '@/composables/useMotion'
</script>

<template>
  <Motion v-bind="entranceAnimations.fadeUp" class="page-container">
    <!-- Page content -->
  </Motion>
</template>
```

### Staggered List Items

```vue
<Motion
  v-for="(item, index) in items"
  :key="item.id"
  :initial="{ opacity: 0, y: 15 }"
  :animate="{ opacity: 1, y: 0 }"
  :transition="{ duration: 0.25, delay: index * 0.05, ease: 'easeOut' }"
>
  <ItemCard :item="item" />
</Motion>
```

### Modal Animations

```vue
<!-- Backdrop -->
<Motion
  :initial="{ opacity: 0 }"
  :animate="{ opacity: 1 }"
  :exit="{ opacity: 0 }"
  :transition="{ duration: 0.2, ease: 'easeOut' }"
  class="modal-backdrop"
>

<!-- Modal Container -->
<Motion
  :initial="{ opacity: 0, scale: 0.95, y: 20 }"
  :animate="{ opacity: 1, scale: 1, y: 0 }"
  :transition="{ type: 'spring', stiffness: 500, damping: 40, mass: 0.8 }"
  class="modal-container"
>
```

### CSS-Only Interactions

```css
/* Active press feedback (iOS-like) */
.active-scale:active {
  transform: scale(0.97);
}

/* Hover states */
.btn:hover {
  background-color: theme('colors.primary.light');
}

/* Focus states */
.input:focus {
  ring: 4px;
  ring-color: theme('colors.primary.100');
}
```

---

## 🧩 Component Development

### Base Component Structure

```vue
<script setup>
import { ref, computed } from 'vue'
import { Motion } from 'motion-v'

// Props
const props = defineProps({
  title: {
    type: String,
    required: true
  },
  variant: {
    type: String,
    default: 'primary'
  }
})

// Emits
const emit = defineEmits(['click', 'update'])

// State
const isActive = ref(false)

// Computed
const buttonClasses = computed(() => ({
  'btn-primary': props.variant === 'primary',
  'btn-secondary': props.variant === 'secondary',
}))

// Methods
const handleClick = () => {
  emit('click')
}
</script>

<template>
  <Motion
    :initial="{ opacity: 0 }"
    :animate="{ opacity: 1 }"
    class="component-wrapper"
  >
    <button 
      :class="buttonClasses"
      @click="handleClick"
      class="active-scale"
    >
      {{ title }}
    </button>
  </Motion>
</template>
```

### Composables Pattern

```javascript
// src/composables/useToast.js
import { ref } from 'vue'

const toasts = ref([])

export function useToast() {
  const add = (message, type = 'info') => {
    const id = Date.now()
    toasts.value.push({ id, message, type })
    
    setTimeout(() => {
      remove(id)
    }, 3000)
  }
  
  const remove = (id) => {
    const index = toasts.value.findIndex(t => t.id === id)
    if (index > -1) {
      toasts.value.splice(index, 1)
    }
  }
  
  return {
    toasts,
    success: (msg) => add(msg, 'success'),
    error: (msg) => add(msg, 'error'),
    info: (msg) => add(msg, 'info'),
    warning: (msg) => add(msg, 'warning'),
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
    await api.post('/api/users', formData.value)
    toast.success('User berhasil dibuat')
  } catch (error) {
    toast.error('Gagal membuat user')
  }
}
</script>
```

---

## 🌐 API Integration

### Axios Setup

```javascript
// src/utils/axios.js
import axios from 'axios'
import { useToast } from '@/composables/useToast'

const api = axios.create({
  baseURL: import.meta.env.VITE_API_URL || 'http://localhost:8080',
  headers: {
    'Content-Type': 'application/json'
  }
})

// Request interceptor
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

// Response interceptor
api.interceptors.response.use(
  (response) => response,
  (error) => {
    const toast = useToast()
    
    if (error.response?.status === 401) {
      toast.error('Session expired. Please login again')
      window.location.href = '/auth/login'
    } else if (error.response?.data?.error) {
      toast.error(error.response.data.error)
    } else {
      toast.error('An error occurred')
    }
    
    return Promise.reject(error)
  }
)

export default api
```

### Making API Calls

```vue
<script setup>
import { ref, onMounted } from 'vue'
import api from '@/utils/axios'

const users = ref([])
const loading = ref(false)
const errors = ref({})

const fetchUsers = async () => {
  loading.value = true
  
  try {
    const response = await api.get('/api/users')
    users.value = response.data.data.users
  } catch (error) {
    // Error handled by interceptor
    console.error('Failed to fetch users:', error)
  } finally {
    loading.value = false
  }
}

const createUser = async (userData) => {
  loading.value = true
  errors.value = {}
  
  try {
    const response = await api.post('/api/users', userData)
    toast.success('User created successfully')
    return response.data.data
  } catch (error) {
    if (error.response?.status === 400) {
      errors.value = error.response.data.details || {}
    }
    throw error
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchUsers()
})
</script>
```

---

## 📱 Mobile-First Development

### Responsive Design Principles

1. **Start with mobile layout** (320px+)
2. **Add breakpoints gradually**:
   - `sm`: 640px
   - `md`: 768px
   - `lg`: 1024px
   - `xl`: 1280px

3. **Touch-friendly targets** (min 44x44px)
4. **Test on real devices**

### Responsive Component Example

```vue
<template>
  <div class="container">
    <!-- Mobile: Stack vertically -->
    <!-- Desktop: Side by side -->
    <div class="flex flex-col md:flex-row gap-4">
      <div class="w-full md:w-1/2">
        <!-- Left content -->
      </div>
      <div class="w-full md:w-1/2">
        <!-- Right content -->
      </div>
    </div>
  </div>
</template>
```

### Haptic Feedback (Mobile)

```javascript
// src/utils/haptics.js
export const triggerHaptic = (pattern = 'light') => {
  if ('vibrate' in navigator) {
    const patterns = {
      light: [10],
      medium: [20],
      heavy: [30],
      success: [10, 50, 10],
      error: [30, 100, 30],
    }
    
    navigator.vibrate(patterns[pattern] || patterns.light)
  }
}
```

**Usage:**
```vue
<script setup>
import { triggerHaptic } from '@/utils/haptics'

const handleClick = () => {
  triggerHaptic('light')
  // ... handle click
}
</script>
```

---

## ✅ Best Practices

### 1. Component Organization

- **Keep components small** (<200 lines)
- **Single responsibility** per component
- **Reusable components** in `components/base/`
- **Feature-specific** in `components/features/`

### 2. State Management

- **Use composables** untuk shared state
- **Props down, events up** pattern
- **Avoid prop drilling** (use provide/inject if needed)

### 3. Performance

- **Lazy load** components dengan `defineAsyncComponent`
- **Use `v-show` vs `v-if`** appropriately
- **Optimize images** (WebP format, lazy loading)
- **Minimize bundle size** (code splitting)

### 4. Naming Conventions

```javascript
// Components: PascalCase
UserCard.vue
ProfileModal.vue

// Composables: camelCase with 'use' prefix
useAuth.js
useToast.js

// Utils: camelCase
formatDate.js
validateEmail.js
```

---

## 🧪 Testing

### Component Testing

```javascript
// UserCard.test.js
import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import UserCard from '@/components/UserCard.vue'

describe('UserCard', () => {
  it('renders user name', () => {
    const wrapper = mount(UserCard, {
      props: {
        user: {
          id: 1,
          full_name: 'John Doe'
        }
      }
    })
    
    expect(wrapper.text()).toContain('John Doe')
  })
  
  it('emits click event', async () => {
    const wrapper = mount(UserCard, {
      props: { user: { id: 1, full_name: 'John' } }
    })
    
    await wrapper.find('button').trigger('click')
    expect(wrapper.emitted('click')).toBeTruthy()
  })
})
```

---

## 📚 Additional Resources

### Documentation
- [Vue 3 Docs](https://vuejs.org/)
- [Tailwind CSS](https://tailwindcss.com/)
- [Motion-V](https://motion-v.com/)
- [Vite](https://vitejs.dev/)

### Related Guides
- [Backend Guide](./backend-guide.md)
- [API Documentation](../04-api-reference/README.md)
- [Design Standards](.cursor/rules/) - Apple-inspired design

---

**Last Updated:** 28 Desember 2025  
**Status:** ✅ Production Ready
