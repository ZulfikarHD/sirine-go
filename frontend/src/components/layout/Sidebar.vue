<template>
  <!-- Mobile Backdrop -->
  <div 
    v-if="isOpen" 
    class="fixed inset-0 z-40 bg-gray-900/50 backdrop-blur-sm md:hidden"
    @click="$emit('close')"
  ></div>

  <!-- Sidebar Container -->
  <aside 
    class="fixed inset-y-0 left-0 z-50 w-64 bg-white/80 backdrop-blur-xl border-r border-gray-200/50 transform transition-transform duration-300 ease-in-out md:translate-x-0 md:static md:inset-auto md:flex md:flex-col"
    :class="isOpen ? 'translate-x-0' : '-translate-x-full'"
  >
    <!-- Logo Section -->
    <div class="h-16 flex items-center px-6 border-b border-gray-200/50">
      <div class="flex items-center space-x-3">
        <div class="w-8 h-8 rounded-lg bg-linear-to-br from-indigo-500 to-fuchsia-600 flex items-center justify-center shadow-lg">
          <Siren class="w-5 h-5 text-white" />
        </div>
        <span class="text-lg font-bold text-gray-900">Sirine Go</span>
      </div>
    </div>

    <!-- Navigation -->
    <nav class="flex-1 overflow-y-auto py-4 px-3 space-y-1 custom-scrollbar">
      <template v-for="(group, groupIndex) in navigationGroups" :key="groupIndex">
        <div v-if="group.title" class="px-3 mt-4 mb-2">
          <p class="text-xs font-semibold text-gray-400 uppercase tracking-wider">{{ group.title }}</p>
        </div>
        
        <router-link
          v-for="item in group.items"
          :key="item.name"
          :to="item.href"
          class="group flex items-center px-3 py-2.5 text-sm font-medium rounded-xl transition-all duration-200"
          :class="[
            isActive(item.href) 
              ? 'bg-linear-to-r from-indigo-50 to-fuchsia-50 text-indigo-700 shadow-sm ring-1 ring-indigo-200' 
              : 'text-gray-600 hover:bg-gray-50 hover:text-gray-900'
          ]"
          @click="$emit('close')"
        >
          <div 
            class="mr-3 shrink-0 transition-colors duration-200"
            :class="isActive(item.href) ? 'text-indigo-600' : 'text-gray-400 group-hover:text-gray-500'"
          >
            <component :is="item.icon" class="w-5 h-5" />
          </div>
          {{ item.name }}
          
          <!-- Active Indicator -->
          <div 
            v-if="isActive(item.href)" 
            class="ml-auto w-1.5 h-1.5 rounded-full bg-linear-to-br from-indigo-500 to-fuchsia-600"
          ></div>
        </router-link>
      </template>
    </nav>

    <!-- User Profile (Bottom) -->
    <div class="p-4 border-t border-gray-200/50 bg-gray-50/50">
      <div class="flex items-center space-x-3">
        <div class="w-9 h-9 rounded-full bg-linear-to-br from-indigo-500 to-fuchsia-600 flex items-center justify-center text-white text-sm font-bold shadow-md">
          {{ userInitial }}
        </div>
        <div class="flex-1 min-w-0">
          <p class="text-sm font-medium text-gray-900 truncate">{{ user?.full_name || 'User' }}</p>
          <p class="text-xs text-gray-500 truncate">{{ user?.role || 'Guest' }}</p>
        </div>
      </div>
    </div>
  </aside>
</template>

<script setup>
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { useAuthStore } from '../../stores/auth'
import { 
  Home, 
  Users, 
  ClipboardList, 
  Box, 
  BarChart3,
  Settings,
  FileText,
  Siren
} from 'lucide-vue-next'

const props = defineProps({
  isOpen: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['close'])

const route = useRoute()
const authStore = useAuthStore()
const user = computed(() => authStore.user)

const userInitial = computed(() => {
  if (!user.value?.full_name) return '?'
  return user.value.full_name
    .split(' ')
    .map(n => n[0])
    .slice(0, 2)
    .join('')
    .toUpperCase()
})

const isActive = (path) => {
  return route.path === path || route.path.startsWith(path + '/')
}

const navigationGroups = [
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
  },
  {
    title: 'Manajemen',
    items: [
      { name: 'Pengguna', href: '/dashboard/users', icon: Users },
      { name: 'Statistik', href: '/dashboard/stats', icon: BarChart3 },
    ]
  },
  {
    title: 'Sistem',
    items: [
      { name: 'Pengaturan', href: '/dashboard/settings', icon: Settings },
    ]
  }
]
</script>
