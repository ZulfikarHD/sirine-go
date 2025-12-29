<template>
  <Motion
    :initial="{ opacity: 0, y: 15 }"
    :animate="{ opacity: 1, y: 0 }"
    :transition="{ duration: 0.25, delay: index * 0.05, ease: 'easeOut' }"
    class="glass-card rounded-2xl p-4 cursor-pointer hover:shadow-lg active-scale"
    style="transition: box-shadow 0.15s ease-out"
    @click="$emit('click', staff)"
  >
    <div class="flex items-start gap-4">
      <!-- Avatar dengan Status Indicator -->
      <div class="relative flex-shrink-0">
        <div class="w-12 h-12 rounded-full bg-gradient-to-br from-indigo-500 to-fuchsia-500 
                    flex items-center justify-center">
          <span class="text-white text-lg font-bold">
            {{ staffInitial }}
          </span>
        </div>
        <!-- Status Dot -->
        <div 
          class="absolute -bottom-0.5 -right-0.5 w-4 h-4 rounded-full border-2 border-white"
          :class="isActive ? 'bg-emerald-500' : 'bg-gray-400'"
        />
      </div>

      <!-- Info -->
      <div class="flex-1 min-w-0">
        <div class="flex items-center justify-between gap-2 mb-1">
          <h4 class="text-base font-semibold text-gray-900 truncate">
            {{ staff.name }}
          </h4>
          <span 
            :class="statusBadgeClass"
            class="px-2 py-0.5 rounded-full text-xs font-medium flex-shrink-0"
          >
            {{ statusLabel }}
          </span>
        </div>

        <!-- Current PO jika active -->
        <div v-if="isActive && staff.current_po" class="mb-2">
          <p class="text-sm font-medium text-indigo-600 truncate">
            PO #{{ staff.current_po }}
          </p>
          <p class="text-xs text-gray-500 truncate">
            {{ staff.product_name }}
          </p>
        </div>

        <!-- Duration / Last Activity -->
        <div class="flex items-center gap-2 text-sm">
          <Clock class="w-4 h-4 text-gray-400" />
          <span v-if="isActive" class="text-gray-600">
            {{ formatDuration(staff.duration_mins) }} berlangsung
          </span>
          <span v-else class="text-gray-500">
            Tidak ada aktivitas
          </span>
        </div>
      </div>
    </div>
  </Motion>
</template>

<script setup>
/**
 * StaffActivityCard - Card untuk menampilkan aktivitas staff
 * dengan status indicator dan current PO info
 */
import { computed } from 'vue'
import { Motion } from 'motion-v'
import { Clock } from 'lucide-vue-next'

const props = defineProps({
  /**
   * Staff activity data
   */
  staff: {
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

defineEmits(['click'])

/**
 * Check jika staff sedang aktif
 */
const isActive = computed(() => {
  return props.staff.status === 'active' && props.staff.current_po
})

/**
 * Get staff initials
 */
const staffInitial = computed(() => {
  if (!props.staff.name) return '?'
  return props.staff.name
    .split(' ')
    .map(n => n[0])
    .slice(0, 2)
    .join('')
    .toUpperCase()
})

/**
 * Status badge class
 */
const statusBadgeClass = computed(() => {
  if (isActive.value) {
    return 'bg-emerald-100 text-emerald-700'
  }
  return 'bg-gray-100 text-gray-600'
})

/**
 * Status label
 */
const statusLabel = computed(() => {
  return isActive.value ? 'Aktif' : 'Idle'
})

/**
 * Format duration untuk display
 */
const formatDuration = (minutes) => {
  if (!minutes) return '0 menit'
  
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
</script>
