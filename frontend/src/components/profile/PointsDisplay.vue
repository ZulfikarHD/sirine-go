<template>
  <Motion v-bind="entranceAnimations.fadeUp" class="space-y-4">
    <!-- Points Card -->
    <div class="bg-white/95 border border-gray-200/30 rounded-2xl p-6 sm:p-5 shadow-sm">
      <div class="flex items-center justify-between mb-4">
        <div class="flex items-center gap-2">
          <span class="text-2xl">⭐</span>
          <span class="text-sm font-medium text-gray-600">Total Points</span>
        </div>
        <div :class="[
          'px-3 py-1 rounded-full text-xs font-semibold transition-all duration-200',
          level === 'Bronze' ? 'bg-orange-100 text-orange-700' : '',
          level === 'Silver' ? 'bg-gray-200 text-gray-700' : '',
          level === 'Gold' ? 'bg-yellow-100 text-yellow-700' : '',
          level === 'Platinum' ? 'bg-gradient-to-r from-indigo-100 to-fuchsia-100 text-indigo-700' : ''
        ]">
          {{ level }}
        </div>
      </div>

      <!-- Animated Points Counter -->
      <div class="text-5xl sm:text-4xl font-bold mb-4 bg-gradient-to-r from-indigo-600 to-fuchsia-600 bg-clip-text text-transparent" style="font-variant-numeric: tabular-nums">
        {{ displayPoints }}
      </div>

      <!-- Progress to Next Level -->
      <div v-if="nextLevel && pointsToNext > 0" class="space-y-2">
        <div class="flex items-center justify-between text-sm">
          <span class="text-gray-600 font-medium">{{ pointsToNext }} pts to {{ nextLevel }}</span>
          <span class="text-indigo-600 font-semibold">{{ progressPercentage }}%</span>
        </div>
        
        <!-- Progress Bar -->
        <div class="relative h-2 bg-gray-100 rounded-full overflow-hidden">
          <Motion
            class="absolute inset-y-0 left-0 bg-gradient-to-r from-indigo-500 to-fuchsia-500 rounded-full"
            :initial="{ width: 0 }"
            :animate="{ width: `${progressPercentage}%` }"
            :transition="{ duration: 1, ease: 'easeOut', delay: 0.3 }"
          />
        </div>
      </div>

      <!-- Max Level Message -->
      <div v-else-if="level === 'Platinum'" class="text-center py-2 font-medium">
        <span class="text-indigo-600">🎉 Max Level Reached!</span>
      </div>
    </div>

    <!-- Stats Grid -->
    <div v-if="showStats" class="grid grid-cols-3 gap-4 sm:gap-3">
      <div class="bg-white/95 border border-gray-200/30 rounded-xl p-4 sm:p-3 text-center">
        <div class="text-2xl sm:text-xl font-bold bg-gradient-to-r from-indigo-600 to-fuchsia-600 bg-clip-text text-transparent mb-1">
          {{ achievementsUnlocked }}
        </div>
        <div class="text-xs text-gray-600 font-medium">Achievements</div>
      </div>
      <div class="bg-white/95 border border-gray-200/30 rounded-xl p-4 sm:p-3 text-center">
        <div class="text-2xl sm:text-xl font-bold bg-gradient-to-r from-indigo-600 to-fuchsia-600 bg-clip-text text-transparent mb-1">
          {{ achievementPercentage }}%
        </div>
        <div class="text-xs text-gray-600 font-medium">Completion</div>
      </div>
      <div class="bg-white/95 border border-gray-200/30 rounded-xl p-4 sm:p-3 text-center">
        <div class="text-2xl sm:text-xl font-bold bg-gradient-to-r from-indigo-600 to-fuchsia-600 bg-clip-text text-transparent mb-1">
          {{ totalAchievements }}
        </div>
        <div class="text-xs text-gray-600 font-medium">Total Available</div>
      </div>
    </div>
  </Motion>
</template>

<script setup>
import { ref, computed, watch, onMounted } from 'vue'
import { Motion } from 'motion-v'
import { entranceAnimations } from '@/composables/useMotion'

const props = defineProps({
  points: {
    type: Number,
    default: 0,
  },
  level: {
    type: String,
    default: 'Bronze',
  },
  nextLevel: {
    type: String,
    default: 'Silver',
  },
  pointsToNext: {
    type: Number,
    default: 0,
  },
  achievementsUnlocked: {
    type: Number,
    default: 0,
  },
  totalAchievements: {
    type: Number,
    default: 0,
  },
  showStats: {
    type: Boolean,
    default: true,
  },
  animateOnMount: {
    type: Boolean,
    default: true,
  },
})

// Animated points value dengan number tween
const displayPoints = ref(0)

/**
 * Calculate progress percentage ke next level
 */
const progressPercentage = computed(() => {
  if (!props.pointsToNext) return 100
  
  const levelThresholds = {
    'Bronze': { current: 0, next: 100 },
    'Silver': { current: 100, next: 500 },
    'Gold': { current: 500, next: 1000 },
    'Platinum': { current: 1000, next: 1000 },
  }
  
  const threshold = levelThresholds[props.level]
  if (!threshold) return 0
  
  const currentProgress = props.points - threshold.current
  const totalRequired = threshold.next - threshold.current
  
  return Math.min(Math.round((currentProgress / totalRequired) * 100), 100)
})

/**
 * Calculate achievement completion percentage
 */
const achievementPercentage = computed(() => {
  if (props.totalAchievements === 0) return 0
  return Math.round((props.achievementsUnlocked / props.totalAchievements) * 100)
})

/**
 * Animate points counter dengan tween effect
 */
const animatePoints = (from, to, duration = 1000) => {
  const startTime = Date.now()
  const difference = to - from
  
  const animate = () => {
    const currentTime = Date.now()
    const elapsed = currentTime - startTime
    const progress = Math.min(elapsed / duration, 1)
    
    // Easing function untuk smooth animation (easeOutQuart)
    const eased = 1 - Math.pow(1 - progress, 4)
    
    displayPoints.value = Math.floor(from + (difference * eased))
    
    if (progress < 1) {
      requestAnimationFrame(animate)
    } else {
      displayPoints.value = to
    }
  }
  
  requestAnimationFrame(animate)
}

/**
 * Watch for points changes dan animate
 */
watch(() => props.points, (newValue, oldValue) => {
  if (newValue !== oldValue) {
    animatePoints(oldValue || 0, newValue)
  }
})

/**
 * Initialize animation on mount
 */
onMounted(() => {
  if (props.animateOnMount) {
    animatePoints(0, props.points)
  } else {
    displayPoints.value = props.points
  }
})
</script>
