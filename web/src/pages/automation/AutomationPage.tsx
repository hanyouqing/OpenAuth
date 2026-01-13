import { useState, useEffect } from 'react'
import {
  Table,
  Button,
  Space,
  message,
  Tag,
  Modal,
  Form,
  Input,
  Switch,
  InputNumber,
} from 'antd'
import {
  PlusOutlined,
  EditOutlined,
  DeleteOutlined,
  PlayCircleOutlined,
  PauseCircleOutlined,
} from '@ant-design/icons'
import type { ColumnsType } from 'antd/es/table'
import {
  automationService,
  AutomationWorkflow,
  AutomationExecution,
  CreateWorkflowRequest,
  UpdateWorkflowRequest,
  AutomationTrigger,
  AutomationAction,
} from '@/services/automationService'
import dayjs from 'dayjs'

const { TextArea } = Input

function AutomationPage() {
  const [workflows, setWorkflows] = useState<AutomationWorkflow[]>([])
  const [loading, setLoading] = useState(false)
  const [modalVisible, setModalVisible] = useState(false)
  const [executionModalVisible, setExecutionModalVisible] = useState(false)
  const [selectedWorkflow, setSelectedWorkflow] = useState<AutomationWorkflow | null>(null)
  const [executions, setExecutions] = useState<AutomationExecution[]>([])
  const [form] = Form.useForm()
  const [editingWorkflow, setEditingWorkflow] = useState<AutomationWorkflow | null>(null)

  const fetchWorkflows = async () => {
    setLoading(true)
    try {
      const data = await automationService.list()
      setWorkflows(data)
    } catch (error: any) {
      message.error(error.response?.data?.error || 'Failed to load workflows')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchWorkflows()
  }, [])

  const handleCreate = () => {
    setEditingWorkflow(null)
    form.resetFields()
    form.setFieldsValue({
      trigger: { type: 'event', event: 'user.created' },
      actions: [{ type: 'send_email', config: {}, on_error: 'continue' }],
      priority: 0,
      enabled: true,
    })
    setModalVisible(true)
  }

  const handleEdit = (workflow: AutomationWorkflow) => {
    setEditingWorkflow(workflow)
    form.setFieldsValue({
      name: workflow.name,
      description: workflow.description,
      trigger: workflow.trigger,
      actions: workflow.actions,
      priority: workflow.priority,
      enabled: workflow.enabled,
    })
    setModalVisible(true)
  }

  const handleSubmit = async (values: any) => {
    try {
      if (editingWorkflow) {
        const updateData: UpdateWorkflowRequest = {
          name: values.name,
          description: values.description,
          trigger: values.trigger,
          actions: values.actions,
          enabled: values.enabled,
          priority: values.priority,
        }
        await automationService.update(editingWorkflow.id, updateData)
        message.success('Workflow updated successfully')
      } else {
        const createData: CreateWorkflowRequest = {
          name: values.name,
          description: values.description,
          trigger: values.trigger,
          actions: values.actions,
          priority: values.priority || 0,
        }
        await automationService.create(createData)
        message.success('Workflow created successfully')
      }
      setModalVisible(false)
      fetchWorkflows()
    } catch (error: any) {
      message.error(error.response?.data?.error || 'Failed to save workflow')
    }
  }

  const handleDelete = async (id: number) => {
    try {
      await automationService.delete(id)
      message.success('Workflow deleted successfully')
      fetchWorkflows()
    } catch (error: any) {
      message.error(error.response?.data?.error || 'Failed to delete workflow')
    }
  }

  const handleToggleEnabled = async (workflow: AutomationWorkflow) => {
    try {
      await automationService.update(workflow.id, { enabled: !workflow.enabled })
      message.success(`Workflow ${workflow.enabled ? 'disabled' : 'enabled'} successfully`)
      fetchWorkflows()
    } catch (error: any) {
      message.error(error.response?.data?.error || 'Failed to update workflow')
    }
  }

  const handleViewExecutions = async (workflow: AutomationWorkflow) => {
    setSelectedWorkflow(workflow)
    try {
      const data = await automationService.getExecutions(workflow.id, 20)
      setExecutions(data)
      setExecutionModalVisible(true)
    } catch (error: any) {
      message.error(error.response?.data?.error || 'Failed to load executions')
    }
  }

  const columns: ColumnsType<AutomationWorkflow> = [
    {
      title: 'Name',
      dataIndex: 'name',
      key: 'name',
      render: (text, record) => (
        <div>
          <div style={{ fontWeight: 500 }}>{text}</div>
          {record.description && (
            <div style={{ fontSize: '12px', color: '#666' }}>{record.description}</div>
          )}
        </div>
      ),
    },
    {
      title: 'Trigger',
      key: 'trigger',
      render: (_, record) => {
        const trigger = record.trigger as AutomationTrigger
        return (
          <Tag color="blue">
            {trigger.type === 'event' ? `Event: ${trigger.event}` : trigger.type}
          </Tag>
        )
      },
    },
    {
      title: 'Actions',
      key: 'actions',
      render: (_, record) => {
        const actions = record.actions as AutomationAction[]
        return <span>{actions.length} action(s)</span>
      },
    },
    {
      title: 'Priority',
      dataIndex: 'priority',
      key: 'priority',
      align: 'center',
    },
    {
      title: 'Status',
      key: 'status',
      render: (_, record) => (
        <Tag color={record.enabled ? 'green' : 'default'}>
          {record.enabled ? 'Enabled' : 'Disabled'}
        </Tag>
      ),
    },
    {
      title: 'Created',
      key: 'created_at',
      render: (_, record) => dayjs(record.created_at).format('YYYY-MM-DD HH:mm'),
    },
    {
      title: 'Actions',
      key: 'actions',
      render: (_, record) => (
        <Space>
          <Button
            size="small"
            icon={record.enabled ? <PauseCircleOutlined /> : <PlayCircleOutlined />}
            onClick={() => handleToggleEnabled(record)}
          >
            {record.enabled ? 'Disable' : 'Enable'}
          </Button>
          <Button size="small" icon={<EditOutlined />} onClick={() => handleEdit(record)}>
            Edit
          </Button>
          <Button
            size="small"
            onClick={() => handleViewExecutions(record)}
          >
            Executions
          </Button>
          <Button
            size="small"
            danger
            icon={<DeleteOutlined />}
            onClick={() => handleDelete(record.id)}
          >
            Delete
          </Button>
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
          Automation Workflows
        </h1>
        <Button type="primary" icon={<PlusOutlined />} onClick={handleCreate}>
          Create Workflow
        </Button>
      </div>
      <Table
        columns={columns}
        dataSource={workflows}
        rowKey="id"
        loading={loading}
        pagination={{
          pageSize: 20,
          showSizeChanger: true,
          showTotal: (total) => `Total ${total} workflows`,
        }}
      />

      <Modal
        title={editingWorkflow ? 'Edit Workflow' : 'Create Workflow'}
        open={modalVisible}
        onCancel={() => setModalVisible(false)}
        onOk={() => form.submit()}
        width={800}
      >
        <Form form={form} onFinish={handleSubmit} layout="vertical">
          <Form.Item name="name" label="Name" rules={[{ required: true }]}>
            <Input />
          </Form.Item>
          <Form.Item name="description" label="Description">
            <TextArea rows={2} />
          </Form.Item>
          <Form.Item name="priority" label="Priority">
            <InputNumber min={0} />
          </Form.Item>
          <Form.Item name="enabled" label="Enabled" valuePropName="checked">
            <Switch />
          </Form.Item>
          <Form.Item
            name="trigger"
            label="Trigger"
            rules={[{ required: true }]}
            tooltip="Trigger configuration (JSON format)"
          >
            <TextArea rows={4} placeholder='{"type":"event","event":"user.created"}' />
          </Form.Item>
          <Form.Item
            name="actions"
            label="Actions"
            rules={[{ required: true }]}
            tooltip="Actions array (JSON format)"
          >
            <TextArea
              rows={6}
              placeholder='[{"type":"send_email","config":{"subject":"Welcome","body":"Hello"}}]'
            />
          </Form.Item>
        </Form>
      </Modal>

      <Modal
        title={`Executions - ${selectedWorkflow?.name}`}
        open={executionModalVisible}
        onCancel={() => setExecutionModalVisible(false)}
        footer={null}
        width={1000}
      >
        <Table
          columns={[
            { title: 'ID', dataIndex: 'id', key: 'id' },
            {
              title: 'Status',
              key: 'status',
              render: (_, record) => (
                <Tag
                  color={
                    record.status === 'completed'
                      ? 'green'
                      : record.status === 'failed'
                      ? 'red'
                      : 'blue'
                  }
                >
                  {record.status}
                </Tag>
              ),
            },
            {
              title: 'Started',
              key: 'started_at',
              render: (_, record) =>
                record.started_at ? dayjs(record.started_at).format('YYYY-MM-DD HH:mm:ss') : '-',
            },
            {
              title: 'Completed',
              key: 'completed_at',
              render: (_, record) =>
                record.completed_at
                  ? dayjs(record.completed_at).format('YYYY-MM-DD HH:mm:ss')
                  : '-',
            },
            {
              title: 'Error',
              dataIndex: 'error',
              key: 'error',
              render: (text) => (text ? <span style={{ color: 'red' }}>{text}</span> : '-'),
            },
          ]}
          dataSource={executions}
          rowKey="id"
          pagination={false}
          size="small"
        />
      </Modal>
    </div>
  )
}

export default AutomationPage
