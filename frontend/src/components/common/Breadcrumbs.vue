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
            v-else-if="crumb.isClickable && crumb.path"
            :to="crumb.path" 
            class="ml-1 text-sm font-medium text-gray-500 hover:text-indigo-600 md:ml-2 transition-colors"
          >
            {{ crumb.name }}
          </router-link>
          <span
            v-else
            class="ml-1 text-sm font-medium text-gray-400 md:ml-2 cursor-not-allowed"
          >
            {{ crumb.name }}
          </span>
        </div>
      </li>
    </ol>
  </nav>
</template>

<script setup>
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Home, ChevronRight } from 'lucide-vue-next'

const route = useRoute()
const router = useRouter()

/**
 * Generate breadcrumb items dari current route path
 * dengan membangun path secara incremental dan validate
 * bahwa route tersebut exist di router dengan exact match
 */
const crumbs = computed(() => {
  // Split path dan filter empty strings
  const pathSegments = route.path.split('/').filter(p => p)
  
  // Filter out 'dashboard' karena sudah ada di home link
  const filteredSegments = pathSegments.filter(p => p !== 'dashboard')
  
  // Get all defined routes untuk validation
  const allRoutes = router.getRoutes()
  
  return filteredSegments.map((segment, index) => {
    // Build full path dari awal sampai segment saat ini
    const fullPath = '/' + pathSegments.slice(0, pathSegments.indexOf(segment) + 1).join('/')
    
    // Format name: capitalize dan replace hyphens dengan spaces
    const name = segment
      .split('-')
      .map(word => word.charAt(0).toUpperCase() + word.slice(1))
      .join(' ')

    // Check apakah exact path ini exist sebagai route (bukan dynamic route)
    // Hanya buat clickable jika route benar-benar exist
    const routeExists = allRoutes.some(r => {
      // Match exact path atau path tanpa dynamic segments
      return r.path === fullPath && r.name !== 'NotFound'
    })

    return {
      name,
      path: routeExists ? fullPath : null,
      isClickable: routeExists && index < filteredSegments.length - 1
    }
  })
})
</script>
