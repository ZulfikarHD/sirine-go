<template>
  <BaseModal
    v-model="isOpen"
    title="Import Users dari CSV"
    size="lg"
    :show-footer="false"
    :dismissible="!uploading"
    @close="handleClose"
  >
    <div class="csv-import-container">
      <!-- Step 1: File Upload -->
      <div v-if="step === 1" class="upload-section">
        <!-- Drag & Drop Area -->
        <div
          class="dropzone"
          :class="{ 'dropzone-active': isDragging, 'dropzone-disabled': uploading }"
          @dragover.prevent="handleDragOver"
          @dragleave.prevent="handleDragLeave"
          @drop.prevent="handleDrop"
          @click="triggerFileInput"
        >
          <svg class="dropzone-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12" />
          </svg>
          
          <p class="dropzone-text">
            <span class="font-semibold">Klik untuk upload</span> atau drag & drop
          </p>
          <p class="dropzone-subtext">CSV file (maksimal 10MB)</p>
        </div>

        <!-- File Info jika sudah ada -->
        <Transition name="fade">
          <div v-if="selectedFile" class="file-info">
            <div class="file-details">
              <svg class="w-5 h-5 text-indigo-600" fill="currentColor" viewBox="0 0 20 20">
                <path fill-rule="evenodd" d="M8 4a3 3 0 00-3 3v4a5 5 0 0010 0V7a1 1 0 112 0v4a7 7 0 11-14 0V7a5 5 0 0110 0v4a3 3 0 11-6 0V7a1 1 0 012 0v4a1 1 0 102 0V7a3 3 0 00-3-3z" clip-rule="evenodd" />
              </svg>
              <div class="file-name-size">
                <span class="file-name">{{ selectedFile.name }}</span>
                <span class="file-size">{{ formatFileSize(selectedFile.size) }}</span>
              </div>
            </div>
            <button @click.stop="clearFile" class="remove-button" :disabled="uploading">
              <svg class="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
                <path fill-rule="evenodd" d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z" clip-rule="evenodd" />
              </svg>
            </button>
          </div>
        </Transition>

        <!-- Hidden File Input -->
        <input
          ref="fileInput"
          type="file"
          accept=".csv,text/csv,application/vnd.ms-excel"
          class="hidden"
          @change="handleFileSelect"
        />

        <!-- CSV Format Info -->
        <div class="format-info">
          <h4 class="format-title">Format CSV yang Diharapkan:</h4>
          <div class="format-example">
            <code class="format-code">
              NIP,Full Name,Email,Phone,Role,Department<br>
              12345,John Doe,john@example.com,08123456789,STAFF_KHAZWAL,KHAZWAL
            </code>
          </div>
          <ul class="format-notes">
            <li>Header harus sesuai urutan di atas</li>
            <li>NIP maksimal 5 digit dan harus unique</li>
            <li>Email harus valid dan unique</li>
            <li>Password akan di-generate otomatis</li>
          </ul>
        </div>

        <!-- Action Buttons -->
        <div class="action-buttons">
          <button @click="handleClose" class="btn-secondary" :disabled="uploading">
            Batal
          </button>
          <button 
            @click="handleUpload" 
            class="btn-primary"
            :disabled="!selectedFile || uploading"
          >
            <span v-if="!uploading">Import Users</span>
            <span v-else class="flex items-center gap-2">
              <span class="spinner w-4 h-4 border-2 border-white/30 border-t-white rounded-full"></span>
              Importing...
            </span>
          </button>
        </div>
      </div>

      <!-- Step 2: Import Result -->
      <div v-if="step === 2" class="result-section">
        <!-- Success Summary -->
        <div class="result-summary">
          <div class="result-icon" :class="result.failed === 0 ? 'bg-green-100' : 'bg-yellow-100'">
            <svg v-if="result.failed === 0" class="w-12 h-12 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            <svg v-else class="w-12 h-12 text-yellow-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
            </svg>
          </div>

          <h3 class="result-title">Import Selesai</h3>
          
          <div class="result-stats">
            <div class="result-stat">
              <span class="stat-value text-green-600">{{ result.imported }}</span>
              <span class="stat-label">Berhasil</span>
            </div>
            <div class="result-stat">
              <span class="stat-value text-red-600">{{ result.failed }}</span>
              <span class="stat-label">Gagal</span>
            </div>
          </div>
        </div>

        <!-- Error List jika ada -->
        <div v-if="result.errors && result.errors.length > 0" class="error-list">
          <h4 class="error-list-title">
            <svg class="w-5 h-5 text-red-600" fill="currentColor" viewBox="0 0 20 20">
              <path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7 4a1 1 0 11-2 0 1 1 0 012 0zm-1-9a1 1 0 00-1 1v4a1 1 0 102 0V6a1 1 0 00-1-1z" clip-rule="evenodd" />
            </svg>
            Detail Error ({{ result.errors.length }} item)
          </h4>
          <div class="error-items">
            <Motion
              v-for="(error, index) in result.errors"
              :key="index"
              :initial="{ opacity: 0, x: -10 }"
              :animate="{ opacity: 1, x: 0 }"
              :transition="{ duration: 0.2, delay: index * 0.05 }"
              class="error-item"
            >
              <span class="error-row">Row {{ error.row }}</span>
              <span class="error-nip">NIP: {{ error.nip }}</span>
              <span class="error-reason">{{ error.reason }}</span>
            </Motion>
          </div>
        </div>

        <!-- Action Buttons -->
        <div class="action-buttons">
          <button @click="handleReset" class="btn-secondary">
            Import Lagi
          </button>
          <button @click="handleFinish" class="btn-primary">
            Selesai
          </button>
        </div>
      </div>
    </div>
  </BaseModal>
</template>

<script setup>
import { ref, computed } from 'vue'
import { Motion } from 'motion-v'
import { BaseModal } from '@/components/common'
import { useApi } from '@/composables/useApi'
import { useHaptic } from '@/composables/useHaptic'

const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false,
  },
})

const emit = defineEmits(['update:modelValue', 'import-success'])

const { post } = useApi()
const haptic = useHaptic()

const isOpen = computed({
  get: () => props.modelValue,
  set: (value) => emit('update:modelValue', value),
})

const fileInput = ref(null)
const selectedFile = ref(null)
const isDragging = ref(false)
const uploading = ref(false)
const step = ref(1) // 1: Upload, 2: Result
const result = ref({
  imported: 0,
  failed: 0,
  errors: [],
})

/**
 * Trigger file input click
 */
const triggerFileInput = () => {
  if (uploading.value) return
  haptic.light()
  fileInput.value?.click()
}

/**
 * Handle file selection
 */
const handleFileSelect = (event) => {
  const file = event.target.files?.[0]
  if (file) {
    selectFile(file)
  }
  // Reset input untuk allow re-upload same file
  event.target.value = ''
}

/**
 * Handle drag over
 */
const handleDragOver = (event) => {
  if (uploading.value) return
  isDragging.value = true
}

/**
 * Handle drag leave
 */
const handleDragLeave = () => {
  isDragging.value = false
}

/**
 * Handle file drop
 */
const handleDrop = (event) => {
  isDragging.value = false
  
  if (uploading.value) return
  
  const file = event.dataTransfer.files?.[0]
  if (file) {
    selectFile(file)
  }
}

/**
 * Select file dan validate
 */
const selectFile = (file) => {
  // Validate file type
  if (!file.name.endsWith('.csv') && file.type !== 'text/csv') {
    alert('File harus berformat CSV')
    return
  }
  
  // Validate file size (max 10MB)
  if (file.size > 10 * 1024 * 1024) {
    alert('Ukuran file terlalu besar. Maksimal 10MB')
    return
  }
  
  selectedFile.value = file
  haptic.success()
}

/**
 * Clear selected file
 */
const clearFile = () => {
  selectedFile.value = null
  haptic.light()
}

/**
 * Format file size untuk display
 */
const formatFileSize = (bytes) => {
  if (bytes === 0) return '0 Bytes'
  const k = 1024
  const sizes = ['Bytes', 'KB', 'MB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return Math.round(bytes / Math.pow(k, i) * 100) / 100 + ' ' + sizes[i]
}

/**
 * Handle upload CSV
 */
const handleUpload = async () => {
  if (!selectedFile.value || uploading.value) return
  
  uploading.value = true
  haptic.medium()
  
  try {
    // Create FormData
    const formData = new FormData()
    formData.append('csv_file', selectedFile.value)
    
    // Upload dengan multipart/form-data
    const response = await post('/users/import', formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    })
    
    if (response.success) {
      result.value = response.data
      step.value = 2
      haptic.success()
    } else {
      throw new Error(response.message || 'Gagal import users')
    }
  } catch (error) {
    alert(error.response?.data?.message || error.message || 'Gagal import users')
    haptic.error()
  } finally {
    uploading.value = false
  }
}

/**
 * Handle reset untuk import lagi
 */
const handleReset = () => {
  selectedFile.value = null
  step.value = 1
  result.value = { imported: 0, failed: 0, errors: [] }
  haptic.light()
}

/**
 * Handle finish dan close modal
 */
const handleFinish = () => {
  emit('import-success', result.value)
  handleClose()
}

/**
 * Handle close modal
 */
const handleClose = () => {
  if (uploading.value) return
  
  // Reset state
  selectedFile.value = null
  step.value = 1
  result.value = { imported: 0, failed: 0, errors: [] }
  
  isOpen.value = false
}
</script>
