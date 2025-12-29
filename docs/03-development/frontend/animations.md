# 🎬 Animation System (Motion-V)

Panduan penggunaan Motion-V untuk iOS-like smooth animations dalam Sirine Go App, yaitu physics-based animations yang memberikan premium feel.

---

## 📋 Overview

Animation System dalam Sirine Go App menggunakan **Motion-V exclusively**, bukan CSS animations, untuk consistent dan smooth transitions dengan iOS-like feel.

### Core Principle

**SEMUA animasi menggunakan Motion-V**  
CSS hanya untuk: hover, focus, active states, dan transform feedback.

### Why Motion-V?

**Benefits:**
- Physics-based spring animations untuk natural movement
- Better performance dengan GPU acceleration
- Easier to create complex sequential animations
- Declarative API yang mudah digunakan
- Consistent timing across all devices

---

## 🚀 Quick Start

### Installation & Import

Motion-V sudah included dalam project dependencies.

```vue
<script setup>
// Import Motion component
import { Motion } from 'motion-v'

// Import animation presets dari composable
import { entranceAnimations, springPresets, iconAnimations } from '@/composables/useMotion'
</script>
```

### Basic Usage

```vue
<template>
  <Motion
    :initial="{ opacity: 0, y: 20 }"
    :animate="{ opacity: 1, y: 0 }"
    :transition="{ duration: 0.3, ease: 'easeOut' }"
  >
    <div class="card">
      <!-- Content -->
    </div>
  </Motion>
</template>
```

---

## 🎨 Animation Presets

Project menyediakan presets di `src/composables/useMotion.js` untuk consistency.

### Spring Presets

```javascript
// src/composables/useMotion.js
export const springPresets = {
  // Default - natural, balanced (recommended untuk most cases)
  default: {
    type: 'spring',
    stiffness: 400,
    damping: 30,
    mass: 0.8
  },
  
  // Snappy - quick responsive (untuk buttons, small elements)
  snappy: {
    type: 'spring',
    stiffness: 500,
    damping: 35,
    mass: 0.6
  },
  
  // Gentle - subtle, soft (untuk large modals, page transitions)
  gentle: {
    type: 'spring',
    stiffness: 300,
    damping: 25,
    mass: 1
  }
}
```

### Entrance Animations

```javascript
export const entranceAnimations = {
  // Fade up - most common
  fadeUp: {
    initial: { opacity: 0, y: 20 },
    animate: { opacity: 1, y: 0 },
    transition: { duration: 0.3, ease: 'easeOut' }
  },
  
  // Fade with scale - untuk emphasis
  fadeScale: {
    initial: { opacity: 0, scale: 0.95 },
    animate: { opacity: 1, scale: 1 },
    transition: springPresets.default
  },
  
  // Simple fade
  fade: {
    initial: { opacity: 0 },
    animate: { opacity: 1 },
    transition: { duration: 0.2, ease: 'easeOut' }
  },
  
  // Slide from left
  slideLeft: {
    initial: { opacity: 0, x: -20 },
    animate: { opacity: 1, x: 0 },
    transition: { duration: 0.25, ease: 'easeOut' }
  },
  
  // Slide from right
  slideRight: {
    initial: { opacity: 0, x: 20 },
    animate: { opacity: 1, x: 0 },
    transition: { duration: 0.25, ease: 'easeOut' }
  }
}
```

### Icon Animations

```javascript
export const iconAnimations = {
  // Pop in - untuk icons di modals
  popIn: {
    initial: { scale: 0 },
    animate: { scale: 1 },
    transition: {
      type: 'spring',
      stiffness: 600,
      damping: 30,
      mass: 0.8
    }
  },
  
  // Bounce - untuk success indicators
  bounce: {
    initial: { scale: 0 },
    animate: { scale: [0, 1.2, 1] },
    transition: {
      duration: 0.5,
      times: [0, 0.6, 1],
      ease: 'easeOut'
    }
  }
}
```

---

## 🎯 Common Patterns

### Page/Section Entrance

```vue
<script setup>
import { Motion } from 'motion-v'
import { entranceAnimations } from '@/composables/useMotion'
</script>

<template>
  <Motion v-bind="entranceAnimations.fadeUp" class="page-container">
    <h1>Page Title</h1>
    <p>Page content...</p>
  </Motion>
</template>
```

### Staggered List Items

```vue
<script setup>
import { Motion } from 'motion-v'

const items = ref([...])
</script>

<template>
  <div class="grid gap-4">
    <Motion
      v-for="(item, index) in items"
      :key="item.id"
      :initial="{ opacity: 0, y: 15 }"
      :animate="{ opacity: 1, y: 0 }"
      :transition="{
        duration: 0.25,
        delay: index * 0.05,  // Stagger delay: 50ms per item
        ease: 'easeOut'
      }"
    >
      <ItemCard :item="item" />
    </Motion>
  </div>
</template>
```

**Best Practices:**
- Maksimal delay: `0.05s` (50ms) per item
- Total stagger tidak lebih dari `1s` (max 20 items)

### Modal/Dialog Animations

```vue
<script setup>
import { Motion } from 'motion-v'
import { springPresets } from '@/composables/useMotion'

const props = defineProps({
  show: Boolean
})
</script>

<template>
  <Teleport to="body">
    <div v-if="show" class="fixed inset-0 z-50 flex items-center justify-center">
      <!-- Backdrop fade -->
      <Motion
        :initial="{ opacity: 0 }"
        :animate="{ opacity: 1 }"
        :exit="{ opacity: 0 }"
        :transition="{ duration: 0.2, ease: 'easeOut' }"
        class="fixed inset-0 bg-black/50 backdrop-blur-sm"
      />
      
      <!-- Modal spring animation -->
      <Motion
        :initial="{ opacity: 0, scale: 0.95, y: 20 }"
        :animate="{ opacity: 1, scale: 1, y: 0 }"
        :exit="{ opacity: 0, scale: 0.95, y: 20 }"
        :transition="springPresets.default"
        class="relative bg-white rounded-2xl shadow-2xl"
      >
        <!-- Modal content -->
      </Motion>
    </div>
  </Teleport>
</template>
```

### Button with Icon Pop

```vue
<template>
  <button class="btn active-scale">
    <Motion v-bind="iconAnimations.popIn">
      <CheckIcon class="w-5 h-5" />
    </Motion>
    <span>Success</span>
  </button>
</template>
```

### Collapsible/Accordion

```vue
<script setup>
import { ref } from 'vue'
import { Motion } from 'motion-v'

const isOpen = ref(false)
</script>

<template>
  <div>
    <button @click="isOpen = !isOpen" class="active-scale">
      Toggle Content
    </button>
    
    <Motion
      v-if="isOpen"
      :initial="{ height: 0, opacity: 0 }"
      :animate="{ height: 'auto', opacity: 1 }"
      :exit="{ height: 0, opacity: 0 }"
      :transition="{ duration: 0.3, ease: 'easeOut' }"
      style="overflow: hidden"
    >
      <div class="p-4">
        Collapsible content here...
      </div>
    </Motion>
  </div>
</template>
```

### Loading Skeleton

```vue
<template>
  <Motion
    :animate="{
      opacity: [0.5, 1, 0.5]
    }"
    :transition="{
      duration: 1.5,
      repeat: Infinity,
      ease: 'easeInOut'
    }"
    class="bg-gray-200 rounded h-20"
  />
</template>
```

---

## 💅 CSS-Only Interactions

Some interactions SHOULD use CSS for better performance:

### Active Press Feedback (iOS-like)

```css
/* src/style.css */
.active-scale:active {
  transform: scale(0.97);
  transition: transform 0.1s ease-out;
}
```

**Usage:**
```vue
<template>
  <button class="active-scale">
    Click Me
  </button>
</template>
```

### Hover States

```css
/* Prefer Tailwind utilities */
<button class="hover:bg-indigo-700 transition-colors">
  Hover Me
</button>
```

### Focus States

```css
<input class="focus:ring-4 focus:ring-indigo-100 transition-all" />
```

---

## ⚡ Performance Optimization

### 1. Use `will-change` Sparingly

```vue
<template>
  <!-- Do: Only on animated elements -->
  <Motion
    :initial="{ opacity: 0 }"
    :animate="{ opacity: 1 }"
    style="will-change: transform, opacity"
  >
</template>
```

### 2. Prefer Transform & Opacity

```vue
<!-- Do: GPU-accelerated properties -->
<Motion
  :animate="{ opacity: 1, scale: 1, x: 0, y: 0 }"
>

<!-- Don't: Triggers layout recalculation -->
<Motion
  :animate="{ width: '300px', height: '200px', top: '50px' }"
>
```

### 3. Limit Stagger Count

```vue
<!-- Do: Max 20 items atau limit delay -->
<Motion
  v-for="(item, index) in items.slice(0, 20)"
  :transition="{ delay: index * 0.05 }"
>

<!-- Don't: Stagger 100+ items -->
<Motion
  v-for="(item, index) in allItems"  // Could be 100+ items
  :transition="{ delay: index * 0.05 }"
>
```

### 4. Avoid Nested Animations

```vue
<!-- Do: Animate container only -->
<Motion :animate="{ opacity: 1 }">
  <div>
    <ChildComponent />
  </div>
</Motion>

<!-- Don't: Nested Motion components -->
<Motion :animate="{ opacity: 1 }">
  <Motion :animate="{ scale: 1 }">
    <Motion :animate="{ y: 0 }">
      <!-- Too many layers -->
    </Motion>
  </Motion>
</Motion>
```

---

## 🎨 Animation Timing

### Recommended Durations

| Element Type | Duration | Reason |
|-------------|----------|--------|
| Small elements (buttons, icons) | 0.15-0.2s | Quick feedback |
| Medium elements (cards, items) | 0.25-0.3s | Balanced |
| Large elements (modals, pages) | 0.3-0.4s | Smooth transition |
| Backdrop overlays | 0.2s | Fast background fade |

### Easing Functions

```javascript
// Entrance animations
ease: 'easeOut'  // Starts fast, ends slow

// Exit animations
ease: 'easeIn'   // Starts slow, ends fast

// Both directions
ease: 'easeInOut' // Smooth both sides

// Spring physics (recommended)
type: 'spring'
stiffness: 400
damping: 30
```

---

## 🚫 What NOT to Animate with Motion-V

Use CSS for these instead:

### 1. Spinner Animations

```css
/* Do: Use CSS animation for infinite spin */
.spinner {
  animation: spin 0.6s linear infinite;
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}
```

### 2. Hover States

```vue
<!-- Do: Use Tailwind hover utilities -->
<div class="hover:scale-105 transition-transform">
  Hover me
</div>

<!-- Don't: Use Motion-V for hover -->
<Motion :whileHover="{ scale: 1.05 }">
  Hover me
</Motion>
```

### 3. Persistent Animations

```css
/* Do: Use CSS for persistent effects */
.pulse-dot {
  animation: pulse 2s ease-in-out infinite;
}

@keyframes pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.5; }
}
```

---

## 📱 Mobile Considerations

### Haptic Feedback

Combine animations dengan haptic feedback untuk better UX:

```vue
<script setup>
import { triggerHaptic } from '@/utils/haptics'

const handleClick = () => {
  triggerHaptic('light')  // Vibrate on tap
  // ... handle action
}
</script>

<template>
  <Motion :animate="{ scale: [1, 0.97, 1] }">
    <button @click="handleClick">
      Tap Me
    </button>
  </Motion>
</template>
```

### Reduced Motion

Respect user's motion preferences:

```vue
<script setup>
import { ref, onMounted } from 'vue'

const prefersReducedMotion = ref(false)

onMounted(() => {
  prefersReducedMotion.value = window.matchMedia('(prefers-reduced-motion: reduce)').matches
})

const transition = computed(() => {
  if (prefersReducedMotion.value) {
    return { duration: 0.01 }  // Instant
  }
  return { duration: 0.3, ease: 'easeOut' }
})
</script>

<template>
  <Motion
    :initial="{ opacity: 0 }"
    :animate="{ opacity: 1 }"
    :transition="transition"
  >
    <!-- Content -->
  </Motion>
</template>
```

---

## ✅ Best Practices

### 1. Keep Animations Short

```vue
<!-- Do: 0.2-0.4s duration -->
:transition="{ duration: 0.3 }"

<!-- Don't: Too slow -->
:transition="{ duration: 1.5 }"
```

### 2. Use Presets

```vue
<!-- Do: Use existing presets -->
<Motion v-bind="entranceAnimations.fadeUp">

<!-- Don't: Recreate inline -->
<Motion
  :initial="{ opacity: 0, y: 20 }"
  :animate="{ opacity: 1, y: 0 }"
  :transition="{ duration: 0.3, ease: 'easeOut' }"
>
```

### 3. Animate on Mount, Not on Data Change

```vue
<!-- Do: Animate component mount -->
<Motion :initial="{ opacity: 0 }" :animate="{ opacity: 1 }">
  <div>{{ data }}</div>
</Motion>

<!-- Don't: Re-animate on every data change -->
<Motion :animate="{ opacity: data ? 1 : 0 }">
  <div>{{ data }}</div>
</Motion>
```

### 4. Exit Animations

```vue
<!-- Do: Define exit for smooth removal -->
<Motion
  :initial="{ opacity: 0 }"
  :animate="{ opacity: 1 }"
  :exit="{ opacity: 0 }"
>
```

---

## 📚 Related Documentation

- [Getting Started Guide](./getting-started.md) - Setup environment
- [Component Guide](./components.md) - Component development
- [Styling Guide](./styling.md) - Design system
- [Motion-V Official Docs](https://motion-v.com/)

---

**Last Updated:** 28 Desember 2025  
**Status:** ✅ Production Ready
