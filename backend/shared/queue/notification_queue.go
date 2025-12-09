package queue

import (
	"context"
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
	manager       QueueManager
	emailQueueKey string
	smsQueueKey   string
}

// NewRedisNotificationQueue tạo mới Redis notification queue
func NewRedisNotificationQueue(redisClient *redis.Client, queuePrefix string) NotificationQueue {
	if queuePrefix == "" {
		queuePrefix = "notifications"
	}

	return &RedisNotificationQueue{
		manager:       NewRedisQueueManager(redisClient),
		emailQueueKey: queuePrefix + ":email",
		smsQueueKey:   queuePrefix + ":sms",
	}
}

// PushEmailNotification đẩy email notification vào Redis queue
func (q *RedisNotificationQueue) PushEmailNotification(ctx context.Context, notification *EmailNotification) error {
	err := q.manager.Push(ctx, q.emailQueueKey, notification)
	if err != nil {
		return err
	}

	log.Info().
		Str("queue", q.emailQueueKey).
		Str("to", notification.To).
		Str("template", notification.TemplateName).
		Msg("Email notification pushed to queue")

	return nil
}

// PushSMSNotification đẩy SMS notification vào Redis queue
func (q *RedisNotificationQueue) PushSMSNotification(ctx context.Context, notification *SMSNotification) error {
	err := q.manager.Push(ctx, q.smsQueueKey, notification)
	if err != nil {
		return err
	}

	log.Info().
		Str("queue", q.smsQueueKey).
		Str("to", notification.To).
		Msg("SMS notification pushed to queue")

	return nil
}

// PopEmailNotification lấy email notification từ queue (blocking với timeout)
func (q *RedisNotificationQueue) PopEmailNotification(ctx context.Context, timeout time.Duration) (*EmailNotification, error) {
	var notification EmailNotification
	err := q.manager.Pop(ctx, q.emailQueueKey, timeout, &notification)
	if err != nil {
		return nil, err
	}

	// Check if struct is empty (meaning no item popped)
	if notification.To == "" && notification.Subject == "" {
		return nil, nil
	}

	return &notification, nil
}

// PopSMSNotification lấy SMS notification từ queue (blocking với timeout)
func (q *RedisNotificationQueue) PopSMSNotification(ctx context.Context, timeout time.Duration) (*SMSNotification, error) {
	var notification SMSNotification
	err := q.manager.Pop(ctx, q.smsQueueKey, timeout, &notification)
	if err != nil {
		return nil, err
	}

	if notification.To == "" && notification.Message == "" {
		return nil, nil
	}

	return &notification, nil
}

// GetQueueLength lấy số lượng message trong queue
func (q *RedisNotificationQueue) GetQueueLength(ctx context.Context, queueName string) (int64, error) {
	return q.manager.Length(ctx, queueName)
}
