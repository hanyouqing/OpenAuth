import { apiClient } from './api'

export interface CPUInfo {
  usage_percent: number
  count: number
}

export interface MemoryInfo {
  total: number
  used: number
  available: number
  usage_percent: number
}

export interface DiskInfo {
  total: number
  used: number
  free: number
  usage_percent: number
}

export interface NetworkInfo {
  bytes_sent: number
  bytes_recv: number
  packets_sent: number
  packets_recv: number
}

export interface SystemResources {
  cpu: CPUInfo
  memory: MemoryInfo
  disk: DiskInfo
  network: NetworkInfo
}

export interface HealthStatus {
  status: string
  service: string
  timestamp: string
  uptime?: string
  resources?: SystemResources
}

export interface VersionInfo {
  version: string
  build_time: string
  git_commit: string
  go_version: string
  service: string
}

export interface MetricsData {
  raw: string
}

export const systemService = {
  async getHealth(): Promise<HealthStatus> {
    const response = await apiClient.get('/health')
    return response.data
  },

  async getVersion(): Promise<VersionInfo> {
    const response = await apiClient.get('/version')
    return response.data
  },

  async getMetrics(): Promise<MetricsData> {
    // Use axios directly for text/plain response
    const response = await apiClient.get('/metrics', {
      headers: {
        Accept: 'text/plain',
      },
      responseType: 'text',
      transformResponse: [(data) => data],
    })
    return { raw: typeof response.data === 'string' ? response.data : String(response.data) }
  },
}
