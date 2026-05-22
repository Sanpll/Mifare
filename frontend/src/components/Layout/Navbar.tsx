import { Link, useNavigate } from 'react-router-dom'
import { useAuth } from '../../context/AuthContext'

const Navbar = () => {
  const { user, logout, isAuthenticated } = useAuth()
  const navigate = useNavigate()

  const handleLogout = () => {
    logout()
    navigate('/login')
  }

  return (
    <nav className="bg-white shadow-md">
      <div className="max-w-7xl mx-auto px-4">
        <div className="flex justify-between items-center h-16">
          <div className="flex items-center space-x-8">
            <Link to="/" className="text-xl font-bold text-blue-600">Mifare</Link>
            
            {isAuthenticated && (
              <div className="flex space-x-6">
                <Link to="/cards" className="hover:text-blue-600">Cards</Link>
                <Link to="/terminal/auth" className="hover:text-blue-600">Terminal Auth</Link>
                {user?.isAdmin && (
                  <>
                    <Link to="/keys" className="hover:text-blue-600">Keys</Link>
                    <Link to="/terminals" className="hover:text-blue-600">Terminals</Link>
                    <Link to="/transactions" className="hover:text-blue-600">Transactions</Link>
                    <Link to="/users" className="hover:text-blue-600">Users</Link>
                  </>
                )}
              </div>
            )}
          </div>

          <div>
            {isAuthenticated ? (
              <div className="flex items-center gap-4">
                <span className="text-sm text-gray-600">
                  {user?.username} {user?.isAdmin && '(Admin)'}
                </span>
                <button
                  onClick={handleLogout}
                  className="text-red-600 hover:text-red-700 text-sm font-medium"
                >
                  Logout
                </button>
              </div>
            ) : (
              <Link to="/login" className="text-blue-600 hover:text-blue-700">Login</Link>
            )}
          </div>
        </div>
      </div>
    </nav>
  )
}

export default Navbar