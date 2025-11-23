// Mock data for dashboard
export type UserRole = 'admin' | 'manager' | 'user';

export interface SummaryCard {
  id: string;
  title: string;
  value: number | string;
  unit?: string;
  change?: number;
  icon: string;
  color: 'blue' | 'green' | 'orange' | 'red';
}

export interface Activity {
  id: string;
  type: 'booking' | 'payment' | 'cancellation' | 'registration';
  description: string;
  user: string;
  timestamp: Date;
  status: 'success' | 'pending' | 'failed';
}

export interface ChartData {
  date: string;
  value: number;
  label?: string;
}

// Admin Dashboard Mock Data
export const adminSummaryCards: SummaryCard[] = [
  {
    id: 'total-bookings',
    title: 'Total Bookings',
    value: 1234,
    change: 12.5,
    icon: '📅',
    color: 'blue',
  },
  {
    id: 'revenue',
    title: 'Revenue',
    value: '$48,500',
    unit: 'USD',
    change: 8.2,
    icon: '💰',
    color: 'green',
  },
  {
    id: 'active-users',
    title: 'Active Users',
    value: 856,
    change: 5.1,
    icon: '👥',
    color: 'orange',
  },
  {
    id: 'system-status',
    title: 'System Status',
    value: 'Healthy',
    icon: '✅',
    color: 'green',
  },
];

// Manager Dashboard Mock Data
export const managerSummaryCards: SummaryCard[] = [
  {
    id: 'trips-today',
    title: 'Trips Today',
    value: 45,
    change: 3.2,
    icon: '🚌',
    color: 'blue',
  },
  {
    id: 'daily-revenue',
    title: 'Daily Revenue',
    value: '$2,450',
    unit: 'USD',
    change: -2.5,
    icon: '💵',
    color: 'green',
  },
  {
    id: 'occupancy',
    title: 'Occupancy Rate',
    value: '78%',
    change: 4.8,
    icon: '📊',
    color: 'orange',
  },
];

// User Dashboard Mock Data
export const userSummaryCards: SummaryCard[] = [
  {
    id: 'total-bookings',
    title: 'My Bookings',
    value: 12,
    icon: '🎫',
    color: 'blue',
  },
  {
    id: 'upcoming-trips',
    title: 'Upcoming Trips',
    value: 2,
    icon: '🛣️',
    color: 'green',
  },
  {
    id: 'total-spent',
    title: 'Total Spent',
    value: '$350',
    unit: 'USD',
    icon: '💳',
    color: 'orange',
  },
];

// Mock Activity Data
export const mockActivities: Activity[] = [
  {
    id: '1',
    type: 'booking',
    description: 'New booking from Ha Noi to Ho Chi Minh',
    user: 'John Doe',
    timestamp: new Date(Date.now() - 5 * 60000),
    status: 'success',
  },
  {
    id: '2',
    type: 'payment',
    description: 'Payment of $150 received',
    user: 'Jane Smith',
    timestamp: new Date(Date.now() - 15 * 60000),
    status: 'success',
  },
  {
    id: '3',
    type: 'cancellation',
    description: 'Booking cancelled - refund processed',
    user: 'Mike Johnson',
    timestamp: new Date(Date.now() - 30 * 60000),
    status: 'success',
  },
  {
    id: '4',
    type: 'registration',
    description: 'New user registered',
    user: 'Sarah Williams',
    timestamp: new Date(Date.now() - 45 * 60000),
    status: 'success',
  },
  {
    id: '5',
    type: 'booking',
    description: 'Booking attempt failed - payment error',
    user: 'Tom Brown',
    timestamp: new Date(Date.now() - 60 * 60000),
    status: 'failed',
  },
  {
    id: '6',
    type: 'payment',
    description: 'Payment verification pending',
    user: 'Emily Davis',
    timestamp: new Date(Date.now() - 90 * 60000),
    status: 'pending',
  },
];

// Mock Chart Data - Revenue Trend
export const revenueChartData: ChartData[] = [
  { date: 'Jan 01', value: 4000 },
  { date: 'Jan 02', value: 3000 },
  { date: 'Jan 03', value: 2000 },
  { date: 'Jan 04', value: 2780 },
  { date: 'Jan 05', value: 1890 },
  { date: 'Jan 06', value: 2390 },
  { date: 'Jan 07', value: 3490 },
];

// Mock Chart Data - Booking Status
export const bookingStatusData = [
  { name: 'Completed', value: 400, color: '#10b981' },
  { name: 'Pending', value: 300, color: '#f59e0b' },
  { name: 'Cancelled', value: 200, color: '#ef4444' },
  { name: 'No Show', value: 50, color: '#8b5cf6' },
];

// Mock Chart Data - Trip Performance
export const tripPerformanceData: ChartData[] = [
  { date: 'Mon', value: 45 },
  { date: 'Tue', value: 52 },
  { date: 'Wed', value: 48 },
  { date: 'Thu', value: 61 },
  { date: 'Fri', value: 55 },
  { date: 'Sat', value: 67 },
  { date: 'Sun', value: 42 },
];

// Mock System Status Data
export const systemStatusData = [
  { name: 'API Server', status: 'online' as const, uptime: 99.99 },
  { name: 'Database', status: 'online' as const, uptime: 99.95 },
  { name: 'Cache Server', status: 'online' as const, uptime: 99.87 },
  { name: 'Payment Gateway', status: 'online' as const, uptime: 99.92 },
];

// Get summary cards based on role
export const getSummaryCardsByRole = (role: UserRole): SummaryCard[] => {
  switch (role) {
    case 'admin':
      return adminSummaryCards;
    case 'manager':
      return managerSummaryCards;
    case 'user':
      return userSummaryCards;
    default:
      return userSummaryCards;
  }
};
