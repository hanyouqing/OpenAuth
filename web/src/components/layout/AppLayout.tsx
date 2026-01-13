import { useState, useEffect, useMemo } from 'react'
import { Layout, Menu, Avatar, Dropdown, Button, Alert } from 'antd'
import {
  MenuFoldOutlined,
  MenuUnfoldOutlined,
  DashboardOutlined,
  UserOutlined,
  AppstoreOutlined,
  LogoutOutlined,
  GlobalOutlined,
  InfoCircleOutlined,
  GithubOutlined,
  ExclamationCircleOutlined,
  LockOutlined,
  ProfileOutlined,
  MobileOutlined,
  ThunderboltOutlined,
} from '@ant-design/icons'
import { useTranslation } from 'react-i18next'
import { useNavigate, useLocation } from 'react-router-dom'
import type { MenuProps } from 'antd'
import { authService, PasswordExpirationWarning } from '@/services/authService'
import { storage } from '@/utils/storage'
import { systemService, VersionInfo } from '@/services/systemService'
import ChangePasswordModal from '@/components/ChangePasswordModal'

const { Header, Sider, Content, Footer } = Layout

interface AppLayoutProps {
  children: React.ReactNode
}

function AppLayout({ children }: AppLayoutProps) {
  const [collapsed, setCollapsed] = useState(false)
  const [versionInfo, setVersionInfo] = useState<VersionInfo | null>(null)
  const [passwordWarning, setPasswordWarning] = useState<PasswordExpirationWarning | null>(null)
  const [isAdmin, setIsAdmin] = useState(false)
  const [changePasswordModalOpen, setChangePasswordModalOpen] = useState(false)
  const { t, i18n } = useTranslation()
  const navigate = useNavigate()
  const location = useLocation()

  // Check admin status on mount and when location changes
  useEffect(() => {
    const checkAdminStatus = async () => {
      // First check from storage
      const userInfo = storage.getUserInfo()
      if (userInfo && userInfo.is_admin !== undefined) {
        const admin = userInfo.is_admin === true || userInfo.is_admin === 'true'
        setIsAdmin(admin)
        // Debug log (can be removed in production)
        if (import.meta.env.DEV) {
          console.log('User info from storage:', userInfo, 'isAdmin:', admin)
        }
      } else {
        // If no user info in storage or is_admin not set, fetch from API
        try {
          const { userService } = await import('@/services/userService')
          const me = await userService.getMe()
          // Check if user has admin role or is_admin flag
          const admin = me.is_admin === true || me.roles?.includes('admin') || false
          setIsAdmin(admin)
          // Update storage with latest info including is_admin
          storage.setUserInfo({ ...me, is_admin: admin })
          if (import.meta.env.DEV) {
            console.log('User info from API:', me, 'isAdmin:', admin)
          }
        } catch (error) {
          console.error('Failed to fetch user info:', error)
          setIsAdmin(false)
        }
      }
    }
    
    checkAdminStatus()
  }, [location])

  useEffect(() => {
    const fetchVersion = async () => {
      try {
        const version = await systemService.getVersion()
        setVersionInfo(version)
      } catch (error) {
        console.error('Failed to fetch version info:', error)
      }
    }
    fetchVersion()
  }, [])

  useEffect(() => {
    // Check for password expiration warning from sessionStorage
    const warningStr = sessionStorage.getItem('password_expiration_warning')
    if (warningStr) {
      try {
        const warning = JSON.parse(warningStr) as PasswordExpirationWarning
        if (warning.warning) {
          setPasswordWarning(warning)
        }
      } catch (error) {
        console.error('Failed to parse password expiration warning:', error)
      }
    }
  }, [])

  const handleDismissPasswordWarning = () => {
    setPasswordWarning(null)
    sessionStorage.removeItem('password_expiration_warning')
  }

  const handleLogout = async () => {
    const refreshToken = storage.getRefreshToken()
    if (refreshToken) {
      try {
        await authService.logout(refreshToken)
      } catch (error) {
        console.error('Logout error:', error)
      }
    }
    navigate('/login')
  }

  const handleLanguageChange = (lang: string) => {
    i18n.changeLanguage(lang)
    storage.setLanguage(lang)
  }

  const menuItems: MenuProps['items'] = useMemo(() => {
    // Always check fresh from storage to ensure we have latest state
    const userInfo = storage.getUserInfo()
    const currentIsAdmin = userInfo?.is_admin === true || userInfo?.is_admin === 'true'
    
    const items: MenuProps['items'] = [
      {
        key: '/',
        icon: <DashboardOutlined />,
        label: t('dashboard.title'),
      },
      {
        key: '/applications',
        icon: <AppstoreOutlined />,
        label: t('applications.title'),
      },
      {
        key: '/devices',
        icon: <MobileOutlined />,
        label: 'Devices',
      },
    ]

    // Admin-only menu items - always check fresh from storage
    if (currentIsAdmin) {
      items.push(
        {
          key: '/users',
          icon: <UserOutlined />,
          label: t('users.title'),
        },
        {
          key: '/automation',
          icon: <ThunderboltOutlined />,
          label: 'Automation',
        },
        {
          key: '/system',
          icon: <InfoCircleOutlined />,
          label: t('system.title'),
        }
      )
    }

    return items
  }, [isAdmin, t, location.pathname]) // Include location to force re-computation

  const userMenuItems: MenuProps['items'] = [
    {
      key: 'profile',
      icon: <ProfileOutlined />,
      label: t('settings.profile.title'),
      onClick: () => navigate('/settings?tab=profile'),
    },
    {
      key: 'changePassword',
      icon: <LockOutlined />,
      label: t('settings.security.changePassword'),
      onClick: () => setChangePasswordModalOpen(true),
    },
    {
      type: 'divider',
    },
    {
      key: 'language',
      icon: <GlobalOutlined />,
      label: 'Language',
      children: [
        {
          key: 'en',
          label: 'English',
          onClick: () => handleLanguageChange('en'),
        },
        {
          key: 'zh-CN',
          label: '中文',
          onClick: () => handleLanguageChange('zh-CN'),
        },
      ],
    },
    {
      type: 'divider',
    },
    {
      key: 'logout',
      icon: <LogoutOutlined />,
      label: t('common.logout'),
      onClick: handleLogout,
    },
  ]

  return (
    <Layout style={{ minHeight: '100vh' }}>
      <Sider
        trigger={null}
        collapsible
        collapsed={collapsed}
        style={{
          background: 'linear-gradient(180deg, #1e293b 0%, #0f172a 100%)',
        }}
      >
        <div
          style={{
            height: 64,
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            color: '#fff',
            fontSize: 18,
            fontWeight: 'bold',
            background: 'linear-gradient(135deg, #6366f1 0%, #8b5cf6 100%)',
            gap: 8,
          }}
        >
          <img
            src="/logo.svg"
            alt="OpenAuth Logo"
            style={{
              width: collapsed ? 32 : 40,
              height: collapsed ? 32 : 40,
            }}
          />
          {!collapsed && <span>OpenAuth</span>}
        </div>
        <Menu
          theme="dark"
          mode="inline"
          selectedKeys={[location.pathname]}
          items={menuItems}
          onClick={({ key }) => navigate(key)}
          style={{
            background: 'transparent',
            borderRight: 'none',
          }}
        />
      </Sider>
      <Layout>
        <Header
          style={{
            padding: '0 24px',
            background: '#fff',
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'space-between',
            boxShadow: '0 1px 4px rgba(0,21,41,.08)',
          }}
        >
          <Button
            type="text"
            icon={collapsed ? <MenuUnfoldOutlined /> : <MenuFoldOutlined />}
            onClick={() => setCollapsed(!collapsed)}
            style={{
              fontSize: 16,
              width: 64,
              height: 64,
            }}
          />
          <div style={{ display: 'flex', alignItems: 'center', gap: 16 }}>
            <a
              href="https://github.com/hanyouqing/OpenAuth"
              target="_blank"
              rel="noopener noreferrer"
              style={{
                fontSize: 20,
                color: '#1e293b',
                display: 'flex',
                alignItems: 'center',
                transition: 'color 0.2s',
              }}
              onMouseEnter={(e) => {
                e.currentTarget.style.color = '#6366f1'
              }}
              onMouseLeave={(e) => {
                e.currentTarget.style.color = '#1e293b'
              }}
            >
              <GithubOutlined />
            </a>
            <Dropdown menu={{ items: userMenuItems }} placement="bottomRight">
              <Avatar
                src="/logo.svg"
                style={{ cursor: 'pointer' }}
                alt="User Avatar"
              />
            </Dropdown>
          </div>
        </Header>
        <Content
          style={{
            margin: '24px 16px',
            padding: 24,
            minHeight: 280,
            background: '#f8fafc',
          }}
        >
          {passwordWarning && (
            <Alert
              message={
                passwordWarning.days_remaining > 0
                  ? `密码将在 ${passwordWarning.days_remaining} 天后过期`
                  : t('auth.passwordExpired')
              }
              description={
                passwordWarning.days_remaining > 0
                  ? `您的密码将在 ${passwordWarning.days_remaining} 天后过期，为了账户安全，请及时修改密码。`
                  : t('auth.passwordExpiredDescription')
              }
              type={passwordWarning.days_remaining > 0 ? 'warning' : 'error'}
              icon={<ExclamationCircleOutlined />}
              closable
              onClose={handleDismissPasswordWarning}
              showIcon
              action={
                <Button
                  size="small"
                  type="primary"
                  onClick={() => {
                    handleDismissPasswordWarning()
                    navigate('/settings?tab=security')
                  }}
                >
                  {t('settings.security.changePassword')}
                </Button>
              }
              style={{ marginBottom: 16 }}
            />
          )}
          {children}
        </Content>
        <ChangePasswordModal
          open={changePasswordModalOpen}
          onCancel={() => setChangePasswordModalOpen(false)}
        />
        <Footer
          style={{
            textAlign: 'center',
            padding: '12px 24px',
            background: '#fff',
            borderTop: '1px solid #e5e7eb',
            fontSize: '12px',
            color: '#64748b',
          }}
        >
          {versionInfo ? (
            <div
              style={{
                display: 'flex',
                justifyContent: 'center',
                alignItems: 'center',
                gap: 20,
                flexWrap: 'wrap',
              }}
            >
              <span>
                <strong style={{ color: '#1e293b' }}>Version:</strong>{' '}
                <code style={{ background: '#f1f5f9', padding: '2px 6px', borderRadius: 3 }}>
                  {versionInfo.version || 'N/A'}
                </code>
              </span>
              <span>
                <strong style={{ color: '#1e293b' }}>Build:</strong>{' '}
                {versionInfo.build_time || 'N/A'}
              </span>
              <span>
                <strong style={{ color: '#1e293b' }}>Commit:</strong>{' '}
                <code style={{ background: '#f1f5f9', padding: '2px 6px', borderRadius: 3 }}>
                  {versionInfo.git_commit
                    ? versionInfo.git_commit.substring(0, 7)
                    : 'N/A'}
                </code>
              </span>
              <span>
                <strong style={{ color: '#1e293b' }}>Go:</strong>{' '}
                {versionInfo.go_version || 'N/A'}
              </span>
            </div>
          ) : (
            <span>Loading version information...</span>
          )}
        </Footer>
      </Layout>
    </Layout>
  )
}

export default AppLayout
