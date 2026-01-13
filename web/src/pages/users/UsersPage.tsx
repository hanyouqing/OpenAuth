import { useState, useEffect } from 'react'
import { Table, Button, Space, message, Modal, Form, Input } from 'antd'
import { PlusOutlined } from '@ant-design/icons'
import { useTranslation } from 'react-i18next'
import { userService, User, CreateUserRequest } from '@/services/userService'

function UsersPage() {
  const { t } = useTranslation()
  const [users, setUsers] = useState<User[]>([])
  const [loading, setLoading] = useState(false)
  const [createModalVisible, setCreateModalVisible] = useState(false)
  const [form] = Form.useForm()
  const [pagination, setPagination] = useState({
    current: 1,
    pageSize: 20,
    total: 0,
  })

  const fetchUsers = async (page = 1, pageSize = 20) => {
    setLoading(true)
    try {
      const response = await userService.list(page, pageSize)
      setUsers(response.data)
      setPagination({
        current: response.page,
        pageSize: response.size,
        total: response.total,
      })
    } catch (error: any) {
      message.error(error.response?.data?.error || 'Failed to load users')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchUsers()
  }, [])

  const handleDelete = async (id: number) => {
    try {
      await userService.delete(id)
      message.success(t('common.success'))
      fetchUsers(pagination.current, pagination.pageSize)
    } catch (error: any) {
      message.error(error.response?.data?.error || 'Failed to delete user')
    }
  }

  const handleCreate = (e?: React.MouseEvent) => {
    e?.preventDefault()
    e?.stopPropagation()
    setCreateModalVisible(true)
    form.resetFields()
  }

  const handleCreateSubmit = async (values: CreateUserRequest) => {
    try {
      await userService.create(values)
      message.success(t('common.success'))
      setCreateModalVisible(false)
      form.resetFields()
      fetchUsers(pagination.current, pagination.pageSize)
    } catch (error: any) {
      message.error(error.response?.data?.error || 'Failed to create user')
    }
  }

  const columns = [
    {
      title: t('users.username'),
      dataIndex: 'username',
      key: 'username',
    },
    {
      title: t('users.email'),
      dataIndex: 'email',
      key: 'email',
    },
    {
      title: t('users.status'),
      dataIndex: 'status',
      key: 'status',
    },
    {
      title: t('users.actions'),
      key: 'actions',
      render: (_: any, record: User) => (
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
    },
  ]

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
          {t('users.title')}
        </h1>
        <Button type="primary" icon={<PlusOutlined />} onClick={handleCreate}>
          {t('users.createUser')}
        </Button>
      </div>
      <Table
        columns={columns}
        dataSource={users}
        loading={loading}
        rowKey="id"
        pagination={{
          current: pagination.current,
          pageSize: pagination.pageSize,
          total: pagination.total,
          onChange: (page, pageSize) => {
            fetchUsers(page, pageSize)
          },
        }}
      />
      <Modal
        title={t('users.createUser')}
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
            name="username"
            label={t('users.username')}
            rules={[{ required: true, message: t('users.username') + ' ' + t('common.error') }]}
          >
            <Input />
          </Form.Item>
          <Form.Item
            name="email"
            label={t('users.email')}
            rules={[
              { required: true, message: t('users.email') + ' ' + t('common.error') },
              { type: 'email', message: t('users.email') + ' ' + t('common.error') },
            ]}
          >
            <Input />
          </Form.Item>
          <Form.Item
            name="password"
            label={t('users.password')}
            rules={[
              { required: true, message: t('users.password') + ' ' + t('common.error') },
              { min: 8, message: t('users.password') + ' ' + t('common.error') },
            ]}
          >
            <Input.Password />
          </Form.Item>
          <Form.Item name="first_name" label={t('users.firstName')}>
            <Input />
          </Form.Item>
          <Form.Item name="last_name" label={t('users.lastName')}>
            <Input />
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

export default UsersPage
