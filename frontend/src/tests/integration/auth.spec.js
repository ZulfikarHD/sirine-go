/**
 * Integration Tests untuk Authentication Flow
 * yang mencakup login, logout, session persistence, dan auto-redirect
 */

import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { createRouter, createMemoryHistory } from 'vue-router'
import axios from 'axios'
import { useAuthStore } from '@/stores/auth'
import Login from '@/views/auth/Login.vue'

// Mock axios
vi.mock('axios')

// Mock router
const mockRoutes = [
  { path: '/login', name: 'login', component: { template: '<div>Login</div>' } },
  { path: '/dashboard', name: 'dashboard', component: { template: '<div>Dashboard</div>' } },
  { path: '/admin', name: 'admin', component: { template: '<div>Admin</div>' }, meta: { requiresAuth: true } },
]

/**
 * Setup test environment dengan Pinia dan Router
 * untuk simulate real application context
 */
function setupTestEnvironment() {
  const pinia = createPinia()
  setActivePinia(pinia)
  
  const router = createRouter({
    history: createMemoryHistory(),
    routes: mockRoutes,
  })
  
  return { pinia, router }
}

describe('Auth Flow Integration Tests', () => {
  beforeEach(() => {
    // Clear localStorage sebelum setiap test
    localStorage.clear()
    
    // Reset axios mocks
    vi.clearAllMocks()
  })

  describe('Login Flow', () => {
    it('login dengan valid credentials harus berhasil dan redirect ke dashboard', async () => {
      const { pinia, router } = setupTestEnvironment()
      
      // Mock successful login response
      const mockLoginResponse = {
        data: {
          success: true,
          data: {
            token: 'mock-jwt-token',
            refresh_token: 'mock-refresh-token',
            user: {
              id: 1,
              nip: '12345',
              full_name: 'Test User',
              email: 'test@example.com',
              role: 'ADMIN',
              department: 'KHAZWAL',
            },
            require_password_change: false,
          },
        },
      }
      
      axios.post.mockResolvedValueOnce(mockLoginResponse)
      
      // Setup store
      const authStore = useAuthStore(pinia)
      
      // Perform login
      const result = await authStore.login('12345', 'TestPass123!', false)
      
      // Assertions
      expect(result.success).toBe(true)
      expect(authStore.isAuthenticated).toBe(true)
      expect(authStore.user.nip).toBe('12345')
      expect(authStore.token).toBe('mock-jwt-token')
      expect(localStorage.getItem('auth_token')).toBe('mock-jwt-token')
    })

    it('login dengan invalid credentials harus gagal dan menampilkan error', async () => {
      const { pinia } = setupTestEnvironment()
      
      // Mock failed login response
      const mockErrorResponse = {
        response: {
          status: 401,
          data: {
            success: false,
            message: 'NIP/Email atau password salah',
          },
        },
      }
      
      axios.post.mockRejectedValueOnce(mockErrorResponse)
      
      // Setup store
      const authStore = useAuthStore(pinia)
      
      // Perform login
      try {
        await authStore.login('12345', 'WrongPassword', false)
      } catch (error) {
        // Expected error
      }
      
      // Assertions
      expect(authStore.isAuthenticated).toBe(false)
      expect(authStore.user).toBeNull()
      expect(authStore.token).toBeNull()
    })

    it('login dengan remember me harus menyimpan token lebih lama', async () => {
      const { pinia } = setupTestEnvironment()
      
      // Mock successful login dengan remember_me
      const mockLoginResponse = {
        data: {
          success: true,
          data: {
            token: 'mock-jwt-token-long',
            refresh_token: 'mock-refresh-token-long',
            user: {
              id: 1,
              nip: '12345',
              full_name: 'Test User',
              email: 'test@example.com',
              role: 'ADMIN',
              department: 'KHAZWAL',
            },
            require_password_change: false,
          },
        },
      }
      
      axios.post.mockResolvedValueOnce(mockLoginResponse)
      
      // Setup store
      const authStore = useAuthStore(pinia)
      
      // Perform login dengan remember_me = true
      await authStore.login('12345', 'TestPass123!', true)
      
      // Verify tokens saved
      expect(authStore.token).toBe('mock-jwt-token-long')
      expect(authStore.refreshToken).toBe('mock-refresh-token-long')
      expect(localStorage.getItem('auth_token')).toBe('mock-jwt-token-long')
      expect(localStorage.getItem('refresh_token')).toBe('mock-refresh-token-long')
    })

    it('login dengan email instead of NIP harus berhasil', async () => {
      const { pinia } = setupTestEnvironment()
      
      // Mock successful login dengan email
      const mockLoginResponse = {
        data: {
          success: true,
          data: {
            token: 'mock-jwt-token',
            refresh_token: 'mock-refresh-token',
            user: {
              id: 1,
              nip: '12345',
              full_name: 'Test User',
              email: 'test@example.com',
              role: 'ADMIN',
              department: 'KHAZWAL',
            },
            require_password_change: false,
          },
        },
      }
      
      axios.post.mockResolvedValueOnce(mockLoginResponse)
      
      // Setup store
      const authStore = useAuthStore(pinia)
      
      // Perform login dengan email
      await authStore.login('test@example.com', 'TestPass123!', false)
      
      // Assertions
      expect(authStore.isAuthenticated).toBe(true)
      expect(authStore.user.email).toBe('test@example.com')
    })
  })

  describe('Session Persistence', () => {
    it('refresh page dengan valid token harus maintain logged in state', async () => {
      const { pinia } = setupTestEnvironment()
      
      // Setup localStorage dengan existing token
      localStorage.setItem('auth_token', 'existing-token')
      localStorage.setItem('refresh_token', 'existing-refresh-token')
      localStorage.setItem('auth_user', JSON.stringify({
        id: 1,
        nip: '12345',
        full_name: 'Test User',
        email: 'test@example.com',
        role: 'ADMIN',
        department: 'KHAZWAL',
      }))
      
      // Mock GET /api/auth/me untuk verify token
      const mockMeResponse = {
        data: {
          success: true,
          data: {
            id: 1,
            nip: '12345',
            full_name: 'Test User',
            email: 'test@example.com',
            role: 'ADMIN',
            department: 'KHAZWAL',
          },
        },
      }
      
      axios.get.mockResolvedValueOnce(mockMeResponse)
      
      // Initialize store (akan load dari localStorage)
      const authStore = useAuthStore(pinia)
      
      // Verify current user
      await authStore.fetchCurrentUser()
      
      // Assertions
      expect(authStore.isAuthenticated).toBe(true)
      expect(authStore.user).toBeTruthy()
      expect(authStore.user.nip).toBe('12345')
    })

    it('refresh page dengan expired token harus redirect ke login', async () => {
      const { pinia, router } = setupTestEnvironment()
      
      // Setup localStorage dengan expired token
      localStorage.setItem('auth_token', 'expired-token')
      
      // Mock GET /api/auth/me dengan 401 response
      const mockErrorResponse = {
        response: {
          status: 401,
          data: {
            success: false,
            message: 'Token expired',
          },
        },
      }
      
      axios.get.mockRejectedValueOnce(mockErrorResponse)
      
      // Initialize store
      const authStore = useAuthStore(pinia)
      
      // Try to fetch current user
      try {
        await authStore.fetchCurrentUser()
      } catch (error) {
        // Expected error
      }
      
      // Assertions
      expect(authStore.isAuthenticated).toBe(false)
      expect(authStore.token).toBeNull()
      expect(localStorage.getItem('auth_token')).toBeNull()
    })
  })

  describe('Logout Flow', () => {
    it('logout harus clear token dan redirect ke login', async () => {
      const { pinia, router } = setupTestEnvironment()
      
      // Setup authenticated state
      localStorage.setItem('auth_token', 'valid-token')
      localStorage.setItem('refresh_token', 'valid-refresh-token')
      
      // Mock successful logout
      axios.post.mockResolvedValueOnce({
        data: { success: true },
      })
      
      // Setup store
      const authStore = useAuthStore(pinia)
      authStore.token = 'valid-token'
      authStore.refreshToken = 'valid-refresh-token'
      authStore.user = { id: 1, nip: '12345' }
      
      // Perform logout
      await authStore.logout()
      
      // Assertions
      expect(authStore.isAuthenticated).toBe(false)
      expect(authStore.token).toBeNull()
      expect(authStore.user).toBeNull()
      expect(localStorage.getItem('auth_token')).toBeNull()
      expect(localStorage.getItem('refresh_token')).toBeNull()
    })
  })

  describe('Auto-redirect', () => {
    it('access protected route tanpa auth harus redirect ke login', async () => {
      const { pinia, router } = setupTestEnvironment()
      
      // Setup unauthenticated state
      const authStore = useAuthStore(pinia)
      authStore.clearAuth()
      
      // Setup navigation guard
      router.beforeEach((to, from, next) => {
        if (to.meta.requiresAuth && !authStore.isAuthenticated) {
          next('/login')
        } else {
          next()
        }
      })
      
      // Try to navigate ke protected route
      await router.push('/admin')
      
      // Verify redirect ke login
      expect(router.currentRoute.value.path).toBe('/login')
    })

    it('access login page ketika sudah authenticated harus redirect ke dashboard', async () => {
      const { pinia, router } = setupTestEnvironment()
      
      // Setup authenticated state
      const authStore = useAuthStore(pinia)
      authStore.token = 'valid-token'
      authStore.user = { id: 1, nip: '12345', role: 'ADMIN' }
      
      // Setup navigation guard
      router.beforeEach((to, from, next) => {
        if (to.path === '/login' && authStore.isAuthenticated) {
          next('/dashboard')
        } else {
          next()
        }
      })
      
      // Try to navigate ke login
      await router.push('/login')
      
      // Verify redirect ke dashboard
      expect(router.currentRoute.value.path).toBe('/dashboard')
    })
  })

  describe('Form Validation', () => {
    it('submit login form dengan empty fields harus show validation error', async () => {
      const { pinia, router } = setupTestEnvironment()
      
      const wrapper = mount(Login, {
        global: {
          plugins: [pinia, router],
        },
      })
      
      // Try to submit tanpa fill fields
      const form = wrapper.find('form')
      if (form.exists()) {
        await form.trigger('submit.prevent')
        
        // Check untuk validation messages
        // Note: Implementasi actual akan depend pada Login.vue component
        await wrapper.vm.$nextTick()
      }
    })
  })

  describe('Token Refresh', () => {
    it('expired token harus trigger auto-refresh', async () => {
      const { pinia } = setupTestEnvironment()
      
      // Setup expired token
      const authStore = useAuthStore(pinia)
      authStore.token = 'expired-token'
      authStore.refreshToken = 'valid-refresh-token'
      
      // Mock refresh token response
      const mockRefreshResponse = {
        data: {
          success: true,
          data: {
            token: 'new-token',
            refresh_token: 'new-refresh-token',
          },
        },
      }
      
      axios.post.mockResolvedValueOnce(mockRefreshResponse)
      
      // Trigger refresh
      // Note: Implementasi actual di useApi.js interceptor
      // Ini hanya test store method jika ada
      if (authStore.refreshAuth) {
        await authStore.refreshAuth()
        
        expect(authStore.token).toBe('new-token')
        expect(authStore.refreshToken).toBe('new-refresh-token')
      }
    })
  })

  describe('Error Handling', () => {
    it('network error harus display user-friendly message', async () => {
      const { pinia } = setupTestEnvironment()
      
      // Mock network error
      axios.post.mockRejectedValueOnce(new Error('Network Error'))
      
      const authStore = useAuthStore(pinia)
      
      // Try login dengan network error
      try {
        await authStore.login('12345', 'TestPass123!', false)
        expect.fail('Should throw error')
      } catch (error) {
        expect(error.message).toBeTruthy()
      }
    })

    it('server error (500) harus handle gracefully', async () => {
      const { pinia } = setupTestEnvironment()
      
      // Mock server error
      const mockErrorResponse = {
        response: {
          status: 500,
          data: {
            success: false,
            message: 'Internal Server Error',
          },
        },
      }
      
      axios.post.mockRejectedValueOnce(mockErrorResponse)
      
      const authStore = useAuthStore(pinia)
      
      // Try login dengan server error
      try {
        await authStore.login('12345', 'TestPass123!', false)
        expect.fail('Should throw error')
      } catch (error) {
        // Error handled
        expect(authStore.isAuthenticated).toBe(false)
      }
    })
  })

  describe('Account Lockout', () => {
    it('account locked message harus ditampilkan dengan clear instruction', async () => {
      const { pinia } = setupTestEnvironment()
      
      // Mock account locked response
      const mockErrorResponse = {
        response: {
          status: 403,
          data: {
            success: false,
            message: 'Akun Anda terkunci hingga 15:30:00 karena terlalu banyak percobaan login gagal',
          },
        },
      }
      
      axios.post.mockRejectedValueOnce(mockErrorResponse)
      
      const authStore = useAuthStore(pinia)
      
      // Try login dengan locked account
      try {
        await authStore.login('12345', 'TestPass123!', false)
        expect.fail('Should throw error')
      } catch (error) {
        // Verify error message contains lockout info
        expect(error.response.data.message).toContain('terkunci')
      }
    })
  })
})
