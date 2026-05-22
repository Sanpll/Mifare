import { Routes, Route, Navigate } from 'react-router-dom'
import { useAuth } from './context/AuthContext'
import Login from './pages/auth/Login'
import Register from './pages/auth/Register'
import Dashboard from './pages/dashboard/Dashboard'
import TerminalAuth from './pages/terminal-auth/TerminalAuth'

const PrivateRoute = ({ children, adminOnly = false }: { 
  children: React.ReactNode
  adminOnly?: boolean 
}) => {
  const { isAuthenticated, user } = useAuth()
  
  if (!isAuthenticated) {
    return <Navigate to="/login" />
  }
  
  if (adminOnly && !user?.isAdmin) {
    return <Navigate to="/" />
  }
  
  return <>{children}</>
}

const AppRouter = () => {
  return (
    <Routes>
      <Route path="/login" element={<Login />} />
      <Route path="/register" element={<Register />} />
      
      <Route path="/" element={
        <PrivateRoute>
          <Dashboard />
        </PrivateRoute>
      } />
      
      <Route path="/terminal/auth" element={
        <PrivateRoute>
          <TerminalAuth />
        </PrivateRoute>
      } />

      {/* Здесь позже добавим другие страницы */}
      <Route path="*" element={<Navigate to="/" />} />
    </Routes>
  )
}

export default AppRouter