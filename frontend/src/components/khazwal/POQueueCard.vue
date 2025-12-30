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
          PO #{{ poItem.po_number }}
        </h3>
        <p class="text-sm text-indigo-600 font-semibold truncate">
          OBC: {{ obcNumber }}
        </p>
      </div>
      
      <PriorityBadge :priority="poItem.priority" size="sm" />
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
        <span v-if="obcMaster.personalization" class="px-2 py-0.5 bg-purple-50 rounded-md text-purple-700 border border-purple-200">
          {{ obcMaster.personalization }}
        </span>
      </div>
    </div>

    <!-- Info Grid -->
    <div class="grid grid-cols-2 gap-3 mb-3">
      <!-- Due Date -->
      <div class="flex items-center gap-2">
        <div :class="dueDateIconClass" class="p-1.5 rounded-lg">
          <Calendar class="w-4 h-4" />
        </div>
        <div class="min-w-0 flex-1">
          <p class="text-xs text-gray-500">Deadline</p>
          <p :class="dueDateTextClass" class="text-sm font-semibold truncate">
            {{ formattedDueDate }}
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
            {{ poItem.quantity_ordered.toLocaleString() }}
          </p>
        </div>
      </div>
    </div>

    <!-- Status Badge -->
    <div class="flex items-center justify-between pt-3 border-t border-gray-200/50">
      <div :class="statusBadgeClass" class="px-3 py-1 rounded-lg text-xs font-semibold">
        {{ statusLabel }}
      </div>
      
      <!-- Days indicator -->
      <div v-if="!poItem.is_past_due" class="text-xs text-gray-500">
        {{ Math.abs(poItem.days_until_due) }} hari lagi
      </div>
      <div v-else class="text-xs font-semibold text-red-600">
        Terlambat {{ Math.abs(poItem.days_until_due) }} hari
      </div>
    </div>
  </Motion>
</template>

<script setup>
/**
 * POQueueCard - Card component untuk display single PO di queue list
 * dengan priority indicator, due date, dan status information
 */
import { computed } from 'vue'
import { Motion } from 'motion-v'
import { Calendar, Package } from 'lucide-vue-next'
import PriorityBadge from '@/components/common/PriorityBadge.vue'

const props = defineProps({
  /**
   * PO item data dari API queue response
   */
  poItem: {
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
 * dengan fallback ke denormalized fields untuk backward compatibility
 */
const obcMaster = computed(() => {
  return props.poItem.obc_master || null
})

/**
 * Product name dengan fallback dari OBC Master atau denormalized field
 */
const productName = computed(() => {
  return obcMaster.value?.material_description || props.poItem.product_name || 'N/A'
})

/**
 * OBC Number dengan fallback dari OBC Master atau denormalized field
 */
const obcNumber = computed(() => {
  return obcMaster.value?.obc_number || props.poItem.obc_number || 'N/A'
})

/**
 * Handle card click untuk navigate ke detail page
 */
const handleClick = () => {
  emit('click', props.poItem)
  
  // Haptic feedback untuk mobile
  if ('vibrate' in navigator) {
    navigator.vibrate(10)
  }
}

/**
 * Format due date untuk display yang user-friendly
 */
const formattedDueDate = computed(() => {
  try {
    const date = new Date(props.poItem.due_date)
    return date.toLocaleDateString('id-ID', { 
      day: 'numeric', 
      month: 'short',
      year: 'numeric'
    })
  } catch (e) {
    return props.poItem.due_date
  }
})

/**
 * Due date icon classes berdasarkan urgency
 * dengan color coding untuk visual indicator
 */
const dueDateIconClass = computed(() => {
  if (props.poItem.is_past_due) {
    return 'bg-red-100 text-red-600'
  }
  if (props.poItem.days_until_due <= 3) {
    return 'bg-yellow-100 text-yellow-600'
  }
  return 'bg-emerald-100 text-emerald-600'
})

/**
 * Due date text color berdasarkan urgency
 */
const dueDateTextClass = computed(() => {
  if (props.poItem.is_past_due) {
    return 'text-red-700'
  }
  if (props.poItem.days_until_due <= 3) {
    return 'text-yellow-700'
  }
  return 'text-emerald-700'
})

/**
 * Status badge styling berdasarkan current status
 */
const statusBadgeClass = computed(() => {
  const statusMap = {
    'WAITING_MATERIAL_PREP': 'bg-blue-50 text-blue-700 border border-blue-200/50',
    'MATERIAL_PREP_IN_PROGRESS': 'bg-yellow-50 text-yellow-700 border border-yellow-200/50',
    'READY_FOR_CETAK': 'bg-green-50 text-green-700 border border-green-200/50'
  }
  return statusMap[props.poItem.current_status] || 'bg-gray-50 text-gray-700 border border-gray-200/50'
})

/**
 * Human-readable status label
 */
const statusLabel = computed(() => {
  const labelMap = {
    'WAITING_MATERIAL_PREP': 'Menunggu Persiapan',
    'MATERIAL_PREP_IN_PROGRESS': 'Sedang Dipersiapkan',
    'READY_FOR_CETAK': 'Siap Cetak'
  }
  return labelMap[props.poItem.current_status] || props.poItem.current_status
})
</script>
