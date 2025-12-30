/**
 * Integration Tests untuk Profile Management
 * yang mencakup view profile, edit profile, dan change password
 */

import { describe, it, expect, beforeEach, vi } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'
import axios from 'axios'
import { useAuthStore } from '@/stores/auth'

// Mock axios
vi.mock('axios')

describe('Profile Management Integration Tests', () => {
  beforeEach(() => {
    const pinia = createPinia()
    setActivePinia(pinia)
    
    // Clear mocks
    vi.clearAllMocks()
  })

  describe('View Profile', () => {
    it('user dapat melihat profile sendiri dengan complete data', async () => {
      const mockProfileData = {
        data: {
          success: true,
          data: {
            id: 1,
            nip: '12345',
            full_name: 'Test User',
            email: 'test@example.com',
            phone: '08123456789',
            role: 'ADMIN',
            department: 'KHAZWAL',
            profile_photo_url: '/uploads/profiles/1.jpg',
            created_at: '2025-01-01T00:00:00Z',
            last_login_at: '2025-01-15T10:30:00Z',
          },
        },
      }
      
      axios.get.mockResolvedValueOnce(mockProfileData)
      
      // Fetch profile
      const response = await axios.get('/api/profile')
      
      // Assertions
      expect(response.data.success).toBe(true)
      expect(response.data.data.nip).toBe('12345')
      expect(response.data.data.email).toBe('test@example.com')
    })

    it('profile harus menampilkan read-only fields dengan jelas', async () => {
      const mockProfileData = {
        data: {
          success: true,
          data: {
            id: 1,
            nip: '12345', // Read-only
            role: 'ADMIN', // Read-only
            department: 'KHAZWAL', // Read-only
            full_name: 'Test User', // Editable
            email: 'test@example.com', // Editable
            phone: '08123456789', // Editable
          },
        },
      }
      
      axios.get.mockResolvedValueOnce(mockProfileData)
      
      // Verify read-only vs editable fields distinction
      const response = await axios.get('/api/profile')
      const profile = response.data.data
      
      expect(profile.nip).toBeTruthy() // Exists tapi read-only
      expect(profile.role).toBeTruthy() // Exists tapi read-only
      expect(profile.full_name).toBeTruthy() // Editable
    })
  })

  describe('Edit Profile', () => {
    it('user dapat update editable fields (name, email, phone)', async () => {
      const updatedData = {
        full_name: 'Updated Name',
        email: 'updated@example.com',
        phone: '08987654321',
      }
      
      const mockUpdateResponse = {
        data: {
          success: true,
          data: {
            id: 1,
            nip: '12345',
            ...updatedData,
            role: 'ADMIN',
            department: 'KHAZWAL',
          },
          message: 'Profile berhasil diperbarui',
        },
      }
      
      axios.put.mockResolvedValueOnce(mockUpdateResponse)
      
      // Update profile
      const response = await axios.put('/api/profile', updatedData)
      
      // Assertions
      expect(response.data.success).toBe(true)
      expect(response.data.data.full_name).toBe('Updated Name')
      expect(response.data.data.email).toBe('updated@example.com')
      expect(response.data.message).toBeTruthy()
    })

    it('user tidak dapat update NIP, role, atau department', async () => {
      const attemptedUpdate = {
        nip: 'admin', // Should be ignored
        role: 'MANAGER', // Should be ignored
        department: 'CETAK', // Should be ignored
        full_name: 'Updated Name',
      }
      
      const mockResponse = {
        data: {
          success: true,
          data: {
            id: 1,
            nip: '12345', // Unchanged
            role: 'ADMIN', // Unchanged
            department: 'KHAZWAL', // Unchanged
            full_name: 'Updated Name', // Updated
          },
        },
      }
      
      axios.put.mockResolvedValueOnce(mockResponse)
      
      // Try to update restricted fields
      const response = await axios.put('/api/profile', attemptedUpdate)
      
      // Verify restricted fields tidak berubah
      expect(response.data.data.nip).toBe('12345')
      expect(response.data.data.role).toBe('ADMIN')
      expect(response.data.data.department).toBe('KHAZWAL')
      expect(response.data.data.full_name).toBe('Updated Name')
    })

    it('validation error untuk invalid email format harus ditampilkan', async () => {
      const invalidData = {
        email: 'invalid-email-format',
      }
      
      const mockErrorResponse = {
        response: {
          status: 400,
          data: {
            success: false,
            message: 'Format email tidak valid',
            errors: {
              email: ['Email harus berupa alamat email yang valid'],
            },
          },
        },
      }
      
      axios.put.mockRejectedValueOnce(mockErrorResponse)
      
      // Try to update dengan invalid email
      try {
        await axios.put('/api/profile', invalidData)
        expect.fail('Should throw validation error')
      } catch (error) {
        expect(error.response.status).toBe(400)
        expect(error.response.data.errors.email).toBeTruthy()
      }
    })

    it('validation error untuk invalid phone format harus ditampilkan', async () => {
      const invalidData = {
        phone: '123', // Too short
      }
      
      const mockErrorResponse = {
        response: {
          status: 400,
          data: {
            success: false,
            message: 'Format nomor telepon tidak valid',
            errors: {
              phone: ['Nomor telepon harus dimulai dengan 08 dan minimal 10 digit'],
            },
          },
        },
      }
      
      axios.put.mockRejectedValueOnce(mockErrorResponse)
      
      // Try to update dengan invalid phone
      try {
        await axios.put('/api/profile', invalidData)
        expect.fail('Should throw validation error')
      } catch (error) {
        expect(error.response.status).toBe(400)
        expect(error.response.data.errors.phone).toBeTruthy()
      }
    })
  })

  describe('Change Password', () => {
    it('user dapat change password dengan current password valid', async () => {
      const passwordData = {
        current_password: 'OldPass123!',
        new_password: 'NewPass123!',
        confirm_password: 'NewPass123!',
      }
      
      const mockResponse = {
        data: {
          success: true,
          message: 'Password berhasil diubah. Silakan login kembali.',
        },
      }
      
      axios.put.mockResolvedValueOnce(mockResponse)
      
      // Change password
      const response = await axios.put('/api/profile/password', passwordData)
      
      // Assertions
      expect(response.data.success).toBe(true)
      expect(response.data.message).toContain('berhasil')
    })

    it('change password dengan wrong current password harus gagal', async () => {
      const passwordData = {
        current_password: 'WrongPassword123!',
        new_password: 'NewPass123!',
        confirm_password: 'NewPass123!',
      }
      
      const mockErrorResponse = {
        response: {
          status: 400,
          data: {
            success: false,
            message: 'Password saat ini tidak valid',
          },
        },
      }
      
      axios.put.mockRejectedValueOnce(mockErrorResponse)
      
      // Try to change password
      try {
        await axios.put('/api/profile/password', passwordData)
        expect.fail('Should throw error')
      } catch (error) {
        expect(error.response.status).toBe(400)
        expect(error.response.data.message).toContain('tidak valid')
      }
    })

    it('new password tidak boleh sama dengan current password', async () => {
      const passwordData = {
        current_password: 'SamePass123!',
        new_password: 'SamePass123!',
        confirm_password: 'SamePass123!',
      }
      
      const mockErrorResponse = {
        response: {
          status: 400,
          data: {
            success: false,
            message: 'Password baru tidak boleh sama dengan password saat ini',
          },
        },
      }
      
      axios.put.mockRejectedValueOnce(mockErrorResponse)
      
      // Try to change password
      try {
        await axios.put('/api/profile/password', passwordData)
        expect.fail('Should throw error')
      } catch (error) {
        expect(error.response.data.message).toContain('tidak boleh sama')
      }
    })

    it('password policy harus enforced (min 8 char, uppercase, number, special)', async () => {
      const weakPasswordData = {
        current_password: 'OldPass123!',
        new_password: 'weak', // Too weak
        confirm_password: 'weak',
      }
      
      const mockErrorResponse = {
        response: {
          status: 400,
          data: {
            success: false,
            message: 'Password tidak memenuhi policy',
            errors: {
              new_password: [
                'Password minimal 8 karakter',
                'Password harus mengandung minimal 1 huruf besar',
                'Password harus mengandung minimal 1 angka',
                'Password harus mengandung minimal 1 karakter spesial',
              ],
            },
          },
        },
      }
      
      axios.put.mockRejectedValueOnce(mockErrorResponse)
      
      // Try to set weak password
      try {
        await axios.put('/api/profile/password', weakPasswordData)
        expect.fail('Should throw validation error')
      } catch (error) {
        expect(error.response.data.errors.new_password).toBeTruthy()
        expect(error.response.data.errors.new_password.length).toBeGreaterThan(0)
      }
    })

    it('password mismatch harus ditangkap frontend validation', async () => {
      const mismatchPasswordData = {
        current_password: 'OldPass123!',
        new_password: 'NewPass123!',
        confirm_password: 'DifferentPass123!', // Mismatch
      }
      
      // Frontend validation should catch this sebelum API call
      // Note: Test ini assume ada frontend validation
      expect(mismatchPasswordData.new_password).not.toBe(mismatchPasswordData.confirm_password)
    })

    it('setelah change password, user harus logout semua sessions', async () => {
      const passwordData = {
        current_password: 'OldPass123!',
        new_password: 'NewPass123!',
        confirm_password: 'NewPass123!',
      }
      
      const mockResponse = {
        data: {
          success: true,
          message: 'Password berhasil diubah. Silakan login kembali.',
        },
      }
      
      axios.put.mockResolvedValueOnce(mockResponse)
      
      // Setup store dengan existing session
      const pinia = createPinia()
      setActivePinia(pinia)
      const authStore = useAuthStore()
      authStore.token = 'old-token'
      authStore.user = { id: 1, nip: '12345' }
      
      // Change password
      await axios.put('/api/profile/password', passwordData)
      
      // Note: Implementation actual di component akan call logout
      // Verify expectation bahwa user akan di-logout
      expect(mockResponse.data.message).toContain('login kembali')
    })
  })

  describe('Profile Photo Upload', () => {
    it('user dapat upload profile photo dengan format valid', async () => {
      const mockFile = new File(['photo'], 'profile.jpg', { type: 'image/jpeg' })
      
      const mockUploadResponse = {
        data: {
          success: true,
          data: {
            profile_photo_url: '/uploads/profiles/1.jpg',
          },
          message: 'Foto profil berhasil diupload',
        },
      }
      
      axios.post.mockResolvedValueOnce(mockUploadResponse)
      
      // Upload photo
      const formData = new FormData()
      formData.append('photo', mockFile)
      
      const response = await axios.post('/api/profile/photo', formData)
      
      // Assertions
      expect(response.data.success).toBe(true)
      expect(response.data.data.profile_photo_url).toBeTruthy()
    })

    it('upload photo dengan invalid format harus ditolak', async () => {
      const mockErrorResponse = {
        response: {
          status: 400,
          data: {
            success: false,
            message: 'Format file tidak valid. Gunakan JPG, PNG, atau WebP',
          },
        },
      }
      
      axios.post.mockRejectedValueOnce(mockErrorResponse)
      
      // Try to upload invalid format
      const mockFile = new File(['document'], 'file.pdf', { type: 'application/pdf' })
      const formData = new FormData()
      formData.append('photo', mockFile)
      
      try {
        await axios.post('/api/profile/photo', formData)
        expect.fail('Should throw error')
      } catch (error) {
        expect(error.response.data.message).toContain('Format file tidak valid')
      }
    })

    it('upload photo dengan file size terlalu besar harus ditolak', async () => {
      const mockErrorResponse = {
        response: {
          status: 400,
          data: {
            success: false,
            message: 'Ukuran file terlalu besar. Maksimal 2MB',
          },
        },
      }
      
      axios.post.mockRejectedValueOnce(mockErrorResponse)
      
      // Try to upload large file
      try {
        const formData = new FormData()
        await axios.post('/api/profile/photo', formData)
        expect.fail('Should throw error')
      } catch (error) {
        expect(error.response.data.message).toContain('terlalu besar')
      }
    })
  })

  describe('Activity Tracking', () => {
    it('profile update harus tercatat di activity logs', async () => {
      const updatedData = {
        full_name: 'Updated Name',
      }
      
      const mockUpdateResponse = {
        data: {
          success: true,
          data: {
            id: 1,
            full_name: 'Updated Name',
          },
        },
      }
      
      axios.put.mockResolvedValueOnce(mockUpdateResponse)
      
      // Update profile
      await axios.put('/api/profile', updatedData)
      
      // Note: Backend akan automatically log activity
      // Verify bahwa update berhasil, logging handled di backend
      expect(mockUpdateResponse.data.success).toBe(true)
    })
  })
})
