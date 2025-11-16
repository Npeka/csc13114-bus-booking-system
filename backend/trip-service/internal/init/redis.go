package appinit

import (
	"fmt"

	sharedDB "bus-booking/shared/db"
	"bus-booking/trip-service/config"
)

func InitRedis(cfg *config.Config) (*sharedDB.RedisManager, error) {
	redisManager, err := sharedDB.NewRedisConnection(&cfg.BaseConfig.Redis)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}
	return redisManager, nil
}
