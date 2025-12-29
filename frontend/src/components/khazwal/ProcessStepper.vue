<template>
  <div class="w-full">
    <!-- Mobile: Vertical Stepper -->
    <div class="block sm:hidden">
      <div class="space-y-3">
        <Motion
          v-for="(step, index) in steps"
          :key="step.id"
          :initial="{ opacity: 0, x: -10 }"
          :animate="{ opacity: 1, x: 0 }"
          :transition="{ duration: 0.2, delay: index * 0.05, ease: 'easeOut' }"
        >
          <button
            @click="handleStepClick(index)"
            :disabled="!canNavigateToStep(index)"
            class="w-full flex items-center gap-3 p-3 rounded-xl transition-colors"
            :class="getStepMobileClass(index)"
          >
            <!-- Step Icon/Number -->
            <div
              class="flex-shrink-0 w-10 h-10 rounded-full flex items-center justify-center font-semibold transition-all"
              :class="getStepIconClass(index)"
            >
              <CheckCircle v-if="step.status === 'completed'" class="w-5 h-5" />
              <span v-else class="text-sm">{{ index + 1 }}</span>
            </div>

            <!-- Step Content -->
            <div class="flex-1 text-left">
              <p class="text-sm font-semibold" :class="getStepTextClass(index)">
                {{ step.title }}
              </p>
              <p v-if="step.description" class="text-xs mt-0.5" :class="getStepDescClass(index)">
                {{ step.description }}
              </p>
            </div>

            <!-- Arrow (only for active) -->
            <ChevronRight
              v-if="currentStep === index"
              class="w-5 h-5 text-indigo-600"
            />
          </button>
        </Motion>
      </div>
    </div>

    <!-- Desktop: Horizontal Stepper -->
    <div class="hidden sm:block">
      <div class="flex items-center justify-between">
        <Motion
          v-for="(step, index) in steps"
          :key="step.id"
          :initial="{ opacity: 0, y: 10 }"
          :animate="{ opacity: 1, y: 0 }"
          :transition="{ duration: 0.2, delay: index * 0.08, ease: 'easeOut' }"
          class="flex-1 relative"
        >
          <!-- Step Button -->
          <button
            @click="handleStepClick(index)"
            :disabled="!canNavigateToStep(index)"
            class="relative z-10 w-full flex flex-col items-center gap-2 p-3 rounded-xl transition-all active-scale"
            :class="getStepDesktopClass(index)"
          >
            <!-- Step Icon/Number -->
            <div
              class="w-12 h-12 rounded-full flex items-center justify-center font-bold transition-all"
              :class="getStepIconClass(index)"
            >
              <CheckCircle v-if="step.status === 'completed'" class="w-6 h-6" />
              <span v-else>{{ index + 1 }}</span>
            </div>

            <!-- Step Title -->
            <p class="text-sm font-semibold text-center" :class="getStepTextClass(index)">
              {{ step.title }}
            </p>

            <!-- Step Description -->
            <p v-if="step.description" class="text-xs text-center" :class="getStepDescClass(index)">
              {{ step.description }}
            </p>
          </button>

          <!-- Connector Line -->
          <div
            v-if="index < steps.length - 1"
            class="absolute top-9 left-[calc(50%+24px)] right-[calc(-50%+24px)] h-1 -z-0 transition-colors"
            :class="getConnectorClass(index)"
          ></div>
        </Motion>
      </div>
    </div>
  </div>
</template>

<script setup>
/**
 * ProcessStepper Component
 * Reusable stepper untuk workflow processes dengan iOS-inspired design
 * yang mendukung navigation ke completed steps dan visual progress indicators
 */
import { computed } from 'vue'
import { Motion } from 'motion-v'
import { CheckCircle, ChevronRight } from 'lucide-vue-next'

const props = defineProps({
  /**
   * Array steps dengan structure:
   * { id: string, title: string, description: string, status: 'pending'|'active'|'completed' }
   */
  steps: {
    type: Array,
    required: true,
    validator: (steps) => {
      return steps.every(step => 
        step.id && step.title && ['pending', 'active', 'completed'].includes(step.status)
      )
    }
  },
  /**
   * Current active step index (0-based)
   */
  currentStep: {
    type: Number,
    required: true
  },
  /**
   * Allow navigation ke completed steps
   */
  allowBackNavigation: {
    type: Boolean,
    default: true
  }
})

const emit = defineEmits(['step-click'])

/**
 * Check apakah user dapat navigate ke step tertentu
 */
const canNavigateToStep = (stepIndex) => {
  if (stepIndex === props.currentStep) return false
  if (!props.allowBackNavigation) return false
  
  // Hanya bisa navigate ke completed steps
  return props.steps[stepIndex].status === 'completed'
}

/**
 * Handle step click dengan emit event ke parent
 */
const handleStepClick = (stepIndex) => {
  if (canNavigateToStep(stepIndex)) {
    emit('step-click', stepIndex)
    
    // Haptic feedback
    if ('vibrate' in navigator) {
      navigator.vibrate(10)
    }
  }
}

/**
 * Get step mobile button class berdasarkan status
 */
const getStepMobileClass = (index) => {
  const step = props.steps[index]
  const isActive = props.currentStep === index
  const isCompleted = step.status === 'completed'
  const isDisabled = !canNavigateToStep(index) && !isActive
  
  if (isActive) {
    return 'bg-indigo-50 border-2 border-indigo-200'
  }
  if (isCompleted) {
    return 'bg-gray-50 border border-gray-200 hover:bg-gray-100'
  }
  return 'bg-white border border-gray-200 opacity-50'
}

/**
 * Get step desktop button class berdasarkan status
 */
const getStepDesktopClass = (index) => {
  const step = props.steps[index]
  const isActive = props.currentStep === index
  const isCompleted = step.status === 'completed'
  
  if (isActive) {
    return 'bg-gradient-to-br from-indigo-50 to-fuchsia-50'
  }
  if (isCompleted) {
    return 'hover:bg-gray-50'
  }
  return 'opacity-50 cursor-not-allowed'
}

/**
 * Get step icon class berdasarkan status
 */
const getStepIconClass = (index) => {
  const step = props.steps[index]
  const isActive = props.currentStep === index
  const isCompleted = step.status === 'completed'
  
  if (isActive) {
    return 'bg-gradient-to-br from-indigo-600 to-fuchsia-600 text-white shadow-lg'
  }
  if (isCompleted) {
    return 'bg-green-500 text-white'
  }
  return 'bg-gray-200 text-gray-400'
}

/**
 * Get step text class berdasarkan status
 */
const getStepTextClass = (index) => {
  const step = props.steps[index]
  const isActive = props.currentStep === index
  const isCompleted = step.status === 'completed'
  
  if (isActive) {
    return 'text-indigo-900'
  }
  if (isCompleted) {
    return 'text-gray-700'
  }
  return 'text-gray-400'
}

/**
 * Get step description class berdasarkan status
 */
const getStepDescClass = (index) => {
  const step = props.steps[index]
  const isActive = props.currentStep === index
  const isCompleted = step.status === 'completed'
  
  if (isActive) {
    return 'text-indigo-600'
  }
  if (isCompleted) {
    return 'text-gray-500'
  }
  return 'text-gray-400'
}

/**
 * Get connector line class berdasarkan status next step
 */
const getConnectorClass = (index) => {
  const nextStep = props.steps[index + 1]
  if (nextStep.status === 'completed' || nextStep.status === 'active') {
    return 'bg-gradient-to-r from-indigo-500 to-fuchsia-500'
  }
  return 'bg-gray-200'
}
</script>
