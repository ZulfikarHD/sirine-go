<template>
  <Teleport to="body">
    <Motion
      v-if="isVisible"
      :initial="{ opacity: 0 }"
      :animate="{ opacity: 1 }"
      :exit="{ opacity: 0 }"
      :transition="{ duration: 0.2, ease: 'easeOut' }"
      class="fixed inset-0 z-[9999] flex items-center justify-center bg-black/60 backdrop-blur-sm p-4"
      @click="handleBackdropClick"
    >
      <Motion
        :initial="{ opacity: 0, scale: 0.95, y: 20 }"
        :animate="{ opacity: 1, scale: 1, y: 0 }"
        :transition="springPresets.default"
        class="relative w-full max-w-md bg-white rounded-2xl shadow-2xl p-6 flex flex-col items-center text-center space-y-4"
        @click.stop
      >
        <!-- Icon -->
        <Motion v-bind="iconAnimations.popIn" class="w-16 h-16 rounded-full bg-gradient-to-br from-yellow-100 to-orange-100 flex items-center justify-center mb-2">
          <svg
            class="w-10 h-10 text-orange-600"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"
            />
          </svg>
        </Motion>

        <!-- Content -->
        <div class="space-y-2">
          <h3 class="text-xl font-bold text-gray-900">Sesi Anda Telah Berakhir</h3>
          <p class="text-sm text-gray-600 leading-relaxed">
            Untuk keamanan akun Anda, sesi login telah berakhir. Silakan login kembali untuk melanjutkan.
          </p>
        </div>

        <!-- Action Button -->
        <button
          class="w-full mt-4 px-6 py-3 rounded-xl bg-gradient-to-r from-indigo-600 to-fuchsia-600 text-white font-semibold hover:from-indigo-700 hover:to-fuchsia-700 focus:outline-none focus:ring-4 focus:ring-indigo-100 transition-all duration-200 active:scale-[0.97]"
          @click="handleLoginClick"
        >
          Login Kembali
        </button>
      </Motion>
    </Motion>
  </Teleport>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { Motion } from 'motion-v'
import { springPresets, iconAnimations } from '@/composables/useMotion'

/**
 * SessionExpired component untuk menampilkan modal saat token expired
 * dengan blocking overlay yang tidak bisa di-dismiss
 */
const props = defineProps({
  show: {
    type: Boolean,
    default: false,
  },
  allowDismiss: {
    type: Boolean,
    default: false,
  },
})

const emit = defineEmits(['login'])

const router = useRouter()
const isVisible = ref(props.show)

/**
 * Handle backdrop click - hanya dismiss jika allowDismiss = true
 */
const handleBackdropClick = () => {
  if (props.allowDismiss) {
    isVisible.value = false
  }
}

/**
 * Handle login button click
 */
const handleLoginClick = () => {
  isVisible.value = false
  emit('login')
  router.push('/login')
}

/**
 * Watch props.show untuk update visibility
 */
onMounted(() => {
  isVisible.value = props.show
})
</script>
