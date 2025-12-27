<template>
  <div class="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100 py-8 px-4">
    <div class="max-w-6xl mx-auto">
      <!-- Header -->
      <Motion
        :initial="{ opacity: 0, y: -20 }"
        :animate="{ opacity: 1, y: 0 }"
        :transition="{ duration: 0.5 }"
        class="text-center mb-8"
      >
        <h1 class="text-4xl md:text-5xl font-bold text-gray-800 mb-2">
          Sirine Go App
        </h1>
        <p class="text-gray-600">
          Web app offline dengan Gin, Vue, dan MySQL
        </p>
        
        <!-- Status Offline/Online -->
        <div class="mt-4 inline-flex items-center gap-2 px-4 py-2 rounded-full bg-white shadow-sm">
          <div 
            class="w-3 h-3 rounded-full"
            :class="isOnline ? 'bg-green-500' : 'bg-red-500'"
          ></div>
          <span class="text-sm font-medium">
            {{ isOnline ? 'Online' : 'Offline' }}
          </span>
        </div>
      </Motion>

      <!-- Add Button -->
      <Motion
        :initial="{ opacity: 0 }"
        :animate="{ opacity: 1 }"
        :transition="{ duration: 0.5, delay: 0.1 }"
        class="mb-6"
      >
        <button
          v-if="!showForm"
          @click="showForm = true"
          class="btn-primary w-full md:w-auto"
        >
          + Tambah Data Baru
        </button>
      </Motion>

      <!-- Form -->
      <div v-if="showForm" class="mb-8">
        <ExampleForm
          :initial-data="editingExample"
          :loading="loading"
          @submit="handleSubmit"
          @cancel="cancelForm"
        />
      </div>

      <!-- Loading State -->
      <div v-if="loading && !showForm" class="text-center py-12">
        <div class="inline-block animate-spin rounded-full h-12 w-12 border-4 border-blue-600 border-t-transparent"></div>
        <p class="mt-4 text-gray-600">Memuat data...</p>
      </div>

      <!-- Error State -->
      <Motion
        v-if="error && !loading"
        :initial="{ opacity: 0, scale: 0.9 }"
        :animate="{ opacity: 1, scale: 1 }"
        class="card bg-red-50 border-2 border-red-200 mb-6"
      >
        <p class="text-red-800">{{ error }}</p>
      </Motion>

      <!-- Empty State -->
      <Motion
        v-if="!loading && examples.length === 0 && !showForm"
        :initial="{ opacity: 0 }"
        :animate="{ opacity: 1 }"
        :transition="{ duration: 0.5 }"
        class="text-center py-12"
      >
        <div class="text-6xl mb-4">📝</div>
        <h3 class="text-xl font-semibold text-gray-700 mb-2">
          Belum ada data
        </h3>
        <p class="text-gray-500">
          Klik tombol "Tambah Data Baru" untuk memulai
        </p>
      </Motion>

      <!-- Examples Grid -->
      <div 
        v-if="!loading && examples.length > 0"
        class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6"
      >
        <ExampleCard
          v-for="example in examples"
          :key="example.id"
          :example="example"
          @edit="handleEdit"
          @delete="handleDelete"
        />
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { Motion } from '@motionone/vue'
import { useOnline } from '@vueuse/core'
import { useExamples } from '../composables/useExamples'
import ExampleCard from '../components/ExampleCard.vue'
import ExampleForm from '../components/ExampleForm.vue'

const isOnline = useOnline()
const showForm = ref(false)
const editingExample = ref(null)

const {
  examples,
  loading,
  error,
  fetchExamples,
  createExample,
  updateExample,
  deleteExample
} = useExamples()

onMounted(() => {
  fetchExamples()
})

const handleSubmit = async (formData) => {
  let success = false
  
  if (editingExample.value) {
    success = await updateExample(editingExample.value.id, formData)
  } else {
    success = await createExample(formData)
  }
  
  if (success) {
    cancelForm()
  }
}

const handleEdit = (example) => {
  editingExample.value = { ...example }
  showForm.value = true
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

const handleDelete = async (id) => {
  if (confirm('Apakah Anda yakin ingin menghapus data ini?')) {
    await deleteExample(id)
  }
}

const cancelForm = () => {
  showForm.value = false
  editingExample.value = null
}
</script>
