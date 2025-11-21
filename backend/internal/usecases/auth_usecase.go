package usecases

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/yourusername/bus-booking/internal/entities"
	"github.com/yourusername/bus-booking/internal/repositories"
)

type AuthUsecase struct {
	userRepo           repositories.UserRepository
	jwtSecret          string
	accessTokenExpiry  time.Duration
	refreshTokenExpiry time.Duration
}

func NewAuthUsecase(
	userRepo repositories.UserRepository,
	jwtSecret string,
	accessTokenExpiry, refreshTokenExpiry time.Duration,
) *AuthUsecase {
	return &AuthUsecase{
		userRepo:           userRepo,
		jwtSecret:          jwtSecret,
		accessTokenExpiry:  accessTokenExpiry,
		refreshTokenExpiry: refreshTokenExpiry,
	}
}

type RegisterInput struct {
	Name     string
	Email    string
	Password string
	Phone    string
}

type LoginInput struct {
	Email    string
	Password string
}

type AuthTokens struct {
	AccessToken  string
	RefreshToken string
	User         *entities.User
}

// Register creates a new user account
func (uc *AuthUsecase) Register(ctx context.Context, input RegisterInput) (*AuthTokens, error) {
	// Check if user already exists
	existing, _ := uc.userRepo.GetByEmail(ctx, input.Email)
	if existing != nil {
		return nil, errors.New("user with this email already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user := &entities.User{
		ID:           uuid.New(),
		Name:         input.Name,
		Email:        input.Email,
		PasswordHash: string(hashedPassword),
		Phone:        input.Phone,
		Role:         entities.RolePassenger,
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := uc.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Generate tokens
	accessToken, err := uc.generateAccessToken(user)
	if err != nil {
		return nil, err
	}

	refreshToken, err := uc.generateRefreshToken(user)
	if err != nil {
		return nil, err
	}

	// Remove password from response
	user.PasswordHash = ""

	return &AuthTokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
	}, nil
}

// Login authenticates a user
func (uc *AuthUsecase) Login(ctx context.Context, input LoginInput) (*AuthTokens, error) {
	// Get user by email
	user, err := uc.userRepo.GetByEmail(ctx, input.Email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	if !user.IsActive {
		return nil, errors.New("account is inactive")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Generate tokens
	accessToken, err := uc.generateAccessToken(user)
	if err != nil {
		return nil, err
	}

	refreshToken, err := uc.generateRefreshToken(user)
	if err != nil {
		return nil, err
	}

	// Remove password from response
	user.PasswordHash = ""

	return &AuthTokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
	}, nil
}

// generateAccessToken creates a JWT access token
func (uc *AuthUsecase) generateAccessToken(user *entities.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID.String(),
		"email":   user.Email,
		"role":    user.Role,
		"exp":     time.Now().Add(uc.accessTokenExpiry).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(uc.jwtSecret))
}

// generateRefreshToken creates a JWT refresh token
func (uc *AuthUsecase) generateRefreshToken(user *entities.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID.String(),
		"exp":     time.Now().Add(uc.refreshTokenExpiry).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(uc.jwtSecret))
}

// RefreshAccessToken generates a new access token from refresh token
func (uc *AuthUsecase) RefreshAccessToken(ctx context.Context, refreshToken string) (*AuthTokens, error) {
	// Parse and validate refresh token
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(uc.jwtSecret), nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid refresh token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	userIDStr, ok := claims["user_id"].(string)
	if !ok {
		return nil, errors.New("invalid user ID in token")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, errors.New("invalid user ID format")
	}

	// Get user
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if !user.IsActive {
		return nil, errors.New("account is inactive")
	}

	// Generate new tokens
	newAccessToken, err := uc.generateAccessToken(user)
	if err != nil {
		return nil, err
	}

	newRefreshToken, err := uc.generateRefreshToken(user)
	if err != nil {
		return nil, err
	}

	// Remove password from response
	user.PasswordHash = ""

	return &AuthTokens{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
		User:         user,
	}, nil
}
