<template>
  <span 
    :class="badgeClass"
    class="px-2.5 py-1 rounded-full text-xs font-semibold inline-flex items-center space-x-1"
  >
    <component :is="roleIcon" class="w-3.5 h-3.5" />
    <span>{{ displayRole }}</span>
  </span>
</template>

<script setup>
import { computed } from 'vue'
import { Shield, Users, Package, Printer, CheckCircle, FileCheck, Archive } from 'lucide-vue-next'

const props = defineProps({
  role: {
    type: String,
    required: true
  }
})

/**
 * Badge class berdasarkan role dengan color-coded system
 * untuk visual hierarchy dan quick identification
 */
const badgeClass = computed(() => {
  const classes = {
    'ADMIN': 'bg-indigo-100 text-indigo-700',
    'MANAGER': 'bg-purple-100 text-purple-700',
    'STAFF_KHAZWAL': 'bg-blue-100 text-blue-700',
    'OPERATOR_CETAK': 'bg-fuchsia-100 text-fuchsia-700',
    'QC_INSPECTOR': 'bg-emerald-100 text-emerald-700',
    'VERIFIKATOR': 'bg-amber-100 text-amber-700',
    'STAFF_KHAZKHIR': 'bg-cyan-100 text-cyan-700'
  }
  return classes[props.role] || 'bg-gray-100 text-gray-700'
})

/**
 * Display role dengan format yang lebih readable
 */
const displayRole = computed(() => {
  const labels = {
    'ADMIN': 'Admin',
    'MANAGER': 'Manager',
    'STAFF_KHAZWAL': 'Staff Khazwal',
    'OPERATOR_CETAK': 'Operator Cetak',
    'QC_INSPECTOR': 'QC Inspector',
    'VERIFIKATOR': 'Verifikator',
    'STAFF_KHAZKHIR': 'Staff Khazkhir'
  }
  return labels[props.role] || props.role
})

/**
 * Icon untuk setiap role untuk visual distinction
 */
const roleIcon = computed(() => {
  const icons = {
    'ADMIN': Shield,
    'MANAGER': Users,
    'STAFF_KHAZWAL': Package,
    'OPERATOR_CETAK': Printer,
    'QC_INSPECTOR': CheckCircle,
    'VERIFIKATOR': FileCheck,
    'STAFF_KHAZKHIR': Archive
  }
  return icons[props.role] || Users
})
</script>
