import { Link } from 'react-router-dom'
import { Bus, LogOut } from 'lucide-react'

export default function Navbar() {
  const isAuthenticated = !!localStorage.getItem('access_token')

  const handleLogout = () => {
    localStorage.removeItem('access_token')
    localStorage.removeItem('refresh_token')
    window.location.href = '/login'
  }

  return (
    <nav className="bg-white shadow-sm border-b">
      <div className="container-custom">
        <div className="flex h-16 items-center justify-between">
          <Link to="/" className="flex items-center gap-2 font-bold text-xl text-primary-600">
            <Bus className="h-6 w-6" />
            Bus Booking
          </Link>

          <div className="flex items-center gap-6">
            <Link to="/search" className="text-gray-700 hover:text-primary-600">
              Search Trips
            </Link>

            {isAuthenticated ? (
              <>
                <Link to="/my-bookings" className="text-gray-700 hover:text-primary-600">
                  My Bookings
                </Link>
                <button
                  onClick={handleLogout}
                  className="flex items-center gap-2 text-gray-700 hover:text-primary-600"
                >
                  <LogOut className="h-4 w-4" />
                  Logout
                </button>
              </>
            ) : (
              <>
                <Link to="/login" className="text-gray-700 hover:text-primary-600">
                  Login
                </Link>
                <Link to="/register" className="btn-primary px-4 py-2 text-sm">
                  Register
                </Link>
              </>
            )}
          </div>
        </div>
      </div>
    </nav>
  )
}
