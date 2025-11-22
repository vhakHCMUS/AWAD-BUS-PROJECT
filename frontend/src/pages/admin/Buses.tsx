import { useEffect, useState } from 'react'
import { adminBusesAPI, Bus } from '../../lib/api'
import { Plus, Edit, Trash2, Loader2 } from 'lucide-react'
import toast from 'react-hot-toast'

export default function AdminBuses() {
  const [buses, setBuses] = useState<Bus[]>([])
  const [loading, setLoading] = useState(true)
  const [showForm, setShowForm] = useState(false)
  const [editingBus, setEditingBus] = useState<Bus | null>(null)
  const [formData, setFormData] = useState<Partial<Bus>>({
    license_plate: '',
    bus_type: 'seater',
    manufacturer: '',
    model: '',
    year: new Date().getFullYear(),
    operator_name: '',
    status: 'active',
    amenities: [],
    seat_layout: {
      rows: 10,
      columns: 4,
      total_seats: 40,
      floors: 1,
      layout: [],
    },
  })

  useEffect(() => {
    loadBuses()
  }, [])

  const loadBuses = async () => {
    try {
      setLoading(true)
      const response = await adminBusesAPI.list()
      setBuses(response.data.buses || [])
    } catch (error) {
      console.error('Failed to load buses:', error)
      toast.error('Failed to load buses')
    } finally {
      setLoading(false)
    }
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    try {
      if (editingBus) {
        await adminBusesAPI.update(editingBus.id, formData)
        toast.success('Bus updated successfully!')
      } else {
        await adminBusesAPI.create(formData)
        toast.success('Bus created successfully!')
      }
      setShowForm(false)
      setEditingBus(null)
      resetForm()
      loadBuses()
    } catch (error: any) {
      toast.error(error.response?.data?.error || 'Failed to save bus')
    }
  }

  const handleDelete = async (id: string) => {
    toast((t) => (
      <div className="flex flex-col gap-3">
        <p className="font-medium">Are you sure you want to delete this bus?</p>
        <div className="flex gap-2 justify-end">
          <button
            onClick={() => {
              toast.dismiss(t.id)
            }}
            className="px-3 py-1.5 text-sm bg-gray-200 hover:bg-gray-300 rounded-md transition-colors"
          >
            Cancel
          </button>
          <button
            onClick={async () => {
              toast.dismiss(t.id)
              try {
                await adminBusesAPI.delete(id)
                toast.success('Bus deleted successfully!')
                loadBuses()
              } catch (error: any) {
                toast.error(error.response?.data?.error || 'Failed to delete bus')
              }
            }}
            className="px-3 py-1.5 text-sm bg-red-600 text-white hover:bg-red-700 rounded-md transition-colors"
          >
            Delete
          </button>
        </div>
      </div>
    ), {
      duration: 10000,
    })
  }

  const handleEdit = (bus: Bus) => {
    setEditingBus(bus)
    setFormData(bus)
    setShowForm(true)
  }

  const resetForm = () => {
    setFormData({
      license_plate: '',
      bus_type: 'seater',
      manufacturer: '',
      model: '',
      year: new Date().getFullYear(),
      operator_name: '',
      status: 'active',
      amenities: [],
      seat_layout: {
        rows: 10,
        columns: 4,
        total_seats: 40,
        floors: 1,
        layout: [],
      },
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
        <h1 className="text-3xl font-bold">Manage Buses</h1>
        <button
          onClick={() => {
            resetForm()
            setEditingBus(null)
            setShowForm(true)
          }}
          className="btn bg-primary-600 text-white hover:bg-primary-700 flex items-center px-4 py-2 font-semibold rounded-lg shadow-md hover:shadow-lg transition-all duration-200"
        >
          <Plus className="h-5 w-5 mr-2" />
          Add Bus
        </button>
      </div>

      {showForm && (
        <div className="card p-6 mb-6">
          <h2 className="text-xl font-semibold mb-4">
            {editingBus ? 'Edit Bus' : 'Add New Bus'}
          </h2>
          <form onSubmit={handleSubmit} className="space-y-4">
            <div className="grid md:grid-cols-2 gap-4">
              <div>
                <label className="block text-sm font-medium mb-1">License Plate *</label>
                <input
                  type="text"
                  value={formData.license_plate}
                  onChange={(e) => setFormData({ ...formData, license_plate: e.target.value })}
                  className="input"
                  required
                />
              </div>
              <div>
                <label className="block text-sm font-medium mb-1">Bus Type *</label>
                <select
                  value={formData.bus_type}
                  onChange={(e) => setFormData({ ...formData, bus_type: e.target.value })}
                  className="input"
                  required
                >
                  <option value="seater">Seater</option>
                  <option value="sleeper">Sleeper</option>
                  <option value="semi-sleeper">Semi-Sleeper</option>
                  <option value="limousine">Limousine</option>
                </select>
              </div>
              <div>
                <label className="block text-sm font-medium mb-1">Manufacturer</label>
                <input
                  type="text"
                  value={formData.manufacturer}
                  onChange={(e) => setFormData({ ...formData, manufacturer: e.target.value })}
                  className="input"
                />
              </div>
              <div>
                <label className="block text-sm font-medium mb-1">Model</label>
                <input
                  type="text"
                  value={formData.model}
                  onChange={(e) => setFormData({ ...formData, model: e.target.value })}
                  className="input"
                />
              </div>
              <div>
                <label className="block text-sm font-medium mb-1">Year</label>
                <input
                  type="number"
                  value={formData.year}
                  onChange={(e) => setFormData({ ...formData, year: parseInt(e.target.value) })}
                  className="input"
                />
              </div>
              <div>
                <label className="block text-sm font-medium mb-1">Operator Name</label>
                <input
                  type="text"
                  value={formData.operator_name}
                  onChange={(e) => setFormData({ ...formData, operator_name: e.target.value })}
                  className="input"
                />
              </div>
              <div>
                <label className="block text-sm font-medium mb-1">Status</label>
                <select
                  value={formData.status}
                  onChange={(e) => setFormData({ ...formData, status: e.target.value })}
                  className="input"
                >
                  <option value="active">Active</option>
                  <option value="maintenance">Maintenance</option>
                  <option value="inactive">Inactive</option>
                </select>
              </div>
              <div>
                <label className="block text-sm font-medium mb-1">Total Seats</label>
                <input
                  type="number"
                  value={formData.seat_layout?.total_seats}
                  onChange={(e) =>
                    setFormData({
                      ...formData,
                      seat_layout: {
                        ...formData.seat_layout!,
                        total_seats: parseInt(e.target.value),
                      },
                    })
                  }
                  className="input"
                />
              </div>
            </div>
            <div className="flex gap-3">
              <button 
                type="submit" 
                className="btn bg-primary-600 text-white hover:bg-primary-700 px-6 py-2.5 font-semibold rounded-lg shadow-md hover:shadow-lg transition-all duration-200"
              >
                {editingBus ? 'Update' : 'Create'}
              </button>
              <button
                type="button"
                onClick={() => {
                  setShowForm(false)
                  setEditingBus(null)
                  resetForm()
                }}
                className="btn bg-gray-200 hover:bg-gray-300 px-6 py-2.5 font-semibold rounded-lg shadow-sm hover:shadow-md transition-all duration-200"
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
                License Plate
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                Type
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                Model
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                Seats
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
            {buses.map((bus) => (
              <tr key={bus.id}>
                <td className="px-6 py-4 whitespace-nowrap font-medium">{bus.license_plate}</td>
                <td className="px-6 py-4 whitespace-nowrap">{bus.bus_type}</td>
                <td className="px-6 py-4 whitespace-nowrap">
                  {bus.manufacturer} {bus.model}
                </td>
                <td className="px-6 py-4 whitespace-nowrap">{bus.seat_layout.total_seats}</td>
                <td className="px-6 py-4 whitespace-nowrap">
                  <span
                    className={`px-2 py-1 text-xs rounded-full ${
                      bus.status === 'active'
                        ? 'bg-green-100 text-green-800'
                        : bus.status === 'maintenance'
                        ? 'bg-yellow-100 text-yellow-800'
                        : 'bg-gray-100 text-gray-800'
                    }`}
                  >
                    {bus.status}
                  </span>
                </td>
                <td className="px-6 py-4 whitespace-nowrap">
                  <button
                    onClick={() => handleEdit(bus)}
                    className="text-blue-600 hover:text-blue-800 mr-3"
                  >
                    <Edit className="h-5 w-5" />
                  </button>
                  <button
                    onClick={() => handleDelete(bus.id)}
                    className="text-red-600 hover:text-red-800"
                  >
                    <Trash2 className="h-5 w-5" />
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
        {buses.length === 0 && (
          <div className="text-center py-8 text-gray-500">No buses found. Add one to get started!</div>
        )}
      </div>
    </div>
  )
}
