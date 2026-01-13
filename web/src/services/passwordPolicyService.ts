import { apiClient } from './api'

export interface PasswordPolicy {
  id: number
  min_length: number
  require_uppercase: boolean
  require_lowercase: boolean
  require_number: boolean
  require_special_char: boolean
  max_age: number // days, 0 means no expiration
  expiration_warning_days: number
  history_count: number
  lockout_threshold: number
  lockout_duration: number
}

export interface UpdatePasswordPolicyRequest {
  min_length?: number
  require_uppercase?: boolean
  require_lowercase?: boolean
  require_number?: boolean
  require_special_char?: boolean
  max_age?: number // days, 0 means no expiration, max 365
  expiration_warning_days?: number
  history_count?: number
  lockout_threshold?: number
  lockout_duration?: number
}

export const passwordPolicyService = {
  async getPolicy(): Promise<PasswordPolicy> {
    const response = await apiClient.get('/password-policy')
    return response.data
  },

  async updatePolicy(policy: UpdatePasswordPolicyRequest): Promise<PasswordPolicy> {
    const response = await apiClient.put('/password-policy', policy)
    return response.data
  },

  async validatePassword(password: string): Promise<{ valid: boolean }> {
    const response = await apiClient.post('/password-policy/validate', { password })
    return response.data
  },
}
