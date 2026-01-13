import { useState, useEffect } from 'react'
import {
  Card,
  Form,
  Input,
  Button,
  Switch,
  Table,
  Space,
  Modal,
  message,
  Popconfirm,
  Select,
  Tag,
  Typography,
  Divider,
} from 'antd'
import {
  PlusOutlined,
  EditOutlined,
  DeleteOutlined,
  SafetyOutlined,
} from '@ant-design/icons'
import { useTranslation } from 'react-i18next'
import {
  whitelistService,
  WhitelistEntry,
  WhitelistType,
  CreateWhitelistEntryRequest,
} from '@/services/whitelistService'
import type { ColumnsType } from 'antd/es/table'

const { Title, Text } = Typography
const { TextArea } = Input

function SecuritySettingsPage() {
  const { t } = useTranslation()
  const [whitelistEnabled, setWhitelistEnabled] = useState(false)
  const [whitelistLoading, setWhitelistLoading] = useState(false)
  const [entries, setEntries] = useState<WhitelistEntry[]>([])
  const [entriesLoading, setEntriesLoading] = useState(false)
  const [modalVisible, setModalVisible] = useState(false)
  const [editingEntry, setEditingEntry] = useState<WhitelistEntry | null>(null)
  const [form] = Form.useForm()

  useEffect(() => {
    fetchWhitelistPolicy()
    fetchEntries()
  }, [])

  const fetchWhitelistPolicy = async () => {
    setWhitelistLoading(true)
    try {
      const policy = await whitelistService.getPolicy()
      setWhitelistEnabled(policy.enabled)
    } catch (error: any) {
      message.error(error.response?.data?.error || t('common.error'))
    } finally {
      setWhitelistLoading(false)
    }
  }

  const fetchEntries = async () => {
    setEntriesLoading(true)
    try {
      const data = await whitelistService.listEntries()
      setEntries(data)
    } catch (error: any) {
      message.error(error.response?.data?.error || t('common.error'))
    } finally {
      setEntriesLoading(false)
    }
  }

  const handleWhitelistToggle = async (enabled: boolean) => {
    setWhitelistLoading(true)
    try {
      await whitelistService.updatePolicy(enabled)
      setWhitelistEnabled(enabled)
      message.success(t('system.security.whitelist.policyUpdated'))
    } catch (error: any) {
      message.error(error.response?.data?.error || t('common.error'))
    } finally {
      setWhitelistLoading(false)
    }
  }

  const handleCreate = () => {
    setEditingEntry(null)
    form.resetFields()
    setModalVisible(true)
  }

  const handleEdit = (entry: WhitelistEntry) => {
    setEditingEntry(entry)
    form.setFieldsValue({
      type: entry.type,
      value: entry.value,
      description: entry.description,
      enabled: entry.enabled,
    })
    setModalVisible(true)
  }

  const handleDelete = async (id: number) => {
    try {
      await whitelistService.deleteEntry(id)
      message.success(t('common.success'))
      fetchEntries()
    } catch (error: any) {
      message.error(error.response?.data?.error || t('common.error'))
    }
  }

  const handleSubmit = async (values: any) => {
    try {
      if (editingEntry) {
        await whitelistService.updateEntry(editingEntry.id, values)
        message.success(t('common.success'))
      } else {
        await whitelistService.createEntry(values as CreateWhitelistEntryRequest)
        message.success(t('common.success'))
      }
      setModalVisible(false)
      form.resetFields()
      fetchEntries()
    } catch (error: any) {
      message.error(error.response?.data?.error || t('common.error'))
    }
  }

  const columns: ColumnsType<WhitelistEntry> = [
    {
      title: t('system.security.whitelist.type'),
      dataIndex: 'type',
      key: 'type',
      render: (type: WhitelistType) => (
        <Tag color={type === 'ip' ? 'blue' : 'green'}>
          {type === 'ip' ? t('system.security.whitelist.ip') : t('system.security.whitelist.region')}
        </Tag>
      ),
    },
    {
      title: t('system.security.whitelist.value'),
      dataIndex: 'value',
      key: 'value',
    },
    {
      title: t('system.security.whitelist.description'),
      dataIndex: 'description',
      key: 'description',
    },
    {
      title: t('system.security.whitelist.enabled'),
      dataIndex: 'enabled',
      key: 'enabled',
      render: (enabled: boolean) => (
        <Tag color={enabled ? 'success' : 'default'}>
          {enabled ? t('common.yes') : t('common.no')}
        </Tag>
      ),
    },
    {
      title: t('common.edit'),
      key: 'actions',
      render: (_, record) => (
        <Space>
          <Button
            type="link"
            icon={<EditOutlined />}
            onClick={() => handleEdit(record)}
          >
            {t('common.edit')}
          </Button>
          <Popconfirm
            title={t('common.delete')}
            description={t('system.security.whitelist.deleteConfirm')}
            onConfirm={() => handleDelete(record.id)}
            okText={t('common.yes')}
            cancelText={t('common.no')}
          >
            <Button type="link" danger icon={<DeleteOutlined />}>
              {t('common.delete')}
            </Button>
          </Popconfirm>
        </Space>
      ),
    },
  ]

  return (
    <Space direction="vertical" style={{ width: '100%' }} size="large">
      <Card>
        <Title level={4}>
          <SafetyOutlined /> {t('system.security.whitelist.title')}
        </Title>
        <Space direction="vertical" style={{ width: '100%' }} size="middle">
          <div>
            <Text strong>{t('system.security.whitelist.enableWhitelist')}</Text>
            <Switch
              checked={whitelistEnabled}
              onChange={handleWhitelistToggle}
              loading={whitelistLoading}
              style={{ marginLeft: 16 }}
            />
            <Text type="secondary" style={{ marginLeft: 16 }}>
              {t('system.security.whitelist.enableDescription')}
            </Text>
          </div>
          <Divider />
          <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
            <Title level={5}>{t('system.security.whitelist.entries')}</Title>
            <Button
              type="primary"
              icon={<PlusOutlined />}
              onClick={handleCreate}
              disabled={!whitelistEnabled}
            >
              {t('system.security.whitelist.addEntry')}
            </Button>
          </div>
          <Table
            columns={columns}
            dataSource={entries}
            loading={entriesLoading}
            rowKey="id"
            pagination={{ pageSize: 10 }}
          />
        </Space>
      </Card>

      <Modal
        title={
          editingEntry
            ? t('system.security.whitelist.editEntry')
            : t('system.security.whitelist.addEntry')
        }
        open={modalVisible}
        onCancel={() => {
          setModalVisible(false)
          form.resetFields()
        }}
        footer={null}
        width={600}
      >
        <Form
          form={form}
          layout="vertical"
          onFinish={handleSubmit}
          initialValues={{ enabled: true }}
        >
          <Form.Item
            name="type"
            label={t('system.security.whitelist.type')}
            rules={[{ required: true }]}
          >
            <Select>
              <Select.Option value="ip">{t('system.security.whitelist.ip')}</Select.Option>
              <Select.Option value="region">{t('system.security.whitelist.region')}</Select.Option>
            </Select>
          </Form.Item>

          <Form.Item
            name="value"
            label={t('system.security.whitelist.value')}
            rules={[{ required: true }]}
            extra={
              form.getFieldValue('type') === 'ip'
                ? t('system.security.whitelist.ipHint')
                : t('system.security.whitelist.regionHint')
            }
          >
            <Input placeholder={form.getFieldValue('type') === 'ip' ? '192.168.1.1 or 192.168.1.0/24' : 'CN, US, etc.'} />
          </Form.Item>

          <Form.Item name="description" label={t('system.security.whitelist.description')}>
            <TextArea rows={3} />
          </Form.Item>

          <Form.Item name="enabled" valuePropName="checked">
            <Switch checkedChildren={t('common.yes')} unCheckedChildren={t('common.no')} />
          </Form.Item>

          <Form.Item>
            <Space>
              <Button onClick={() => setModalVisible(false)}>{t('common.cancel')}</Button>
              <Button type="primary" htmlType="submit">
                {t('common.save')}
              </Button>
            </Space>
          </Form.Item>
        </Form>
      </Modal>
    </Space>
  )
}

export default SecuritySettingsPage
