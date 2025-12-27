import { describe, it, expect, beforeEach, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useAuthStore } from './auth'

describe('Auth Store', () => {
  beforeEach(() => {
    // Create fresh pinia instance untuk setiap test
    setActivePinia(createPinia())
    // Clear localStorage mock
    localStorage.clear()
  })

  describe('Initial State', () => {
    it('harus memiliki user null di awal', () => {
      const store = useAuthStore()
      expect(store.user).toBeNull()
    })

    it('harus memiliki token null jika tidak ada di localStorage', () => {
      const store = useAuthStore()
      expect(store.token).toBeNull()
    })

    it('harus tidak authenticated di awal', () => {
      const store = useAuthStore()
      expect(store.isAuthenticated).toBe(false)
    })
  })

  describe('setAuth', () => {
    it('harus set auth data dengan benar', () => {
      const store = useAuthStore()
      const authData = {
        token: 'test-token',
        refresh_token: 'test-refresh-token',
        user: {
          id: 1,
          nip: '99999',
          full_name: 'Test User',
          role: 'ADMIN',
        },
      }

      store.setAuth(authData)

      expect(store.token).toBe('test-token')
      expect(store.refreshToken).toBe('test-refresh-token')
      expect(store.user).toEqual(authData.user)
      expect(store.isAuthenticated).toBe(true)
    })

    it('harus menyimpan ke localStorage', () => {
      const store = useAuthStore()
      const authData = {
        token: 'test-token',
        refresh_token: 'test-refresh-token',
        user: { id: 1, nip: '99999' },
      }

      store.setAuth(authData)

      expect(localStorage.setItem).toHaveBeenCalledWith('auth_token', 'test-token')
      expect(localStorage.setItem).toHaveBeenCalledWith('refresh_token', 'test-refresh-token')
      expect(localStorage.setItem).toHaveBeenCalledWith('user_data', JSON.stringify(authData.user))
    })
  })

  describe('clearAuth', () => {
    it('harus clear semua auth data', () => {
      const store = useAuthStore()
      // Set auth first
      store.setAuth({
        token: 'test-token',
        refresh_token: 'test-refresh-token',
        user: { id: 1 },
      })

      // Clear auth
      store.clearAuth()

      expect(store.token).toBeNull()
      expect(store.refreshToken).toBeNull()
      expect(store.user).toBeNull()
      expect(store.isAuthenticated).toBe(false)
    })

    it('harus remove dari localStorage', () => {
      const store = useAuthStore()
      store.clearAuth()

      expect(localStorage.removeItem).toHaveBeenCalledWith('auth_token')
      expect(localStorage.removeItem).toHaveBeenCalledWith('refresh_token')
      expect(localStorage.removeItem).toHaveBeenCalledWith('user_data')
    })
  })

  describe('restoreAuth', () => {
    it('harus restore auth dari localStorage', () => {
      const mockUser = { id: 1, nip: '99999', full_name: 'Test' }
      localStorage.getItem = vi.fn((key) => {
        if (key === 'auth_token') return 'stored-token'
        if (key === 'refresh_token') return 'stored-refresh'
        if (key === 'user_data') return JSON.stringify(mockUser)
        return null
      })

      const store = useAuthStore()
      store.restoreAuth()

      expect(store.token).toBe('stored-token')
      expect(store.refreshToken).toBe('stored-refresh')
      expect(store.user).toEqual(mockUser)
    })

    it('harus handle invalid JSON di localStorage', () => {
      localStorage.getItem = vi.fn((key) => {
        if (key === 'auth_token') return 'token'
        if (key === 'user_data') return 'invalid-json'
        return null
      })

      const store = useAuthStore()
      store.restoreAuth()

      // Should clear auth jika JSON invalid
      expect(store.user).toBeNull()
      expect(store.token).toBeNull()
    })
  })

  describe('hasRole', () => {
    it('harus return true jika user memiliki role', () => {
      const store = useAuthStore()
      store.setAuth({
        token: 'token',
        refresh_token: 'refresh',
        user: { id: 1, role: 'ADMIN' },
      })

      expect(store.hasRole('ADMIN')).toBe(true)
      expect(store.hasRole('ADMIN', 'MANAGER')).toBe(true)
    })

    it('harus return false jika user tidak memiliki role', () => {
      const store = useAuthStore()
      store.setAuth({
        token: 'token',
        refresh_token: 'refresh',
        user: { id: 1, role: 'STAFF_KHAZWAL' },
      })

      expect(store.hasRole('ADMIN')).toBe(false)
      expect(store.hasRole('ADMIN', 'MANAGER')).toBe(false)
    })

    it('harus return false jika user null', () => {
      const store = useAuthStore()
      expect(store.hasRole('ADMIN')).toBe(false)
    })
  })

  describe('isAdmin', () => {
    it('harus return true untuk ADMIN role', () => {
      const store = useAuthStore()
      store.setAuth({
        token: 'token',
        refresh_token: 'refresh',
        user: { id: 1, role: 'ADMIN' },
      })

      expect(store.isAdmin()).toBe(true)
    })

    it('harus return true untuk MANAGER role', () => {
      const store = useAuthStore()
      store.setAuth({
        token: 'token',
        refresh_token: 'refresh',
        user: { id: 1, role: 'MANAGER' },
      })

      expect(store.isAdmin()).toBe(true)
    })

    it('harus return false untuk non-admin roles', () => {
      const store = useAuthStore()
      store.setAuth({
        token: 'token',
        refresh_token: 'refresh',
        user: { id: 1, role: 'STAFF_KHAZWAL' },
      })

      expect(store.isAdmin()).toBe(false)
    })
  })

  describe('Computed Properties', () => {
    it('userRole harus return role dari user', () => {
      const store = useAuthStore()
      store.setAuth({
        token: 'token',
        refresh_token: 'refresh',
        user: { id: 1, role: 'ADMIN' },
      })

      expect(store.userRole).toBe('ADMIN')
    })

    it('userDepartment harus return department dari user', () => {
      const store = useAuthStore()
      store.setAuth({
        token: 'token',
        refresh_token: 'refresh',
        user: { id: 1, department: 'KHAZWAL' },
      })

      expect(store.userDepartment).toBe('KHAZWAL')
    })

    it('requirePasswordChange harus return must_change_password flag', () => {
      const store = useAuthStore()
      store.setAuth({
        token: 'token',
        refresh_token: 'refresh',
        user: { id: 1, must_change_password: true },
      })

      expect(store.requirePasswordChange).toBe(true)
    })
  })
})
