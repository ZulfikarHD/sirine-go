<template>
  <!-- Table Skeleton -->
  <div v-if="type === 'table'" class="space-y-3">
    <div v-for="i in rows" :key="i" class="flex items-center gap-4 p-4 bg-gray-50 rounded-xl animate-pulse">
      <div class="w-12 h-12 bg-gray-200 rounded-full"></div>
      <div class="flex-1 space-y-2">
        <div class="h-4 bg-gray-200 rounded w-1/3"></div>
        <div class="h-3 bg-gray-200 rounded w-1/2"></div>
      </div>
      <div class="flex gap-2">
        <div class="w-16 h-8 bg-gray-200 rounded"></div>
        <div class="w-16 h-8 bg-gray-200 rounded"></div>
      </div>
    </div>
  </div>

  <!-- Card Skeleton -->
  <div v-else-if="type === 'card'" class="glass-card p-6 rounded-2xl animate-pulse">
    <div class="space-y-4">
      <div class="h-6 bg-gray-200 rounded w-1/3"></div>
      <div class="space-y-2">
        <div class="h-4 bg-gray-200 rounded"></div>
        <div class="h-4 bg-gray-200 rounded w-5/6"></div>
        <div class="h-4 bg-gray-200 rounded w-4/6"></div>
      </div>
    </div>
  </div>

  <!-- List Skeleton -->
  <div v-else-if="type === 'list'" class="space-y-2">
    <div v-for="i in rows" :key="i" class="flex items-center gap-3 p-3 bg-gray-50 rounded-xl animate-pulse">
      <div class="w-8 h-8 bg-gray-200 rounded-full"></div>
      <div class="flex-1 space-y-2">
        <div class="h-3 bg-gray-200 rounded w-2/3"></div>
        <div class="h-2 bg-gray-200 rounded w-1/2"></div>
      </div>
    </div>
  </div>

  <!-- Grid Skeleton (for achievements, etc) -->
  <div v-else-if="type === 'grid'" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
    <div v-for="i in rows" :key="i" class="glass-card p-4 rounded-xl animate-pulse">
      <div class="flex items-start gap-3">
        <div class="w-12 h-12 bg-gray-200 rounded-full"></div>
        <div class="flex-1 space-y-2">
          <div class="h-4 bg-gray-200 rounded w-3/4"></div>
          <div class="h-3 bg-gray-200 rounded"></div>
          <div class="h-2 bg-gray-200 rounded w-1/2"></div>
        </div>
      </div>
    </div>
  </div>

  <!-- Profile Skeleton -->
  <div v-else-if="type === 'profile'" class="space-y-6">
    <div class="glass-card p-8 rounded-2xl animate-pulse">
      <div class="flex items-center gap-6">
        <div class="w-24 h-24 bg-gray-200 rounded-full"></div>
        <div class="flex-1 space-y-3">
          <div class="h-8 bg-gray-200 rounded w-1/3"></div>
          <div class="flex gap-2">
            <div class="w-20 h-6 bg-gray-200 rounded-full"></div>
            <div class="w-20 h-6 bg-gray-200 rounded-full"></div>
            <div class="w-20 h-6 bg-gray-200 rounded-full"></div>
          </div>
        </div>
      </div>
    </div>
    <div class="glass-card p-6 rounded-2xl animate-pulse">
      <div class="space-y-4">
        <div class="h-6 bg-gray-200 rounded w-1/4"></div>
        <div class="grid grid-cols-2 gap-4">
          <div class="h-20 bg-gray-100 rounded-xl"></div>
          <div class="h-20 bg-gray-100 rounded-xl"></div>
          <div class="h-20 bg-gray-100 rounded-xl"></div>
          <div class="h-20 bg-gray-100 rounded-xl"></div>
        </div>
      </div>
    </div>
  </div>

  <!-- Default Simple Skeleton -->
  <div v-else class="space-y-3 animate-pulse">
    <div v-for="i in rows" :key="i" class="h-12 bg-gray-200 rounded"></div>
  </div>
</template>

<script setup>
defineProps({
  type: {
    type: String,
    default: 'default',
    validator: (value) => ['default', 'table', 'card', 'list', 'grid', 'profile'].includes(value)
  },
  rows: {
    type: Number,
    default: 3
  }
})
</script>
