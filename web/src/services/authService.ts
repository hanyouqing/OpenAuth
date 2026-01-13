import { apiClient } from './api'
import { storage } from '@/utils/storage'

export interface LoginRequest {
  username: string
  password: string
}

export interface PasswordExpirationWarning {
  warning: boolean
  days_remaining: number
  message: string
}

export interface LoginResponse {
  user: {
    id: number
    username: string
    email: string
    status: string
    roles?: string[]
    is_admin?: boolean
  }
  access_token: string
  refresh_token: string
  password_expiration_warning?: PasswordExpirationWarning
}

export interface RefreshRequest {
  refresh_token: string
}

export interface RefreshResponse {
  access_token: string
}

export const authService = {
  async login(data: LoginRequest): Promise<LoginResponse> {
    const response = await apiClient.post<{ code: number; message: string; data: LoginResponse }>('/auth/login', data)
    // Backend returns {code, message, data} format
    // response.data is the axios response data, which is {code, message, data}
    const responseData = response.data as any
    const loginData = responseData.data || responseData
    
    if (!loginData) {
      throw new Error('Invalid response format')
    }
    
    if (!loginData.access_token) {
      console.error('Login response:', loginData)
      throw new Error('No access token in response')
    }
    
    storage.setToken(loginData.access_token)
    storage.setRefreshToken(loginData.refresh_token)
    // Store user info including roles
    if (loginData.user) {
      storage.setUserInfo(loginData.user)
    }
    
    return loginData
  },

  async logout(refreshToken: string): Promise<void> {
    await apiClient.post('/auth/logout', { refresh_token: refreshToken })
    storage.clear()
  },

  async refresh(refreshToken: string): Promise<RefreshResponse> {
    const response = await apiClient.post<{ code: number; message: string; data: RefreshResponse }>('/auth/refresh', {
      refresh_token: refreshToken,
    })
    // Backend returns {code, message, data} format
    const refreshData = response.data.data || response.data
    if (refreshData.access_token) {
      storage.setToken(refreshData.access_token)
    }
    return refreshData
  },
}
