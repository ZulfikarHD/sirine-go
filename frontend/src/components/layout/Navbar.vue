<template>
  <nav class="navbar">
    <div class="navbar-inner">
      <!-- Left Section: Mobile Menu & Brand -->
      <div class="flex items-center">
        <!-- Mobile Menu Button -->
        <button 
          @click="$emit('toggleSidebar')"
          class="navbar-menu-btn md:hidden"
          aria-label="Toggle sidebar"
        >
          <Menu class="w-5 h-5" />
        </button>

        <!-- Logo & Title (Mobile only) -->
        <div class="flex items-center space-x-2.5 md:hidden">
          <div class="navbar-logo">
            <Siren class="w-5 h-5 text-white" />
          </div>
          <div>
            <h1 class="text-base font-bold text-gray-900 tracking-tight">Sirine Go</h1>
          </div>
        </div>
        
        <!-- Desktop Breadcrumbs -->
        <div class="hidden md:flex items-center ml-2">
          <Breadcrumbs />
        </div>
      </div>

      <!-- Right Section: User Menu -->
      <div class="flex items-center gap-2">
        <!-- Notification Bell (Placeholder) -->
        <button 
          class="navbar-icon-btn hidden sm:flex"
          aria-label="Notifikasi"
        >
          <Bell class="w-5 h-5" />
          <!-- Notification Badge -->
          <span class="notification-badge">2</span>
        </button>

        <!-- User Menu Trigger -->
        <div class="relative" ref="menuRef">
          <button
            @click="toggleMenu"
            class="user-menu-btn"
            :class="{ 'user-menu-btn-active': isMenuOpen }"
          >
            <!-- Avatar -->
            <div class="user-avatar-sm">
              {{ userInitial }}
            </div>
            <!-- User Info (hidden on mobile) -->
            <div class="hidden md:block text-left">
              <p class="text-sm font-semibold text-gray-900 leading-tight">{{ user?.full_name }}</p>
              <p class="text-xs text-gray-500 leading-tight">{{ formatRole(user?.role) }}</p>
            </div>
            <!-- Dropdown Icon -->
            <ChevronDown 
              class="w-4 h-4 text-gray-400" 
              :class="{ 'rotate-180': isMenuOpen }"
              :style="{ transition: 'transform 0.2s ease-out' }"
            />
          </button>

          <!-- Dropdown Menu -->
          <Teleport to="body">
            <div v-if="isMenuOpen" class="dropdown-backdrop" @click="closeMenu" />
            <Motion
              v-if="isMenuOpen"
              :initial="{ opacity: 0, scale: 0.95, y: -8 }"
              :animate="{ opacity: 1, scale: 1, y: 0 }"
              :exit="{ opacity: 0, scale: 0.95, y: -8 }"
              :transition="{ type: 'spring', stiffness: 500, damping: 35, mass: 0.6 }"
              class="dropdown-menu"
              :style="dropdownPosition"
            >
              <!-- User Info Header -->
              <div class="dropdown-header">
                <div class="flex items-center gap-3">
                  <div class="user-avatar-md">
                    {{ userInitial }}
                  </div>
                  <div class="flex-1 min-w-0">
                    <p class="text-sm font-semibold text-gray-900 truncate">{{ user?.full_name }}</p>
                    <p class="text-xs text-gray-500 truncate">{{ user?.email }}</p>
                  </div>
                </div>
                <div class="mt-2 flex items-center gap-2">
                  <span class="dropdown-role-badge">
                    {{ formatRole(user?.role) }}
                  </span>
                  <span v-if="user?.department" class="dropdown-dept-badge">
                    {{ user?.department }}
                  </span>
                </div>
              </div>

              <!-- Menu Items -->
              <div class="dropdown-body">
                <router-link
                  to="/profile"
                  class="dropdown-item"
                  @click="closeMenu"
                >
                  <div class="dropdown-item-icon dropdown-item-icon-default">
                    <User class="w-4 h-4" />
                  </div>
                  <div>
                    <span class="dropdown-item-label">Profile</span>
                    <span class="dropdown-item-desc">Kelola akun Anda</span>
                  </div>
                </router-link>

                <router-link
                  to="/profile/change-password"
                  class="dropdown-item"
                  @click="closeMenu"
                >
                  <div class="dropdown-item-icon dropdown-item-icon-default">
                    <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 7a2 2 0 012 2m4 0a6 6 0 01-7.743 5.743L11 17H9v2H7v2H4a1 1 0 01-1-1v-2.586a1 1 0 01.293-.707l5.964-5.964A6 6 0 1121 9z" />
                    </svg>
                  </div>
                  <div>
                    <span class="dropdown-item-label">Ganti Password</span>
                    <span class="dropdown-item-desc">Ubah password Anda</span>
                  </div>
                </router-link>

                <router-link
                  to="/dashboard/settings"
                  class="dropdown-item"
                  @click="closeMenu"
                >
                  <div class="dropdown-item-icon dropdown-item-icon-default">
                    <Settings class="w-4 h-4" />
                  </div>
                  <div>
                    <span class="dropdown-item-label">Pengaturan</span>
                    <span class="dropdown-item-desc">Preferensi aplikasi</span>
                  </div>
                </router-link>
              </div>

              <!-- Logout Section -->
              <div class="dropdown-footer">
                <button
                  @click="handleLogout"
                  class="dropdown-item dropdown-item-danger"
                  :disabled="isLoggingOut"
                >
                  <div class="dropdown-item-icon dropdown-item-icon-danger">
                    <LogOut v-if="!isLoggingOut" class="w-4 h-4" />
                    <div v-else class="w-4 h-4 border-2 border-red-500 border-t-transparent rounded-full animate-spin" />
                  </div>
                  <span class="dropdown-item-label">
                    {{ isLoggingOut ? 'Keluar...' : 'Keluar' }}
                  </span>
                </button>
              </div>
            </Motion>
          </Teleport>
        </div>
      </div>
    </div>
  </nav>
</template>

<script setup>
/**
 * Navbar Component
 * Top navigation bar dengan user dropdown menu
 * menggunakan Motion-V untuk smooth animations
 */
import { ref, computed, onMounted, onUnmounted, nextTick } from 'vue'
import { Motion } from 'motion-v'
import { useAuthStore } from '@/stores/auth'
import { useAuth } from '@/composables/useAuth'
import { Menu, Siren, ChevronDown, User, LogOut, Bell, Settings } from 'lucide-vue-next'
import Breadcrumbs from '@/components/common/Breadcrumbs.vue'

defineProps({
  showBrand: {
    type: Boolean,
    default: true
  }
})

defineEmits(['toggleSidebar'])

const authStore = useAuthStore()
const { logout } = useAuth()

const user = computed(() => authStore.user)
const isMenuOpen = ref(false)
const isLoggingOut = ref(false)
const menuRef = ref(null)
const dropdownPosition = ref({})

/**
 * Generate user initials untuk avatar
 */
const userInitial = computed(() => {
  if (!user.value?.full_name) return '?'
  return user.value.full_name
    .split(' ')
    .map(n => n[0])
    .slice(0, 2)
    .join('')
    .toUpperCase()
})

/**
 * Format role untuk display yang lebih readable
 */
const formatRole = (role) => {
  const roleMap = {
    'ADMIN': 'Administrator',
    'MANAGER': 'Manager',
    'OPERATOR': 'Operator',
    'VIEWER': 'Viewer'
  }
  return roleMap[role] || role || 'Guest'
}

/**
 * Calculate dropdown position relative to trigger button
 */
const updateDropdownPosition = async () => {
  await nextTick()
  if (menuRef.value) {
    const rect = menuRef.value.getBoundingClientRect()
    dropdownPosition.value = {
      position: 'fixed',
      top: `${rect.bottom + 8}px`,
      right: `${window.innerWidth - rect.right}px`,
      zIndex: 9999
    }
  }
}

const toggleMenu = async () => {
  isMenuOpen.value = !isMenuOpen.value
  if (isMenuOpen.value) {
    await updateDropdownPosition()
  }
}

const closeMenu = () => {
  isMenuOpen.value = false
}

/**
 * Handle logout dengan loading state
 */
const handleLogout = async () => {
  if (isLoggingOut.value) return
  isLoggingOut.value = true
  
  try {
    await logout()
  } finally {
    isLoggingOut.value = false
    closeMenu()
  }
}

/**
 * Keyboard & resize handlers
 */
const handleEscape = (e) => {
  if (e.key === 'Escape') closeMenu()
}

const handleResize = () => {
  if (isMenuOpen.value) updateDropdownPosition()
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
