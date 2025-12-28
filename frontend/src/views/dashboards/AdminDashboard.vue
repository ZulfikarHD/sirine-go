<template>
  <AppLayout>
    <!-- Welcome Section -->
    <Motion
      :initial="{ opacity: 0, y: 10 }"
      :animate="{ opacity: 1, y: 0 }"
      :transition="{ duration: 0.3, ease: 'easeOut' }"
      class="mb-8"
    >
      <h1 class="text-3xl font-bold text-gray-900 mb-2">
        Selamat Datang, {{ user?.full_name }}!
      </h1>
      <p class="text-gray-600">Dashboard Administrator - Sistem Produksi Pita Cukai</p>
    </Motion>

    <!-- Stats Cards -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
      <Motion 
        v-for="(stat, index) in stats" 
        :key="index"
        :initial="{ opacity: 0, y: 15 }"
        :animate="{ opacity: 1, y: 0 }"
        :transition="{ 
          duration: 0.25, 
          delay: index * 0.05,
          ease: 'easeOut'
        }"
        class="glass-card p-6 rounded-2xl active-scale cursor-pointer"
      >
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm font-medium text-gray-600 mb-1">{{ stat.label }}</p>
            <p class="text-3xl font-bold text-gray-900">{{ stat.value }}</p>
            <p class="text-xs text-gray-500 mt-1">{{ stat.change }}</p>
          </div>
          <div :class="`w-12 h-12 rounded-xl flex items-center justify-center ${stat.bgColor}`">
            <component :is="stat.icon" class="w-6 h-6 text-white" />
          </div>
        </div>
      </Motion>
    </div>

    <!-- Quick Actions -->
    <Motion
      :initial="{ opacity: 0 }"
      :animate="{ opacity: 1 }"
      :transition="{ duration: 0.25, delay: 0.2 }"
      class="mb-8"
    >
      <h2 class="text-xl font-semibold text-gray-900 mb-4">Quick Actions</h2>
      <div class="grid grid-cols-2 md:grid-cols-4 gap-4">
        <button 
          v-for="(action, index) in quickActions" 
          :key="index"
          class="glass-card p-4 rounded-xl flex flex-col items-center justify-center space-y-2 hover:bg-indigo-50 active-scale"
        >
          <div :class="`w-10 h-10 rounded-lg ${action.bgColor} flex items-center justify-center`">
            <component :is="action.icon" class="w-5 h-5 text-white" />
          </div>
          <span class="text-sm font-medium text-gray-700">{{ action.label }}</span>
        </button>
      </div>
    </Motion>

    <!-- Recent Activity -->
    <Motion
      :initial="{ opacity: 0, y: 10 }"
      :animate="{ opacity: 1, y: 0 }"
      :transition="{ duration: 0.25, delay: 0.25 }"
      class="glass-card p-6 rounded-2xl"
    >
      <h2 class="text-xl font-semibold text-gray-900 mb-4">Aktivitas Terkini</h2>
      <div class="space-y-4">
        <div 
          v-for="(activity, index) in recentActivities" 
          :key="index"
          class="flex items-center space-x-4 p-3 rounded-xl hover:bg-gray-50"
        >
          <div :class="`w-10 h-10 rounded-full ${activity.bgColor} flex items-center justify-center shrink-0`">
            <span class="text-white font-semibold">{{ activity.initial }}</span>
          </div>
          <div class="flex-1 min-w-0">
            <p class="text-sm font-medium text-gray-900">{{ activity.title }}</p>
            <p class="text-xs text-gray-500">{{ activity.description }}</p>
          </div>
          <div class="text-xs text-gray-400 shrink-0">
            {{ activity.time }}
          </div>
        </div>
      </div>
    </Motion>
  </AppLayout>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useAuthStore } from '../../stores/auth'
import AppLayout from '../../components/layout/AppLayout.vue'
import { Motion } from 'motion-v'
import { 
  Users, 
  ClipboardList, 
  Box, 
  CheckCircle,
  UserPlus,
  FileText,
  Settings
} from 'lucide-vue-next'

const authStore = useAuthStore()
const user = computed(() => authStore.user)

// Stats data
const stats = ref([
  {
    label: 'Total Users',
    value: '24',
    change: '+3 bulan ini',
    icon: Users,
    bgColor: 'bg-linear-to-br from-indigo-500 to-indigo-600',
  },
  {
    label: 'PO Aktif',
    value: '12',
    change: '8 pending approval',
    icon: ClipboardList,
    bgColor: 'bg-linear-to-br from-fuchsia-500 to-fuchsia-600',
  },
  {
    label: 'Produksi Hari Ini',
    value: '850',
    change: 'Target: 1000 unit',
    icon: Box,
    bgColor: 'bg-linear-to-br from-emerald-500 to-emerald-600',
  },
  {
    label: 'QC Pass Rate',
    value: '98.5%',
    change: '+2.1% dari kemarin',
    icon: CheckCircle,
    bgColor: 'bg-linear-to-br from-amber-500 to-amber-600',
  },
])

// Quick actions
const quickActions = ref([
  {
    label: 'Tambah User',
    icon: UserPlus,
    bgColor: 'bg-indigo-500',
  },
  {
    label: 'Buat PO',
    icon: ClipboardList,
    bgColor: 'bg-fuchsia-500',
  },
  {
    label: 'Laporan',
    icon: FileText,
    bgColor: 'bg-emerald-500',
  },
  {
    label: 'Settings',
    icon: Settings,
    bgColor: 'bg-gray-500',
  },
])

// Recent activities
const recentActivities = ref([
  {
    initial: 'ZH',
    title: 'User Baru Ditambahkan',
    description: 'Operator Cetak - Muhammad Rizki telah ditambahkan',
    time: '5 menit lalu',
    bgColor: 'bg-indigo-500',
  },
  {
    initial: 'PO',
    title: 'PO-2025-001 Disetujui',
    description: 'Production Order untuk 1000 unit pita cukai',
    time: '1 jam lalu',
    bgColor: 'bg-fuchsia-500',
  },
  {
    initial: 'QC',
    title: 'QC Inspection Selesai',
    description: 'Batch #B-2025-045 - Pass 985/1000',
    time: '2 jam lalu',
    bgColor: 'bg-emerald-500',
  },
])
</script>
