package entities

import (
	"time"

	"github.com/google/uuid"
)

// PaymentGateway represents the payment provider
type PaymentGateway string

const (
	PaymentGatewayPayOS PaymentGateway = "payos"
	PaymentGatewayMoMo  PaymentGateway = "momo"
)

// PaymentStatus represents the status of a payment
type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "pending"
	PaymentStatusCompleted PaymentStatus = "completed"
	PaymentStatusFailed    PaymentStatus = "failed"
	PaymentStatusRefunded  PaymentStatus = "refunded"
	PaymentStatusCancelled PaymentStatus = "cancelled"
)

// Payment represents a payment transaction
type Payment struct {
	ID               uuid.UUID      `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	BookingID        uuid.UUID      `json:"booking_id" gorm:"type:uuid;not null;uniqueIndex"`
	Gateway          PaymentGateway `json:"gateway" gorm:"type:varchar(20);not null"`
	GatewayPaymentID string         `json:"gateway_payment_id" gorm:"uniqueIndex"` // Payment ID from gateway
	Amount           float64        `json:"amount" gorm:"not null"`
	Currency         string         `json:"currency" gorm:"type:varchar(3);default:'VND'"`
	Status           PaymentStatus  `json:"status" gorm:"type:varchar(20);not null;default:'pending';index"`
	PaymentMethod    string         `json:"payment_method,omitempty"` // "credit_card", "wallet", etc.
	TransactionID    string         `json:"transaction_id,omitempty"` // Bank transaction ID
	FailureReason    *string        `json:"failure_reason,omitempty"`
	IdempotencyKey   string         `json:"idempotency_key" gorm:"uniqueIndex;not null"` // For webhook deduplication

	// Associations
	Booking *Booking `json:"booking,omitempty" gorm:"foreignKey:BookingID"`

	PaidAt    *time.Time `json:"paid_at,omitempty"`
	CreatedAt time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName overrides the table name
func (Payment) TableName() string {
	return "payments"
}
