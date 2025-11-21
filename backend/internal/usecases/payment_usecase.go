package usecases

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/yourusername/bus-booking/internal/entities"
	"github.com/yourusername/bus-booking/internal/infrastructure/payment"
	"github.com/yourusername/bus-booking/internal/repositories"
)

type PaymentUsecase struct {
	paymentRepo repositories.PaymentRepository
	bookingRepo repositories.BookingRepository
	gateways    map[entities.PaymentGateway]payment.Gateway
}

func NewPaymentUsecase(
	paymentRepo repositories.PaymentRepository,
	bookingRepo repositories.BookingRepository,
	gateways map[entities.PaymentGateway]payment.Gateway,
) *PaymentUsecase {
	return &PaymentUsecase{
		paymentRepo: paymentRepo,
		bookingRepo: bookingRepo,
		gateways:    gateways,
	}
}

// CreatePayment initiates a payment for a booking
func (uc *PaymentUsecase) CreatePayment(ctx context.Context, bookingID uuid.UUID, gateway entities.PaymentGateway) (*entities.Payment, string, error) {
	// Get booking details
	booking, err := uc.bookingRepo.GetByIDWithDetails(ctx, bookingID)
	if err != nil {
		return nil, "", fmt.Errorf("booking not found: %w", err)
	}

	if booking.Status != entities.BookingStatusPending {
		return nil, "", fmt.Errorf("booking is not in pending status")
	}

	// Check if payment already exists
	existingPayment, _ := uc.paymentRepo.GetByBookingID(ctx, bookingID)
	if existingPayment != nil && existingPayment.Status != entities.PaymentStatusFailed {
		return nil, "", fmt.Errorf("payment already exists for this booking")
	}

	// Get payment gateway
	gw, ok := uc.gateways[gateway]
	if !ok {
		return nil, "", fmt.Errorf("unsupported payment gateway: %s", gateway)
	}

	// Generate idempotency key
	idempotencyKey := fmt.Sprintf("%s_%s_%d", bookingID, gateway, time.Now().Unix())

	// Create payment request
	req := payment.PaymentRequest{
		BookingID:   bookingID,
		Amount:      booking.TotalPrice,
		Currency:    "VND",
		Description: fmt.Sprintf("Thanh toán vé xe - %s", booking.BookingCode),
		ReturnURL:   fmt.Sprintf("https://vietbusbooking.com/booking/%s/payment-success", bookingID),
		CancelURL:   fmt.Sprintf("https://vietbusbooking.com/booking/%s/payment-cancel", bookingID),
		WebhookURL:  "https://api.vietbusbooking.com/api/v1/payments/webhook",
		CustomerInfo: payment.CustomerInfo{
			Name:  booking.ContactName,
			Email: booking.ContactEmail,
			Phone: booking.ContactPhone,
		},
	}

	// Call gateway to create payment
	resp, err := gw.CreatePayment(ctx, req)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create payment: %w", err)
	}

	// Create payment record
	pmt := &entities.Payment{
		BookingID:        bookingID,
		Gateway:          gateway,
		GatewayPaymentID: resp.GatewayPaymentID,
		Amount:           booking.TotalPrice,
		Currency:         "VND",
		Status:           entities.PaymentStatusPending,
		IdempotencyKey:   idempotencyKey,
	}

	if err := uc.paymentRepo.Create(ctx, pmt); err != nil {
		return nil, "", fmt.Errorf("failed to save payment: %w", err)
	}

	return pmt, resp.PaymentURL, nil
}

// HandleWebhook processes payment webhook notifications
func (uc *PaymentUsecase) HandleWebhook(ctx context.Context, gateway entities.PaymentGateway, signature string, payload []byte) error {
	// Get gateway
	gw, ok := uc.gateways[gateway]
	if !ok {
		return fmt.Errorf("unsupported gateway: %s", gateway)
	}

	// Verify webhook signature
	if err := gw.VerifyWebhook(ctx, signature, payload); err != nil {
		return fmt.Errorf("invalid webhook signature: %w", err)
	}

	// Parse webhook payload
	webhookData, err := payment.ParseWebhookPayload(payload)
	if err != nil {
		return fmt.Errorf("failed to parse webhook: %w", err)
	}

	// Get payment by gateway payment ID
	pmt, err := uc.paymentRepo.GetByGatewayPaymentID(ctx, webhookData.GatewayPaymentID)
	if err != nil {
		return fmt.Errorf("payment not found: %w", err)
	}

	// Check idempotency - if already processed, return success
	if pmt.Status == entities.PaymentStatusCompleted {
		return nil
	}

	// Update payment status
	switch webhookData.Status {
	case "success", "completed":
		pmt.Status = entities.PaymentStatusCompleted
		pmt.TransactionID = webhookData.TransactionID
		paidAt := webhookData.PaidAt
		pmt.PaidAt = &paidAt
	case "failed", "error":
		pmt.Status = entities.PaymentStatusFailed
		failureReason := "Payment failed"
		pmt.FailureReason = &failureReason
	case "cancelled":
		pmt.Status = entities.PaymentStatusCancelled
	default:
		return fmt.Errorf("unknown payment status: %s", webhookData.Status)
	}

	if err := uc.paymentRepo.Update(ctx, pmt); err != nil {
		return fmt.Errorf("failed to update payment: %w", err)
	}

	// Update booking status if payment succeeded
	if pmt.Status == entities.PaymentStatusCompleted {
		if err := uc.bookingRepo.UpdateStatus(ctx, pmt.BookingID, entities.BookingStatusPaid); err != nil {
			return fmt.Errorf("failed to update booking: %w", err)
		}
	}

	return nil
}

// CheckPaymentStatus checks the current status of a payment
func (uc *PaymentUsecase) CheckPaymentStatus(ctx context.Context, paymentID uuid.UUID) (*entities.Payment, error) {
	pmt, err := uc.paymentRepo.GetByID(ctx, paymentID)
	if err != nil {
		return nil, err
	}

	// If payment is not pending, return current status
	if pmt.Status != entities.PaymentStatusPending {
		return pmt, nil
	}

	// Query gateway for latest status
	gw, ok := uc.gateways[pmt.Gateway]
	if !ok {
		return pmt, nil
	}

	status, err := gw.CheckPaymentStatus(ctx, pmt.GatewayPaymentID)
	if err != nil {
		return pmt, nil // Return current status if gateway check fails
	}

	// Update if status changed
	if status.Status != pmt.Status {
		pmt.Status = status.Status
		pmt.TransactionID = status.TransactionID
		pmt.PaidAt = status.PaidAt
		uc.paymentRepo.Update(ctx, pmt)
	}

	return pmt, nil
}

// RefundPayment processes a refund for a payment
func (uc *PaymentUsecase) RefundPayment(ctx context.Context, paymentID uuid.UUID) error {
	pmt, err := uc.paymentRepo.GetByID(ctx, paymentID)
	if err != nil {
		return err
	}

	if pmt.Status != entities.PaymentStatusCompleted {
		return fmt.Errorf("payment must be completed to refund")
	}

	// Call gateway to refund
	gw, ok := uc.gateways[pmt.Gateway]
	if !ok {
		return fmt.Errorf("unsupported gateway")
	}

	if err := gw.RefundPayment(ctx, pmt.GatewayPaymentID, pmt.Amount); err != nil {
		return fmt.Errorf("refund failed: %w", err)
	}

	// Update payment status
	pmt.Status = entities.PaymentStatusRefunded
	if err := uc.paymentRepo.Update(ctx, pmt); err != nil {
		return err
	}

	// Update booking status
	if err := uc.bookingRepo.UpdateStatus(ctx, pmt.BookingID, entities.BookingStatusRefunded); err != nil {
		return err
	}

	return nil
}
