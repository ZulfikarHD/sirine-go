<template>
  <!-- Mobile Backdrop dengan Motion-V animation -->
  <Teleport to="body">
    <Motion
      v-if="isOpen"
      :initial="modalAnimations.backdrop.initial"
      :animate="modalAnimations.backdrop.animate"
      :exit="modalAnimations.backdrop.exit"
      :transition="modalAnimations.backdrop.transition"
      class="fixed inset-0 z-40 bg-gray-900/60 md:hidden"
      @click="$emit('close')"
    />
  </Teleport>

  <!-- Sidebar Container -->
  <aside 
    class="sidebar-container"
    :class="[
      isOpen ? 'translate-x-0' : '-translate-x-full',
      'md:translate-x-0'
    ]"
  >
    <!-- Logo Section dengan gradient accent -->
    <div class="sidebar-header">
      <div class="flex items-center space-x-3">
        <div class="sidebar-logo">
          <Siren class="w-5 h-5 text-white" />
        </div>
        <div>
          <span class="text-lg font-bold text-gray-900 tracking-tight">Sirine Go</span>
          <span class="block text-[10px] font-medium text-gray-400 -mt-0.5">Production System</span>
        </div>
      </div>
    </div>

    <!-- Navigation dengan staggered animation -->
    <nav class="flex-1 overflow-y-auto py-4 px-3 space-y-1 custom-scrollbar">
      <template v-for="(group, groupIndex) in navigationGroups" :key="groupIndex">
        <!-- Group Title dengan fade animation -->
        <Motion
          v-if="group.title"
          :initial="{ opacity: 0, x: -10 }"
          :animate="{ opacity: 1, x: 0 }"
          :transition="{ duration: 0.2, delay: groupIndex * 0.08 }"
          class="px-3 mt-5 mb-2 first:mt-0"
        >
          <p class="nav-group-title">{{ group.title }}</p>
        </Motion>
        
        <!-- Navigation Items dengan stagger -->
        <Motion
          v-for="(item, itemIndex) in group.items"
          :key="item.name"
          v-bind="listItemAnimation(groupIndex * 3 + itemIndex, 0.05)"
          class="block"
        >
          <router-link
            :to="item.href"
            class="nav-item group"
            :class="[
              isActive(item.href) 
                ? 'nav-item-active' 
                : 'nav-item-inactive'
            ]"
            @click="$emit('close')"
          >
            <!-- Icon Container -->
            <div 
              class="nav-icon-container"
              :class="isActive(item.href) ? 'nav-icon-active' : 'nav-icon-inactive'"
            >
              <component :is="item.icon" class="w-[18px] h-[18px]" />
            </div>
            
            <!-- Label -->
            <span class="flex-1">{{ item.name }}</span>
            
            <!-- Active Indicator Dot -->
            <div 
              v-if="isActive(item.href)" 
              class="nav-active-indicator"
            />
            
            <!-- Hover Arrow untuk inactive items -->
            <ChevronRight 
              v-if="!isActive(item.href)"
              class="w-4 h-4 text-gray-300 opacity-0 group-hover:opacity-100 group-hover:translate-x-0.5"
              :style="{ transition: 'all 0.15s ease-out' }"
            />
          </router-link>
        </Motion>
      </template>
    </nav>

    <!-- User Profile Section (Bottom) -->
    <div class="sidebar-footer">
      <div class="user-profile-card">
        <!-- Avatar dengan gradient -->
        <div class="user-avatar">
          {{ userInitial }}
        </div>
        
        <!-- User Info -->
        <div class="flex-1 min-w-0">
          <p class="text-sm font-semibold text-gray-900 truncate">
            {{ user?.full_name || 'User' }}
          </p>
          <p class="text-xs text-gray-500 truncate flex items-center gap-1">
            <span class="w-1.5 h-1.5 rounded-full bg-emerald-400" />
            {{ formatRole(user?.role) }}
          </p>
        </div>
        
        <!-- Quick Settings -->
        <button 
          class="p-1.5 rounded-lg hover:bg-gray-100 text-gray-400 hover:text-gray-600 active-scale"
          aria-label="Pengaturan cepat"
        >
          <MoreVertical class="w-4 h-4" />
        </button>
      </div>
    </div>
  </aside>
</template>

<script setup>
/**
 * Sidebar Navigation Component
 * Komponen navigasi utama dengan Motion-V animations
 * untuk smooth transitions dan staggered menu items
 */
import { computed } from 'vue'
import { Motion } from 'motion-v'
import { useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { listItemAnimation, modalAnimations } from '@/composables/useMotion'
import { 
  Home, 
  Users, 
  ClipboardList, 
  Box, 
  BarChart3,
  Settings,
  FileText,
  Siren,
  UserCircle,
  ChevronRight,
  MoreVertical,
  Bell,
  ScrollText,
  Package,
  History,
  Activity,
  Printer,
  Calculator,
  Scissors
} from 'lucide-vue-next'

defineProps({
  isOpen: {
    type: Boolean,
    default: false
  }
})

defineEmits(['close'])

const route = useRoute()
const authStore = useAuthStore()
const user = computed(() => authStore.user)

/**
 * Generate user initials dari full_name
 * untuk avatar placeholder
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
    'STAFF_KHAZWAL': 'Staff Khazanah Awal',
    'SUPERVISOR_KHAZWAL': 'Supervisor Khazwal',
    'OPERATOR_CETAK': 'Operator Cetak',
    'SUPERVISOR_CETAK': 'Supervisor Cetak',
    'QC_INSPECTOR': 'QC Inspector',
    'VERIFIKATOR': 'Verifikator',
    'STAFF_KHAZKHIR': 'Staff Khazanah Akhir'
  }
  return roleMap[role] || role || 'Guest'
}

/**
 * Check active route dengan smart matching
 * Mencegah double-active pada sibling routes dengan level yang sama
 */
const isActive = (path) => {
  const currentPath = route.path
  
  // Exact match untuk dashboard
  if (path === '/dashboard') {
    return currentPath === '/dashboard'
  }
  
  // Special handling untuk Khazwal Material Prep routes
  if (path === '/khazwal/material-prep') {
    // Active jika di queue page atau detail/process pages
    // Tapi TIDAK active jika di history atau monitoring
    return currentPath === '/khazwal/material-prep' || 
           (currentPath.startsWith('/khazwal/material-prep/') && 
            !currentPath.startsWith('/khazwal/material-prep/history'))
  }
  
  if (path === '/khazwal/counting') {
    // Active untuk counting queue dan work pages
    return currentPath === '/khazwal/counting' || currentPath.startsWith('/khazwal/counting/')
  }
  
  if (path === '/khazwal/cutting') {
    // Active untuk cutting queue dan start pages
    return currentPath === '/khazwal/cutting' || currentPath.startsWith('/khazwal/cutting/')
  }
  
  if (path === '/khazwal/material-prep/history') {
    // Exact match untuk history
    return currentPath === '/khazwal/material-prep/history'
  }
  
  if (path === '/khazwal/monitoring') {
    // Exact match untuk monitoring
    return currentPath === '/khazwal/monitoring'
  }
  
  // For other routes, support nested children
  return currentPath === path || currentPath.startsWith(path + '/')
}

/**
 * Navigation structure dengan role-based visibility
 */
const navigationGroups = computed(() => {
  const isAdmin = user.value?.role === 'ADMIN' || user.value?.role === 'MANAGER'
  const isKhazwal = user.value?.role === 'STAFF_KHAZWAL'
  const isKhazwalSupervisor = user.value?.role === 'SUPERVISOR_KHAZWAL'
  const isCetak = user.value?.role === 'OPERATOR_CETAK' || user.value?.role === 'SUPERVISOR_CETAK'
  
  const groups = [
    {
      title: '',
      items: [
        { name: 'Dashboard', href: '/dashboard', icon: Home },
      ]
    },
    {
      title: 'Produksi',
      items: [
        { name: 'Pesanan (PO)', href: '/dashboard/orders', icon: ClipboardList },
        { name: 'Produksi', href: '/dashboard/production', icon: Box },
        { name: 'Laporan', href: '/dashboard/reports', icon: FileText },
      ]
    }
  ]

  // Khazwal Material Preparation - visible untuk STAFF_KHAZWAL, SUPERVISOR_KHAZWAL, ADMIN, MANAGER
  if (isKhazwal || isKhazwalSupervisor || isAdmin) {
    const khazwalItems = [
      { name: 'Persiapan Material', href: '/khazwal/material-prep', icon: Package },
      { name: 'Penghitungan', href: '/khazwal/counting', icon: Calculator },
      { name: 'Pemotongan', href: '/khazwal/cutting', icon: Scissors },
      { name: 'Riwayat', href: '/khazwal/material-prep/history', icon: History },
    ]
    
    // Monitoring hanya untuk Supervisor dan Admin
    if (isKhazwalSupervisor || isAdmin) {
      khazwalItems.push({ name: 'Monitoring', href: '/khazwal/monitoring', icon: Activity })
    }
    
    groups.push({
      title: 'Khazanah Awal',
      items: khazwalItems
    })
  }

  // Unit Cetak - visible untuk OPERATOR_CETAK, SUPERVISOR_CETAK, ADMIN, MANAGER
  if (isCetak || isAdmin) {
    groups.push({
      title: 'Unit Cetak',
      items: [
        { name: 'Antrian Cetak', href: '/cetak/queue', icon: Printer },
      ]
    })
  }

  if (isAdmin) {
    groups.push({
      title: 'Manajemen',
      items: [
        { name: 'Manajemen User', href: '/admin/users', icon: Users },
        { name: 'Audit Logs', href: '/admin/audit', icon: ScrollText },
        { name: 'Statistik', href: '/dashboard/stats', icon: BarChart3 },
      ]
    })
  }

  groups.push({
    title: 'Akun',
    items: [
      { name: 'Profile', href: '/profile', icon: UserCircle },
      { name: 'Notifikasi', href: '/notifications', icon: Bell },
      { name: 'Pengaturan', href: '/dashboard/settings', icon: Settings },
    ]
  })

  return groups
})
</script>
