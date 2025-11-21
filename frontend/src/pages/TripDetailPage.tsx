import { useEffect, useState } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { tripsAPI, bookingsAPI, Trip, Seat as APISeat } from '../lib/api'
import { Loader2, ArrowLeft, MapPin, Clock, DollarSign } from 'lucide-react'

export default function TripDetailPage() {
  const { id } = useParams<{ id: string }>()
  const navigate = useNavigate()
  const [trip, setTrip] = useState<Trip | null>(null)
  const [seats, setSeats] = useState<APISeat[]>([])
  const [selectedSeats, setSelectedSeats] = useState<string[]>([])
  const [loading, setLoading] = useState(true)
  const [bookingLoading, setBookingLoading] = useState(false)
  const [contactInfo, setContactInfo] = useState({
    name: '',
    email: '',
    phone: '',
  })

  useEffect(() => {
    if (id) {
      loadTripDetails()
    }
  }, [id])

  const loadTripDetails = async () => {
    try {
      setLoading(true)
      const [tripRes, seatsRes] = await Promise.all([
        tripsAPI.getById(id!),
        tripsAPI.getSeats(id!),
      ])
      setTrip(tripRes.data)
      setSeats(seatsRes.data.seats)
    } catch (error) {
      console.error('Failed to load trip:', error)
      alert('Failed to load trip details')
      navigate('/search')
    } finally {
      setLoading(false)
    }
  }

  const handleSeatSelect = (seatNumber: string) => {
    setSelectedSeats((prev) =>
      prev.includes(seatNumber)
        ? prev.filter((s) => s !== seatNumber)
        : [...prev, seatNumber]
    )
  }

  const handleBooking = async (e: React.FormEvent) => {
    e.preventDefault()
    if (selectedSeats.length === 0) {
      alert('Please select at least one seat')
      return
    }

    try {
      setBookingLoading(true)
      const response = await bookingsAPI.initiate({
        trip_id: id!,
        seat_numbers: selectedSeats,
        contact_name: contactInfo.name,
        contact_email: contactInfo.email,
        contact_phone: contactInfo.phone,
      })

      // Redirect to payment page or booking confirmation
      navigate(`/bookings/${response.data.id}`)
    } catch (error: any) {
      alert(error.response?.data?.error || 'Failed to create booking')
    } finally {
      setBookingLoading(false)
    }
  }

  if (loading) {
    return (
      <div className="flex justify-center items-center h-screen">
        <Loader2 className="h-8 w-8 animate-spin text-primary-600" />
      </div>
    )
  }

  if (!trip) {
    return (
      <div className="container-custom py-8">
        <div className="card p-8 text-center">
          <h2 className="text-2xl font-bold mb-4">Trip not found</h2>
          <button onClick={() => navigate('/search')} className="btn bg-primary-600 text-white">
            Back to Search
          </button>
        </div>
      </div>
    )
  }

  const totalPrice = selectedSeats.length * trip.price

  return (
    <div className="container-custom py-8">
      <button
        onClick={() => navigate(-1)}
        className="flex items-center text-gray-600 hover:text-gray-900 mb-6"
      >
        <ArrowLeft className="h-5 w-5 mr-2" />
        Back
      </button>

      <div className="grid lg:grid-cols-3 gap-8">
        {/* Trip Details */}
        <div className="lg:col-span-2 space-y-6">
          <div className="card p-6">
            <h1 className="text-2xl font-bold mb-4">
              {trip.route?.from_city} → {trip.route?.to_city}
            </h1>
            <div className="grid md:grid-cols-2 gap-4 text-sm">
              <div className="flex items-start gap-3">
                <MapPin className="h-5 w-5 text-primary-600 flex-shrink-0 mt-0.5" />
                <div>
                  <div className="font-medium">Route</div>
                  <div className="text-gray-600">{trip.route?.name}</div>
                  <div className="text-gray-500">{trip.route?.distance} km</div>
                </div>
              </div>
              <div className="flex items-start gap-3">
                <Clock className="h-5 w-5 text-primary-600 flex-shrink-0 mt-0.5" />
                <div>
                  <div className="font-medium">Departure</div>
                  <div className="text-gray-600">
                    {new Date(trip.departure_time).toLocaleString('vi-VN')}
                  </div>
                </div>
              </div>
              <div className="flex items-start gap-3">
                <Clock className="h-5 w-5 text-primary-600 flex-shrink-0 mt-0.5" />
                <div>
                  <div className="font-medium">Arrival</div>
                  <div className="text-gray-600">
                    {new Date(trip.arrival_time).toLocaleString('vi-VN')}
                  </div>
                </div>
              </div>
              <div className="flex items-start gap-3">
                <DollarSign className="h-5 w-5 text-primary-600 flex-shrink-0 mt-0.5" />
                <div>
                  <div className="font-medium">Price per seat</div>
                  <div className="text-xl font-bold text-primary-600">
                    {trip.price.toLocaleString('vi-VN')} VND
                  </div>
                </div>
              </div>
            </div>
          </div>

          {/* Seat Selection */}
          <div className="card p-6">
            <h2 className="text-xl font-bold mb-4">Select Your Seats</h2>
            <div className="grid grid-cols-8 gap-2">
              {seats.map((seat) => (
                <button
                  key={seat.id}
                  onClick={() => handleSeatSelect(seat.seat_number)}
                  disabled={seat.status === 'booked' || seat.status === 'locked'}
                  className={`
                    p-4 rounded-lg border-2 font-medium text-sm transition-all
                    ${
                      selectedSeats.includes(seat.seat_number)
                        ? 'bg-primary-500 border-primary-600 text-white'
                        : seat.status === 'available'
                        ? 'bg-white border-gray-300 hover:border-primary-500 hover:bg-primary-50'
                        : seat.status === 'booked'
                        ? 'bg-gray-200 border-gray-300 text-gray-400 cursor-not-allowed'
                        : 'bg-yellow-100 border-yellow-400 text-yellow-700 cursor-not-allowed'
                    }
                  `}
                >
                  {seat.seat_number}
                </button>
              ))}
            </div>
            <div className="flex flex-wrap gap-4 mt-6 text-sm">
              <div className="flex items-center gap-2">
                <div className="w-8 h-10 bg-white border-2 border-gray-300 rounded"></div>
                <span>Available</span>
              </div>
              <div className="flex items-center gap-2">
                <div className="w-8 h-10 bg-primary-500 border-2 border-primary-600 rounded"></div>
                <span>Selected</span>
              </div>
              <div className="flex items-center gap-2">
                <div className="w-8 h-10 bg-gray-200 border-2 border-gray-300 rounded"></div>
                <span>Booked</span>
              </div>
              <div className="flex items-center gap-2">
                <div className="w-8 h-10 bg-yellow-100 border-2 border-yellow-400 rounded"></div>
                <span>Locked</span>
              </div>
            </div>
          </div>
        </div>

        {/* Booking Form */}
        <div className="lg:col-span-1">
          <div className="card p-6 sticky top-4">
            <h2 className="text-xl font-bold mb-4">Booking Details</h2>

            <div className="mb-4 p-4 bg-gray-50 rounded">
              <div className="flex justify-between mb-2">
                <span className="text-gray-600">Selected Seats:</span>
                <span className="font-medium">
                  {selectedSeats.length > 0 ? selectedSeats.join(', ') : 'None'}
                </span>
              </div>
              <div className="flex justify-between text-lg font-bold">
                <span>Total:</span>
                <span className="text-primary-600">{totalPrice.toLocaleString('vi-VN')} VND</span>
              </div>
            </div>

            <form onSubmit={handleBooking} className="space-y-4">
              <div>
                <label className="block text-sm font-medium mb-1">Full Name *</label>
                <input
                  type="text"
                  value={contactInfo.name}
                  onChange={(e) => setContactInfo({ ...contactInfo, name: e.target.value })}
                  className="input"
                  required
                />
              </div>
              <div>
                <label className="block text-sm font-medium mb-1">Email *</label>
                <input
                  type="email"
                  value={contactInfo.email}
                  onChange={(e) => setContactInfo({ ...contactInfo, email: e.target.value })}
                  className="input"
                  required
                />
              </div>
              <div>
                <label className="block text-sm font-medium mb-1">Phone Number *</label>
                <input
                  type="tel"
                  value={contactInfo.phone}
                  onChange={(e) => setContactInfo({ ...contactInfo, phone: e.target.value })}
                  className="input"
                  required
                />
              </div>
              <button
                type="submit"
                disabled={selectedSeats.length === 0 || bookingLoading}
                className="btn bg-primary-600 text-white hover:bg-primary-700 w-full disabled:opacity-50 disabled:cursor-not-allowed"
              >
                {bookingLoading ? (
                  <>
                    <Loader2 className="inline h-5 w-5 mr-2 animate-spin" />
                    Processing...
                  </>
                ) : (
                  'Continue to Payment'
                )}
              </button>
            </form>
          </div>
        </div>
      </div>
    </div>
  )
}
