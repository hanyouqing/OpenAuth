import { useState } from 'react'
import { Modal, Form, Input, Button, message } from 'antd'
import { LockOutlined } from '@ant-design/icons'
import { useTranslation } from 'react-i18next'
import { userService } from '@/services/userService'

interface ChangePasswordModalProps {
  open: boolean
  onCancel: () => void
  onSuccess?: () => void
}

function ChangePasswordModal({ open, onCancel, onSuccess }: ChangePasswordModalProps) {
  const { t } = useTranslation()
  const [form] = Form.useForm()
  const [loading, setLoading] = useState(false)

  const handleSubmit = async (values: {
    current_password: string
    new_password: string
    confirm_password: string
  }) => {
    if (values.new_password !== values.confirm_password) {
      message.error(t('settings.security.passwordMismatch'))
      return
    }

    setLoading(true)
    try {
      await userService.changePassword({
        current_password: values.current_password,
        new_password: values.new_password,
      })
      message.success(t('settings.security.passwordChanged'))
      form.resetFields()
      onCancel()
      onSuccess?.()
    } catch (error: any) {
      message.error(error.response?.data?.error || t('common.error'))
    } finally {
      setLoading(false)
    }
  }

  return (
    <Modal
      title={t('settings.security.changePassword')}
      open={open}
      onCancel={onCancel}
      footer={null}
      width={500}
    >
      <Form
        form={form}
        layout="vertical"
        onFinish={handleSubmit}
        autoComplete="off"
      >
        <Form.Item
          label={t('settings.security.currentPassword')}
          name="current_password"
          rules={[
            {
              required: true,
              message: t('settings.security.currentPassword') + ' ' + t('common.error'),
            },
          ]}
        >
          <Input.Password
            prefix={<LockOutlined />}
            placeholder={t('settings.security.currentPassword')}
          />
        </Form.Item>

        <Form.Item
          label={t('settings.security.newPassword')}
          name="new_password"
          rules={[
            {
              required: true,
              message: t('settings.security.newPassword') + ' ' + t('common.error'),
            },
            {
              min: 8,
              message: t('settings.security.newPassword') + ' ' + t('common.error'),
            },
          ]}
        >
          <Input.Password
            prefix={<LockOutlined />}
            placeholder={t('settings.security.newPassword')}
          />
        </Form.Item>

        <Form.Item
          label={t('settings.security.confirmPassword')}
          name="confirm_password"
          dependencies={['new_password']}
          rules={[
            {
              required: true,
              message: t('settings.security.confirmPassword') + ' ' + t('common.error'),
            },
            ({ getFieldValue }) => ({
              validator(_, value) {
                if (!value || getFieldValue('new_password') === value) {
                  return Promise.resolve()
                }
                return Promise.reject(
                  new Error(t('settings.security.passwordMismatch'))
                )
              },
            }),
          ]}
        >
          <Input.Password
            prefix={<LockOutlined />}
            placeholder={t('settings.security.confirmPassword')}
          />
        </Form.Item>

        <Form.Item style={{ marginBottom: 0, textAlign: 'right' }}>
          <Button onClick={onCancel} style={{ marginRight: 8 }}>
            {t('common.cancel')}
          </Button>
          <Button type="primary" htmlType="submit" loading={loading}>
            {t('settings.security.changePassword')}
          </Button>
        </Form.Item>
      </Form>
    </Modal>
  )
}

export default ChangePasswordModal
