# Bus Ticket Booking System

A production-grade, full-stack real-time bus ticket booking platform with high concurrency handling, WebSocket-powered seat selection, payment gateway integration, OAuth2 authentication, and AI chatbot assistance.

## 🚀 Features

### Core Functionality
- **Real-time Seat Selection**: WebSocket-based live seat availability updates
- **Distributed Seat Locking**: Redis + PostgreSQL row-level locking (SELECT FOR UPDATE) prevents double booking
- **Secure Payments**: PayOS/MoMo gateway integration with idempotent webhook handling
- **OAuth 2.0 Authentication**: Google + GitHub login with JWT + refresh tokens
- **E-Ticket Generation**: Automatic PDF generation with QR codes
- **Email Notifications**: Automated booking confirmations and ticket delivery
- **Admin Dashboard**: Comprehensive fleet, route, and trip management
- **AI Chatbot**: Claude 3.5/GPT-4o powered booking assistant with function calling

### Technical Highlights
- **Clean Architecture**: 4-layer separation (entities → usecases → repositories → delivery)
- **High Concurrency**: Handles 1000+ concurrent bookings with distributed locking
- **Production-Ready**: Docker multi-stage builds, CI/CD, health checks
- **Real-time Updates**: Redis Pub/Sub → WebSocket broadcast per trip
- **Type Safety**: Go 1.23 + TypeScript with strict type checking
- **API Documentation**: Auto-generated Swagger/OpenAPI docs

## 📋 Tech Stack

### Backend
- **Language**: Go 1.23
- **Framework**: Gin (HTTP), Gorilla WebSocket
- **Database**: PostgreSQL 16 with GORM
- **Cache**: Redis 7 (locking + pub/sub)
- **Architecture**: Clean Architecture with dependency injection
- **Auth**: JWT, OAuth2 (Google, GitHub)
- **Payments**: PayOS/MoMo SDK
- **Documentation**: Swagger/swaggo
- **Email**: SMTP with gomail
- **PDF/QR**: gofpdf + go-qrcode

### Frontend
- **Framework**: React 18 + Vite + TypeScript
- **Routing**: React Router v6
- **State**: TanStack Query + Zustand
- **Styling**: TailwindCSS + Headless UI
- **Real-time**: Native WebSocket with reconnection
- **Icons**: Lucide React

### Infrastructure
- **Containerization**: Docker + Docker Compose
- **CI/CD**: GitHub Actions
- **Deployment**: Railway/Render ready
- **Database Migrations**: SQL scripts + GORM AutoMigrate

## 🏗️ Architecture

```
backend/
├── cmd/api/main.go              # Application entry point
├── internal/
│   ├── entities/                # Domain models (User, Trip, Booking, etc.)
│   ├── usecases/                # Business logic (BookingUsecase, etc.)
│   ├── repositories/            # Data access interfaces & implementations
│   │   ├── interfaces.go
│   │   ├── postgres/            # PostgreSQL implementations
│   │   └── cache/               # Redis cache & pub/sub
│   ├── delivery/
│   │   ├── http/                # REST API handlers & middleware
│   │   └── websocket/           # WebSocket hub & client management
│   └── infrastructure/          # External services (email, PDF, AI)
├── go.mod
└── Dockerfile

frontend/
├── src/
│   ├── pages/                   # Route components
│   ├── components/              # Reusable UI components
│   ├── lib/                     # API client, WebSocket, utilities
│   ├── config/                  # Constants & environment config
│   └── App.tsx
├── package.json
└── Dockerfile

database/
├── schema.sql                   # PostgreSQL schema with indexes
└── seed.sql                     # Sample data for development
```

## 🚦 Getting Started

### Prerequisites
- Go 1.23+
- Node.js 20+
- PostgreSQL 16
- Redis 7
- Docker & Docker Compose (optional)

### Development Setup

#### 1. Clone the repository
```bash
git clone https://github.com/yourusername/bus-booking.git
cd AWAD-BUS-PROJECT
```

#### 2. Backend Setup
```bash
cd backend

# Install Go dependencies
go mod download

# Copy environment file
cp .env.example .env

# Edit .env with your configuration
# Required: DB credentials, Redis, JWT secret, OAuth keys

# Run database migrations
psql -U postgres -d bus_booking -f ../database/schema.sql
psql -U postgres -d bus_booking -f ../database/seed.sql

# Run the server
go run cmd/api/main.go
```

The API will be available at `http://localhost:8080`
Swagger docs at `http://localhost:8080/swagger/index.html`

#### 3. Frontend Setup
```bash
cd frontend

# Install dependencies
npm install

# Run development server
npm run dev
```

The app will be available at `http://localhost:5173`

### 🐳 Docker Setup (Recommended)

```bash
# Development mode (with hot-reload)
docker-compose -f docker-compose.dev.yml up

# Production mode
docker-compose up --build
```

Services:
- Frontend: http://localhost:3000
- Backend API: http://localhost:8080
- WebSocket: ws://localhost:8081
- PostgreSQL: localhost:5432
- Redis: localhost:6379

## 📚 API Documentation

### Authentication Endpoints
```
POST   /api/v1/auth/register          Register new user
POST   /api/v1/auth/login             Login with email/password
POST   /api/v1/auth/refresh           Refresh access token
GET    /api/v1/auth/google            OAuth with Google
GET    /api/v1/auth/github            OAuth with GitHub
```

### Trip Management
```
GET    /api/v1/trips                  Search trips
GET    /api/v1/trips/:id              Get trip details
GET    /api/v1/trips/:id/seats        Get seat availability
```

### Booking Flow
```
POST   /api/v1/bookings               Initiate booking (lock seats)
GET    /api/v1/bookings/:id           Get booking details
GET    /api/v1/bookings               Get user's bookings
POST   /api/v1/bookings/:id/cancel    Cancel booking
```

### Payment
```
POST   /api/v1/payments               Create payment
POST   /api/v1/payments/webhook       Payment webhook (PayOS/MoMo)
```

### Admin Routes (requires admin role)
```
POST   /api/v1/admin/buses            Create bus
GET    /api/v1/admin/buses            List buses
PUT    /api/v1/admin/buses/:id        Update bus
DELETE /api/v1/admin/buses/:id        Delete bus

# Similar CRUD for routes and trips
```

Full API documentation available at `/swagger` when running the server.

## 🔐 Environment Variables

### Backend (.env)
```bash
# Server
PORT=8080
ENV=development

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=bus_booking

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379

# JWT
JWT_SECRET=your-secret-key
JWT_ACCESS_EXPIRY=15m
JWT_REFRESH_EXPIRY=7d

# OAuth
GOOGLE_CLIENT_ID=your-client-id
GOOGLE_CLIENT_SECRET=your-secret
GITHUB_CLIENT_ID=your-client-id
GITHUB_CLIENT_SECRET=your-secret

# Payment
PAYOS_CLIENT_ID=your-client-id
PAYOS_API_KEY=your-api-key

# Email
SMTP_HOST=smtp.gmail.com
SMTP_USER=your-email@gmail.com
SMTP_PASSWORD=your-app-password

# AI Chatbot
OPENAI_API_KEY=your-api-key
```

### Frontend (.env)
```bash
VITE_API_URL=http://localhost:8080/api/v1
VITE_WS_URL=ws://localhost:8081/ws
```

## 🧪 Testing

### Backend Tests
```bash
cd backend
go test ./... -v
```

### Load Testing with k6
```bash
cd backend
k6 run tests/load/booking_test.js
```

## 🚀 Deployment

### Railway Deployment
1. Push to GitHub
2. Connect repository to Railway
3. Set environment variables
4. Deploy!

### Render Deployment
1. Create new Web Service
2. Connect GitHub repository
3. Configure build command: `docker build`
4. Set environment variables

### Manual Docker Deployment
```bash
# Build images
docker build -t bus-booking-api ./backend
docker build -t bus-booking-frontend ./frontend

# Push to registry
docker push your-registry/bus-booking-api
docker push your-registry/bus-booking-frontend

# Deploy on server
docker-compose -f docker-compose.yml up -d
```

## 📊 Key Implementation Details

### Distributed Seat Locking
```go
// Two-phase locking strategy
// 1. Redis SETNX for fast distributed lock
err := cache.LockSeat(tripID, seatNumber, userID, 10*time.Minute)

// 2. PostgreSQL row-level lock for data consistency
err := db.Clauses(clause.Locking{Strength: "UPDATE"}).
    Where("trip_id = ? AND seat_number = ?", tripID, seatNumber).
    Find(&seat).Error
```

### Real-time Updates
```go
// Publish seat changes to Redis Pub/Sub
cache.PublishSeatUpdate(ctx, tripID, seatInfo)

// WebSocket hub subscribes and broadcasts to clients
hub.broadcast(tripID, seatUpdateMessage)
```

### Payment Webhook Idempotency
```go
// Use idempotency key to prevent double processing
payment := &entities.Payment{
    IdempotencyKey: gatewayPaymentID,
    // ... other fields
}
// Database unique constraint ensures single processing
```

## 🛠️ Development Tools

### Generate Swagger Docs
```bash
cd backend
swag init -g cmd/api/main.go
```

### Database Migrations
```bash
# Create migration
migrate create -ext sql -dir database/migrations -seq add_users_table

# Run migrations
migrate -path database/migrations -database "postgresql://user:pass@localhost/dbname" up
```

## 🤝 Contributing

1. Fork the repository
2. Create feature branch (`git checkout -b feature/amazing-feature`)
3. Commit changes (`git commit -m 'Add amazing feature'`)
4. Push to branch (`git push origin feature/amazing-feature`)
5. Open Pull Request

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 👥 Authors

- Your Name - Initial work

## 🙏 Acknowledgments

- Clean Architecture principles by Robert C. Martin
- Go community for excellent libraries
- React ecosystem for modern frontend tools

## 📞 Support

For support, email support@busbooking.com or create an issue in the GitHub repository.

---

**Built with ❤️ using Go, React, PostgreSQL, and Redis**
