<template>
  <AppLayout>
    <!-- Header Section -->
    <Motion v-bind="entranceAnimations.fadeUp">
      <div class="mb-6 flex items-center gap-4">
        <!-- Back Button -->
        <button
          @click="handleBack"
          class="p-2 rounded-xl hover:bg-gray-100 active-scale text-gray-600"
        >
          <ArrowLeft class="w-5 h-5" />
        </button>

        <div class="flex-1">
          <h1 class="text-3xl font-bold text-gray-900 mb-2">
            Proses Persiapan Material
          </h1>
          <p v-if="poDetail" class="text-gray-600">
            PO #{{ poDetail.po_number }} - {{ poDetail.obc_number }}
          </p>
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
      <p class="text-red-700 mb-4">{{ error }}</p>
      <button
        @click="fetchDetail"
        class="px-6 py-3 bg-red-600 text-white font-semibold rounded-xl hover:bg-red-700 active-scale"
      >
        <RefreshCw class="w-4 h-4 inline-block mr-2" />
        Coba Lagi
      </button>
    </Motion>

    <!-- Main Content -->
    <div v-else-if="poDetail && materialPrep" class="space-y-6">
      <!-- Process Stepper -->
      <Motion
        :initial="{ opacity: 0, y: 15 }"
        :animate="{ opacity: 1, y: 0 }"
        :transition="{ duration: 0.25, delay: 0.1, ease: 'easeOut' }"
        class="glass-card rounded-2xl p-6"
      >
        <ProcessStepper
          :steps="processSteps"
          :current-step="currentStep"
          :allow-back-navigation="true"
          @step-click="navigateToStep"
        />
      </Motion>

      <!-- Step Content -->
      <Motion
        :key="currentStep"
        :initial="{ opacity: 0, x: 20 }"
        :animate="{ opacity: 1, x: 0 }"
        :transition="springPresets.snappy"
        class="glass-card rounded-2xl p-6"
      >
        <!-- Step 1: Confirm Plat -->
        <div v-if="currentStep === 0">
          <div class="mb-6">
            <h2 class="text-2xl font-bold text-gray-900 mb-2">
              Konfirmasi Plat
            </h2>
            <p class="text-gray-600">
              Scan barcode pada plat atau masukkan kode secara manual
            </p>
          </div>

          <!-- Expected Plat Code -->
          <div class="mb-6 p-4 bg-blue-50 border-2 border-blue-200 rounded-xl">
            <p class="text-sm text-blue-700 mb-2">Kode Plat yang Diharapkan:</p>
            <p class="text-2xl font-bold text-blue-900 font-mono">
              {{ materialPrep.sap_plat_code }}
            </p>
          </div>

          <!-- Barcode Scanner -->
          <BarcodeScanner
            :expected-code="materialPrep.sap_plat_code"
            :auto-start="false"
            :show-manual-input="true"
            @scan-success="handlePlatConfirmed"
            @scan-error="handleScanError"
          />
        </div>

        <!-- Step 2: Input Kertas -->
        <div v-else-if="currentStep === 1">
          <div class="mb-6">
            <h2 class="text-2xl font-bold text-gray-900 mb-2">
              Input Kertas Blanko
            </h2>
            <p class="text-gray-600">
              Masukkan jumlah kertas blanko yang diambil dari gudang
            </p>
          </div>

          <KertasInputForm
            :target-quantity="materialPrep.kertas_blanko_quantity"
            :loading="stepLoading"
            @submit="handleKertasSubmit"
          />
        </div>

        <!-- Step 3: Input Tinta -->
        <div v-else-if="currentStep === 2">
          <div class="mb-6">
            <h2 class="text-2xl font-bold text-gray-900 mb-2">
              Checklist Tinta
            </h2>
            <p class="text-gray-600">
              Centang tinta yang sudah diambil dan masukkan quantity
            </p>
          </div>

          <TintaChecklist
            :requirements="tintaRequirements"
            :loading="stepLoading"
            @submit="handleTintaSubmit"
          />
        </div>

        <!-- Step 4: Review & Finalize -->
        <div v-else-if="currentStep === 3 && !isFinalized">
          <div class="mb-6">
            <h2 class="text-2xl font-bold text-gray-900 mb-2">
              Review & Finalize
            </h2>
            <p class="text-gray-600">
              Review material yang sudah disiapkan dan selesaikan proses
            </p>
          </div>

          <!-- Summary Checklist -->
          <div class="space-y-4 mb-6">
            <!-- Plat Summary -->
            <div class="flex items-start gap-4 p-4 bg-green-50 border border-green-200 rounded-xl">
              <div class="flex-shrink-0 w-10 h-10 rounded-full bg-green-500 text-white flex items-center justify-center">
                <CheckCircle class="w-5 h-5" />
              </div>
              <div class="flex-1 min-w-0">
                <p class="text-sm font-medium text-gray-500">Plat Dikonfirmasi</p>
                <p class="text-lg font-bold text-gray-900 font-mono">{{ materialPrep?.sap_plat_code }}</p>
                <p v-if="materialPrep?.plat_retrieved_at" class="text-xs text-gray-500 mt-1">
                  {{ formatDateTime(materialPrep.plat_retrieved_at) }}
                </p>
              </div>
            </div>

            <!-- Kertas Summary -->
            <div class="flex items-start gap-4 p-4 bg-green-50 border border-green-200 rounded-xl">
              <div class="flex-shrink-0 w-10 h-10 rounded-full bg-green-500 text-white flex items-center justify-center">
                <FileText class="w-5 h-5" />
              </div>
              <div class="flex-1 min-w-0">
                <p class="text-sm font-medium text-gray-500">Kertas Blanko</p>
                <p class="text-lg font-bold text-gray-900">
                  {{ materialPrep?.kertas_blanko_actual || 0 }} lembar
                  <span v-if="materialPrep?.kertas_blanko_variance && materialPrep.kertas_blanko_variance !== 0" 
                        :class="materialPrep.kertas_blanko_variance > 0 ? 'text-green-600' : 'text-red-600'"
                        class="text-sm font-normal ml-2">
                    ({{ materialPrep.kertas_blanko_variance > 0 ? '+' : '' }}{{ materialPrep.kertas_blanko_variance }})
                  </span>
                </p>
                <p v-if="materialPrep?.kertas_blanko_variance_reason" class="text-xs text-gray-500 mt-1">
                  Catatan: {{ materialPrep.kertas_blanko_variance_reason }}
                </p>
              </div>
            </div>

            <!-- Tinta Summary -->
            <div class="flex items-start gap-4 p-4 bg-green-50 border border-green-200 rounded-xl">
              <div class="flex-shrink-0 w-10 h-10 rounded-full bg-green-500 text-white flex items-center justify-center">
                <Palette class="w-5 h-5" />
              </div>
              <div class="flex-1 min-w-0">
                <p class="text-sm font-medium text-gray-500">Tinta</p>
                <p class="text-lg font-bold text-gray-900">{{ parsedTintaActual.length }} warna</p>
                <div v-if="parsedTintaActual.length > 0" class="flex flex-wrap gap-2 mt-2">
                  <span
                    v-for="tinta in parsedTintaActual"
                    :key="tinta.color"
                    class="px-2 py-1 rounded-full text-xs font-medium bg-white border border-gray-200 text-gray-700"
                  >
                    {{ tinta.color }}: {{ tinta.quantity }}
                  </span>
                </div>
              </div>
            </div>
          </div>

          <!-- Photo Upload Section -->
          <div class="mb-6">
            <label class="block text-sm font-semibold text-gray-700 mb-3">
              Foto Material (Opsional)
            </label>
            <PhotoUploader
              v-model="uploadedPhotos"
              :max-photos="5"
              :max-size-k-b="2048"
            />
          </div>

          <!-- Notes Section -->
          <div class="mb-6">
            <label class="block text-sm font-semibold text-gray-700 mb-2">
              Catatan Tambahan (Opsional)
            </label>
            <textarea
              v-model="finalizeNotes"
              rows="3"
              placeholder="Tambahkan catatan jika diperlukan..."
              class="input-field resize-none"
            ></textarea>
          </div>

          <!-- Finalize Button -->
          <button
            @click="handleFinalize"
            :disabled="stepLoading"
            class="w-full py-4 bg-gradient-to-r from-indigo-600 to-fuchsia-600 text-white font-bold rounded-xl hover:shadow-lg active-scale disabled:opacity-50 disabled:cursor-not-allowed"
          >
            <span v-if="stepLoading" class="flex items-center justify-center gap-2">
              <Loader2 class="w-5 h-5 animate-spin" />
              Memproses...
            </span>
            <span v-else class="flex items-center justify-center gap-2">
              <Send class="w-5 h-5" />
              Selesai & Kirim ke Unit Cetak
            </span>
          </button>
        </div>

        <!-- Success Screen -->
        <div v-else-if="currentStep === 3 && isFinalized">
          <Motion v-bind="entranceAnimations.fadeScale" class="text-center py-8">
            <!-- Success Icon dengan Animation -->
            <Motion v-bind="iconAnimations.popIn" class="mb-6">
              <div class="w-24 h-24 mx-auto rounded-full bg-gradient-to-br from-green-400 to-emerald-600 flex items-center justify-center shadow-lg">
                <CheckCircle class="w-12 h-12 text-white" />
              </div>
            </Motion>

            <!-- Success Message -->
            <h2 class="text-2xl font-bold text-gray-900 mb-2">
              Berhasil Diselesaikan!
            </h2>
            <p class="text-gray-600 mb-6">
              Material untuk <strong>PO #{{ finalizeResult?.obc_number }}</strong> telah siap dan notifikasi telah dikirim ke Unit Cetak.
            </p>

            <!-- Duration Display -->
            <Motion
              :initial="{ opacity: 0, y: 10 }"
              :animate="{ opacity: 1, y: 0 }"
              :transition="{ duration: 0.3, delay: 0.2 }"
              class="duration-card"
            >
              <Clock class="w-6 h-6 text-indigo-600" />
              <div>
                <p class="text-sm text-gray-500">Waktu Persiapan</p>
                <p class="text-2xl font-bold text-gray-900">
                  {{ formatDuration(finalizeResult?.duration_minutes || 0) }}
                </p>
              </div>
            </Motion>

            <!-- Completion Details -->
            <Motion
              :initial="{ opacity: 0, y: 10 }"
              :animate="{ opacity: 1, y: 0 }"
              :transition="{ duration: 0.3, delay: 0.3 }"
              class="completion-details"
            >
              <div class="completion-detail-item">
                <span class="text-gray-500">Diselesaikan oleh</span>
                <span class="font-medium">{{ finalizeResult?.prepared_by_name || '-' }}</span>
              </div>
              <div class="completion-detail-item">
                <span class="text-gray-500">Waktu selesai</span>
                <span class="font-medium">{{ finalizeResult?.completed_at || '-' }}</span>
              </div>
              <div v-if="finalizeResult?.photos_count > 0" class="completion-detail-item">
                <span class="text-gray-500">Foto dilampirkan</span>
                <span class="font-medium">{{ finalizeResult.photos_count }} foto</span>
              </div>
            </Motion>

            <!-- Action Buttons -->
            <Motion
              :initial="{ opacity: 0, y: 10 }"
              :animate="{ opacity: 1, y: 0 }"
              :transition="{ duration: 0.3, delay: 0.4 }"
              class="flex flex-col gap-3 mt-8"
            >
              <button
                @click="$router.push('/khazwal/material-prep')"
                class="w-full py-4 bg-gradient-to-r from-indigo-600 to-fuchsia-600 text-white font-bold rounded-xl hover:shadow-lg active-scale"
              >
                <ArrowLeft class="w-5 h-5 inline-block mr-2" />
                Kembali ke Queue
              </button>
              <button
                @click="$router.push(`/khazwal/material-prep/${$route.params.id}`)"
                class="w-full py-3 bg-white border-2 border-gray-200 text-gray-700 font-semibold rounded-xl hover:bg-gray-50 active-scale"
              >
                Lihat Detail PO
              </button>
            </Motion>
          </Motion>
        </div>
      </Motion>

      <!-- Progress Info Card (hidden when finalized) -->
      <Motion
        v-if="!isFinalized"
        :initial="{ opacity: 0, y: 15 }"
        :animate="{ opacity: 1, y: 0 }"
        :transition="{ duration: 0.25, delay: 0.2, ease: 'easeOut' }"
        class="glass-card rounded-2xl p-4"
      >
        <div class="flex items-center gap-3">
          <Info class="w-5 h-5 text-blue-600" />
          <p class="text-sm text-gray-600">
            <strong>Step {{ currentStep + 1 }} dari {{ processSteps.length }}</strong> - 
            {{ processSteps[currentStep].description }}
          </p>
        </div>
      </Motion>
    </div>

    <!-- Alert Dialog -->
    <AlertDialog
      v-model="alertDialog.isOpen.value"
      :title="alertDialog.config.value.title"
      :message="alertDialog.config.value.message"
      :detail="alertDialog.config.value.detail"
      :variant="alertDialog.config.value.variant"
      @close="alertDialog.handleClose()"
    />

    <!-- Confirm Dialog -->
    <ConfirmDialog
      v-model="confirmDialog.isOpen.value"
      :title="confirmDialog.config.value.title"
      :message="confirmDialog.config.value.message"
      :detail="confirmDialog.config.value.detail"
      :variant="confirmDialog.config.value.variant"
      :confirm-text="confirmDialog.config.value.confirmText"
      :cancel-text="confirmDialog.config.value.cancelText"
      :loading="confirmDialog.loading.value"
      @confirm="confirmDialog.handleConfirm()"
      @cancel="confirmDialog.handleCancel()"
    />
  </AppLayout>
</template>

<script setup>
/**
 * MaterialPrepProcessPage (Sprint 3 & 4)
 * Halaman workflow untuk material preparation process dengan stepper
 * yang mencakup: Confirm Plat → Input Kertas → Input Tinta → Review & Finalize
 */
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Motion } from 'motion-v'
import { entranceAnimations, springPresets, iconAnimations } from '@/composables/useMotion'
import { useKhazwalApi } from '@/composables/useKhazwalApi'
import { useAlertDialog, useConfirmDialog } from '@/composables/useModal'
import AppLayout from '@/components/layout/AppLayout.vue'
import ProcessStepper from '@/components/khazwal/ProcessStepper.vue'
import BarcodeScanner from '@/components/common/BarcodeScanner.vue'
import KertasInputForm from '@/components/khazwal/KertasInputForm.vue'
import TintaChecklist from '@/components/khazwal/TintaChecklist.vue'
import PhotoUploader from '@/components/common/PhotoUploader.vue'
import AlertDialog from '@/components/common/AlertDialog.vue'
import ConfirmDialog from '@/components/common/ConfirmDialog.vue'
import { 
  ArrowLeft, 
  RefreshCw, 
  AlertTriangle, 
  Info,
  CheckCircle,
  FileText,
  Palette,
  Send,
  Loader2,
  Clock
} from 'lucide-vue-next'

const route = useRoute()
const router = useRouter()
const khazwalApi = useKhazwalApi()
const alertDialog = useAlertDialog()
const confirmDialog = useConfirmDialog()

// State management
const loading = ref(false)
const error = ref(null)
const poDetail = ref(null)
const currentStep = ref(0)
const stepLoading = ref(false)

// Step 4 (Finalize) states
const uploadedPhotos = ref([])
const finalizeNotes = ref('')
const isFinalized = ref(false)
const finalizeResult = ref(null)

// Process steps configuration
const processSteps = ref([
  {
    id: 'plat',
    title: 'Konfirmasi Plat',
    description: 'Scan barcode plat',
    status: 'active'
  },
  {
    id: 'kertas',
    title: 'Input Kertas',
    description: 'Input jumlah kertas blanko',
    status: 'pending'
  },
  {
    id: 'tinta',
    title: 'Input Tinta',
    description: 'Checklist tinta per warna',
    status: 'pending'
  },
  {
    id: 'review',
    title: 'Review',
    description: 'Review & finalize (Sprint 5)',
    status: 'pending'
  }
])

/**
 * Computed: Material prep data
 */
const materialPrep = computed(() => {
  return poDetail.value?.khazwal_material_prep || null
})

/**
 * Computed: Tinta requirements list
 */
const tintaRequirements = computed(() => {
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
 * Computed: Parsed tinta actual untuk summary display
 */
const parsedTintaActual = computed(() => {
  if (!materialPrep.value?.tinta_actual) return []
  
  try {
    const tinta = materialPrep.value.tinta_actual
    if (Array.isArray(tinta)) return tinta
    if (typeof tinta === 'string') return JSON.parse(tinta)
    if (typeof tinta === 'object') return Object.values(tinta)
    return []
  } catch (e) {
    return []
  }
})

/**
 * Format datetime untuk display
 */
const formatDateTime = (datetime) => {
  if (!datetime) return '-'
  try {
    const date = new Date(datetime)
    return date.toLocaleString('id-ID', {
      day: '2-digit',
      month: 'short',
      year: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    })
  } catch {
    return datetime
  }
}

/**
 * Format duration dalam menit ke format yang readable
 */
const formatDuration = (minutes) => {
  if (!minutes || minutes === 0) return '0 menit'
  
  const hours = Math.floor(minutes / 60)
  const mins = minutes % 60
  
  if (hours === 0) {
    return `${mins} menit`
  } else if (mins === 0) {
    return `${hours} jam`
  } else {
    return `${hours} jam ${mins} menit`
  }
}

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
      
      // Validate status
      if (poDetail.value.current_status !== 'MATERIAL_PREP_IN_PROGRESS') {
        error.value = 'PO tidak dalam status yang valid untuk proses persiapan'
      }
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
 * Navigate to specific step (untuk back navigation)
 */
const navigateToStep = (stepIndex) => {
  currentStep.value = stepIndex
  
  // Update step status
  processSteps.value.forEach((step, index) => {
    if (index < stepIndex) {
      step.status = 'completed'
    } else if (index === stepIndex) {
      step.status = 'active'
    } else {
      step.status = 'pending'
    }
  })
}

/**
 * Move to next step
 */
const nextStep = () => {
  if (currentStep.value < processSteps.value.length - 1) {
    processSteps.value[currentStep.value].status = 'completed'
    currentStep.value++
    processSteps.value[currentStep.value].status = 'active'
  }
}

/**
 * Handle plat confirmation success
 */
const handlePlatConfirmed = async (platCode) => {
  stepLoading.value = true

  try {
    const response = await khazwalApi.confirmPlat(materialPrep.value.id, platCode)

    if (response.success) {
      await alertDialog.success('Plat Berhasil Dikonfirmasi!', {
        title: 'Berhasil',
        detail: `Kode plat ${platCode} telah terverifikasi.`,
        autoDismiss: true,
        autoDismissDelay: 2000
      })
      
      // Move to next step
      nextStep()
    } else {
      throw new Error(response.message || 'Gagal mengkonfirmasi plat')
    }
  } catch (err) {
    console.error('Error confirming plat:', err)
    await alertDialog.error('Gagal Mengkonfirmasi Plat', {
      title: 'Terjadi Kesalahan',
      detail: err.response?.data?.message || err.message
    })
  } finally {
    stepLoading.value = false
  }
}

/**
 * Handle scan error
 */
const handleScanError = async (errorData) => {
  await alertDialog.error('Barcode Tidak Sesuai', {
    title: 'Scan Error',
    detail: errorData.message || 'Kode barcode tidak sesuai dengan yang diharapkan'
  })
}

/**
 * Handle kertas input submit
 */
const handleKertasSubmit = async (data) => {
  stepLoading.value = true

  try {
    const response = await khazwalApi.updateKertas(
      materialPrep.value.id,
      data.actualQty,
      data.varianceReason
    )

    if (response.success) {
      await alertDialog.success('Kertas Blanko Berhasil Diupdate!', {
        title: 'Berhasil',
        detail: `Jumlah actual: ${data.actualQty} lembar`,
        autoDismiss: true,
        autoDismissDelay: 2000
      })
      
      // Move to next step
      nextStep()
    } else {
      throw new Error(response.message || 'Gagal mengupdate kertas blanko')
    }
  } catch (err) {
    console.error('Error updating kertas:', err)
    await alertDialog.error('Gagal Mengupdate Kertas', {
      title: 'Terjadi Kesalahan',
      detail: err.response?.data?.message || err.message
    })
  } finally {
    stepLoading.value = false
  }
}

/**
 * Handle tinta checklist submit
 */
const handleTintaSubmit = async (tintaActual) => {
  stepLoading.value = true

  try {
    const response = await khazwalApi.updateTinta(materialPrep.value.id, tintaActual)

    if (response.success) {
      await alertDialog.success('Tinta Berhasil Diupdate!', {
        title: 'Berhasil',
        detail: `${tintaActual.length} warna tinta telah dicatat.`,
        autoDismiss: true,
        autoDismissDelay: 2000
      })
      
      // Move to next step
      nextStep()
    } else {
      throw new Error(response.message || 'Gagal mengupdate tinta')
    }
  } catch (err) {
    console.error('Error updating tinta:', err)
    await alertDialog.error('Gagal Mengupdate Tinta', {
      title: 'Terjadi Kesalahan',
      detail: err.response?.data?.message || err.message
    })
  } finally {
    stepLoading.value = false
  }
}

/**
 * Handle finalize material preparation
 */
const handleFinalize = async () => {
  // Show confirmation dialog
  const confirmed = await confirmDialog.confirm({
    title: 'Selesaikan Persiapan',
    message: 'Yakin ingin menyelesaikan persiapan material?',
    detail: 'Material akan dikirim ke Unit Cetak untuk proses selanjutnya.',
    variant: 'default',
    confirmText: 'Ya, Selesaikan',
    cancelText: 'Batal'
  })

  if (!confirmed) return

  stepLoading.value = true

  try {
    // Get base64 photos array
    const photosBase64 = uploadedPhotos.value.map(photo => photo.base64)

    // Call API finalize dengan material prep ID
    const prepId = materialPrep.value?.id
    if (!prepId) {
      throw new Error('Material prep ID tidak ditemukan')
    }

    const response = await khazwalApi.finalize(prepId, {
      photos: photosBase64,
      notes: finalizeNotes.value
    })

    if (response.success) {
      // Store result untuk success screen
      finalizeResult.value = response.data
      isFinalized.value = true

      // Mark step as completed
      processSteps.value[currentStep.value].status = 'completed'

      // Haptic feedback untuk success
      if ('vibrate' in navigator) {
        navigator.vibrate([30, 50, 30])
      }
    } else {
      throw new Error(response.message || 'Gagal menyelesaikan persiapan material')
    }
  } catch (err) {
    console.error('Error finalizing:', err)
    await alertDialog.error('Gagal Menyelesaikan Persiapan', {
      title: 'Terjadi Kesalahan',
      detail: err.response?.data?.message || err.message
    })
  } finally {
    stepLoading.value = false
  }
}

/**
 * Handle back navigation dengan confirmation
 */
const handleBack = () => {
  if (currentStep.value > 0) {
    // Navigate to previous step jika masih di tengah proses
    navigateToStep(currentStep.value - 1)
  } else {
    // Go back to detail page jika di step pertama
    router.push(`/khazwal/material-prep/${route.params.id}`)
  }
}

// Fetch detail on mount
onMounted(() => {
  fetchDetail()
})
</script>
