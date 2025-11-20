# Quick Start Guide

## Prerequisites Check

Before you begin, ensure you have the following installed:

- [ ] Go 1.23 or higher
- [ ] Node.js 20 or higher
- [ ] PostgreSQL 16 or higher
- [ ] Redis 7 or higher
- [ ] Docker & Docker Compose (optional but recommended)

## 🚀 Quick Start with Docker (Recommended)

This is the fastest way to get the entire stack running:

```bash
# 1. Clone the repository
git clone https://github.com/yourusername/AWAD-BUS-PROJECT.git
cd AWAD-BUS-PROJECT

# 2. Start all services with Docker Compose
docker-compose up --build

# That's it! The application is now running:
# - Frontend: http://localhost:3000
# - Backend API: http://localhost:8080
# - API Docs: http://localhost:8080/swagger/index.html
# - WebSocket: ws://localhost:8081
```

## 📋 Manual Setup (Development)

### Step 1: Database Setup

```bash
# Start PostgreSQL (if not already running)
# On Windows with PostgreSQL installed:
# Services → PostgreSQL → Start

# Create database
psql -U postgres -c "CREATE DATABASE bus_booking;"

# Run migrations
psql -U postgres -d bus_booking -f database/schema.sql

# Seed sample data
psql -U postgres -d bus_booking -f database/seed.sql
```

### Step 2: Redis Setup

```bash
# Start Redis
# On Windows: redis-server
# On macOS: brew services start redis
# On Linux: sudo systemctl start redis

# Verify Redis is running
redis-cli ping
# Should return: PONG
```

### Step 3: Backend Setup

```bash
cd backend

# Install Go dependencies
go mod download

# Copy and configure environment
cp .env.example .env

# Edit .env and update the following:
# - DB_PASSWORD (your PostgreSQL password)
# - JWT_SECRET (generate a random secret)
# - OAuth credentials (if you have them)
# - Payment gateway credentials (if you have them)

# Start the backend server
go run cmd/api/main.go

# Server will start on http://localhost:8080
# Visit http://localhost:8080/health to verify it's running
```

### Step 4: Frontend Setup

```bash
cd frontend

# Install dependencies
npm install

# Start development server
npm run dev

# Frontend will start on http://localhost:5173
```

## 🔧 Configuration

### Backend Environment Variables

Create `backend/.env`:

```env
# Server
PORT=8080
ENV=development

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password_here
DB_NAME=bus_booking
DB_SSL_MODE=disable

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=

# JWT (Generate a strong random secret)
JWT_SECRET=your-super-secret-jwt-key-change-this-now
JWT_ACCESS_EXPIRY=15m
JWT_REFRESH_EXPIRY=7d

# OAuth (Get from Google/GitHub Developer Console)
GOOGLE_CLIENT_ID=your-google-client-id
GOOGLE_CLIENT_SECRET=your-google-client-secret
GOOGLE_REDIRECT_URL=http://localhost:8080/api/v1/auth/google/callback

GITHUB_CLIENT_ID=your-github-client-id
GITHUB_CLIENT_SECRET=your-github-client-secret
GITHUB_REDIRECT_URL=http://localhost:8080/api/v1/auth/github/callback

# Payment Gateway (Get from PayOS/MoMo)
PAYOS_CLIENT_ID=your-payos-client-id
PAYOS_API_KEY=your-payos-api-key
PAYOS_CHECKSUM_KEY=your-payos-checksum-key

# Email (Use Gmail with App Password)
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your-email@gmail.com
SMTP_PASSWORD=your-app-password

# AI Chatbot (Optional - for chatbot features)
OPENAI_API_KEY=your-openai-api-key
```

### Frontend Environment Variables

Create `frontend/.env`:

```env
VITE_API_URL=http://localhost:8080/api/v1
VITE_WS_URL=ws://localhost:8081/ws
```

## 🧪 Testing the Setup

### 1. Test Backend API

```bash
# Health check
curl http://localhost:8080/health

# Should return: {"status":"ok"}
```

### 2. Test Database Connection

```bash
# List trips
curl http://localhost:8080/api/v1/trips

# Should return sample trips from seed data
```

### 3. Test Frontend

1. Open http://localhost:5173 in your browser
2. You should see the homepage
3. Click "Search Trips" to see sample routes

### 4. Test Real-time Features

1. Open the browser console
2. Navigate to a trip detail page
3. Watch for WebSocket connection messages

## 👤 Default Credentials

The seed data includes these test accounts:

**Admin Account:**
- Email: admin@busbooking.com
- Password: admin123

**Passenger Account:**
- Email: passenger@example.com
- Password: passenger123

## 🐛 Troubleshooting

### Backend won't start

```bash
# Check if PostgreSQL is running
psql -U postgres -c "SELECT 1;"

# Check if Redis is running
redis-cli ping

# Check Go version
go version  # Should be 1.23+

# Clear Go cache
go clean -cache
```

### Database connection error

```bash
# Verify database exists
psql -U postgres -l | grep bus_booking

# Check connection
psql -U postgres -d bus_booking -c "SELECT COUNT(*) FROM users;"
```

### Frontend errors

```bash
# Clear node modules and reinstall
rm -rf node_modules package-lock.json
npm install

# Check Node version
node --version  # Should be 20+
```

### Port already in use

```bash
# On Windows (PowerShell):
netstat -ano | findstr :8080
# Kill the process using the port

# On macOS/Linux:
lsof -ti:8080 | xargs kill -9
```

## 📚 Next Steps

Once your setup is working:

1. **Explore the API**: Visit http://localhost:8080/swagger/index.html
2. **Review Architecture**: Read the main README.md
3. **Customize**: Modify configs in `.env` files
4. **Add Features**: Follow the clean architecture pattern
5. **Deploy**: Use Docker Compose for production deployment

## 🆘 Getting Help

If you encounter issues:

1. Check the logs:
   - Backend: Console output from `go run cmd/api/main.go`
   - Frontend: Browser console (F12)
   - Docker: `docker-compose logs -f`

2. Verify all services:
   ```bash
   # PostgreSQL
   pg_isready
   
   # Redis
   redis-cli ping
   
   # Backend
   curl http://localhost:8080/health
   ```

3. Create an issue on GitHub with:
   - Error message
   - Steps to reproduce
   - Your environment (OS, versions)

## 🎉 Success!

If everything is working, you should be able to:
- ✅ Browse trips
- ✅ Select seats in real-time
- ✅ Create bookings
- ✅ Login with test accounts
- ✅ Access admin panel (with admin account)

Happy coding! 🚀
