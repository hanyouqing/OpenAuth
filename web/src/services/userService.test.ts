import { describe, it, expect, vi, beforeEach } from 'vitest'
import { userService } from './userService'
import { apiClient } from './api'

vi.mock('./api')

describe('userService', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  describe('list', () => {
    it('should get user list successfully', async () => {
      const mockResponse = {
        data: [
          { id: 1, username: 'user1', email: 'user1@example.com', status: 'active', created_at: '2025-01-01', updated_at: '2025-01-01' },
          { id: 2, username: 'user2', email: 'user2@example.com', status: 'active', created_at: '2025-01-01', updated_at: '2025-01-01' },
        ],
        total: 2,
        page: 1,
        size: 20,
      }

      vi.mocked(apiClient.get).mockResolvedValue({ data: mockResponse })

      const result = await userService.list(1, 20)

      expect(apiClient.get).toHaveBeenCalledWith('/users', { params: { page: 1, page_size: 20 } })
      expect(result).toEqual(mockResponse)
    })
  })

  describe('get', () => {
    it('should get user by id successfully', async () => {
      const mockResponse = {
        id: 1,
        username: 'testuser',
        email: 'test@example.com',
        status: 'active',
        created_at: '2025-01-01',
        updated_at: '2025-01-01',
      }

      vi.mocked(apiClient.get).mockResolvedValue({ data: mockResponse })

      const result = await userService.get(1)

      expect(apiClient.get).toHaveBeenCalledWith('/users/1')
      expect(result).toEqual(mockResponse)
    })
  })

  describe('create', () => {
    it('should create user successfully', async () => {
      const mockResponse = {
        id: 1,
        username: 'newuser',
        email: 'newuser@example.com',
        status: 'active',
        created_at: '2025-01-01',
        updated_at: '2025-01-01',
      }

      vi.mocked(apiClient.post).mockResolvedValue({ data: mockResponse })

      const result = await userService.create({
        username: 'newuser',
        email: 'newuser@example.com',
        password: 'password123',
      })

      expect(apiClient.post).toHaveBeenCalledWith('/users', {
        username: 'newuser',
        email: 'newuser@example.com',
        password: 'password123',
      })
      expect(result).toEqual(mockResponse)
    })
  })

  describe('update', () => {
    it('should update user successfully', async () => {
      const mockResponse = {
        id: 1,
        username: 'updateduser',
        email: 'updated@example.com',
        status: 'active',
        created_at: '2025-01-01',
        updated_at: '2025-01-01',
      }

      vi.mocked(apiClient.put).mockResolvedValue({ data: mockResponse })

      const result = await userService.update(1, { email: 'updated@example.com' })

      expect(apiClient.put).toHaveBeenCalledWith('/users/1', { email: 'updated@example.com' })
      expect(result).toEqual(mockResponse)
    })
  })

  describe('delete', () => {
    it('should delete user successfully', async () => {
      vi.mocked(apiClient.delete).mockResolvedValue({ data: {} })

      await userService.delete(1)

      expect(apiClient.delete).toHaveBeenCalledWith('/users/1')
    })
  })

  describe('getMe', () => {
    it('should get current user successfully', async () => {
      const mockResponse = {
        id: 1,
        username: 'currentuser',
        email: 'current@example.com',
        status: 'active',
        created_at: '2025-01-01',
        updated_at: '2025-01-01',
      }

      vi.mocked(apiClient.get).mockResolvedValue({ data: mockResponse })

      const result = await userService.getMe()

      expect(apiClient.get).toHaveBeenCalledWith('/users/me')
      expect(result).toEqual(mockResponse)
    })
  })

  describe('updateMe', () => {
    it('should update current user successfully', async () => {
      const mockResponse = {
        id: 1,
        username: 'updateduser',
        email: 'updated@example.com',
        status: 'active',
        created_at: '2025-01-01',
        updated_at: '2025-01-01',
      }

      vi.mocked(apiClient.put).mockResolvedValue({ data: mockResponse })

      const result = await userService.updateMe({ email: 'updated@example.com' })

      expect(apiClient.put).toHaveBeenCalledWith('/users/me', { email: 'updated@example.com' })
      expect(result).toEqual(mockResponse)
    })
  })

  describe('changePassword', () => {
    it('should change password successfully', async () => {
      vi.mocked(apiClient.put).mockResolvedValue({ data: {} })

      await userService.changePassword({
        current_password: 'oldpassword',
        new_password: 'newpassword123',
      })

      expect(apiClient.put).toHaveBeenCalledWith('/users/me/password', {
        current_password: 'oldpassword',
        new_password: 'newpassword123',
      })
    })
  })
})
