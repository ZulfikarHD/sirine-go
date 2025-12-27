import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '../stores/auth'

/**
 * Route definitions dengan lazy loading untuk code splitting
 */
const routes = [
  {
    path: '/',
    redirect: '/dashboard',
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('../views/auth/Login.vue'),
    meta: { requiresAuth: false, guestOnly: true },
  },
  {
    path: '/dashboard',
    name: 'Dashboard',
    redirect: (to) => {
      const authStore = useAuthStore()
      if (!authStore.user) return '/login'
      
      // Redirect based on role
      const role = authStore.user.role
      if (role === 'ADMIN' || role === 'MANAGER') {
        return '/dashboard/admin'
      }
      return '/dashboard/staff'
    },
    meta: { requiresAuth: true },
  },
  {
    path: '/dashboard/admin',
    name: 'AdminDashboard',
    component: () => import('../views/dashboards/AdminDashboard.vue'),
    meta: { 
      requiresAuth: true, 
      roles: ['ADMIN', 'MANAGER'],
    },
  },
  {
    path: '/dashboard/staff',
    name: 'StaffDashboard',
    component: () => import('../views/dashboards/StaffDashboard.vue'),
    meta: { 
      requiresAuth: true,
    },
  },
  {
    path: '/profile',
    name: 'Profile',
    component: () => import('../views/profile/Profile.vue'),
    meta: { requiresAuth: true },
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: () => import('../views/NotFound.vue'),
  },
]

/**
 * Router instance dengan history mode
 */
const router = createRouter({
  history: createWebHistory(),
  routes,
})

/**
 * Navigation guard untuk authentication dan authorization
 */
router.beforeEach((to, from, next) => {
  const authStore = useAuthStore()
  
  // Restore auth dari localStorage jika ada
  if (!authStore.user && authStore.token) {
    authStore.restoreAuth()
  }

  const requiresAuth = to.matched.some(record => record.meta.requiresAuth)
  const guestOnly = to.matched.some(record => record.meta.guestOnly)
  const requiredRoles = to.meta.roles

  // Guest only pages (seperti login) - redirect jika sudah login
  if (guestOnly && authStore.isAuthenticated) {
    const dashboardRoute = getDashboardRouteByRole(authStore.user.role)
    next(dashboardRoute)
    return
  }

  // Protected pages - require authentication
  if (requiresAuth && !authStore.isAuthenticated) {
    next({
      path: '/login',
      query: { redirect: to.fullPath },
    })
    return
  }

  // Role-based access control
  if (requiresAuth && requiredRoles && !authStore.hasRole(...requiredRoles)) {
    // User tidak memiliki role yang diperlukan
    next('/dashboard') // Redirect ke dashboard default mereka
    return
  }

  next()
})

/**
 * Get dashboard route berdasarkan user role
 */
const getDashboardRouteByRole = (role) => {
  if (role === 'ADMIN' || role === 'MANAGER') {
    return '/dashboard/admin'
  }
  return '/dashboard/staff'
}

/**
 * After navigation hook untuk analytics atau logging
 */
router.afterEach((to, from) => {
  // Update document title
  const baseTitle = 'Sirine Go'
  document.title = to.meta.title ? `${to.meta.title} - ${baseTitle}` : baseTitle

  // Scroll to top
  window.scrollTo(0, 0)
})

export default router
