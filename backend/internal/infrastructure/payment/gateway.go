package payment

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/yourusername/bus-booking/internal/entities"
)

// Gateway interface defines common payment operations
type Gateway interface {
	CreatePayment(ctx context.Context, req PaymentRequest) (*PaymentResponse, error)
	VerifyWebhook(ctx context.Context, signature string, payload []byte) error
	CheckPaymentStatus(ctx context.Context, gatewayPaymentID string) (*PaymentStatus, error)
	RefundPayment(ctx context.Context, gatewayPaymentID string, amount float64) error
}

// PaymentRequest represents a payment creation request
type PaymentRequest struct {
	BookingID    uuid.UUID
	Amount       float64
	Currency     string
	Description  string
	ReturnURL    string
	CancelURL    string
	WebhookURL   string
	CustomerInfo CustomerInfo
}

// CustomerInfo represents customer details
type CustomerInfo struct {
	Name  string
	Email string
	Phone string
}

// PaymentResponse represents the gateway's response
type PaymentResponse struct {
	GatewayPaymentID string
	PaymentURL       string
	QRCodeURL        string
	ExpiresAt        time.Time
}

// PaymentStatus represents payment status check result
type PaymentStatus struct {
	Status        entities.PaymentStatus
	TransactionID string
	PaidAt        *time.Time
	FailureReason string
}

// WebhookPayload represents a generic webhook payload
type WebhookPayload struct {
	GatewayPaymentID string
	Status           string
	TransactionID    string
	Amount           float64
	PaidAt           time.Time
}

// GenerateSignature creates HMAC-SHA256 signature
func GenerateSignature(secret string, data string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

// VerifySignature verifies HMAC-SHA256 signature
func VerifySignature(secret string, data string, signature string) bool {
	expected := GenerateSignature(secret, data)
	return hmac.Equal([]byte(expected), []byte(signature))
}

// MockGateway is a simulated payment gateway for testing
// Can be easily replaced with real implementations
type MockGateway struct {
	Name      string
	SecretKey string
}

func NewMockGateway(name, secretKey string) *MockGateway {
	return &MockGateway{
		Name:      name,
		SecretKey: secretKey,
	}
}

func (g *MockGateway) CreatePayment(ctx context.Context, req PaymentRequest) (*PaymentResponse, error) {
	// Simulate payment creation
	gatewayPaymentID := fmt.Sprintf("%s_%s", g.Name, uuid.New().String()[:8])

	return &PaymentResponse{
		GatewayPaymentID: gatewayPaymentID,
		PaymentURL:       fmt.Sprintf("https://mock-%s.vn/pay/%s", g.Name, gatewayPaymentID),
		QRCodeURL:        fmt.Sprintf("https://mock-%s.vn/qr/%s", g.Name, gatewayPaymentID),
		ExpiresAt:        time.Now().Add(15 * time.Minute),
	}, nil
}

func (g *MockGateway) VerifyWebhook(ctx context.Context, signature string, payload []byte) error {
	// Simulate webhook verification
	expected := GenerateSignature(g.SecretKey, string(payload))
	if !hmac.Equal([]byte(expected), []byte(signature)) {
		return fmt.Errorf("invalid webhook signature")
	}
	return nil
}

func (g *MockGateway) CheckPaymentStatus(ctx context.Context, gatewayPaymentID string) (*PaymentStatus, error) {
	// Simulate status check - in real implementation, call gateway API
	now := time.Now()
	return &PaymentStatus{
		Status:        entities.PaymentStatusCompleted,
		TransactionID: fmt.Sprintf("TXN_%s", gatewayPaymentID),
		PaidAt:        &now,
	}, nil
}

func (g *MockGateway) RefundPayment(ctx context.Context, gatewayPaymentID string, amount float64) error {
	// Simulate refund - in real implementation, call gateway API
	return nil
}

// ParseWebhookPayload parses webhook JSON payload
func ParseWebhookPayload(payload []byte) (*WebhookPayload, error) {
	var wp WebhookPayload
	if err := json.Unmarshal(payload, &wp); err != nil {
		return nil, err
	}
	return &wp, nil
}
