<template>
  <AppLayout>
    <!-- Profile Header -->
    <div class="glass-card p-8 rounded-2xl mb-6">
      <div class="flex flex-col md:flex-row items-center md:items-start space-y-4 md:space-y-0 md:space-x-6">
        <!-- Avatar -->
        <div class="relative">
          <div class="w-24 h-24 rounded-full bg-linear-to-br from-indigo-500 to-fuchsia-600 flex items-center justify-center text-white text-3xl font-bold shadow-xl">
            {{ userInitial }}
          </div>
          <button class="absolute bottom-0 right-0 w-8 h-8 bg-indigo-600 rounded-full flex items-center justify-center text-white shadow-lg hover:bg-indigo-700 transition-colors">
            <Camera class="w-4 h-4" />
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
              <Mail class="w-4 h-4 mr-2" />
              {{ user?.email }}
            </p>
            <p class="flex items-center justify-center md:justify-start">
              <IdCard class="w-4 h-4 mr-2" />
              NIP: {{ user?.nip }}
            </p>
            <p v-if="user?.phone" class="flex items-center justify-center md:justify-start">
              <Phone class="w-4 h-4 mr-2" />
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
            <UserCog class="w-6 h-6 text-white" />
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
            <Lock class="w-6 h-6 text-white" />
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
  </AppLayout>
</template>

<script setup>
import { computed } from 'vue'
import { useAuthStore } from '../../stores/auth'
import AppLayout from '../../components/layout/AppLayout.vue'
import { Camera, Mail, IdCard, Phone, UserCog, Lock } from 'lucide-vue-next'

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
</style>
