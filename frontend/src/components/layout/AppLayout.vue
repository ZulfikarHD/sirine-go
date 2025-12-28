<template>
  <div class="app-layout">
    <!-- Sidebar Navigation -->
    <Sidebar :is-open="isSidebarOpen" @close="closeSidebar" />

    <!-- Main Content Area -->
    <div class="app-main">
      <!-- Top Navbar -->
      <Navbar @toggle-sidebar="toggleSidebar" />
      
      <!-- Page Content dengan route transition -->
      <main class="app-content custom-scrollbar">
        <div class="app-content-inner">
          <Motion
            :key="$route.path"
            v-bind="entranceAnimations.fadeUp"
          >
            <slot />
          </Motion>
        </div>
      </main>
    </div>
  </div>
</template>

<script setup>
/**
 * AppLayout Component
 * Layout utama aplikasi yang menggabungkan Sidebar, Navbar,
 * dan content area dengan route transition animations
 */
import { ref, watch, onMounted, onUnmounted } from 'vue'
import { useRoute } from 'vue-router'
import { Motion } from 'motion-v'
import { entranceAnimations } from '@/composables/useMotion'
import Navbar from './Navbar.vue'
import Sidebar from './Sidebar.vue'

const route = useRoute()
const isSidebarOpen = ref(false)

/**
 * Toggle sidebar visibility (mobile)
 */
const toggleSidebar = () => {
  isSidebarOpen.value = !isSidebarOpen.value
}

/**
 * Close sidebar dengan optional haptic feedback
 */
const closeSidebar = () => {
  isSidebarOpen.value = false
  // Light haptic feedback on close (iOS-like)
  if ('vibrate' in navigator) {
    navigator.vibrate(5)
  }
}

/**
 * Auto-close sidebar on route change (mobile)
 */
watch(() => route.path, () => {
  if (isSidebarOpen.value && window.innerWidth < 768) {
    closeSidebar()
  }
})

/**
 * Handle escape key untuk close sidebar
 */
const handleEscape = (e) => {
  if (e.key === 'Escape' && isSidebarOpen.value) {
    closeSidebar()
  }
}

/**
 * Handle resize - close sidebar when transitioning to desktop
 */
const handleResize = () => {
  if (window.innerWidth >= 768 && isSidebarOpen.value) {
    isSidebarOpen.value = false
  }
}

onMounted(() => {
  document.addEventListener('keydown', handleEscape)
  window.addEventListener('resize', handleResize)
})

onUnmounted(() => {
  document.removeEventListener('keydown', handleEscape)
  window.removeEventListener('resize', handleResize)
})
</script>
