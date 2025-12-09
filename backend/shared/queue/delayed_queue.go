package queue

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

// DelayedItem đại diện cho một item trong hàng đợi trễ
type DelayedItem struct {
	Type    string      `json:"type"`    // Loại job (e.g., "booking_expiry")
	Payload interface{} `json:"payload"` // Dữ liệu của job
}

// DelayedQueueManager định nghĩa interface cho hàng đợi trễ
type DelayedQueueManager interface {
	// Schedule lên lịch một công việc sẽ được thực thi tại thời điểm `executeAt`
	Schedule(ctx context.Context, queueName string, item *DelayedItem, executeAt time.Time) error

	// Poll tìm và lấy các item đã đến hạn xử lý.
	// Sử dụng cơ chế nguyên tử (atomic) để đảm bảo một item chỉ được xử lý bởi một worker.
	Poll(ctx context.Context, queueName string, limit int) ([]*DelayedItem, error)
}

// RedisDelayedQueueManager implementation sử dụng Redis Sorted Sets (ZSET)
type RedisDelayedQueueManager struct {
	client *redis.Client
}

// NewRedisDelayedQueueManager tạo mới một RedisDelayedQueueManager
func NewRedisDelayedQueueManager(client *redis.Client) DelayedQueueManager {
	return &RedisDelayedQueueManager{
		client: client,
	}
}

// Schedule thêm item vào ZSET với score là timestamp
func (m *RedisDelayedQueueManager) Schedule(ctx context.Context, queueName string, item *DelayedItem, executeAt time.Time) error {
	data, err := json.Marshal(item)
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal delayed item")
		return err
	}

	// Score là Unix timestamp
	score := float64(executeAt.Unix())

	err = m.client.ZAdd(ctx, queueName, redis.Z{
		Score:  score,
		Member: data,
	}).Err()

	if err != nil {
		log.Error().Err(err).Str("queue", queueName).Msg("Failed to schedule delayed item")
		return err
	}

	// Logging an toàn (không log toàn bộ data nếu nhạy cảm, ở đây log type cho debug)
	log.Debug().
		Str("queue", queueName).
		Str("type", item.Type).
		Time("execute_at", executeAt).
		Msg("Scheduled delayed job")

	return nil
}

// Poll sử dụng Lua script để atomic fetch-and-delete các item đã đến hạn
func (m *RedisDelayedQueueManager) Poll(ctx context.Context, queueName string, limit int) ([]*DelayedItem, error) {
	now := float64(time.Now().Unix())

	// Lua script:
	// 1. ZRANGEBYSCORE key -inf now LIMIT 0 limit
	// 2. Nếu có item, ZREM key item...
	// 3. Return items
	script := redis.NewScript(`
		local items = redis.call('ZRANGEBYSCORE', KEYS[1], '-inf', ARGV[1], 'LIMIT', 0, ARGV[2])
		if #items > 0 then
			redis.call('ZREM', KEYS[1], unpack(items))
		end
		return items
	`)

	// Run script
	result, err := script.Run(ctx, m.client, []string{queueName}, now, limit).Result()
	if err != nil {
		if err == redis.Nil {
			return []*DelayedItem{}, nil
		}
		log.Error().Err(err).Str("queue", queueName).Msg("Failed to poll delayed queue")
		return nil, err
	}

	// Parse result
	paramsInterface, ok := result.([]interface{})
	if !ok {
		// Trường hợp trả về rỗng hoặc format lạ
		return []*DelayedItem{}, nil
	}

	var items []*DelayedItem
	for _, p := range paramsInterface {
		str, ok := p.(string)
		if !ok {
			continue
		}

		var item DelayedItem
		if err := json.Unmarshal([]byte(str), &item); err != nil {
			log.Error().Err(err).Str("raw", str).Msg("Failed to unmarshal delayed item")
			continue // Skip bad items
		}
		items = append(items, &item)
	}

	return items, nil
}
