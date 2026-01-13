import axios, { AxiosInstance, AxiosError } from 'axios'
import { storage } from '@/utils/storage'

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || '/api/v1'

class ApiClient {
  private client: AxiosInstance

  constructor() {
    this.client = axios.create({
      baseURL: API_BASE_URL,
      headers: {
        'Content-Type': 'application/json',
      },
    })

    // Request interceptor
    this.client.interceptors.request.use(
      (config) => {
        const token = storage.getToken()
        if (token) {
          config.headers.Authorization = `Bearer ${token}`
        }
        return config
      },
      (error) => {
        return Promise.reject(error)
      }
    )

    // Response interceptor
    this.client.interceptors.response.use(
      (response) => response,
      async (error: AxiosError) => {
        const originalRequest = error.config as any

        // Avoid infinite loop on refresh token endpoint
        if (error.response?.status === 401 && !originalRequest._retry) {
          // Skip refresh for login endpoint
          if (originalRequest.url?.includes('/auth/login')) {
            return Promise.reject(error)
          }

          originalRequest._retry = true
          const refreshToken = storage.getRefreshToken()
          
          if (refreshToken) {
            try {
              // Use apiClient to ensure proxy is used
              const response = await this.client.post<{ code: number; message: string; data: { access_token: string } }>('/auth/refresh', {
                refresh_token: refreshToken,
              })
              // Backend returns {code, message, data} format
              const refreshData = response.data.data || response.data
              const access_token = refreshData.access_token || (response.data as any).access_token
              storage.setToken(access_token)
              
              // Retry original request
              if (originalRequest) {
                originalRequest.headers.Authorization = `Bearer ${access_token}`
                return this.client.request(originalRequest)
              }
            } catch (refreshError) {
              // Refresh failed, clear tokens and redirect to login
              storage.clear()
              // Only redirect if not already on login page
              if (window.location.pathname !== '/login') {
                window.location.href = '/login'
              }
              return Promise.reject(refreshError)
            }
          } else {
            // No refresh token, clear and redirect
            storage.clear()
            if (window.location.pathname !== '/login') {
              window.location.href = '/login'
            }
          }
        }
        return Promise.reject(error)
      }
    )
  }

  get instance() {
    return this.client
  }
}

export const apiClient = new ApiClient().instance
