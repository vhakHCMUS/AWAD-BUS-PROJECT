package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/yourusername/bus-booking/internal/entities"
)

// RedisCache provides caching and distributed locking using Redis
type RedisCache struct {
	client *redis.Client
}

// NewRedisCache creates a new Redis cache instance
func NewRedisCache(client *redis.Client) *RedisCache {
	return &RedisCache{client: client}
}

// Seat Lock Keys
func seatLockKey(tripID uuid.UUID, seatNumber string) string {
	return fmt.Sprintf("seat:lock:%s:%s", tripID.String(), seatNumber)
}

func tripSeatsKey(tripID uuid.UUID) string {
	return fmt.Sprintf("trip:seats:%s", tripID.String())
}

// LockSeat attempts to lock a seat using Redis SETNX with expiry
func (c *RedisCache) LockSeat(ctx context.Context, tripID uuid.UUID, seatNumber string, lockedBy uuid.UUID, duration time.Duration) error {
	key := seatLockKey(tripID, seatNumber)
	success, err := c.client.SetNX(ctx, key, lockedBy.String(), duration).Result()
	if err != nil {
		return err
	}
	if !success {
		return fmt.Errorf("seat is already locked")
	}
	return nil
}

// UnlockSeat releases a seat lock
func (c *RedisCache) UnlockSeat(ctx context.Context, tripID uuid.UUID, seatNumber string) error {
	key := seatLockKey(tripID, seatNumber)
	return c.client.Del(ctx, key).Err()
}

// IsS eatLocked checks if a seat is currently locked
func (c *RedisCache) IsSeatLocked(ctx context.Context, tripID uuid.UUID, seatNumber string) (bool, error) {
	key := seatLockKey(tripID, seatNumber)
	exists, err := c.client.Exists(ctx, key).Result()
	return exists > 0, err
}

// CacheTripSeats caches all seat information for a trip
func (c *RedisCache) CacheTripSeats(ctx context.Context, tripID uuid.UUID, seats []*entities.SeatInfo, ttl time.Duration) error {
	key := tripSeatsKey(tripID)
	data, err := json.Marshal(seats)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, key, data, ttl).Err()
}

// GetTripSeats retrieves cached seat information
func (c *RedisCache) GetTripSeats(ctx context.Context, tripID uuid.UUID) ([]*entities.SeatInfo, error) {
	key := tripSeatsKey(tripID)
	data, err := c.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil // Cache miss
	}
	if err != nil {
		return nil, err
	}

	var seats []*entities.SeatInfo
	err = json.Unmarshal([]byte(data), &seats)
	return seats, err
}

// InvalidateTripSeats removes cached seat data for a trip
func (c *RedisCache) InvalidateTripSeats(ctx context.Context, tripID uuid.UUID) error {
	key := tripSeatsKey(tripID)
	return c.client.Del(ctx, key).Err()
}

// PublishSeatUpdate publishes seat update event to Redis Pub/Sub
func (c *RedisCache) PublishSeatUpdate(ctx context.Context, tripID uuid.UUID, seat *entities.SeatInfo) error {
	channel := fmt.Sprintf("trip:%s:seats", tripID.String())
	data, err := json.Marshal(seat)
	if err != nil {
		return err
	}
	return c.client.Publish(ctx, channel, data).Err()
}

// SubscribeTripUpdates subscribes to seat updates for a trip
func (c *RedisCache) SubscribeTripUpdates(ctx context.Context, tripID uuid.UUID) *redis.PubSub {
	channel := fmt.Sprintf("trip:%s:seats", tripID.String())
	return c.client.Subscribe(ctx, channel)
}

// Rate Limiting
func rateLimitKey(identifier string, window string) string {
	return fmt.Sprintf("ratelimit:%s:%s", identifier, window)
}

// CheckRateLimit checks if request is within rate limit
func (c *RedisCache) CheckRateLimit(ctx context.Context, identifier string, maxRequests int, window time.Duration) (bool, error) {
	key := rateLimitKey(identifier, time.Now().Format("2006-01-02-15:04"))

	pipe := c.client.Pipeline()
	incr := pipe.Incr(ctx, key)
	pipe.Expire(ctx, key, window)
	_, err := pipe.Exec(ctx)
	if err != nil {
		return false, err
	}

	count, err := incr.Result()
	if err != nil {
		return false, err
	}

	return count <= int64(maxRequests), nil
}

// Session/Token Management
func sessionKey(sessionID string) string {
	return fmt.Sprintf("session:%s", sessionID)
}

// SetSession stores session data
func (c *RedisCache) SetSession(ctx context.Context, sessionID string, data interface{}, ttl time.Duration) error {
	key := sessionKey(sessionID)
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, key, jsonData, ttl).Err()
}

// GetSession retrieves session data
func (c *RedisCache) GetSession(ctx context.Context, sessionID string, dest interface{}) error {
	key := sessionKey(sessionID)
	data, err := c.client.Get(ctx, key).Result()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(data), dest)
}

// DeleteSession removes session data
func (c *RedisCache) DeleteSession(ctx context.Context, sessionID string) error {
	key := sessionKey(sessionID)
	return c.client.Del(ctx, key).Err()
}
