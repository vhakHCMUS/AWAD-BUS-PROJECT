package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/yourusername/bus-booking/internal/entities"
	"gorm.io/gorm"
)

type tripRepository struct {
	db *gorm.DB
}

func NewTripRepository(db *gorm.DB) *tripRepository {
	return &tripRepository{db: db}
}

func (r *tripRepository) Create(ctx context.Context, trip *entities.Trip) error {
	return r.db.WithContext(ctx).Create(trip).Error
}

func (r *tripRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.Trip, error) {
	var trip entities.Trip
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&trip).Error
	if err != nil {
		return nil, err
	}
	return &trip, nil
}

func (r *tripRepository) GetByIDWithDetails(ctx context.Context, id uuid.UUID) (*entities.Trip, error) {
	var trip entities.Trip
	err := r.db.WithContext(ctx).
		Preload("Route").
		Preload("Bus").
		Where("id = ?", id).
		First(&trip).Error
	if err != nil {
		return nil, err
	}
	return &trip, nil
}

func (r *tripRepository) Update(ctx context.Context, trip *entities.Trip) error {
	return r.db.WithContext(ctx).Save(trip).Error
}

func (r *tripRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entities.Trip{}, id).Error
}

func (r *tripRepository) List(ctx context.Context, status entities.TripStatus, limit, offset int) ([]*entities.Trip, error) {
	var trips []*entities.Trip
	query := r.db.WithContext(ctx).Preload("Route").Preload("Bus")

	if status != "" {
		query = query.Where("status = ?", status)
	}

	err := query.Limit(limit).Offset(offset).Find(&trips).Error
	return trips, err
}

func (r *tripRepository) Search(ctx context.Context, fromCity, toCity string, date time.Time, limit, offset int) ([]*entities.Trip, error) {
	var trips []*entities.Trip
	query := r.db.WithContext(ctx).
		Preload("Route").
		Preload("Bus").
		Joins("JOIN routes ON routes.id = trips.route_id").
		Where("routes.from_city ILIKE ? AND routes.to_city ILIKE ?", "%"+fromCity+"%", "%"+toCity+"%").
		Where("DATE(trips.departure_time) = DATE(?)", date).
		Where("trips.status = ?", entities.TripStatusScheduled)

	err := query.Limit(limit).Offset(offset).Find(&trips).Error
	return trips, err
}

func (r *tripRepository) GetUpcomingTrips(ctx context.Context, limit int) ([]*entities.Trip, error) {
	var trips []*entities.Trip
	err := r.db.WithContext(ctx).
		Preload("Route").
		Preload("Bus").
		Where("departure_time > ? AND status = ?", time.Now(), entities.TripStatusScheduled).
		Order("departure_time ASC").
		Limit(limit).
		Find(&trips).Error
	return trips, err
}
