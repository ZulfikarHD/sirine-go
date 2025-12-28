<template>
  <AppLayout>
    <!-- Header dengan Search dan Actions -->
    <Motion
      v-bind="entranceAnimations.fadeUp"
      class="glass-card p-6 rounded-2xl mb-6"
    >
      <div class="flex flex-col md:flex-row md:items-center md:justify-between space-y-4 md:space-y-0">
        <div>
          <h1 class="text-2xl font-bold text-gray-900">Manajemen User</h1>
          <p class="text-sm text-gray-600 mt-1">Kelola data user sistem</p>
        </div>
        <button 
          @click="openCreateModal"
          class="btn-primary flex items-center space-x-2 active-scale"
        >
          <UserPlus class="w-5 h-5" />
          <span>Tambah User Baru</span>
        </button>
      </div>

      <!-- Search dan Filters -->
      <div class="mt-6 grid grid-cols-1 md:grid-cols-4 gap-4">
        <!-- Search -->
        <div class="md:col-span-2">
          <div class="relative">
            <Search class="absolute left-3 top-1/2 transform -translate-y-1/2 w-5 h-5 text-gray-400" />
            <input
              v-model="searchQuery"
              @input="debouncedSearch"
              type="text"
              placeholder="Cari berdasarkan NIP atau Nama..."
              class="w-full pl-10 pr-4 py-2.5 rounded-xl border border-gray-200 focus:border-indigo-500 focus:ring-2 focus:ring-indigo-200"
            />
          </div>
        </div>

        <!-- Role Filter -->
        <select
          v-model="selectedRole"
          @change="applyFilters"
          class="px-4 py-2.5 rounded-xl border border-gray-200 focus:border-indigo-500 focus:ring-2 focus:ring-indigo-200"
        >
          <option value="">Semua Role</option>
          <option value="ADMIN">Admin</option>
          <option value="MANAGER">Manager</option>
          <option value="STAFF_KHAZWAL">Staff Khazwal</option>
          <option value="OPERATOR_CETAK">Operator Cetak</option>
          <option value="QC_INSPECTOR">QC Inspector</option>
          <option value="VERIFIKATOR">Verifikator</option>
          <option value="STAFF_KHAZKHIR">Staff Khazkhir</option>
        </select>

        <!-- Department Filter -->
        <select
          v-model="selectedDepartment"
          @change="applyFilters"
          class="px-4 py-2.5 rounded-xl border border-gray-200 focus:border-indigo-500 focus:ring-2 focus:ring-indigo-200"
        >
          <option value="">Semua Department</option>
          <option value="KHAZWAL">Khazwal</option>
          <option value="CETAK">Cetak</option>
          <option value="VERIFIKASI">Verifikasi</option>
          <option value="KHAZKHIR">Khazkhir</option>
        </select>
      </div>

      <!-- Active Filters -->
      <div v-if="userStore.hasFilters" class="mt-4 flex items-center space-x-2">
        <span class="text-sm text-gray-600">Filter aktif:</span>
        <button
          @click="clearFilters"
          class="px-3 py-1 text-xs bg-gray-100 text-gray-700 rounded-full hover:bg-gray-200 flex items-center space-x-1"
        >
          <X class="w-3 h-3" />
          <span>Clear All</span>
        </button>
      </div>
    </Motion>

    <!-- Users Table -->
    <Motion
      :initial="{ opacity: 0, y: 15 }"
      :animate="{ opacity: 1, y: 0 }"
      :transition="{ duration: 0.25, delay: 0.1, ease: 'easeOut' }"
      class="glass-card rounded-2xl overflow-hidden"
    >
      <!-- Loading State -->
      <div v-if="userStore.loading" class="p-12 text-center">
        <div class="inline-block animate-spin rounded-full h-12 w-12 border-4 border-indigo-500 border-t-transparent"></div>
        <p class="mt-4 text-gray-600">Memuat data users...</p>
      </div>

      <!-- Empty State -->
      <div v-else-if="userStore.isEmpty" class="p-12 text-center">
        <Users class="w-16 h-16 text-gray-300 mx-auto mb-4" />
        <h3 class="text-lg font-semibold text-gray-900 mb-2">Tidak ada user</h3>
        <p class="text-gray-600 mb-4">
          {{ userStore.hasFilters ? 'Tidak ada user yang sesuai dengan filter' : 'Belum ada user dalam sistem' }}
        </p>
        <button v-if="!userStore.hasFilters" @click="openCreateModal" class="btn-primary">
          Tambah User Pertama
        </button>
      </div>

      <!-- Table -->
      <div v-else class="overflow-x-auto">
        <table class="w-full">
          <thead class="bg-gray-50 border-b border-gray-200">
            <tr>
              <th class="px-6 py-4 text-left text-xs font-semibold text-gray-600 uppercase tracking-wider">User</th>
              <th class="px-6 py-4 text-left text-xs font-semibold text-gray-600 uppercase tracking-wider">Role</th>
              <th class="px-6 py-4 text-left text-xs font-semibold text-gray-600 uppercase tracking-wider">Department</th>
              <th class="px-6 py-4 text-left text-xs font-semibold text-gray-600 uppercase tracking-wider">Status</th>
              <th class="px-6 py-4 text-left text-xs font-semibold text-gray-600 uppercase tracking-wider">Shift</th>
              <th class="px-6 py-4 text-right text-xs font-semibold text-gray-600 uppercase tracking-wider">Actions</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-200">
            <Motion
              v-for="(user, index) in userStore.users" 
              :key="user.id"
              as="tr"
              :initial="{ opacity: 0, x: -10 }"
              :animate="{ opacity: 1, x: 0 }"
              :transition="{ duration: 0.2, delay: index * 0.03, ease: 'easeOut' }"
              class="hover:bg-gray-50"
            >
              <td class="px-6 py-4">
                <div class="flex items-center space-x-3">
                  <div class="w-10 h-10 rounded-full bg-linear-to-br from-indigo-500 to-fuchsia-600 flex items-center justify-center text-white font-semibold text-sm">
                    {{ getUserInitial(user.full_name) }}
                  </div>
                  <div>
                    <p class="font-semibold text-gray-900">{{ user.full_name }}</p>
                    <p class="text-sm text-gray-500">NIP: {{ user.nip }}</p>
                  </div>
                </div>
              </td>
              <td class="px-6 py-4">
                <RoleBadge :role="user.role" />
              </td>
              <td class="px-6 py-4">
                <span class="text-sm text-gray-900">{{ user.department }}</span>
              </td>
              <td class="px-6 py-4">
                <span 
                  :class="{
                    'bg-emerald-100 text-emerald-700': user.status === 'ACTIVE',
                    'bg-gray-100 text-gray-700': user.status === 'INACTIVE',
                    'bg-red-100 text-red-700': user.status === 'SUSPENDED'
                  }"
                  class="px-2.5 py-1 rounded-full text-xs font-semibold"
                >
                  {{ user.status }}
                </span>
              </td>
              <td class="px-6 py-4">
                <span class="text-sm text-gray-600">{{ user.shift }}</span>
              </td>
              <td class="px-6 py-4">
                <div class="flex items-center justify-end space-x-2">
                  <button
                    @click="openEditModal(user)"
                    class="p-2 text-indigo-600 hover:bg-indigo-50 rounded-lg active-scale"
                    title="Edit"
                  >
                    <Edit class="w-4 h-4" />
                  </button>
                  <button
                    @click="confirmDelete(user)"
                    class="p-2 text-red-600 hover:bg-red-50 rounded-lg active-scale"
                    title="Hapus"
                  >
                    <Trash2 class="w-4 h-4" />
                  </button>
                </div>
              </td>
            </Motion>
          </tbody>
        </table>
      </div>

      <!-- Pagination -->
      <div v-if="userStore.totalPages > 1" class="px-6 py-4 border-t border-gray-200 flex items-center justify-between">
        <div class="text-sm text-gray-600">
          Menampilkan {{ userStore.users.length }} dari {{ userStore.total }} users
        </div>
        <div class="flex items-center space-x-2">
          <button
            @click="goToPage(userStore.currentPage - 1)"
            :disabled="userStore.currentPage === 1"
            class="px-3 py-2 rounded-lg border border-gray-200 text-sm font-medium text-gray-700 hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            Previous
          </button>
          <div class="flex items-center space-x-1">
            <button
              v-for="page in visiblePages"
              :key="page"
              @click="goToPage(page)"
              :class="{
                'bg-indigo-600 text-white': page === userStore.currentPage,
                'text-gray-700 hover:bg-gray-50': page !== userStore.currentPage
              }"
              class="px-3 py-2 rounded-lg text-sm font-medium"
            >
              {{ page }}
            </button>
          </div>
          <button
            @click="goToPage(userStore.currentPage + 1)"
            :disabled="userStore.currentPage === userStore.totalPages"
            class="px-3 py-2 rounded-lg border border-gray-200 text-sm font-medium text-gray-700 hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            Next
          </button>
        </div>
      </div>
    </Motion>

    <!-- User Form Modal -->
    <UserFormModal
      v-if="showFormModal"
      :user="selectedUser"
      @close="closeFormModal"
      @success="handleFormSuccess"
    />

    <!-- Confirm Dialog untuk Delete -->
    <ConfirmDialog
      v-model="confirmDialog.isOpen.value"
      :title="confirmDialog.config.value.title"
      :message="confirmDialog.config.value.message"
      :variant="confirmDialog.config.value.variant"
      :loading="confirmDialog.loading.value"
      :show-warning="confirmDialog.config.value.showWarning"
      :confirm-text="confirmDialog.config.value.confirmText"
      @confirm="confirmDialog.handleConfirm"
      @cancel="confirmDialog.handleCancel"
    />

    <!-- Alert Dialog untuk Success/Error -->
    <AlertDialog
      v-model="alertDialog.isOpen.value"
      :title="alertDialog.config.value.title"
      :message="alertDialog.config.value.message"
      :variant="alertDialog.config.value.variant"
      :detail="alertDialog.config.value.detail"
      @close="alertDialog.handleClose"
    />
  </AppLayout>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { Motion } from 'motion-v'
import { useUserStore } from '../../../stores/user'
import { useRouter } from 'vue-router'
import AppLayout from '../../../components/layout/AppLayout.vue'
import RoleBadge from '../../../components/admin/RoleBadge.vue'
import UserFormModal from '../../../components/admin/UserFormModal.vue'
import { ConfirmDialog, AlertDialog } from '../../../components/common'
import { useConfirmDialog, useAlertDialog } from '../../../composables/useModal'
import { entranceAnimations } from '../../../composables/useMotion'
import { UserPlus, Search, Users, Edit, Trash2, X } from 'lucide-vue-next'

const userStore = useUserStore()
const router = useRouter()

// Modal composables
const confirmDialog = useConfirmDialog()
const alertDialog = useAlertDialog()

// State
const searchQuery = ref('')
const selectedRole = ref('')
const selectedDepartment = ref('')
const showFormModal = ref(false)
const selectedUser = ref(null)

// Debounced search
let searchTimeout = null
const debouncedSearch = () => {
  clearTimeout(searchTimeout)
  searchTimeout = setTimeout(() => {
    applyFilters()
  }, 300)
}

// Apply filters
const applyFilters = async () => {
  await userStore.setFilters({
    search: searchQuery.value,
    role: selectedRole.value,
    department: selectedDepartment.value
  })
}

// Clear filters
const clearFilters = () => {
  searchQuery.value = ''
  selectedRole.value = ''
  selectedDepartment.value = ''
  userStore.clearFilters()
}

// Pagination
const goToPage = async (page) => {
  if (page >= 1 && page <= userStore.totalPages) {
    await userStore.fetchUsers(page)
  }
}

const visiblePages = computed(() => {
  const current = userStore.currentPage
  const total = userStore.totalPages
  const pages = []
  
  let start = Math.max(1, current - 2)
  let end = Math.min(total, start + 4)
  
  if (end - start < 4) {
    start = Math.max(1, end - 4)
  }
  
  for (let i = start; i <= end; i++) {
    pages.push(i)
  }
  
  return pages
})

// Modal actions
const openCreateModal = () => {
  selectedUser.value = null
  showFormModal.value = true
}

const openEditModal = (user) => {
  selectedUser.value = user
  showFormModal.value = true
}

const closeFormModal = () => {
  showFormModal.value = false
  selectedUser.value = null
}

const handleFormSuccess = async () => {
  closeFormModal()
  await userStore.fetchUsers(userStore.currentPage)
}

// Delete user
const confirmDelete = async (user) => {
  const confirmed = await confirmDialog.confirm({
    title: 'Hapus User',
    message: `Yakin ingin menghapus user "${user.full_name}"?`,
    detail: `NIP: ${user.nip} | Role: ${user.role}`,
    variant: 'danger',
    confirmText: 'Ya, Hapus',
    showWarning: true,
    warningMessage: 'Data user yang dihapus tidak dapat dikembalikan.'
  })

  if (confirmed) {
    confirmDialog.loading.value = true
    try {
      await userStore.deleteUser(user.id)
      confirmDialog.close()
      
      await alertDialog.success('User berhasil dihapus!', {
        detail: `${user.full_name} telah dihapus dari sistem.`,
        autoDismiss: true,
        autoDismissDelay: 3000
      })
      
      await userStore.fetchUsers(userStore.currentPage)
    } catch (error) {
      confirmDialog.close()
      await alertDialog.error('Gagal menghapus user', {
        detail: error.message || 'Terjadi kesalahan saat menghapus user.'
      })
    }
  }
}

// Helper
const getUserInitial = (fullName) => {
  if (!fullName) return '?'
  return fullName
    .split(' ')
    .map(n => n[0])
    .slice(0, 2)
    .join('')
    .toUpperCase()
}

// Lifecycle
onMounted(async () => {
  await userStore.fetchUsers(1)
})
</script>
