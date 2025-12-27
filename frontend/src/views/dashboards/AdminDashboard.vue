<template>
  <div class="min-h-screen bg-gray-50">
    <Navbar />
    
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <!-- Welcome Section -->
      <div class="mb-8 animate-fade-in">
        <h1 class="text-3xl font-bold text-gray-900 mb-2">
          Selamat Datang, {{ user?.full_name }}!
        </h1>
        <p class="text-gray-600">Dashboard Administrator - Sistem Produksi Pita Cukai</p>
      </div>

      <!-- Stats Cards -->
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
        <div 
          v-for="(stat, index) in stats" 
          :key="index"
          class="glass-card p-6 rounded-2xl transform transition-all duration-300 hover:scale-105 cursor-pointer"
          :style="{ animationDelay: `${index * 100}ms` }"
        >
          <div class="flex items-center justify-between">
            <div>
              <p class="text-sm font-medium text-gray-600 mb-1">{{ stat.label }}</p>
              <p class="text-3xl font-bold text-gray-900">{{ stat.value }}</p>
              <p class="text-xs text-gray-500 mt-1">{{ stat.change }}</p>
            </div>
            <div :class="`w-12 h-12 rounded-xl flex items-center justify-center ${stat.bgColor}`">
              <svg class="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" :d="stat.icon" />
              </svg>
            </div>
          </div>
        </div>
      </div>

      <!-- Quick Actions -->
      <div class="mb-8">
        <h2 class="text-xl font-semibold text-gray-900 mb-4">Quick Actions</h2>
        <div class="grid grid-cols-2 md:grid-cols-4 gap-4">
          <button 
            v-for="(action, index) in quickActions" 
            :key="index"
            class="glass-card p-4 rounded-xl flex flex-col items-center justify-center space-y-2 hover:bg-indigo-50 transition-colors active-scale"
          >
            <div :class="`w-10 h-10 rounded-lg ${action.bgColor} flex items-center justify-center`">
              <svg class="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" :d="action.icon" />
              </svg>
            </div>
            <span class="text-sm font-medium text-gray-700">{{ action.label }}</span>
          </button>
        </div>
      </div>

      <!-- Recent Activity -->
      <div class="glass-card p-6 rounded-2xl">
        <h2 class="text-xl font-semibold text-gray-900 mb-4">Aktivitas Terkini</h2>
        <div class="space-y-4">
          <div 
            v-for="(activity, index) in recentActivities" 
            :key="index"
            class="flex items-center space-x-4 p-3 rounded-xl hover:bg-gray-50 transition-colors"
          >
            <div :class="`w-10 h-10 rounded-full ${activity.bgColor} flex items-center justify-center flex-shrink-0`">
              <span class="text-white font-semibold">{{ activity.initial }}</span>
            </div>
            <div class="flex-1 min-w-0">
              <p class="text-sm font-medium text-gray-900">{{ activity.title }}</p>
              <p class="text-xs text-gray-500">{{ activity.description }}</p>
            </div>
            <div class="text-xs text-gray-400 flex-shrink-0">
              {{ activity.time }}
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useAuthStore } from '../../stores/auth'
import Navbar from '../../components/layout/Navbar.vue'
import { animate, stagger } from 'motion-dom'

const authStore = useAuthStore()
const user = computed(() => authStore.user)

// Stats data
const stats = ref([
  {
    label: 'Total Users',
    value: '24',
    change: '+3 bulan ini',
    icon: 'M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z',
    bgColor: 'bg-gradient-to-br from-indigo-500 to-indigo-600',
  },
  {
    label: 'PO Aktif',
    value: '12',
    change: '8 pending approval',
    icon: 'M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z',
    bgColor: 'bg-gradient-to-br from-fuchsia-500 to-fuchsia-600',
  },
  {
    label: 'Produksi Hari Ini',
    value: '850',
    change: 'Target: 1000 unit',
    icon: 'M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z',
    bgColor: 'bg-gradient-to-br from-emerald-500 to-emerald-600',
  },
  {
    label: 'QC Pass Rate',
    value: '98.5%',
    change: '+2.1% dari kemarin',
    icon: 'M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z',
    bgColor: 'bg-gradient-to-br from-amber-500 to-amber-600',
  },
])

// Quick actions
const quickActions = ref([
  {
    label: 'Tambah User',
    icon: 'M18 9v3m0 0v3m0-3h3m-3 0h-3m-2-5a4 4 0 11-8 0 4 4 0 018 0zM3 20a6 6 0 0112 0v1H3v-1z',
    bgColor: 'bg-indigo-500',
  },
  {
    label: 'Buat PO',
    icon: 'M12 4v16m8-8H4',
    bgColor: 'bg-fuchsia-500',
  },
  {
    label: 'Laporan',
    icon: 'M9 17v-2m3 2v-4m3 4v-6m2 10H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z',
    bgColor: 'bg-emerald-500',
  },
  {
    label: 'Settings',
    icon: 'M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z',
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

onMounted(() => {
  // Staggered animation untuk stats cards
  animate(
    '.glass-card',
    { opacity: [0, 1], transform: ['translateY(20px)', 'translateY(0)'] },
    { duration: 0.5, delay: stagger(0.1) }
  )
})
</script>

<style scoped>
@import "tailwindcss" reference;

.glass-card {
  backdrop-filter: blur(16px) saturate(180%);
  -webkit-backdrop-filter: blur(16px) saturate(180%);
  background-color: rgba(255, 255, 255, 0.9);
  border: 1px solid rgba(209, 213, 219, 0.3);
}

.active-scale {
  @apply transform transition-transform duration-150 active:scale-95;
}

.animate-fade-in {
  animation: fadeIn 0.6s ease-in-out;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(-10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}
</style>
