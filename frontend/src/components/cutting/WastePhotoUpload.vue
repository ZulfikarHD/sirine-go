<template>
  <div>
    <label class="block text-sm font-semibold text-gray-700 mb-2">
      <svg class="w-4 h-4 inline mr-1" fill="currentColor" viewBox="0 0 20 20">
        <path fill-rule="evenodd" d="M4 3a2 2 0 00-2 2v10a2 2 0 002 2h12a2 2 0 002-2V5a2 2 0 00-2-2H4zm12 12H4l4-8 3 6 2-4 3 6z" clip-rule="evenodd"/>
      </svg>
      Foto Waste <span v-if="required" class="text-red-500">*</span>
    </label>

    <!-- Upload Area (if no photo yet) -->
    <div
      v-if="!previewUrl"
      @click="triggerFileInput"
      @dragover.prevent="isDragging = true"
      @dragleave.prevent="isDragging = false"
      @drop.prevent="handleDrop"
      class="relative border-2 border-dashed rounded-xl p-8 text-center cursor-pointer transition-all"
      :class="[
        isDragging ? 'border-indigo-500 bg-indigo-50' : 'border-gray-300 hover:border-indigo-400 hover:bg-gray-50',
        disabled ? 'opacity-50 cursor-not-allowed' : ''
      ]"
    >
      <input
        ref="fileInput"
        type="file"
        accept="image/*"
        capture="environment"
        @change="handleFileSelect"
        :disabled="disabled"
        class="hidden"
      />

      <div class="space-y-3">
        <!-- Icon -->
        <div class="w-16 h-16 mx-auto rounded-full bg-gradient-to-br from-indigo-100 to-fuchsia-100 flex items-center justify-center">
          <svg class="w-8 h-8 text-indigo-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 9a2 2 0 012-2h.93a2 2 0 001.664-.89l.812-1.22A2 2 0 0110.07 4h3.86a2 2 0 011.664.89l.812 1.22A2 2 0 0018.07 7H19a2 2 0 012 2v9a2 2 0 01-2 2H5a2 2 0 01-2-2V9z"/>
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 13a3 3 0 11-6 0 3 3 0 016 0z"/>
          </svg>
        </div>

        <!-- Text -->
        <div>
          <p class="text-sm font-semibold text-gray-700 mb-1">
            {{ isMobile ? 'Ambil Foto atau Pilih dari Galeri' : 'Klik untuk Upload atau Drag & Drop' }}
          </p>
          <p class="text-xs text-gray-500">
            Format: JPG, PNG, WEBP (Max {{ maxSizeMB }}MB)
          </p>
        </div>

        <!-- Mobile Actions -->
        <div v-if="isMobile" class="flex gap-2 justify-center pt-2">
          <button
            type="button"
            @click.stop="openCamera"
            :disabled="disabled"
            class="px-4 py-2 bg-gradient-to-r from-indigo-600 to-fuchsia-600 text-white text-sm font-semibold rounded-lg hover:from-indigo-700 hover:to-fuchsia-700 transition-all active-scale"
          >
            <svg class="w-4 h-4 inline mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 9a2 2 0 012-2h.93a2 2 0 001.664-.89l.812-1.22A2 2 0 0110.07 4h3.86a2 2 0 011.664.89l.812 1.22A2 2 0 0018.07 7H19a2 2 0 012 2v9a2 2 0 01-2 2H5a2 2 0 01-2-2V9z"/>
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 13a3 3 0 11-6 0 3 3 0 016 0z"/>
            </svg>
            Kamera
          </button>
          <button
            type="button"
            @click.stop="openGallery"
            :disabled="disabled"
            class="px-4 py-2 bg-white border-2 border-gray-300 text-gray-700 text-sm font-semibold rounded-lg hover:bg-gray-50 transition-all active-scale"
          >
            <svg class="w-4 h-4 inline mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"/>
            </svg>
            Galeri
          </button>
        </div>
      </div>

      <!-- Uploading Indicator -->
      <div v-if="isUploading" class="absolute inset-0 bg-white/90 rounded-xl flex items-center justify-center">
        <div class="text-center">
          <svg class="animate-spin w-12 h-12 text-indigo-600 mx-auto mb-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
          </svg>
          <p class="text-sm font-semibold text-gray-700">Uploading...</p>
          <p class="text-xs text-gray-500">{{ uploadProgress }}%</p>
        </div>
      </div>
    </div>

    <!-- Preview Area (if photo uploaded) -->
    <div v-else class="relative">
      <div class="border-2 border-gray-300 rounded-xl overflow-hidden">
        <img
          :src="previewUrl"
          alt="Waste photo preview"
          class="w-full h-64 object-cover"
        />
      </div>

      <!-- Remove Button -->
      <button
        type="button"
        @click="handleRemove"
        :disabled="disabled"
        class="absolute top-2 right-2 w-10 h-10 bg-red-500 hover:bg-red-600 text-white rounded-full flex items-center justify-center transition-all active-scale shadow-lg"
        :class="{ 'opacity-50 cursor-not-allowed': disabled }"
      >
        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
        </svg>
      </button>

      <!-- Photo Info -->
      <div class="mt-2 text-xs text-gray-500 flex items-center gap-2">
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/>
        </svg>
        <span>Foto berhasil diupload</span>
      </div>
    </div>

    <!-- Error Message -->
    <div v-if="errorMessage" class="mt-2 text-sm text-red-600 flex items-center gap-2">
      <svg class="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
        <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd"/>
      </svg>
      {{ errorMessage }}
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { useApi } from '@/composables/useApi'
import { useHaptic } from '@/composables/useHaptic'

const props = defineProps({
  modelValue: {
    type: String,
    default: ''
  },
  required: {
    type: Boolean,
    default: false
  },
  disabled: {
    type: Boolean,
    default: false
  },
  maxSize: {
    type: Number,
    default: 2 * 1024 * 1024 // 2MB in bytes
  }
})

const emit = defineEmits(['update:modelValue', 'upload', 'remove'])

const api = useApi()
const haptic = useHaptic()

// Refs
const fileInput = ref(null)
const isDragging = ref(false)
const isUploading = ref(false)
const uploadProgress = ref(0)
const previewUrl = ref('')
const errorMessage = ref('')

// Computed
const maxSizeMB = computed(() => {
  return Math.round(props.maxSize / (1024 * 1024))
})

const isMobile = computed(() => {
  return /Android|webOS|iPhone|iPad|iPod|BlackBerry|IEMobile|Opera Mini/i.test(navigator.userAgent)
})

// Watch modelValue to update preview
watch(() => props.modelValue, (newValue) => {
  if (newValue && newValue.startsWith('http')) {
    previewUrl.value = newValue
  } else if (newValue && newValue.startsWith('/uploads')) {
    previewUrl.value = `${import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080'}${newValue}`
  } else if (!newValue) {
    previewUrl.value = ''
  }
}, { immediate: true })

/**
 * Trigger file input click
 */
const triggerFileInput = () => {
  if (!props.disabled) {
    fileInput.value?.click()
  }
}

/**
 * Open camera (for mobile)
 */
const openCamera = () => {
  if (!props.disabled) {
    const input = document.createElement('input')
    input.type = 'file'
    input.accept = 'image/*'
    input.capture = 'environment'
    input.onchange = (e) => {
      const file = e.target.files[0]
      if (file) {
        processFile(file)
      }
    }
    input.click()
  }
}

/**
 * Open gallery (for mobile)
 */
const openGallery = () => {
  if (!props.disabled) {
    const input = document.createElement('input')
    input.type = 'file'
    input.accept = 'image/*'
    input.onchange = (e) => {
      const file = e.target.files[0]
      if (file) {
        processFile(file)
      }
    }
    input.click()
  }
}

/**
 * Handle file select from input
 */
const handleFileSelect = (event) => {
  const file = event.target.files[0]
  if (file) {
    processFile(file)
  }
}

/**
 * Handle drag and drop
 */
const handleDrop = (event) => {
  isDragging.value = false
  
  if (props.disabled) return
  
  const file = event.dataTransfer.files[0]
  if (file) {
    processFile(file)
  }
}

/**
 * Process file - validate and upload
 */
const processFile = async (file) => {
  errorMessage.value = ''
  
  // Validate file type
  if (!file.type.startsWith('image/')) {
    errorMessage.value = 'File harus berupa gambar (JPG, PNG, WEBP)'
    haptic.error()
    return
  }
  
  // Validate file size
  if (file.size > props.maxSize) {
    errorMessage.value = `Ukuran file maksimal ${maxSizeMB.value}MB`
    haptic.error()
    return
  }
  
  // Create preview
  const reader = new FileReader()
  reader.onload = (e) => {
    previewUrl.value = e.target.result
  }
  reader.readAsDataURL(file)
  
  // Upload file
  await uploadFile(file)
}

/**
 * Upload file to server
 */
const uploadFile = async (file) => {
  isUploading.value = true
  uploadProgress.value = 0
  
  try {
    const formData = new FormData()
    formData.append('photo', file)
    
    // Upload to profile photo endpoint (or create dedicated waste photo endpoint)
    // For now, using profile photo endpoint as it handles image uploads
    const response = await api.post('/profile/photo', formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      },
      onUploadProgress: (progressEvent) => {
        uploadProgress.value = Math.round((progressEvent.loaded * 100) / progressEvent.total)
      }
    })
    
    // Get photo URL from response
    const photoUrl = response.photo_url || response.url || ''
    
    // Emit update
    emit('update:modelValue', photoUrl)
    emit('upload', photoUrl)
    
    // Success haptic
    haptic.success()
    
  } catch (error) {
    console.error('Failed to upload photo:', error)
    errorMessage.value = 'Gagal mengupload foto. Silakan coba lagi.'
    previewUrl.value = ''
    haptic.error()
  } finally {
    isUploading.value = false
    uploadProgress.value = 0
  }
}

/**
 * Handle remove photo
 */
const handleRemove = () => {
  if (props.disabled) return
  
  previewUrl.value = ''
  errorMessage.value = ''
  
  // Emit update
  emit('update:modelValue', '')
  emit('remove')
  
  // Light haptic
  haptic.light()
  
  // Reset file input
  if (fileInput.value) {
    fileInput.value.value = ''
  }
}
</script>
