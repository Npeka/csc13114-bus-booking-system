package db

import (
	"context"
	"crypto/tls"
	"fmt"
	"time"

	"bus-booking/shared/config"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

type RedisManager struct {
	Client *redis.Client
	Config *config.RedisConfig
}

func NewRedisConnection(cfg *config.RedisConfig) (*RedisManager, error) {
	// Redis client options
	options := &redis.Options{
		Addr:         fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Username:     cfg.Username,
		Password:     cfg.Password,
		DB:           cfg.DB,
		PoolSize:     cfg.PoolSize,
		MinIdleConns: cfg.MinIdleConns,
		MaxRetries:   cfg.MaxRetries,
		DialTimeout:  cfg.DialTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		PoolTimeout:  30 * time.Second,
	}

	if cfg.TLS {
		options.TLSConfig = &tls.Config{
			InsecureSkipVerify: false,
		}
	}

	client := redis.NewClient(options)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	log.Info().
		Str("host", cfg.Host).
		Int("port", cfg.Port).
		Int("db", cfg.DB).
		Msg("Successfully connected to Redis")

	return &RedisManager{
		Client: client,
		Config: cfg,
	}, nil
}

func (rm *RedisManager) Close() error {
	if rm.Client != nil {
		return rm.Client.Close()
	}
	return nil
}

func (rm *RedisManager) HealthCheck() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := rm.Client.Ping(ctx).Result()
	return err
}

func (rm *RedisManager) GetStats() *redis.PoolStats {
	if rm.Client == nil {
		return nil
	}
	return rm.Client.PoolStats()
}

func (rm *RedisManager) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return rm.Client.Set(ctx, key, value, expiration).Err()
}

func (rm *RedisManager) Get(ctx context.Context, key string) (string, error) {
	return rm.Client.Get(ctx, key).Result()
}

func (rm *RedisManager) Del(ctx context.Context, keys ...string) error {
	return rm.Client.Del(ctx, keys...).Err()
}

func (rm *RedisManager) Exists(ctx context.Context, keys ...string) (int64, error) {
	return rm.Client.Exists(ctx, keys...).Result()
}

func (rm *RedisManager) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return rm.Client.Expire(ctx, key, expiration).Err()
}

func (rm *RedisManager) TTL(ctx context.Context, key string) (time.Duration, error) {
	return rm.Client.TTL(ctx, key).Result()
}

func (rm *RedisManager) HSet(ctx context.Context, key string, values ...interface{}) error {
	return rm.Client.HSet(ctx, key, values...).Err()
}

func (rm *RedisManager) HGet(ctx context.Context, key, field string) (string, error) {
	return rm.Client.HGet(ctx, key, field).Result()
}

func (rm *RedisManager) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	return rm.Client.HGetAll(ctx, key).Result()
}

func (rm *RedisManager) HDel(ctx context.Context, key string, fields ...string) error {
	return rm.Client.HDel(ctx, key, fields...).Err()
}

func (rm *RedisManager) LPush(ctx context.Context, key string, values ...interface{}) error {
	return rm.Client.LPush(ctx, key, values...).Err()
}

func (rm *RedisManager) RPush(ctx context.Context, key string, values ...interface{}) error {
	return rm.Client.RPush(ctx, key, values...).Err()
}

func (rm *RedisManager) LPop(ctx context.Context, key string) (string, error) {
	return rm.Client.LPop(ctx, key).Result()
}

func (rm *RedisManager) RPop(ctx context.Context, key string) (string, error) {
	return rm.Client.RPop(ctx, key).Result()
}

func (rm *RedisManager) LLen(ctx context.Context, key string) (int64, error) {
	return rm.Client.LLen(ctx, key).Result()
}

func (rm *RedisManager) SAdd(ctx context.Context, key string, members ...interface{}) error {
	return rm.Client.SAdd(ctx, key, members...).Err()
}

func (rm *RedisManager) SMembers(ctx context.Context, key string) ([]string, error) {
	return rm.Client.SMembers(ctx, key).Result()
}

func (rm *RedisManager) SIsMember(ctx context.Context, key string, member interface{}) (bool, error) {
	return rm.Client.SIsMember(ctx, key, member).Result()
}

func (rm *RedisManager) SRem(ctx context.Context, key string, members ...interface{}) error {
	return rm.Client.SRem(ctx, key, members...).Err()
}

func (rm *RedisManager) ZAdd(ctx context.Context, key string, members ...redis.Z) error {
	return rm.Client.ZAdd(ctx, key, members...).Err()
}

func (rm *RedisManager) ZRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return rm.Client.ZRange(ctx, key, start, stop).Result()
}

func (rm *RedisManager) ZRangeWithScores(ctx context.Context, key string, start, stop int64) ([]redis.Z, error) {
	return rm.Client.ZRangeWithScores(ctx, key, start, stop).Result()
}

func (rm *RedisManager) ZRem(ctx context.Context, key string, members ...interface{}) error {
	return rm.Client.ZRem(ctx, key, members...).Err()
}

func (rm *RedisManager) Incr(ctx context.Context, key string) (int64, error) {
	return rm.Client.Incr(ctx, key).Result()
}

func (rm *RedisManager) IncrBy(ctx context.Context, key string, value int64) (int64, error) {
	return rm.Client.IncrBy(ctx, key, value).Result()
}

func (rm *RedisManager) Decr(ctx context.Context, key string) (int64, error) {
	return rm.Client.Decr(ctx, key).Result()
}

func (rm *RedisManager) DecrBy(ctx context.Context, key string, value int64) (int64, error) {
	return rm.Client.DecrBy(ctx, key, value).Result()
}

func (rm *RedisManager) FlushDB(ctx context.Context) error {
	return rm.Client.FlushDB(ctx).Err()
}

func (rm *RedisManager) FlushAll(ctx context.Context) error {
	return rm.Client.FlushAll(ctx).Err()
}
