import { useState, useEffect, useMemo } from 'react'
import { Card, Tabs, Descriptions, Tag, Spin, Alert, Typography, Space, Input, Progress, Divider, Statistic, Row, Col } from 'antd'
import {
  CheckCircleOutlined,
  InfoCircleOutlined,
  BarChartOutlined,
  LineChartOutlined,
  FileTextOutlined,
  SafetyOutlined,
} from '@ant-design/icons'
import { useTranslation } from 'react-i18next'
import { systemService, HealthStatus, VersionInfo, MetricsData } from '@/services/systemService'
import SecuritySettingsPage from './SecuritySettingsPage'

const { Title, Text } = Typography
const { TextArea } = Input

function SystemInfoPage() {
  const { t } = useTranslation()
  const [health, setHealth] = useState<HealthStatus | null>(null)
  const [version, setVersion] = useState<VersionInfo | null>(null)
  const [metrics, setMetrics] = useState<MetricsData | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const formatBytes = (bytes: number): string => {
    if (bytes === 0) return '0 B'
    const k = 1024
    const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
    const i = Math.floor(Math.log(bytes) / Math.log(k))
    return Math.round((bytes / Math.pow(k, i)) * 100) / 100 + ' ' + sizes[i]
  }

  // Parse Prometheus metrics format
  const parseMetrics = (raw: string) => {
    const metrics: Record<string, number> = {}
    const lines = raw.split('\n')
    
    for (const line of lines) {
      // Skip comments and empty lines
      const trimmed = line.trim()
      if (trimmed === '' || trimmed.startsWith('#')) {
        continue
      }
      
      // Parse metric line: metric_name value or metric_name{labels} value
      // Handle both simple format and format with labels
      // For metrics with labels, we'll aggregate by metric name (sum if multiple)
      const simpleMatch = trimmed.match(/^([a-zA-Z_:][a-zA-Z0-9_:]*)\s+([0-9.eE+-]+)$/)
      if (simpleMatch) {
        const [, name, value] = simpleMatch
        const numValue = parseFloat(value)
        if (!isNaN(numValue)) {
          metrics[name] = (metrics[name] || 0) + numValue
        }
        continue
      }
      
      // Handle metrics with labels: metric_name{label1="value1",label2="value2"} value
      const labelMatch = trimmed.match(/^([a-zA-Z_:][a-zA-Z0-9_:]*)\{.*\}\s+([0-9.eE+-]+)$/)
      if (labelMatch) {
        const [, name, value] = labelMatch
        const numValue = parseFloat(value)
        if (!isNaN(numValue)) {
          metrics[name] = (metrics[name] || 0) + numValue
        }
      }
    }
    
    return metrics
  }

  const parsedMetrics = useMemo(() => {
    if (!metrics?.raw) return null
    return parseMetrics(metrics.raw)
  }, [metrics])

  // Render metrics visualization
  const renderMetricsChart = () => {
    if (!parsedMetrics) {
      return <Spin />
    }

    const systemMetrics = [
      { name: 'openauth_uptime_seconds', label: 'Uptime (seconds)' },
      { name: 'openauth_memory_alloc_bytes', label: 'Memory Allocated' },
      { name: 'openauth_goroutines', label: 'Goroutines' },
      { name: 'openauth_gc_runs_total', label: 'GC Runs' },
    ]

    const httpMetrics = [
      { name: 'openauth_http_requests_total', label: 'Total Requests' },
      { name: 'openauth_http_requests_in_flight', label: 'Requests In Flight' },
      { name: 'openauth_http_request_errors_total', label: 'Request Errors' },
    ]

    const businessMetrics = [
      { name: 'openauth_login_attempts_total', label: 'Login Attempts' },
      { name: 'openauth_login_success_total', label: 'Login Success' },
      { name: 'openauth_login_failure_total', label: 'Login Failures' },
      { name: 'openauth_user_creations_total', label: 'User Creations' },
      { name: 'openauth_session_creations_total', label: 'Session Creations' },
      { name: 'openauth_sessions_active', label: 'Active Sessions' },
    ]

    const performanceMetrics = [
      { name: 'openauth_db_connections_active', label: 'DB Active Connections' },
      { name: 'openauth_db_connections_idle', label: 'DB Idle Connections' },
      { name: 'openauth_db_queries_total', label: 'DB Total Queries' },
      { name: 'openauth_cache_hits_total', label: 'Cache Hits' },
      { name: 'openauth_cache_misses_total', label: 'Cache Misses' },
    ]

    const renderMetricGroup = (title: string, metricList: Array<{ name: string; label: string }>) => {
      const availableMetrics = metricList.filter(m => parsedMetrics[m.name] !== undefined)
      if (availableMetrics.length === 0) return null

      return (
        <div key={title} style={{ marginBottom: 24 }}>
          <Divider orientation="left">{title}</Divider>
          <Row gutter={[16, 16]}>
            {availableMetrics.map((metric) => {
              const value = parsedMetrics[metric.name]
              return (
                <Col xs={24} sm={12} md={8} lg={6} key={metric.name}>
                  <Card>
                    <Statistic
                      title={metric.label}
                      value={value}
                      valueStyle={{ fontSize: '20px' }}
                      formatter={(val) => {
                        if (metric.name.includes('bytes') || metric.name.includes('memory')) {
                          return formatBytes(Number(val))
                        }
                        return Number(val).toLocaleString()
                      }}
                    />
                  </Card>
                </Col>
              )
            })}
          </Row>
        </div>
      )
    }

    return (
      <Space direction="vertical" style={{ width: '100%' }} size="large">
        {renderMetricGroup(t('system.metrics.system'), systemMetrics)}
        {renderMetricGroup(t('system.metrics.http'), httpMetrics)}
        {renderMetricGroup(t('system.metrics.business'), businessMetrics)}
        {renderMetricGroup(t('system.metrics.performance'), performanceMetrics)}
      </Space>
    )
  }

  const fetchData = async () => {
    setLoading(true)
    setError(null)
    try {
      const [healthData, versionData, metricsData] = await Promise.all([
        systemService.getHealth(),
        systemService.getVersion(),
        systemService.getMetrics(),
      ])
      setHealth(healthData)
      setVersion(versionData)
      setMetrics(metricsData)
    } catch (err: any) {
      setError(err.response?.data?.error || err.message || t('system.error'))
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchData()
    const interval = setInterval(fetchData, 30000) // Refresh every 30 seconds
    return () => clearInterval(interval)
  }, [])

  const tabItems = [
    {
      key: 'health',
      label: (
        <Space>
          <CheckCircleOutlined />
          {t('system.health.title')}
        </Space>
      ),
      children: (
        <Card>
          {health ? (
            <Space direction="vertical" style={{ width: '100%' }} size="large">
              <Descriptions bordered column={1}>
                <Descriptions.Item label={t('system.health.status')}>
                  <Tag
                    color={health.status === 'ok' ? 'success' : 'error'}
                    icon={<CheckCircleOutlined />}
                  >
                    {health.status.toUpperCase()}
                  </Tag>
                </Descriptions.Item>
                <Descriptions.Item label={t('system.health.service')}>{health.service}</Descriptions.Item>
                <Descriptions.Item label={t('system.health.timestamp')}>{health.timestamp}</Descriptions.Item>
                {health.uptime && (
                  <Descriptions.Item label={t('system.health.uptime')}>{health.uptime}</Descriptions.Item>
                )}
              </Descriptions>

              {health.resources && (
                <>
                  <Divider orientation="left">{t('system.resources.title')}</Divider>
                  
                  <div>
                    <Text strong>{t('system.resources.cpu')}</Text>
                    <Progress
                      percent={Math.round(health.resources.cpu.usage_percent)}
                      status={health.resources.cpu.usage_percent > 80 ? 'exception' : 'active'}
                      format={(percent) => `${percent}%`}
                    />
                    <Text type="secondary" style={{ fontSize: '12px' }}>
                      {t('system.resources.cpuCores')}: {health.resources.cpu.count}
                    </Text>
                  </div>

                  <div>
                    <Text strong>{t('system.resources.memory')}</Text>
                    <Progress
                      percent={Math.round(health.resources.memory.usage_percent)}
                      status={health.resources.memory.usage_percent > 80 ? 'exception' : 'active'}
                      format={(percent) => `${percent}%`}
                    />
                    <Text type="secondary" style={{ fontSize: '12px' }}>
                      {t('system.resources.used')}: {formatBytes(health.resources.memory.used)} / {formatBytes(health.resources.memory.total)}
                    </Text>
                  </div>

                  <div>
                    <Text strong>{t('system.resources.disk')}</Text>
                    <Progress
                      percent={Math.round(health.resources.disk.usage_percent)}
                      status={health.resources.disk.usage_percent > 80 ? 'exception' : 'active'}
                      format={(percent) => `${percent}%`}
                    />
                    <Text type="secondary" style={{ fontSize: '12px' }}>
                      {t('system.resources.used')}: {formatBytes(health.resources.disk.used)} / {formatBytes(health.resources.disk.total)} ({t('system.resources.free')}: {formatBytes(health.resources.disk.free)})
                    </Text>
                  </div>

                  <div>
                    <Text strong>{t('system.resources.network')}</Text>
                    <Descriptions bordered column={2} size="small">
                      <Descriptions.Item label={t('system.resources.bytesSent')}>
                        {formatBytes(health.resources.network.bytes_sent)}
                      </Descriptions.Item>
                      <Descriptions.Item label={t('system.resources.bytesRecv')}>
                        {formatBytes(health.resources.network.bytes_recv)}
                      </Descriptions.Item>
                      <Descriptions.Item label={t('system.resources.packetsSent')}>
                        {health.resources.network.packets_sent.toLocaleString()}
                      </Descriptions.Item>
                      <Descriptions.Item label={t('system.resources.packetsRecv')}>
                        {health.resources.network.packets_recv.toLocaleString()}
                      </Descriptions.Item>
                    </Descriptions>
                  </div>
                </>
              )}
            </Space>
          ) : (
            <Spin />
          )}
        </Card>
      ),
    },
    {
      key: 'version',
      label: (
        <Space>
          <InfoCircleOutlined />
          {t('system.version.title')}
        </Space>
      ),
      children: (
        <Card>
          {version ? (
            <Descriptions bordered column={1}>
              <Descriptions.Item label={t('system.version.version')}>
                <Tag color="blue">{version.version}</Tag>
              </Descriptions.Item>
              <Descriptions.Item label={t('system.version.buildTime')}>{version.build_time}</Descriptions.Item>
              <Descriptions.Item label={t('system.version.gitCommit')}>
                <Text code>{version.git_commit}</Text>
              </Descriptions.Item>
              <Descriptions.Item label={t('system.version.goVersion')}>
                <Text code>{version.go_version}</Text>
              </Descriptions.Item>
              <Descriptions.Item label={t('system.version.service')}>{version.service}</Descriptions.Item>
            </Descriptions>
          ) : (
            <Spin />
          )}
        </Card>
      ),
    },
    {
      key: 'metrics',
      label: (
        <Space>
          <BarChartOutlined />
          {t('system.metrics.title')}
        </Space>
      ),
      children: (
        <Card>
          {metrics ? (
            <Tabs
              items={[
                {
                  key: 'chart',
                  label: (
                    <Space>
                      <LineChartOutlined />
                      {t('system.metrics.chart')}
                    </Space>
                  ),
                  children: (
                    <div style={{ padding: '16px 0' }}>
                      {renderMetricsChart()}
                    </div>
                  ),
                },
                {
                  key: 'raw',
                  label: (
                    <Space>
                      <FileTextOutlined />
                      {t('system.metrics.raw')}
                    </Space>
                  ),
                  children: (
                    <Space direction="vertical" style={{ width: '100%' }} size="large">
                      <Text type="secondary">
                        {t('system.metrics.description')}
                      </Text>
                      <TextArea
                        value={metrics.raw}
                        rows={20}
                        readOnly
                        style={{
                          fontFamily: 'monospace',
                          fontSize: '12px',
                          backgroundColor: '#f5f5f5',
                        }}
                      />
                    </Space>
                  ),
                },
              ]}
            />
          ) : (
            <Spin />
          )}
        </Card>
      ),
    },
    {
      key: 'security',
      label: (
        <Space>
          <SafetyOutlined />
          {t('system.security.title')}
        </Space>
      ),
      children: <SecuritySettingsPage />,
    },
  ]

  return (
    <div>
      <div
        style={{
          marginBottom: 24,
          display: 'flex',
          justifyContent: 'space-between',
          alignItems: 'center',
        }}
      >
        <Title level={2} style={{ margin: 0 }}>
          {t('system.title')}
        </Title>
      </div>

      {error && (
        <Alert
          message={t('system.error')}
          description={error}
          type="error"
          showIcon
          closable
          onClose={() => setError(null)}
          style={{ marginBottom: 16 }}
        />
      )}

      {loading && !health && !version && !metrics ? (
        <Card>
          <Spin size="large" style={{ display: 'block', textAlign: 'center', padding: 40 }} />
        </Card>
      ) : (
        <Tabs items={tabItems} />
      )}
    </div>
  )
}

export default SystemInfoPage
