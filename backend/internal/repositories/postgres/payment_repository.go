package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/yourusername/bus-booking/internal/entities"
	"gorm.io/gorm"
)

type paymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) *paymentRepository {
	return &paymentRepository{db: db}
}

func (r *paymentRepository) Create(ctx context.Context, payment *entities.Payment) error {
	return r.db.WithContext(ctx).Create(payment).Error
}

func (r *paymentRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.Payment, error) {
	var payment entities.Payment
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&payment).Error
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

func (r *paymentRepository) GetByBookingID(ctx context.Context, bookingID uuid.UUID) (*entities.Payment, error) {
	var payment entities.Payment
	err := r.db.WithContext(ctx).Where("booking_id = ?", bookingID).First(&payment).Error
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

func (r *paymentRepository) GetByTransactionID(ctx context.Context, transactionID string) (*entities.Payment, error) {
	var payment entities.Payment
	err := r.db.WithContext(ctx).Where("transaction_id = ?", transactionID).First(&payment).Error
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

func (r *paymentRepository) GetByGatewayPaymentID(ctx context.Context, gatewayPaymentID string) (*entities.Payment, error) {
	var payment entities.Payment
	err := r.db.WithContext(ctx).Where("gateway_payment_id = ?", gatewayPaymentID).First(&payment).Error
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

func (r *paymentRepository) GetByIdempotencyKey(ctx context.Context, key string) (*entities.Payment, error) {
	var payment entities.Payment
	err := r.db.WithContext(ctx).Where("idempotency_key = ?", key).First(&payment).Error
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

func (r *paymentRepository) Update(ctx context.Context, payment *entities.Payment) error {
	return r.db.WithContext(ctx).Save(payment).Error
}

func (r *paymentRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status entities.PaymentStatus) error {
	return r.db.WithContext(ctx).
		Model(&entities.Payment{}).
		Where("id = ?", id).
		Update("status", status).Error
}

func (r *paymentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entities.Payment{}, id).Error
}

func (r *paymentRepository) List(ctx context.Context, limit, offset int) ([]*entities.Payment, error) {
	var payments []*entities.Payment
	err := r.db.WithContext(ctx).
		Preload("Booking").
		Limit(limit).
		Offset(offset).
		Find(&payments).Error
	return payments, err
}

func (r *paymentRepository) GetByStatus(ctx context.Context, status entities.PaymentStatus, limit, offset int) ([]*entities.Payment, error) {
	var payments []*entities.Payment
	err := r.db.WithContext(ctx).
		Preload("Booking").
		Where("status = ?", status).
		Limit(limit).
		Offset(offset).
		Find(&payments).Error
	return payments, err
}
