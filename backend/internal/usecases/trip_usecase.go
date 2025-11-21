package usecases

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/yourusername/bus-booking/internal/entities"
	"github.com/yourusername/bus-booking/internal/repositories"
)

type TripUsecase struct {
	tripRepo  repositories.TripRepository
	busRepo   repositories.BusRepository
	routeRepo repositories.RouteRepository
	seatRepo  repositories.SeatRepository
}

func NewTripUsecase(
	tripRepo repositories.TripRepository,
	busRepo repositories.BusRepository,
	routeRepo repositories.RouteRepository,
	seatRepo repositories.SeatRepository,
) *TripUsecase {
	return &TripUsecase{
		tripRepo:  tripRepo,
		busRepo:   busRepo,
		routeRepo: routeRepo,
		seatRepo:  seatRepo,
	}
}

func (uc *TripUsecase) CreateTrip(ctx context.Context, trip *entities.Trip) error {
	// Validate bus exists and is active
	bus, err := uc.busRepo.GetByID(ctx, trip.BusID)
	if err != nil {
		return fmt.Errorf("bus not found: %w", err)
	}
	if bus.Status != entities.BusStatusActive {
		return fmt.Errorf("bus is not active")
	}

	// Validate route exists and is active
	route, err := uc.routeRepo.GetByID(ctx, trip.RouteID)
	if err != nil {
		return fmt.Errorf("route not found: %w", err)
	}
	if !route.IsActive {
		return fmt.Errorf("route is not active")
	}

	// Calculate duration if not provided
	if trip.Duration == 0 {
		trip.Duration = int(trip.ArrivalTime.Sub(trip.DepartureTime).Minutes())
	}

	// Set default price if not provided
	if trip.Price == 0 {
		trip.Price = route.BasePrice
	}

	// Create trip
	err = uc.tripRepo.Create(ctx, trip)
	if err != nil {
		return fmt.Errorf("failed to create trip: %w", err)
	}

	// Initialize seats for this trip
	seatNumbers := generateSeatNumbers(bus.SeatLayout)
	err = uc.seatRepo.InitializeSeatsForTrip(ctx, trip.ID, seatNumbers)
	if err != nil {
		return fmt.Errorf("failed to initialize seats: %w", err)
	}

	return nil
}

func (uc *TripUsecase) GetTripByID(ctx context.Context, id uuid.UUID) (*entities.Trip, error) {
	return uc.tripRepo.GetByIDWithDetails(ctx, id)
}

func (uc *TripUsecase) UpdateTrip(ctx context.Context, trip *entities.Trip) error {
	// Validate trip exists
	existingTrip, err := uc.tripRepo.GetByID(ctx, trip.ID)
	if err != nil {
		return fmt.Errorf("trip not found: %w", err)
	}

	// Don't allow updating certain fields if trip has bookings
	if existingTrip.Status != entities.TripStatusScheduled {
		return fmt.Errorf("cannot update trip that is not in scheduled status")
	}

	return uc.tripRepo.Update(ctx, trip)
}

func (uc *TripUsecase) DeleteTrip(ctx context.Context, id uuid.UUID) error {
	trip, err := uc.tripRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("trip not found: %w", err)
	}

	// Don't allow deleting trips with active bookings
	if trip.Status == entities.TripStatusInTransit || trip.Status == entities.TripStatusBoarding {
		return fmt.Errorf("cannot delete trip in progress")
	}

	return uc.tripRepo.Delete(ctx, id)
}

func (uc *TripUsecase) SearchTrips(ctx context.Context, fromCity, toCity string, date time.Time, page, limit int) ([]*entities.Trip, error) {
	offset := (page - 1) * limit
	return uc.tripRepo.Search(ctx, fromCity, toCity, date, limit, offset)
}

func (uc *TripUsecase) ListTrips(ctx context.Context, status entities.TripStatus, page, limit int) ([]*entities.Trip, error) {
	offset := (page - 1) * limit
	return uc.tripRepo.List(ctx, status, limit, offset)
}

func (uc *TripUsecase) GetUpcomingTrips(ctx context.Context, limit int) ([]*entities.Trip, error) {
	return uc.tripRepo.GetUpcomingTrips(ctx, limit)
}

func generateSeatNumbers(layout entities.SeatLayout) []string {
	seats := make([]string, 0, layout.TotalSeats)
	for _, row := range layout.Layout {
		for _, seat := range row {
			if seat != "aisle" && seat != "empty" && seat != "" {
				seats = append(seats, seat)
			}
		}
	}
	return seats
}
