import { useState, useEffect } from 'react'
import { Form, Input, Button, Card, message, Upload, Avatar, Space, Spin } from 'antd'
import { UserOutlined, UploadOutlined } from '@ant-design/icons'
import { useTranslation } from 'react-i18next'
import { userService, User } from '@/services/userService'

function ProfileSettingsPage() {
  const { t } = useTranslation()
  const [form] = Form.useForm()
  const [loading, setLoading] = useState(false)
  const [user, setUser] = useState<User | null>(null)
  const [avatarUrl, setAvatarUrl] = useState<string>('')

  useEffect(() => {
    fetchUser()
  }, [])

  const fetchUser = async () => {
    try {
      setLoading(true)
      const userData = await userService.getMe()
      setUser(userData)
      form.setFieldsValue({
        username: userData.username,
        email: userData.email,
        first_name: userData.first_name || '',
        last_name: userData.last_name || '',
      })
      // Use logo as default avatar if no avatar is set
      setAvatarUrl(userData.avatar && userData.avatar.trim() !== '' ? userData.avatar : '/logo.svg')
    } catch (error: any) {
      console.error('Failed to fetch user:', error)
      const errorMessage = error.response?.data?.error || error.message || t('common.error')
      message.error(errorMessage)
    } finally {
      setLoading(false)
    }
  }

  const handleSubmit = async (values: any) => {
    setLoading(true)
    try {
      const updatedUser = await userService.updateMe({
        email: values.email,
        first_name: values.first_name,
        last_name: values.last_name,
      })
      setUser(updatedUser)
      message.success(t('settings.profile.profileUpdated'))
    } catch (error: any) {
      message.error(error.response?.data?.error || t('common.error'))
    } finally {
      setLoading(false)
    }
  }

  const handleAvatarChange = async (file: File) => {
    try {
      // Convert file to base64
      const reader = new FileReader()
      reader.onloadend = async () => {
        const base64String = reader.result as string
        try {
          const updatedUser = await userService.uploadAvatar(base64String)
          setUser(updatedUser)
          setAvatarUrl(base64String)
          message.success(t('settings.profile.avatarUpdated'))
        } catch (error: any) {
          message.error(error.response?.data?.error || t('common.error'))
        }
      }
      reader.readAsDataURL(file)
    } catch (error) {
      message.error(t('common.error'))
    }
  }

  const beforeUpload = (file: File) => {
    const isImage = file.type.startsWith('image/')
    if (!isImage) {
      message.error(t('settings.profile.invalidImage'))
      return false
    }
    const isLt2M = file.size / 1024 / 1024 < 2
    if (!isLt2M) {
      message.error(t('settings.profile.imageTooLarge'))
      return false
    }
    handleAvatarChange(file)
    return false // Prevent auto upload
  }

  if (loading && !user) {
    return (
      <Card>
        <div style={{ textAlign: 'center', padding: 40 }}>
          <Spin size="large" />
        </div>
      </Card>
    )
  }

  return (
    <Card>
      <Space direction="vertical" size="large" style={{ width: '100%' }}>
        <div>
          <h3 style={{ marginBottom: 16 }}>{t('settings.profile.avatar')}</h3>
          <Space direction="vertical" align="center">
            <Avatar
              src={avatarUrl || '/logo.svg'}
              size={120}
              icon={<UserOutlined />}
              style={{ border: '2px solid #e5e7eb' }}
            />
            <Upload
              beforeUpload={beforeUpload}
              showUploadList={false}
              accept="image/*"
            >
              <Button icon={<UploadOutlined />}>{t('settings.profile.uploadAvatar')}</Button>
            </Upload>
          </Space>
        </div>

        <Form
          form={form}
          layout="vertical"
          onFinish={handleSubmit}
          autoComplete="off"
        >
          <Form.Item
            label={t('settings.profile.username')}
            name="username"
          >
            <Input disabled value={user?.username} />
          </Form.Item>

          <Form.Item
            label={t('settings.profile.email')}
            name="email"
            rules={[
              { required: true, message: t('settings.profile.email') + ' ' + t('common.error') },
              { type: 'email', message: t('settings.profile.email') + ' ' + t('common.error') },
            ]}
          >
            <Input placeholder={t('settings.profile.email')} />
          </Form.Item>

          <Form.Item
            label={t('settings.profile.firstName')}
            name="first_name"
          >
            <Input placeholder={t('settings.profile.firstName')} />
          </Form.Item>

          <Form.Item
            label={t('settings.profile.lastName')}
            name="last_name"
          >
            <Input placeholder={t('settings.profile.lastName')} />
          </Form.Item>

          <Form.Item>
            <Button type="primary" htmlType="submit" loading={loading}>
              {t('settings.profile.saveChanges')}
            </Button>
          </Form.Item>
        </Form>
      </Space>
    </Card>
  )
}

export default ProfileSettingsPage
