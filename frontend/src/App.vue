<script setup>
import { onMounted } from 'vue'
import { useAuthStore } from './stores/auth'

const authStore = useAuthStore()

// Restore auth state saat app dimuat
onMounted(() => {
  authStore.restoreAuth()
})
</script>

<template>
  <router-view />
</template>
