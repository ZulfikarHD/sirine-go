<template>
  <Teleport to="body">
    <!-- Backdrop -->
    <Motion
      v-if="modelValue"
      :initial="{ opacity: 0 }"
      :animate="{ opacity: 1 }"
      :exit="{ opacity: 0 }"
      :transition="{ duration: 0.2, ease: 'easeOut' }"
      class="fixed inset-0 z-[100] bg-black/95"
      @click="handleBackdropClick"
    >
      <!-- Container -->
      <div class="relative h-full flex flex-col">
        <!-- Header -->
        <div class="flex items-center justify-between p-4 z-10">
          <div class="text-white">
            <p class="font-semibold">{{ title }}</p>
            <p class="text-sm text-white/70">
              {{ currentIndex + 1 }} dari {{ photos.length }} foto
            </p>
          </div>
          
          <button
            @click="close"
            class="p-2 rounded-full bg-white/10 hover:bg-white/20 active-scale"
          >
            <X class="w-6 h-6 text-white" />
          </button>
        </div>

        <!-- Main Photo Area -->
        <div 
          ref="photoContainerRef"
          class="flex-1 relative overflow-hidden flex items-center justify-center"
          @touchstart="handleTouchStart"
          @touchmove="handleTouchMove"
          @touchend="handleTouchEnd"
        >
          <!-- Navigation Button Left -->
          <button
            v-if="currentIndex > 0"
            @click.stop="prevPhoto"
            class="absolute left-4 p-3 rounded-full bg-white/10 hover:bg-white/20 
                   active-scale z-10 hidden md:block"
          >
            <ChevronLeft class="w-8 h-8 text-white" />
          </button>

          <!-- Photo -->
          <Motion
            :key="currentIndex"
            :initial="slideDirection === 'left' ? { opacity: 0, x: 100 } : { opacity: 0, x: -100 }"
            :animate="{ opacity: 1, x: 0 }"
            :transition="{ duration: 0.25, ease: 'easeOut' }"
            class="w-full h-full flex items-center justify-center p-4"
            @click.stop
          >
            <img
              ref="imageRef"
              :src="currentPhoto"
              :alt="`Photo ${currentIndex + 1}`"
              class="max-w-full max-h-full object-contain select-none"
              :style="imageStyle"
              @load="handleImageLoad"
              @error="handleImageError"
              draggable="false"
            />
          </Motion>

          <!-- Navigation Button Right -->
          <button
            v-if="currentIndex < photos.length - 1"
            @click.stop="nextPhoto"
            class="absolute right-4 p-3 rounded-full bg-white/10 hover:bg-white/20 
                   active-scale z-10 hidden md:block"
          >
            <ChevronRight class="w-8 h-8 text-white" />
          </button>

          <!-- Loading Indicator -->
          <div 
            v-if="imageLoading"
            class="absolute inset-0 flex items-center justify-center"
          >
            <div class="w-10 h-10 border-3 border-white/20 border-t-white rounded-full animate-spin" />
          </div>

          <!-- Error State -->
          <div 
            v-if="imageError"
            class="absolute inset-0 flex flex-col items-center justify-center text-white/60"
          >
            <ImageOff class="w-16 h-16 mb-4" />
            <p>Gagal memuat foto</p>
          </div>
        </div>

        <!-- Thumbnail Strip -->
        <div 
          v-if="photos.length > 1"
          class="flex items-center justify-center gap-2 p-4 overflow-x-auto custom-scrollbar"
        >
          <button
            v-for="(photo, idx) in photos"
            :key="idx"
            @click="goToPhoto(idx)"
            class="flex-shrink-0 w-16 h-16 rounded-lg overflow-hidden border-2 
                   transition-all duration-150"
            :class="currentIndex === idx 
              ? 'border-white scale-110' 
              : 'border-white/30 opacity-60 hover:opacity-100'"
          >
            <img
              :src="photo"
              :alt="`Thumbnail ${idx + 1}`"
              class="w-full h-full object-cover"
            />
          </button>
        </div>

        <!-- Mobile Swipe Indicator -->
        <div 
          v-if="photos.length > 1"
          class="flex items-center justify-center gap-1.5 pb-4 md:hidden"
        >
          <div
            v-for="(_, idx) in photos"
            :key="idx"
            class="w-2 h-2 rounded-full transition-all duration-150"
            :class="currentIndex === idx ? 'bg-white w-4' : 'bg-white/40'"
          />
        </div>

        <!-- Zoom Controls -->
        <div class="absolute bottom-20 right-4 flex flex-col gap-2 z-10 hidden md:flex">
          <button
            @click="zoomIn"
            class="p-2 rounded-full bg-white/10 hover:bg-white/20 active-scale"
            :disabled="scale >= 3"
          >
            <ZoomIn class="w-5 h-5 text-white" />
          </button>
          <button
            @click="zoomOut"
            class="p-2 rounded-full bg-white/10 hover:bg-white/20 active-scale"
            :disabled="scale <= 1"
          >
            <ZoomOut class="w-5 h-5 text-white" />
          </button>
          <button
            v-if="scale !== 1"
            @click="resetZoom"
            class="p-2 rounded-full bg-white/10 hover:bg-white/20 active-scale"
          >
            <Maximize2 class="w-5 h-5 text-white" />
          </button>
        </div>
      </div>
    </Motion>
  </Teleport>
</template>

<script setup>
/**
 * MaterialPhotoViewer - Fullscreen photo viewer dengan swipe navigation
 * dan zoom capability untuk material photos
 */
import { ref, computed, watch, onMounted, onUnmounted } from 'vue'
import { Motion } from 'motion-v'
import { 
  X, 
  ChevronLeft, 
  ChevronRight, 
  ZoomIn, 
  ZoomOut,
  Maximize2,
  ImageOff
} from 'lucide-vue-next'

const props = defineProps({
  /**
   * V-model untuk visibility
   */
  modelValue: {
    type: Boolean,
    default: false
  },
  
  /**
   * Array of photo URLs
   */
  photos: {
    type: Array,
    default: () => []
  },
  
  /**
   * Initial photo index
   */
  initialIndex: {
    type: Number,
    default: 0
  },
  
  /**
   * Title untuk header
   */
  title: {
    type: String,
    default: 'Foto Material'
  }
})

const emit = defineEmits(['update:modelValue', 'close'])

// State
const currentIndex = ref(props.initialIndex)
const slideDirection = ref('left')
const imageLoading = ref(true)
const imageError = ref(false)
const scale = ref(1)
const photoContainerRef = ref(null)
const imageRef = ref(null)

// Touch handling
let touchStartX = 0
let touchStartY = 0
let isSwiping = false

/**
 * Computed untuk current photo
 */
const currentPhoto = computed(() => {
  return props.photos[currentIndex.value] || ''
})

/**
 * Computed untuk image style dengan zoom
 */
const imageStyle = computed(() => {
  return {
    transform: `scale(${scale.value})`,
    transition: 'transform 0.2s ease-out'
  }
})

/**
 * Watch initial index changes
 */
watch(() => props.initialIndex, (newVal) => {
  currentIndex.value = newVal
})

/**
 * Watch visibility untuk reset state
 */
watch(() => props.modelValue, (newVal) => {
  if (newVal) {
    currentIndex.value = props.initialIndex
    scale.value = 1
    imageLoading.value = true
    imageError.value = false
    
    // Lock body scroll
    document.body.style.overflow = 'hidden'
  } else {
    // Unlock body scroll
    document.body.style.overflow = ''
  }
})

/**
 * Close viewer
 */
const close = () => {
  emit('update:modelValue', false)
  emit('close')
  
  // Haptic feedback
  if ('vibrate' in navigator) {
    navigator.vibrate(10)
  }
}

/**
 * Handle backdrop click
 */
const handleBackdropClick = (e) => {
  if (e.target === e.currentTarget) {
    close()
  }
}

/**
 * Navigate to previous photo
 */
const prevPhoto = () => {
  if (currentIndex.value > 0) {
    slideDirection.value = 'right'
    currentIndex.value--
    resetImageState()
    
    if ('vibrate' in navigator) {
      navigator.vibrate(5)
    }
  }
}

/**
 * Navigate to next photo
 */
const nextPhoto = () => {
  if (currentIndex.value < props.photos.length - 1) {
    slideDirection.value = 'left'
    currentIndex.value++
    resetImageState()
    
    if ('vibrate' in navigator) {
      navigator.vibrate(5)
    }
  }
}

/**
 * Go to specific photo
 */
const goToPhoto = (index) => {
  if (index === currentIndex.value) return
  
  slideDirection.value = index > currentIndex.value ? 'left' : 'right'
  currentIndex.value = index
  resetImageState()
  
  if ('vibrate' in navigator) {
    navigator.vibrate(5)
  }
}

/**
 * Reset image state
 */
const resetImageState = () => {
  imageLoading.value = true
  imageError.value = false
  scale.value = 1
}

/**
 * Handle image load
 */
const handleImageLoad = () => {
  imageLoading.value = false
}

/**
 * Handle image error
 */
const handleImageError = () => {
  imageLoading.value = false
  imageError.value = true
}

/**
 * Zoom in
 */
const zoomIn = () => {
  if (scale.value < 3) {
    scale.value = Math.min(3, scale.value + 0.5)
  }
}

/**
 * Zoom out
 */
const zoomOut = () => {
  if (scale.value > 1) {
    scale.value = Math.max(1, scale.value - 0.5)
  }
}

/**
 * Reset zoom
 */
const resetZoom = () => {
  scale.value = 1
}

/**
 * Touch handlers untuk swipe navigation
 */
const handleTouchStart = (e) => {
  if (scale.value > 1) return // Disable swipe when zoomed
  
  touchStartX = e.touches[0].clientX
  touchStartY = e.touches[0].clientY
  isSwiping = false
}

const handleTouchMove = (e) => {
  if (scale.value > 1) return
  
  const touchX = e.touches[0].clientX
  const touchY = e.touches[0].clientY
  const diffX = touchStartX - touchX
  const diffY = Math.abs(touchStartY - touchY)
  
  // Only consider horizontal swipes
  if (Math.abs(diffX) > diffY && Math.abs(diffX) > 10) {
    isSwiping = true
    e.preventDefault()
  }
}

const handleTouchEnd = (e) => {
  if (scale.value > 1 || !isSwiping) return
  
  const touchX = e.changedTouches[0].clientX
  const diffX = touchStartX - touchX
  const threshold = 50 // Minimum swipe distance
  
  if (diffX > threshold) {
    nextPhoto()
  } else if (diffX < -threshold) {
    prevPhoto()
  }
  
  isSwiping = false
}

/**
 * Keyboard navigation
 */
const handleKeydown = (e) => {
  if (!props.modelValue) return
  
  switch (e.key) {
    case 'ArrowLeft':
      prevPhoto()
      break
    case 'ArrowRight':
      nextPhoto()
      break
    case 'Escape':
      close()
      break
    case '+':
    case '=':
      zoomIn()
      break
    case '-':
      zoomOut()
      break
    case '0':
      resetZoom()
      break
  }
}

onMounted(() => {
  document.addEventListener('keydown', handleKeydown)
})

onUnmounted(() => {
  document.removeEventListener('keydown', handleKeydown)
  document.body.style.overflow = ''
})
</script>
