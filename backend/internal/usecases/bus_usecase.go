package usecases

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/yourusername/bus-booking/internal/entities"
	"github.com/yourusername/bus-booking/internal/repositories"
)

type BusUsecase struct {
	busRepo repositories.BusRepository
}

func NewBusUsecase(busRepo repositories.BusRepository) *BusUsecase {
	return &BusUsecase{busRepo: busRepo}
}

func (uc *BusUsecase) CreateBus(ctx context.Context, bus *entities.Bus) error {
	// Check for duplicate license plate
	existing, err := uc.busRepo.GetByLicensePlate(ctx, bus.LicensePlate)
	if err == nil && existing != nil {
		return fmt.Errorf("bus with license plate %s already exists", bus.LicensePlate)
	}

	// Validate seat layout
	if err := validateSeatLayout(&bus.SeatLayout); err != nil {
		return fmt.Errorf("invalid seat layout: %w", err)
	}

	return uc.busRepo.Create(ctx, bus)
}

func (uc *BusUsecase) GetBusByID(ctx context.Context, id uuid.UUID) (*entities.Bus, error) {
	return uc.busRepo.GetByID(ctx, id)
}

func (uc *BusUsecase) UpdateBus(ctx context.Context, bus *entities.Bus) error {
	// Validate bus exists
	_, err := uc.busRepo.GetByID(ctx, bus.ID)
	if err != nil {
		return fmt.Errorf("bus not found: %w", err)
	}

	// Validate seat layout
	if err := validateSeatLayout(&bus.SeatLayout); err != nil {
		return fmt.Errorf("invalid seat layout: %w", err)
	}

	return uc.busRepo.Update(ctx, bus)
}

func (uc *BusUsecase) DeleteBus(ctx context.Context, id uuid.UUID) error {
	return uc.busRepo.Delete(ctx, id)
}

func (uc *BusUsecase) ListBuses(ctx context.Context, status entities.BusStatus, page, limit int) ([]*entities.Bus, error) {
	offset := (page - 1) * limit
	return uc.busRepo.List(ctx, status, limit, offset)
}

func validateSeatLayout(layout *entities.SeatLayout) error {
	if layout.Rows <= 0 || layout.Columns <= 0 {
		return fmt.Errorf("rows and columns must be positive")
	}
	if layout.TotalSeats <= 0 {
		return fmt.Errorf("total seats must be positive")
	}
	if layout.Floors < 1 || layout.Floors > 2 {
		return fmt.Errorf("floors must be 1 or 2")
	}
	if len(layout.Layout) == 0 {
		return fmt.Errorf("layout cannot be empty")
	}
	return nil
}
