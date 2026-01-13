import { useState } from 'react'
import { Form, Input, Button, Card, message } from 'antd'
import { UserOutlined, LockOutlined } from '@ant-design/icons'
import { useTranslation } from 'react-i18next'
import { authService } from '@/services/authService'

function LoginPage() {
  const { t } = useTranslation()
  const [loading, setLoading] = useState(false)

  const onFinish = async (values: { username: string; password: string }) => {
    setLoading(true)
    try {
      const response = await authService.login(values)
      if (response && response.access_token) {
        message.success(t('common.success') || 'Login successful')
        
        // Show password expiration warning if present
        if (response.password_expiration_warning?.warning) {
          // Store warning in sessionStorage to show in the main app
          sessionStorage.setItem(
            'password_expiration_warning',
            JSON.stringify(response.password_expiration_warning)
          )
        }
        
        // Use window.location to force a full page reload and re-check auth
        window.location.href = '/'
      } else {
        message.error('Login failed: No token received')
      }
    } catch (error: any) {
      console.error('Login error:', error)
      const errorMessage =
        error.response?.data?.error || error.message || t('common.error') || 'Login failed'
      message.error(errorMessage)
    } finally {
      setLoading(false)
    }
  }

  return (
    <div
      style={{
        display: 'flex',
        flexDirection: 'column',
        justifyContent: 'center',
        alignItems: 'center',
        minHeight: '100vh',
        background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
        padding: '20px',
      }}
    >
      <Card
        style={{
          width: '100%',
          maxWidth: 400,
          boxShadow: '0 20px 25px -5px rgba(0, 0, 0, 0.1), 0 10px 10px -5px rgba(0, 0, 0, 0.04)',
        }}
      >
        {/* Logo and Title Section inside Card */}
        <div
          style={{
            display: 'flex',
            flexDirection: 'column',
            alignItems: 'center',
            gap: 12,
            marginBottom: 32,
            paddingTop: 8,
          }}
        >
          <img
            src="/logo.svg"
            alt="OpenAuth Logo"
            style={{
              width: 64,
              height: 64,
            }}
          />
          <h1
            style={{
              fontSize: '28px',
              fontWeight: 700,
              margin: 0,
              color: '#1e293b',
              letterSpacing: '-0.5px',
            }}
          >
            OpenAuth
          </h1>
          <div
            style={{
              fontSize: '16px',
              fontWeight: 500,
              color: '#64748b',
              marginTop: -4,
            }}
          >
            {t('auth.login')}
          </div>
        </div>

        <Form
          name="login"
          onFinish={onFinish}
          autoComplete="off"
          size="large"
        >
          <Form.Item
            name="username"
            rules={[{ required: true, message: 'Please input your username!' }]}
          >
            <Input
              prefix={<UserOutlined />}
              placeholder={t('auth.username')}
            />
          </Form.Item>

          <Form.Item
            name="password"
            rules={[{ required: true, message: 'Please input your password!' }]}
          >
            <Input.Password
              prefix={<LockOutlined />}
              placeholder={t('auth.password')}
            />
          </Form.Item>

          <Form.Item>
            <Button type="primary" htmlType="submit" block loading={loading}>
              {t('auth.login')}
            </Button>
          </Form.Item>
        </Form>
      </Card>
    </div>
  )
}

export default LoginPage
