<template>
  <AppLayout>
    <div class="max-w-4xl mx-auto px-4 py-6 space-y-6">
      <!-- Header -->
      <Motion v-bind="entranceAnimations.fadeUp">
        <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
          <div>
            <h1 class="text-2xl sm:text-3xl font-bold bg-gradient-to-r from-indigo-600 to-fuchsia-600 bg-clip-text text-transparent">
              Penghitungan
            </h1>
            <p class="text-gray-600 mt-1">Daftar PO menunggu penghitungan hasil cetak</p>
          </div>

          <!-- Refresh Button -->
          <button
            @click="refreshQueue"
            :disabled="isLoadingQueue"
            class="inline-flex items-center gap-2 px-4 py-2 bg-white border-2 border-gray-300 rounded-xl font-semibold text-gray-700 hover:border-indigo-500 hover:text-indigo-600 transition-all active-scale"
            :class="{ 'opacity-50 cursor-not-allowed': isLoadingQueue }"
          >
            <svg 
              class="w-5 h-5"
              :class="{ 'animate-spin': isLoadingQueue }"
              fill="none" 
              stroke="currentColor" 
              viewBox="0 0 24 24"
            >
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"/>
            </svg>
            <span class="hidden sm:inline">Refresh</span>
          </button>
        </div>
      </Motion>

      <!-- Stats Summary -->
      <Motion v-bind="entranceAnimations.fadeScale">
        <div class="grid grid-cols-1 sm:grid-cols-3 gap-4">
          <div class="glass-card p-4">
            <div class="flex items-center gap-3">
              <div class="w-12 h-12 rounded-full bg-gradient-to-br from-indigo-500 to-fuchsia-500 flex items-center justify-center">
                <svg class="w-6 h-6 text-white" fill="currentColor" viewBox="0 0 20 20">
                  <path d="M9 2a1 1 0 000 2h2a1 1 0 100-2H9z"/>
                  <path fill-rule="evenodd" d="M4 5a2 2 0 012-2 3 3 0 003 3h2a3 3 0 003-3 2 2 0 012 2v11a2 2 0 01-2 2H6a2 2 0 01-2-2V5zm3 4a1 1 0 000 2h.01a1 1 0 100-2H7zm3 0a1 1 0 000 2h3a1 1 0 100-2h-3zm-3 4a1 1 0 100 2h.01a1 1 0 100-2H7zm3 0a1 1 0 100 2h3a1 1 0 100-2h-3z" clip-rule="evenodd"/>
                </svg>
              </div>
              <div>
                <p class="text-xs text-gray-500">Total Queue</p>
                <p class="text-2xl font-bold text-gray-900">{{ queueMeta.total }}</p>
              </div>
            </div>
          </div>

          <div class="glass-card p-4">
            <div class="flex items-center gap-3">
              <div class="w-12 h-12 rounded-full bg-red-500 flex items-center justify-center">
                <svg class="w-6 h-6 text-white" fill="currentColor" viewBox="0 0 20 20">
                  <path fill-rule="evenodd" d="M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.213 2.98-1.742 2.98H4.42c-1.53 0-2.493-1.646-1.743-2.98l5.58-9.92zM11 13a1 1 0 11-2 0 1 1 0 012 0zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z" clip-rule="evenodd"/>
                </svg>
              </div>
              <div>
                <p class="text-xs text-gray-500">Overdue</p>
                <p class="text-2xl font-bold text-red-600">{{ queueMeta.overdue_count }}</p>
              </div>
            </div>
          </div>

          <div class="glass-card p-4">
            <div class="flex items-center gap-3">
              <div class="w-12 h-12 rounded-full bg-green-500 flex items-center justify-center">
                <svg class="w-6 h-6 text-white" fill="currentColor" viewBox="0 0 20 20">
                  <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm1-12a1 1 0 10-2 0v4a1 1 0 00.293.707l2.828 2.829a1 1 0 101.415-1.415L11 9.586V6z" clip-rule="evenodd"/>
                </svg>
              </div>
              <div>
                <p class="text-xs text-gray-500">On Time</p>
                <p class="text-2xl font-bold text-green-600">{{ queueMeta.total - queueMeta.overdue_count }}</p>
              </div>
            </div>
          </div>
        </div>
      </Motion>

      <!-- Filter Section (Optional - collapsed by default) -->
      <!-- Can be added later if needed -->

      <!-- Queue List -->
      <div v-if="isLoadingQueue" class="space-y-4">
        <LoadingSkeleton v-for="i in 3" :key="i" type="card" />
      </div>

      <div v-else-if="queue.length === 0" class="text-center py-12">
        <Motion v-bind="entranceAnimations.fadeScale">
          <div class="glass-card p-8">
            <div class="w-20 h-20 mx-auto mb-4 rounded-full bg-gradient-to-br from-gray-100 to-gray-200 flex items-center justify-center">
              <svg class="w-10 h-10 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"/>
              </svg>
            </div>
            <h3 class="text-lg font-bold text-gray-900 mb-2">Tidak Ada PO Menunggu</h3>
            <p class="text-gray-600">Semua PO sudah dihitung atau belum ada yang selesai cetak</p>
          </div>
        </Motion>
      </div>

      <div v-else class="space-y-4">
        <CountingQueueCard
          v-for="(item, index) in queue"
          :key="item.po_id"
          :item="item"
          :index="index"
          @click="navigateToWork(item)"
        />
      </div>

      <!-- Error Message -->
      <div v-if="error" class="glass-card p-4 bg-red-50 border-2 border-red-200">
        <div class="flex items-start gap-3">
          <svg class="w-6 h-6 text-red-600 flex-shrink-0" fill="currentColor" viewBox="0 0 20 20">
            <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd"/>
          </svg>
          <div>
            <p class="font-semibold text-red-900">Terjadi Kesalahan</p>
            <p class="text-sm text-red-700">{{ error }}</p>
          </div>
        </div>
      </div>
    </div>
  </AppLayout>
</template>

<script setup>
import { onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { storeToRefs } from 'pinia'
import { Motion } from 'motion-v'
import { entranceAnimations } from '@/composables/useMotion'
import { useCountingStore } from '@/stores/counting'
import AppLayout from '@/components/layout/AppLayout.vue'
import CountingQueueCard from '@/components/counting/CountingQueueCard.vue'
import LoadingSkeleton from '@/components/common/LoadingSkeleton.vue'

const router = useRouter()
const countingStore = useCountingStore()
const { queue, queueMeta, isLoadingQueue, error } = storeToRefs(countingStore)

onMounted(() => {
  countingStore.fetchQueue()
  
  // Setup auto-refresh every 30 seconds
  const intervalId = setInterval(() => {
    countingStore.fetchQueue()
  }, 30000)

  // Cleanup on unmount
  onBeforeUnmount(() => {
    clearInterval(intervalId)
  })
})

const refreshQueue = () => {
  countingStore.fetchQueue()
}

const navigateToWork = (item) => {
  router.push({
    name: 'counting-work',
    params: { poId: item.po_id }
  })
}
</script>
