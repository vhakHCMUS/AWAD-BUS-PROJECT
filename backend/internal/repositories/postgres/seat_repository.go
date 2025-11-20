package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/yourusername/bus-booking/internal/entities"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type seatRepository struct {
	db *gorm.DB
}

// NewSeatRepository creates a new seat repository
func NewSeatRepository(db *gorm.DB) *seatRepository {
	return &seatRepository{db: db}
}

func (r *seatRepository) InitializeSeatsForTrip(ctx context.Context, tripID uuid.UUID, seatNumbers []string) error {
	seats := make([]*entities.SeatInfo, len(seatNumbers))
	for i, seatNum := range seatNumbers {
		seats[i] = &entities.SeatInfo{
			TripID:     tripID,
			SeatNumber: seatNum,
			Status:     entities.SeatStatusAvailable,
		}
	}
	return r.db.WithContext(ctx).CreateInBatches(seats, 100).Error
}

func (r *seatRepository) GetByTripAndSeat(ctx context.Context, tripID uuid.UUID, seatNumber string) (*entities.SeatInfo, error) {
	var seat entities.SeatInfo
	err := r.db.WithContext(ctx).
		Where("trip_id = ? AND seat_number = ?", tripID, seatNumber).
		First(&seat).Error
	if err != nil {
		return nil, err
	}
	return &seat, nil
}

func (r *seatRepository) GetAllByTrip(ctx context.Context, tripID uuid.UUID) ([]*entities.SeatInfo, error) {
	var seats []*entities.SeatInfo
	err := r.db.WithContext(ctx).
		Where("trip_id = ?", tripID).
		Order("seat_number ASC").
		Find(&seats).Error
	return seats, err
}

// LockSeats implements distributed locking with PostgreSQL row-level locks (SELECT FOR UPDATE)
func (r *seatRepository) LockSeats(ctx context.Context, tripID uuid.UUID, seatNumbers []string, lockedBy uuid.UUID, duration time.Duration) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var seats []*entities.SeatInfo

		// Acquire row-level locks using SELECT FOR UPDATE
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("trip_id = ? AND seat_number IN ?", tripID, seatNumbers).
			Find(&seats).Error
		if err != nil {
			return err
		}

		if len(seats) != len(seatNumbers) {
			return fmt.Errorf("some seats not found")
		}

		// Check if all seats are available or lock expired
		now := time.Now()
		for _, seat := range seats {
			if seat.Status == entities.SeatStatusBooked {
				return fmt.Errorf("seat %s is already booked", seat.SeatNumber)
			}
			if seat.Status == entities.SeatStatusLocked {
				if seat.LockedUntil != nil && now.Before(*seat.LockedUntil) {
					return fmt.Errorf("seat %s is currently locked", seat.SeatNumber)
				}
			}
		}

		// Lock the seats
		lockedUntil := now.Add(duration)
		return tx.Model(&entities.SeatInfo{}).
			Where("trip_id = ? AND seat_number IN ?", tripID, seatNumbers).
			Updates(map[string]interface{}{
				"status":       entities.SeatStatusLocked,
				"locked_until": lockedUntil,
				"locked_by":    lockedBy,
			}).Error
	})
}

func (r *seatRepository) UnlockSeats(ctx context.Context, tripID uuid.UUID, seatNumbers []string) error {
	return r.db.WithContext(ctx).Model(&entities.SeatInfo{}).
		Where("trip_id = ? AND seat_number IN ? AND status = ?", tripID, seatNumbers, entities.SeatStatusLocked).
		Updates(map[string]interface{}{
			"status":       entities.SeatStatusAvailable,
			"locked_until": nil,
			"locked_by":    nil,
		}).Error
}

func (r *seatRepository) UnlockExpiredSeats(ctx context.Context) (int, error) {
	result := r.db.WithContext(ctx).Model(&entities.SeatInfo{}).
		Where("status = ? AND locked_until < ?", entities.SeatStatusLocked, time.Now()).
		Updates(map[string]interface{}{
			"status":       entities.SeatStatusAvailable,
			"locked_until": nil,
			"locked_by":    nil,
		})
	return int(result.RowsAffected), result.Error
}

func (r *seatRepository) MarkSeatsAsBooked(ctx context.Context, tripID uuid.UUID, seatNumbers []string, bookingID uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&entities.SeatInfo{}).
		Where("trip_id = ? AND seat_number IN ?", tripID, seatNumbers).
		Updates(map[string]interface{}{
			"status":       entities.SeatStatusBooked,
			"booking_id":   bookingID,
			"locked_until": nil,
			"locked_by":    nil,
		}).Error
}

func (r *seatRepository) ReleaseSeats(ctx context.Context, bookingID uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&entities.SeatInfo{}).
		Where("booking_id = ?", bookingID).
		Updates(map[string]interface{}{
			"status":     entities.SeatStatusAvailable,
			"booking_id": nil,
		}).Error
}

func (r *seatRepository) GetAvailableSeats(ctx context.Context, tripID uuid.UUID) ([]*entities.SeatInfo, error) {
	var seats []*entities.SeatInfo
	err := r.db.WithContext(ctx).
		Where("trip_id = ? AND (status = ? OR (status = ? AND locked_until < ?))",
			tripID,
			entities.SeatStatusAvailable,
			entities.SeatStatusLocked,
			time.Now()).
		Order("seat_number ASC").
		Find(&seats).Error
	return seats, err
}

func (r *seatRepository) CountAvailableSeats(ctx context.Context, tripID uuid.UUID) (int, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&entities.SeatInfo{}).
		Where("trip_id = ? AND (status = ? OR (status = ? AND locked_until < ?))",
			tripID,
			entities.SeatStatusAvailable,
			entities.SeatStatusLocked,
			time.Now()).
		Count(&count).Error
	return int(count), err
}
