<template>
  <Motion
    :initial="{ opacity: 0, y: 20 }"
    :animate="{ opacity: 1, y: 0 }"
    :transition="{ duration: 0.3 }"
    class="card hover:shadow-lg transition-shadow duration-200"
  >
    <div class="flex justify-between items-start mb-4">
      <h3 class="text-xl font-semibold text-gray-800">{{ example.title }}</h3>
      <span 
        class="px-3 py-1 text-sm rounded-full"
        :class="example.is_active ? 'bg-green-100 text-green-800' : 'bg-gray-100 text-gray-800'"
      >
        {{ example.is_active ? 'Aktif' : 'Tidak Aktif' }}
      </span>
    </div>
    
    <p class="text-gray-600 mb-4">{{ example.content || 'Tidak ada konten' }}</p>
    
    <div class="flex gap-2 justify-end">
      <button 
        @click="$emit('edit', example)"
        class="px-3 py-1.5 text-sm bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors duration-200"
      >
        Edit
      </button>
      <button 
        @click="$emit('delete', example.id)"
        class="px-3 py-1.5 text-sm bg-red-600 text-white rounded-lg hover:bg-red-700 transition-colors duration-200"
      >
        Hapus
      </button>
    </div>
  </Motion>
</template>

<script setup>
import { Motion } from '@motionone/vue'

defineProps({
  example: {
    type: Object,
    required: true
  }
})

defineEmits(['edit', 'delete'])
</script>
