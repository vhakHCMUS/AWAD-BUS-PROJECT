package repositories

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/yourusername/bus-booking/internal/entities"
)

// UserRepository defines the interface for user data operations
type UserRepository interface {
	Create(ctx context.Context, user *entities.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.User, error)
	GetByEmail(ctx context.Context, email string) (*entities.User, error)
	GetByOAuthID(ctx context.Context, oauthID string, provider string) (*entities.User, error)
	Update(ctx context.Context, user *entities.User) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, limit, offset int) ([]*entities.User, error)
}

// BusRepository defines the interface for bus data operations
type BusRepository interface {
	Create(ctx context.Context, bus *entities.Bus) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Bus, error)
	GetByLicensePlate(ctx context.Context, licensePlate string) (*entities.Bus, error)
	Update(ctx context.Context, bus *entities.Bus) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, status entities.BusStatus, limit, offset int) ([]*entities.Bus, error)
}

// RouteRepository defines the interface for route data operations
type RouteRepository interface {
	Create(ctx context.Context, route *entities.Route) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Route, error)
	Update(ctx context.Context, route *entities.Route) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, limit, offset int) ([]*entities.Route, error)
	Search(ctx context.Context, fromCity, toCity string) ([]*entities.Route, error)
}

// TripRepository defines the interface for trip data operations
type TripRepository interface {
	Create(ctx context.Context, trip *entities.Trip) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Trip, error)
	GetByIDWithDetails(ctx context.Context, id uuid.UUID) (*entities.Trip, error)
	Update(ctx context.Context, trip *entities.Trip) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, status entities.TripStatus, limit, offset int) ([]*entities.Trip, error)
	Search(ctx context.Context, fromCity, toCity string, date time.Time, limit, offset int) ([]*entities.Trip, error)
	GetUpcomingTrips(ctx context.Context, limit int) ([]*entities.Trip, error)
}

// SeatRepository defines the interface for seat data operations
type SeatRepository interface {
	// Batch operations for trip initialization
	InitializeSeatsForTrip(ctx context.Context, tripID uuid.UUID, seatNumbers []string) error

	// Individual seat operations
	GetByTripAndSeat(ctx context.Context, tripID uuid.UUID, seatNumber string) (*entities.SeatInfo, error)
	GetAllByTrip(ctx context.Context, tripID uuid.UUID) ([]*entities.SeatInfo, error)

	// Lock management with row-level locking
	LockSeats(ctx context.Context, tripID uuid.UUID, seatNumbers []string, lockedBy uuid.UUID, duration time.Duration) error
	UnlockSeats(ctx context.Context, tripID uuid.UUID, seatNumbers []string) error
	UnlockExpiredSeats(ctx context.Context) (int, error)

	// Booking operations
	MarkSeatsAsBooked(ctx context.Context, tripID uuid.UUID, seatNumbers []string, bookingID uuid.UUID) error
	ReleaseSeats(ctx context.Context, bookingID uuid.UUID) error

	// Query operations
	GetAvailableSeats(ctx context.Context, tripID uuid.UUID) ([]*entities.SeatInfo, error)
	CountAvailableSeats(ctx context.Context, tripID uuid.UUID) (int, error)
}

// BookingRepository defines the interface for booking data operations
type BookingRepository interface {
	Create(ctx context.Context, booking *entities.Booking) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Booking, error)
	GetByIDWithDetails(ctx context.Context, id uuid.UUID) (*entities.Booking, error)
	GetByCode(ctx context.Context, code string) (*entities.Booking, error)
	Update(ctx context.Context, booking *entities.Booking) error
	UpdateStatus(ctx context.Context, id uuid.UUID, status entities.BookingStatus) error
	Delete(ctx context.Context, id uuid.UUID) error

	// User bookings
	GetByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entities.Booking, error)

	// Expiry management
	GetExpiredBookings(ctx context.Context) ([]*entities.Booking, error)
	MarkAsExpired(ctx context.Context, id uuid.UUID) error

	// Statistics
	GetTripBookings(ctx context.Context, tripID uuid.UUID) ([]*entities.Booking, error)
	CountBookingsByStatus(ctx context.Context, status entities.BookingStatus) (int64, error)
}

// PaymentRepository defines the interface for payment data operations
type PaymentRepository interface {
	Create(ctx context.Context, payment *entities.Payment) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Payment, error)
	GetByBookingID(ctx context.Context, bookingID uuid.UUID) (*entities.Payment, error)
	GetByGatewayPaymentID(ctx context.Context, gatewayPaymentID string) (*entities.Payment, error)
	GetByIdempotencyKey(ctx context.Context, key string) (*entities.Payment, error)
	Update(ctx context.Context, payment *entities.Payment) error
	UpdateStatus(ctx context.Context, id uuid.UUID, status entities.PaymentStatus) error
	List(ctx context.Context, limit, offset int) ([]*entities.Payment, error)
}

// TicketRepository defines the interface for ticket data operations
type TicketRepository interface {
	Create(ctx context.Context, ticket *entities.Ticket) error
	CreateBatch(ctx context.Context, tickets []*entities.Ticket) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Ticket, error)
	GetByCode(ctx context.Context, code string) (*entities.Ticket, error)
	GetByBookingID(ctx context.Context, bookingID uuid.UUID) ([]*entities.Ticket, error)
	Update(ctx context.Context, ticket *entities.Ticket) error
	CheckIn(ctx context.Context, ticketCode string) error
}

// RefreshTokenRepository defines the interface for refresh token operations
type RefreshTokenRepository interface {
	Create(ctx context.Context, token *entities.RefreshToken) error
	GetByToken(ctx context.Context, token string) (*entities.RefreshToken, error)
	Revoke(ctx context.Context, token string) error
	RevokeAllForUser(ctx context.Context, userID uuid.UUID) error
	DeleteExpired(ctx context.Context) error
}
