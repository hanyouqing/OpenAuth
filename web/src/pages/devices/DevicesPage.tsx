import { useState, useEffect } from 'react'
import { Table, Button, Space, message, Tag, Popconfirm, Tooltip } from 'antd'
import {
  DeleteOutlined,
  SafetyOutlined,
  SafetyCertificateOutlined,
  DesktopOutlined,
  MobileOutlined,
  TabletOutlined,
} from '@ant-design/icons'
import type { ColumnsType } from 'antd/es/table'
import { deviceService, Device } from '@/services/deviceService'
import dayjs from 'dayjs'

function DevicesPage() {
  const [devices, setDevices] = useState<Device[]>([])
  const [loading, setLoading] = useState(false)

  const fetchDevices = async () => {
    setLoading(true)
    try {
      const data = await deviceService.list()
      setDevices(data)
    } catch (error: any) {
      message.error(error.response?.data?.error || 'Failed to load devices')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchDevices()
  }, [])

  const handleTrust = async (deviceId: string) => {
    try {
      await deviceService.trust(deviceId)
      message.success('Device trusted successfully')
      fetchDevices()
    } catch (error: any) {
      message.error(error.response?.data?.error || 'Failed to trust device')
    }
  }

  const handleUntrust = async (deviceId: string) => {
    try {
      await deviceService.untrust(deviceId)
      message.success('Device untrusted successfully')
      fetchDevices()
    } catch (error: any) {
      message.error(error.response?.data?.error || 'Failed to untrust device')
    }
  }

  const handleDelete = async (deviceId: string) => {
    try {
      await deviceService.delete(deviceId)
      message.success('Device deleted successfully')
      fetchDevices()
    } catch (error: any) {
      message.error(error.response?.data?.error || 'Failed to delete device')
    }
  }

  const getDeviceIcon = (type?: string) => {
    switch (type) {
      case 'mobile':
        return <MobileOutlined />
      case 'tablet':
        return <TabletOutlined />
      default:
        return <DesktopOutlined />
    }
  }

  const columns: ColumnsType<Device> = [
    {
      title: 'Device',
      key: 'device',
      render: (_, record) => (
        <Space>
          {getDeviceIcon(record.device_type)}
          <div>
            <div style={{ fontWeight: 500 }}>
              {record.device_name || 'Unknown Device'}
            </div>
            <div style={{ fontSize: '12px', color: '#666' }}>
              {record.device_id.substring(0, 16)}...
            </div>
          </div>
        </Space>
      ),
    },
    {
      title: 'OS / Browser',
      key: 'os_browser',
      render: (_, record) => (
        <div>
          <div>{record.os || 'Unknown OS'}</div>
          <div style={{ fontSize: '12px', color: '#666' }}>
            {record.browser || 'Unknown Browser'}
          </div>
        </div>
      ),
    },
    {
      title: 'IP Address',
      dataIndex: 'ip_address',
      key: 'ip_address',
    },
    {
      title: 'Status',
      key: 'status',
      render: (_, record) => (
        <Tag
          icon={record.trusted ? <SafetyCertificateOutlined /> : <SafetyOutlined />}
          color={record.trusted ? 'green' : 'default'}
        >
          {record.trusted ? 'Trusted' : 'Not Trusted'}
        </Tag>
      ),
    },
    {
      title: 'Login Count',
      dataIndex: 'login_count',
      key: 'login_count',
      align: 'center',
    },
    {
      title: 'Last Seen',
      key: 'last_seen',
      render: (_, record) =>
        record.last_seen_at
          ? dayjs(record.last_seen_at).format('YYYY-MM-DD HH:mm')
          : '-',
    },
    {
      title: 'Actions',
      key: 'actions',
      render: (_, record) => (
        <Space>
          {record.trusted ? (
            <Tooltip title="Untrust Device">
              <Button
                size="small"
                icon={<SafetyOutlined />}
                onClick={() => handleUntrust(record.device_id)}
              >
                Untrust
              </Button>
            </Tooltip>
          ) : (
            <Tooltip title="Trust Device">
              <Button
                size="small"
                type="primary"
                icon={<SafetyCertificateOutlined />}
                onClick={() => handleTrust(record.device_id)}
              >
                Trust
              </Button>
            </Tooltip>
          )}
          <Popconfirm
            title="Are you sure you want to delete this device?"
            onConfirm={() => handleDelete(record.device_id)}
            okText="Yes"
            cancelText="No"
          >
            <Button size="small" danger icon={<DeleteOutlined />}>
              Delete
            </Button>
          </Popconfirm>
        </Space>
      ),
    },
  ]

  return (
    <div>
      <div
        style={{
          display: 'flex',
          justifyContent: 'space-between',
          alignItems: 'center',
          marginBottom: 24,
        }}
      >
        <h1 style={{ fontSize: '28px', fontWeight: 600, color: '#1e293b', margin: 0 }}>
          Devices
        </h1>
      </div>
      <Table
        columns={columns}
        dataSource={devices}
        rowKey="id"
        loading={loading}
        pagination={{
          pageSize: 20,
          showSizeChanger: true,
          showTotal: (total) => `Total ${total} devices`,
        }}
      />
    </div>
  )
}

export default DevicesPage
