import { apiClient } from './api'

export interface User {
  id: number
  username: string
  email: string
  first_name?: string
  last_name?: string
  avatar?: string
  status: string
  roles?: string[]
  is_admin?: boolean
  created_at: string
  updated_at: string
}

export interface CreateUserRequest {
  username: string
  email: string
  password: string
  first_name?: string
  last_name?: string
}

export interface UpdateUserRequest {
  email?: string
  first_name?: string
  last_name?: string
  status?: string
  avatar?: string
}

export interface ChangePasswordRequest {
  current_password: string
  new_password: string
}

export interface ListUsersResponse {
  data: User[]
  total: number
  page: number
  size: number
}

export const userService = {
  async list(page = 1, pageSize = 20): Promise<ListUsersResponse> {
    const response = await apiClient.get<ListUsersResponse>('/users', {
      params: { page, page_size: pageSize },
    })
    return response.data
  },

  async get(id: number): Promise<User> {
    const response = await apiClient.get<User>(`/users/${id}`)
    return response.data
  },

  async create(data: CreateUserRequest): Promise<User> {
    const response = await apiClient.post<User>('/users', data)
    return response.data
  },

  async update(id: number, data: UpdateUserRequest): Promise<User> {
    const response = await apiClient.put<User>(`/users/${id}`, data)
    return response.data
  },

  async delete(id: number): Promise<void> {
    await apiClient.delete(`/users/${id}`)
  },

  async getMe(): Promise<User> {
    const response = await apiClient.get<User>('/users/me')
    return response.data
  },

  async updateMe(data: UpdateUserRequest): Promise<User> {
    const response = await apiClient.put<User>('/users/me', data)
    return response.data
  },

  async changePassword(data: ChangePasswordRequest): Promise<void> {
    await apiClient.put('/users/me/password', data)
  },

  async uploadAvatar(avatar: string): Promise<User> {
    const response = await apiClient.put<User>('/users/me/avatar', { avatar })
    return response.data
  },
}
