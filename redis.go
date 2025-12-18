package main

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client

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

// RedisExists checks whether a key exists.
func RedisExists(ctx context.Context, r *redis.Client, key string) (bool, error) {
	n, err := r.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("redis exists %q: %w", key, err)
	}
	return n > 0, nil
}

// RedisDel deletes one or more keys.
func RedisDel(ctx context.Context, r *redis.Client, keys ...string) (int64, error) {
	n, err := r.Del(ctx, keys...).Result()
	if err != nil {
		return 0, fmt.Errorf("redis del: %w", err)
	}
	return n, nil
}

// RedisIncr increments a key and returns the new value.
func RedisIncr(ctx context.Context, r *redis.Client, key string) (int64, error) {
	v, err := r.Incr(ctx, key).Result()
	if err != nil {
		return 0, fmt.Errorf("redis incr %q: %w", key, err)
	}
	return v, nil
}

// RedisExpire sets a TTL on a key.
func RedisExpire(ctx context.Context, r *redis.Client, key string, expiration time.Duration) (bool, error) {
	ok, err := r.Expire(ctx, key, expiration).Result()
	if err != nil {
		return false, fmt.Errorf("redis expire %q: %w", key, err)
	}
	return ok, nil
}

// RedisMSet sets multiple keys to multiple values.
func RedisMSet(ctx context.Context, r *redis.Client, values map[string]interface{}) error {
	if err := r.MSet(ctx, values).Err(); err != nil {
		return fmt.Errorf("redis mset: %w", err)
	}
	return nil
}

// RedisMGet gets multiple keys.
func RedisMGet(ctx context.Context, r *redis.Client, keys ...string) ([]interface{}, error) {
	vals, err := r.MGet(ctx, keys...).Result()
	if err != nil {
		return nil, fmt.Errorf("redis mget: %w", err)
	}
	return vals, nil
}

// Hash operations
// RedisHSet sets one or more field-value pairs in a hash.
func RedisHSet(ctx context.Context, r *redis.Client, key string, values map[string]interface{}) error {
	if err := r.HSet(ctx, key, values).Err(); err != nil {
		return fmt.Errorf("redis hset %q: %w", key, err)
	}
	return nil
}

// RedisHGet gets a field from a hash.
func RedisHGet(ctx context.Context, r *redis.Client, key, field string) (string, error) {
	v, err := r.HGet(ctx, key, field).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("hash field %q not found in %q", field, key)
	}
	if err != nil {
		return "", fmt.Errorf("redis hget %q %q: %w", key, field, err)
	}
	return v, nil
}

// RedisHGetAll returns all fields and values in a hash.
func RedisHGetAll(ctx context.Context, r *redis.Client, key string) (map[string]string, error) {
	m, err := r.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("redis hgetall %q: %w", key, err)
	}
	return m, nil
}

// RedisHDel deletes one or more hash fields.
func RedisHDel(ctx context.Context, r *redis.Client, key string, fields ...string) (int64, error) {
	n, err := r.HDel(ctx, key, fields...).Result()
	if err != nil {
		return 0, fmt.Errorf("redis hdel %q: %w", key, err)
	}
	return n, nil
}

// Sorted set (zset) operations
// RedisZAdd adds members with scores to a sorted set.
func RedisZAdd(ctx context.Context, r *redis.Client, key string, members ...redis.Z) (int64, error) {
	n, err := r.ZAdd(ctx, key, members...).Result()
	if err != nil {
		return 0, fmt.Errorf("redis zadd %q: %w", key, err)
	}
	return n, nil
}

// RedisZRangeByScore returns members in a score range.
func RedisZRangeByScore(ctx context.Context, r *redis.Client, key string, opt *redis.ZRangeBy) ([]string, error) {
	vals, err := r.ZRangeByScore(ctx, key, opt).Result()
	if err != nil {
		return nil, fmt.Errorf("redis zrangebyscore %q: %w", key, err)
	}
	return vals, nil
}

// RedisZRem removes one or more members from a sorted set.
func RedisZRem(ctx context.Context, r *redis.Client, key string, members ...interface{}) (int64, error) {
	n, err := r.ZRem(ctx, key, members...).Result()
	if err != nil {
		return 0, fmt.Errorf("redis zrem %q: %w", key, err)
	}
	return n, nil
}

// RedisZRangeWithScores returns members with scores.
func RedisZRangeWithScores(ctx context.Context, r *redis.Client, key string, start, stop int64) ([]redis.Z, error) {
	vals, err := r.ZRangeWithScores(ctx, key, start, stop).Result()
	if err != nil {
		return nil, fmt.Errorf("redis zrangewithscores %q: %w", key, err)
	}
	return vals, nil
}
