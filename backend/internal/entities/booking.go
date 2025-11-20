package entities

import (
	"time"

	"github.com/google/uuid"
)

// BookingStatus represents the status of a booking
type BookingStatus string

const (
	BookingStatusPending   BookingStatus = "pending"
	BookingStatusPaid      BookingStatus = "paid"
	BookingStatusConfirmed BookingStatus = "confirmed"
	BookingStatusExpired   BookingStatus = "expired"
	BookingStatusCancelled BookingStatus = "cancelled"
	BookingStatusRefunded  BookingStatus = "refunded"
)

// PassengerInfo represents passenger details for a booking
type PassengerInfo struct {
	Name   string `json:"name"`
	Phone  string `json:"phone"`
	Email  string `json:"email"`
	IDCard string `json:"id_card,omitempty"`
}

// Booking represents a ticket booking
type Booking struct {
	ID           uuid.UUID     `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	TripID       uuid.UUID     `json:"trip_id" gorm:"type:uuid;not null;index"`
	UserID       *uuid.UUID    `json:"user_id,omitempty" gorm:"type:uuid;index"` // Nullable for guest bookings
	ContactEmail string        `json:"contact_email" gorm:"not null"`
	ContactPhone string        `json:"contact_phone" gorm:"not null"`
	ContactName  string        `json:"contact_name" gorm:"not null"`
	Seats        []string      `json:"seats" gorm:"type:text[];not null"` // ["A1", "A2"]
	TotalPrice   float64       `json:"total_price" gorm:"not null"`
	Status       BookingStatus `json:"status" gorm:"type:varchar(20);not null;default:'pending';index"`
	ExpiresAt    *time.Time    `json:"expires_at,omitempty" gorm:"index"`        // For pending bookings
	BookingCode  string        `json:"booking_code" gorm:"uniqueIndex;not null"` // Human-readable code

	// Associations
	Trip    *Trip    `json:"trip,omitempty" gorm:"foreignKey:TripID"`
	User    *User    `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Payment *Payment `json:"payment,omitempty" gorm:"foreignKey:BookingID"`
	Tickets []Ticket `json:"tickets,omitempty" gorm:"foreignKey:BookingID"`

	CreatedAt   time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
	CancelledAt *time.Time `json:"cancelled_at,omitempty"`
}

// TableName overrides the table name
func (Booking) TableName() string {
	return "bookings"
}

// IsExpired checks if the booking has expired
func (b *Booking) IsExpired() bool {
	if b.Status != BookingStatusPending || b.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*b.ExpiresAt)
}
