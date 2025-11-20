package entities

import (
	"time"

	"github.com/google/uuid"
)

// Ticket represents an e-ticket for a passenger
type Ticket struct {
	ID             uuid.UUID  `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	BookingID      uuid.UUID  `json:"booking_id" gorm:"type:uuid;not null;index"`
	PassengerName  string     `json:"passenger_name" gorm:"not null"`
	PassengerPhone string     `json:"passenger_phone"`
	PassengerEmail string     `json:"passenger_email"`
	SeatNumber     string     `json:"seat_number" gorm:"not null"`
	TicketCode     string     `json:"ticket_code" gorm:"uniqueIndex;not null"` // Unique code for QR
	QRCodePath     string     `json:"qr_code_path"`
	PDFPath        string     `json:"pdf_path"`
	IsCheckedIn    bool       `json:"is_checked_in" gorm:"default:false"`
	CheckedInAt    *time.Time `json:"checked_in_at,omitempty"`

	// Associations
	Booking *Booking `json:"booking,omitempty" gorm:"foreignKey:BookingID"`

	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName overrides the table name
func (Ticket) TableName() string {
	return "tickets"
}
