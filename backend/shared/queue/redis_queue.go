package queue

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

// QueueManager defines the interface for a generic queue manager
type QueueManager interface {
	// Push pushes an item to the queue
	Push(ctx context.Context, queueName string, item interface{}) error

	// Pop pops an item from the queue (blocking with timeout)
	// target is a pointer to the struct where the data will be unmarshaled
	Pop(ctx context.Context, queueName string, timeout time.Duration, target interface{}) error

	// Length returns the number of items in the queue
	Length(ctx context.Context, queueName string) (int64, error)
}

// RedisQueueManager implements QueueManager using Redis
type RedisQueueManager struct {
	client *redis.Client
}

// NewRedisQueueManager creates a new RedisQueueManager
func NewRedisQueueManager(client *redis.Client) QueueManager {
	return &RedisQueueManager{
		client: client,
	}
}

// Push pushes an item to the Redis list
func (m *RedisQueueManager) Push(ctx context.Context, queueName string, item interface{}) error {
	data, err := json.Marshal(item)
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal item for queue")
		return err
	}

	err = m.client.LPush(ctx, queueName, data).Err()
	if err != nil {
		log.Error().Err(err).Str("queue", queueName).Msg("Failed to push item to Redis queue")
		return err
	}

	return nil
}

// Pop pops an item from the Redis list (blocking)
func (m *RedisQueueManager) Pop(ctx context.Context, queueName string, timeout time.Duration, target interface{}) error {
	result, err := m.client.BRPop(ctx, timeout, queueName).Result()
	if err != nil {
		if err == redis.Nil {
			return nil // Timeout, no item
		}
		log.Error().Err(err).Str("queue", queueName).Msg("Failed to pop item from Redis queue")
		return err
	}

	if len(result) < 2 {
		return nil
	}

	err = json.Unmarshal([]byte(result[1]), target)
	if err != nil {
		log.Error().Err(err).Msg("Failed to unmarshal item from queue")
		return err
	}

	return nil
}

// Length returns the length of the Redis list
func (m *RedisQueueManager) Length(ctx context.Context, queueName string) (int64, error) {
	length, err := m.client.LLen(ctx, queueName).Result()
	if err != nil {
		log.Error().Err(err).Str("queue", queueName).Msg("Failed to get Redis queue length")
		return 0, err
	}
	return length, nil
}
