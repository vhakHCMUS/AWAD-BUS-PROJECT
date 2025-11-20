package entities

import (
	"time"

	"github.com/google/uuid"
)

// Route represents a bus route from one city to another
type Route struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name        string    `json:"name" gorm:"not null"` // e.g., "Hanoi - Ho Chi Minh"
	FromCity    string    `json:"from_city" gorm:"index;not null"`
	ToCity      string    `json:"to_city" gorm:"index;not null"`
	Distance    float64   `json:"distance" gorm:"not null"`   // in kilometers
	BasePrice   float64   `json:"base_price" gorm:"not null"` // base ticket price
	Description string    `json:"description"`
	IsActive    bool      `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName overrides the table name
func (Route) TableName() string {
	return "routes"
}

// TripStatus represents the status of a trip
type TripStatus string

const (
	TripStatusScheduled TripStatus = "scheduled"
	TripStatusBoarding  TripStatus = "boarding"
	TripStatusInTransit TripStatus = "in_transit"
	TripStatusCompleted TripStatus = "completed"
	TripStatusCancelled TripStatus = "cancelled"
	TripStatusDelayed   TripStatus = "delayed"
)

// Trip represents a scheduled trip for a specific route
type Trip struct {
	ID            uuid.UUID  `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	BusID         uuid.UUID  `json:"bus_id" gorm:"type:uuid;not null;index"`
	RouteID       uuid.UUID  `json:"route_id" gorm:"type:uuid;not null;index"`
	DepartureTime time.Time  `json:"departure_time" gorm:"index;not null"`
	ArrivalTime   time.Time  `json:"arrival_time" gorm:"not null"`
	Duration      int        `json:"duration"` // in minutes
	Price         float64    `json:"price" gorm:"not null"`
	Status        TripStatus `json:"status" gorm:"type:varchar(20);not null;default:'scheduled'"`
	DriverName    string     `json:"driver_name"`
	DriverPhone   string     `json:"driver_phone"`

	// Associations
	Bus   *Bus   `json:"bus,omitempty" gorm:"foreignKey:BusID"`
	Route *Route `json:"route,omitempty" gorm:"foreignKey:RouteID"`

	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName overrides the table name
func (Trip) TableName() string {
	return "trips"
}
