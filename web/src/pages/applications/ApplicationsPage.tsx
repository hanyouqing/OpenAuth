import { useState, useEffect } from 'react'
import { Table, Button, Space, message, Tag, Modal, Form, Input, Select } from 'antd'
import { PlusOutlined } from '@ant-design/icons'
import { useTranslation } from 'react-i18next'
import type { ColumnsType } from 'antd/es/table'
import { storage } from '@/utils/storage'
import {
  applicationService,
  Application,
  CreateApplicationRequest,
} from '@/services/applicationService'

function ApplicationsPage() {
  const { t } = useTranslation()
  const [applications, setApplications] = useState<Application[]>([])
  const [loading, setLoading] = useState(false)
  const [isAdmin, setIsAdmin] = useState(false)
  const [createModalVisible, setCreateModalVisible] = useState(false)
  const [form] = Form.useForm()
  const [pagination, setPagination] = useState({
    current: 1,
    pageSize: 20,
    total: 0,
  })

  useEffect(() => {
    // Check if user is admin
    const userInfo = storage.getUserInfo()
    const admin = userInfo?.is_admin === true
    setIsAdmin(admin)
    fetchApplications()
  }, [])

  const fetchApplications = async (page = 1, pageSize = 20) => {
    setLoading(true)
    try {
      // For regular users, don't pass pagination params (they get all authorized apps)
      // For admins, pass pagination params
      const userInfo = storage.getUserInfo()
      const admin = userInfo?.is_admin === true
      const response = admin
        ? await applicationService.list(page, pageSize)
        : await applicationService.list()
      setApplications(response.data)
      setPagination({
        current: response.page || 1,
        pageSize: response.size || response.data.length,
        total: response.total || response.data.length,
      })
    } catch (error: any) {
      message.error(
        error.response?.data?.error || 'Failed to load applications'
      )
    } finally {
      setLoading(false)
    }
  }

  const handleDelete = async (id: number) => {
    try {
      await applicationService.delete(id)
      message.success(t('common.success'))
      fetchApplications(pagination.current, pagination.pageSize)
    } catch (error: any) {
      message.error(
        error.response?.data?.error || 'Failed to delete application'
      )
    }
  }

  const handleCreate = (e?: React.MouseEvent) => {
    e?.preventDefault()
    e?.stopPropagation()
    setCreateModalVisible(true)
    form.resetFields()
  }

  const handleCreateSubmit = async (values: CreateApplicationRequest) => {
    try {
      await applicationService.create(values)
      message.success(t('common.success'))
      setCreateModalVisible(false)
      form.resetFields()
      fetchApplications(pagination.current, pagination.pageSize)
    } catch (error: any) {
      message.error(
        error.response?.data?.error || 'Failed to create application'
      )
    }
  }

  const getTypeColor = (type: string) => {
    switch (type) {
      case 'saml':
        return 'blue'
      case 'oidc':
        return 'green'
      case 'oauth2':
        return 'orange'
      default:
        return 'default'
    }
  }

  // Get columns based on admin status
  const getColumns = (): ColumnsType<Application> => {
    const userInfo = storage.getUserInfo()
    const admin = userInfo?.is_admin === true

    const baseColumns: ColumnsType<Application> = [
      {
        title: t('applications.name'),
        dataIndex: 'name',
        key: 'name',
      },
      {
        title: t('applications.type'),
        dataIndex: 'type',
        key: 'type',
        render: (type: string) => (
          <Tag color={getTypeColor(type)}>{type.toUpperCase()}</Tag>
        ),
      },
      {
        title: t('applications.status'),
        dataIndex: 'status',
        key: 'status',
        render: (status: string) => (
          <Tag color={status === 'active' ? 'green' : 'red'}>
            {status}
          </Tag>
        ),
      },
    ]

    if (admin) {
      baseColumns.push({
        title: t('applications.actions'),
        key: 'actions',
        render: (_: any, record: Application) => (
          <Space>
            <Button type="link">{t('common.edit')}</Button>
            <Button
              type="link"
              danger
              onClick={() => handleDelete(record.id)}
            >
              {t('common.delete')}
            </Button>
          </Space>
        ),
      })
    }

    return baseColumns
  }

  return (
    <div>
      <div
        style={{
          marginBottom: 16,
          display: 'flex',
          justifyContent: 'space-between',
          alignItems: 'center',
        }}
      >
        <h1 style={{ fontSize: '28px', fontWeight: 600, color: '#1e293b', margin: 0 }}>
          {t('applications.title')}
        </h1>
        {isAdmin && (
          <Button type="primary" icon={<PlusOutlined />} onClick={handleCreate}>
            {t('applications.createApplication')}
          </Button>
        )}
      </div>
      <Table
        columns={getColumns()}
        dataSource={applications}
        loading={loading}
        rowKey="id"
        pagination={
          isAdmin
            ? {
                current: pagination.current,
                pageSize: pagination.pageSize,
                total: pagination.total,
                onChange: (page, pageSize) => {
                  fetchApplications(page, pageSize)
                },
              }
            : false
        }
      />
      <Modal
        title={t('applications.createApplication')}
        open={createModalVisible}
        onCancel={() => {
          setCreateModalVisible(false)
          form.resetFields()
        }}
        footer={null}
        width={600}
      >
        <Form
          form={form}
          layout="vertical"
          onFinish={handleCreateSubmit}
          autoComplete="off"
        >
          <Form.Item
            name="name"
            label={t('applications.name')}
            rules={[{ required: true, message: t('applications.name') + ' ' + t('common.error') }]}
          >
            <Input />
          </Form.Item>
          <Form.Item
            name="type"
            label={t('applications.type')}
            rules={[{ required: true, message: t('applications.type') + ' ' + t('common.error') }]}
          >
            <Select>
              <Select.Option value="saml">SAML</Select.Option>
              <Select.Option value="oidc">OIDC</Select.Option>
              <Select.Option value="oauth2">OAuth2</Select.Option>
            </Select>
          </Form.Item>
          <Form.Item name="description" label={t('applications.description')}>
            <Input.TextArea rows={3} />
          </Form.Item>
          <Form.Item name="icon" label={t('applications.icon')}>
            <Input placeholder="URL or icon name" />
          </Form.Item>
          <Form.Item>
            <Space>
              <Button onClick={() => {
                setCreateModalVisible(false)
                form.resetFields()
              }}>
                {t('common.cancel')}
              </Button>
              <Button type="primary" htmlType="submit">
                {t('common.create')}
              </Button>
            </Space>
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}

export default ApplicationsPage
