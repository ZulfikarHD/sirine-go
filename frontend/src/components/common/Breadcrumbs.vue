<template>
  <nav class="flex" aria-label="Breadcrumb">
    <ol class="inline-flex items-center space-x-1 md:space-x-3">
      <li class="inline-flex items-center">
        <router-link to="/dashboard" class="inline-flex items-center text-sm font-medium text-gray-500 hover:text-indigo-600 transition-colors">
          <Home class="w-4 h-4 mr-2" />
          Dashboard
        </router-link>
      </li>
      <li v-for="(crumb, index) in crumbs" :key="index">
        <div class="flex items-center">
          <ChevronRight class="w-4 h-4 text-gray-400" />
          <span 
            v-if="index === crumbs.length - 1"
            class="ml-1 text-sm font-medium text-gray-900 md:ml-2"
          >
            {{ crumb.name }}
          </span>
          <router-link 
            v-else
            :to="crumb.path" 
            class="ml-1 text-sm font-medium text-gray-500 hover:text-indigo-600 md:ml-2 transition-colors"
          >
            {{ crumb.name }}
          </router-link>
        </div>
      </li>
    </ol>
  </nav>
</template>

<script setup>
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { Home, ChevronRight } from 'lucide-vue-next'

const route = useRoute()

const crumbs = computed(() => {
  const pathArray = route.path.split('/').filter(p => p && p !== 'dashboard')
  
  return pathArray.map((path, index) => {
    // Construct the path for the link
    const to = `/dashboard/${pathArray.slice(0, index + 1).join('/')}`
    
    // Format the name (capitalize, replace hyphens with spaces)
    const name = path
      .split('-')
      .map(word => word.charAt(0).toUpperCase() + word.slice(1))
      .join(' ')

    return {
      name,
      path: to
    }
  })
})
</script>
