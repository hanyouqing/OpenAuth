import { useState, useEffect, useCallback } from 'react'
import { useNavigate } from 'react-router-dom'

interface User {
  id: number
  username: string
  email: string
  status: string
}

export function useAuth() {
  const [user, setUser] = useState<User | null>(null)
  const [loading, setLoading] = useState(true)
  const navigate = useNavigate()

  const checkAuth = useCallback(() => {
    const token = localStorage.getItem('access_token')
    if (token) {
      // Token exists, user is authenticated
      setLoading(false)
      return true
    }
    setLoading(false)
    return false
  }, [])

  useEffect(() => {
    checkAuth()
  }, [checkAuth])

  const login = (userData: User, accessToken: string, refreshToken: string) => {
    localStorage.setItem('access_token', accessToken)
    localStorage.setItem('refresh_token', refreshToken)
    setUser(userData)
    setLoading(false)
  }

  const logout = () => {
    localStorage.removeItem('access_token')
    localStorage.removeItem('refresh_token')
    setUser(null)
    setLoading(false)
    navigate('/login')
  }

  const isAuthenticated = useCallback(() => {
    return !!localStorage.getItem('access_token')
  }, [])

  return {
    user,
    loading,
    login,
    logout,
    isAuthenticated,
  }
}
