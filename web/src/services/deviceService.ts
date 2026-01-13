import { apiClient } from './api'

export interface Device {
  id: number
  user_id: number
  device_id: string
  device_name?: string
  device_type?: string
  os?: string
  browser?: string
  ip_address?: string
  user_agent?: string
  trusted: boolean
  last_seen_at: string
  first_seen_at: string
  login_count: number
  failed_login_count: number
  created_at: string
  updated_at: string
}

export const deviceService = {
  async list(): Promise<Device[]> {
    const response = await apiClient.get<{ code: number; data: Device[] }>('/devices')
    if (response.data.code === 200) {
      return response.data.data
    }
    throw new Error('Failed to fetch devices')
  },

  async trust(deviceId: string): Promise<void> {
    const response = await apiClient.post<{ code: number; message: string }>(
      `/devices/${deviceId}/trust`
    )
    if (response.data.code !== 200) {
      throw new Error(response.data.message || 'Failed to trust device')
    }
  },

  async untrust(deviceId: string): Promise<void> {
    const response = await apiClient.post<{ code: number; message: string }>(
      `/devices/${deviceId}/untrust`
    )
    if (response.data.code !== 200) {
      throw new Error(response.data.message || 'Failed to untrust device')
    }
  },

  async delete(deviceId: string): Promise<void> {
    const response = await apiClient.delete<{ code: number; message: string }>(
      `/devices/${deviceId}`
    )
    if (response.data.code !== 200) {
      throw new Error(response.data.message || 'Failed to delete device')
    }
  },
}
