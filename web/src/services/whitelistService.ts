import { apiClient } from './api'

export type WhitelistType = 'ip' | 'region'

export interface WhitelistEntry {
  id: number
  type: WhitelistType
  value: string
  description?: string
  enabled: boolean
  created_at: string
  updated_at: string
}

export interface WhitelistPolicy {
  id: number
  enabled: boolean
  created_at: string
  updated_at: string
}

export interface CreateWhitelistEntryRequest {
  type: WhitelistType
  value: string
  description?: string
  enabled?: boolean
}

export interface UpdateWhitelistEntryRequest {
  type?: WhitelistType
  value?: string
  description?: string
  enabled?: boolean
}

export const whitelistService = {
  async getPolicy(): Promise<WhitelistPolicy> {
    const response = await apiClient.get('/admin/whitelist/policy')
    return response.data
  },

  async updatePolicy(enabled: boolean): Promise<void> {
    await apiClient.put('/admin/whitelist/policy', { enabled })
  },

  async listEntries(type?: WhitelistType): Promise<WhitelistEntry[]> {
    const params = type ? { type } : {}
    const response = await apiClient.get('/admin/whitelist/entries', { params })
    return response.data
  },

  async createEntry(data: CreateWhitelistEntryRequest): Promise<WhitelistEntry> {
    const response = await apiClient.post('/admin/whitelist/entries', data)
    return response.data
  },

  async updateEntry(id: number, data: UpdateWhitelistEntryRequest): Promise<WhitelistEntry> {
    const response = await apiClient.put(`/admin/whitelist/entries/${id}`, data)
    return response.data
  },

  async deleteEntry(id: number): Promise<void> {
    await apiClient.delete(`/admin/whitelist/entries/${id}`)
  },
}
