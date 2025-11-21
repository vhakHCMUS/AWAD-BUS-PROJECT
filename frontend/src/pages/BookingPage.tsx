import { useEffect, useState } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { bookingsAPI, paymentsAPI, Booking } from '../lib/api'
import { Loader2, CheckCircle, Clock, AlertCircle } from 'lucide-react'
import { PaymentSelector } from '../components/PaymentSelector'

export default function BookingPage() {
  const { id } = useParams<{ id: string }>()
  const navigate = useNavigate()
  const [booking, setBooking] = useState<Booking | null>(null)
  const [loading, setLoading] = useState(true)
  const [paymentLoading, setPaymentLoading] = useState(false)

  useEffect(() => {
    if (id) {
      loadBooking()
    }
  }, [id])

  const loadBooking = async () => {
    try {
      setLoading(true)
      const response = await bookingsAPI.getById(id!)
      setBooking(response.data)
    } catch (error) {
      console.error('Failed to load booking:', error)
      alert('Booking not found')
      navigate('/search')
    } finally {
      setLoading(false)
    }
  }

  const handlePayment = async (gateway: string) => {
    try {
      setPaymentLoading(true)
      const response = await paymentsAPI.create({
        booking_id: id!,
        gateway: gateway,
      })

      // In a real app, redirect to payment gateway URL
      console.log('Payment created:', response.data)
      alert('Payment gateway integration would redirect here')
      
      // For demo, just refresh booking status
      await loadBooking()
    } catch (error: any) {
      alert(error.response?.data?.error || 'Failed to create payment')
    } finally {
      setPaymentLoading(false)
    }
  }

  if (loading) {
    return (
      <div className="flex justify-center items-center h-screen">
        <Loader2 className="h-8 w-8 animate-spin text-primary-600" />
      </div>
    )
  }

  if (!booking) {
    return (
      <div className="container-custom py-8">
        <div className="card p-8 text-center">
          <AlertCircle className="h-16 w-16 text-red-500 mx-auto mb-4" />
          <h2 className="text-2xl font-bold mb-4">Booking not found</h2>
          <button onClick={() => navigate('/search')} className="btn bg-primary-600 text-white">
            Back to Search
          </button>
        </div>
      </div>
    )
  }

  const isExpired = booking.expires_at && new Date(booking.expires_at) < new Date()

  return (
    <div className="container-custom py-8 max-w-4xl">
      <div className="space-y-6">
        {/* Booking Status */}
        <div className="card p-6">
          <div className="flex items-center justify-between mb-4">
            <h1 className="text-2xl font-bold">Booking Details</h1>
            <div
              className={`px-4 py-2 rounded-full text-sm font-medium ${
                booking.status === 'confirmed'
                  ? 'bg-green-100 text-green-800'
                  : booking.status === 'pending'
                  ? 'bg-yellow-100 text-yellow-800'
                  : 'bg-red-100 text-red-800'
              }`}
            >
              {booking.status.toUpperCase()}
            </div>
          </div>

          <div className="grid md:grid-cols-2 gap-6">
            <div>
              <p className="text-sm text-gray-500 mb-1">Booking Code</p>
              <p className="text-lg font-mono font-bold">{booking.booking_code}</p>
            </div>
            <div>
              <p className="text-sm text-gray-500 mb-1">Total Amount</p>
              <p className="text-2xl font-bold text-primary-600">
                {booking.total_price.toLocaleString('vi-VN')} VND
              </p>
            </div>
            <div>
              <p className="text-sm text-gray-500 mb-1">Selected Seats</p>
              <p className="font-medium">{booking.seats.join(', ')}</p>
            </div>
            <div>
              <p className="text-sm text-gray-500 mb-1">Passenger Count</p>
              <p className="font-medium">{booking.seats.length} seat(s)</p>
            </div>
          </div>

          {booking.expires_at && booking.status === 'pending' && (
            <div className="mt-4 p-4 bg-yellow-50 border border-yellow-200 rounded-lg flex items-start gap-3">
              <Clock className="h-5 w-5 text-yellow-600 flex-shrink-0 mt-0.5" />
              <div>
                <p className="font-medium text-yellow-900">
                  {isExpired ? 'Booking Expired' : 'Complete payment before:'}
                </p>
                <p className="text-sm text-yellow-700">
                  {new Date(booking.expires_at).toLocaleString('vi-VN')}
                </p>
              </div>
            </div>
          )}
        </div>

        {/* Trip Details */}
        {booking.trip && (
          <div className="card p-6">
            <h2 className="text-xl font-bold mb-4">Trip Information</h2>
            <div className="grid md:grid-cols-2 gap-4">
              <div>
                <p className="text-sm text-gray-500">Route</p>
                <p className="font-medium">{booking.trip.route?.name}</p>
              </div>
              <div>
                <p className="text-sm text-gray-500">Bus</p>
                <p className="font-medium">{booking.trip.bus?.license_plate}</p>
              </div>
              <div>
                <p className="text-sm text-gray-500">Departure</p>
                <p className="font-medium">
                  {new Date(booking.trip.departure_time).toLocaleString('vi-VN')}
                </p>
              </div>
              <div>
                <p className="text-sm text-gray-500">Arrival</p>
                <p className="font-medium">
                  {new Date(booking.trip.arrival_time).toLocaleString('vi-VN')}
                </p>
              </div>
            </div>
          </div>
        )}

        {/* Contact Information */}
        <div className="card p-6">
          <h2 className="text-xl font-bold mb-4">Contact Information</h2>
          <div className="grid md:grid-cols-3 gap-4">
            <div>
              <p className="text-sm text-gray-500">Name</p>
              <p className="font-medium">{booking.contact_name}</p>
            </div>
            <div>
              <p className="text-sm text-gray-500">Email</p>
              <p className="font-medium">{booking.contact_email}</p>
            </div>
            <div>
              <p className="text-sm text-gray-500">Phone</p>
              <p className="font-medium">{booking.contact_phone}</p>
            </div>
          </div>
        </div>

        {/* Payment Section */}
        {booking.status === 'pending' && !isExpired && (
          <div className="card p-6">
            <h2 className="text-xl font-bold mb-4">Complete Payment</h2>
            <p className="text-gray-600 mb-4">
              Total Amount: <span className="font-bold text-primary-600">{booking.total_price.toLocaleString('vi-VN')} VND</span>
            </p>
            <PaymentSelector onSelect={handlePayment} />
            {paymentLoading && (
              <div className="mt-4 text-center">
                <Loader2 className="inline h-5 w-5 animate-spin text-primary-600" />
                <span className="ml-2">Processing payment...</span>
              </div>
            )}
          </div>
        )}

        {booking.status === 'confirmed' && (
          <div className="card p-6 bg-green-50 border-green-200">
            <div className="flex items-center gap-4">
              <CheckCircle className="h-12 w-12 text-green-600" />
              <div>
                <h3 className="text-xl font-bold text-green-900">Payment Successful!</h3>
                <p className="text-green-700">
                  Your booking is confirmed. Check your email for the ticket details.
                </p>
              </div>
            </div>
            <div className="mt-4 flex gap-3">
              <button
                onClick={() => navigate('/bookings')}
                className="btn bg-primary-600 text-white hover:bg-primary-700"
              >
                View My Bookings
              </button>
              <button
                onClick={() => navigate('/search')}
                className="btn bg-gray-200 hover:bg-gray-300"
              >
                Book Another Trip
              </button>
            </div>
          </div>
        )}
      </div>
    </div>
  )
}
