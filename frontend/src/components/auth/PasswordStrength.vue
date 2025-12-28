<template>
  <div class="w-full space-y-2">
    <!-- Progress Bar -->
    <div class="w-full h-2 bg-gray-200 rounded-full overflow-hidden">
      <div 
        class="h-full transition-all duration-300 ease-out rounded-full"
        :class="{
          'bg-red-500': strengthClass === 'weak',
          'bg-yellow-500': strengthClass === 'medium',
          'bg-green-500': strengthClass === 'strong',
          'bg-emerald-600': strengthClass === 'very-strong'
        }"
        :style="{ width: strengthPercentage + '%' }"
      ></div>
    </div>

    <!-- Strength Label -->
    <div 
      class="text-sm font-medium text-center"
      :class="{
        'text-red-600': strengthClass === 'weak',
        'text-yellow-600': strengthClass === 'medium',
        'text-green-600': strengthClass === 'strong',
        'text-emerald-600': strengthClass === 'very-strong'
      }"
    >
      {{ strengthLabel }}
    </div>

    <!-- Requirements Checklist -->
    <div v-if="showRequirements" class="mt-3 space-y-2">
      <div 
        v-for="req in requirements" 
        :key="req.label"
        class="flex items-center gap-2 text-sm transition-colors duration-200"
        :class="req.met ? 'text-green-600' : 'text-gray-600'"
      >
        <svg 
          class="w-5 h-5 shrink-0 transition-colors duration-200"
          :class="req.met ? 'text-green-600' : 'text-gray-400'"
          fill="none" 
          viewBox="0 0 24 24" 
          stroke="currentColor"
        >
          <path 
            v-if="req.met"
            stroke-linecap="round" 
            stroke-linejoin="round" 
            stroke-width="2" 
            d="M5 13l4 4L19 7"
          />
          <circle 
            v-else
            cx="12" 
            cy="12" 
            r="10" 
            stroke-width="2"
          />
        </svg>
        <span class="leading-tight">{{ req.label }}</span>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'

/**
 * PasswordStrength component untuk menampilkan strength indicator
 * dengan visual progress bar dan requirements checklist
 */
const props = defineProps({
  password: {
    type: String,
    default: '',
  },
  showRequirements: {
    type: Boolean,
    default: true,
  },
})

/**
 * Calculate password strength berdasarkan requirements
 * Returns: 0 = Lemah, 1 = Sedang, 2 = Kuat, 3 = Sangat Kuat
 */
const calculateStrength = computed(() => {
  const pwd = props.password
  let strength = 0

  if (!pwd) return 0

  // Length check
  if (pwd.length >= 8) strength++
  if (pwd.length >= 12) strength++

  // Character diversity
  if (/[a-z]/.test(pwd)) strength++
  if (/[A-Z]/.test(pwd)) strength++
  if (/[0-9]/.test(pwd)) strength++
  if (/[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]/.test(pwd)) strength++

  // Normalize to 0-3 scale
  if (strength <= 2) return 0 // Lemah
  if (strength <= 3) return 1 // Sedang
  if (strength <= 4) return 2 // Kuat
  return 3 // Sangat Kuat
})

/**
 * Requirements checklist untuk password policy
 */
const requirements = computed(() => {
  const pwd = props.password
  return [
    {
      label: 'Minimal 8 karakter',
      met: pwd.length >= 8,
    },
    {
      label: 'Mengandung huruf besar',
      met: /[A-Z]/.test(pwd),
    },
    {
      label: 'Mengandung angka',
      met: /[0-9]/.test(pwd),
    },
    {
      label: 'Mengandung karakter spesial (!@#$%^&*, dll)',
      met: /[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]/.test(pwd),
    },
  ]
})

/**
 * Check apakah semua requirements terpenuhi
 */
const allRequirementsMet = computed(() => {
  return requirements.value.every(req => req.met)
})

/**
 * Strength label text
 */
const strengthLabel = computed(() => {
  const strength = calculateStrength.value
  const labels = ['Lemah', 'Sedang', 'Kuat', 'Sangat Kuat']
  return labels[strength] || 'Lemah'
})

/**
 * Strength class untuk styling
 */
const strengthClass = computed(() => {
  const strength = calculateStrength.value
  const classes = ['weak', 'medium', 'strong', 'very-strong']
  return classes[strength] || 'weak'
})

/**
 * Strength percentage untuk progress bar
 */
const strengthPercentage = computed(() => {
  const strength = calculateStrength.value
  const percentages = [25, 50, 75, 100]
  return percentages[strength] || 0
})

defineExpose({
  allRequirementsMet,
  calculateStrength,
})
</script>
