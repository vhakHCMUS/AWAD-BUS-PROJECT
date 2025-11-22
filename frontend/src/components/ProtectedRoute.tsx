import { ReactNode } from 'react'
import { Navigate } from 'react-router-dom'

interface ProtectedRouteProps {
  children: ReactNode
  requireAdmin?: boolean
}

export default function ProtectedRoute({ children, requireAdmin = false }: ProtectedRouteProps) {
  const isAuthenticated = !!localStorage.getItem('access_token')
  
  if (!isAuthenticated) {
    return <Navigate to="/login" replace />
  }

  if (requireAdmin) {
    const userStr = localStorage.getItem('user')
    let userRole = 'passenger'
    
    if (userStr) {
      try {
        const user = JSON.parse(userStr)
        userRole = user.role || 'passenger'
      } catch (e) {
        console.error('Failed to parse user data:', e)
      }
    }
    
    if (userRole !== 'admin') {
      return <Navigate to="/" replace />
    }
  }

  return <>{children}</>
}
