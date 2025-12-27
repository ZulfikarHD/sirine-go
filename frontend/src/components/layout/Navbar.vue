<template>
  <nav class="glass-navbar sticky top-0 z-50 backdrop-blur-xl bg-white/80 border-b border-gray-200/50">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <div class="flex justify-between items-center h-16">
        <!-- Logo & Title -->
        <div class="flex items-center space-x-3">
          <div class="w-10 h-10 rounded-xl bg-gradient-to-br from-indigo-500 to-fuchsia-600 flex items-center justify-center shadow-lg">
            <svg class="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z" />
            </svg>
          </div>
          <div>
            <h1 class="text-xl font-bold text-gray-900">Sirine Go</h1>
            <p class="text-xs text-gray-500">Sistem Produksi Pita Cukai</p>
          </div>
        </div>

        <!-- User Menu -->
        <div class="relative">
          <button
            @click="toggleMenu"
            class="flex items-center space-x-3 px-3 py-2 rounded-xl hover:bg-gray-100 transition-colors active-scale"
          >
            <!-- Avatar -->
            <div class="w-10 h-10 rounded-full bg-gradient-to-br from-indigo-500 to-fuchsia-600 flex items-center justify-center text-white font-semibold shadow-lg">
              {{ userInitial }}
            </div>
            <!-- User Info (hidden on mobile) -->
            <div class="hidden md:block text-left">
              <p class="text-sm font-semibold text-gray-900">{{ user?.full_name }}</p>
              <p class="text-xs text-gray-500">{{ user?.role }}</p>
            </div>
            <!-- Dropdown Icon -->
            <svg 
              class="w-5 h-5 text-gray-500 transition-transform duration-200" 
              :class="{ 'rotate-180': isMenuOpen }"
              fill="none" 
              stroke="currentColor" 
              viewBox="0 0 24 24"
            >
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
            </svg>
          </button>

          <!-- Dropdown Menu -->
          <transition
            enter-active-class="transition ease-out duration-200"
            enter-from-class="transform opacity-0 scale-95"
            enter-to-class="transform opacity-100 scale-100"
            leave-active-class="transition ease-in duration-150"
            leave-from-class="transform opacity-100 scale-100"
            leave-to-class="transform opacity-0 scale-95"
          >
            <div
              v-if="isMenuOpen"
              class="absolute right-0 mt-2 w-56 rounded-2xl shadow-xl bg-white border border-gray-200 overflow-hidden"
              @click.stop
            >
              <!-- User Info (mobile only) -->
              <div class="md:hidden px-4 py-3 border-b border-gray-200 bg-gray-50">
                <p class="text-sm font-semibold text-gray-900">{{ user?.full_name }}</p>
                <p class="text-xs text-gray-500">{{ user?.email }}</p>
                <p class="text-xs text-gray-500">{{ user?.role }} - {{ user?.department }}</p>
              </div>

              <!-- Menu Items -->
              <div class="py-2">
                <router-link
                  to="/profile"
                  class="flex items-center space-x-3 px-4 py-3 hover:bg-gray-50 transition-colors"
                  @click="closeMenu"
                >
                  <svg class="w-5 h-5 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
                  </svg>
                  <span class="text-sm font-medium text-gray-700">Profile</span>
                </router-link>

                <button
                  @click="handleLogout"
                  class="w-full flex items-center space-x-3 px-4 py-3 hover:bg-red-50 transition-colors text-left"
                >
                  <svg class="w-5 h-5 text-red-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
                  </svg>
                  <span class="text-sm font-medium text-red-600">Keluar</span>
                </button>
              </div>
            </div>
          </transition>
        </div>
      </div>
    </div>

    <!-- Click outside to close -->
    <div
      v-if="isMenuOpen"
      class="fixed inset-0 z-40"
      @click="closeMenu"
    ></div>
  </nav>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../../stores/auth'
import { useAuth } from '../../composables/useAuth'

const router = useRouter()
const authStore = useAuthStore()
const { logout } = useAuth()

const user = computed(() => authStore.user)
const isMenuOpen = ref(false)

const userInitial = computed(() => {
  if (!user.value?.full_name) return '?'
  return user.value.full_name
    .split(' ')
    .map(n => n[0])
    .slice(0, 2)
    .join('')
    .toUpperCase()
})

const toggleMenu = () => {
  isMenuOpen.value = !isMenuOpen.value
}

const closeMenu = () => {
  isMenuOpen.value = false
}

const handleLogout = async () => {
  closeMenu()
  await logout()
}

// Close menu on escape key
const handleEscape = (e) => {
  if (e.key === 'Escape') {
    closeMenu()
  }
}

onMounted(() => {
  document.addEventListener('keydown', handleEscape)
})

onUnmounted(() => {
  document.removeEventListener('keydown', handleEscape)
})
</script>

<style scoped>
.glass-navbar {
  backdrop-filter: blur(16px) saturate(180%);
  -webkit-backdrop-filter: blur(16px) saturate(180%);
}

.active-scale {
  @apply transform transition-transform duration-150 active:scale-95;
}
</style>
