import { useState, useEffect } from 'react'
import {
  Card,
  Button,
  List,
  Modal,
  Input,
  QRCode,
  message,
  Space,
  Typography,
  Switch,
  Divider,
} from 'antd'
import { PlusOutlined, DeleteOutlined } from '@ant-design/icons'
import { useTranslation } from 'react-i18next'
import { mfaService, MFADevice } from '@/services/mfaService'
import { mfaPolicyService } from '@/services/mfaPolicyService'

const { Title, Text } = Typography

function MFASettingsPage() {
  const { t } = useTranslation()
  const [devices, setDevices] = useState<MFADevice[]>([])
  const [loading, setLoading] = useState(false)
  const [forceMFA, setForceMFA] = useState(false)
  const [policyLoading, setPolicyLoading] = useState(false)
  const [totpModalVisible, setTotpModalVisible] = useState(false)
  const [qrCode, setQrCode] = useState('')
  const [totpSecret, setTotpSecret] = useState('')
  const [verifyToken, setVerifyToken] = useState('')
  const [deviceName, setDeviceName] = useState('')

  useEffect(() => {
    fetchDevices()
    fetchMFAPolicy()
  }, [])

  const fetchMFAPolicy = async () => {
    setPolicyLoading(true)
    try {
      const policy = await mfaPolicyService.getPolicy()
      setForceMFA(policy.force_mfa)
    } catch (error: any) {
      message.error(
        error.response?.data?.error || 'Failed to load MFA policy'
      )
    } finally {
      setPolicyLoading(false)
    }
  }

  const handleForceMFAToggle = async (checked: boolean) => {
    try {
      await mfaPolicyService.updatePolicy(checked)
      setForceMFA(checked)
      message.success(
        checked
          ? t('settings.mfa.forceMFAEnabled')
          : t('settings.mfa.forceMFADisabled')
      )
    } catch (error: any) {
      message.error(
        error.response?.data?.error || 'Failed to update MFA policy'
      )
    }
  }

  const fetchDevices = async () => {
    setLoading(true)
    try {
      const response = await mfaService.listDevices()
      setDevices(response.data)
    } catch (error: any) {
      message.error(
        error.response?.data?.error || 'Failed to load MFA devices'
      )
    } finally {
      setLoading(false)
    }
  }

  const handleGenerateTOTP = async () => {
    try {
      const response = await mfaService.generateTOTP(deviceName || 'My Device')
      setTotpSecret(response.secret)
      setQrCode(response.qr_code)
      setTotpModalVisible(true)
    } catch (error: any) {
      message.error(
        error.response?.data?.error || 'Failed to generate TOTP'
      )
    }
  }

  const handleVerifyTOTP = async () => {
    try {
      await mfaService.verifyTOTP(verifyToken)
      message.success('TOTP verified successfully')
      setTotpModalVisible(false)
      setVerifyToken('')
      setDeviceName('')
      fetchDevices()
    } catch (error: any) {
      message.error(error.response?.data?.error || 'Invalid token')
    }
  }

  const handleDeleteDevice = async (deviceId: number) => {
    try {
      await mfaService.deleteDevice(deviceId)
      message.success('Device deleted successfully')
      fetchDevices()
    } catch (error: any) {
      message.error(
        error.response?.data?.error || 'Failed to delete device'
      )
    }
  }

  return (
    <div>
      <Card>
        <Space direction="vertical" style={{ width: '100%' }} size="large">
          <div
            style={{
              display: 'flex',
              justifyContent: 'space-between',
              alignItems: 'center',
            }}
          >
            <Title level={3}>{t('settings.mfa.title')}</Title>
            <Button
              type="primary"
              icon={<PlusOutlined />}
              onClick={() => setTotpModalVisible(true)}
            >
              {t('settings.mfa.addDevice')}
            </Button>
          </div>

          <Divider />

          <div
            style={{
              display: 'flex',
              justifyContent: 'space-between',
              alignItems: 'center',
            }}
          >
            <div>
              <Text strong>{t('settings.mfa.forceMFA')}</Text>
              <br />
              <Text type="secondary">{t('settings.mfa.forceMFADescription')}</Text>
            </div>
            <Switch
              checked={forceMFA}
              onChange={handleForceMFAToggle}
              loading={policyLoading}
            />
          </div>

          <Divider />

          <Title level={4}>{t('settings.mfa.devices')}</Title>

          <List
            loading={loading}
            dataSource={devices}
            renderItem={(device) => (
              <List.Item
                actions={[
                  <Button
                    danger
                    icon={<DeleteOutlined />}
                    onClick={() => handleDeleteDevice(device.id)}
                  >
                    {t('common.delete')}
                  </Button>,
                ]}
              >
                <List.Item.Meta
                  title={device.name}
                  description={`${t('settings.mfa.type')}: ${device.type.toUpperCase()} | ${t('settings.mfa.verified')}: ${
                    device.verified ? t('common.yes') : t('common.no')
                  }`}
                />
              </List.Item>
            )}
          />
        </Space>
      </Card>

      <Modal
        title={t('settings.mfa.setupTOTP')}
        open={totpModalVisible}
        onCancel={() => {
          setTotpModalVisible(false)
          setQrCode('')
          setTotpSecret('')
          setVerifyToken('')
          setDeviceName('')
        }}
        footer={null}
      >
        {!qrCode ? (
          <Space direction="vertical" style={{ width: '100%' }}>
            <Input
              placeholder={t('settings.mfa.deviceName')}
              value={deviceName}
              onChange={(e) => setDeviceName(e.target.value)}
            />
            <Button type="primary" onClick={handleGenerateTOTP} block>
              {t('settings.mfa.generateQRCode')}
            </Button>
          </Space>
        ) : (
          <Space direction="vertical" style={{ width: '100%' }} size="large">
            <div style={{ textAlign: 'center' }}>
              <Text>{t('settings.mfa.scanQRCode')}</Text>
              <QRCode value={qrCode} style={{ marginTop: 16 }} />
              <Text type="secondary" style={{ display: 'block', marginTop: 8 }}>
                {t('settings.mfa.enterSecret')}: {totpSecret}
              </Text>
            </div>
            <Input
              placeholder={t('settings.mfa.enter6DigitCode')}
              value={verifyToken}
              onChange={(e) => setVerifyToken(e.target.value)}
              maxLength={6}
            />
            <Button type="primary" onClick={handleVerifyTOTP} block>
              {t('settings.mfa.verifyAndEnable')}
            </Button>
          </Space>
        )}
      </Modal>
    </div>
  )
}

export default MFASettingsPage
