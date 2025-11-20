package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/yourusername/bus-booking/internal/entities"
	"gorm.io/gorm"
)

type bookingRepository struct {
	db *gorm.DB
}

// NewBookingRepository creates a new booking repository
func NewBookingRepository(db *gorm.DB) *bookingRepository {
	return &bookingRepository{db: db}
}

func (r *bookingRepository) Create(ctx context.Context, booking *entities.Booking) error {
	return r.db.WithContext(ctx).Create(booking).Error
}

func (r *bookingRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.Booking, error) {
	var booking entities.Booking
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&booking).Error
	if err != nil {
		return nil, err
	}
	return &booking, nil
}

func (r *bookingRepository) GetByIDWithDetails(ctx context.Context, id uuid.UUID) (*entities.Booking, error) {
	var booking entities.Booking
	err := r.db.WithContext(ctx).
		Preload("Trip.Route").
		Preload("Trip.Bus").
		Preload("User").
		Preload("Payment").
		Preload("Tickets").
		Where("id = ?", id).
		First(&booking).Error
	if err != nil {
		return nil, err
	}
	return &booking, nil
}

func (r *bookingRepository) GetByCode(ctx context.Context, code string) (*entities.Booking, error) {
	var booking entities.Booking
	err := r.db.WithContext(ctx).
		Preload("Trip.Route").
		Preload("Trip.Bus").
		Preload("Payment").
		Preload("Tickets").
		Where("booking_code = ?", code).
		First(&booking).Error
	if err != nil {
		return nil, err
	}
	return &booking, nil
}

func (r *bookingRepository) Update(ctx context.Context, booking *entities.Booking) error {
	return r.db.WithContext(ctx).Save(booking).Error
}

func (r *bookingRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status entities.BookingStatus) error {
	return r.db.WithContext(ctx).Model(&entities.Booking{}).
		Where("id = ?", id).
		Update("status", status).Error
}

func (r *bookingRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entities.Booking{}, id).Error
}

func (r *bookingRepository) GetByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entities.Booking, error) {
	var bookings []*entities.Booking
	err := r.db.WithContext(ctx).
		Preload("Trip.Route").
		Preload("Trip.Bus").
		Preload("Payment").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&bookings).Error
	return bookings, err
}

func (r *bookingRepository) GetExpiredBookings(ctx context.Context) ([]*entities.Booking, error) {
	var bookings []*entities.Booking
	err := r.db.WithContext(ctx).
		Where("status = ? AND expires_at < NOW()", entities.BookingStatusPending).
		Find(&bookings).Error
	return bookings, err
}

func (r *bookingRepository) MarkAsExpired(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&entities.Booking{}).
		Where("id = ?", id).
		Update("status", entities.BookingStatusExpired).Error
}

func (r *bookingRepository) GetTripBookings(ctx context.Context, tripID uuid.UUID) ([]*entities.Booking, error) {
	var bookings []*entities.Booking
	err := r.db.WithContext(ctx).
		Preload("User").
		Where("trip_id = ? AND status IN ?", tripID, []entities.BookingStatus{
			entities.BookingStatusPaid,
			entities.BookingStatusConfirmed,
		}).
		Find(&bookings).Error
	return bookings, err
}

func (r *bookingRepository) CountBookingsByStatus(ctx context.Context, status entities.BookingStatus) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&entities.Booking{}).
		Where("status = ?", status).
		Count(&count).Error
	return count, err
}
