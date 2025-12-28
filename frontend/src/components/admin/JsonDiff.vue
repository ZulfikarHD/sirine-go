<template>
  <div v-if="changes" class="space-y-3">
    <div class="text-xs font-semibold text-gray-700 uppercase tracking-wider mb-2">
      Perubahan Data
    </div>

    <!-- Before/After Comparison -->
    <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
      <!-- Before -->
      <div v-if="changes.before" class="rounded-lg border border-red-200 bg-red-50/50 p-3">
        <div class="flex items-center gap-2 mb-2">
          <svg class="w-4 h-4 text-red-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 12H4" />
          </svg>
          <span class="text-xs font-semibold text-red-700 uppercase">Sebelum</span>
        </div>
        <div class="bg-white rounded-md p-2 border border-red-100">
          <pre class="text-xs text-gray-800 overflow-x-auto whitespace-pre-wrap break-words">{{ formatJson(changes.before) }}</pre>
        </div>
      </div>

      <!-- After -->
      <div v-if="changes.after" class="rounded-lg border border-green-200 bg-green-50/50 p-3">
        <div class="flex items-center gap-2 mb-2">
          <svg class="w-4 h-4 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
          </svg>
          <span class="text-xs font-semibold text-green-700 uppercase">Sesudah</span>
        </div>
        <div class="bg-white rounded-md p-2 border border-green-100">
          <pre class="text-xs text-gray-800 overflow-x-auto whitespace-pre-wrap break-words">{{ formatJson(changes.after) }}</pre>
        </div>
      </div>
    </div>

    <!-- Field-by-Field Comparison untuk objects -->
    <div v-if="isObject(changes.before) && isObject(changes.after)" class="mt-4">
      <div class="text-xs font-semibold text-gray-700 uppercase tracking-wider mb-2">
        Detail Perubahan Field
      </div>
      <div class="space-y-2">
        <Motion
          v-for="(field, index) in changedFields"
          :key="field"
          :initial="{ opacity: 0, x: -10 }"
          :animate="{ opacity: 1, x: 0 }"
          :transition="{ duration: 0.2, delay: index * 0.05, ease: 'easeOut' }"
          class="rounded-lg border border-gray-200 bg-white p-3"
        >
          <div class="font-medium text-sm text-gray-900 mb-1">{{ field }}</div>
          <div class="grid grid-cols-2 gap-3 text-xs">
            <div>
              <span class="text-red-600 font-medium">-</span>
              <span class="text-gray-700 ml-1">{{ formatValue(changes.before[field]) }}</span>
            </div>
            <div>
              <span class="text-green-600 font-medium">+</span>
              <span class="text-gray-700 ml-1">{{ formatValue(changes.after[field]) }}</span>
            </div>
          </div>
        </Motion>
      </div>
    </div>
  </div>

  <!-- Empty State -->
  <div v-else class="text-center py-4">
    <p class="text-xs text-gray-500">Tidak ada data perubahan</p>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { Motion } from 'motion-v'

// Props
const props = defineProps({
  changes: {
    type: Object,
    default: null,
  },
})

/**
 * isObject mengecek apakah value adalah object
 */
const isObject = (value) => {
  return value !== null && typeof value === 'object' && !Array.isArray(value)
}

/**
 * changedFields mendapatkan list field yang berubah
 * dengan comparison antara before dan after
 */
const changedFields = computed(() => {
  if (!props.changes?.before || !props.changes?.after) return []
  if (!isObject(props.changes.before) || !isObject(props.changes.after)) return []

  const before = props.changes.before
  const after = props.changes.after
  const allFields = new Set([...Object.keys(before), ...Object.keys(after)])
  
  return Array.from(allFields).filter(field => {
    // Skip fields yang sama atau fields internal (dimulai dengan _)
    if (field.startsWith('_')) return false
    return JSON.stringify(before[field]) !== JSON.stringify(after[field])
  })
})

/**
 * formatJson format object menjadi pretty JSON string
 * dengan indentation untuk readability
 */
const formatJson = (value) => {
  if (value === null) return 'null'
  if (value === undefined) return 'undefined'
  
  try {
    if (typeof value === 'object') {
      return JSON.stringify(value, null, 2)
    }
    return String(value)
  } catch (error) {
    return String(value)
  }
}

/**
 * formatValue format single value untuk display
 * dengan handling special cases
 */
const formatValue = (value) => {
  if (value === null) return '(kosong)'
  if (value === undefined) return '(tidak diset)'
  if (value === '') return '(string kosong)'
  if (typeof value === 'boolean') return value ? 'Ya' : 'Tidak'
  if (typeof value === 'object') return JSON.stringify(value)
  
  // Format password fields
  if (String(value).length > 50 && String(value).includes('$2')) {
    return '••••••••' // Mask password hashes
  }
  
  return String(value)
}
</script>
