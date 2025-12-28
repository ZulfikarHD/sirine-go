<template>
  <div class="flex flex-col items-center gap-3">
    <!-- Photo Display dengan hover overlay -->
    <Motion v-bind="entranceAnimations.fadeScale" class="relative">
      <div 
        class="relative w-32 h-32 sm:w-28 sm:h-28 rounded-full overflow-hidden cursor-pointer active-scale bg-gray-100 border-4 border-white shadow-lg hover:shadow-xl transition-shadow duration-200"
        style="will-change: transform"
        @click="triggerFileInput"
        @mouseenter="showOverlay = true"
        @mouseleave="showOverlay = false"
      >
        <!-- Current Photo atau Placeholder -->
        <img 
          v-if="currentPhoto" 
          :src="photoUrl" 
          :alt="altText"
          class="w-full h-full object-cover"
        />
        <div v-else class="w-full h-full flex items-center justify-center bg-gradient-to-br from-gray-100 to-gray-200">
          <svg class="w-16 h-16 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
          </svg>
        </div>

        <!-- Overlay untuk "Change Photo" on hover -->
        <Transition>
          <div v-if="showOverlay || !currentPhoto" class="absolute inset-0 flex flex-col items-center justify-center bg-gray-900/70 text-white transition-opacity duration-200">
            <svg class="w-6 h-6 mb-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 9a2 2 0 012-2h.93a2 2 0 001.664-.89l.812-1.22A2 2 0 0110.07 4h3.86a2 2 0 011.664.89l.812 1.22A2 2 0 0018.07 7H19a2 2 0 012 2v9a2 2 0 01-2 2H5a2 2 0 01-2-2V9z" />
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 13a3 3 0 11-6 0 3 3 0 016 0z" />
            </svg>
            <span class="text-xs font-medium">
              {{ currentPhoto ? 'Ganti Foto' : 'Upload Foto' }}
            </span>
          </div>
        </Transition>

        <!-- Loading overlay saat upload -->
        <Transition>
          <div v-if="loading" class="absolute inset-0 flex flex-col items-center justify-center bg-white/90">
            <div class="spinner"></div>
            <span class="mt-2 text-xs font-medium text-gray-600">Uploading...</span>
          </div>
        </Transition>
      </div>

      <!-- Delete button jika ada photo -->
      <Motion
        v-if="currentPhoto && !loading"
        :initial="{ opacity: 0, scale: 0 }"
        :animate="{ opacity: 1, scale: 1 }"
        :transition="{ type: 'spring', stiffness: 500, damping: 35 }"
        class="absolute -top-1 -right-1 w-7 h-7 rounded-full bg-red-500 hover:bg-red-600 text-white flex items-center justify-center cursor-pointer active-scale shadow-md hover:shadow-lg transition-all duration-200"
        @click="handleDelete"
      >
        <svg class="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
          <path fill-rule="evenodd" d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z" clip-rule="evenodd" />
        </svg>
      </Motion>
    </Motion>

    <!-- Hidden file input -->
    <input
      ref="fileInput"
      type="file"
      accept="image/jpeg,image/jpg,image/png,image/webp"
      class="hidden"
      @change="handleFileSelect"
    />

    <!-- Info text -->
    <p v-if="showInfo" class="text-xs text-gray-500 text-center">
      Format: JPG, PNG, WebP. Maksimal 5MB
    </p>

    <!-- Error message -->
    <Transition>
      <div v-if="error" class="text-sm text-red-600 text-center bg-red-50 border border-red-200 rounded-lg px-3 py-2 transition-opacity duration-200">
        {{ error }}
      </div>
    </Transition>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { Motion } from 'motion-v'
import { entranceAnimations } from '@/composables/useMotion'
import { useHaptic } from '@/composables/useHaptic'
import { useApi } from '@/composables/useApi'

const props = defineProps({
  currentPhoto: {
    type: String,
    default: '',
  },
  altText: {
    type: String,
    default: 'Profile Photo',
  },
  showInfo: {
    type: Boolean,
    default: true,
  },
  maxSizeInMB: {
    type: Number,
    default: 5,
  },
})

const emit = defineEmits(['upload-success', 'upload-error', 'delete-success'])

const { post, del } = useApi()
const haptic = useHaptic()

const fileInput = ref(null)
const showOverlay = ref(false)
const loading = ref(false)
const error = ref('')

/**
 * Compute full photo URL dengan base URL
 */
const photoUrl = computed(() => {
  if (!props.currentPhoto) return ''
  
  // Jika sudah full URL, return as is
  if (props.currentPhoto.startsWith('http')) {
    return props.currentPhoto
  }
  
  // Jika relative path, tambahkan base URL
  const baseURL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080'
  return `${baseURL.replace('/api', '')}${props.currentPhoto}`
})

/**
 * Trigger file input click
 */
const triggerFileInput = () => {
  if (loading.value) return
  haptic.light()
  fileInput.value?.click()
}

/**
 * Validate file sebelum upload
 */
const validateFile = (file) => {
  // Check file type
  const allowedTypes = ['image/jpeg', 'image/jpg', 'image/png', 'image/webp']
  if (!allowedTypes.includes(file.type)) {
    return 'Format file tidak didukung. Gunakan JPG, PNG, atau WebP'
  }
  
  // Check file size
  const maxSize = props.maxSizeInMB * 1024 * 1024 // Convert to bytes
  if (file.size > maxSize) {
    return `Ukuran file terlalu besar. Maksimal ${props.maxSizeInMB}MB`
  }
  
  return null
}

/**
 * Handle file selection dan upload
 */
const handleFileSelect = async (event) => {
  const file = event.target.files?.[0]
  if (!file) return
  
  // Reset error
  error.value = ''
  
  // Validate file
  const validationError = validateFile(file)
  if (validationError) {
    error.value = validationError
    haptic.error()
    return
  }
  
  // Upload file
  await uploadPhoto(file)
  
  // Reset input untuk allow re-upload same file
  event.target.value = ''
}

/**
 * Upload photo ke server
 */
const uploadPhoto = async (file) => {
  loading.value = true
  
  try {
    // Create FormData
    const formData = new FormData()
    formData.append('photo', file)
    
    // Upload dengan multipart/form-data
    const response = await post('/profile/photo', formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    })
    
    if (response.success) {
      haptic.success()
      emit('upload-success', response.data.photo_url)
    } else {
      throw new Error(response.message || 'Gagal upload photo')
    }
  } catch (err) {
    error.value = err.response?.data?.message || err.message || 'Gagal upload photo'
    haptic.error()
    emit('upload-error', error.value)
  } finally {
    loading.value = false
  }
}

/**
 * Handle delete photo
 */
const handleDelete = async () => {
  if (loading.value) return
  
  // Confirm delete
  if (!confirm('Yakin ingin menghapus foto profil?')) return
  
  loading.value = true
  haptic.medium()
  
  try {
    const response = await del('/profile/photo')
    
    if (response.success) {
      haptic.success()
      emit('delete-success')
    } else {
      throw new Error(response.message || 'Gagal menghapus photo')
    }
  } catch (err) {
    error.value = err.response?.data?.message || err.message || 'Gagal menghapus photo'
    haptic.error()
  } finally {
    loading.value = false
  }
}
</script>
