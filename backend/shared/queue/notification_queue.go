package queue

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

// EmailNotification đại diện cho email notification message
type EmailNotification struct {
	To           string                 `json:"to"`
	Subject      string                 `json:"subject"`
	TemplateName string                 `json:"template_name"`
	TemplateData map[string]interface{} `json:"template_data"`
	Priority     int                    `json:"priority"` // 1 = cao, 2 = trung bình, 3 = thấp
	ScheduledAt  *time.Time             `json:"scheduled_at,omitempty"`
}

// SMSNotification đại diện cho SMS notification message
type SMSNotification struct {
	To           string                 `json:"to"`
	Message      string                 `json:"message"`
	TemplateName string                 `json:"template_name,omitempty"`
	TemplateData map[string]interface{} `json:"template_data,omitempty"`
	Priority     int                    `json:"priority"`
	ScheduledAt  *time.Time             `json:"scheduled_at,omitempty"`
}

// NotificationQueue định nghĩa interface cho notification queue
type NotificationQueue interface {
	// PushEmailNotification đẩy email notification vào queue
	PushEmailNotification(ctx context.Context, notification *EmailNotification) error

	// PushSMSNotification đẩy SMS notification vào queue
	PushSMSNotification(ctx context.Context, notification *SMSNotification) error

	// PopEmailNotification lấy email notification từ queue (dành cho worker)
	PopEmailNotification(ctx context.Context, timeout time.Duration) (*EmailNotification, error)

	// PopSMSNotification lấy SMS notification từ queue (dành cho worker)
	PopSMSNotification(ctx context.Context, timeout time.Duration) (*SMSNotification, error)

	// GetQueueLength lấy số lượng message trong queue
	GetQueueLength(ctx context.Context, queueName string) (int64, error)
}

// RedisNotificationQueue implementation sử dụng Redis
type RedisNotificationQueue struct {
	client        *redis.Client
	emailQueueKey string
	smsQueueKey   string
}

// NewRedisNotificationQueue tạo mới Redis notification queue
func NewRedisNotificationQueue(redisClient *redis.Client, queuePrefix string) NotificationQueue {
	if queuePrefix == "" {
		queuePrefix = "notifications"
	}

	return &RedisNotificationQueue{
		client:        redisClient,
		emailQueueKey: queuePrefix + ":email",
		smsQueueKey:   queuePrefix + ":sms",
	}
}

// PushEmailNotification đẩy email notification vào Redis queue
// NOTE: Hiện tại chưa có notification service, nên method này chỉ để chuẩn bị
func (q *RedisNotificationQueue) PushEmailNotification(ctx context.Context, notification *EmailNotification) error {
	// Serialize notification to JSON
	data, err := json.Marshal(notification)
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal email notification")
		return err
	}

	// Push to Redis list (LPUSH for FIFO with BRPOP)
	err = q.client.LPush(ctx, q.emailQueueKey, data).Err()
	if err != nil {
		log.Error().Err(err).Str("queue", q.emailQueueKey).Msg("Failed to push email notification to queue")
		return err
	}

	log.Info().
		Str("queue", q.emailQueueKey).
		Str("to", notification.To).
		Str("template", notification.TemplateName).
		Msg("Email notification pushed to queue")

	// TODO: Set expiration if needed
	// q.client.Expire(ctx, q.emailQueueKey, 24*time.Hour)

	return nil
}

// PushSMSNotification đẩy SMS notification vào Redis queue
// NOTE: Hiện tại chưa có notification service, nên method này chỉ để chuẩn bị
func (q *RedisNotificationQueue) PushSMSNotification(ctx context.Context, notification *SMSNotification) error {
	// Serialize notification to JSON
	data, err := json.Marshal(notification)
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal SMS notification")
		return err
	}

	// Push to Redis list
	err = q.client.LPush(ctx, q.smsQueueKey, data).Err()
	if err != nil {
		log.Error().Err(err).Str("queue", q.smsQueueKey).Msg("Failed to push SMS notification to queue")
		return err
	}

	log.Info().
		Str("queue", q.smsQueueKey).
		Str("to", notification.To).
		Msg("SMS notification pushed to queue")

	return nil
}

// PopEmailNotification lấy email notification từ queue (blocking với timeout)
// NOTE: Method này sẽ được notification service worker sử dụng
func (q *RedisNotificationQueue) PopEmailNotification(ctx context.Context, timeout time.Duration) (*EmailNotification, error) {
	// BRPOP: Blocking right pop với timeout
	result, err := q.client.BRPop(ctx, timeout, q.emailQueueKey).Result()
	if err != nil {
		if err == redis.Nil {
			// Queue empty, timeout reached
			return nil, nil
		}
		log.Error().Err(err).Str("queue", q.emailQueueKey).Msg("Failed to pop email notification from queue")
		return nil, err
	}

	// result[0] is the queue key, result[1] is the value
	if len(result) < 2 {
		return nil, nil
	}

	var notification EmailNotification
	err = json.Unmarshal([]byte(result[1]), &notification)
	if err != nil {
		log.Error().Err(err).Msg("Failed to unmarshal email notification")
		return nil, err
	}

	return &notification, nil
}

// PopSMSNotification lấy SMS notification từ queue (blocking với timeout)
// NOTE: Method này sẽ được notification service worker sử dụng
func (q *RedisNotificationQueue) PopSMSNotification(ctx context.Context, timeout time.Duration) (*SMSNotification, error) {
	result, err := q.client.BRPop(ctx, timeout, q.smsQueueKey).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		log.Error().Err(err).Str("queue", q.smsQueueKey).Msg("Failed to pop SMS notification from queue")
		return nil, err
	}

	if len(result) < 2 {
		return nil, nil
	}

	var notification SMSNotification
	err = json.Unmarshal([]byte(result[1]), &notification)
	if err != nil {
		log.Error().Err(err).Msg("Failed to unmarshal SMS notification")
		return nil, err
	}

	return &notification, nil
}

// GetQueueLength lấy số lượng message trong queue
func (q *RedisNotificationQueue) GetQueueLength(ctx context.Context, queueName string) (int64, error) {
	length, err := q.client.LLen(ctx, queueName).Result()
	if err != nil {
		log.Error().Err(err).Str("queue", queueName).Msg("Failed to get queue length")
		return 0, err
	}
	return length, nil
}
