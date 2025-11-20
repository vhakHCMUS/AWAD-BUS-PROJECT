package entities

import (
	"time"

	"github.com/google/uuid"
)

// SeatStatus represents the status of a seat
type SeatStatus string

const (
	SeatStatusAvailable SeatStatus = "available"
	SeatStatusLocked    SeatStatus = "locked"
	SeatStatusBooked    SeatStatus = "booked"
)

// SeatInfo represents seat availability for a specific trip
type SeatInfo struct {
	ID          uuid.UUID  `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	TripID      uuid.UUID  `json:"trip_id" gorm:"type:uuid;not null;index:idx_trip_seat,unique"`
	SeatNumber  string     `json:"seat_number" gorm:"not null;index:idx_trip_seat,unique"` // e.g., "A1", "B2"
	Status      SeatStatus `json:"status" gorm:"type:varchar(20);not null;default:'available';index"`
	LockedUntil *time.Time `json:"locked_until,omitempty" gorm:"index"`
	LockedBy    *uuid.UUID `json:"locked_by,omitempty" gorm:"type:uuid"` // User ID or session ID
	BookingID   *uuid.UUID `json:"booking_id,omitempty" gorm:"type:uuid;index"`
	CreatedAt   time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName overrides the table name
func (SeatInfo) TableName() string {
	return "seats_status"
}

// IsLockExpired checks if the seat lock has expired
func (s *SeatInfo) IsLockExpired() bool {
	if s.Status != SeatStatusLocked || s.LockedUntil == nil {
		return false
	}
	return time.Now().After(*s.LockedUntil)
}
