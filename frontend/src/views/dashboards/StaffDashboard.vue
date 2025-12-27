<template>
  <div class="min-h-screen bg-gray-50">
    <Navbar />
    
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <!-- Welcome Section -->
      <div class="mb-8 animate-fade-in">
        <h1 class="text-3xl font-bold text-gray-900 mb-2">
          Selamat Datang, {{ user?.full_name }}!
        </h1>
        <p class="text-gray-600">Dashboard {{ user?.department }} - Shift {{ user?.shift }}</p>
      </div>

      <!-- My Tasks -->
      <div class="mb-8">
        <h2 class="text-xl font-semibold text-gray-900 mb-4">Tugas Hari Ini</h2>
        <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
          <div 
            v-for="(task, index) in tasks" 
            :key="index"
            class="glass-card p-6 rounded-2xl hover:shadow-lg transition-all active-scale"
          >
            <div class="flex items-start justify-between mb-4">
              <div :class="`w-12 h-12 rounded-xl ${task.bgColor} flex items-center justify-center`">
                <svg class="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" :d="task.icon" />
                </svg>
              </div>
              <span :class="`px-3 py-1 rounded-full text-xs font-semibold ${task.statusColor}`">
                {{ task.status }}
              </span>
            </div>
            <h3 class="text-lg font-semibold text-gray-900 mb-2">{{ task.title }}</h3>
            <p class="text-sm text-gray-600 mb-4">{{ task.description }}</p>
            <div class="flex items-center justify-between text-xs text-gray-500">
              <span>{{ task.progress }}</span>
              <span>{{ task.deadline }}</span>
            </div>
          </div>
        </div>
      </div>

      <!-- Performance -->
      <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <!-- Today's Performance -->
        <div class="glass-card p-6 rounded-2xl">
          <h2 class="text-xl font-semibold text-gray-900 mb-4">Performa Hari Ini</h2>
          <div class="space-y-4">
            <div v-for="(metric, index) in performance" :key="index">
              <div class="flex justify-between text-sm mb-2">
                <span class="font-medium text-gray-700">{{ metric.label }}</span>
                <span class="font-semibold text-gray-900">{{ metric.value }}</span>
              </div>
              <div class="w-full bg-gray-200 rounded-full h-2">
                <div 
                  :class="`h-2 rounded-full ${metric.color}`" 
                  :style="{ width: metric.percentage }"
                ></div>
              </div>
            </div>
          </div>
        </div>

        <!-- Notifications -->
        <div class="glass-card p-6 rounded-2xl">
          <h2 class="text-xl font-semibold text-gray-900 mb-4">Notifikasi</h2>
          <div class="space-y-3">
            <div 
              v-for="(notif, index) in notifications" 
              :key="index"
              class="flex items-start space-x-3 p-3 rounded-xl hover:bg-gray-50 transition-colors"
            >
              <div :class="`w-2 h-2 rounded-full ${notif.color} mt-2 flex-shrink-0`"></div>
              <div class="flex-1 min-w-0">
                <p class="text-sm font-medium text-gray-900">{{ notif.title }}</p>
                <p class="text-xs text-gray-500">{{ notif.time }}</p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useAuthStore } from '../../stores/auth'
import Navbar from '../../components/layout/Navbar.vue'

const authStore = useAuthStore()
const user = computed(() => authStore.user)

// Tasks data
const tasks = ref([
  {
    title: 'PO-2025-045',
    description: 'Cetak 500 unit pita cukai jenis A',
    status: 'In Progress',
    statusColor: 'bg-blue-100 text-blue-700',
    progress: '250/500',
    deadline: '14:00 WIB',
    icon: 'M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z',
    bgColor: 'bg-indigo-500',
  },
  {
    title: 'QC Inspection',
    description: 'Inspeksi kualitas batch #B-2025-044',
    status: 'Pending',
    statusColor: 'bg-yellow-100 text-yellow-700',
    progress: '0/200',
    deadline: '16:00 WIB',
    icon: 'M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z',
    bgColor: 'bg-emerald-500',
  },
  {
    title: 'Verifikasi Data',
    description: 'Verifikasi data produksi shift pagi',
    status: 'Completed',
    statusColor: 'bg-green-100 text-green-700',
    progress: '100/100',
    deadline: '12:00 WIB',
    icon: 'M5 13l4 4L19 7',
    bgColor: 'bg-fuchsia-500',
  },
])

// Performance metrics
const performance = ref([
  {
    label: 'Target Harian',
    value: '850/1000',
    percentage: '85%',
    color: 'bg-gradient-to-r from-indigo-500 to-indigo-600',
  },
  {
    label: 'Quality Rate',
    value: '98.5%',
    percentage: '98.5%',
    color: 'bg-gradient-to-r from-emerald-500 to-emerald-600',
  },
  {
    label: 'Efficiency',
    value: '92%',
    percentage: '92%',
    color: 'bg-gradient-to-r from-fuchsia-500 to-fuchsia-600',
  },
])

// Notifications
const notifications = ref([
  {
    title: 'PO Baru tersedia untuk dikerjakan',
    time: '5 menit lalu',
    color: 'bg-blue-500',
  },
  {
    title: 'Shift berakhir dalam 2 jam',
    time: '10 menit lalu',
    color: 'bg-yellow-500',
  },
  {
    title: 'QC Inspection approved untuk batch sebelumnya',
    time: '1 jam lalu',
    color: 'bg-green-500',
  },
])
</script>

<style scoped>
.glass-card {
  backdrop-filter: blur(16px) saturate(180%);
  -webkit-backdrop-filter: blur(16px) saturate(180%);
  background-color: rgba(255, 255, 255, 0.9);
  border: 1px solid rgba(209, 213, 219, 0.3);
}

.active-scale {
  @apply transform transition-transform duration-150 active:scale-[0.98];
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
