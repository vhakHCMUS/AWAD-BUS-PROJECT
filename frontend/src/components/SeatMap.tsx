import { useState } from 'react'
import { useTranslation } from 'react-i18next'
import { motion } from 'framer-motion'
import { cn } from '../lib/utils'

export interface Seat {
  id: string
  number: string
  status: 'available' | 'selected' | 'booked' | 'locked'
  floor?: 1 | 2
  row: number
  col: number
}

export interface BusLayout {
  type: '29_limousine' | '45_sleeper' | '40_seater' | '16_limousine'
  floors: number
  rows: number
  cols: number
  seats: Seat[]
}

interface SeatMapProps {
  layout: BusLayout
  selectedSeats: string[]
  onSeatSelect: (seatNumber: string) => void
  maxSeats?: number
}

export const BUS_LAYOUTS: Record<string, Partial<BusLayout>> = {
  '29_limousine': {
    floors: 2,
    rows: 5,
    cols: 3,
  },
  '45_sleeper': {
    floors: 2,
    rows: 12,
    cols: 2,
  },
  '40_seater': {
    floors: 1,
    rows: 10,
    cols: 4,
  },
  '16_limousine': {
    floors: 1,
    rows: 4,
    cols: 4,
  },
}

export default function SeatMap({ layout, selectedSeats, onSeatSelect, maxSeats = 6 }: SeatMapProps) {
  const { t } = useTranslation()
  const [activeFloor, setActiveFloor] = useState<1 | 2>(1)

  const handleSeatClick = (seat: Seat) => {
    if (seat.status === 'booked' || seat.status === 'locked') return

    if (seat.status === 'selected') {
      onSeatSelect(seat.number)
      return
    }

    if (selectedSeats.length >= maxSeats) {
      return
    }

    onSeatSelect(seat.number)
  }

  const getSeatStatus = (seat: Seat): Seat['status'] => {
    if (selectedSeats.includes(seat.number)) return 'selected'
    return seat.status
  }

  const renderSeat = (seat: Seat) => {
    const status = getSeatStatus(seat)
    const isClickable = status === 'available' || status === 'selected'

    return (
      <motion.button
        key={seat.id}
        whileHover={isClickable ? { scale: 1.05 } : {}}
        whileTap={isClickable ? { scale: 0.95 } : {}}
        onClick={() => handleSeatClick(seat)}
        disabled={!isClickable}
        className={cn(
          'relative w-12 h-16 rounded-lg border-2 transition-all duration-200 font-medium text-sm',
          {
            'bg-white border-gray-300 hover:border-primary-500 hover:bg-primary-50 cursor-pointer':
              status === 'available',
            'bg-primary-500 border-primary-600 text-white shadow-md':
              status === 'selected',
            'bg-gray-200 border-gray-300 text-gray-400 cursor-not-allowed':
              status === 'booked',
            'bg-yellow-100 border-yellow-400 text-yellow-700 cursor-not-allowed animate-pulse':
              status === 'locked',
          }
        )}
      >
        <div className="flex flex-col items-center justify-center h-full">
          <span className="text-xs">{seat.number}</span>
        </div>
      </motion.button>
    )
  }

  const renderFloor = (floorNumber: 1 | 2) => {
    const floorSeats = layout.seats.filter(
      (s) => (layout.floors === 1 ? true : s.floor === floorNumber)
    )

    const seatGrid = Array.from({ length: layout.rows }, (_, rowIndex) => {
      return floorSeats.filter((s) => s.row === rowIndex)
    })

    return (
      <div className="space-y-2">
        {/* Driver seat indicator */}
        {floorNumber === 1 && (
          <div className="flex justify-end mb-4">
            <div className="w-12 h-12 bg-gray-700 rounded-lg flex items-center justify-center text-white text-xs">
              {t('seat.driver')}
            </div>
          </div>
        )}

        {/* Seat grid */}
        <div className="space-y-2">
          {seatGrid.map((row, rowIndex) => (
            <div key={rowIndex} className="flex justify-center gap-2">
              {row.map((seat) => renderSeat(seat))}
            </div>
          ))}
        </div>
      </div>
    )
  }

  return (
    <div className="bg-white rounded-xl shadow-lg p-6">
      {/* Floor selector for multi-floor buses */}
      {layout.floors === 2 && (
        <div className="flex gap-2 mb-6">
          <button
            onClick={() => setActiveFloor(1)}
            className={cn(
              'flex-1 py-2 px-4 rounded-lg font-medium transition-colors',
              activeFloor === 1
                ? 'bg-primary-500 text-white'
                : 'bg-gray-100 text-gray-600 hover:bg-gray-200'
            )}
          >
            {t('seat.floor_1')}
          </button>
          <button
            onClick={() => setActiveFloor(2)}
            className={cn(
              'flex-1 py-2 px-4 rounded-lg font-medium transition-colors',
              activeFloor === 2
                ? 'bg-primary-500 text-white'
                : 'bg-gray-100 text-gray-600 hover:bg-gray-200'
            )}
          >
            {t('seat.floor_2')}
          </button>
        </div>
      )}

      {/* Seat map */}
      <div className="bg-gray-50 rounded-lg p-4 border-2 border-gray-200">
        {layout.floors === 1 ? (
          renderFloor(1)
        ) : (
          <div>{renderFloor(activeFloor)}</div>
        )}
      </div>

      {/* Legend */}
      <div className="flex flex-wrap gap-4 mt-6 text-sm">
        <div className="flex items-center gap-2">
          <div className="w-8 h-10 bg-white border-2 border-gray-300 rounded"></div>
          <span>{t('seat.available')}</span>
        </div>
        <div className="flex items-center gap-2">
          <div className="w-8 h-10 bg-primary-500 border-2 border-primary-600 rounded"></div>
          <span>{t('seat.selected')}</span>
        </div>
        <div className="flex items-center gap-2">
          <div className="w-8 h-10 bg-gray-200 border-2 border-gray-300 rounded"></div>
          <span>{t('seat.booked')}</span>
        </div>
        <div className="flex items-center gap-2">
          <div className="w-8 h-10 bg-yellow-100 border-2 border-yellow-400 rounded"></div>
          <span>{t('seat.locked')}</span>
        </div>
      </div>

      {/* Selected seats info */}
      {selectedSeats.length > 0 && (
        <div className="mt-4 p-4 bg-primary-50 rounded-lg border border-primary-200">
          <p className="text-sm font-medium text-primary-900">
            {t('booking.seats_selected', { count: selectedSeats.length })}:{' '}
            <span className="font-bold">{selectedSeats.join(', ')}</span>
          </p>
        </div>
      )}
    </div>
  )
}
