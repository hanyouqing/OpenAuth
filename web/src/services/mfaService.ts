import { apiClient } from './api'

export interface MFADevice {
  id: number
  type: 'totp' | 'sms' | 'email' | 'webauthn'
  name: string
  verified: boolean
  created_at: string
}

export interface GenerateTOTPResponse {
  secret: string
  qr_code: string
}

export const mfaService = {
  async generateTOTP(name: string): Promise<GenerateTOTPResponse> {
    const response = await apiClient.post<GenerateTOTPResponse>(
      '/mfa/totp/generate',
      { name }
    )
    return response.data
  },

  async verifyTOTP(token: string): Promise<void> {
    await apiClient.post('/mfa/totp/verify', { token })
  },

  async generateSMSCode(phoneNumber: string): Promise<void> {
    await apiClient.post('/mfa/sms/generate', { phone_number: phoneNumber })
  },

  async verifySMSCode(code: string): Promise<void> {
    await apiClient.post('/mfa/sms/verify', { code })
  },

  async listDevices(): Promise<{ data: MFADevice[] }> {
    const response = await apiClient.get<{ data: MFADevice[] }>('/mfa/devices')
    return response.data
  },

  async deleteDevice(deviceId: number): Promise<void> {
    await apiClient.delete(`/mfa/devices/${deviceId}`)
  },
}
