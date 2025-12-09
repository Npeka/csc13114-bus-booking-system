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

type RedisManager interface {
	GetClient() *redis.Client
	Close() error
	HealthCheck() error
	GetStats() *redis.PoolStats
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Del(ctx context.Context, keys ...string) error
	Exists(ctx context.Context, keys ...string) (int64, error)
	Expire(ctx context.Context, key string, expiration time.Duration) error
	TTL(ctx context.Context, key string) (time.Duration, error)
	HSet(ctx context.Context, key string, values ...interface{}) error
	HGet(ctx context.Context, key, field string) (string, error)
	HGetAll(ctx context.Context, key string) (map[string]string, error)
	HDel(ctx context.Context, key string, fields ...string) error
	LPush(ctx context.Context, key string, values ...interface{}) error
	RPush(ctx context.Context, key string, values ...interface{}) error
	LPop(ctx context.Context, key string) (string, error)
	RPop(ctx context.Context, key string) (string, error)
	LLen(ctx context.Context, key string) (int64, error)
	SAdd(ctx context.Context, key string, members ...interface{}) error
	SMembers(ctx context.Context, key string) ([]string, error)
	SIsMember(ctx context.Context, key string, member interface{}) (bool, error)
	SRem(ctx context.Context, key string, members ...interface{}) error
	ZAdd(ctx context.Context, key string, members ...redis.Z) error
	ZRange(ctx context.Context, key string, start, stop int64) ([]string, error)
	ZRangeWithScores(ctx context.Context, key string, start, stop int64) ([]redis.Z, error)
	ZRem(ctx context.Context, key string, members ...interface{}) error
	Incr(ctx context.Context, key string) (int64, error)
	IncrBy(ctx context.Context, key string, value int64) (int64, error)
	Decr(ctx context.Context, key string) (int64, error)
}

type RedisManagerImpl struct {
	Client *redis.Client
}

func MustNewRedisConnection(cfg *config.RedisConfig) RedisManager {
	rm, err := NewRedisConnection(cfg)
	if err != nil {
		panic(err)
	}
	return rm
}

func NewRedisConnection(cfg *config.RedisConfig) (RedisManager, error) {
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

	return &RedisManagerImpl{Client: client}, nil
}

func (rm *RedisManagerImpl) GetClient() *redis.Client {
	return rm.Client
}

func (rm *RedisManagerImpl) Close() error {
	if rm.Client != nil {
		return rm.Client.Close()
	}
	return nil
}

func (rm *RedisManagerImpl) HealthCheck() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := rm.Client.Ping(ctx).Result()
	return err
}

func (rm *RedisManagerImpl) GetStats() *redis.PoolStats {
	if rm.Client == nil {
		return nil
	}
	return rm.Client.PoolStats()
}

func (rm *RedisManagerImpl) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return rm.Client.Set(ctx, key, value, expiration).Err()
}

func (rm *RedisManagerImpl) Get(ctx context.Context, key string) (string, error) {
	return rm.Client.Get(ctx, key).Result()
}

func (rm *RedisManagerImpl) Del(ctx context.Context, keys ...string) error {
	return rm.Client.Del(ctx, keys...).Err()
}

func (rm *RedisManagerImpl) Exists(ctx context.Context, keys ...string) (int64, error) {
	return rm.Client.Exists(ctx, keys...).Result()
}

func (rm *RedisManagerImpl) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return rm.Client.Expire(ctx, key, expiration).Err()
}

func (rm *RedisManagerImpl) TTL(ctx context.Context, key string) (time.Duration, error) {
	return rm.Client.TTL(ctx, key).Result()
}

func (rm *RedisManagerImpl) HSet(ctx context.Context, key string, values ...interface{}) error {
	return rm.Client.HSet(ctx, key, values...).Err()
}

func (rm *RedisManagerImpl) HGet(ctx context.Context, key, field string) (string, error) {
	return rm.Client.HGet(ctx, key, field).Result()
}

func (rm *RedisManagerImpl) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	return rm.Client.HGetAll(ctx, key).Result()
}

func (rm *RedisManagerImpl) HDel(ctx context.Context, key string, fields ...string) error {
	return rm.Client.HDel(ctx, key, fields...).Err()
}

func (rm *RedisManagerImpl) LPush(ctx context.Context, key string, values ...interface{}) error {
	return rm.Client.LPush(ctx, key, values...).Err()
}

func (rm *RedisManagerImpl) RPush(ctx context.Context, key string, values ...interface{}) error {
	return rm.Client.RPush(ctx, key, values...).Err()
}

func (rm *RedisManagerImpl) LPop(ctx context.Context, key string) (string, error) {
	return rm.Client.LPop(ctx, key).Result()
}

func (rm *RedisManagerImpl) RPop(ctx context.Context, key string) (string, error) {
	return rm.Client.RPop(ctx, key).Result()
}

func (rm *RedisManagerImpl) LLen(ctx context.Context, key string) (int64, error) {
	return rm.Client.LLen(ctx, key).Result()
}

func (rm *RedisManagerImpl) SAdd(ctx context.Context, key string, members ...interface{}) error {
	return rm.Client.SAdd(ctx, key, members...).Err()
}

func (rm *RedisManagerImpl) SMembers(ctx context.Context, key string) ([]string, error) {
	return rm.Client.SMembers(ctx, key).Result()
}

func (rm *RedisManagerImpl) SIsMember(ctx context.Context, key string, member interface{}) (bool, error) {
	return rm.Client.SIsMember(ctx, key, member).Result()
}

func (rm *RedisManagerImpl) SRem(ctx context.Context, key string, members ...interface{}) error {
	return rm.Client.SRem(ctx, key, members...).Err()
}

func (rm *RedisManagerImpl) ZAdd(ctx context.Context, key string, members ...redis.Z) error {
	return rm.Client.ZAdd(ctx, key, members...).Err()
}

func (rm *RedisManagerImpl) ZRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return rm.Client.ZRange(ctx, key, start, stop).Result()
}

func (rm *RedisManagerImpl) ZRangeWithScores(ctx context.Context, key string, start, stop int64) ([]redis.Z, error) {
	return rm.Client.ZRangeWithScores(ctx, key, start, stop).Result()
}

func (rm *RedisManagerImpl) ZRem(ctx context.Context, key string, members ...interface{}) error {
	return rm.Client.ZRem(ctx, key, members...).Err()
}

func (rm *RedisManagerImpl) Incr(ctx context.Context, key string) (int64, error) {
	return rm.Client.Incr(ctx, key).Result()
}

func (rm *RedisManagerImpl) IncrBy(ctx context.Context, key string, value int64) (int64, error) {
	return rm.Client.IncrBy(ctx, key, value).Result()
}

func (rm *RedisManagerImpl) Decr(ctx context.Context, key string) (int64, error) {
	return rm.Client.Decr(ctx, key).Result()
}

func (rm *RedisManagerImpl) DecrBy(ctx context.Context, key string, value int64) (int64, error) {
	return rm.Client.DecrBy(ctx, key, value).Result()
}

func (rm *RedisManagerImpl) FlushDB(ctx context.Context) error {
	return rm.Client.FlushDB(ctx).Err()
}

func (rm *RedisManagerImpl) FlushAll(ctx context.Context) error {
	return rm.Client.FlushAll(ctx).Err()
}
