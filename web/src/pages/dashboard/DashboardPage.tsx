import { Card, Row, Col, Statistic } from 'antd'
import { UserOutlined, AppstoreOutlined, TeamOutlined } from '@ant-design/icons'
import { useTranslation } from 'react-i18next'

function DashboardPage() {
  const { t } = useTranslation()

  return (
    <div>
      <h1 style={{ marginBottom: 24, fontSize: '28px', fontWeight: 600, color: '#1e293b' }}>
        {t('dashboard.title')}
      </h1>
      <Row gutter={[16, 16]}>
        <Col xs={24} sm={12} lg={8}>
          <Card
            hoverable
            style={{
              background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
              border: 'none',
            }}
            styles={{ body: { color: '#fff' } }}
          >
            <Statistic
              title={<span style={{ color: '#fff', opacity: 0.9 }}>{t('dashboard.totalUsers')}</span>}
              value={0}
              valueStyle={{ color: '#fff' }}
              prefix={<UserOutlined style={{ color: '#fff', opacity: 0.8 }} />}
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} lg={8}>
          <Card
            hoverable
            style={{
              background: 'linear-gradient(135deg, #f093fb 0%, #f5576c 100%)',
              border: 'none',
            }}
            styles={{ body: { color: '#fff' } }}
          >
            <Statistic
              title={<span style={{ color: '#fff', opacity: 0.9 }}>{t('dashboard.totalApplications')}</span>}
              value={0}
              valueStyle={{ color: '#fff' }}
              prefix={<AppstoreOutlined style={{ color: '#fff', opacity: 0.8 }} />}
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} lg={8}>
          <Card
            hoverable
            style={{
              background: 'linear-gradient(135deg, #4facfe 0%, #00f2fe 100%)',
              border: 'none',
            }}
            styles={{ body: { color: '#fff' } }}
          >
            <Statistic
              title={<span style={{ color: '#fff', opacity: 0.9 }}>{t('dashboard.activeSessions')}</span>}
              value={0}
              valueStyle={{ color: '#fff' }}
              prefix={<TeamOutlined style={{ color: '#fff', opacity: 0.8 }} />}
            />
          </Card>
        </Col>
      </Row>
    </div>
  )
}

export default DashboardPage
