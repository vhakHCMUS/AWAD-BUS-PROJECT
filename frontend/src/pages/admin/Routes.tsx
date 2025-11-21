import { useEffect, useState } from 'react'
import { adminRoutesAPI, Route } from '../../lib/api'
import { Plus, Edit, Trash2, Loader2 } from 'lucide-react'

export default function AdminRoutes() {
  const [routes, setRoutes] = useState<Route[]>([])
  const [loading, setLoading] = useState(true)
  const [showForm, setShowForm] = useState(false)
  const [editingRoute, setEditingRoute] = useState<Route | null>(null)
  const [formData, setFormData] = useState<Partial<Route>>({
    name: '',
    from_city: '',
    to_city: '',
    distance: 0,
    base_price: 0,
    description: '',
    is_active: true,
  })

  useEffect(() => {
    loadRoutes()
  }, [])

  const loadRoutes = async () => {
    try {
      setLoading(true)
      const response = await adminRoutesAPI.list()
      setRoutes(response.data.routes || [])
    } catch (error) {
      console.error('Failed to load routes:', error)
      alert('Failed to load routes')
    } finally {
      setLoading(false)
    }
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    try {
      if (editingRoute) {
        await adminRoutesAPI.update(editingRoute.id, formData)
        alert('Route updated successfully')
      } else {
        await adminRoutesAPI.create(formData)
        alert('Route created successfully')
      }
      setShowForm(false)
      setEditingRoute(null)
      resetForm()
      loadRoutes()
    } catch (error: any) {
      alert(error.response?.data?.error || 'Failed to save route')
    }
  }

  const handleDelete = async (id: string) => {
    if (!confirm('Are you sure you want to delete this route?')) return
    try {
      await adminRoutesAPI.delete(id)
      alert('Route deleted successfully')
      loadRoutes()
    } catch (error: any) {
      alert(error.response?.data?.error || 'Failed to delete route')
    }
  }

  const handleEdit = (route: Route) => {
    setEditingRoute(route)
    setFormData(route)
    setShowForm(true)
  }

  const resetForm = () => {
    setFormData({
      name: '',
      from_city: '',
      to_city: '',
      distance: 0,
      base_price: 0,
      description: '',
      is_active: true,
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
        <h1 className="text-3xl font-bold">Manage Routes</h1>
        <button
          onClick={() => {
            resetForm()
            setEditingRoute(null)
            setShowForm(true)
          }}
          className="btn bg-primary-600 text-white hover:bg-primary-700"
        >
          <Plus className="h-5 w-5 mr-2" />
          Add Route
        </button>
      </div>

      {showForm && (
        <div className="card p-6 mb-6">
          <h2 className="text-xl font-semibold mb-4">
            {editingRoute ? 'Edit Route' : 'Add New Route'}
          </h2>
          <form onSubmit={handleSubmit} className="space-y-4">
            <div className="grid md:grid-cols-2 gap-4">
              <div>
                <label className="block text-sm font-medium mb-1">Route Name</label>
                <input
                  type="text"
                  value={formData.name}
                  onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                  className="input"
                  placeholder="e.g., Hanoi - Ho Chi Minh"
                />
              </div>
              <div>
                <label className="block text-sm font-medium mb-1">From City *</label>
                <input
                  type="text"
                  value={formData.from_city}
                  onChange={(e) => setFormData({ ...formData, from_city: e.target.value })}
                  className="input"
                  required
                />
              </div>
              <div>
                <label className="block text-sm font-medium mb-1">To City *</label>
                <input
                  type="text"
                  value={formData.to_city}
                  onChange={(e) => setFormData({ ...formData, to_city: e.target.value })}
                  className="input"
                  required
                />
              </div>
              <div>
                <label className="block text-sm font-medium mb-1">Distance (km) *</label>
                <input
                  type="number"
                  step="0.1"
                  value={formData.distance}
                  onChange={(e) => setFormData({ ...formData, distance: parseFloat(e.target.value) })}
                  className="input"
                  required
                />
              </div>
              <div>
                <label className="block text-sm font-medium mb-1">Base Price (VND) *</label>
                <input
                  type="number"
                  step="1000"
                  value={formData.base_price}
                  onChange={(e) => setFormData({ ...formData, base_price: parseFloat(e.target.value) })}
                  className="input"
                  required
                />
              </div>
              <div>
                <label className="block text-sm font-medium mb-1">Status</label>
                <select
                  value={formData.is_active ? 'true' : 'false'}
                  onChange={(e) => setFormData({ ...formData, is_active: e.target.value === 'true' })}
                  className="input"
                >
                  <option value="true">Active</option>
                  <option value="false">Inactive</option>
                </select>
              </div>
            </div>
            <div>
              <label className="block text-sm font-medium mb-1">Description</label>
              <textarea
                value={formData.description}
                onChange={(e) => setFormData({ ...formData, description: e.target.value })}
                className="input"
                rows={3}
              />
            </div>
            <div className="flex gap-2">
              <button type="submit" className="btn bg-primary-600 text-white hover:bg-primary-700">
                {editingRoute ? 'Update' : 'Create'}
              </button>
              <button
                type="button"
                onClick={() => {
                  setShowForm(false)
                  setEditingRoute(null)
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
                Route Name
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                From → To
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                Distance
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                Base Price
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
            {routes.map((route) => (
              <tr key={route.id}>
                <td className="px-6 py-4 whitespace-nowrap font-medium">{route.name}</td>
                <td className="px-6 py-4 whitespace-nowrap">
                  {route.from_city} → {route.to_city}
                </td>
                <td className="px-6 py-4 whitespace-nowrap">{route.distance} km</td>
                <td className="px-6 py-4 whitespace-nowrap">
                  {route.base_price.toLocaleString('vi-VN')} VND
                </td>
                <td className="px-6 py-4 whitespace-nowrap">
                  <span
                    className={`px-2 py-1 text-xs rounded-full ${
                      route.is_active
                        ? 'bg-green-100 text-green-800'
                        : 'bg-gray-100 text-gray-800'
                    }`}
                  >
                    {route.is_active ? 'Active' : 'Inactive'}
                  </span>
                </td>
                <td className="px-6 py-4 whitespace-nowrap">
                  <button
                    onClick={() => handleEdit(route)}
                    className="text-blue-600 hover:text-blue-800 mr-3"
                  >
                    <Edit className="h-5 w-5" />
                  </button>
                  <button
                    onClick={() => handleDelete(route.id)}
                    className="text-red-600 hover:text-red-800"
                  >
                    <Trash2 className="h-5 w-5" />
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
        {routes.length === 0 && (
          <div className="text-center py-8 text-gray-500">No routes found. Add one to get started!</div>
        )}
      </div>
    </div>
  )
}
