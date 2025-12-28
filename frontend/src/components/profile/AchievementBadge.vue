<template>
  <Motion
    v-bind="badgeAnimation"
    @click="handleClick"
    :class="[
      'relative flex items-start gap-3 p-4 sm:p-3 rounded-xl bg-white/95 border border-gray-200/30 transition-all duration-200 ease-out',
      achievement.is_unlocked ? 'shadow-sm' : 'opacity-60',
      achievement.is_unlocked ? 'cursor-pointer active-scale hover:shadow-md hover:border-indigo-200/50 hover:bg-white hover:-translate-y-px' : ''
    ]"
    style="will-change: transform"
  >
    <!-- Icon dengan pop-in animation -->
    <div :class="[
      'relative flex-shrink-0 w-12 h-12 sm:w-10 sm:h-10 rounded-full flex items-center justify-center',
      achievement.is_unlocked ? 'bg-gradient-to-br from-indigo-50 to-fuchsia-50' : 'bg-gray-100'
    ]">
      <span 
        :class="[
          'text-3xl sm:text-2xl transition-[filter] duration-300 ease-out',
          !achievement.is_unlocked ? 'grayscale opacity-50' : ''
        ]"
        :style="!achievement.is_unlocked ? 'filter: grayscale(100%) opacity(0.5)' : ''"
      >
        {{ achievement.icon }}
      </span>
      
      <!-- Lock overlay untuk locked achievements -->
      <div v-if="!achievement.is_unlocked" class="absolute inset-0 flex items-center justify-center bg-gray-900/20 rounded-full text-gray-600">
        <svg class="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
          <path fill-rule="evenodd" d="M5 9V7a5 5 0 0110 0v2a2 2 0 012 2v5a2 2 0 01-2 2H5a2 2 0 01-2-2v-5a2 2 0 012-2zm8-2v2H7V7a3 3 0 016 0z" clip-rule="evenodd" />
        </svg>
      </div>
    </div>

    <!-- Badge content -->
    <div class="flex-1 min-w-0">
      <h4 :class="[
        'text-sm font-semibold mb-1 line-clamp-1',
        achievement.is_unlocked ? 'text-gray-900' : 'text-gray-400'
      ]">
        {{ achievement.name }}
      </h4>
      <p :class="[
        'text-xs mb-2 line-clamp-2',
        achievement.is_unlocked ? 'text-gray-600' : 'text-gray-400'
      ]">
        {{ achievement.description }}
      </p>
      
      <!-- Points dan unlock date -->
      <div class="flex items-center justify-between gap-2 text-xs">
        <span :class="[
          'font-medium',
          achievement.is_unlocked 
            ? 'bg-gradient-to-r from-indigo-600 to-fuchsia-600 bg-clip-text text-transparent' 
            : 'text-gray-400'
        ]">
          +{{ achievement.points }} pts
        </span>
        <span v-if="achievement.is_unlocked && achievement.unlocked_at" class="text-gray-500">
          {{ formatDate(achievement.unlocked_at) }}
        </span>
      </div>
    </div>

    <!-- Unlock effect dengan Motion-V -->
    <Motion
      v-if="showUnlockEffect"
      :initial="{ opacity: 0, scale: 0.8 }"
      :animate="{ opacity: [0, 1, 0], scale: [0.8, 1.05, 1.1] }"
      :transition="{ duration: 1, ease: 'easeOut' }"
      class="absolute inset-0 rounded-xl bg-gradient-to-r from-indigo-400/30 to-fuchsia-400/30 pointer-events-none"
    />
  </Motion>
</template>

<script setup>
import { computed, ref } from 'vue'
import { Motion } from 'motion-v'
import { listItemAnimation, springPresets } from '@/composables/useMotion'
import { useHaptic } from '@/composables/useHaptic'

const props = defineProps({
  achievement: {
    type: Object,
    required: true,
  },
  index: {
    type: Number,
    default: 0,
  },
  showUnlockAnimation: {
    type: Boolean,
    default: false,
  },
})

const emit = defineEmits(['click'])

const haptic = useHaptic()
const showUnlockEffect = ref(props.showUnlockAnimation)

/**
 * Animation untuk badge entrance dengan stagger
 */
const badgeAnimation = computed(() => {
  return listItemAnimation(props.index)
})

/**
 * Format date untuk display unlock date
 */
const formatDate = (dateString) => {
  if (!dateString) return ''
  
  const date = new Date(dateString)
  const now = new Date()
  const diffTime = Math.abs(now - date)
  const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24))
  
  if (diffDays === 0) return 'Hari ini'
  if (diffDays === 1) return 'Kemarin'
  if (diffDays < 7) return `${diffDays} hari lalu`
  if (diffDays < 30) return `${Math.floor(diffDays / 7)} minggu lalu`
  if (diffDays < 365) return `${Math.floor(diffDays / 30)} bulan lalu`
  
  return date.toLocaleDateString('id-ID', { 
    year: 'numeric', 
    month: 'short', 
    day: 'numeric' 
  })
}

/**
 * Handle click pada badge untuk show detail atau tooltip
 */
const handleClick = () => {
  if (props.achievement.is_unlocked) {
    haptic.light()
    emit('click', props.achievement)
  }
}

/**
 * Trigger unlock animation effect (called from parent)
 */
defineExpose({
  triggerUnlockEffect: () => {
    showUnlockEffect.value = true
    haptic.achievement()
    
    setTimeout(() => {
      showUnlockEffect.value = false
    }, 1000)
  },
})
</script>
