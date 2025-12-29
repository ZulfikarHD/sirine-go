<template>
  <div class="animate-pulse">
    <!-- Table Skeleton untuk loading tables -->
    <div v-if="type === 'table'" class="space-y-3">
      <div 
        v-for="i in rows" 
        :key="i" 
        class="flex gap-4"
      >
        <div class="h-12 bg-gray-200 rounded flex-1"></div>
        <div class="h-12 bg-gray-200 rounded w-32"></div>
        <div class="h-12 bg-gray-200 rounded w-24"></div>
      </div>
    </div>

    <!-- Card Skeleton untuk loading cards -->
    <div v-else-if="type === 'card'" class="bg-white rounded-lg border border-gray-200 p-6 space-y-4">
      <div class="h-4 bg-gray-200 rounded w-3/4"></div>
      <div class="h-4 bg-gray-200 rounded w-1/2"></div>
      <div class="h-4 bg-gray-200 rounded w-full"></div>
    </div>

    <!-- Profile Skeleton untuk loading profile -->
    <div v-else-if="type === 'profile'" class="bg-white rounded-lg border border-gray-200 p-6">
      <div class="flex items-center gap-6 mb-6">
        <div class="w-24 h-24 bg-gray-200 rounded-full"></div>
        <div class="flex-1 space-y-3">
          <div class="h-6 bg-gray-200 rounded w-48"></div>
          <div class="h-4 bg-gray-200 rounded w-32"></div>
        </div>
      </div>
      <div class="space-y-3">
        <div class="h-4 bg-gray-200 rounded w-full"></div>
        <div class="h-4 bg-gray-200 rounded w-3/4"></div>
        <div class="h-4 bg-gray-200 rounded w-2/3"></div>
      </div>
    </div>

    <!-- List Skeleton untuk loading lists -->
    <div v-else-if="type === 'list'" class="space-y-3">
      <div 
        v-for="i in rows" 
        :key="i" 
        class="flex items-center gap-4 p-4 bg-white rounded-lg border border-gray-200"
      >
        <div class="w-12 h-12 bg-gray-200 rounded-full"></div>
        <div class="flex-1 space-y-2">
          <div class="h-4 bg-gray-200 rounded w-3/4"></div>
          <div class="h-3 bg-gray-200 rounded w-1/2"></div>
        </div>
      </div>
    </div>

    <!-- Text Skeleton untuk loading text -->
    <div v-else-if="type === 'text'" class="space-y-2">
      <div 
        v-for="i in rows" 
        :key="i" 
        class="h-4 bg-gray-200 rounded"
        :class="i === rows ? 'w-2/3' : 'w-full'"
      ></div>
    </div>

    <!-- Default Skeleton (simple rectangle) -->
    <div v-else class="h-20 bg-gray-200 rounded"></div>
  </div>
</template>

<script setup>
/**
 * LoadingSkeleton component untuk menampilkan placeholder loading
 * dengan animated pulse effect untuk better UX
 */

defineProps({
  /**
   * Tipe skeleton yang akan ditampilkan
   * @values table, card, profile, list, text, default
   */
  type: {
    type: String,
    default: 'default',
    validator: (value) => ['table', 'card', 'profile', 'list', 'text', 'default'].includes(value),
  },
  
  /**
   * Jumlah rows untuk skeleton (applicable untuk table, list, text)
   */
  rows: {
    type: Number,
    default: 3,
  },
})
</script>

<style scoped>
/* 
  Pulse animation sudah included di Tailwind CSS
  .animate-pulse akan automatically apply skeleton loading effect
*/
</style>
