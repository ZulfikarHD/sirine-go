<template>
  <div class="space-y-4">
    <!-- Upload Area -->
    <Motion v-bind="entranceAnimations.fadeUp">
      <div
        v-if="photos.length < maxPhotos"
        class="border-2 border-dashed border-gray-300 rounded-2xl p-8 cursor-pointer transition-colors duration-200 hover:border-indigo-400 hover:bg-indigo-50/50"
        :class="{ 'border-indigo-500 bg-indigo-50': isDragging }"
        @click="triggerFileInput"
        @dragover.prevent="handleDragOver"
        @dragleave.prevent="handleDragLeave"
        @drop.prevent="handleDrop"
      >
        <input
          ref="fileInputRef"
          type="file"
          accept="image/*"
          multiple
          capture="environment"
          class="hidden"
          @change="handleFileSelect"
        />
        
        <div class="flex flex-col items-center justify-center text-center">
          <Camera class="w-10 h-10 text-indigo-500 mb-3" />
          <p class="text-lg font-semibold text-gray-900">
            Tambah Foto Material
          </p>
          <p class="text-sm text-gray-600 mt-1">
            Ketuk untuk kamera atau pilih file
          </p>
          <p class="text-xs text-gray-400 mt-2">
            Maksimal {{ maxPhotos }} foto • Max {{ maxSizeKB / 1024 }}MB per foto
          </p>
        </div>
      </div>
    </Motion>

    <!-- Photo Grid Preview -->
    <div v-if="photos.length > 0" class="grid grid-cols-3 sm:grid-cols-4 gap-3 sm:gap-4">
      <Motion
        v-for="(photo, index) in photos"
        :key="photo.id"
        :initial="{ opacity: 0, scale: 0.8 }"
        :animate="{ opacity: 1, scale: 1 }"
        :transition="{ ...springPresets.snappy, delay: index * 0.05 }"
        class="relative"
      >
        <!-- Photo Thumbnail -->
        <div class="relative aspect-square rounded-xl overflow-hidden bg-gray-100 ring-2 ring-transparent hover:ring-indigo-400 transition-all">
          <img
            :src="photo.preview"
            :alt="`Foto ${index + 1}`"
            class="w-full h-full object-cover cursor-pointer"
            @click="openPreview(index)"
          />
          
          <!-- Delete Button -->
          <button
            type="button"
            class="absolute top-1.5 right-1.5 p-1.5 rounded-full bg-black/60 text-white backdrop-blur-sm hover:bg-red-500 transition-colors active-scale"
            @click.stop="removePhoto(index)"
          >
            <X class="w-4 h-4" />
          </button>

          <!-- Photo Index Badge -->
          <div class="absolute bottom-1.5 left-1.5 px-2 py-0.5 rounded-full bg-black/60 text-white text-xs font-medium backdrop-blur-sm">
            {{ index + 1 }}
          </div>

          <!-- Upload Progress -->
          <div v-if="photo.uploading" class="absolute bottom-0 left-0 right-0 h-1 bg-gray-200">
            <div 
              class="h-full bg-gradient-to-r from-indigo-500 to-fuchsia-500 transition-all" 
              :style="{ width: `${photo.progress}%` }"
            ></div>
          </div>
        </div>
      </Motion>

      <!-- Add More Button (jika masih bisa tambah) -->
      <Motion
        v-if="photos.length < maxPhotos"
        :initial="{ opacity: 0, scale: 0.8 }"
        :animate="{ opacity: 1, scale: 1 }"
        :transition="{ ...springPresets.snappy, delay: photos.length * 0.05 }"
        class="aspect-square rounded-xl border-2 border-dashed border-gray-300 flex flex-col items-center justify-center cursor-pointer hover:border-indigo-400 hover:bg-indigo-50/50 transition-colors active-scale"
        @click="triggerFileInput"
      >
        <Plus class="w-8 h-8 text-gray-400" />
        <span class="text-xs text-gray-500 mt-1">Tambah</span>
      </Motion>
    </div>

    <!-- Photo Count Info -->
    <div v-if="photos.length > 0" class="flex items-center gap-2 text-sm text-gray-500">
      <ImageIcon class="w-4 h-4" />
      <span>{{ photos.length }} / {{ maxPhotos }} foto</span>
    </div>

    <!-- Full Preview Modal -->
    <BaseModal
      v-model="previewModal"
      title="Preview Foto"
      size="lg"
      :show-footer="false"
      no-padding
    >
      <div class="relative bg-black">
        <img
          v-if="previewIndex !== null && photos[previewIndex]"
          :src="photos[previewIndex].preview"
          :alt="`Foto ${previewIndex + 1}`"
          class="w-full max-h-[70vh] object-contain"
        />
        
        <!-- Navigation Buttons -->
        <div v-if="photos.length > 1" class="absolute inset-y-0 left-0 right-0 flex items-center justify-between px-4">
          <button
            type="button"
            class="p-2 rounded-full bg-white/90 text-gray-800 hover:bg-white disabled:opacity-30 disabled:cursor-not-allowed transition-all shadow-lg active-scale"
            :disabled="previewIndex === 0"
            @click="navigatePreview(-1)"
          >
            <ChevronLeft class="w-6 h-6" />
          </button>
          <button
            type="button"
            class="p-2 rounded-full bg-white/90 text-gray-800 hover:bg-white disabled:opacity-30 disabled:cursor-not-allowed transition-all shadow-lg active-scale"
            :disabled="previewIndex === photos.length - 1"
            @click="navigatePreview(1)"
          >
            <ChevronRight class="w-6 h-6" />
          </button>
        </div>

        <!-- Preview Info -->
        <div class="absolute bottom-4 left-1/2 -translate-x-1/2 px-4 py-2 rounded-full bg-black/60 text-white text-sm backdrop-blur-sm">
          {{ previewIndex + 1 }} / {{ photos.length }}
        </div>
      </div>
    </BaseModal>

    <!-- Error Alert -->
    <AlertDialog
      v-model="alertDialog.isOpen.value"
      :title="alertDialog.config.value.title"
      :message="alertDialog.config.value.message"
      :detail="alertDialog.config.value.detail"
      :variant="alertDialog.config.value.variant"
      @close="alertDialog.handleClose()"
    />
  </div>
</template>

<script setup>
/**
 * PhotoUploader Component (Sprint 4)
 * Komponen untuk multiple photo upload dengan preview, compress, dan delete
 * yang mengimplementasikan iOS-inspired design dengan Motion-V animations
 */
import { ref, computed } from 'vue'
import { Motion } from 'motion-v'
import { entranceAnimations, springPresets } from '@/composables/useMotion'
import { useAlertDialog } from '@/composables/useModal'
import BaseModal from './BaseModal.vue'
import AlertDialog from './AlertDialog.vue'
import { 
  Camera, 
  X, 
  Plus, 
  Image as ImageIcon,
  ChevronLeft,
  ChevronRight
} from 'lucide-vue-next'

// Props definition
const props = defineProps({
  /**
   * Array of photo objects with preview URLs
   */
  modelValue: {
    type: Array,
    default: () => []
  },
  /**
   * Maximum number of photos allowed
   */
  maxPhotos: {
    type: Number,
    default: 5
  },
  /**
   * Maximum file size per photo in KB
   */
  maxSizeKB: {
    type: Number,
    default: 2048 // 2MB
  }
})

// Emits definition
const emit = defineEmits([
  'update:modelValue',
  'upload-start',
  'upload-complete',
  'upload-error'
])

// Composables
const alertDialog = useAlertDialog()

// State
const fileInputRef = ref(null)
const isDragging = ref(false)
const previewModal = ref(false)
const previewIndex = ref(null)

// Computed: photos array dengan two-way binding
const photos = computed({
  get: () => props.modelValue,
  set: (value) => emit('update:modelValue', value)
})

/**
 * Trigger file input click
 */
const triggerFileInput = () => {
  fileInputRef.value?.click()
}

/**
 * Handle file selection dari input
 */
const handleFileSelect = async (event) => {
  const files = Array.from(event.target.files || [])
  await processFiles(files)
  
  // Reset input untuk allow re-select same file
  if (fileInputRef.value) {
    fileInputRef.value.value = ''
  }
}

/**
 * Handle drag over event
 */
const handleDragOver = () => {
  isDragging.value = true
}

/**
 * Handle drag leave event
 */
const handleDragLeave = () => {
  isDragging.value = false
}

/**
 * Handle drop event untuk drag & drop upload
 */
const handleDrop = async (event) => {
  isDragging.value = false
  const files = Array.from(event.dataTransfer?.files || [])
  await processFiles(files)
}

/**
 * Process uploaded files dengan validation dan compression
 */
const processFiles = async (files) => {
  // Filter hanya image files
  const imageFiles = files.filter(file => file.type.startsWith('image/'))
  
  if (imageFiles.length === 0) {
    await alertDialog.error('File Tidak Valid', {
      title: 'Upload Error',
      detail: 'Hanya file gambar (JPG, PNG, GIF, WEBP) yang diperbolehkan.'
    })
    return
  }

  // Check remaining slots
  const remainingSlots = props.maxPhotos - photos.value.length
  const filesToProcess = imageFiles.slice(0, remainingSlots)

  if (imageFiles.length > remainingSlots) {
    await alertDialog.warning('Batas Foto Tercapai', {
      title: 'Peringatan',
      detail: `Hanya ${remainingSlots} foto lagi yang dapat ditambahkan. ${imageFiles.length - remainingSlots} foto diabaikan.`
    })
  }

  emit('upload-start')

  // Process each file
  for (const file of filesToProcess) {
    try {
      const processedPhoto = await processImage(file)
      photos.value = [...photos.value, processedPhoto]
    } catch (error) {
      console.error('Error processing image:', error)
      emit('upload-error', { file, error })
    }
  }

  emit('upload-complete', photos.value)
  
  // Haptic feedback untuk success
  if ('vibrate' in navigator) {
    navigator.vibrate(10)
  }
}

/**
 * Process single image dengan compression
 */
const processImage = (file) => {
  return new Promise((resolve, reject) => {
    const reader = new FileReader()
    
    reader.onload = async (e) => {
      try {
        const originalBase64 = e.target.result
        
        // Check file size
        const sizeKB = file.size / 1024
        let finalBase64 = originalBase64

        // Compress if needed
        if (sizeKB > props.maxSizeKB) {
          finalBase64 = await compressImage(originalBase64, props.maxSizeKB)
        }

        resolve({
          id: `photo_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`,
          preview: finalBase64,
          base64: finalBase64.split(',')[1], // Remove data:image prefix untuk upload
          fileName: file.name,
          fileType: file.type,
          originalSize: sizeKB,
          uploading: false,
          progress: 100
        })
      } catch (error) {
        reject(error)
      }
    }

    reader.onerror = () => {
      reject(new Error('Gagal membaca file'))
    }

    reader.readAsDataURL(file)
  })
}

/**
 * Compress image to target size
 */
const compressImage = (base64, maxSizeKB) => {
  return new Promise((resolve) => {
    const img = new Image()
    img.src = base64

    img.onload = () => {
      const canvas = document.createElement('canvas')
      const ctx = canvas.getContext('2d')

      // Calculate new dimensions (max 1920px width)
      let width = img.width
      let height = img.height
      const maxWidth = 1920

      if (width > maxWidth) {
        height = (height * maxWidth) / width
        width = maxWidth
      }

      canvas.width = width
      canvas.height = height

      // Draw image
      ctx.drawImage(img, 0, 0, width, height)

      // Try different quality levels
      let quality = 0.8
      let result = canvas.toDataURL('image/jpeg', quality)
      
      // Reduce quality until under max size
      while (result.length / 1024 > maxSizeKB && quality > 0.1) {
        quality -= 0.1
        result = canvas.toDataURL('image/jpeg', quality)
      }

      resolve(result)
    }
  })
}

/**
 * Remove photo by index
 */
const removePhoto = (index) => {
  photos.value = photos.value.filter((_, i) => i !== index)
  
  // Haptic feedback
  if ('vibrate' in navigator) {
    navigator.vibrate(10)
  }
}

/**
 * Open preview modal
 */
const openPreview = (index) => {
  previewIndex.value = index
  previewModal.value = true
}

/**
 * Navigate preview carousel
 */
const navigatePreview = (direction) => {
  const newIndex = previewIndex.value + direction
  if (newIndex >= 0 && newIndex < photos.value.length) {
    previewIndex.value = newIndex
  }
}

/**
 * Get base64 array untuk API submission
 */
const getBase64Array = () => {
  return photos.value.map(photo => photo.base64)
}

// Expose methods untuk parent component
defineExpose({
  getBase64Array,
  triggerFileInput
})
</script>
