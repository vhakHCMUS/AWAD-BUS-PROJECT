import { useState, useEffect } from 'react';
import { RefreshCw } from 'lucide-react';
import SummaryCards from '../components/SummaryCards';
import ActivityList from '../components/ActivityList';
import StatisticChart from '../components/StatisticChart';
import { getSummaryCardsByRole, mockActivities, UserRole } from '../lib/mockData';

export default function DashboardPage() {
  const [role, setRole] = useState<UserRole>('admin');
  const [isLoading, setIsLoading] = useState(false);
  const [activities, setActivities] = useState(mockActivities);
  const [lastRefresh, setLastRefresh] = useState(new Date());

  // Get user role from localStorage (you can update this based on your auth system)
  useEffect(() => {
    const user = localStorage.getItem('user');
    if (user) {
      try {
        const userData = JSON.parse(user);
        setRole(userData.role || 'user');
      } catch {
        setRole('user');
      }
    }
  }, []);

  const summaryCards = getSummaryCardsByRole(role);

  const handleRefresh = async () => {
    setIsLoading(true);
    // Simulate API call
    await new Promise(resolve => setTimeout(resolve, 1500));

    // Shuffle activities to simulate new data
    const shuffled = [...activities].sort(() => Math.random() - 0.5);
    setActivities(shuffled);
    setLastRefresh(new Date());

    setIsLoading(false);
  };

  const handleRoleChange = (newRole: UserRole) => {
    setRole(newRole);
  };

  return (
    <div className="min-h-screen bg-gray-100 py-8 px-4 sm:px-6 lg:px-8">
      <div className="max-w-7xl mx-auto">
        {/* Header */}
        <div className="mb-8">
          <div className="flex flex-col md:flex-row md:items-center md:justify-between gap-4 mb-6">
            <div>
              <h1 className="text-3xl font-bold text-gray-900">Dashboard</h1>
              <p className="text-gray-600 mt-1">
                Welcome back! Last refreshed at {lastRefresh.toLocaleTimeString()}
              </p>
            </div>

            {/* Refresh Button */}
            <button
              onClick={handleRefresh}
              disabled={isLoading}
              className="flex items-center gap-2 px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors disabled:opacity-50 disabled:cursor-not-allowed w-fit"
            >
              <RefreshCw size={18} className={isLoading ? 'animate-spin' : ''} />
              {isLoading ? 'Refreshing...' : 'Refresh All'}
            </button>
          </div>

          {/* Role Selector */}
          <div className="flex gap-2 flex-wrap">
            <span className="text-sm text-gray-600 self-center">View as:</span>
            {(['admin', 'manager', 'user'] as const).map(r => (
              <button
                key={r}
                onClick={() => handleRoleChange(r)}
                className={`px-4 py-2 rounded-lg text-sm font-medium transition-colors ${
                  role === r
                    ? 'bg-blue-600 text-white'
                    : 'bg-white text-gray-700 border border-gray-300 hover:bg-gray-50'
                }`}
              >
                {r.charAt(0).toUpperCase() + r.slice(1)}
              </button>
            ))}
          </div>
        </div>

        {/* Summary Cards */}
        <SummaryCards cards={summaryCards} />

        {/* Main Content Grid */}
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
          {/* Activity List - takes 2 columns on large screens */}
          <div className="lg:col-span-2">
            <ActivityList activities={activities} onRefresh={handleRefresh} />
          </div>

          {/* Statistic Chart - takes 1 column on large screens */}
          <div>
            <StatisticChart role={role} />
          </div>
        </div>

        {/* Additional Info Section */}
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6 mt-8">
          {/* Quick Stats */}
          <div className="bg-white rounded-lg border border-gray-200 shadow-sm p-6">
            <h3 className="text-lg font-semibold text-gray-900 mb-4">Quick Stats</h3>
            <div className="space-y-3">
              <div className="flex justify-between items-center pb-3 border-b border-gray-200">
                <span className="text-gray-600">Total Activities</span>
                <span className="font-semibold text-gray-900">{activities.length}</span>
              </div>
              <div className="flex justify-between items-center pb-3 border-b border-gray-200">
                <span className="text-gray-600">Successful</span>
                <span className="font-semibold text-green-600">
                  {activities.filter(a => a.status === 'success').length}
                </span>
              </div>
              <div className="flex justify-between items-center pb-3 border-b border-gray-200">
                <span className="text-gray-600">Pending</span>
                <span className="font-semibold text-yellow-600">
                  {activities.filter(a => a.status === 'pending').length}
                </span>
              </div>
              <div className="flex justify-between items-center">
                <span className="text-gray-600">Failed</span>
                <span className="font-semibold text-red-600">
                  {activities.filter(a => a.status === 'failed').length}
                </span>
              </div>
            </div>
          </div>

          {/* User Guide */}
          <div className="bg-gradient-to-br from-blue-50 to-indigo-50 rounded-lg border border-blue-200 shadow-sm p-6">
            <h3 className="text-lg font-semibold text-gray-900 mb-4">Dashboard Features</h3>
            <ul className="space-y-2 text-sm text-gray-700">
              <li>✓ Real-time activity monitoring</li>
              <li>✓ Advanced filtering and sorting</li>
              <li>✓ Role-based data visualization</li>
              <li>✓ Interactive charts and statistics</li>
              <li>✓ System health monitoring</li>
              <li>✓ One-click data refresh</li>
            </ul>
          </div>
        </div>
      </div>
    </div>
  );
}
