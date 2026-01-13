import { useState, useEffect } from 'react'
import { Tabs } from 'antd'
import { UserOutlined, LockOutlined, SafetyOutlined } from '@ant-design/icons'
import { useTranslation } from 'react-i18next'
import { useSearchParams } from 'react-router-dom'
import MFASettingsPage from './MFASettingsPage'
import ProfileSettingsPage from './ProfileSettingsPage'
import SecuritySettingsPage from './SecuritySettingsPage'

function SettingsPage() {
  const { t } = useTranslation()
  const [searchParams, setSearchParams] = useSearchParams()
  const [activeTab, setActiveTab] = useState('profile')

  useEffect(() => {
    const tab = searchParams.get('tab')
    if (tab && ['profile', 'security', 'mfa'].includes(tab)) {
      setActiveTab(tab)
    }
  }, [searchParams])

  const handleTabChange = (key: string) => {
    setActiveTab(key)
    setSearchParams({ tab: key })
  }

  const items = [
    {
      key: 'profile',
      label: (
        <span>
          <UserOutlined />
          {t('settings.profile.title')}
        </span>
      ),
      children: <ProfileSettingsPage />,
    },
    {
      key: 'security',
      label: (
        <span>
          <LockOutlined />
          {t('settings.security.title')}
        </span>
      ),
      children: <SecuritySettingsPage />,
    },
    {
      key: 'mfa',
      label: (
        <span title={t('settings.mfa.title')}>
          <SafetyOutlined />
          MFA
        </span>
      ),
      children: <MFASettingsPage />,
    },
  ]

  return (
    <div>
      <h1 style={{ fontSize: '28px', fontWeight: 600, color: '#1e293b', marginBottom: 24 }}>
        {t('settings.title')}
      </h1>
      <Tabs items={items} activeKey={activeTab} onChange={handleTabChange} />
    </div>
  )
}

export default SettingsPage
