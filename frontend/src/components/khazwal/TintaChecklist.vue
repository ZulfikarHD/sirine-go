<template>
  <div class="tinta-checklist">
    <!-- Header -->
    <Motion v-bind="entranceAnimations.fadeUp">
      <div class="mb-6 p-4 bg-purple-50 border-2 border-purple-200 rounded-xl">
        <div class="flex items-center gap-3 mb-2">
          <Droplet class="w-6 h-6 text-purple-600" />
          <p class="text-sm font-semibold text-purple-900">Daftar Tinta yang Dibutuhkan</p>
        </div>
        <p class="text-sm text-purple-700">
          Centang setiap warna tinta yang telah diambil dan masukkan quantity (kg)
        </p>
      </div>
    </Motion>

    <!-- Tinta List -->
    <div class="space-y-3 mb-6">
      <Motion
        v-for="(tinta, index) in tintaItems"
        :key="tinta.color"
        :initial="{ opacity: 0, x: -10 }"
        :animate="{ opacity: 1, x: 0 }"
        :transition="{ duration: 0.2, delay: index * 0.05, ease: 'easeOut' }"
      >
        <div
          class="p-4 rounded-xl border-2 transition-all"
          :class="getTintaCardClass(tinta)"
        >
          <!-- Header: Checkbox + Color Name -->
          <div class="flex items-center gap-3 mb-3">
            <button
              @click="toggleTinta(tinta.color)"
              class="flex-shrink-0 w-6 h-6 rounded-lg border-2 flex items-center justify-center transition-all active-scale"
              :class="tinta.checked 
                ? 'bg-purple-600 border-purple-600' 
                : 'border-gray-300 bg-white hover:border-purple-400'"
            >
              <Check v-if="tinta.checked" class="w-4 h-4 text-white" />
            </button>

            <!-- Color Indicator -->
            <div
              class="w-8 h-8 rounded-lg border-2 border-gray-300 flex-shrink-0"
              :style="{ backgroundColor: getColorHex(tinta.color) }"
            ></div>

            <!-- Color Name -->
            <div class="flex-1">
              <p class="font-semibold text-gray-900">
                {{ tinta.color }}
              </p>
              <p v-if="tinta.requirement" class="text-xs text-gray-500">
                Kebutuhan: {{ tinta.requirement }} kg
              </p>
            </div>
          </div>

          <!-- Quantity Input (shown when checked) -->
          <Motion
            v-if="tinta.checked"
            :initial="{ opacity: 0, height: 0 }"
            :animate="{ opacity: 1, height: 'auto' }"
            :transition="springPresets.default"
          >
            <div class="pl-9">
              <label class="block text-xs font-semibold text-gray-700 mb-2">
                Quantity (kg) <span class="text-red-500">*</span>
              </label>
              <input
                v-model.number="tinta.quantity"
                type="number"
                inputmode="decimal"
                step="0.1"
                min="0"
                placeholder="0.0"
                class="w-full px-3 py-2 border-2 border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-purple-500 focus:border-transparent"
                :class="{ 'border-red-300 bg-red-50': tinta.checked && !tinta.quantity }"
                @input="updateTinta(tinta.color, tinta.quantity)"
              />
              
              <!-- Low Stock Warning -->
              <Motion
                v-if="tinta.quantity && tinta.quantity < 10"
                :initial="{ opacity: 0, y: -5 }"
                :animate="{ opacity: 1, y: 0 }"
                :transition="{ duration: 0.2, ease: 'easeOut' }"
                class="mt-2 flex items-start gap-1 text-xs text-yellow-700"
              >
                <AlertTriangle class="w-3 h-3 flex-shrink-0 mt-0.5" />
                <span>Stok tinta rendah (< 10 kg)</span>
              </Motion>
            </div>
          </Motion>
        </div>
      </Motion>
    </div>

    <!-- Progress Bar -->
    <Motion
      :initial="{ opacity: 0, y: 10 }"
      :animate="{ opacity: 1, y: 0 }"
      :transition="{ duration: 0.2, delay: 0.1, ease: 'easeOut' }"
      class="mb-6"
    >
      <div class="p-4 bg-gray-50 border border-gray-200 rounded-xl">
        <div class="flex items-center justify-between mb-2">
          <p class="text-sm font-semibold text-gray-700">Progress</p>
          <p class="text-sm font-bold text-purple-600">
            {{ checkedCount }} / {{ tintaItems.length }}
          </p>
        </div>
        <div class="w-full h-2 bg-gray-200 rounded-full overflow-hidden">
          <div
            class="h-full bg-gradient-to-r from-purple-500 to-fuchsia-500 transition-all duration-300"
            :style="{ width: `${progressPercentage}%` }"
          ></div>
        </div>
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
        class="w-full px-6 py-4 bg-gradient-to-r from-purple-600 to-fuchsia-600 text-white font-semibold text-lg rounded-xl hover:shadow-lg active-scale disabled:opacity-50 disabled:cursor-not-allowed transition-all"
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
          Pastikan semua tinta yang diperlukan sudah dicentang dan quantity telah diisi.
        </p>
      </div>
    </Motion>
  </div>
</template>

<script setup>
/**
 * TintaChecklist Component
 * Checklist component untuk tinta dengan quantity input per warna
 * dan low stock warning untuk inventory management
 */
import { ref, computed, watch } from 'vue'
import { Motion } from 'motion-v'
import { entranceAnimations, springPresets } from '@/composables/useMotion'
import { 
  Droplet, 
  Check, 
  AlertTriangle, 
  Info,
  Loader
} from 'lucide-vue-next'

const props = defineProps({
  /**
   * Array tinta requirements dari PO
   * Format: [{ color: string, requirement?: number }]
   */
  requirements: {
    type: Array,
    required: true
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

// State: Initialize tinta items dari requirements
const tintaItems = ref([])

/**
 * Initialize tinta items dari requirements prop
 */
const initializeTintaItems = () => {
  tintaItems.value = props.requirements.map(req => {
    // Handle both object and string format
    const color = typeof req === 'string' ? req : req.color || req
    const requirement = typeof req === 'object' ? req.requirement : null
    
    return {
      color,
      requirement,
      checked: false,
      quantity: null
    }
  })
}

// Initialize on mount
watch(() => props.requirements, () => {
  initializeTintaItems()
}, { immediate: true })

// Computed: Checked count
const checkedCount = computed(() => {
  return tintaItems.value.filter(t => t.checked).length
})

// Computed: Progress percentage
const progressPercentage = computed(() => {
  if (tintaItems.value.length === 0) return 0
  return (checkedCount.value / tintaItems.value.length) * 100
})

// Computed: Validation
const isValid = computed(() => {
  // Minimal 1 tinta harus dicek
  if (checkedCount.value === 0) return false
  
  // Semua tinta yang dicek harus punya quantity > 0
  const allCheckedHaveQuantity = tintaItems.value
    .filter(t => t.checked)
    .every(t => t.quantity && t.quantity > 0)
  
  return allCheckedHaveQuantity
})

/**
 * Toggle tinta checkbox
 */
const toggleTinta = (color) => {
  const tinta = tintaItems.value.find(t => t.color === color)
  if (tinta) {
    tinta.checked = !tinta.checked
    
    // Reset quantity saat uncheck
    if (!tinta.checked) {
      tinta.quantity = null
    }
    
    // Haptic feedback
    if ('vibrate' in navigator) {
      navigator.vibrate(10)
    }
  }
}

/**
 * Update tinta quantity
 */
const updateTinta = (color, quantity) => {
  const tinta = tintaItems.value.find(t => t.color === color)
  if (tinta) {
    tinta.quantity = quantity
    
    // Haptic feedback
    if ('vibrate' in navigator) {
      navigator.vibrate(5)
    }
  }
}

/**
 * Get tinta card class berdasarkan state
 */
const getTintaCardClass = (tinta) => {
  if (tinta.checked) {
    return 'bg-purple-50 border-purple-300'
  }
  return 'bg-white border-gray-200 opacity-70 hover:opacity-100'
}

/**
 * Get color hex untuk visual indicator
 * Map common color names ke hex values
 */
const getColorHex = (colorName) => {
  const colorMap = {
    // Basic colors
    'cyan': '#00FFFF',
    'magenta': '#FF00FF',
    'yellow': '#FFFF00',
    'black': '#000000',
    'white': '#FFFFFF',
    
    // Extended colors
    'red': '#FF0000',
    'blue': '#0000FF',
    'green': '#00FF00',
    'orange': '#FFA500',
    'purple': '#800080',
    'pink': '#FFC0CB',
    'brown': '#A52A2A',
    'grey': '#808080',
    'gray': '#808080',
    'gold': '#FFD700',
    'silver': '#C0C0C0',
    
    // Process colors (CMYK)
    'process cyan': '#00B0F0',
    'process magenta': '#EC008C',
    'process yellow': '#FFF200',
    'process black': '#000000',
    
    // Pantone-like
    'reflex blue': '#001489',
    'warm red': '#F9423A',
    'bright orange': '#FF6600'
  }
  
  const normalized = colorName.toLowerCase().trim()
  return colorMap[normalized] || '#6B7280' // Default gray if not found
}

/**
 * Handle form submit
 */
const handleSubmit = () => {
  if (!isValid.value) return
  
  // Filter hanya tinta yang checked
  const checkedTinta = tintaItems.value
    .filter(t => t.checked)
    .map(t => ({
      color: t.color,
      quantity: t.quantity,
      checked: true
    }))
  
  emit('submit', checkedTinta)
  
  // Haptic feedback
  if ('vibrate' in navigator) {
    navigator.vibrate([10, 50, 10])
  }
}
</script>
