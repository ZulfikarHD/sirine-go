<template>
  <div class="barcode-scanner">
    <!-- Camera View -->
    <Motion
      v-if="!scanSuccess"
      v-bind="entranceAnimations.fadeScale"
      class="relative"
    >
      <!-- Scanner Container -->
      <div
        v-show="isScanning"
        id="barcode-scanner-container"
        class="relative w-full aspect-video bg-gray-900 rounded-2xl overflow-hidden"
      >
        <!-- Scanner overlay dengan gradient borders -->
        <div class="absolute inset-0 pointer-events-none">
          <!-- Top gradient -->
          <div class="absolute top-0 left-0 right-0 h-20 bg-gradient-to-b from-black/50 to-transparent"></div>
          <!-- Bottom gradient -->
          <div class="absolute bottom-0 left-0 right-0 h-20 bg-gradient-to-t from-black/50 to-transparent"></div>
          
          <!-- Scan line animation -->
          <div class="absolute inset-x-0 top-1/2 -translate-y-1/2 h-1 bg-gradient-to-r from-transparent via-indigo-500 to-transparent animate-pulse"></div>
          
          <!-- Corner guides -->
          <div class="absolute top-1/4 left-1/4 w-12 h-12 border-t-4 border-l-4 border-indigo-500 rounded-tl-xl"></div>
          <div class="absolute top-1/4 right-1/4 w-12 h-12 border-t-4 border-r-4 border-indigo-500 rounded-tr-xl"></div>
          <div class="absolute bottom-1/4 left-1/4 w-12 h-12 border-b-4 border-l-4 border-indigo-500 rounded-bl-xl"></div>
          <div class="absolute bottom-1/4 right-1/4 w-12 h-12 border-b-4 border-r-4 border-indigo-500 rounded-br-xl"></div>
        </div>
      </div>

      <!-- Camera Permission / Error State -->
      <div
        v-if="!isScanning && !scanSuccess"
        class="w-full aspect-video bg-gray-100 rounded-2xl flex flex-col items-center justify-center p-6 text-center"
      >
        <Camera v-if="!error" class="w-16 h-16 text-gray-400 mb-4" />
        <AlertTriangle v-else class="w-16 h-16 text-red-500 mb-4" />
        
        <p v-if="!error" class="text-gray-600 mb-4">
          Klik tombol di bawah untuk mengaktifkan kamera
        </p>
        <p v-else class="text-red-600 mb-4">
          {{ error }}
        </p>

        <button
          v-if="!isScanning"
          @click="startScanning"
          :disabled="loading"
          class="px-6 py-3 bg-gradient-to-r from-indigo-600 to-fuchsia-600 text-white font-semibold rounded-xl hover:shadow-lg active-scale disabled:opacity-50"
        >
          <Loader v-if="loading" class="w-4 h-4 inline-block mr-2 animate-spin" />
          <Camera v-else class="w-4 h-4 inline-block mr-2" />
          {{ loading ? 'Memulai...' : 'Aktifkan Kamera' }}
        </button>
      </div>

      <!-- Stop Scanner Button -->
      <button
        v-if="isScanning"
        @click="stopScanning"
        class="mt-4 w-full px-4 py-3 bg-red-500 text-white font-semibold rounded-xl hover:bg-red-600 active-scale"
      >
        <X class="w-4 h-4 inline-block mr-2" />
        Hentikan Scanner
      </button>
    </Motion>

    <!-- Manual Input Fallback -->
    <Motion
      v-if="!isScanning || showManualInput"
      :initial="{ opacity: 0, y: 10 }"
      :animate="{ opacity: 1, y: 0 }"
      :transition="{ duration: 0.2, delay: 0.1, ease: 'easeOut' }"
      class="mt-4"
    >
      <div class="p-4 bg-blue-50 border border-blue-200 rounded-xl">
        <p class="text-sm text-blue-800 mb-3 flex items-start gap-2">
          <Info class="w-4 h-4 flex-shrink-0 mt-0.5" />
          <span>Jika scanner tidak berfungsi, masukkan kode barcode secara manual</span>
        </p>
        
        <div class="flex gap-2">
          <input
            v-model="manualCode"
            type="text"
            placeholder="Masukkan kode barcode..."
            class="flex-1 px-4 py-2 border border-blue-300 rounded-xl focus:outline-none focus:ring-2 focus:ring-indigo-500"
            @keyup.enter="handleManualSubmit"
          />
          <button
            @click="handleManualSubmit"
            :disabled="!manualCode.trim()"
            class="px-6 py-2 bg-indigo-600 text-white font-semibold rounded-xl hover:bg-indigo-700 active-scale disabled:opacity-50 disabled:cursor-not-allowed"
          >
            Submit
          </button>
        </div>
      </div>
    </Motion>

    <!-- Success State -->
    <Motion
      v-if="scanSuccess"
      v-bind="entranceAnimations.fadeScale"
      class="w-full aspect-video bg-gradient-to-br from-green-50 to-emerald-50 rounded-2xl flex flex-col items-center justify-center p-6 border-2 border-green-200"
    >
      <Motion v-bind="iconAnimations.popIn">
        <div class="w-20 h-20 bg-green-500 rounded-full flex items-center justify-center mb-4">
          <CheckCircle class="w-12 h-12 text-white" />
        </div>
      </Motion>
      
      <p class="text-xl font-bold text-green-900 mb-2">
        Barcode Berhasil Discan!
      </p>
      <p class="text-green-700 font-mono text-lg">
        {{ scannedCode }}
      </p>
    </Motion>
  </div>
</template>

<script setup>
/**
 * BarcodeScanner Component
 * Component untuk scan barcode dengan camera access dan manual input fallback
 * menggunakan html5-qrcode library untuk barcode detection
 */
import { ref, onMounted, onBeforeUnmount, watch } from 'vue'
import { Motion } from 'motion-v'
import { entranceAnimations, iconAnimations } from '@/composables/useMotion'
import { Camera, X, CheckCircle, AlertTriangle, Info, Loader } from 'lucide-vue-next'

const props = defineProps({
  /**
   * Expected barcode code untuk validation
   */
  expectedCode: {
    type: String,
    default: null
  },
  /**
   * Auto start scanning on mount
   */
  autoStart: {
    type: Boolean,
    default: false
  },
  /**
   * Show manual input option
   */
  showManualInput: {
    type: Boolean,
    default: true
  }
})

const emit = defineEmits(['scan-success', 'scan-error', 'code-scanned'])

// State
const isScanning = ref(false)
const loading = ref(false)
const error = ref(null)
const scanSuccess = ref(false)
const scannedCode = ref('')
const manualCode = ref('')
const html5QrCode = ref(null)

/**
 * Start barcode scanning dengan camera access
 */
const startScanning = async () => {
  loading.value = true
  error.value = null

  try {
    // Dynamic import html5-qrcode
    const { Html5Qrcode } = await import('html5-qrcode')
    
    html5QrCode.value = new Html5Qrcode('barcode-scanner-container')
    
    // Start scanning dengan optimal config
    await html5QrCode.value.start(
      { facingMode: 'environment' }, // Use back camera
      {
        fps: 10, // Frames per second untuk scanning
        qrbox: { width: 250, height: 250 }, // Scan area
        aspectRatio: 16/9
      },
      onScanSuccess,
      onScanFailure
    )
    
    isScanning.value = true
    loading.value = false
    
    // Haptic feedback
    if ('vibrate' in navigator) {
      navigator.vibrate(10)
    }
  } catch (err) {
    console.error('Error starting scanner:', err)
    error.value = 'Gagal mengakses kamera. Pastikan izin kamera telah diberikan.'
    loading.value = false
    
    // Emit error untuk parent component
    emit('scan-error', err)
  }
}

/**
 * Stop barcode scanning dan release camera
 */
const stopScanning = async () => {
  if (html5QrCode.value && isScanning.value) {
    try {
      await html5QrCode.value.stop()
      html5QrCode.value.clear()
      isScanning.value = false
      
      // Haptic feedback
      if ('vibrate' in navigator) {
        navigator.vibrate(10)
      }
    } catch (err) {
      console.error('Error stopping scanner:', err)
    }
  }
}

/**
 * Handle successful barcode scan
 */
const onScanSuccess = (decodedText, decodedResult) => {
  // Stop scanning
  stopScanning()
  
  scannedCode.value = decodedText
  
  // Validate jika expectedCode diberikan
  if (props.expectedCode && decodedText !== props.expectedCode) {
    error.value = `Kode tidak sesuai! Expected: ${props.expectedCode}, Got: ${decodedText}`
    emit('scan-error', { 
      message: 'Kode tidak sesuai', 
      expected: props.expectedCode, 
      scanned: decodedText 
    })
    
    // Heavy haptic untuk error
    if ('vibrate' in navigator) {
      navigator.vibrate([30, 100, 30])
    }
    return
  }
  
  // Success
  scanSuccess.value = true
  emit('scan-success', decodedText)
  emit('code-scanned', decodedText)
  
  // Success haptic
  if ('vibrate' in navigator) {
    navigator.vibrate([10, 50, 10])
  }
}

/**
 * Handle scan failure (called continuously, tidak perlu log)
 */
const onScanFailure = (errorMessage) => {
  // Jangan log setiap failure karena terlalu banyak
  // Scanner akan terus mencoba scan
}

/**
 * Handle manual barcode input submit
 */
const handleManualSubmit = () => {
  const code = manualCode.value.trim()
  if (!code) return
  
  // Validate jika expectedCode diberikan
  if (props.expectedCode && code !== props.expectedCode) {
    error.value = `Kode tidak sesuai! Expected: ${props.expectedCode}, Got: ${code}`
    emit('scan-error', { 
      message: 'Kode tidak sesuai', 
      expected: props.expectedCode, 
      scanned: code 
    })
    
    // Heavy haptic untuk error
    if ('vibrate' in navigator) {
      navigator.vibrate([30, 100, 30])
    }
    return
  }
  
  // Success
  scannedCode.value = code
  scanSuccess.value = true
  emit('scan-success', code)
  emit('code-scanned', code)
  
  // Success haptic
  if ('vibrate' in navigator) {
    navigator.vibrate([10, 50, 10])
  }
}

/**
 * Reset scanner state untuk scan ulang
 */
const reset = () => {
  scanSuccess.value = false
  scannedCode.value = ''
  manualCode.value = ''
  error.value = null
  
  if (props.autoStart) {
    startScanning()
  }
}

// Expose reset method untuk parent component
defineExpose({
  reset,
  startScanning,
  stopScanning
})

// Lifecycle hooks
onMounted(() => {
  if (props.autoStart) {
    startScanning()
  }
})

onBeforeUnmount(() => {
  stopScanning()
})
</script>

<style scoped>
/* Custom styles untuk scanner container */
#barcode-scanner-container {
  position: relative;
}

/* Override html5-qrcode default styles jika perlu */
:deep(#barcode-scanner-container video) {
  border-radius: 1rem;
  object-fit: cover;
}
</style>
