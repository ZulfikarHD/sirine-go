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
    path: '/forgot-password',
    name: 'ForgotPassword',
    component: () => import('../views/auth/ForgotPassword.vue'),
    meta: { requiresAuth: false, guestOnly: true },
  },
  {
    path: '/reset-password',
    name: 'ResetPassword',
    component: () => import('../views/auth/ResetPassword.vue'),
    meta: { requiresAuth: false, guestOnly: true },
  },
  {
    path: '/force-change-password',
    name: 'ForceChangePassword',
    component: () => import('../views/auth/ForceChangePassword.vue'),
    meta: { requiresAuth: true, skipPasswordCheck: true },
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
    path: '/profile/edit',
    name: 'EditProfile',
    component: () => import('../views/profile/EditProfile.vue'),
    meta: { requiresAuth: true },
  },
  {
    path: '/profile/change-password',
    name: 'ChangePassword',
    component: () => import('../views/profile/ChangePassword.vue'),
    meta: { requiresAuth: true },
  },
  {
    path: '/profile/achievements',
    name: 'Achievements',
    component: () => import('../views/profile/Achievements.vue'),
    meta: { requiresAuth: true },
  },
  {
    path: '/admin/users',
    name: 'UserManagement',
    component: () => import('../views/admin/users/UserList.vue'),
    meta: { 
      requiresAuth: true,
      roles: ['ADMIN', 'MANAGER'],
    },
  },
  {
    path: '/admin/audit',
    name: 'ActivityLogs',
    component: () => import('../views/admin/audit/ActivityLogs.vue'),
    meta: { 
      requiresAuth: true,
      roles: ['ADMIN', 'MANAGER'],
    },
  },
  {
    path: '/notifications',
    name: 'NotificationCenter',
    component: () => import('../views/notifications/NotificationCenter.vue'),
    meta: { 
      requiresAuth: true,
    },
  },
  {
    path: '/khazwal/material-prep',
    name: 'MaterialPrepQueue',
    component: () => import('../views/khazwal/MaterialPrepQueuePage.vue'),
    meta: { 
      requiresAuth: true,
      roles: ['STAFF_KHAZWAL', 'ADMIN', 'MANAGER'],
      title: 'Persiapan Material'
    }
  },
  {
    path: '/khazwal/material-prep/:id',
    name: 'MaterialPrepDetail',
    component: () => import('../views/khazwal/MaterialPrepDetailPage.vue'),
    meta: { 
      requiresAuth: true,
      roles: ['STAFF_KHAZWAL', 'ADMIN', 'MANAGER'],
      title: 'Detail Persiapan'
    }
  },
  {
    path: '/khazwal/material-prep/:id/process',
    name: 'MaterialPrepProcess',
    component: () => import('../views/khazwal/MaterialPrepProcessPage.vue'),
    meta: { 
      requiresAuth: true,
      roles: ['STAFF_KHAZWAL', 'ADMIN', 'MANAGER'],
      title: 'Proses Persiapan'
    }
  },
  {
    path: '/khazwal/material-prep/history',
    name: 'MaterialPrepHistory',
    component: () => import('../views/khazwal/MaterialPrepHistoryPage.vue'),
    meta: { 
      requiresAuth: true,
      roles: ['STAFF_KHAZWAL', 'SUPERVISOR_KHAZWAL', 'ADMIN', 'MANAGER'],
      title: 'Riwayat Persiapan'
    }
  },
  {
    path: '/khazwal/monitoring',
    name: 'KhazwalMonitoring',
    component: () => import('../views/khazwal/SupervisorMonitoringPage.vue'),
    meta: { 
      requiresAuth: true,
      roles: ['SUPERVISOR_KHAZWAL', 'ADMIN', 'MANAGER'],
      title: 'Monitoring Khazwal'
    }
  },
  {
    path: '/cetak/queue',
    name: 'CetakQueue',
    component: () => import('../views/cetak/CetakQueuePage.vue'),
    meta: { 
      requiresAuth: true,
      roles: ['OPERATOR_CETAK', 'SUPERVISOR_CETAK', 'ADMIN', 'MANAGER'],
      title: 'Antrian Cetak'
    }
  },
  {
    path: '/khazwal/counting',
    name: 'counting-queue',
    component: () => import('../views/khazwal/counting/CountingQueuePage.vue'),
    meta: { 
      requiresAuth: true,
      roles: ['STAFF_KHAZWAL', 'ADMIN', 'MANAGER'],
      title: 'Penghitungan'
    }
  },
  {
    path: '/khazwal/counting/:poId',
    name: 'counting-work',
    component: () => import('../views/khazwal/counting/CountingWorkPage.vue'),
    meta: { 
      requiresAuth: true,
      roles: ['STAFF_KHAZWAL', 'ADMIN', 'MANAGER'],
      title: 'Proses Penghitungan'
    }
  },
  {
    path: '/khazwal/cutting',
    name: 'cutting-queue',
    component: () => import('../views/khazwal/cutting/CuttingQueuePage.vue'),
    meta: { 
      requiresAuth: true,
      roles: ['STAFF_KHAZWAL', 'ADMIN', 'MANAGER'],
      title: 'Pemotongan'
    }
  },
  {
    path: '/khazwal/cutting/start/:poId',
    name: 'cutting-start',
    component: () => import('../views/khazwal/cutting/CuttingStartPage.vue'),
    meta: { 
      requiresAuth: true,
      roles: ['STAFF_KHAZWAL', 'ADMIN', 'MANAGER'],
      title: 'Mulai Pemotongan'
    }
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
  const skipPasswordCheck = to.matched.some(record => record.meta.skipPasswordCheck)

  // Guest only pages (seperti login) - redirect jika sudah login
  if (guestOnly && authStore.isAuthenticated) {
    // Check jika user harus ganti password
    if (authStore.user?.must_change_password) {
      next('/force-change-password')
      return
    }
    
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

  // Check force change password (kecuali untuk route yang skip check)
  if (requiresAuth && !skipPasswordCheck && authStore.user?.must_change_password) {
    if (to.path !== '/force-change-password') {
      next('/force-change-password')
      return
    }
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
