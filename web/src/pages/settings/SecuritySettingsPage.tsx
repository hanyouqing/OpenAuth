import { Alert } from 'antd'
import { LockOutlined } from '@ant-design/icons'
import { useTranslation } from 'react-i18next'

function SecuritySettingsPage() {
  const { t } = useTranslation()

  return (
    <Alert
      message={t('settings.security.passwordChangeMoved')}
      description={t('settings.security.passwordChangeMovedDescription')}
      type="info"
      icon={<LockOutlined />}
      showIcon
    />
  )
}

export default SecuritySettingsPage
