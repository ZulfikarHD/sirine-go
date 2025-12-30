<template>
  <div class="kertas-input-form">
    <!-- Target Info Card -->
    <Motion v-bind="entranceAnimations.fadeUp">
      <div class="p-4 bg-amber-50 border-2 border-amber-200 rounded-xl mb-6">
        <div class="flex items-center gap-3 mb-2">
          <FileStack class="w-6 h-6 text-amber-600" />
          <p class="text-sm font-semibold text-amber-900">Target Kertas Blanko</p>
        </div>
        <p class="text-3xl font-bold text-amber-700">
          {{ targetQuantity.toLocaleString() }} <span class="text-lg">lembar</span>
        </p>
      </div>
    </Motion>

    <!-- Actual Quantity Input -->
    <Motion
      :initial="{ opacity: 0, y: 10 }"
      :animate="{ opacity: 1, y: 0 }"
      :transition="{ duration: 0.2, delay: 0.1, ease: 'easeOut' }"
    >
      <div class="mb-6">
        <label class="block text-sm font-semibold text-gray-700 mb-2">
          Jumlah Actual <span class="text-red-500">*</span>
        </label>
        <input
          v-model.number="actualQty"
          type="number"
          inputmode="numeric"
          min="0"
          placeholder="Masukkan jumlah actual..."
          class="w-full px-4 py-3 text-lg border-2 border-gray-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent transition-all"
          :class="{ 'border-red-300 bg-red-50': hasVariance && needsReason }"
          @input="calculateVariance"
        />
      </div>
    </Motion>

    <!-- Variance Display -->
    <Motion
      v-if="actualQty !== null && actualQty !== ''"
      :initial="{ opacity: 0, scale: 0.95 }"
      :animate="{ opacity: 1, scale: 1 }"
      :transition="springPresets.default"
    >
      <div class="mb-6 p-4 rounded-xl border-2" :class="varianceCardClass">
        <div class="flex items-center justify-between mb-2">
          <div class="flex items-center gap-2">
            <component :is="varianceIcon" :class="varianceIconClass" class="w-5 h-5" />
            <p class="text-sm font-semibold" :class="varianceTextClass">
              {{ varianceLabel }}
            </p>
          </div>
          <p class="text-2xl font-bold" :class="varianceTextClass">
            {{ varianceDisplay }}
          </p>
        </div>
        
        <!-- Variance Percentage -->
        <div class="flex items-center gap-2">
          <div class="flex-1 h-2 bg-gray-200 rounded-full overflow-hidden">
            <div
              class="h-full transition-all duration-300"
              :class="varianceBarClass"
              :style="{ width: `${Math.min(Math.abs(variancePercentage), 100)}%` }"
            ></div>
          </div>
          <p class="text-sm font-semibold" :class="varianceTextClass">
            {{ variancePercentage.toFixed(1) }}%
          </p>
        </div>

        <!-- Variance Warning (if > 5%) -->
        <Motion
          v-if="needsReason"
          :initial="{ opacity: 0, y: -10 }"
          :animate="{ opacity: 1, y: 0 }"
          :transition="{ duration: 0.2, ease: 'easeOut' }"
          class="mt-3 p-3 bg-yellow-50 border border-yellow-300 rounded-lg"
        >
          <div class="flex items-start gap-2">
            <AlertTriangle class="w-4 h-4 text-yellow-600 flex-shrink-0 mt-0.5" />
            <p class="text-xs text-yellow-800">
              <strong>Variance melebihi 5%!</strong><br />
              Wajib memberikan alasan untuk variance yang signifikan.
            </p>
          </div>
        </Motion>
      </div>
    </Motion>

    <!-- Variance Reason (Required if > 5%) -->
    <Motion
      v-if="needsReason"
      :initial="{ opacity: 0, y: 10 }"
      :animate="{ opacity: 1, y: 0 }"
      :transition="{ duration: 0.2, delay: 0.1, ease: 'easeOut' }"
    >
      <div class="mb-6">
        <label class="block text-sm font-semibold text-gray-700 mb-2">
          Alasan Variance <span class="text-red-500">*</span>
        </label>
        <textarea
          v-model="varianceReason"
          rows="3"
          placeholder="Jelaskan alasan perbedaan antara target dan actual..."
          class="w-full px-4 py-3 border-2 border-gray-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent transition-all resize-none"
          :class="{ 'border-red-300 bg-red-50': !varianceReason.trim() }"
        ></textarea>
        <p class="text-xs text-gray-500 mt-1">
          Minimal 10 karakter
        </p>
      </div>
    </Motion>

    <!-- Submit Button -->
    <Motion
      :initial="{ opacity: 0, y: 10 }"
      :animate="{ opacity: 1, y: 0 }"
      :transition="{ duration: 0.2, delay: 0.2, ease: 'easeOut' }"
    >
      <button
        @click="handleSubmit"
        :disabled="!isValid || loading"
        class="w-full px-6 py-4 bg-gradient-to-r from-indigo-600 to-fuchsia-600 text-white font-semibold text-lg rounded-xl hover:shadow-lg active-scale disabled:opacity-50 disabled:cursor-not-allowed transition-all"
      >
        <Loader v-if="loading" class="w-5 h-5 inline-block mr-2 animate-spin" />
        <Check v-else class="w-5 h-5 inline-block mr-2" />
        {{ loading ? 'Menyimpan...' : 'Simpan & Lanjutkan' }}
      </button>
    </Motion>

    <!-- Info Box -->
    <Motion
      :initial="{ opacity: 0 }"
      :animate="{ opacity: 1 }"
      :transition="{ duration: 0.3, delay: 0.3, ease: 'easeOut' }"
      class="mt-4 p-3 bg-blue-50 border border-blue-200 rounded-xl"
    >
      <div class="flex items-start gap-2">
        <Info class="w-4 h-4 text-blue-600 flex-shrink-0 mt-0.5" />
        <p class="text-xs text-blue-800">
          Masukkan jumlah kertas blanko yang benar-benar diambil dari gudang.
        </p>
      </div>
    </Motion>
  </div>
</template>

<script setup>
/**
 * KertasInputForm Component
 * Form untuk input kertas blanko actual dengan real-time variance calculation,
 * conditional reason input untuk variance > 5%, dan support untuk restore data
 * saat navigate back ke step ini
 */
import { ref, computed, watch } from 'vue'
import { Motion } from 'motion-v'
import { entranceAnimations, springPresets } from '@/composables/useMotion'
import { 
  FileStack, 
  TrendingUp, 
  TrendingDown, 
  Minus, 
  AlertTriangle, 
  Check,
  Info,
  Loader
} from 'lucide-vue-next'

const props = defineProps({
  /**
   * Target quantity dari planning
   */
  targetQuantity: {
    type: Number,
    required: true
  },
  /**
   * Initial actual quantity (untuk restore data saat navigate back)
   */
  initialActualQty: {
    type: Number,
    default: null
  },
  /**
   * Initial variance reason (untuk restore data saat navigate back)
   */
  initialVarianceReason: {
    type: String,
    default: ''
  },
  /**
   * Loading state
   */
  loading: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['submit'])

// State - Initialize dengan data existing jika ada
const actualQty = ref(props.initialActualQty)
const varianceReason = ref(props.initialVarianceReason)

// Watch for prop changes untuk handle navigation back
watch(() => props.initialActualQty, (newVal) => {
  if (newVal !== null && newVal !== actualQty.value) {
    actualQty.value = newVal
  }
}, { immediate: true })

watch(() => props.initialVarianceReason, (newVal) => {
  if (newVal && newVal !== varianceReason.value) {
    varianceReason.value = newVal
  }
}, { immediate: true })

// Computed: Variance calculation
const variance = computed(() => {
  if (actualQty.value === null || actualQty.value === '') return 0
  return actualQty.value - props.targetQuantity
})

const variancePercentage = computed(() => {
  if (props.targetQuantity === 0) return 0
  return (variance.value / props.targetQuantity) * 100
})

const hasVariance = computed(() => {
  return variance.value !== 0
})

const needsReason = computed(() => {
  return Math.abs(variancePercentage.value) > 5
})

// Computed: Validation
const isValid = computed(() => {
  // Actual qty harus diisi
  if (actualQty.value === null || actualQty.value === '' || actualQty.value < 0) {
    return false
  }
  
  // Jika variance > 5%, reason harus diisi minimal 10 karakter
  if (needsReason.value && varianceReason.value.trim().length < 10) {
    return false
  }
  
  return true
})

// Computed: Variance display
const varianceDisplay = computed(() => {
  if (variance.value === 0) return '0'
  const sign = variance.value > 0 ? '+' : ''
  return `${sign}${variance.value.toLocaleString()}`
})

const varianceLabel = computed(() => {
  if (variance.value > 0) return 'Kelebihan'
  if (variance.value < 0) return 'Kekurangan'
  return 'Sesuai Target'
})

const varianceIcon = computed(() => {
  if (variance.value > 0) return TrendingUp
  if (variance.value < 0) return TrendingDown
  return Minus
})

// Computed: Styling classes
const varianceCardClass = computed(() => {
  if (variance.value === 0) {
    return 'bg-green-50 border-green-200'
  }
  if (Math.abs(variancePercentage.value) > 5) {
    return 'bg-red-50 border-red-300'
  }
  if (variance.value > 0) {
    return 'bg-yellow-50 border-yellow-200'
  }
  return 'bg-orange-50 border-orange-200'
})

const varianceTextClass = computed(() => {
  if (variance.value === 0) return 'text-green-700'
  if (Math.abs(variancePercentage.value) > 5) return 'text-red-700'
  if (variance.value > 0) return 'text-yellow-700'
  return 'text-orange-700'
})

const varianceIconClass = computed(() => {
  if (variance.value === 0) return 'text-green-600'
  if (Math.abs(variancePercentage.value) > 5) return 'text-red-600'
  if (variance.value > 0) return 'text-yellow-600'
  return 'text-orange-600'
})

const varianceBarClass = computed(() => {
  if (variance.value === 0) return 'bg-green-500'
  if (Math.abs(variancePercentage.value) > 5) return 'bg-red-500'
  if (variance.value > 0) return 'bg-yellow-500'
  return 'bg-orange-500'
})

/**
 * Calculate variance saat input berubah
 */
const calculateVariance = () => {
  // Haptic feedback untuk iOS-like experience
  if ('vibrate' in navigator && actualQty.value !== null) {
    navigator.vibrate(5)
  }
}

/**
 * Handle form submit
 */
const handleSubmit = () => {
  if (!isValid.value) return
  
  emit('submit', {
    actualQty: actualQty.value,
    varianceReason: varianceReason.value.trim()
  })
  
  // Haptic feedback
  if ('vibrate' in navigator) {
    navigator.vibrate([10, 50, 10])
  }
}

// Reset variance reason saat tidak diperlukan
watch(needsReason, (newValue) => {
  if (!newValue) {
    varianceReason.value = ''
  }
})
</script>
