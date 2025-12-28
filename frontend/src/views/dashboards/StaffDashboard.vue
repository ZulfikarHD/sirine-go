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
      <p class="text-gray-600">Dashboard {{ user?.department }} - Shift {{ user?.shift }}</p>
    </Motion>

    <!-- My Tasks -->
    <Motion
      :initial="{ opacity: 0 }"
      :animate="{ opacity: 1 }"
      :transition="{ duration: 0.25, delay: 0.1 }"
      class="mb-8"
    >
      <h2 class="text-xl font-semibold text-gray-900 mb-4">Tugas Hari Ini</h2>
      <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
        <Motion 
          v-for="(task, index) in tasks" 
          :key="index"
          :initial="{ opacity: 0, y: 15 }"
          :animate="{ opacity: 1, y: 0 }"
          :transition="{ 
            duration: 0.25, 
            delay: 0.15 + (index * 0.05),
            ease: 'easeOut'
          }"
          class="glass-card p-6 rounded-2xl active-scale"
        >
          <div class="flex items-start justify-between mb-4">
            <div :class="`w-12 h-12 rounded-xl ${task.bgColor} flex items-center justify-center`">
              <component :is="task.icon" class="w-6 h-6 text-white" />
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
        </Motion>
      </div>
    </Motion>

    <!-- Performance -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <!-- Today's Performance -->
      <Motion
        :initial="{ opacity: 0, y: 10 }"
        :animate="{ opacity: 1, y: 0 }"
        :transition="{ duration: 0.25, delay: 0.3 }"
        class="glass-card p-6 rounded-2xl"
      >
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
      </Motion>

      <!-- Notifications -->
      <Motion
        :initial="{ opacity: 0, y: 10 }"
        :animate="{ opacity: 1, y: 0 }"
        :transition="{ duration: 0.25, delay: 0.35 }"
        class="glass-card p-6 rounded-2xl"
      >
        <h2 class="text-xl font-semibold text-gray-900 mb-4">Notifikasi</h2>
        <div class="space-y-3">
          <div 
            v-for="(notif, index) in notifications" 
            :key="index"
            class="flex items-start space-x-3 p-3 rounded-xl hover:bg-gray-50"
          >
            <div :class="`w-2 h-2 rounded-full ${notif.color} mt-2 shrink-0`"></div>
            <div class="flex-1 min-w-0">
              <p class="text-sm font-medium text-gray-900">{{ notif.title }}</p>
              <p class="text-xs text-gray-500">{{ notif.time }}</p>
            </div>
          </div>
        </div>
      </Motion>
    </div>
  </AppLayout>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useAuthStore } from '../../stores/auth'
import AppLayout from '../../components/layout/AppLayout.vue'
import { Motion } from 'motion-v'
import { ClipboardList, CheckCircle, CheckSquare } from 'lucide-vue-next'

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
    icon: ClipboardList,
    bgColor: 'bg-indigo-500',
  },
  {
    title: 'QC Inspection',
    description: 'Inspeksi kualitas batch #B-2025-044',
    status: 'Pending',
    statusColor: 'bg-yellow-100 text-yellow-700',
    progress: '0/200',
    deadline: '16:00 WIB',
    icon: CheckCircle,
    bgColor: 'bg-emerald-500',
  },
  {
    title: 'Verifikasi Data',
    description: 'Verifikasi data produksi shift pagi',
    status: 'Completed',
    statusColor: 'bg-green-100 text-green-700',
    progress: '100/100',
    deadline: '12:00 WIB',
    icon: CheckSquare,
    bgColor: 'bg-fuchsia-500',
  },
])

// Performance metrics
const performance = ref([
  {
    label: 'Target Harian',
    value: '850/1000',
    percentage: '85%',
    color: 'bg-linear-to-r from-indigo-500 to-indigo-600',
  },
  {
    label: 'Quality Rate',
    value: '98.5%',
    percentage: '98.5%',
    color: 'bg-linear-to-r from-emerald-500 to-emerald-600',
  },
  {
    label: 'Efficiency',
    value: '92%',
    percentage: '92%',
    color: 'bg-linear-to-r from-fuchsia-500 to-fuchsia-600',
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
