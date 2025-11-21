package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/yourusername/bus-booking/internal/entities"
	"gorm.io/gorm"
)

type RouteRepository struct {
	db *gorm.DB
}

func NewRouteRepository(db *gorm.DB) *RouteRepository {
	return &RouteRepository{db: db}
}

func (r *RouteRepository) Create(ctx context.Context, route *entities.Route) error {
	return r.db.WithContext(ctx).Create(route).Error
}

func (r *RouteRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.Route, error) {
	var route entities.Route
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&route).Error
	if err != nil {
		return nil, err
	}
	return &route, nil
}

func (r *RouteRepository) Update(ctx context.Context, route *entities.Route) error {
	return r.db.WithContext(ctx).Save(route).Error
}

func (r *RouteRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entities.Route{}, "id = ?", id).Error
}

func (r *RouteRepository) List(ctx context.Context, limit, offset int) ([]*entities.Route, error) {
	var routes []*entities.Route
	err := r.db.WithContext(ctx).
		Where("is_active = ?", true).
		Limit(limit).
		Offset(offset).
		Order("name ASC").
		Find(&routes).Error
	return routes, err
}

func (r *RouteRepository) Search(ctx context.Context, fromCity, toCity string) ([]*entities.Route, error) {
	var routes []*entities.Route
	query := r.db.WithContext(ctx).Where("is_active = ?", true)

	if fromCity != "" {
		query = query.Where("LOWER(from_city) LIKE LOWER(?)", "%"+fromCity+"%")
	}
	if toCity != "" {
		query = query.Where("LOWER(to_city) LIKE LOWER(?)", "%"+toCity+"%")
	}

	err := query.Order("name ASC").Find(&routes).Error
	return routes, err
}
