<template>
  <div class="glass-card rounded-2xl overflow-hidden">
    <div class="overflow-x-auto">
      <table class="w-full divide-y divide-gray-200">
        <!-- Table Header -->
        <thead class="bg-gray-50">
          <tr>
            <th scope="col" class="px-4 py-3 text-left text-xs font-semibold text-gray-700 uppercase tracking-wider whitespace-nowrap">
              Waktu
            </th>
            <th scope="col" class="px-4 py-3 text-left text-xs font-semibold text-gray-700 uppercase tracking-wider whitespace-nowrap">
              User
            </th>
            <th scope="col" class="px-4 py-3 text-left text-xs font-semibold text-gray-700 uppercase tracking-wider whitespace-nowrap">
              Aksi
            </th>
            <th scope="col" class="px-4 py-3 text-left text-xs font-semibold text-gray-700 uppercase tracking-wider whitespace-nowrap">
              Entity
            </th>
            <th scope="col" class="px-4 py-3 text-right text-xs font-semibold text-gray-700 uppercase tracking-wider whitespace-nowrap">
              Detail
            </th>
          </tr>
        </thead>

        <!-- Table Body dengan staggered animation -->
        <tbody class="bg-white divide-y divide-gray-200">
          <!-- Single v-for dengan template untuk grouping main row + detail row -->
          <template v-for="(log, index) in logs" :key="log.id">
            <!-- Main Row dengan Motion-V animation -->
            <Motion
              as="tr"
              :initial="{ opacity: 0, x: -10 }"
              :animate="{ opacity: 1, x: 0 }"
              :transition="{ duration: 0.2, delay: index * 0.03, ease: 'easeOut' }"
              class="hover:bg-gray-50 transition-colors duration-150"
            >
              <!-- Timestamp -->
              <td class="px-4 py-3 whitespace-nowrap text-sm text-gray-600">
                <div>{{ formatDate(log.created_at) }}</div>
                <div class="text-xs text-gray-400">{{ formatTime(log.created_at) }}</div>
              </td>

              <!-- User Info -->
              <td class="px-4 py-3 whitespace-nowrap">
                <div class="flex items-center gap-2">
                  <div class="w-8 h-8 rounded-full bg-linear-to-br from-indigo-500 to-fuchsia-500 flex items-center justify-center text-white text-xs font-semibold">
                    {{ getUserInitials(log.user) }}
                  </div>
                  <div>
                    <div class="text-sm font-medium text-gray-900">{{ log.user?.full_name || 'Unknown' }}</div>
                    <div class="text-xs text-gray-500">NIP: {{ log.user?.nip || '-' }}</div>
                  </div>
                </div>
              </td>

              <!-- Action Badge -->
              <td class="px-4 py-3 whitespace-nowrap">
                <span
                  class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium"
                  :class="getActionBadgeClass(log.action)"
                >
                  {{ getActionLabel(log.action) }}
                </span>
              </td>

              <!-- Entity Info -->
              <td class="px-4 py-3 whitespace-nowrap">
                <div class="text-sm text-gray-900">{{ log.entity_type }}</div>
                <div v-if="log.entity_id" class="text-xs text-gray-500">ID: {{ log.entity_id }}</div>
              </td>

              <!-- Detail Button -->
              <td class="px-4 py-3 whitespace-nowrap text-right">
                <button
                  @click.stop="toggleExpanded(log.id)"
                  class="inline-flex items-center justify-center px-3 py-1.5 text-xs font-medium text-gray-600 hover:text-indigo-600 hover:bg-gray-100 rounded-lg transition-all duration-150 active-scale border border-gray-200"
                >
                  <svg
                    class="w-4 h-4 mr-1 transition-transform duration-200"
                    :class="{ 'rotate-180': expandedRows.has(log.id) }"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                  >
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
                  </svg>
                  {{ expandedRows.has(log.id) ? 'Tutup' : 'Detail' }}
                </button>
              </td>
            </Motion>

            <!-- Expanded Detail Row - muncul langsung setelah main row -->
            <tr v-if="expandedRows.has(log.id)" class="bg-gray-50/50">
              <td colspan="5" class="px-4 py-4">
                <div class="pl-4 md:pl-12">
                  <JsonDiff v-if="log.changes" :changes="parseChanges(log.changes)" />
                  <div v-else class="text-xs text-gray-500 italic mb-3">Tidak ada perubahan data</div>
                  
                  <!-- Additional Info -->
                  <div class="mt-3 pt-3 border-t border-gray-200 text-xs text-gray-500 grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div v-if="log.ip_address" class="flex items-center gap-2">
                      <svg class="w-3.5 h-3.5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 12a9 9 0 01-9 9m9-9a9 9 0 00-9-9m9 9H3m9 9a9 9 0 01-9-9m9 9c1.657 0 3-4.03 3-9s-1.343-9-3-9m0 18c-1.657 0-3-4.03-3-9s1.343-9 3-9m-9 9a9 9 0 019-9" />
                      </svg>
                      <span class="font-medium">IP:</span> {{ log.ip_address }}
                    </div>
                    <div v-if="log.user_agent" class="flex items-center gap-2">
                      <svg class="w-3.5 h-3.5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9.75 17L9 20l-1 1h8l-1-1-.75-3M3 13h18M5 17h14a2 2 0 002-2V5a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" />
                      </svg>
                      <span class="font-medium">Agent:</span> <span class="truncate max-w-xs">{{ log.user_agent }}</span>
                    </div>
                  </div>
                </div>
              </td>
            </tr>
          </template>

          <!-- Empty State -->
          <tr v-if="logs.length === 0">
            <td colspan="5" class="px-4 py-12 text-center">
              <svg class="w-12 h-12 mx-auto text-gray-300 mb-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
              </svg>
              <p class="text-sm text-gray-500">Tidak ada activity logs</p>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup>
import { reactive } from 'vue'
import { Motion } from 'motion-v'
import { format } from 'date-fns'
import { id as idLocale } from 'date-fns/locale'
import JsonDiff from './JsonDiff.vue'

// Props
defineProps({
  logs: {
    type: Array,
    required: true,
    default: () => [],
  },
})

/**
 * State untuk expanded rows menggunakan reactive object
 * untuk memastikan Vue dapat track perubahan state
 */
const expandedRows = reactive(new Set())

/**
 * toggleExpanded toggle expand/collapse untuk row detail
 * dengan haptic feedback dan force reactivity update
 */
const toggleExpanded = (logId) => {
  if (expandedRows.has(logId)) {
    expandedRows.delete(logId)
  } else {
    expandedRows.add(logId)
  }
  
  // Haptic feedback untuk mobile UX
  if ('vibrate' in navigator) {
    navigator.vibrate(10)
  }
}

/**
 * parseChanges parse JSON changes untuk display
 */
const parseChanges = (changes) => {
  try {
    if (typeof changes === 'string') {
      return JSON.parse(changes)
    }
    return changes
  } catch (error) {
    console.error('Error parsing changes:', error)
    return null
  }
}

/**
 * formatDate format date dengan Indonesian locale
 */
const formatDate = (timestamp) => {
  return format(new Date(timestamp), 'dd MMM yyyy', { locale: idLocale })
}

/**
 * formatTime format time
 */
const formatTime = (timestamp) => {
  return format(new Date(timestamp), 'HH:mm:ss')
}

/**
 * getUserInitials mendapatkan inisial dari nama user
 */
const getUserInitials = (user) => {
  if (!user?.full_name) return '??'
  const names = user.full_name.split(' ')
  if (names.length === 1) return names[0].substring(0, 2).toUpperCase()
  return (names[0][0] + names[names.length - 1][0]).toUpperCase()
}

/**
 * getActionLabel mengembalikan label untuk action type
 */
const getActionLabel = (action) => {
  const labels = {
    CREATE: 'Buat',
    UPDATE: 'Update',
    DELETE: 'Hapus',
    LOGIN: 'Login',
    LOGOUT: 'Logout',
    PASSWORD_CHANGE: 'Ganti Password',
  }
  return labels[action] || action
}

/**
 * getActionBadgeClass mengembalikan class untuk action badge
 * dengan color-coded berdasarkan action type
 */
const getActionBadgeClass = (action) => {
  const classes = {
    CREATE: 'bg-green-100 text-green-800',
    UPDATE: 'bg-blue-100 text-blue-800',
    DELETE: 'bg-red-100 text-red-800',
    LOGIN: 'bg-purple-100 text-purple-800',
    LOGOUT: 'bg-gray-100 text-gray-800',
    PASSWORD_CHANGE: 'bg-yellow-100 text-yellow-800',
  }
  return classes[action] || 'bg-gray-100 text-gray-800'
}
</script>
