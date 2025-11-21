import { useEffect, useState } from 'react'
import { adminTripsAPI, adminBusesAPI, adminRoutesAPI, tripsAPI, Trip, Bus, Route } from '../../lib/api'
import { Plus, Edit, Trash2, Loader2 } from 'lucide-react'

export default function AdminTrips() {
  const [trips, setTrips] = useState<Trip[]>([])
  const [buses, setBuses] = useState<Bus[]>([])
  const [routes, setRoutes] = useState<Route[]>([])
  const [loading, setLoading] = useState(true)
  const [showForm, setShowForm] = useState(false)
  const [editingTrip, setEditingTrip] = useState<Trip | null>(null)
  const [formData, setFormData] = useState({
    bus_id: '',
    route_id: '',
    departure_time: '',
    arrival_time: '',
    price: 0,
    driver_name: '',
    driver_phone: '',
  })

  useEffect(() => {
    loadData()
  }, [])

  const loadData = async () => {
    try {
      setLoading(true)
      const [busesRes, routesRes] = await Promise.all([
        adminBusesAPI.list(),
        adminRoutesAPI.list(),
      ])
      setBuses(busesRes.data.buses || [])
      setRoutes(routesRes.data.routes || [])
      
      // Load upcoming trips
      const tripsRes = await tripsAPI.search({
        from_city: '',
        to_city: '',
        date: new Date().toISOString().split('T')[0],
        limit: 100,
      })
      setTrips(tripsRes.data.trips || [])
    } catch (error) {
      console.error('Failed to load data:', error)
    } finally {
      setLoading(false)
    }
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    try {
      if (editingTrip) {
        await adminTripsAPI.update(editingTrip.id, formData)
        alert('Trip updated successfully')
      } else {
        await adminTripsAPI.create(formData)
        alert('Trip created successfully')
      }
      setShowForm(false)
      setEditingTrip(null)
      resetForm()
      loadData()
    } catch (error: any) {
      alert(error.response?.data?.error || 'Failed to save trip')
    }
  }

  const handleDelete = async (id: string) => {
    if (!confirm('Are you sure you want to delete this trip?')) return
    try {
      await adminTripsAPI.delete(id)
      alert('Trip deleted successfully')
      loadData()
    } catch (error: any) {
      alert(error.response?.data?.error || 'Failed to delete trip')
    }
  }

  const handleEdit = (trip: Trip) => {
    setEditingTrip(trip)
    setFormData({
      bus_id: trip.bus_id,
      route_id: trip.route_id,
      departure_time: new Date(trip.departure_time).toISOString().slice(0, 16),
      arrival_time: new Date(trip.arrival_time).toISOString().slice(0, 16),
      price: trip.price,
      driver_name: trip.driver_name || '',
      driver_phone: trip.driver_phone || '',
    })
    setShowForm(true)
  }

  const resetForm = () => {
    setFormData({
      bus_id: '',
      route_id: '',
      departure_time: '',
      arrival_time: '',
      price: 0,
      driver_name: '',
      driver_phone: '',
    })
  }

  if (loading) {
    return (
      <div className="flex justify-center items-center h-64">
        <Loader2 className="h-8 w-8 animate-spin text-primary-600" />
      </div>
    )
  }

  return (
    <div>
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-3xl font-bold">Manage Trips</h1>
        <button
          onClick={() => {
            resetForm()
            setEditingTrip(null)
            setShowForm(true)
          }}
          className="btn bg-primary-600 text-white hover:bg-primary-700"
        >
          <Plus className="h-5 w-5 mr-2" />
          Schedule Trip
        </button>
      </div>

      {showForm && (
        <div className="card p-6 mb-6">
          <h2 className="text-xl font-semibold mb-4">
            {editingTrip ? 'Edit Trip' : 'Schedule New Trip'}
          </h2>
          <form onSubmit={handleSubmit} className="space-y-4">
            <div className="grid md:grid-cols-2 gap-4">
              <div>
                <label className="block text-sm font-medium mb-1">Bus *</label>
                <select
                  value={formData.bus_id}
                  onChange={(e) => setFormData({ ...formData, bus_id: e.target.value })}
                  className="input"
                  required
                >
                  <option value="">Select a bus</option>
                  {buses.filter(b => b.status === 'active').map((bus) => (
                    <option key={bus.id} value={bus.id}>
                      {bus.license_plate} - {bus.bus_type} ({bus.seat_layout.total_seats} seats)
                    </option>
                  ))}
                </select>
              </div>
              <div>
                <label className="block text-sm font-medium mb-1">Route *</label>
                <select
                  value={formData.route_id}
                  onChange={(e) => {
                    const route = routes.find(r => r.id === e.target.value)
                    setFormData({
                      ...formData,
                      route_id: e.target.value,
                      price: route?.base_price || 0,
                    })
                  }}
                  className="input"
                  required
                >
                  <option value="">Select a route</option>
                  {routes.filter(r => r.is_active).map((route) => (
                    <option key={route.id} value={route.id}>
                      {route.name} ({route.distance} km)
                    </option>
                  ))}
                </select>
              </div>
              <div>
                <label className="block text-sm font-medium mb-1">Departure Time *</label>
                <input
                  type="datetime-local"
                  value={formData.departure_time}
                  onChange={(e) => setFormData({ ...formData, departure_time: e.target.value })}
                  className="input"
                  required
                />
              </div>
              <div>
                <label className="block text-sm font-medium mb-1">Arrival Time *</label>
                <input
                  type="datetime-local"
                  value={formData.arrival_time}
                  onChange={(e) => setFormData({ ...formData, arrival_time: e.target.value })}
                  className="input"
                  required
                />
              </div>
              <div>
                <label className="block text-sm font-medium mb-1">Price (VND) *</label>
                <input
                  type="number"
                  step="1000"
                  value={formData.price}
                  onChange={(e) => setFormData({ ...formData, price: parseFloat(e.target.value) })}
                  className="input"
                  required
                />
              </div>
              <div>
                <label className="block text-sm font-medium mb-1">Driver Name</label>
                <input
                  type="text"
                  value={formData.driver_name}
                  onChange={(e) => setFormData({ ...formData, driver_name: e.target.value })}
                  className="input"
                />
              </div>
              <div>
                <label className="block text-sm font-medium mb-1">Driver Phone</label>
                <input
                  type="tel"
                  value={formData.driver_phone}
                  onChange={(e) => setFormData({ ...formData, driver_phone: e.target.value })}
                  className="input"
                />
              </div>
            </div>
            <div className="flex gap-2">
              <button type="submit" className="btn bg-primary-600 text-white hover:bg-primary-700">
                {editingTrip ? 'Update' : 'Create'}
              </button>
              <button
                type="button"
                onClick={() => {
                  setShowForm(false)
                  setEditingTrip(null)
                  resetForm()
                }}
                className="btn bg-gray-200 hover:bg-gray-300"
              >
                Cancel
              </button>
            </div>
          </form>
        </div>
      )}

      <div className="card overflow-hidden">
        <table className="w-full">
          <thead className="bg-gray-50">
            <tr>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                Route
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                Bus
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                Departure
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                Price
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                Status
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                Actions
              </th>
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-200">
            {trips.map((trip) => (
              <tr key={trip.id}>
                <td className="px-6 py-4">
                  <div className="font-medium">{trip.route?.name || 'N/A'}</div>
                  <div className="text-sm text-gray-500">
                    {trip.route?.from_city} → {trip.route?.to_city}
                  </div>
                </td>
                <td className="px-6 py-4 whitespace-nowrap">
                  {trip.bus?.license_plate || 'N/A'}
                </td>
                <td className="px-6 py-4 whitespace-nowrap">
                  {new Date(trip.departure_time).toLocaleString('vi-VN')}
                </td>
                <td className="px-6 py-4 whitespace-nowrap">
                  {trip.price.toLocaleString('vi-VN')} VND
                </td>
                <td className="px-6 py-4 whitespace-nowrap">
                  <span
                    className={`px-2 py-1 text-xs rounded-full ${
                      trip.status === 'scheduled'
                        ? 'bg-blue-100 text-blue-800'
                        : trip.status === 'in_transit'
                        ? 'bg-green-100 text-green-800'
                        : trip.status === 'completed'
                        ? 'bg-gray-100 text-gray-800'
                        : 'bg-red-100 text-red-800'
                    }`}
                  >
                    {trip.status}
                  </span>
                </td>
                <td className="px-6 py-4 whitespace-nowrap">
                  <button
                    onClick={() => handleEdit(trip)}
                    className="text-blue-600 hover:text-blue-800 mr-3"
                  >
                    <Edit className="h-5 w-5" />
                  </button>
                  <button
                    onClick={() => handleDelete(trip.id)}
                    className="text-red-600 hover:text-red-800"
                  >
                    <Trash2 className="h-5 w-5" />
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
        {trips.length === 0 && (
          <div className="text-center py-8 text-gray-500">
            No trips scheduled. Create one to get started!
          </div>
        )}
      </div>
    </div>
  )
}
