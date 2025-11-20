package usecases

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/yourusername/bus-booking/internal/entities"
	"github.com/yourusername/bus-booking/internal/repositories"
	"github.com/yourusername/bus-booking/internal/repositories/cache"
)

// BookingUsecase handles booking business logic
type BookingUsecase struct {
	bookingRepo repositories.BookingRepository
	seatRepo    repositories.SeatRepository
	tripRepo    repositories.TripRepository
	paymentRepo repositories.PaymentRepository
	ticketRepo  repositories.TicketRepository
	cache       *cache.RedisCache

	seatLockDuration time.Duration
	bookingExpiry    time.Duration
}

// NewBookingUsecase creates a new booking usecase
func NewBookingUsecase(
	bookingRepo repositories.BookingRepository,
	seatRepo repositories.SeatRepository,
	tripRepo repositories.TripRepository,
	paymentRepo repositories.PaymentRepository,
	ticketRepo repositories.TicketRepository,
	cache *cache.RedisCache,
	seatLockDuration time.Duration,
	bookingExpiry time.Duration,
) *BookingUsecase {
	return &BookingUsecase{
		bookingRepo:      bookingRepo,
		seatRepo:         seatRepo,
		tripRepo:         tripRepo,
		paymentRepo:      paymentRepo,
		ticketRepo:       ticketRepo,
		cache:            cache,
		seatLockDuration: seatLockDuration,
		bookingExpiry:    bookingExpiry,
	}
}

// InitiateBooking starts the booking process by locking seats
func (uc *BookingUsecase) InitiateBooking(ctx context.Context, tripID uuid.UUID, seatNumbers []string, userID *uuid.UUID) (*entities.Booking, error) {
	// Validate trip exists
	trip, err := uc.tripRepo.GetByIDWithDetails(ctx, tripID)
	if err != nil {
		return nil, fmt.Errorf("trip not found: %w", err)
	}

	if trip.Status != entities.TripStatusScheduled {
		return nil, fmt.Errorf("trip is not available for booking")
	}

	// Generate lock ID (user ID or session ID)
	lockID := uuid.New()
	if userID != nil {
		lockID = *userID
	}

	// Attempt to lock seats in both Redis and PostgreSQL
	// 1. Redis lock for fast distributed locking
	for _, seatNum := range seatNumbers {
		err := uc.cache.LockSeat(ctx, tripID, seatNum, lockID, uc.seatLockDuration)
		if err != nil {
			// Rollback Redis locks
			for i := 0; i < len(seatNumbers); i++ {
				_ = uc.cache.UnlockSeat(ctx, tripID, seatNumbers[i])
			}
			return nil, fmt.Errorf("failed to lock seat %s: %w", seatNum, err)
		}
	}

	// 2. PostgreSQL lock with SELECT FOR UPDATE
	err = uc.seatRepo.LockSeats(ctx, tripID, seatNumbers, lockID, uc.seatLockDuration)
	if err != nil {
		// Rollback Redis locks
		for _, seatNum := range seatNumbers {
			_ = uc.cache.UnlockSeat(ctx, tripID, seatNum)
		}
		return nil, fmt.Errorf("failed to lock seats in database: %w", err)
	}

	// Create pending booking
	booking := &entities.Booking{
		TripID:      tripID,
		UserID:      userID,
		Seats:       seatNumbers,
		TotalPrice:  trip.Price * float64(len(seatNumbers)),
		Status:      entities.BookingStatusPending,
		BookingCode: generateBookingCode(),
	}

	expiresAt := time.Now().Add(uc.bookingExpiry)
	booking.ExpiresAt = &expiresAt

	err = uc.bookingRepo.Create(ctx, booking)
	if err != nil {
		return nil, fmt.Errorf("failed to create booking: %w", err)
	}

	// Invalidate cache and publish seat update
	_ = uc.cache.InvalidateTripSeats(ctx, tripID)
	uc.publishSeatUpdates(ctx, tripID, seatNumbers, entities.SeatStatusLocked)

	return booking, nil
}

// ConfirmBooking confirms a booking after successful payment
func (uc *BookingUsecase) ConfirmBooking(ctx context.Context, bookingID uuid.UUID, paymentID uuid.UUID) error {
	booking, err := uc.bookingRepo.GetByIDWithDetails(ctx, bookingID)
	if err != nil {
		return fmt.Errorf("booking not found: %w", err)
	}

	if booking.Status != entities.BookingStatusPending {
		return fmt.Errorf("booking is not in pending status")
	}

	// Verify payment
	payment, err := uc.paymentRepo.GetByID(ctx, paymentID)
	if err != nil {
		return fmt.Errorf("payment not found: %w", err)
	}

	if payment.Status != entities.PaymentStatusCompleted {
		return fmt.Errorf("payment is not completed")
	}

	// Mark seats as booked
	err = uc.seatRepo.MarkSeatsAsBooked(ctx, booking.TripID, booking.Seats, bookingID)
	if err != nil {
		return fmt.Errorf("failed to mark seats as booked: %w", err)
	}

	// Update booking status
	booking.Status = entities.BookingStatusConfirmed
	err = uc.bookingRepo.Update(ctx, booking)
	if err != nil {
		return fmt.Errorf("failed to update booking: %w", err)
	}

	// Generate tickets
	err = uc.generateTickets(ctx, booking)
	if err != nil {
		return fmt.Errorf("failed to generate tickets: %w", err)
	}

	// Clear Redis locks and invalidate cache
	for _, seatNum := range booking.Seats {
		_ = uc.cache.UnlockSeat(ctx, booking.TripID, seatNum)
	}
	_ = uc.cache.InvalidateTripSeats(ctx, booking.TripID)

	// Publish seat updates
	uc.publishSeatUpdates(ctx, booking.TripID, booking.Seats, entities.SeatStatusBooked)

	return nil
}

// CancelBooking cancels a booking and releases seats
func (uc *BookingUsecase) CancelBooking(ctx context.Context, bookingID uuid.UUID) error {
	booking, err := uc.bookingRepo.GetByID(ctx, bookingID)
	if err != nil {
		return fmt.Errorf("booking not found: %w", err)
	}

	if booking.Status == entities.BookingStatusCancelled || booking.Status == entities.BookingStatusExpired {
		return fmt.Errorf("booking is already cancelled or expired")
	}

	// Release seats
	err = uc.seatRepo.ReleaseSeats(ctx, bookingID)
	if err != nil {
		return fmt.Errorf("failed to release seats: %w", err)
	}

	// Update booking status
	now := time.Now()
	booking.Status = entities.BookingStatusCancelled
	booking.CancelledAt = &now
	err = uc.bookingRepo.Update(ctx, booking)
	if err != nil {
		return fmt.Errorf("failed to update booking: %w", err)
	}

	// Clear Redis locks
	for _, seatNum := range booking.Seats {
		_ = uc.cache.UnlockSeat(ctx, booking.TripID, seatNum)
	}
	_ = uc.cache.InvalidateTripSeats(ctx, booking.TripID)

	// Publish updates
	uc.publishSeatUpdates(ctx, booking.TripID, booking.Seats, entities.SeatStatusAvailable)

	return nil
}

// GetBookingByID retrieves a booking by its ID
func (uc *BookingUsecase) GetBookingByID(ctx context.Context, bookingID uuid.UUID) (*entities.Booking, error) {
	booking, err := uc.bookingRepo.GetByIDWithDetails(ctx, bookingID)
	if err != nil {
		return nil, fmt.Errorf("booking not found: %w", err)
	}
	return booking, nil
}

// GetUserBookings retrieves all bookings for a specific user
func (uc *BookingUsecase) GetUserBookings(ctx context.Context, userID uuid.UUID, page, limit int) ([]*entities.Booking, error) {
	offset := (page - 1) * limit
	bookings, err := uc.bookingRepo.GetByUser(ctx, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get user bookings: %w", err)
	}
	return bookings, nil
}

// GetAvailableSeats retrieves available seats for a trip (with caching)
func (uc *BookingUsecase) GetAvailableSeats(ctx context.Context, tripID uuid.UUID) ([]*entities.SeatInfo, error) {
	// Try cache first
	seats, err := uc.cache.GetTripSeats(ctx, tripID)
	if err == nil && seats != nil {
		return seats, nil
	}

	// Cache miss - fetch from database
	seats, err = uc.seatRepo.GetAllByTrip(ctx, tripID)
	if err != nil {
		return nil, err
	}

	// Cache for 30 seconds
	_ = uc.cache.CacheTripSeats(ctx, tripID, seats, 30*time.Second)

	return seats, nil
}

// Helper functions
func generateBookingCode() string {
	return fmt.Sprintf("BK%d%s", time.Now().Unix(), uuid.New().String()[:8])
}

func (uc *BookingUsecase) generateTickets(ctx context.Context, booking *entities.Booking) error {
	tickets := make([]*entities.Ticket, len(booking.Seats))

	for i, seatNum := range booking.Seats {
		tickets[i] = &entities.Ticket{
			BookingID:      booking.ID,
			PassengerName:  booking.ContactName,
			PassengerPhone: booking.ContactPhone,
			PassengerEmail: booking.ContactEmail,
			SeatNumber:     seatNum,
			TicketCode:     fmt.Sprintf("TK%d%s", time.Now().Unix(), uuid.New().String()[:8]),
		}
	}

	return uc.ticketRepo.CreateBatch(ctx, tickets)
}

func (uc *BookingUsecase) publishSeatUpdates(ctx context.Context, tripID uuid.UUID, seatNumbers []string, status entities.SeatStatus) {
	for _, seatNum := range seatNumbers {
		seat := &entities.SeatInfo{
			TripID:     tripID,
			SeatNumber: seatNum,
			Status:     status,
		}
		_ = uc.cache.PublishSeatUpdate(ctx, tripID, seat)
	}
}

// ExpireOldBookings is a cleanup job that expires pending bookings
func (uc *BookingUsecase) ExpireOldBookings(ctx context.Context) error {
	expiredBookings, err := uc.bookingRepo.GetExpiredBookings(ctx)
	if err != nil {
		return err
	}

	for _, booking := range expiredBookings {
		_ = uc.CancelBooking(ctx, booking.ID)
	}

	return nil
}
