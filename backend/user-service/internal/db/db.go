package db

import (
	sharedConfig "bus-booking/shared/config"
	sharedDB "bus-booking/shared/db"
	"bus-booking/user-service/config"
)

func NewPostgresConnection(cfg *config.DatabaseConfig, env string) (*sharedDB.DatabaseManager, error) {
	sharedCfg := &sharedConfig.DatabaseConfig{
		Host:            cfg.Host,
		Port:            cfg.Port,
		Name:            cfg.Name,
		Username:        cfg.Username,
		Password:        cfg.Password,
		SSLMode:         cfg.SSLMode,
		TimeZone:        cfg.TimeZone,
		MaxOpenConns:    cfg.MaxOpenConns,
		MaxIdleConns:    cfg.MaxIdleConns,
		ConnMaxLifetime: cfg.ConnMaxLifetime,
		ConnMaxIdleTime: cfg.ConnMaxIdleTime,
	}

	return sharedDB.NewPostgresConnection(sharedCfg, env)
}

// NewRedisConnection creates a new Redis connection
func NewRedisConnection(cfg *config.RedisConfig) (*sharedDB.RedisManager, error) {
	// Convert local config to shared config
	sharedCfg := &sharedConfig.RedisConfig{
		Host:         cfg.Host,
		Port:         cfg.Port,
		Password:     cfg.Password,
		DB:           cfg.DB,
		PoolSize:     cfg.PoolSize,
		MinIdleConns: cfg.MinIdleConns,
		MaxRetries:   cfg.MaxRetries,
		DialTimeout:  cfg.DialTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	}

	return sharedDB.NewRedisConnection(sharedCfg)
}
