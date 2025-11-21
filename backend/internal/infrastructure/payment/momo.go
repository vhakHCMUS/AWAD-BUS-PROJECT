package payment

import (
	"context"
	"fmt"
	"time"

	"github.com/yourusername/bus-booking/internal/entities"
)

// MoMoGateway handles MoMo e-wallet payments
// Documentation: https://developers.momo.vn/
type MoMoGateway struct {
	PartnerCode string
	AccessKey   string
	SecretKey   string
	Endpoint    string
	UseMock     bool // Set to false when using real API
}

func NewMoMoGateway(partnerCode, accessKey, secretKey, endpoint string, useMock bool) *MoMoGateway {
	return &MoMoGateway{
		PartnerCode: partnerCode,
		AccessKey:   accessKey,
		SecretKey:   secretKey,
		Endpoint:    endpoint,
		UseMock:     useMock,
	}
}

func (g *MoMoGateway) CreatePayment(ctx context.Context, req PaymentRequest) (*PaymentResponse, error) {
	if g.UseMock {
		// Mock implementation
		return &PaymentResponse{
			GatewayPaymentID: fmt.Sprintf("MOMO_%s", req.BookingID.String()[:8]),
			PaymentURL:       fmt.Sprintf("https://test-payment.momo.vn/gw_payment/transactionProcessor?partnerCode=%s&orderId=%s", g.PartnerCode, req.BookingID),
			QRCodeURL:        fmt.Sprintf("https://test-payment.momo.vn/qr/%s", req.BookingID),
			ExpiresAt:        time.Now().Add(15 * time.Minute),
		}, nil
	}

	// Real MoMo API implementation
	// TODO: Replace with actual MoMo API calls
	/*
		requestID := uuid.New().String()
		orderID := req.BookingID.String()
		amount := fmt.Sprintf("%.0f", req.Amount)
		orderInfo := req.Description
		redirectURL := req.ReturnURL
		ipnURL := req.WebhookURL
		extraData := ""

		// Build raw signature
		rawSignature := fmt.Sprintf("accessKey=%s&amount=%s&extraData=%s&ipnUrl=%s&orderId=%s&orderInfo=%s&partnerCode=%s&redirectUrl=%s&requestId=%s&requestType=captureWallet",
			g.AccessKey, amount, extraData, ipnURL, orderID, orderInfo, g.PartnerCode, redirectURL, requestID)

		signature := GenerateSignature(g.SecretKey, rawSignature)

		// Make HTTP request to MoMo endpoint
		payload := map[string]string{
			"partnerCode":  g.PartnerCode,
			"accessKey":    g.AccessKey,
			"requestId":    requestID,
			"amount":       amount,
			"orderId":      orderID,
			"orderInfo":    orderInfo,
			"redirectUrl":  redirectURL,
			"ipnUrl":       ipnURL,
			"extraData":    extraData,
			"requestType":  "captureWallet",
			"signature":    signature,
			"lang":         "vi",
		}

		// Send HTTP POST request and parse response
	*/

	return nil, fmt.Errorf("real MoMo API not yet implemented")
}

func (g *MoMoGateway) VerifyWebhook(ctx context.Context, signature string, payload []byte) error {
	// Verify MoMo webhook signature
	expected := GenerateSignature(g.SecretKey, string(payload))
	if expected != signature {
		return fmt.Errorf("invalid MoMo webhook signature")
	}
	return nil
}

func (g *MoMoGateway) CheckPaymentStatus(ctx context.Context, gatewayPaymentID string) (*PaymentStatus, error) {
	if g.UseMock {
		now := time.Now()
		return &PaymentStatus{
			Status:        entities.PaymentStatusCompleted,
			TransactionID: fmt.Sprintf("MOMO_TXN_%s", gatewayPaymentID),
			PaidAt:        &now,
		}, nil
	}

	// Real API: Query transaction status from MoMo
	return nil, fmt.Errorf("real MoMo status check not yet implemented")
}

func (g *MoMoGateway) RefundPayment(ctx context.Context, gatewayPaymentID string, amount float64) error {
	if g.UseMock {
		return nil
	}

	// Real API: Call MoMo refund endpoint
	return fmt.Errorf("real MoMo refund not yet implemented")
}
