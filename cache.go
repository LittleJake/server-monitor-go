package main

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// NewRedisClient returns a configured Redis client.
func NewRedisClient(addr, password string, db int) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     password,
		DB:           db,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolSize:     10,
	})
}

// RedisPing verifies connectivity to Redis.
func RedisPing(ctx context.Context, r *redis.Client) error {
	if err := r.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("redis ping: %w", err)
	}
	return nil
}

// RedisSet sets a key with an expiration.
func RedisSet(ctx context.Context, r *redis.Client, key string, value interface{}, expiration time.Duration) error {
	if err := r.Set(ctx, key, value, expiration).Err(); err != nil {
		return fmt.Errorf("redis set %q: %w", key, err)
	}
	return nil
}

// RedisGet retrieves a string value for a key.
func RedisGet(ctx context.Context, r *redis.Client, key string) (string, error) {
	val, err := r.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("key %q not found", key)
	}
	if err != nil {
		return "", fmt.Errorf("redis get %q: %w", key, err)
	}
	return val, nil
}

// RedisSubscribe subscribes to a channel and calls handler for each message.
// The subscription runs in a goroutine; caller is responsible for context cancellation.
func RedisSubscribe(ctx context.Context, r *redis.Client, channel string, handler func(msg *redis.Message)) error {
	sub := r.Subscribe(ctx, channel)
	// Confirm subscription
	if _, err := sub.Receive(ctx); err != nil {
		return fmt.Errorf("subscribe receive: %w", err)
	}
	ch := sub.Channel()

	go func() {
		for msg := range ch {
			handler(msg)
		}
	}()

	return nil
}

// RedisClose closes the client connection.
func RedisClose(r *redis.Client) error {
	return r.Close()
}
