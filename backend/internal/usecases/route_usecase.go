package usecases

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/yourusername/bus-booking/internal/entities"
	"github.com/yourusername/bus-booking/internal/repositories"
)

type RouteUsecase struct {
	routeRepo repositories.RouteRepository
}

func NewRouteUsecase(routeRepo repositories.RouteRepository) *RouteUsecase {
	return &RouteUsecase{routeRepo: routeRepo}
}

func (uc *RouteUsecase) CreateRoute(ctx context.Context, route *entities.Route) error {
	// Validate route data
	if route.FromCity == "" || route.ToCity == "" {
		return fmt.Errorf("from_city and to_city are required")
	}
	if route.Distance <= 0 {
		return fmt.Errorf("distance must be positive")
	}
	if route.BasePrice <= 0 {
		return fmt.Errorf("base price must be positive")
	}

	// Generate route name if not provided
	if route.Name == "" {
		route.Name = fmt.Sprintf("%s - %s", route.FromCity, route.ToCity)
	}

	return uc.routeRepo.Create(ctx, route)
}

func (uc *RouteUsecase) GetRouteByID(ctx context.Context, id uuid.UUID) (*entities.Route, error) {
	return uc.routeRepo.GetByID(ctx, id)
}

func (uc *RouteUsecase) UpdateRoute(ctx context.Context, route *entities.Route) error {
	// Validate route exists
	_, err := uc.routeRepo.GetByID(ctx, route.ID)
	if err != nil {
		return fmt.Errorf("route not found: %w", err)
	}

	// Validate route data
	if route.Distance <= 0 {
		return fmt.Errorf("distance must be positive")
	}
	if route.BasePrice <= 0 {
		return fmt.Errorf("base price must be positive")
	}

	return uc.routeRepo.Update(ctx, route)
}

func (uc *RouteUsecase) DeleteRoute(ctx context.Context, id uuid.UUID) error {
	return uc.routeRepo.Delete(ctx, id)
}

func (uc *RouteUsecase) ListRoutes(ctx context.Context, page, limit int) ([]*entities.Route, error) {
	offset := (page - 1) * limit
	return uc.routeRepo.List(ctx, limit, offset)
}

func (uc *RouteUsecase) SearchRoutes(ctx context.Context, fromCity, toCity string) ([]*entities.Route, error) {
	return uc.routeRepo.Search(ctx, fromCity, toCity)
}
