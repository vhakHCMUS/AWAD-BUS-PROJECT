package payment

import (
	"context"
	"fmt"
	"time"

	"github.com/yourusername/bus-booking/internal/entities"
)

// ZaloPayGateway handles ZaloPay e-wallet payments
// Documentation: https://docs.zalopay.vn/
type ZaloPayGateway struct {
	AppID    string
	Key1     string
	Key2     string
	Endpoint string
	UseMock  bool
}

func NewZaloPayGateway(appID, key1, key2, endpoint string, useMock bool) *ZaloPayGateway {
	return &ZaloPayGateway{
		AppID:    appID,
		Key1:     key1,
		Key2:     key2,
		Endpoint: endpoint,
		UseMock:  useMock,
	}
}

func (g *ZaloPayGateway) CreatePayment(ctx context.Context, req PaymentRequest) (*PaymentResponse, error) {
	if g.UseMock {
		return &PaymentResponse{
			GatewayPaymentID: fmt.Sprintf("ZALOPAY_%s", req.BookingID.String()[:8]),
			PaymentURL:       fmt.Sprintf("https://sbgateway.zalopay.vn/api/gateway/pay/%s", req.BookingID),
			QRCodeURL:        fmt.Sprintf("https://sbgateway.zalopay.vn/qr/%s", req.BookingID),
			ExpiresAt:        time.Now().Add(15 * time.Minute),
		}, nil
	}

	// Real ZaloPay API implementation
	// TODO: Replace with actual ZaloPay API calls
	/*
		appTransID := fmt.Sprintf("%s_%d_%s", time.Now().Format("060102"), time.Now().Unix(), req.BookingID.String()[:8])
		embedData := fmt.Sprintf(`{"redirecturl":"%s"}`, req.ReturnURL)

		// Build order data
		orderData := map[string]interface{}{
			"app_id":       g.AppID,
			"app_trans_id": appTransID,
			"app_user":     req.CustomerInfo.Email,
			"app_time":     time.Now().UnixMilli(),
			"amount":       int64(req.Amount),
			"item":         fmt.Sprintf(`[{"itemid":"%s","itemname":"%s","itemprice":%d,"itemquantity":1}]`, req.BookingID, req.Description, int64(req.Amount)),
			"embed_data":   embedData,
			"bank_code":    "zalopayapp",
			"callback_url": req.WebhookURL,
		}

		// Generate MAC signature
		data := fmt.Sprintf("%s|%s|%s|%d|%s|%s|%s",
			g.AppID, appTransID, req.CustomerInfo.Email, int64(req.Amount),
			time.Now().UnixMilli(), embedData, "[]")
		mac := GenerateSignature(g.Key1, data)
		orderData["mac"] = mac

		// Send HTTP POST request to ZaloPay create order endpoint
	*/

	return nil, fmt.Errorf("real ZaloPay API not yet implemented")
}

func (g *ZaloPayGateway) VerifyWebhook(ctx context.Context, signature string, payload []byte) error {
	expected := GenerateSignature(g.Key2, string(payload))
	if expected != signature {
		return fmt.Errorf("invalid ZaloPay webhook signature")
	}
	return nil
}

func (g *ZaloPayGateway) CheckPaymentStatus(ctx context.Context, gatewayPaymentID string) (*PaymentStatus, error) {
	if g.UseMock {
		now := time.Now()
		return &PaymentStatus{
			Status:        entities.PaymentStatusCompleted,
			TransactionID: fmt.Sprintf("ZALOPAY_TXN_%s", gatewayPaymentID),
			PaidAt:        &now,
		}, nil
	}

	// Real API: Query transaction status
	return nil, fmt.Errorf("real ZaloPay status check not yet implemented")
}

func (g *ZaloPayGateway) RefundPayment(ctx context.Context, gatewayPaymentID string, amount float64) error {
	if g.UseMock {
		return nil
	}

	// Real API: Call ZaloPay refund endpoint
	return fmt.Errorf("real ZaloPay refund not yet implemented")
}
