package handlers

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/yourusername/bus-booking/internal/entities"
	"github.com/yourusername/bus-booking/internal/usecases"
)

type PaymentHandler struct {
	paymentUsecase *usecases.PaymentUsecase
}

func NewPaymentHandler(paymentUsecase *usecases.PaymentUsecase) *PaymentHandler {
	return &PaymentHandler{
		paymentUsecase: paymentUsecase,
	}
}

// CreatePayment godoc
// @Summary Create a payment for booking
// @Description Initialize payment with selected gateway (MoMo, ZaloPay, PayOS)
// @Tags payments
// @Accept json
// @Produce json
// @Param request body CreatePaymentRequest true "Payment request"
// @Success 200 {object} CreatePaymentResponse
// @Failure 400 {object} ErrorResponse
// @Security BearerAuth
// @Router /payments [post]
func (h *PaymentHandler) CreatePayment(c *gin.Context) {
	var req CreatePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	bookingID, err := uuid.Parse(req.BookingID)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid booking ID"})
		return
	}

	// Validate gateway
	var gateway entities.PaymentGateway
	switch req.Gateway {
	case "momo":
		gateway = entities.PaymentGatewayMoMo
	case "zalopay":
		gateway = "zalopay"
	case "payos":
		gateway = entities.PaymentGatewayPayOS
	default:
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid payment gateway"})
		return
	}

	payment, paymentURL, err := h.paymentUsecase.CreatePayment(c.Request.Context(), bookingID, gateway)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, CreatePaymentResponse{
		PaymentID:  payment.ID.String(),
		PaymentURL: paymentURL,
		ExpiresIn:  900, // 15 minutes
	})
}

// WebhookMoMo godoc
// @Summary MoMo payment webhook
// @Description Handle MoMo payment notification webhook
// @Tags payments
// @Accept json
// @Produce json
// @Success 200 {object} WebhookResponse
// @Router /payments/webhook/momo [post]
func (h *PaymentHandler) WebhookMoMo(c *gin.Context) {
	signature := c.GetHeader("X-Signature")
	payload, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}

	err = h.paymentUsecase.HandleWebhook(c.Request.Context(), entities.PaymentGatewayMoMo, signature, payload)
	if err != nil {
		c.JSON(http.StatusBadRequest, WebhookResponse{Status: "error", Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, WebhookResponse{Status: "success"})
}

// WebhookZaloPay godoc
// @Summary ZaloPay payment webhook
// @Description Handle ZaloPay payment notification webhook
// @Tags payments
// @Accept json
// @Produce json
// @Success 200 {object} WebhookResponse
// @Router /payments/webhook/zalopay [post]
func (h *PaymentHandler) WebhookZaloPay(c *gin.Context) {
	signature := c.GetHeader("X-Mac")
	payload, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}

	err = h.paymentUsecase.HandleWebhook(c.Request.Context(), "zalopay", signature, payload)
	if err != nil {
		c.JSON(http.StatusBadRequest, WebhookResponse{Status: "error", Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, WebhookResponse{Status: "success"})
}

// WebhookPayOS godoc
// @Summary PayOS payment webhook
// @Description Handle PayOS payment notification webhook
// @Tags payments
// @Accept json
// @Produce json
// @Success 200 {object} WebhookResponse
// @Router /payments/webhook/payos [post]
func (h *PaymentHandler) WebhookPayOS(c *gin.Context) {
	signature := c.GetHeader("X-Signature")
	payload, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}

	err = h.paymentUsecase.HandleWebhook(c.Request.Context(), entities.PaymentGatewayPayOS, signature, payload)
	if err != nil {
		c.JSON(http.StatusBadRequest, WebhookResponse{Status: "error", Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, WebhookResponse{Status: "success"})
}

// GetPaymentStatus godoc
// @Summary Get payment status
// @Description Check current status of a payment
// @Tags payments
// @Produce json
// @Param id path string true "Payment ID"
// @Success 200 {object} entities.Payment
// @Failure 404 {object} ErrorResponse
// @Security BearerAuth
// @Router /payments/{id} [get]
func (h *PaymentHandler) GetPaymentStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid payment ID"})
		return
	}

	payment, err := h.paymentUsecase.CheckPaymentStatus(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Payment not found"})
		return
	}

	c.JSON(http.StatusOK, payment)
}

// Request/Response types
type CreatePaymentRequest struct {
	BookingID string `json:"booking_id" binding:"required"`
	Gateway   string `json:"gateway" binding:"required,oneof=momo zalopay payos"`
}

type CreatePaymentResponse struct {
	PaymentID  string `json:"payment_id"`
	PaymentURL string `json:"payment_url"`
	ExpiresIn  int    `json:"expires_in"` // seconds
}

type WebhookResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}
