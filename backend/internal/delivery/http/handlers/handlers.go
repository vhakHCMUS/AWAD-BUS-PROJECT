package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/yourusername/bus-booking/internal/entities"
	"github.com/yourusername/bus-booking/internal/usecases"
)

type TripHandler struct {
	tripUsecase    *usecases.TripUsecase
	bookingUsecase *usecases.BookingUsecase
}

func NewTripHandler(tripUsecase *usecases.TripUsecase, bookingUsecase *usecases.BookingUsecase) *TripHandler {
	return &TripHandler{
		tripUsecase:    tripUsecase,
		bookingUsecase: bookingUsecase,
	}
}

type SearchTripsRequest struct {
	FromCity string `form:"from_city" binding:"required"`
	ToCity   string `form:"to_city" binding:"required"`
	Date     string `form:"date" binding:"required"` // YYYY-MM-DD format
	Page     int    `form:"page"`
	Limit    int    `form:"limit"`
}

func (h *TripHandler) Search(c *gin.Context) {
	var req SearchTripsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	// Parse date
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid date format, use YYYY-MM-DD"})
		return
	}

	// Set defaults
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 || req.Limit > 100 {
		req.Limit = 20
	}

	trips, err := h.tripUsecase.SearchTrips(c.Request.Context(), req.FromCity, req.ToCity, date, req.Page, req.Limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to search trips"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"trips": trips})
}

func (h *TripHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid trip ID"})
		return
	}

	trip, err := h.tripUsecase.GetTripByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Trip not found"})
		return
	}

	c.JSON(http.StatusOK, trip)
}

func (h *TripHandler) GetSeats(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid trip ID"})
		return
	}

	seats, err := h.bookingUsecase.GetAvailableSeats(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to get seats"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"seats": seats})
}

type CreateTripRequest struct {
	BusID         string  `json:"bus_id" binding:"required"`
	RouteID       string  `json:"route_id" binding:"required"`
	DepartureTime string  `json:"departure_time" binding:"required"`
	ArrivalTime   string  `json:"arrival_time" binding:"required"`
	Price         float64 `json:"price" binding:"required,gt=0"`
	DriverName    string  `json:"driver_name"`
	DriverPhone   string  `json:"driver_phone"`
}

func (h *TripHandler) Create(c *gin.Context) {
	var req CreateTripRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	busID, err := uuid.Parse(req.BusID)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid bus ID"})
		return
	}

	routeID, err := uuid.Parse(req.RouteID)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid route ID"})
		return
	}

	departureTime, err := time.Parse(time.RFC3339, req.DepartureTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid departure time format"})
		return
	}

	arrivalTime, err := time.Parse(time.RFC3339, req.ArrivalTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid arrival time format"})
		return
	}

	trip := &entities.Trip{
		BusID:         busID,
		RouteID:       routeID,
		DepartureTime: departureTime,
		ArrivalTime:   arrivalTime,
		Price:         req.Price,
		DriverName:    req.DriverName,
		DriverPhone:   req.DriverPhone,
		Status:        entities.TripStatusScheduled,
	}

	err = h.tripUsecase.CreateTrip(c.Request.Context(), trip)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, trip)
}

type UpdateTripRequest struct {
	DepartureTime string              `json:"departure_time"`
	ArrivalTime   string              `json:"arrival_time"`
	Price         float64             `json:"price,omitempty"`
	Status        entities.TripStatus `json:"status,omitempty"`
	DriverName    string              `json:"driver_name"`
	DriverPhone   string              `json:"driver_phone"`
}

func (h *TripHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid trip ID"})
		return
	}

	var req UpdateTripRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	// Get existing trip
	trip, err := h.tripUsecase.GetTripByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Trip not found"})
		return
	}

	// Update fields
	if req.DepartureTime != "" {
		departureTime, err := time.Parse(time.RFC3339, req.DepartureTime)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid departure time format"})
			return
		}
		trip.DepartureTime = departureTime
	}

	if req.ArrivalTime != "" {
		arrivalTime, err := time.Parse(time.RFC3339, req.ArrivalTime)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid arrival time format"})
			return
		}
		trip.ArrivalTime = arrivalTime
	}

	if req.Price > 0 {
		trip.Price = req.Price
	}

	if req.Status != "" {
		trip.Status = req.Status
	}

	if req.DriverName != "" {
		trip.DriverName = req.DriverName
	}

	if req.DriverPhone != "" {
		trip.DriverPhone = req.DriverPhone
	}

	err = h.tripUsecase.UpdateTrip(c.Request.Context(), trip)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, trip)
}

func (h *TripHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid trip ID"})
		return
	}

	err = h.tripUsecase.DeleteTrip(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{Message: "Trip deleted successfully"})
}

type TicketHandler struct{}

func NewTicketHandler() *TicketHandler {
	return &TicketHandler{}
}

func (h *TicketHandler) GetTicket(c *gin.Context) {
	c.JSON(200, gin.H{"ticket": nil})
}

func (h *TicketHandler) CheckIn(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Checked in"})
}

type BusHandler struct {
	busUsecase *usecases.BusUsecase
}

func NewBusHandler(busUsecase *usecases.BusUsecase) *BusHandler {
	return &BusHandler{busUsecase: busUsecase}
}

func (h *BusHandler) Create(c *gin.Context) {
	var bus entities.Bus
	if err := c.ShouldBindJSON(&bus); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	err := h.busUsecase.CreateBus(c.Request.Context(), &bus)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, bus)
}

func (h *BusHandler) List(c *gin.Context) {
	page := 1
	limit := 20
	if pageStr := c.Query("page"); pageStr != "" {
		if p, err := parsePositiveInt(pageStr); err == nil {
			page = p
		}
	}
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := parsePositiveInt(limitStr); err == nil && l <= 100 {
			limit = l
		}
	}

	status := entities.BusStatus(c.Query("status"))

	buses, err := h.busUsecase.ListBuses(c.Request.Context(), status, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to list buses"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"buses": buses})
}

func (h *BusHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid bus ID"})
		return
	}

	bus, err := h.busUsecase.GetBusByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Bus not found"})
		return
	}

	c.JSON(http.StatusOK, bus)
}

func (h *BusHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid bus ID"})
		return
	}

	var bus entities.Bus
	if err := c.ShouldBindJSON(&bus); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	bus.ID = id
	err = h.busUsecase.UpdateBus(c.Request.Context(), &bus)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, bus)
}

func (h *BusHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid bus ID"})
		return
	}

	err = h.busUsecase.DeleteBus(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{Message: "Bus deleted successfully"})
}

type RouteHandler struct {
	routeUsecase *usecases.RouteUsecase
}

func NewRouteHandler(routeUsecase *usecases.RouteUsecase) *RouteHandler {
	return &RouteHandler{routeUsecase: routeUsecase}
}

func (h *RouteHandler) Create(c *gin.Context) {
	var route entities.Route
	if err := c.ShouldBindJSON(&route); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	err := h.routeUsecase.CreateRoute(c.Request.Context(), &route)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, route)
}

func (h *RouteHandler) List(c *gin.Context) {
	page := 1
	limit := 20
	if pageStr := c.Query("page"); pageStr != "" {
		if p, err := parsePositiveInt(pageStr); err == nil {
			page = p
		}
	}
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := parsePositiveInt(limitStr); err == nil && l <= 100 {
			limit = l
		}
	}

	routes, err := h.routeUsecase.ListRoutes(c.Request.Context(), page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to list routes"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"routes": routes})
}

func (h *RouteHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid route ID"})
		return
	}

	route, err := h.routeUsecase.GetRouteByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Route not found"})
		return
	}

	c.JSON(http.StatusOK, route)
}

func (h *RouteHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid route ID"})
		return
	}

	var route entities.Route
	if err := c.ShouldBindJSON(&route); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	route.ID = id
	err = h.routeUsecase.UpdateRoute(c.Request.Context(), &route)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, route)
}

func (h *RouteHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid route ID"})
		return
	}

	err = h.routeUsecase.DeleteRoute(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{Message: "Route deleted successfully"})
}

func parsePositiveInt(s string) (int, error) {
	var i int
	_, err := fmt.Sscanf(s, "%d", &i)
	if err != nil || i <= 0 {
		return 0, fmt.Errorf("invalid positive integer")
	}
	return i, nil
}
