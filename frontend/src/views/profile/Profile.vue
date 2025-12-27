<template>
  <div class="min-h-screen bg-gray-50">
    <Navbar />
    
    <div class="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <!-- Profile Header -->
      <div class="glass-card p-8 rounded-2xl mb-6">
        <div class="flex flex-col md:flex-row items-center md:items-start space-y-4 md:space-y-0 md:space-x-6">
          <!-- Avatar -->
          <div class="relative">
            <div class="w-24 h-24 rounded-full bg-gradient-to-br from-indigo-500 to-fuchsia-600 flex items-center justify-center text-white text-3xl font-bold shadow-xl">
              {{ userInitial }}
            </div>
            <button class="absolute bottom-0 right-0 w-8 h-8 bg-indigo-600 rounded-full flex items-center justify-center text-white shadow-lg hover:bg-indigo-700 transition-colors">
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 9a2 2 0 012-2h.93a2 2 0 001.664-.89l.812-1.22A2 2 0 0110.07 4h3.86a2 2 0 011.664.89l.812 1.22A2 2 0 0018.07 7H19a2 2 0 012 2v9a2 2 0 01-2 2H5a2 2 0 01-2-2V9z" />
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 13a3 3 0 11-6 0 3 3 0 016 0z" />
              </svg>
            </button>
          </div>

          <!-- User Info -->
          <div class="flex-1 text-center md:text-left">
            <h1 class="text-3xl font-bold text-gray-900 mb-2">{{ user?.full_name }}</h1>
            <div class="flex flex-wrap items-center justify-center md:justify-start gap-2 mb-4">
              <span class="px-3 py-1 rounded-full text-xs font-semibold bg-indigo-100 text-indigo-700">
                {{ user?.role }}
              </span>
              <span class="px-3 py-1 rounded-full text-xs font-semibold bg-fuchsia-100 text-fuchsia-700">
                {{ user?.department }}
              </span>
              <span class="px-3 py-1 rounded-full text-xs font-semibold bg-emerald-100 text-emerald-700">
                Shift {{ user?.shift }}
              </span>
            </div>
            <div class="space-y-1 text-sm text-gray-600">
              <p class="flex items-center justify-center md:justify-start">
                <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" />
                </svg>
                {{ user?.email }}
              </p>
              <p class="flex items-center justify-center md:justify-start">
                <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 6H5a2 2 0 00-2 2v9a2 2 0 002 2h14a2 2 0 002-2V8a2 2 0 00-2-2h-5m-4 0V5a2 2 0 114 0v1m-4 0a2 2 0 104 0m-5 8a2 2 0 100-4 2 2 0 000 4zm0 0c1.306 0 2.417.835 2.83 2M9 14a3.001 3.001 0 00-2.83 2M15 11h3m-3 4h2" />
                </svg>
                NIP: {{ user?.nip }}
              </p>
              <p v-if="user?.phone" class="flex items-center justify-center md:justify-start">
                <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 5a2 2 0 012-2h3.28a1 1 0 01.948.684l1.498 4.493a1 1 0 01-.502 1.21l-2.257 1.13a11.042 11.042 0 005.516 5.516l1.13-2.257a1 1 0 011.21-.502l4.493 1.498a1 1 0 01.684.949V19a2 2 0 01-2 2h-1C9.716 21 3 14.284 3 6V5z" />
                </svg>
                {{ user?.phone }}
              </p>
            </div>
          </div>
        </div>
      </div>

      <!-- Actions -->
      <div class="grid grid-cols-1 md:grid-cols-2 gap-4 mb-6">
        <button class="glass-card p-6 rounded-2xl hover:shadow-lg transition-all text-left active-scale">
          <div class="flex items-center space-x-4">
            <div class="w-12 h-12 rounded-xl bg-indigo-500 flex items-center justify-center">
              <svg class="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
              </svg>
            </div>
            <div>
              <h3 class="font-semibold text-gray-900">Edit Profile</h3>
              <p class="text-sm text-gray-600">Ubah informasi profile Anda</p>
            </div>
          </div>
        </button>

        <button class="glass-card p-6 rounded-2xl hover:shadow-lg transition-all text-left active-scale">
          <div class="flex items-center space-x-4">
            <div class="w-12 h-12 rounded-xl bg-fuchsia-500 flex items-center justify-center">
              <svg class="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 7a2 2 0 012 2m4 0a6 6 0 01-7.743 5.743L11 17H9v2H7v2H4a1 1 0 01-1-1v-2.586a1 1 0 01.293-.707l5.964-5.964A6 6 0 1121 9z" />
              </svg>
            </div>
            <div>
              <h3 class="font-semibold text-gray-900">Ganti Password</h3>
              <p class="text-sm text-gray-600">Update password keamanan Anda</p>
            </div>
          </div>
        </button>
      </div>

      <!-- Account Info -->
      <div class="glass-card p-6 rounded-2xl">
        <h2 class="text-xl font-semibold text-gray-900 mb-4">Informasi Akun</h2>
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div class="p-4 bg-gray-50 rounded-xl">
            <p class="text-sm text-gray-600 mb-1">Tanggal Bergabung</p>
            <p class="font-semibold text-gray-900">{{ formatDate(user?.created_at) }}</p>
          </div>
          <div class="p-4 bg-gray-50 rounded-xl">
            <p class="text-sm text-gray-600 mb-1">Last Login</p>
            <p class="font-semibold text-gray-900">{{ formatDate(user?.last_login_at) }}</p>
          </div>
          <div class="p-4 bg-gray-50 rounded-xl">
            <p class="text-sm text-gray-600 mb-1">Status Akun</p>
            <p class="font-semibold text-emerald-600">{{ user?.status }}</p>
          </div>
          <div class="p-4 bg-gray-50 rounded-xl">
            <p class="text-sm text-gray-600 mb-1">ID User</p>
            <p class="font-semibold text-gray-900">#{{ user?.id }}</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useAuthStore } from '../../stores/auth'
import Navbar from '../../components/layout/Navbar.vue'

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

const formatDate = (dateString) => {
  if (!dateString) return '-'
  const date = new Date(dateString)
  return date.toLocaleDateString('id-ID', {
    day: 'numeric',
    month: 'long',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  })
}
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
  @apply transform transition-transform duration-150 active:scale-[0.98];
}
</style>
