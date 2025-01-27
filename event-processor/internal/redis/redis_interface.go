package redis

import (
	"context"
	"event-processor/internal/config"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisClient interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Del(ctx context.Context, key string) error
	HSet(ctx context.Context, key string, values map[string]interface{}) error
	HGetAll(ctx context.Context, key string) (map[string]string, error)
	LockChunk(ctx context.Context, chunkID string, workerID string, expiration time.Duration) (bool, error)
	UnlockChunk(ctx context.Context, chunkID string) error
	Close() error
}

type redisClient struct {
	client *redis.Client
}

// NewRedisClient initializes a new Redis client
func NewRedisClient(config *config.Config) RedisClient {
	addr := config.RedisHost + ":" + config.RedisPort

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: config.RedisPassword,
		DB:       config.RedisDB,
	})

	return &redisClient{
		client: client,
	}
}

// Set sets a key-value pair with an expiration time
func (r *redisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return r.client.Set(ctx, key, value, expiration).Err()
}

// Get retrieves the value of a key
func (r *redisClient) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

// Del deletes a key
func (r *redisClient) Del(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

// HSet sets multiple fields in a hash
func (r *redisClient) HSet(ctx context.Context, key string, values map[string]interface{}) error {
	return r.client.HSet(ctx, key, values).Err()
}

// HGetAll retrieves all fields and values in a hash
func (r *redisClient) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	return r.client.HGetAll(ctx, key).Result()
}

// LockChunk locks a chunk for processing by a worker
func (r *redisClient) LockChunk(ctx context.Context, chunkID string, workerID string, expiration time.Duration) (bool, error) {
	success, err := r.client.SetNX(ctx, chunkID, workerID, expiration).Result()
	if err != nil {
		return false, err
	}
	return success, nil
}

// UnlockChunk unlocks a chunk
func (r *redisClient) UnlockChunk(ctx context.Context, chunkID string) error {
	return r.client.Del(ctx, chunkID).Err()
}

// Close closes the Redis client connection
func (r *redisClient) Close() error {
	return r.client.Close()
}
