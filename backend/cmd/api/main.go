package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	postgresDriver "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/yourusername/bus-booking/internal/delivery/http/handlers"
	"github.com/yourusername/bus-booking/internal/delivery/http/middleware"
	"github.com/yourusername/bus-booking/internal/delivery/websocket"
	"github.com/yourusername/bus-booking/internal/entities"
	"github.com/yourusername/bus-booking/internal/infrastructure"
	"github.com/yourusername/bus-booking/internal/repositories"
	"github.com/yourusername/bus-booking/internal/repositories/cache"
	"github.com/yourusername/bus-booking/internal/repositories/postgres"
	"github.com/yourusername/bus-booking/internal/usecases"
)

// @title Bus Booking API
// @version 1.0
// @description Production-grade bus ticket booking system with real-time seat selection
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@busbooking.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1
// @schemes http https

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Initialize database
	db, err := initDatabase()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Initialize Redis
	redisClient := initRedis()

	// Run migrations
	if err := runMigrations(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize dependencies
	container := initDependencies(db, redisClient)

	// Setup Gin router
	router := setupRouter(container)

	// Start background jobs
	go startBackgroundJobs(container)

	// Start WebSocket server
	go websocket.StartWebSocketServer(container.RedisCache)

	// Start HTTP server
	port := getEnv("PORT", "8080")
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown
	go func() {
		log.Printf("Server starting on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exited")
}

func initDatabase() (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_PORT", "5432"),
		getEnv("DB_USER", "postgres"),
		getEnv("DB_PASSWORD", "postgres"),
		getEnv("DB_NAME", "bus_booking"),
		getEnv("DB_SSL_MODE", "disable"),
	)

	db, err := gorm.Open(postgresDriver.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db, nil
}

func initRedis() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", getEnv("REDIS_HOST", "localhost"), getEnv("REDIS_PORT", "6379")),
		Password: getEnv("REDIS_PASSWORD", ""),
		DB:       0,
	})
}

func runMigrations(db *gorm.DB) error {
	return db.AutoMigrate(
		&entities.User{},
		&entities.Bus{},
		&entities.Route{},
		&entities.Trip{},
		&entities.SeatInfo{},
		&entities.Booking{},
		&entities.Payment{},
		&entities.Ticket{},
		&entities.RefreshToken{},
	)
}

type Container struct {
	// Repositories (using interfaces)
	UserRepo    repositories.UserRepository
	BookingRepo repositories.BookingRepository
	SeatRepo    repositories.SeatRepository
	RedisCache  *cache.RedisCache

	// Usecases
	BookingUsecase *usecases.BookingUsecase

	// Infrastructure
	EmailService *infrastructure.EmailService
	PDFGenerator *infrastructure.PDFGenerator
}

func initDependencies(db *gorm.DB, redisClient *redis.Client) *Container {
	// Repositories
	userRepo := postgres.NewUserRepository(db)
	bookingRepo := postgres.NewBookingRepository(db)
	seatRepo := postgres.NewSeatRepository(db)
	tripRepo := postgres.NewTripRepository(db)
	paymentRepo := postgres.NewPaymentRepository(db)
	ticketRepo := postgres.NewTicketRepository(db)

	// Cache
	redisCache := cache.NewRedisCache(redisClient)

	// Infrastructure
	emailService := infrastructure.NewEmailService()
	pdfGenerator := infrastructure.NewPDFGenerator()

	// Usecases
	seatLockDuration, _ := time.ParseDuration(getEnv("SEAT_LOCK_DURATION", "10m"))
	bookingExpiry, _ := time.ParseDuration(getEnv("BOOKING_EXPIRY", "15m"))

	bookingUsecase := usecases.NewBookingUsecase(
		bookingRepo,
		seatRepo,
		tripRepo,
		paymentRepo,
		ticketRepo,
		redisCache,
		seatLockDuration,
		bookingExpiry,
	)

	return &Container{
		UserRepo:       userRepo,
		BookingRepo:    bookingRepo,
		SeatRepo:       seatRepo,
		RedisCache:     redisCache,
		BookingUsecase: bookingUsecase,
		EmailService:   emailService,
		PDFGenerator:   pdfGenerator,
	}
}

func setupRouter(container *Container) *gin.Engine {
	if getEnv("ENV", "development") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middleware.CORS())

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API v1
	v1 := router.Group("/api/v1")
	{
		// Public routes
		auth := v1.Group("/auth")
		{
			authHandler := handlers.NewAuthHandler()
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.RefreshToken)
			auth.GET("/google", authHandler.GoogleLogin)
			auth.GET("/google/callback", authHandler.GoogleCallback)
			auth.GET("/github", authHandler.GitHubLogin)
			auth.GET("/github/callback", authHandler.GitHubCallback)
		}

		// Trip search (public)
		trips := v1.Group("/trips")
		{
			tripHandler := handlers.NewTripHandler()
			trips.GET("", tripHandler.Search)
			trips.GET("/:id", tripHandler.GetByID)
			trips.GET("/:id/seats", tripHandler.GetSeats)
		}

		// Protected routes
		authorized := v1.Group("")
		authorized.Use(middleware.AuthMiddleware())
		{
			// Bookings
			bookings := authorized.Group("/bookings")
			{
				bookingHandler := handlers.NewBookingHandler(container.BookingUsecase)
				bookings.POST("", bookingHandler.InitiateBooking)
				bookings.GET("/:id", bookingHandler.GetBooking)
				bookings.GET("", bookingHandler.GetUserBookings)
				bookings.POST("/:id/cancel", bookingHandler.CancelBooking)
			}

			// Payments
			payments := authorized.Group("/payments")
			{
				paymentHandler := handlers.NewPaymentHandler()
				payments.POST("", paymentHandler.CreatePayment)
				payments.POST("/webhook", paymentHandler.Webhook)
			}

			// Tickets
			tickets := authorized.Group("/tickets")
			{
				ticketHandler := handlers.NewTicketHandler()
				tickets.GET("/:code", ticketHandler.GetTicket)
				tickets.POST("/:code/checkin", ticketHandler.CheckIn)
			}
		}

		// Admin routes
		admin := v1.Group("/admin")
		admin.Use(middleware.AuthMiddleware(), middleware.AdminMiddleware())
		{
			// Bus management
			buses := admin.Group("/buses")
			{
				busHandler := handlers.NewBusHandler()
				buses.POST("", busHandler.Create)
				buses.GET("", busHandler.List)
				buses.GET("/:id", busHandler.GetByID)
				buses.PUT("/:id", busHandler.Update)
				buses.DELETE("/:id", busHandler.Delete)
			}

			// Route management
			routes := admin.Group("/routes")
			{
				routeHandler := handlers.NewRouteHandler()
				routes.POST("", routeHandler.Create)
				routes.GET("", routeHandler.List)
				routes.GET("/:id", routeHandler.GetByID)
				routes.PUT("/:id", routeHandler.Update)
				routes.DELETE("/:id", routeHandler.Delete)
			}

			// Trip management
			trips := admin.Group("/trips")
			{
				tripHandler := handlers.NewTripHandler()
				trips.POST("", tripHandler.Create)
				trips.PUT("/:id", tripHandler.Update)
				trips.DELETE("/:id", tripHandler.Delete)
			}
		}
	}

	return router
}

func startBackgroundJobs(container *Container) {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		ctx := context.Background()

		// Expire old bookings
		_ = container.BookingUsecase.ExpireOldBookings(ctx)

		// Unlock expired seats
		_, _ = container.SeatRepo.UnlockExpiredSeats(ctx)
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
