import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { tripsAPI, Trip } from '../lib/api'
import { Search, MapPin, Calendar, ArrowRight, Loader2, Bus } from 'lucide-react'

export default function TripSearchPage() {
  const navigate = useNavigate()
  const [searchParams, setSearchParams] = useState({
    from_city: '',
    to_city: '',
    date: new Date().toISOString().split('T')[0],
  })
  const [trips, setTrips] = useState<Trip[]>([])
  const [loading, setLoading] = useState(false)
  const [searched, setSearched] = useState(false)

  const handleSearch = async (e: React.FormEvent) => {
    e.preventDefault()
    try {
      setLoading(true)
      const response = await tripsAPI.search(searchParams)
      setTrips(response.data.trips || [])
      setSearched(true)
    } catch (error) {
      console.error('Search failed:', error)
      alert('Failed to search trips')
    } finally {
      setLoading(false)
    }
  }

  const calculateDuration = (departure: string, arrival: string) => {
    const diff = new Date(arrival).getTime() - new Date(departure).getTime()
    const hours = Math.floor(diff / (1000 * 60 * 60))
    const minutes = Math.floor((diff % (1000 * 60 * 60)) / (1000 * 60))
    return `${hours}h ${minutes}m`
  }

  return (
    <div className="container-custom py-8">
      <h1 className="text-3xl font-bold mb-8">Search Bus Trips</h1>

      {/* Search Form */}
      <div className="card p-6 mb-8">
        <form onSubmit={handleSearch} className="space-y-4">
          <div className="grid md:grid-cols-3 gap-4">
            <div>
              <label className="block text-sm font-medium mb-2">
                <MapPin className="inline h-4 w-4 mr-1" />
                From City
              </label>
              <input
                type="text"
                value={searchParams.from_city}
                onChange={(e) => setSearchParams({ ...searchParams, from_city: e.target.value })}
                className="input"
                placeholder="e.g., Hanoi"
                required
              />
            </div>
            <div>
              <label className="block text-sm font-medium mb-2">
                <MapPin className="inline h-4 w-4 mr-1" />
                To City
              </label>
              <input
                type="text"
                value={searchParams.to_city}
                onChange={(e) => setSearchParams({ ...searchParams, to_city: e.target.value })}
                className="input"
                placeholder="e.g., Ho Chi Minh"
                required
              />
            </div>
            <div>
              <label className="block text-sm font-medium mb-2">
                <Calendar className="inline h-4 w-4 mr-1" />
                Travel Date
              </label>
              <input
                type="date"
                value={searchParams.date}
                onChange={(e) => setSearchParams({ ...searchParams, date: e.target.value })}
                className="input"
                required
                min={new Date().toISOString().split('T')[0]}
              />
            </div>
          </div>
          <button
            type="submit"
            disabled={loading}
            className="btn bg-primary-600 text-white hover:bg-primary-700 w-full md:w-auto disabled:opacity-50"
          >
            {loading ? (
              <>
                <Loader2 className="inline h-5 w-5 mr-2 animate-spin" />
                Searching...
              </>
            ) : (
              <>
                <Search className="inline h-5 w-5 mr-2" />
                Search Trips
              </>
            )}
          </button>
        </form>
      </div>

      {/* Results */}
      {loading && (
        <div className="flex justify-center items-center h-64">
          <Loader2 className="h-8 w-8 animate-spin text-primary-600" />
        </div>
      )}

      {!loading && searched && trips.length === 0 && (
        <div className="card p-12 text-center">
          <Bus className="h-16 w-16 text-gray-300 mx-auto mb-4" />
          <h3 className="text-xl font-semibold text-gray-700 mb-2">No trips found</h3>
          <p className="text-gray-500">
            Try adjusting your search criteria or selecting a different date.
          </p>
        </div>
      )}

      {!loading && trips.length > 0 && (
        <div className="space-y-4">
          <h2 className="text-xl font-semibold mb-4">
            {trips.length} trip{trips.length !== 1 ? 's' : ''} found
          </h2>
          {trips.map((trip) => (
            <div key={trip.id} className="card p-6 hover:shadow-lg transition-shadow">
              <div className="flex flex-col md:flex-row md:items-center md:justify-between gap-4">
                <div className="flex-1">
                  <div className="flex items-center gap-4 mb-3">
                    <div>
                      <div className="text-2xl font-bold">
                        {new Date(trip.departure_time).toLocaleTimeString('vi-VN', {
                          hour: '2-digit',
                          minute: '2-digit',
                        })}
                      </div>
                      <div className="text-sm text-gray-500">{trip.route?.from_city}</div>
                    </div>
                    <div className="flex-1 flex flex-col items-center">
                      <div className="text-xs text-gray-500 mb-1">
                        {calculateDuration(trip.departure_time, trip.arrival_time)}
                      </div>
                      <div className="w-full h-px bg-gray-300 relative">
                        <ArrowRight className="absolute -top-2.5 right-0 h-5 w-5 text-primary-600" />
                      </div>
                    </div>
                    <div>
                      <div className="text-2xl font-bold">
                        {new Date(trip.arrival_time).toLocaleTimeString('vi-VN', {
                          hour: '2-digit',
                          minute: '2-digit',
                        })}
                      </div>
                      <div className="text-sm text-gray-500">{trip.route?.to_city}</div>
                    </div>
                  </div>
                  <div className="flex flex-wrap gap-4 text-sm text-gray-600">
                    <div>
                      <span className="font-medium">Bus:</span> {trip.bus?.license_plate || 'N/A'}
                    </div>
                    <div>
                      <span className="font-medium">Type:</span> {trip.bus?.bus_type || 'N/A'}
                    </div>
                    <div>
                      <span className="font-medium">Seats Available:</span>{' '}
                      {trip.bus?.seat_layout.total_seats || 'N/A'}
                    </div>
                  </div>
                </div>
                <div className="flex flex-col items-end gap-2">
                  <div className="text-right">
                    <div className="text-2xl font-bold text-primary-600">
                      {trip.price.toLocaleString('vi-VN')} VND
                    </div>
                    <div className="text-sm text-gray-500">per seat</div>
                  </div>
                  <button
                    onClick={() => navigate(`/trips/${trip.id}`)}
                    className="btn bg-primary-600 text-white hover:bg-primary-700"
                  >
                    Select Seats
                  </button>
                </div>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  )
}
