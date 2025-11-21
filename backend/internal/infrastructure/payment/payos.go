package payment

import (
	"context"
	"fmt"
	"time"

	"github.com/yourusername/bus-booking/internal/entities"
)

// PayOSGateway handles PayOS payments (ATM/Visa/Mastercard)
// Documentation: https://payos.vn/docs/
type PayOSGateway struct {
	ClientID    string
	APIKey      string
	ChecksumKey string
	Endpoint    string
	UseMock     bool
}

func NewPayOSGateway(clientID, apiKey, checksumKey, endpoint string, useMock bool) *PayOSGateway {
	return &PayOSGateway{
		ClientID:    clientID,
		APIKey:      apiKey,
		ChecksumKey: checksumKey,
		Endpoint:    endpoint,
		UseMock:     useMock,
	}
}

func (g *PayOSGateway) CreatePayment(ctx context.Context, req PaymentRequest) (*PaymentResponse, error) {
	if g.UseMock {
		return &PaymentResponse{
			GatewayPaymentID: fmt.Sprintf("PAYOS_%s", req.BookingID.String()[:8]),
			PaymentURL:       fmt.Sprintf("https://test.payos.vn/checkout/%s", req.BookingID),
			QRCodeURL:        fmt.Sprintf("https://test.payos.vn/qr/%s", req.BookingID),
			ExpiresAt:        time.Now().Add(15 * time.Minute),
		}, nil
	}

	// Real PayOS API implementation
	// TODO: Replace with actual PayOS API calls
	/*
		orderCode := int64(time.Now().Unix())

		paymentData := map[string]interface{}{
			"orderCode":   orderCode,
			"amount":      int(req.Amount),
			"description": req.Description,
			"cancelUrl":   req.CancelURL,
			"returnUrl":   req.ReturnURL,
			"buyerName":   req.CustomerInfo.Name,
			"buyerEmail":  req.CustomerInfo.Email,
			"buyerPhone":  req.CustomerInfo.Phone,
		}

		// Generate checksum signature
		dataStr := fmt.Sprintf("%d|%d|%s|%s|%s",
			orderCode, int(req.Amount), req.Description, req.CancelURL, req.ReturnURL)
		checksum := GenerateSignature(g.ChecksumKey, dataStr)
		paymentData["signature"] = checksum

		// Set authorization header with API key
		headers := map[string]string{
			"x-client-id": g.ClientID,
			"x-api-key":   g.APIKey,
		}

		// Send HTTP POST request to PayOS create payment link endpoint
	*/

	return nil, fmt.Errorf("real PayOS API not yet implemented")
}

func (g *PayOSGateway) VerifyWebhook(ctx context.Context, signature string, payload []byte) error {
	expected := GenerateSignature(g.ChecksumKey, string(payload))
	if expected != signature {
		return fmt.Errorf("invalid PayOS webhook signature")
	}
	return nil
}

func (g *PayOSGateway) CheckPaymentStatus(ctx context.Context, gatewayPaymentID string) (*PaymentStatus, error) {
	if g.UseMock {
		now := time.Now()
		return &PaymentStatus{
			Status:        entities.PaymentStatusCompleted,
			TransactionID: fmt.Sprintf("PAYOS_TXN_%s", gatewayPaymentID),
			PaidAt:        &now,
		}, nil
	}

	// Real API: Query payment status
	return nil, fmt.Errorf("real PayOS status check not yet implemented")
}

func (g *PayOSGateway) RefundPayment(ctx context.Context, gatewayPaymentID string, amount float64) error {
	if g.UseMock {
		return nil
	}

	// Real API: Call PayOS cancel/refund endpoint
	return fmt.Errorf("real PayOS refund not yet implemented")
}
