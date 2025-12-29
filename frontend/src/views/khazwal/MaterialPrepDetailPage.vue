<template>
  <AppLayout>
    <!-- Header Section -->
    <Motion v-bind="entranceAnimations.fadeUp">
      <div class="mb-6 flex items-center gap-4">
        <!-- Back Button -->
        <button
          @click="$router.back()"
          class="p-2 rounded-xl hover:bg-gray-100 active-scale text-gray-600"
        >
          <ArrowLeft class="w-5 h-5" />
        </button>

        <div class="flex-1">
          <h1 class="text-3xl font-bold text-gray-900 mb-2">
            Detail Persiapan Material
          </h1>
          <p v-if="poDetail" class="text-gray-600">
            PO #{{ poDetail.po_number }} - {{ poDetail.obc_number }}
          </p>
          <p v-else class="text-gray-400">
            Loading...
          </p>
        </div>

        <!-- Status & Priority Badges -->
        <div v-if="poDetail" class="hidden sm:flex items-center gap-2">
          <PriorityBadge :priority="poDetail.priority" size="md" />
          <div :class="statusBadgeClass" class="flex items-center gap-2 px-4 py-2 rounded-xl">
            <div class="w-2 h-2 rounded-full animate-pulse" :class="statusDotClass" />
            <span class="text-sm font-semibold">{{ statusLabel }}</span>
          </div>
        </div>
      </div>
    </Motion>

    <!-- Loading State -->
    <div v-if="loading" class="space-y-6">
      <div class="glass-card rounded-2xl p-6 animate-pulse">
        <div class="h-6 bg-gray-200 rounded w-48 mb-4"></div>
        <div class="space-y-3">
          <div class="h-4 bg-gray-200 rounded w-full"></div>
          <div class="h-4 bg-gray-200 rounded w-3/4"></div>
        </div>
      </div>
    </div>

    <!-- Error State -->
    <Motion
      v-else-if="error"
      v-bind="entranceAnimations.fadeScale"
      class="glass-card rounded-2xl p-8 text-center bg-red-50 border-red-200"
    >
      <AlertTriangle class="w-12 h-12 text-red-600 mx-auto mb-4" />
      <h3 class="text-lg font-bold text-red-900 mb-2">
        Terjadi Kesalahan
      </h3>
      <p class="text-red-700 mb-4">
        {{ error }}
      </p>
      <button
        @click="fetchDetail"
        class="px-6 py-3 bg-red-600 text-white font-semibold rounded-xl hover:bg-red-700 active-scale"
      >
        <RefreshCw class="w-4 h-4 inline-block mr-2" />
        Coba Lagi
      </button>
    </Motion>

    <!-- Content Section -->
    <div v-else-if="poDetail" class="grid grid-cols-1 lg:grid-cols-3 gap-6">
      <!-- Main Content -->
      <div class="lg:col-span-2 space-y-6">
        <!-- PO Info Card -->
        <Motion
          :initial="{ opacity: 0, y: 15 }"
          :animate="{ opacity: 1, y: 0 }"
          :transition="{ duration: 0.25, delay: 0.1, ease: 'easeOut' }"
          class="glass-card rounded-2xl p-6"
        >
          <div class="flex items-center gap-3 mb-6">
            <Motion v-bind="iconAnimations.popIn">
              <div class="p-3 rounded-xl bg-gradient-to-br from-indigo-100 to-fuchsia-100">
                <FileText class="w-6 h-6 text-indigo-600" />
              </div>
            </Motion>
            <div>
              <h2 class="text-lg font-bold text-gray-900">Informasi PO</h2>
              <p class="text-sm text-gray-500">Detail production order</p>
            </div>
          </div>

          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <!-- PO Number -->
            <div>
              <p class="text-xs text-gray-500 mb-1">PO Number</p>
              <p class="text-sm font-semibold text-gray-900">{{ poDetail.po_number }}</p>
            </div>

            <!-- OBC Number -->
            <div>
              <p class="text-xs text-gray-500 mb-1">OBC Number</p>
              <p class="text-sm font-semibold text-gray-900">{{ poDetail.obc_number }}</p>
            </div>

            <!-- Product Name -->
            <div class="md:col-span-2">
              <p class="text-xs text-gray-500 mb-1">Nama Produk</p>
              <p class="text-sm font-semibold text-gray-900">{{ poDetail.product_name }}</p>
            </div>

            <!-- Customer Code -->
            <div>
              <p class="text-xs text-gray-500 mb-1">Customer Code</p>
              <p class="text-sm font-semibold text-gray-900">{{ poDetail.sap_customer_code }}</p>
            </div>

            <!-- Product Code -->
            <div>
              <p class="text-xs text-gray-500 mb-1">Product Code</p>
              <p class="text-sm font-semibold text-gray-900">{{ poDetail.sap_product_code }}</p>
            </div>

            <!-- Quantity -->
            <div>
              <p class="text-xs text-gray-500 mb-1">Quantity</p>
              <p class="text-sm font-semibold text-gray-900">{{ poDetail.quantity_ordered.toLocaleString() }} unit</p>
            </div>

            <!-- Target Lembar Besar -->
            <div>
              <p class="text-xs text-gray-500 mb-1">Target Lembar Besar</p>
              <p class="text-sm font-semibold text-gray-900">{{ poDetail.quantity_target_lembar_besar.toLocaleString() }}</p>
            </div>

            <!-- Order Date -->
            <div>
              <p class="text-xs text-gray-500 mb-1">Order Date</p>
              <p class="text-sm font-semibold text-gray-900">{{ formatDate(poDetail.order_date) }}</p>
            </div>

            <!-- Due Date -->
            <div>
              <p class="text-xs text-gray-500 mb-1">Due Date</p>
              <p :class="dueDateClass" class="text-sm font-semibold">
                {{ formatDate(poDetail.due_date) }}
                <span class="text-xs ml-1">({{ daysUntilDueText }})</span>
              </p>
            </div>
          </div>

          <!-- Notes if any -->
          <div v-if="poDetail.notes" class="mt-4 pt-4 border-t border-gray-200">
            <p class="text-xs text-gray-500 mb-1">Catatan</p>
            <p class="text-sm text-gray-700">{{ poDetail.notes }}</p>
          </div>
        </Motion>

        <!-- Material Requirements Card -->
        <Motion
          :initial="{ opacity: 0, y: 15 }"
          :animate="{ opacity: 1, y: 0 }"
          :transition="{ duration: 0.25, delay: 0.2, ease: 'easeOut' }"
          class="glass-card rounded-2xl p-6"
        >
          <div class="flex items-center gap-3 mb-6">
            <Motion v-bind="iconAnimations.popIn">
              <div class="p-3 rounded-xl bg-gradient-to-br from-emerald-100 to-teal-100">
                <ClipboardCheck class="w-6 h-6 text-emerald-600" />
              </div>
            </Motion>
            <div>
              <h2 class="text-lg font-bold text-gray-900">Kebutuhan Material</h2>
              <p class="text-sm text-gray-500">Plat, Kertas, dan Tinta</p>
            </div>
          </div>

          <div v-if="materialPrep" class="space-y-4">
            <!-- Plat Info -->
            <div class="p-4 bg-blue-50 border border-blue-200 rounded-xl">
              <div class="flex items-center gap-2 mb-2">
                <Box class="w-5 h-5 text-blue-600" />
                <p class="font-semibold text-blue-900">Plat Code</p>
              </div>
              <p class="text-lg font-bold text-blue-700">{{ materialPrep.sap_plat_code }}</p>
            </div>

            <!-- Kertas Blanko -->
            <div class="p-4 bg-amber-50 border border-amber-200 rounded-xl">
              <div class="flex items-center gap-2 mb-2">
                <FileStack class="w-5 h-5 text-amber-600" />
                <p class="font-semibold text-amber-900">Kertas Blanko</p>
              </div>
              <p class="text-lg font-bold text-amber-700">{{ materialPrep.kertas_blanko_quantity.toLocaleString() }} lembar</p>
            </div>

            <!-- Tinta Requirements -->
            <div class="p-4 bg-purple-50 border border-purple-200 rounded-xl">
              <div class="flex items-center gap-2 mb-3">
                <Droplet class="w-5 h-5 text-purple-600" />
                <p class="font-semibold text-purple-900">Tinta</p>
              </div>
              <div class="space-y-2">
                <div
                  v-for="(tinta, index) in tintaList"
                  :key="index"
                  class="flex items-center gap-2 text-sm"
                >
                  <div class="w-2 h-2 rounded-full bg-purple-500"></div>
                  <span class="text-purple-900">{{ tinta }}</span>
                </div>
              </div>
            </div>
          </div>
          <div v-else class="text-center py-8 text-gray-500">
            <Package class="w-12 h-12 text-gray-300 mx-auto mb-3" />
            <p>Belum ada data material preparation</p>
          </div>
        </Motion>
      </div>

      <!-- Sidebar -->
      <div class="space-y-6">
        <!-- Actions Card -->
        <Motion
          :initial="{ opacity: 0, y: 15 }"
          :animate="{ opacity: 1, y: 0 }"
          :transition="{ duration: 0.25, delay: 0.3, ease: 'easeOut' }"
          class="glass-card rounded-2xl p-6"
        >
          <h3 class="font-bold text-gray-900 mb-4">Actions</h3>
          
          <div class="space-y-3">
            <!-- Start Prep Button -->
            <button
              v-if="canStartPrep"
              @click="handleStartPrep"
              :disabled="startPrepLoading"
              class="w-full px-4 py-3 bg-gradient-to-r from-indigo-600 to-fuchsia-600 text-white font-semibold rounded-xl hover:shadow-lg active-scale disabled:opacity-50 disabled:cursor-not-allowed"
              style="transition: all 0.15s ease-out"
            >
              <Play v-if="!startPrepLoading" class="w-4 h-4 inline-block mr-2" />
              <RefreshCw v-else class="w-4 h-4 inline-block mr-2 animate-spin" />
              {{ startPrepLoading ? 'Memproses...' : 'Mulai Persiapan' }}
            </button>

            <!-- Continue Process Button (Sprint 3) -->
            <button
              v-if="isInProgress"
              @click="continueProcess"
              class="w-full px-4 py-3 bg-gradient-to-r from-emerald-600 to-teal-600 text-white font-semibold rounded-xl hover:shadow-lg active-scale"
              style="transition: all 0.15s ease-out"
            >
              <ArrowRight class="w-4 h-4 inline-block mr-2" />
              Lanjutkan Proses
            </button>

            <!-- In Progress Info -->
            <div v-if="isInProgress" class="p-4 bg-yellow-50 border border-yellow-200 rounded-xl">
              <div class="flex items-start gap-2">
                <Clock class="w-5 h-5 text-yellow-600 flex-shrink-0 mt-0.5" />
                <div>
                  <p class="text-sm font-semibold text-yellow-900 mb-1">
                    Sedang Dipersiapkan
                  </p>
                  <p class="text-xs text-yellow-700">
                    Material preparation sedang dalam proses
                  </p>
                </div>
              </div>
            </div>

            <!-- Completed Info -->
            <div v-if="isCompleted" class="p-4 bg-green-50 border border-green-200 rounded-xl">
              <div class="flex items-start gap-2">
                <CheckCircle class="w-5 h-5 text-green-600 flex-shrink-0 mt-0.5" />
                <div>
                  <p class="text-sm font-semibold text-green-900 mb-1">
                    Selesai
                  </p>
                  <p class="text-xs text-green-700">
                    Material preparation telah selesai
                  </p>
                </div>
              </div>
            </div>
          </div>

          <!-- Info Box -->
          <div class="mt-6 p-4 bg-blue-50 border border-blue-200 rounded-xl">
            <div class="flex items-start gap-2">
              <Info class="w-5 h-5 text-blue-600 flex-shrink-0 mt-0.5" />
              <p class="text-sm text-blue-800">
                Pastikan semua material tersedia sebelum memulai persiapan
              </p>
            </div>
          </div>
        </Motion>

        <!-- Timeline Card -->
        <Motion
          :initial="{ opacity: 0, y: 15 }"
          :animate="{ opacity: 1, y: 0 }"
          :transition="{ duration: 0.25, delay: 0.4, ease: 'easeOut' }"
          class="glass-card rounded-2xl p-6"
        >
          <h3 class="font-bold text-gray-900 mb-4">Timeline</h3>
          <div v-if="stageTracking && stageTracking.length > 0" class="space-y-3">
            <div
              v-for="(track, index) in stageTracking"
              :key="track.id"
              class="flex gap-3"
            >
              <div class="flex flex-col items-center">
                <div :class="[
                  'w-8 h-8 rounded-full flex items-center justify-center',
                  index === 0 ? 'bg-indigo-100 text-indigo-600' : 'bg-gray-100 text-gray-400'
                ]">
                  <div class="w-2 h-2 rounded-full bg-current"></div>
                </div>
                <div v-if="index < stageTracking.length - 1" class="w-0.5 h-full bg-gray-200 mt-2"></div>
              </div>
              <div class="flex-1 pb-4">
                <p class="text-sm font-semibold text-gray-900">{{ track.status }}</p>
                <p class="text-xs text-gray-500 mt-1">{{ formatDateTime(track.created_at) }}</p>
                <p v-if="track.notes" class="text-xs text-gray-600 mt-1">{{ track.notes }}</p>
              </div>
            </div>
          </div>
          <div v-else class="text-center py-4 text-gray-500">
            <Clock class="w-8 h-8 text-gray-300 mx-auto mb-2" />
            <p class="text-sm">Belum ada tracking</p>
          </div>
        </Motion>
      </div>
    </div>

    <!-- Confirm Dialog -->
    <ConfirmDialog
      v-model="confirmDialog.isOpen.value"
      :title="confirmDialog.config.value.title"
      :message="confirmDialog.config.value.message"
      :detail="confirmDialog.config.value.detail"
      :variant="confirmDialog.config.value.variant"
      :confirm-text="confirmDialog.config.value.confirmText"
      :loading="confirmDialog.loading.value"
      @confirm="confirmDialog.handleConfirm()"
      @cancel="confirmDialog.handleCancel()"
    />

    <!-- Alert Dialog -->
    <AlertDialog
      v-model="alertDialog.isOpen.value"
      :title="alertDialog.config.value.title"
      :message="alertDialog.config.value.message"
      :detail="alertDialog.config.value.detail"
      :variant="alertDialog.config.value.variant"
      @close="alertDialog.handleClose()"
    />
  </AppLayout>
</template>

<script setup>
/**
 * Material Prep Detail Page (Sprint 2: Full Implementation)
 * Halaman detail PO dengan real data dan Start Prep workflow
 * menggunakan ConfirmDialog untuk confirmation dan AlertDialog untuk feedback
 */
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Motion } from 'motion-v'
import { entranceAnimations, iconAnimations } from '@/composables/useMotion'
import { useKhazwalApi } from '@/composables/useKhazwalApi'
import { useConfirmDialog, useAlertDialog } from '@/composables/useModal'
import AppLayout from '@/components/layout/AppLayout.vue'
import PriorityBadge from '@/components/common/PriorityBadge.vue'
import ConfirmDialog from '@/components/common/ConfirmDialog.vue'
import AlertDialog from '@/components/common/AlertDialog.vue'
import { 
  Package, 
  ArrowLeft,
  ArrowRight,
  FileText,
  ClipboardCheck,
  Play,
  Info,
  Clock,
  RefreshCw,
  AlertTriangle,
  Box,
  FileStack,
  Droplet,
  CheckCircle
} from 'lucide-vue-next'

const route = useRoute()
const router = useRouter()
const khazwalApi = useKhazwalApi()
const confirmDialog = useConfirmDialog()
const alertDialog = useAlertDialog()

// State management
const loading = ref(false)
const error = ref(null)
const poDetail = ref(null)
const startPrepLoading = ref(false)

/**
 * Fetch PO detail dari API
 */
const fetchDetail = async () => {
  loading.value = true
  error.value = null

  try {
    const poId = route.params.id
    const response = await khazwalApi.getDetail(poId)

    if (response.success) {
      poDetail.value = response.data
    } else {
      error.value = response.message || 'Gagal memuat detail PO'
    }
  } catch (err) {
    console.error('Error fetching detail:', err)
    error.value = err.response?.data?.message || 'Gagal terhubung ke server'
  } finally {
    loading.value = false
  }
}

/**
 * Continue to process page (Sprint 3)
 */
const continueProcess = () => {
  router.push(`/khazwal/material-prep/${route.params.id}/process`)
  
  // Haptic feedback
  if ('vibrate' in navigator) {
    navigator.vibrate(10)
  }
}

/**
 * Handle Start Prep dengan confirmation dialog
 */
const handleStartPrep = async () => {
  const confirmed = await confirmDialog.confirm({
    title: 'Mulai Persiapan Material',
    message: 'Apakah Anda yakin ingin memulai persiapan material untuk PO ini?',
    detail: 'Status PO akan berubah menjadi "Sedang Dipersiapkan" dan proses tracking akan dimulai.',
    variant: 'info',
    confirmText: 'Ya, Mulai Sekarang',
    cancelText: 'Batal'
  })

  if (!confirmed) return

  startPrepLoading.value = true
  confirmDialog.loading.value = true

  try {
    const response = await khazwalApi.startPrep(route.params.id)

    if (response.success) {
      // Update local state dengan data baru
      poDetail.value = response.data
      
      confirmDialog.close()
      
      // Show success alert
      await alertDialog.success('Berhasil!', {
        title: 'Persiapan Material Dimulai',
        detail: 'Status PO telah diupdate dan tracking dimulai.',
        autoDismiss: true,
        autoDismissDelay: 2000
      })
      
      // Haptic feedback
      if ('vibrate' in navigator) {
        navigator.vibrate([30, 100, 30])
      }
    } else {
      throw new Error(response.message || 'Gagal memulai persiapan')
    }
  } catch (err) {
    console.error('Error starting prep:', err)
    confirmDialog.close()
    
    await alertDialog.error('Gagal Memulai Persiapan', {
      title: 'Terjadi Kesalahan',
      detail: err.response?.data?.message || err.message || 'Gagal memulai persiapan material',
    })
  } finally {
    startPrepLoading.value = false
    confirmDialog.loading.value = false
  }
}

/**
 * Computed: Material prep data
 */
const materialPrep = computed(() => {
  return poDetail.value?.khazwal_material_prep || null
})

/**
 * Computed: Stage tracking data
 */
const stageTracking = computed(() => {
  return poDetail.value?.stage_tracking || []
})

/**
 * Computed: Tinta list dari JSON
 */
const tintaList = computed(() => {
  if (!materialPrep.value?.tinta_requirements) return []
  
  try {
    const tinta = materialPrep.value.tinta_requirements
    if (Array.isArray(tinta)) return tinta
    if (typeof tinta === 'object') return Object.values(tinta)
    return []
  } catch (e) {
    return []
  }
})

/**
 * Computed: Can start prep (status harus WAITING_MATERIAL_PREP)
 */
const canStartPrep = computed(() => {
  return poDetail.value?.current_status === 'WAITING_MATERIAL_PREP'
})

/**
 * Computed: Is in progress
 */
const isInProgress = computed(() => {
  return poDetail.value?.current_status === 'MATERIAL_PREP_IN_PROGRESS'
})

/**
 * Computed: Is completed
 */
const isCompleted = computed(() => {
  return materialPrep.value?.status === 'COMPLETED'
})

/**
 * Computed: Status badge class
 */
const statusBadgeClass = computed(() => {
  if (!poDetail.value) return ''
  
  const statusMap = {
    'WAITING_MATERIAL_PREP': 'bg-blue-50 text-blue-700 border-blue-200',
    'MATERIAL_PREP_IN_PROGRESS': 'bg-yellow-50 text-yellow-700 border-yellow-200',
    'READY_FOR_CETAK': 'bg-green-50 text-green-700 border-green-200'
  }
  
  return statusMap[poDetail.value.current_status] || 'bg-gray-50 text-gray-700 border-gray-200'
})

/**
 * Computed: Status dot class
 */
const statusDotClass = computed(() => {
  if (!poDetail.value) return ''
  
  const statusMap = {
    'WAITING_MATERIAL_PREP': 'bg-blue-400',
    'MATERIAL_PREP_IN_PROGRESS': 'bg-yellow-400',
    'READY_FOR_CETAK': 'bg-green-400'
  }
  
  return statusMap[poDetail.value.current_status] || 'bg-gray-400'
})

/**
 * Computed: Status label
 */
const statusLabel = computed(() => {
  if (!poDetail.value) return ''
  
  const labelMap = {
    'WAITING_MATERIAL_PREP': 'Menunggu Persiapan',
    'MATERIAL_PREP_IN_PROGRESS': 'Sedang Dipersiapkan',
    'READY_FOR_CETAK': 'Siap Cetak'
  }
  
  return labelMap[poDetail.value.current_status] || poDetail.value.current_status
})

/**
 * Computed: Due date text class
 */
const dueDateClass = computed(() => {
  if (!poDetail.value) return 'text-gray-900'
  
  if (poDetail.value.is_past_due) return 'text-red-700'
  if (poDetail.value.days_until_due <= 3) return 'text-yellow-700'
  return 'text-emerald-700'
})

/**
 * Computed: Days until due text
 */
const daysUntilDueText = computed(() => {
  if (!poDetail.value) return ''
  
  if (poDetail.value.is_past_due) {
    return `Terlambat ${Math.abs(poDetail.value.days_until_due)} hari`
  }
  return `${poDetail.value.days_until_due} hari lagi`
})

/**
 * Format date untuk display yang user-friendly
 */
const formatDate = (dateString) => {
  try {
    const date = new Date(dateString)
    return date.toLocaleDateString('id-ID', { 
      day: 'numeric', 
      month: 'long',
      year: 'numeric'
    })
  } catch (e) {
    return dateString
  }
}

/**
 * Format datetime untuk timeline display
 */
const formatDateTime = (dateString) => {
  try {
    const date = new Date(dateString)
    return date.toLocaleDateString('id-ID', { 
      day: 'numeric', 
      month: 'short',
      year: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    })
  } catch (e) {
    return dateString
  }
}

// Fetch detail on mount
onMounted(() => {
  fetchDetail()
})
</script>
