import { describe, it, expect, vi, beforeEach } from 'vitest'
import { authService } from './authService'
import { apiClient } from './api'
import { storage } from '@/utils/storage'

vi.mock('./api')
vi.mock('@/utils/storage')

describe('authService', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  describe('login', () => {
    it('should login successfully with valid credentials', async () => {
      const mockResponse = {
        access_token: 'token123',
        refresh_token: 'refresh123',
        user: {
          id: 1,
          username: 'testuser',
          email: 'test@example.com',
          status: 'active',
        },
      }

      vi.mocked(apiClient.post).mockResolvedValue({ data: mockResponse })

      const result = await authService.login({ username: 'testuser', password: 'password123' })

      expect(apiClient.post).toHaveBeenCalledWith('/auth/login', {
        username: 'testuser',
        password: 'password123',
      })
      expect(storage.setToken).toHaveBeenCalledWith('token123')
      expect(storage.setRefreshToken).toHaveBeenCalledWith('refresh123')
      expect(result).toEqual(mockResponse)
    })

    it('should handle login error', async () => {
      vi.mocked(apiClient.post).mockRejectedValue(new Error('Invalid credentials'))

      await expect(authService.login({ username: 'testuser', password: 'wrongpassword' })).rejects.toThrow('Invalid credentials')
    })
  })

  describe('refresh', () => {
    it('should refresh token successfully', async () => {
      const mockResponse = {
        access_token: 'newtoken123',
      }

      vi.mocked(apiClient.post).mockResolvedValue({ data: mockResponse })

      const result = await authService.refresh('refresh123')

      expect(apiClient.post).toHaveBeenCalledWith('/auth/refresh', {
        refresh_token: 'refresh123',
      })
      expect(storage.setToken).toHaveBeenCalledWith('newtoken123')
      expect(result).toEqual(mockResponse)
    })

    it('should handle refresh token error', async () => {
      vi.mocked(apiClient.post).mockRejectedValue(new Error('Invalid refresh token'))

      await expect(authService.refresh('invalid-token')).rejects.toThrow('Invalid refresh token')
    })
  })

  describe('logout', () => {
    it('should logout successfully', async () => {
      vi.mocked(apiClient.post).mockResolvedValue({ data: {} })

      await authService.logout('refresh123')

      expect(apiClient.post).toHaveBeenCalledWith('/auth/logout', { refresh_token: 'refresh123' })
      expect(storage.clear).toHaveBeenCalled()
    })
  })
})
