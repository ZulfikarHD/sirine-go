<template>
  <AppLayout>
    <!-- Page Header -->
    <Motion v-bind="entranceAnimations.fadeUp" class="mb-6">
      <div class="flex items-center justify-between">
        <div>
          <h1 class="text-3xl font-bold text-gray-900 mb-2">Achievements</h1>
          <p class="text-gray-600">Lacak progress dan pencapaian Anda</p>
        </div>
        <button 
          @click="router.back()"
          class="btn-secondary"
        >
          ← Kembali
        </button>
      </div>
    </Motion>

    <!-- Points Summary Card -->
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

    <!-- Category Tabs -->
    <Motion
      :initial="{ opacity: 0, y: 15 }"
      :animate="{ opacity: 1, y: 0 }"
      :transition="{ duration: 0.25, delay: 0.15 }"
      class="glass-card p-4 rounded-2xl mb-6"
    >
      <div class="flex flex-wrap items-center gap-2">
        <button
          v-for="category in categories"
          :key="category.value"
          @click="selectedCategory = category.value"
          :class="[
            'px-4 py-2 rounded-lg text-sm font-medium transition-all duration-200',
            selectedCategory === category.value
              ? 'bg-indigo-600 text-white shadow-md'
              : 'bg-gray-100 text-gray-600 hover:bg-gray-200'
          ]"
        >
          <span class="mr-2">{{ category.icon }}</span>
          {{ category.label }}
        </button>
      </div>
    </Motion>

    <!-- Achievements Grid -->
    <div v-if="loading" class="text-center py-12">
      <div class="spinner w-12 h-12 border-4 border-indigo-200 border-t-indigo-600 rounded-full mx-auto mb-4"></div>
      <p class="text-gray-600">Loading achievements...</p>
    </div>

    <div v-else-if="filteredAchievements.length === 0" class="glass-card p-12 rounded-2xl text-center">
      <div class="w-20 h-20 rounded-full bg-gray-100 mx-auto mb-4 flex items-center justify-center">
        <Trophy class="w-10 h-10 text-gray-400" />
      </div>
      <h3 class="text-xl font-semibold text-gray-900 mb-2">Tidak ada achievement</h3>
      <p class="text-gray-600">
        {{ selectedCategory === 'all' 
          ? 'Belum ada achievement yang tersedia' 
          : `Tidak ada achievement dalam kategori ${getCategoryLabel(selectedCategory)}`
        }}
      </p>
    </div>

    <div v-else class="grid grid-cols-1 lg:grid-cols-2 gap-4">
      <AchievementBadge
        v-for="(achievement, index) in filteredAchievements"
        :key="achievement.id"
        :achievement="achievement"
        :index="index"
        @click="showAchievementDetail(achievement)"
      />
    </div>

    <!-- Achievement Detail Modal -->
    <BaseModal
      v-model="showDetailModal"
      :title="selectedAchievement?.name || 'Achievement Detail'"
      size="sm"
    >
      <div v-if="selectedAchievement" class="text-center space-y-4">
        <!-- Icon -->
        <div class="w-24 h-24 rounded-full bg-gradient-to-br from-indigo-50 to-fuchsia-50 mx-auto flex items-center justify-center">
          <span class="text-6xl">{{ selectedAchievement.icon }}</span>
        </div>

        <!-- Description -->
        <p class="text-gray-600">{{ selectedAchievement.description }}</p>

        <!-- Points -->
        <div class="flex items-center justify-center gap-2">
          <span class="text-2xl font-bold bg-gradient-to-r from-indigo-600 to-fuchsia-600 bg-clip-text text-transparent">
            +{{ selectedAchievement.points }}
          </span>
          <span class="text-gray-600">points</span>
        </div>

        <!-- Unlock Date jika sudah unlocked -->
        <div v-if="selectedAchievement.is_unlocked" class="pt-4 border-t border-gray-200">
          <p class="text-sm text-gray-500">Unlocked on</p>
          <p class="font-semibold text-gray-900">{{ formatUnlockDate(selectedAchievement.unlocked_at) }}</p>
        </div>

        <!-- Locked Message -->
        <div v-else class="pt-4 border-t border-gray-200">
          <p class="text-sm text-gray-500">🔒 Belum di-unlock</p>
          <p class="text-xs text-gray-400 mt-1">Terus beraktivitas untuk unlock achievement ini!</p>
        </div>
      </div>
    </BaseModal>
  </AppLayout>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { Motion } from 'motion-v'
import { useRouter } from 'vue-router'
import { Trophy } from 'lucide-vue-next'
import { useApi } from '../../composables/useApi'
import { useHaptic } from '../../composables/useHaptic'
import AppLayout from '../../components/layout/AppLayout.vue'
import PointsDisplay from '../../components/profile/PointsDisplay.vue'
import AchievementBadge from '../../components/profile/AchievementBadge.vue'
import { BaseModal } from '../../components/common'
import { entranceAnimations } from '../../composables/useMotion'

const router = useRouter()
const { get } = useApi()
const haptic = useHaptic()

const loading = ref(true)
const achievements = ref([])
const stats = ref(null)
const selectedCategory = ref('all')
const showDetailModal = ref(false)
const selectedAchievement = ref(null)

/**
 * Achievement categories untuk filtering
 */
const categories = [
  { value: 'all', label: 'Semua', icon: '🏆' },
  { value: 'LOGIN', label: 'Login', icon: '🔑' },
  { value: 'PRODUCTIVITY', label: 'Produktivitas', icon: '⚡' },
  { value: 'QUALITY', label: 'Kualitas', icon: '⭐' },
  { value: 'MILESTONE', label: 'Milestone', icon: '🎯' },
]

/**
 * Get category label by value
 */
const getCategoryLabel = (value) => {
  const category = categories.find(c => c.value === value)
  return category ? category.label : value
}

/**
 * Filtered achievements berdasarkan selected category
 */
const filteredAchievements = computed(() => {
  if (selectedCategory.value === 'all') {
    return achievements.value
  }
  return achievements.value.filter(a => a.category === selectedCategory.value)
})

/**
 * Fetch achievements dari API
 */
const fetchAchievements = async () => {
  try {
    const response = await get('/profile/achievements')
    if (response.success) {
      achievements.value = response.data
    }
  } catch (error) {
    console.error('Failed to fetch achievements:', error)
  }
}

/**
 * Fetch user stats untuk points display
 */
const fetchStats = async () => {
  try {
    const response = await get('/profile/stats')
    if (response.success) {
      stats.value = response.data
    }
  } catch (error) {
    console.error('Failed to fetch stats:', error)
  }
}

/**
 * Show achievement detail modal
 */
const showAchievementDetail = (achievement) => {
  if (!achievement.is_unlocked) return // Don't show detail untuk locked achievements
  
  haptic.light()
  selectedAchievement.value = achievement
  showDetailModal.value = true
}

/**
 * Format unlock date untuk display
 */
const formatUnlockDate = (dateString) => {
  if (!dateString) return ''
  
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
  loading.value = true
  await Promise.all([
    fetchAchievements(),
    fetchStats(),
  ])
  loading.value = false
})
</script>
