<template>
  <Motion
    :initial="{ opacity: 0, y: 15 }"
    :animate="{ opacity: 1, y: 0 }"
    :transition="{ duration: 0.25, delay: index * 0.05, ease: 'easeOut' }"
    @click="handleClick"
    class="glass-card rounded-2xl p-4 sm:p-5 cursor-pointer hover:shadow-lg active-scale"
    style="transition: box-shadow 0.15s ease-out"
  >
    <!-- Header Row -->
    <div class="flex items-start justify-between gap-3 mb-3">
      <div class="flex-1 min-w-0">
        <h3 class="text-lg font-bold text-gray-900 mb-1 truncate">
          PO #{{ item.po_number }}
        </h3>
        <p class="text-sm text-gray-600 truncate">
          OBC: {{ item.obc_number }}
        </p>
      </div>
      
      <PriorityBadge :priority="item.priority" size="sm" />
    </div>

    <!-- Product Name -->
    <p class="text-gray-800 font-medium mb-3 line-clamp-2">
      {{ item.product_name }}
    </p>

    <!-- Info Grid -->
    <div class="grid grid-cols-2 gap-3 mb-3">
      <!-- Completion Time -->
      <div class="flex items-center gap-2">
        <div class="p-1.5 rounded-lg bg-emerald-50">
          <CheckCircle class="w-4 h-4 text-emerald-600" />
        </div>
        <div class="min-w-0 flex-1">
          <p class="text-xs text-gray-500">Selesai</p>
          <p class="text-sm font-semibold text-emerald-700 truncate">
            {{ formattedCompletedAt }}
          </p>
        </div>
      </div>

      <!-- Duration -->
      <div class="flex items-center gap-2">
        <div class="p-1.5 rounded-lg bg-indigo-50">
          <Clock class="w-4 h-4 text-indigo-600" />
        </div>
        <div class="min-w-0 flex-1">
          <p class="text-xs text-gray-500">Durasi</p>
          <p class="text-sm font-semibold text-indigo-700 truncate">
            {{ formatDuration(item.duration_minutes) }}
          </p>
        </div>
      </div>
    </div>

    <!-- Footer -->
    <div class="flex items-center justify-between pt-3 border-t border-gray-200/50">
      <!-- Prepared By -->
      <div class="flex items-center gap-2 min-w-0 flex-1">
        <div class="w-6 h-6 rounded-full bg-gradient-to-br from-indigo-500 to-fuchsia-500 flex items-center justify-center flex-shrink-0">
          <span class="text-white text-xs font-bold">
            {{ preparedByInitial }}
          </span>
        </div>
        <span class="text-sm text-gray-600 truncate">
          {{ item.prepared_by_name || 'Unknown' }}
        </span>
      </div>
      
      <!-- Quantity & Photos -->
      <div class="flex items-center gap-2">
        <span class="text-xs text-gray-500 bg-gray-100 px-2 py-1 rounded-lg">
          {{ item.quantity.toLocaleString() }} pcs
        </span>
        <div 
          v-if="item.photos_count > 0" 
          class="flex items-center gap-1 px-2 py-1 rounded-lg bg-gray-100"
        >
          <Image class="w-3.5 h-3.5 text-gray-500" />
          <span class="text-xs font-medium text-gray-600">
            {{ item.photos_count }}
          </span>
        </div>
      </div>
    </div>
  </Motion>
</template>

<script setup>
/**
 * PrepHistoryCard - Card component untuk display single history item
 * dengan completion time, duration, dan prepared by info
 */
import { computed } from 'vue'
import { Motion } from 'motion-v'
import { CheckCircle, Clock, Image } from 'lucide-vue-next'
import PriorityBadge from '@/components/common/PriorityBadge.vue'

const props = defineProps({
  /**
   * History item data dari API response
   */
  item: {
    type: Object,
    required: true
  },
  
  /**
   * Index untuk stagger animation
   */
  index: {
    type: Number,
    default: 0
  }
})

const emit = defineEmits(['click'])

/**
 * Handle card click untuk view detail
 */
const handleClick = () => {
  emit('click', props.item)
  
  // Haptic feedback
  if ('vibrate' in navigator) {
    navigator.vibrate(10)
  }
}

/**
 * Format completed at timestamp
 */
const formattedCompletedAt = computed(() => {
  if (!props.item.completed_at) return '-'
  
  try {
    const date = new Date(props.item.completed_at)
    const today = new Date()
    const yesterday = new Date(today)
    yesterday.setDate(yesterday.getDate() - 1)
    
    const isToday = date.toDateString() === today.toDateString()
    const isYesterday = date.toDateString() === yesterday.toDateString()
    
    if (isToday) {
      return 'Hari ini, ' + date.toLocaleTimeString('id-ID', { hour: '2-digit', minute: '2-digit' })
    }
    if (isYesterday) {
      return 'Kemarin, ' + date.toLocaleTimeString('id-ID', { hour: '2-digit', minute: '2-digit' })
    }
    
    return date.toLocaleDateString('id-ID', { 
      day: 'numeric', 
      month: 'short',
      hour: '2-digit',
      minute: '2-digit'
    })
  } catch (e) {
    return props.item.completed_at
  }
})

/**
 * Format duration untuk display
 */
const formatDuration = (minutes) => {
  if (!minutes) return '-'
  
  if (minutes < 60) {
    return `${minutes} menit`
  }
  
  const hours = Math.floor(minutes / 60)
  const mins = minutes % 60
  
  if (mins === 0) {
    return `${hours} jam`
  }
  
  return `${hours}j ${mins}m`
}

/**
 * Get prepared by initials
 */
const preparedByInitial = computed(() => {
  if (!props.item.prepared_by_name) return '?'
  return props.item.prepared_by_name
    .split(' ')
    .map(n => n[0])
    .slice(0, 2)
    .join('')
    .toUpperCase()
})
</script>
