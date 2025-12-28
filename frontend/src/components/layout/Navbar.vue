<template>
  <nav class="glass-navbar sticky top-0 z-50 backdrop-blur-xl bg-white/80 border-b border-gray-200/50">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <div class="flex justify-between items-center h-16">
        <!-- Left Section: Mobile Menu & Brand -->
        <div class="flex items-center">
          <!-- Mobile Menu Button -->
          <button 
            @click="$emit('toggleSidebar')"
            class="md:hidden mr-3 p-2 -ml-2 rounded-xl text-gray-500 hover:bg-gray-100 hover:text-gray-900 focus:outline-none active:scale-95 transition-all"
          >
            <span class="sr-only">Open sidebar</span>
            <Menu class="w-6 h-6" />
          </button>

          <!-- Logo & Title (Hidden on Desktop if using Sidebar, but keeping for now as it might be used without sidebar) -->
          <div class="flex items-center space-x-3 md:hidden">
            <div class="w-10 h-10 rounded-xl bg-linear-to-br from-indigo-500 to-fuchsia-600 flex items-center justify-center shadow-lg">
              <Siren class="w-6 h-6 text-white" />
            </div>
            <div>
              <h1 class="text-xl font-bold text-gray-900">Sirine Go</h1>
            </div>
          </div>
          
          <!-- Desktop Title/Breadcrumbs -->
          <div class="hidden md:flex items-center ml-4">
             <Breadcrumbs />
          </div>
        </div>

        <!-- User Menu -->
        <div class="relative">
          <button
            @click="toggleMenu"
            class="flex items-center space-x-3 px-3 py-2 rounded-xl hover:bg-gray-100 transition-colors active-scale"
          >
            <!-- Avatar -->
            <div class="w-10 h-10 rounded-full bg-linear-to-br from-indigo-500 to-fuchsia-600 flex items-center justify-center text-white font-semibold shadow-lg">
              {{ userInitial }}
            </div>
            <!-- User Info (hidden on mobile) -->
            <div class="hidden md:block text-left">
              <p class="text-sm font-semibold text-gray-900">{{ user?.full_name }}</p>
              <p class="text-xs text-gray-500">{{ user?.role }}</p>
            </div>
            <!-- Dropdown Icon -->
            <ChevronDown 
              class="w-5 h-5 text-gray-500 transition-transform duration-200" 
              :class="{ 'rotate-180': isMenuOpen }"
            />
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
                  <User class="w-5 h-5 text-gray-500" />
                  <span class="text-sm font-medium text-gray-700">Profile</span>
                </router-link>

                <button
                  @click="handleLogout"
                  class="w-full flex items-center space-x-3 px-4 py-3 hover:bg-red-50 transition-colors text-left"
                >
                  <LogOut class="w-5 h-5 text-red-500" />
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
import { Menu, Siren, ChevronDown, User, LogOut } from 'lucide-vue-next'
import Breadcrumbs from '../common/Breadcrumbs.vue'

const props = defineProps({
  showBrand: {
    type: Boolean,
    default: true
  }
})

const emit = defineEmits(['toggleSidebar'])

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
</style>
