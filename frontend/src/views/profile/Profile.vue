<template>
  <AppLayout>
    <!-- Profile Header -->
    <Motion
      v-bind="entranceAnimations.fadeUp"
      class="glass-card p-8 rounded-2xl mb-6"
    >
      <div class="flex flex-col md:flex-row items-center md:items-start space-y-4 md:space-y-0 md:space-x-6">
        <!-- Photo Upload Component -->
        <PhotoUpload
          :current-photo="user?.profile_photo_url"
          :alt-text="`${user?.full_name} Profile Photo`"
          @upload-success="handlePhotoUploadSuccess"
          @delete-success="handlePhotoDeleteSuccess"
        />

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
    </Motion>

    <!-- Points Display -->
    <PointsDisplay
      v-if="stats"
      :points="stats.total_points"
      :level="stats.level"
      :next-level="stats.next_level"
      :points-to-next="stats.points_to_next"
      :achievements-unlocked="stats.achievements_unlocked"
      :total-achievements="stats.total_achievements"
      class="mb-6"
    />

    <!-- Recent Achievements Preview -->
    <Motion
      v-if="recentAchievements && recentAchievements.length > 0"
      :initial="{ opacity: 0, y: 15 }"
      :animate="{ opacity: 1, y: 0 }"
      :transition="{ duration: 0.25, delay: 0.15, ease: 'easeOut' }"
      class="glass-card p-6 rounded-2xl mb-6"
    >
      <div class="flex items-center justify-between mb-4">
        <h2 class="text-xl font-semibold text-gray-900">Achievement Terbaru</h2>
        <button 
          @click="router.push('/profile/achievements')" 
          class="text-sm font-medium text-indigo-600 hover:text-indigo-700"
        >
          Lihat Semua →
        </button>
      </div>
      
      <div class="grid grid-cols-1 gap-3">
        <AchievementBadge
          v-for="(achievement, index) in recentAchievements.slice(0, 3)"
          :key="achievement.id"
          :achievement="achievement"
          :index="index"
        />
      </div>
    </Motion>

    <!-- Actions -->
    <div class="grid grid-cols-1 md:grid-cols-3 gap-4 mb-6">
      <Motion
        :initial="{ opacity: 0, y: 15 }"
        :animate="{ opacity: 1, y: 0 }"
        :transition="{ duration: 0.25, delay: 0.2, ease: 'easeOut' }"
      >
        <button 
          @click="router.push('/profile/edit')"
          class="glass-card p-6 rounded-2xl hover:shadow-lg text-left active-scale w-full"
        >
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
      </Motion>

      <Motion
        :initial="{ opacity: 0, y: 15 }"
        :animate="{ opacity: 1, y: 0 }"
        :transition="{ duration: 0.25, delay: 0.25, ease: 'easeOut' }"
      >
        <button 
          @click="router.push('/profile/change-password')" 
          class="glass-card p-6 rounded-2xl hover:shadow-lg text-left active-scale w-full"
        >
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
      </Motion>

      <Motion
        :initial="{ opacity: 0, y: 15 }"
        :animate="{ opacity: 1, y: 0 }"
        :transition="{ duration: 0.25, delay: 0.3, ease: 'easeOut' }"
      >
        <button 
          @click="router.push('/profile/achievements')" 
          class="glass-card p-6 rounded-2xl hover:shadow-lg text-left active-scale w-full"
        >
          <div class="flex items-center space-x-4">
            <div class="w-12 h-12 rounded-xl bg-emerald-500 flex items-center justify-center">
              <Trophy class="w-6 h-6 text-white" />
            </div>
            <div>
              <h3 class="font-semibold text-gray-900">Achievements</h3>
              <p class="text-sm text-gray-600">Lihat semua pencapaian</p>
            </div>
          </div>
        </button>
      </Motion>
    </div>

    <!-- Account Info -->
    <Motion
      :initial="{ opacity: 0, y: 15 }"
      :animate="{ opacity: 1, y: 0 }"
      :transition="{ duration: 0.25, delay: 0.35, ease: 'easeOut' }"
      class="glass-card p-6 rounded-2xl"
    >
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
    </Motion>
  </AppLayout>
</template>

<script setup>
import { computed, onMounted, ref } from 'vue'
import { Motion } from 'motion-v'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../../stores/auth'
import { useApi } from '../../composables/useApi'
import AppLayout from '../../components/layout/AppLayout.vue'
import PhotoUpload from '../../components/profile/PhotoUpload.vue'
import PointsDisplay from '../../components/profile/PointsDisplay.vue'
import AchievementBadge from '../../components/profile/AchievementBadge.vue'
import { entranceAnimations } from '../../composables/useMotion'
import { Mail, IdCard, Phone, UserCog, Lock, Trophy } from 'lucide-vue-next'

const router = useRouter()
const authStore = useAuthStore()
const { get } = useApi()

const user = computed(() => authStore.user)
const stats = ref(null)
const recentAchievements = ref([])

/**
 * Fetch user stats untuk points dan level display
 */
const fetchUserStats = async () => {
  try {
    const response = await get('/profile/stats')
    if (response.success) {
      stats.value = response.data
    }
  } catch (error) {
    console.error('Failed to fetch user stats:', error)
  }
}

/**
 * Fetch recent achievements untuk preview
 */
const fetchRecentAchievements = async () => {
  try {
    const response = await get('/profile/achievements')
    if (response.success) {
      // Filter hanya unlocked achievements dan sort by unlock date
      const unlocked = response.data.filter(a => a.is_unlocked)
      unlocked.sort((a, b) => new Date(b.unlocked_at) - new Date(a.unlocked_at))
      recentAchievements.value = unlocked
    }
  } catch (error) {
    console.error('Failed to fetch achievements:', error)
  }
}

/**
 * Handle photo upload success
 */
const handlePhotoUploadSuccess = async (photoUrl) => {
  // Update user profile photo di auth store
  authStore.updateUserField('profile_photo_url', photoUrl)
  
  // Refresh stats (mungkin ada achievement unlock)
  await fetchUserStats()
  await fetchRecentAchievements()
}

/**
 * Handle photo delete success
 */
const handlePhotoDeleteSuccess = () => {
  // Update user profile photo di auth store
  authStore.updateUserField('profile_photo_url', '')
}

/**
 * Format date untuk display
 */
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

/**
 * Initialize data on mount
 */
onMounted(async () => {
  await Promise.all([
    fetchUserStats(),
    fetchRecentAchievements(),
  ])
})
</script>
