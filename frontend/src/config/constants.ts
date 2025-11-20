// Environment variables
export const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1'
export const WS_URL = import.meta.env.VITE_WS_URL || 'ws://localhost:8081/ws'

// App constants
export const APP_NAME = 'Bus Booking'
export const DEFAULT_PAGE_SIZE = 20

// Seat statuses
export const SEAT_STATUS = {
  AVAILABLE: 'available',
  LOCKED: 'locked',
  BOOKED: 'booked',
} as const

// Booking statuses
export const BOOKING_STATUS = {
  PENDING: 'pending',
  PAID: 'paid',
  CONFIRMED: 'confirmed',
  EXPIRED: 'expired',
  CANCELLED: 'cancelled',
  REFUNDED: 'refunded',
} as const

// Trip statuses
export const TRIP_STATUS = {
  SCHEDULED: 'scheduled',
  BOARDING: 'boarding',
  IN_TRANSIT: 'in_transit',
  COMPLETED: 'completed',
  CANCELLED: 'cancelled',
  DELAYED: 'delayed',
} as const

// User roles
export const USER_ROLE = {
  GUEST: 'guest',
  PASSENGER: 'passenger',
  ADMIN: 'admin',
} as const
