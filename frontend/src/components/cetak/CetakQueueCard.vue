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
        <p class="text-sm text-indigo-600 font-semibold truncate">
          OBC: {{ obcNumber }}
        </p>
      </div>
      
      <PriorityBadge :priority="item.priority" size="sm" />
    </div>

    <!-- Product Name -->
    <p class="text-gray-800 font-medium mb-2 line-clamp-2">
      {{ productName }}
    </p>

    <!-- OBC Master Details (if available) -->
    <div v-if="obcMaster" class="mb-3 p-2 bg-gray-50 rounded-lg">
      <div class="flex flex-wrap gap-2 text-xs">
        <span v-if="obcMaster.material" class="px-2 py-0.5 bg-white rounded-md text-gray-700 border border-gray-200">
          Material: {{ obcMaster.material }}
        </span>
        <span v-if="obcMaster.seri" class="px-2 py-0.5 bg-white rounded-md text-gray-700 border border-gray-200">
          Seri: {{ obcMaster.seri }}
        </span>
        <span v-if="obcMaster.warna" class="px-2 py-0.5 bg-white rounded-md text-gray-700 border border-gray-200">
          Warna: {{ obcMaster.warna }}
        </span>
        <span v-if="obcMaster.personalization && obcMaster.personalization !== '-'" class="px-2 py-0.5 bg-purple-50 rounded-md text-purple-700 border border-purple-200">
          {{ obcMaster.personalization }}
        </span>
        <span v-if="obcMaster.plat_number" class="px-2 py-0.5 bg-indigo-50 rounded-md text-indigo-700 border border-indigo-200">
          Plat: {{ obcMaster.plat_number }}
        </span>
      </div>
    </div>

    <!-- Info Grid -->
    <div class="grid grid-cols-2 gap-3 mb-3">
      <!-- Material Ready Info -->
      <div class="flex items-center gap-2">
        <div class="p-1.5 rounded-lg bg-emerald-50">
          <CheckCircle class="w-4 h-4 text-emerald-600" />
        </div>
        <div class="min-w-0 flex-1">
          <p class="text-xs text-gray-500">Material Siap</p>
          <p class="text-sm font-semibold text-emerald-700 truncate">
            {{ formattedMaterialReady }}
          </p>
        </div>
      </div>

      <!-- Quantity -->
      <div class="flex items-center gap-2">
        <div class="p-1.5 rounded-lg bg-indigo-50">
          <Package class="w-4 h-4 text-indigo-600" />
        </div>
        <div class="min-w-0 flex-1">
          <p class="text-xs text-gray-500">Quantity</p>
          <p class="text-sm font-semibold text-gray-900 truncate">
            {{ item.quantity.toLocaleString() }}
          </p>
        </div>
      </div>
    </div>

    <!-- Due Date Row -->
    <div class="flex items-center gap-2 mb-3">
      <div :class="dueDateIconClass" class="p-1.5 rounded-lg">
        <Calendar class="w-4 h-4" />
      </div>
      <div class="min-w-0 flex-1">
        <p class="text-xs text-gray-500">Deadline</p>
        <p :class="dueDateTextClass" class="text-sm font-semibold truncate">
          {{ formattedDueDate }}
        </p>
      </div>
      <!-- Days indicator -->
      <div v-if="!item.is_past_due" class="text-xs text-gray-500 flex-shrink-0">
        {{ Math.abs(item.days_until_due) }} hari lagi
      </div>
      <div v-else class="text-xs font-semibold text-red-600 flex-shrink-0">
        Terlambat {{ Math.abs(item.days_until_due) }} hari
      </div>
    </div>

    <!-- Footer dengan Prepared By dan Photos -->
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
      
      <!-- Photo Count jika ada -->
      <div 
        v-if="hasPhotos" 
        class="flex items-center gap-1 px-2 py-1 rounded-lg bg-gray-100"
      >
        <Image class="w-3.5 h-3.5 text-gray-500" />
        <span class="text-xs font-medium text-gray-600">
          {{ item.material_photos?.length || 0 }}
        </span>
      </div>
    </div>
  </Motion>
</template>

<script setup>
/**
 * CetakQueueCard - Card component untuk display single PO di queue cetak
 * dengan material ready indicator, prepared by info, dan photo preview
 */
import { computed } from 'vue'
import { Motion } from 'motion-v'
import { Calendar, Package, CheckCircle, Image } from 'lucide-vue-next'
import PriorityBadge from '@/components/common/PriorityBadge.vue'

const props = defineProps({
  /**
   * Queue item data dari API response
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
 * OBC Master data dari relationship (jika ada)
 */
const obcMaster = computed(() => {
  return props.item.obc_master || null
})

/**
 * Product name dengan fallback dari OBC Master atau denormalized field
 */
const productName = computed(() => {
  return obcMaster.value?.material_description || props.item.product_name || 'N/A'
})

/**
 * OBC Number dengan fallback dari OBC Master atau denormalized field
 */
const obcNumber = computed(() => {
  return obcMaster.value?.obc_number || props.item.obc_number || 'N/A'
})

/**
 * Handle card click untuk view detail
 */
const handleClick = () => {
  emit('click', props.item)
  
  // Haptic feedback untuk mobile
  if ('vibrate' in navigator) {
    navigator.vibrate(10)
  }
}

/**
 * Format material ready timestamp untuk display
 */
const formattedMaterialReady = computed(() => {
  if (!props.item.material_ready_at) return 'Siap'
  
  try {
    const date = new Date(props.item.material_ready_at)
    const today = new Date()
    const isToday = date.toDateString() === today.toDateString()
    
    if (isToday) {
      return date.toLocaleTimeString('id-ID', { hour: '2-digit', minute: '2-digit' })
    }
    
    return date.toLocaleDateString('id-ID', { 
      day: 'numeric', 
      month: 'short'
    })
  } catch (e) {
    return 'Siap'
  }
})

/**
 * Format due date untuk display
 */
const formattedDueDate = computed(() => {
  try {
    const date = new Date(props.item.due_date)
    return date.toLocaleDateString('id-ID', { 
      day: 'numeric', 
      month: 'short',
      year: 'numeric'
    })
  } catch (e) {
    return props.item.due_date
  }
})

/**
 * Due date icon classes berdasarkan urgency
 */
const dueDateIconClass = computed(() => {
  if (props.item.is_past_due) {
    return 'bg-red-100 text-red-600'
  }
  if (props.item.days_until_due <= 3) {
    return 'bg-yellow-100 text-yellow-600'
  }
  return 'bg-emerald-100 text-emerald-600'
})

/**
 * Due date text color berdasarkan urgency
 */
const dueDateTextClass = computed(() => {
  if (props.item.is_past_due) {
    return 'text-red-700'
  }
  if (props.item.days_until_due <= 3) {
    return 'text-yellow-700'
  }
  return 'text-emerald-700'
})

/**
 * Get prepared by initials untuk avatar
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

/**
 * Check jika ada material photos
 */
const hasPhotos = computed(() => {
  return props.item.material_photos && props.item.material_photos.length > 0
})
</script>
