import { Navigate } from 'react-router-dom'
import { storage } from '@/utils/storage'

interface AdminRouteProps {
  children: React.ReactNode
}

export function AdminRoute({ children }: AdminRouteProps) {
  const userInfo = storage.getUserInfo()
  const isAdmin = userInfo?.is_admin === true

  if (!isAdmin) {
    return <Navigate to="/" replace />
  }

  return <>{children}</>
}
