package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/yourusername/bus-booking/internal/entities"
	"gorm.io/gorm"
)

type ticketRepository struct {
	db *gorm.DB
}

func NewTicketRepository(db *gorm.DB) *ticketRepository {
	return &ticketRepository{db: db}
}

func (r *ticketRepository) Create(ctx context.Context, ticket *entities.Ticket) error {
	return r.db.WithContext(ctx).Create(ticket).Error
}

func (r *ticketRepository) CreateBatch(ctx context.Context, tickets []*entities.Ticket) error {
	return r.db.WithContext(ctx).Create(&tickets).Error
}

func (r *ticketRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.Ticket, error) {
	var ticket entities.Ticket
	err := r.db.WithContext(ctx).
		Preload("Booking").
		Where("id = ?", id).
		First(&ticket).Error
	if err != nil {
		return nil, err
	}
	return &ticket, nil
}

func (r *ticketRepository) GetByCode(ctx context.Context, code string) (*entities.Ticket, error) {
	var ticket entities.Ticket
	err := r.db.WithContext(ctx).
		Preload("Booking").
		Preload("Booking.Trip").
		Preload("Booking.Trip.Route").
		Preload("Booking.Trip.Bus").
		Where("ticket_code = ?", code).
		First(&ticket).Error
	if err != nil {
		return nil, err
	}
	return &ticket, nil
}

func (r *ticketRepository) GetByBookingID(ctx context.Context, bookingID uuid.UUID) ([]*entities.Ticket, error) {
	var tickets []*entities.Ticket
	err := r.db.WithContext(ctx).
		Where("booking_id = ?", bookingID).
		Find(&tickets).Error
	return tickets, err
}

func (r *ticketRepository) Update(ctx context.Context, ticket *entities.Ticket) error {
	return r.db.WithContext(ctx).Save(ticket).Error
}

func (r *ticketRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entities.Ticket{}, id).Error
}

func (r *ticketRepository) List(ctx context.Context, limit, offset int) ([]*entities.Ticket, error) {
	var tickets []*entities.Ticket
	err := r.db.WithContext(ctx).
		Preload("Booking").
		Limit(limit).
		Offset(offset).
		Find(&tickets).Error
	return tickets, err
}

func (r *ticketRepository) MarkAsUsed(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&entities.Ticket{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"used_at": gorm.Expr("CURRENT_TIMESTAMP"),
		}).Error
}

func (r *ticketRepository) CheckIn(ctx context.Context, ticketCode string) error {
	return r.db.WithContext(ctx).
		Model(&entities.Ticket{}).
		Where("ticket_code = ?", ticketCode).
		Updates(map[string]interface{}{
			"checked_in_at": gorm.Expr("CURRENT_TIMESTAMP"),
		}).Error
}
