package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/yourusername/bus-booking/internal/usecases"
)

type BookingHandler struct {
	usecase *usecases.BookingUsecase
}

func NewBookingHandler(usecase *usecases.BookingUsecase) *BookingHandler {
	return &BookingHandler{usecase: usecase}
}

type InitiateBookingRequest struct {
	TripID       string   `json:"trip_id" binding:"required"`
	SeatNumbers  []string `json:"seat_numbers" binding:"required,min=1"`
	ContactName  string   `json:"contact_name" binding:"required"`
	ContactEmail string   `json:"contact_email" binding:"required,email"`
	ContactPhone string   `json:"contact_phone" binding:"required"`
}

// InitiateBooking godoc
// @Summary Initiate a new booking
// @Description Create a pending booking and lock seats
// @Tags bookings
// @Accept json
// @Produce json
// @Param request body InitiateBookingRequest true "Booking details"
// @Success 201 {object} entities.Booking
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse "Seat already locked or booked"
// @Security BearerAuth
// @Router /bookings [post]
func (h *BookingHandler) InitiateBooking(c *gin.Context) {
	var req InitiateBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	tripID, err := uuid.Parse(req.TripID)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid trip ID"})
		return
	}

	// Get user ID from context (set by auth middleware)
	var userID *uuid.UUID
	if uid, exists := c.Get("user_id"); exists {
		if parsed, err := uuid.Parse(uid.(string)); err == nil {
			userID = &parsed
		}
	}

	booking, err := h.usecase.InitiateBooking(c.Request.Context(), tripID, req.SeatNumbers, userID)
	if err != nil {
		c.JSON(http.StatusConflict, ErrorResponse{Error: err.Error()})
		return
	}

	// Update contact details
	booking.ContactName = req.ContactName
	booking.ContactEmail = req.ContactEmail
	booking.ContactPhone = req.ContactPhone

	c.JSON(http.StatusCreated, booking)
}

// GetBooking godoc
// @Summary Get booking details
// @Description Retrieve booking by ID with all details
// @Tags bookings
// @Produce json
// @Param id path string true "Booking ID"
// @Success 200 {object} entities.Booking
// @Failure 404 {object} ErrorResponse
// @Security BearerAuth
// @Router /bookings/{id} [get]
func (h *BookingHandler) GetBooking(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid booking ID"})
		return
	}

	booking, err := h.usecase.GetBookingByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Booking not found"})
		return
	}

	c.JSON(http.StatusOK, booking)
}

// GetUserBookings godoc
// @Summary Get user's bookings
// @Description Retrieve all bookings for the authenticated user
// @Tags bookings
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {array} entities.Booking
// @Security BearerAuth
// @Router /bookings [get]
func (h *BookingHandler) GetUserBookings(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userIDValue, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "User not authenticated"})
		return
	}

	userID, ok := userIDValue.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Invalid user ID"})
		return
	}

	// Pagination parameters
	page := 1
	limit := 10
	if pageStr := c.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	bookings, err := h.usecase.GetUserBookings(c.Request.Context(), userID, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to retrieve bookings"})
		return
	}

	c.JSON(http.StatusOK, bookings)
}

// CancelBooking godoc
// @Summary Cancel a booking
// @Description Cancel a booking and release seats
// @Tags bookings
// @Param id path string true "Booking ID"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Security BearerAuth
// @Router /bookings/{id}/cancel [post]
func (h *BookingHandler) CancelBooking(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid booking ID"})
		return
	}

	err = h.usecase.CancelBooking(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{Message: "Booking cancelled successfully"})
}

// Common response types
type ErrorResponse struct {
	Error string `json:"error"`
}

type SuccessResponse struct {
	Message string `json:"message"`
}
