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

  // TODO: Implement proper role checking from JWT token
  if (requireAdmin) {
    const userRole = localStorage.getItem('user_role') || 'passenger'
    if (userRole !== 'admin') {
      return <Navigate to="/" replace />
    }
  }

  return <>{children}</>
}
