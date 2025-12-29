# 🚀 Frontend Getting Started

Panduan awal untuk memulai frontend development dalam Sirine Go App dengan Vue 3, Vite, dan Motion-V, yaitu setup environment dan basic concepts.

---

## 📋 Prerequisites

### Required Software

Pastikan sudah terinstal:

| Software | Version | Check Command | Purpose |
|----------|---------|---------------|---------|
| **Node.js** | v18+ | `node --version` | JavaScript runtime |
| **Yarn** | Latest | `yarn --version` | Package manager |
| **Git** | Latest | `git --version` | Version control |

### Optional Tools

| Tool | Purpose | Installation |
|------|---------|--------------|
| **Vue DevTools** | Browser extension untuk debugging | [Chrome](https://chrome.google.com/webstore/detail/vuejs-devtools/nhdogjmejiglipccpnnnanhbledajbpd) \| [Firefox](https://addons.mozilla.org/en-US/firefox/addon/vue-js-devtools/) |
| **VS Code** | Recommended code editor | [Download](https://code.visualstudio.com/) |
| **Volar** | Vue 3 language support untuk VS Code | [Extension](https://marketplace.visualstudio.com/items?itemName=Vue.volar) |

### Knowledge Prerequisites

Disarankan memiliki pengetahuan dasar tentang:
- JavaScript (ES6+)
- Vue 3 Composition API
- HTML/CSS fundamentals
- RESTful API concepts

---

## 🔧 Setup Development Environment

### Step 1: Clone Repository

```bash
# Clone project
git clone <repository-url>
cd sirine-go

# Navigate ke frontend folder
cd frontend
```

### Step 2: Install Dependencies

```bash
# Install packages dengan Yarn
yarn install

# Verify installation
yarn --version
```

### Step 3: Configure Environment

```bash
# Copy environment template
cp .env.example .env

# Edit dengan API URL backend
nano .env
```

**Konfigurasi minimal `.env`:**
```env
# API Backend URL
VITE_API_BASE_URL=http://localhost:8080

# App Settings
VITE_APP_NAME=Sirine Go
VITE_TIMEZONE=Asia/Jakarta

# Feature Flags
VITE_ENABLE_GAMIFICATION=true
VITE_ENABLE_NOTIFICATIONS=true
VITE_POLLING_INTERVAL=30000
```

### Step 4: Start Development Server

```bash
# Start dev server dengan hot reload
yarn dev
```

**Output yang diharapkan:**
```
  VITE v5.0.0  ready in 500 ms

  ➜  Local:   http://localhost:5173/
  ➜  Network: use --host to expose
  ➜  press h + enter to show help
```

**Browser:** Buka http://localhost:5173

---

## 🎯 Development Commands

### Basic Commands

```bash
# Start development server
yarn dev

# Build untuk production
yarn build

# Preview production build
yarn preview

# Run linter
yarn lint

# Fix linting issues
yarn lint --fix
```

### Testing Commands

```bash
# Run tests
yarn test

# Run tests dengan watch mode
yarn test:watch

# Generate coverage report
yarn test:coverage
```

---

## 📁 Project Structure

```
frontend/
├── src/
│   ├── assets/                # Static assets
│   │   ├── images/            # Image files
│   │   └── fonts/             # Font files
│   │
│   ├── components/            # Vue components
│   │   ├── base/              # Base reusable components
│   │   │   ├── BaseButton.vue
│   │   │   ├── BaseInput.vue
│   │   │   └── BaseModal.vue
│   │   │
│   │   ├── layout/            # Layout components
│   │   │   ├── Navbar.vue
│   │   │   ├── Sidebar.vue
│   │   │   └── Footer.vue
│   │   │
│   │   └── features/          # Feature-specific components
│   │       ├── UserCard.vue
│   │       ├── NotificationBell.vue
│   │       └── AchievementBadge.vue
│   │
│   ├── composables/           # Composition API composables
│   │   ├── useAuth.js         # Authentication state & methods
│   │   ├── useToast.js        # Toast notifications
│   │   ├── useMotion.js       # Animation presets
│   │   ├── useModal.js        # Modal management
│   │   └── useApi.js          # API calls wrapper
│   │
│   ├── pages/                 # Inertia page components
│   │   ├── Auth/              # Authentication pages
│   │   │   ├── Login.vue
│   │   │   └── ForgotPassword.vue
│   │   │
│   │   ├── Dashboard/         # Dashboard pages
│   │   │   └── Index.vue
│   │   │
│   │   ├── Users/             # User management pages
│   │   │   ├── Index.vue
│   │   │   ├── Create.vue
│   │   │   └── Edit.vue
│   │   │
│   │   └── Profile/           # User profile pages
│   │       ├── Index.vue
│   │       └── Edit.vue
│   │
│   ├── utils/                 # Utility functions
│   │   ├── axios.js           # Axios configuration
│   │   ├── formatters.js      # Data formatting functions
│   │   ├── validators.js      # Form validation helpers
│   │   └── haptics.js         # Haptic feedback utilities
│   │
│   ├── style.css              # Global styles (minimal)
│   └── app.js                 # Main app entry point
│
├── public/                    # Static files (served as-is)
│   └── favicon.ico
│
├── .env                       # Environment variables (git-ignored)
├── .env.example               # Environment template
├── index.html                 # HTML entry point
├── tailwind.config.js         # Tailwind CSS configuration
├── vite.config.js             # Vite configuration
├── package.json               # Dependencies & scripts
└── yarn.lock                  # Yarn lock file
```

### Folder Structure Explanation

**`components/base/`** - Reusable UI components:
- Generic components yang digunakan di multiple pages
- Tidak memiliki business logic
- Fully customizable via props

**`components/features/`** - Feature-specific components:
- Components dengan business logic
- Tied to specific features
- May use composables untuk state management

**`composables/`** - Shared reactive state dan logic:
- Follow `use*` naming convention
- Return reactive refs dan methods
- Can be used across multiple components

**`pages/`** - Inertia page components:
- Top-level components untuk routes
- Receive props dari backend
- Compose smaller components

**`utils/`** - Pure utility functions:
- No Vue dependencies
- Stateless helper functions
- Easy to unit test

---

## 📚 Vue 3 Composition API Basics

### Component Structure

```vue
<script setup>
// 1. Imports
import { ref, computed, onMounted } from 'vue'
import { Motion } from 'motion-v'
import { useToast } from '@/composables/useToast'

// 2. Props definition
const props = defineProps({
  title: {
    type: String,
    required: true
  },
  count: {
    type: Number,
    default: 0
  }
})

// 3. Emits definition
const emit = defineEmits(['update', 'delete'])

// 4. Reactive state
const isActive = ref(false)
const items = ref([])

// 5. Computed properties
const activeCount = computed(() => {
  return items.value.filter(item => item.active).length
})

// 6. Methods
const handleClick = () => {
  emit('update', items.value)
}

// 7. Lifecycle hooks
onMounted(() => {
  console.log('Component mounted')
})
</script>

<template>
  <Motion
    :initial="{ opacity: 0 }"
    :animate="{ opacity: 1 }"
  >
    <div class="component">
      <h2>{{ title }}</h2>
      <p>Active: {{ activeCount }} / {{ count }}</p>
      <button @click="handleClick">Update</button>
    </div>
  </Motion>
</template>
```

### Reactive State

```javascript
import { ref, reactive } from 'vue'

// Ref untuk primitive values
const count = ref(0)
const message = ref('Hello')

// Access value dengan .value
console.log(count.value)  // 0
count.value++              // 1

// Reactive untuk objects
const state = reactive({
  user: null,
  isLoading: false,
  errors: []
})

// Access directly (no .value)
state.isLoading = true
state.user = { name: 'John' }
```

### Computed Properties

```javascript
import { ref, computed } from 'vue'

const firstName = ref('John')
const lastName = ref('Doe')

// Read-only computed
const fullName = computed(() => {
  return `${firstName.value} ${lastName.value}`
})

// Writable computed
const fullNameWritable = computed({
  get: () => `${firstName.value} ${lastName.value}`,
  set: (value) => {
    const parts = value.split(' ')
    firstName.value = parts[0]
    lastName.value = parts[1]
  }
})
```

### Watchers

```javascript
import { ref, watch } from 'vue'

const count = ref(0)
const message = ref('')

// Watch single value
watch(count, (newValue, oldValue) => {
  console.log(`Count changed from ${oldValue} to ${newValue}`)
})

// Watch multiple values
watch([count, message], ([newCount, newMessage]) => {
  console.log('Either count or message changed')
})

// Watch with options
watch(
  count,
  (newValue) => {
    // Do something
  },
  {
    immediate: true,  // Run immediately
    deep: true        // Deep watch for objects
  }
)
```

---

## 🔍 Testing Your Setup

### 1. Check Dev Server

```bash
# Server harus running di http://localhost:5173
curl http://localhost:5173
```

### 2. Test Hot Reload

Edit `src/pages/Dashboard/Index.vue` dan save. Browser harus auto-refresh.

### 3. Test API Connection

Pastikan backend running, lalu coba login dari frontend.

---

## ⚠️ Troubleshooting

### Port Already in Use

```bash
# Error: Port 5173 is in use
# Solution: Kill process atau ubah port di vite.config.js
lsof -ti:5173 | xargs kill -9
```

### Module Not Found

```bash
# Error: Cannot find module '@/components/...'
# Solution: Re-install dependencies
rm -rf node_modules yarn.lock
yarn install
```

### Vite Build Fails

```bash
# Error: Build failed with errors
# Solution: Clear cache dan rebuild
rm -rf node_modules/.vite
yarn build
```

### Hot Reload Not Working

```bash
# Solution: Restart dev server
# Ctrl+C to stop
yarn dev
```

---

## 📚 Next Steps

Setelah setup berhasil, lanjutkan dengan:

1. **[Component Guide](./components.md)** - Belajar membuat Vue components
2. **[Animation Guide](./animations.md)** - Implementasi Motion-V animations
3. **[Styling Guide](./styling.md)** - Apply design system
4. **[API Reference](../../04-api-reference/README.md)** - Explore available endpoints

---

## 📖 Learning Resources

### Official Documentation
- [Vue 3 Docs](https://vuejs.org/)
- [Vite Guide](https://vitejs.dev/guide/)
- [Motion-V Docs](https://motion-v.com/)
- [Tailwind CSS](https://tailwindcss.com/)

### Vue 3 Tutorials
- [Vue Mastery](https://www.vuemastery.com/)
- [Vue School](https://vueschool.io/)

---

**Last Updated:** 28 Desember 2025  
**Status:** ✅ Production Ready
