import { apiClient } from './api'

export interface AutomationTrigger {
  type: 'event' | 'schedule' | 'webhook' | 'manual'
  event?: string
  schedule?: string
  conditions?: Record<string, any>
}

export interface AutomationAction {
  type: string
  config: Record<string, any>
  on_error?: 'continue' | 'stop' | 'retry'
}

export interface AutomationWorkflow {
  id: number
  name: string
  description?: string
  enabled: boolean
  trigger: string | AutomationTrigger
  actions: string | AutomationAction[]
  priority: number
  created_at: string
  updated_at: string
}

export interface AutomationExecution {
  id: number
  workflow_id: number
  status: 'pending' | 'running' | 'completed' | 'failed'
  input: string | Record<string, any>
  output: string | Record<string, any>
  error?: string
  started_at?: string
  completed_at?: string
  created_at: string
}

export interface CreateWorkflowRequest {
  name: string
  description?: string
  trigger: AutomationTrigger
  actions: AutomationAction[]
  priority?: number
}

export interface UpdateWorkflowRequest {
  name?: string
  description?: string
  trigger?: AutomationTrigger
  actions?: AutomationAction[]
  enabled?: boolean
  priority?: number
}

export const automationService = {
  async list(enabled?: boolean): Promise<AutomationWorkflow[]> {
    const params = enabled !== undefined ? { enabled: enabled.toString() } : {}
    const response = await apiClient.get<{ code: number; data: AutomationWorkflow[] }>(
      '/automation/workflows',
      { params }
    )
    if (response.data.code === 200) {
      return response.data.data.map((w) => ({
        ...w,
        trigger: typeof w.trigger === 'string' ? JSON.parse(w.trigger) : w.trigger,
        actions: typeof w.actions === 'string' ? JSON.parse(w.actions) : w.actions,
      }))
    }
    throw new Error('Failed to fetch workflows')
  },

  async get(id: number): Promise<AutomationWorkflow> {
    const response = await apiClient.get<{ code: number; data: AutomationWorkflow }>(
      `/automation/workflows/${id}`
    )
    if (response.data.code === 200) {
      const w = response.data.data
      return {
        ...w,
        trigger: typeof w.trigger === 'string' ? JSON.parse(w.trigger) : w.trigger,
        actions: typeof w.actions === 'string' ? JSON.parse(w.actions) : w.actions,
      }
    }
    throw new Error('Failed to fetch workflow')
  },

  async create(data: CreateWorkflowRequest): Promise<AutomationWorkflow> {
    const response = await apiClient.post<{ code: number; data: AutomationWorkflow }>(
      '/automation/workflows',
      data
    )
    if (response.data.code === 200) {
      const w = response.data.data
      return {
        ...w,
        trigger: typeof w.trigger === 'string' ? JSON.parse(w.trigger) : w.trigger,
        actions: typeof w.actions === 'string' ? JSON.parse(w.actions) : w.actions,
      }
    }
    throw new Error('Failed to create workflow')
  },

  async update(id: number, data: UpdateWorkflowRequest): Promise<void> {
    const response = await apiClient.put<{ code: number; message: string }>(
      `/automation/workflows/${id}`,
      data
    )
    if (response.data.code !== 200) {
      throw new Error(response.data.message || 'Failed to update workflow')
    }
  },

  async delete(id: number): Promise<void> {
    const response = await apiClient.delete<{ code: number; message: string }>(
      `/automation/workflows/${id}`
    )
    if (response.data.code !== 200) {
      throw new Error(response.data.message || 'Failed to delete workflow')
    }
  },

  async getExecutions(workflowId: number, limit?: number): Promise<AutomationExecution[]> {
    const params = limit ? { limit: limit.toString() } : {}
    const response = await apiClient.get<{ code: number; data: AutomationExecution[] }>(
      `/automation/workflows/${workflowId}/executions`,
      { params }
    )
    if (response.data.code === 200) {
      return response.data.data.map((e) => ({
        ...e,
        input: typeof e.input === 'string' ? JSON.parse(e.input) : e.input,
        output: typeof e.output === 'string' ? JSON.parse(e.output) : e.output,
      }))
    }
    throw new Error('Failed to fetch executions')
  },

  async getExecution(id: number): Promise<AutomationExecution> {
    const response = await apiClient.get<{ code: number; data: AutomationExecution }>(
      `/automation/executions/${id}`
    )
    if (response.data.code === 200) {
      const e = response.data.data
      return {
        ...e,
        input: typeof e.input === 'string' ? JSON.parse(e.input) : e.input,
        output: typeof e.output === 'string' ? JSON.parse(e.output) : e.output,
      }
    }
    throw new Error('Failed to fetch execution')
  },
}
