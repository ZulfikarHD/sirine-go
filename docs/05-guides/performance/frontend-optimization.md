# ⚡ Frontend Performance Optimization - Vue 3

Panduan optimasi performa frontend Vue 3 untuk aplikasi Sirine Go dengan fokus pada mobile-first experience.

## 📋 Daftar Isi

1. [Bundle Size Optimization](#bundle-size-optimization)
2. [Code Splitting & Lazy Loading](#code-splitting--lazy-loading)
3. [Component Optimization](#component-optimization)
4. [Animation Performance](#animation-performance)
5. [Asset Optimization](#asset-optimization)
6. [Network Optimization](#network-optimization)
7. [Runtime Performance](#runtime-performance)

---

## 📦 Bundle Size Optimization

### **Analyze Bundle Size**

```bash
# Build with analysis
cd frontend
yarn build --report

# Opens browser with bundle visualization
```

**Target:** Total bundle < 500KB (gzipped)

### **Tree Shaking Configuration**

```javascript
// File: frontend/vite.config.js

export default defineConfig({
  build: {
    rollupOptions: {
      output: {
        manualChunks: {
          'vendor': ['vue', 'vue-router', 'axios'],
          'motion': ['motion-v'],
          'icons': ['lucide-vue-next']
        }
      }
    }
  }
})
```

### **Remove Unused Dependencies**

```bash
# Analyze which dependencies are actually used
npx depcheck

# Remove unused packages
yarn remove <package-name>
```

### **Import Only What You Need**

```javascript
// ❌ BAD: Import entire library
import _ from 'lodash'
import { Icon } from 'lucide-vue-next'

// ✅ GOOD: Import specific functions
import debounce from 'lodash-es/debounce'
import { UserIcon, BellIcon } from 'lucide-vue-next'
```

---

## 🔄 Code Splitting & Lazy Loading

### **Route-Based Code Splitting**

```javascript
// File: frontend/src/router/index.js

const router = createRouter({
  routes: [
    {
      path: '/dashboard',
      // ✅ Lazy load (creates separate chunk)
      component: () => import('@/views/Dashboard.vue')
    },
    {
      path: '/users',
      // ✅ With named chunk
      component: () => import(/* webpackChunkName: "users" */ '@/views/UserManagement.vue')
    },
    {
      path: '/admin',
      // ✅ Prefetch for likely navigation
      component: () => import(/* webpackPrefetch: true */ '@/views/AdminPanel.vue')
    }
  ]
})
```

### **Component-Level Lazy Loading**

```vue
<script setup>
// ✅ Lazy load heavy components
const UserTable = defineAsyncComponent(() => 
  import('@/components/UserTable.vue')
)

const ChartWidget = defineAsyncComponent({
  loader: () => import('@/components/ChartWidget.vue'),
  loadingComponent: LoadingSpinner,
  delay: 200, // Show loading after 200ms
  timeout: 10000
})
</script>

<template>
  <Suspense>
    <template #default>
      <UserTable />
    </template>
    <template #fallback>
      <LoadingSpinner />
    </template>
  </Suspense>
</template>
```

### **Conditional Loading**

```vue
<script setup>
import { ref, defineAsyncComponent } from 'vue'

const showChart = ref(false)

// Only load when needed
const AnalyticsChart = showChart.value 
  ? defineAsyncComponent(() => import('@/components/AnalyticsChart.vue'))
  : null
</script>
```

---

## ⚙️ Component Optimization

### **Use v-once for Static Content**

```vue
<template>
  <!-- Static header - render once -->
  <div v-once class="app-header">
    <img src="/logo.svg" alt="Logo">
    <h1>Sirine Go</h1>
  </div>
</template>
```

### **v-memo for Expensive Lists**

```vue
<template>
  <!-- Only re-render if item.id or item.updated_at changes -->
  <div 
    v-for="item in items" 
    :key="item.id"
    v-memo="[item.id, item.updated_at]"
  >
    <ExpensiveComponent :data="item" />
  </div>
</template>
```

### **Virtual Scrolling for Long Lists**

```bash
yarn add vue-virtual-scroller
```

```vue
<script setup>
import { RecycleScroller } from 'vue-virtual-scroller'
import 'vue-virtual-scroller/dist/vue-virtual-scroller.css'

const items = ref([]) // 1000+ items
</script>

<template>
  <!-- Only renders visible items -->
  <RecycleScroller
    :items="items"
    :item-size="72"
    key-field="id"
  >
    <template #default="{ item }">
      <UserCard :user="item" />
    </template>
  </RecycleScroller>
</template>
```

### **Computed Properties for Expensive Operations**

```vue
<script setup>
import { computed } from 'vue'

const users = ref([])

// ❌ BAD: Recalculates on every render
const filteredUsers = () => {
  return users.value.filter(u => u.active).sort((a, b) => a.name.localeCompare(b.name))
}

// ✅ GOOD: Cached until users changes
const filteredUsers = computed(() => {
  return users.value
    .filter(u => u.active)
    .sort((a, b) => a.name.localeCompare(b.name))
})
</script>
```

### **Debounce User Input**

```vue
<script setup>
import { ref, watch } from 'vue'
import { debounce } from 'lodash-es'

const searchQuery = ref('')

// ✅ Debounce API calls
const debouncedSearch = debounce((query) => {
  searchUsers(query)
}, 300)

watch(searchQuery, (newQuery) => {
  debouncedSearch(newQuery)
})
</script>
```

---

## 🎨 Animation Performance

### **Use Motion-V Properly**

```vue
<script setup>
import { Motion } from 'motion-v'
import { springPresets } from '@/composables/useMotion'

// ✅ GOOD: Animate transform & opacity (GPU-accelerated)
const goodAnimation = {
  initial: { opacity: 0, y: 20 },
  animate: { opacity: 1, y: 0 },
  transition: springPresets.default
}

// ❌ BAD: Animate width, height, margin (triggers layout)
const badAnimation = {
  initial: { width: 0, height: 0 },
  animate: { width: 300, height: 200 },
}
</script>

<template>
  <Motion v-bind="goodAnimation">
    <div class="card">...</div>
  </Motion>
</template>
```

### **Limit Concurrent Animations**

```vue
<script setup>
// ✅ Stagger animations with small delays
const items = ref([...])
</script>

<template>
  <Motion
    v-for="(item, index) in items"
    :key="item.id"
    :initial="{ opacity: 0, y: 15 }"
    :animate="{ opacity: 1, y: 0 }"
    :transition="{ 
      duration: 0.25, 
      delay: index * 0.05,  // Max 0.05s delay
      ease: 'easeOut'
    }"
  >
    <UserCard :user="item" />
  </Motion>
</template>
```

### **Use will-change Sparingly**

```vue
<style scoped>
/* ✅ GOOD: Only on elements that will animate */
.modal-enter-active,
.modal-leave-active {
  will-change: transform, opacity;
}

.modal-enter-to,
.modal-leave-from {
  will-change: auto; /* Remove after animation */
}

/* ❌ BAD: On all elements */
.card {
  will-change: transform; /* Wastes memory! */
}
</style>
```

### **Reduce Backdrop Filters**

```vue
<style scoped>
/* ❌ BAD: Backdrop blur everywhere */
.glass-card {
  backdrop-filter: blur(8px);
}

/* ✅ GOOD: Only on navbar (1 per viewport max) */
.navbar {
  backdrop-filter: blur(8px);
}

.card {
  background: rgba(255, 255, 255, 0.95);
  /* No backdrop-filter */
}
</style>
```

---

## 🖼️ Asset Optimization

### **Image Optimization**

```bash
# Install image optimizer
yarn add -D vite-plugin-imagemin

# Convert to WebP for smaller size
```

```javascript
// vite.config.js
import viteImagemin from 'vite-plugin-imagemin'

export default defineConfig({
  plugins: [
    viteImagemin({
      gifsicle: { optimizationLevel: 7 },
      mozjpeg: { quality: 80 },
      pngquant: { quality: [0.8, 0.9] },
      svgo: {
        plugins: [{ removeViewBox: false }]
      }
    })
  ]
})
```

### **Use Modern Image Formats**

```vue
<template>
  <!-- ✅ Use WebP with fallback -->
  <picture>
    <source srcset="/images/banner.webp" type="image/webp">
    <source srcset="/images/banner.jpg" type="image/jpeg">
    <img src="/images/banner.jpg" alt="Banner" loading="lazy">
  </picture>
</template>
```

### **Lazy Load Images**

```vue
<template>
  <!-- ✅ Native lazy loading -->
  <img src="/avatar.jpg" alt="Avatar" loading="lazy">
  
  <!-- ✅ Blur placeholder -->
  <img 
    src="/avatar-placeholder.jpg" 
    data-src="/avatar.jpg"
    class="lazyload"
    alt="Avatar"
  >
</template>
```

### **Icon Strategy**

```vue
<script setup>
// ✅ Import only needed icons
import { UserIcon, BellIcon, SettingsIcon } from 'lucide-vue-next'

// ❌ Don't import entire icon library
// import * as Icons from 'lucide-vue-next'
</script>
```

---

## 🌐 Network Optimization

### **API Request Caching**

```javascript
// File: frontend/src/services/apiCache.js

const cache = new Map()
const CACHE_DURATION = 5 * 60 * 1000 // 5 minutes

export async function fetchWithCache(url, options = {}) {
  const cacheKey = `${url}_${JSON.stringify(options)}`
  const cached = cache.get(cacheKey)
  
  if (cached && Date.now() - cached.timestamp < CACHE_DURATION) {
    return cached.data
  }
  
  const response = await fetch(url, options)
  const data = await response.json()
  
  cache.set(cacheKey, {
    data,
    timestamp: Date.now()
  })
  
  return data
}
```

### **Request Deduplication**

```javascript
// Prevent duplicate concurrent requests
const pendingRequests = new Map()

export async function fetchWithDedup(url) {
  if (pendingRequests.has(url)) {
    return pendingRequests.get(url)
  }
  
  const promise = fetch(url).then(r => r.json())
  pendingRequests.set(url, promise)
  
  try {
    const data = await promise
    return data
  } finally {
    pendingRequests.delete(url)
  }
}
```

### **Preload Critical Data**

```vue
<script setup>
import { onMounted } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()

onMounted(() => {
  // ✅ Preload likely next pages
  router.beforeEach((to, from, next) => {
    if (to.path === '/dashboard') {
      // Prefetch dashboard data
      fetchDashboardData()
    }
    next()
  })
})
</script>
```

---

## ⚡ Runtime Performance

### **Avoid Memory Leaks**

```vue
<script setup>
import { onMounted, onUnmounted } from 'vue'

let interval

onMounted(() => {
  interval = setInterval(() => {
    updateNotifications()
  }, 30000)
})

// ✅ IMPORTANT: Clean up
onUnmounted(() => {
  if (interval) {
    clearInterval(interval)
  }
})
</script>
```

### **Use Web Workers for Heavy Computations**

```javascript
// File: frontend/src/workers/dataProcessor.worker.js

self.addEventListener('message', (e) => {
  const { data } = e.data
  
  // Heavy computation
  const result = processLargeDataset(data)
  
  self.postMessage({ result })
})
```

```vue
<script setup>
const worker = new Worker('/workers/dataProcessor.worker.js')

worker.postMessage({ data: largeDataset })

worker.onmessage = (e) => {
  const { result } = e.data
  displayResult(result)
}
</script>
```

### **Service Worker for Offline Support**

```javascript
// File: frontend/public/sw.js

const CACHE_NAME = 'sirine-go-v1'
const urlsToCache = [
  '/',
  '/css/app.css',
  '/js/app.js',
  '/images/logo.svg'
]

self.addEventListener('install', (event) => {
  event.waitUntil(
    caches.open(CACHE_NAME)
      .then((cache) => cache.addAll(urlsToCache))
  )
})

self.addEventListener('fetch', (event) => {
  event.respondWith(
    caches.match(event.request)
      .then((response) => response || fetch(event.request))
  )
})
```

---

## ✅ Performance Checklist

### Build Configuration
- [ ] Bundle size < 500KB (gzipped)
- [ ] Route-based code splitting enabled
- [ ] Tree shaking configured
- [ ] Images optimized (WebP)
- [ ] Icons tree-shaken
- [ ] Source maps disabled in production

### Component Optimization
- [ ] Virtual scrolling for long lists (>100 items)
- [ ] v-memo used for expensive list items
- [ ] Computed properties for expensive calculations
- [ ] Debounced search inputs
- [ ] Lazy loaded heavy components

### Animation Performance
- [ ] Only animate transform & opacity
- [ ] Maximum 1 backdrop-filter per viewport
- [ ] Stagger delays ≤ 0.05s
- [ ] will-change used sparingly
- [ ] Animation durations 0.2-0.3s

### Network Optimization
- [ ] API response caching implemented
- [ ] Request deduplication enabled
- [ ] Critical resources preloaded
- [ ] Images lazy loaded
- [ ] Service worker registered

---

## 📊 Target Performance Metrics

### Lighthouse Scores (Mobile)
- Performance: > 90
- Accessibility: > 95
- Best Practices: > 95
- SEO: > 90

### Core Web Vitals
- **LCP (Largest Contentful Paint):** < 2.5s
- **FID (First Input Delay):** < 100ms
- **CLS (Cumulative Layout Shift):** < 0.1

### Bundle Sizes
- Initial bundle: < 200KB (gzipped)
- Total bundle: < 500KB (gzipped)
- Route chunks: < 50KB each (gzipped)

---

## 🔍 Debugging Performance

### Chrome DevTools

```bash
# 1. Open DevTools
# 2. Performance tab
# 3. Record
# 4. Perform actions
# 5. Stop & analyze
```

### Lighthouse Audit

```bash
# Install Lighthouse CI
npm install -g @lhci/cli

# Run audit
lhci autorun --collect.url=http://localhost:3000
```

### Vue DevTools

```bash
# Install browser extension
# Enable Performance tab
# Record component render times
```

---

## 📞 Support

**Developer:** Zulfikar Hidayatullah  
**Phone:** +62 857-1583-8733

## 📖 Related Documentation

- [Frontend Development Guide](../../03-development/frontend/getting-started.md)
- [Animation Guide](../../03-development/frontend/animations.md)
- [Component Guide](../../03-development/frontend/components.md)
- [Backend Optimization](./backend-optimization.md)

---

**Last Updated:** 29 Desember 2025  
**Version:** 1.0.0
