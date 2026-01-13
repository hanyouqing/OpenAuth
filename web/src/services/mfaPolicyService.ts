import { apiClient } from './api'

export interface MFAPolicy {
  id: number
  force_mfa: boolean
}

export interface UpdateMFAPolicyRequest {
  force_mfa: boolean
}

export const mfaPolicyService = {
  async getPolicy(): Promise<MFAPolicy> {
    const response = await apiClient.get('/mfa/policy')
    return response.data
  },

  async updatePolicy(forceMFA: boolean): Promise<MFAPolicy> {
    const response = await apiClient.put('/mfa/policy', { force_mfa: forceMFA })
    return response.data
  },
}
