import { Routes, Route } from 'react-router-dom'
import { Toaster } from 'react-hot-toast'
import Layout from './components/Layout'
import { ChatbotWidget } from './components/ChatbotWidget'
import HomePage from './pages/HomePage'
import TripSearchPage from './pages/TripSearchPage'
import TripDetailPage from './pages/TripDetailPage'
import BookingPage from './pages/BookingPage'
import MyBookingsPage from './pages/MyBookingsPage'
import LoginPage from './pages/LoginPage'
import RegisterPage from './pages/RegisterPage'
import AdminLayout from './components/admin/AdminLayout'
import AdminDashboard from './pages/admin/Dashboard'
import AdminBuses from './pages/admin/Buses'
import AdminRoutes from './pages/admin/Routes'
import AdminTrips from './pages/admin/Trips'
import ProtectedRoute from './components/ProtectedRoute'
import NotFoundPage from './pages/NotFoundPage'

function App() {
  return (
    <>
      <Toaster position="top-right" />
      <ChatbotWidget />
      <Routes>
      {/* Public routes */}
      <Route path="/" element={<Layout />}>
        <Route index element={<HomePage />} />
        <Route path="search" element={<TripSearchPage />} />
        <Route path="trips/:id" element={<TripDetailPage />} />
        <Route path="login" element={<LoginPage />} />
        <Route path="register" element={<RegisterPage />} />
        <Route path="bookings/:id" element={<BookingPage />} />
        
        {/* Protected routes */}
        <Route
          path="booking/:tripId"
          element={
            <ProtectedRoute>
              <BookingPage />
            </ProtectedRoute>
          }
        />
        <Route
          path="my-bookings"
          element={
            <ProtectedRoute>
              <MyBookingsPage />
            </ProtectedRoute>
          }
        />
      </Route>

      {/* Admin routes */}
      <Route
        path="/admin"
        element={
          <ProtectedRoute requireAdmin>
            <AdminLayout />
          </ProtectedRoute>
        }
      >
        <Route index element={<AdminDashboard />} />
        <Route path="buses" element={<AdminBuses />} />
        <Route path="routes" element={<AdminRoutes />} />
        <Route path="trips" element={<AdminTrips />} />
      </Route>

      {/* 404 */}
      <Route path="*" element={<NotFoundPage />} />
    </Routes>
    </>
  )
}

export default App
