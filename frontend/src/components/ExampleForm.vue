<template>
  <Motion
    :initial="{ opacity: 0, scale: 0.95 }"
    :animate="{ opacity: 1, scale: 1 }"
    :transition="{ duration: 0.2 }"
    class="card"
  >
    <h2 class="text-2xl font-bold mb-6 text-gray-800">
      {{ isEdit ? 'Edit Data' : 'Tambah Data Baru' }}
    </h2>
    
    <form @submit.prevent="handleSubmit" class="space-y-4">
      <div>
        <label for="title" class="block text-sm font-medium text-gray-700 mb-2">
          Judul <span class="text-red-500">*</span>
        </label>
        <input
          id="title"
          v-model="form.title"
          type="text"
          required
          class="input"
          placeholder="Masukkan judul"
        />
      </div>
      
      <div>
        <label for="content" class="block text-sm font-medium text-gray-700 mb-2">
          Konten
        </label>
        <textarea
          id="content"
          v-model="form.content"
          rows="4"
          class="input resize-none"
          placeholder="Masukkan konten"
        ></textarea>
      </div>
      
      <div class="flex items-center">
        <input
          id="is_active"
          v-model="form.is_active"
          type="checkbox"
          class="w-4 h-4 text-blue-600 border-gray-300 rounded focus:ring-blue-500"
        />
        <label for="is_active" class="ml-2 text-sm font-medium text-gray-700">
          Aktif
        </label>
      </div>
      
      <div class="flex gap-3 pt-4">
        <button
          type="submit"
          :disabled="loading"
          class="btn-primary flex-1 disabled:opacity-50 disabled:cursor-not-allowed"
        >
          {{ loading ? 'Menyimpan...' : (isEdit ? 'Perbarui' : 'Simpan') }}
        </button>
        <button
          type="button"
          @click="$emit('cancel')"
          class="btn-secondary"
        >
          Batal
        </button>
      </div>
    </form>
  </Motion>
</template>

<script setup>
import { ref, watch } from 'vue'
import { Motion } from 'motion-v'

const props = defineProps({
  initialData: {
    type: Object,
    default: null
  },
  loading: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['submit', 'cancel'])

const isEdit = ref(!!props.initialData)

const form = ref({
  title: '',
  content: '',
  is_active: true
})

watch(() => props.initialData, (newData) => {
  if (newData) {
    form.value = { ...newData }
    isEdit.value = true
  } else {
    form.value = {
      title: '',
      content: '',
      is_active: true
    }
    isEdit.value = false
  }
}, { immediate: true })

const handleSubmit = () => {
  emit('submit', form.value)
}
</script>
