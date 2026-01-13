import { apiClient } from './api'

export interface Application {
  id: number
  name: string
  type: 'saml' | 'oidc' | 'oauth2'
  status: string
  description?: string
  icon?: string
  config: Record<string, any>
  created_at: string
  updated_at: string
}

export interface CreateApplicationRequest {
  name: string
  type: 'saml' | 'oidc' | 'oauth2'
  description?: string
  icon?: string
  config?: Record<string, any>
}

export interface UpdateApplicationRequest {
  name?: string
  description?: string
  icon?: string
  status?: string
  config?: Record<string, any>
}

export interface AssignApplicationRequest {
  user_id?: number
  group_id?: number
}

export interface ListApplicationsResponse {
  data: Application[]
  total: number
  page: number
  size: number
}

export const applicationService = {
  async list(page?: number, pageSize?: number): Promise<ListApplicationsResponse> {
    const params: any = {}
    if (page !== undefined && pageSize !== undefined) {
      params.page = page
      params.page_size = pageSize
    }
    const response = await apiClient.get<ListApplicationsResponse>(
      '/applications',
      Object.keys(params).length > 0 ? { params } : undefined
    )
    return response.data
  },

  async get(id: number): Promise<Application> {
    const response = await apiClient.get<Application>(`/applications/${id}`)
    return response.data
  },

  async create(data: CreateApplicationRequest): Promise<Application> {
    const response = await apiClient.post<Application>('/applications', data)
    return response.data
  },

  async update(
    id: number,
    data: UpdateApplicationRequest
  ): Promise<Application> {
    const response = await apiClient.put<Application>(
      `/applications/${id}`,
      data
    )
    return response.data
  },

  async delete(id: number): Promise<void> {
    await apiClient.delete(`/applications/${id}`)
  },

  async assign(id: number, data: AssignApplicationRequest): Promise<void> {
    await apiClient.post(`/applications/${id}/assign`, data)
  },
}
