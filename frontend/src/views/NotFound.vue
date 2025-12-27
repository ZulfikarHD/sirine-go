<template>
  <div class="min-h-screen flex items-center justify-center bg-gradient-to-br from-indigo-50 via-fuchsia-50 to-indigo-100 p-4">
    <div class="text-center">
      <div class="mb-8">
        <h1 class="text-9xl font-bold text-transparent bg-clip-text bg-gradient-to-r from-indigo-600 to-fuchsia-600">
          404
        </h1>
        <p class="text-2xl font-semibold text-gray-900 mt-4">Halaman Tidak Ditemukan</p>
        <p class="text-gray-600 mt-2">Maaf, halaman yang Anda cari tidak tersedia.</p>
      </div>
      <router-link
        to="/dashboard"
        class="inline-flex items-center space-x-2 px-6 py-3 rounded-xl font-semibold text-white bg-gradient-to-r from-indigo-600 to-fuchsia-600 hover:from-indigo-700 hover:to-fuchsia-700 transition-all active-scale shadow-lg"
      >
        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18" />
        </svg>
        <span>Kembali ke Dashboard</span>
      </router-link>
    </div>
  </div>
</template>

<script setup>
// No script needed for this simple page
</script>

<style scoped>
.active-scale {
  @apply transform transition-transform duration-150 active:scale-95;
}
</style>
