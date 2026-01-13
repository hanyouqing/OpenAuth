import { Routes, Route, Navigate } from 'react-router-dom'
import { ProtectedRoute } from './components/ProtectedRoute'
import { AdminRoute } from './components/AdminRoute'
import AppLayout from './components/layout/AppLayout'
import LoginPage from './pages/auth/LoginPage'
import DashboardPage from './pages/dashboard/DashboardPage'
import UsersPage from './pages/users/UsersPage'
import ApplicationsPage from './pages/applications/ApplicationsPage'
import SettingsPage from './pages/settings/SettingsPage'
import SystemInfoPage from './pages/system/SystemInfoPage'
import DevicesPage from './pages/devices/DevicesPage'
import AutomationPage from './pages/automation/AutomationPage'

function App() {
  return (
    <Routes>
      <Route path="/login" element={<LoginPage />} />
      <Route
        path="/*"
        element={
          <ProtectedRoute>
            <AppLayout>
              <Routes>
                <Route path="/" element={<DashboardPage />} />
                <Route path="/applications" element={<ApplicationsPage />} />
                <Route path="/settings" element={<SettingsPage />} />
                <Route path="/devices" element={<DevicesPage />} />
                <Route
                  path="/users"
                  element={
                    <AdminRoute>
                      <UsersPage />
                    </AdminRoute>
                  }
                />
                <Route
                  path="/automation"
                  element={
                    <AdminRoute>
                      <AutomationPage />
                    </AdminRoute>
                  }
                />
                <Route
                  path="/system"
                  element={
                    <AdminRoute>
                      <SystemInfoPage />
                    </AdminRoute>
                  }
                />
                <Route path="*" element={<Navigate to="/" replace />} />
              </Routes>
            </AppLayout>
          </ProtectedRoute>
        }
      />
    </Routes>
  )
}

export default App
