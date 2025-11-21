package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/yourusername/bus-booking/internal/entities"
	"gorm.io/gorm"
)

type BusRepository struct {
	db *gorm.DB
}

func NewBusRepository(db *gorm.DB) *BusRepository {
	return &BusRepository{db: db}
}

func (r *BusRepository) Create(ctx context.Context, bus *entities.Bus) error {
	return r.db.WithContext(ctx).Create(bus).Error
}

func (r *BusRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.Bus, error) {
	var bus entities.Bus
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&bus).Error
	if err != nil {
		return nil, err
	}
	return &bus, nil
}

func (r *BusRepository) GetByLicensePlate(ctx context.Context, licensePlate string) (*entities.Bus, error) {
	var bus entities.Bus
	err := r.db.WithContext(ctx).Where("license_plate = ?", licensePlate).First(&bus).Error
	if err != nil {
		return nil, err
	}
	return &bus, nil
}

func (r *BusRepository) Update(ctx context.Context, bus *entities.Bus) error {
	return r.db.WithContext(ctx).Save(bus).Error
}

func (r *BusRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entities.Bus{}, "id = ?", id).Error
}

func (r *BusRepository) List(ctx context.Context, status entities.BusStatus, limit, offset int) ([]*entities.Bus, error) {
	var buses []*entities.Bus
	query := r.db.WithContext(ctx)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	err := query.Limit(limit).Offset(offset).Order("created_at DESC").Find(&buses).Error
	return buses, err
}
