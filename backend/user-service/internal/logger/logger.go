package logger

import (
	sharedConfig "bus-booking/shared/config"
	"bus-booking/shared/logger"
	"bus-booking/user-service/config"
)

func SetupLogger(cfg *config.LogConfig) error {
	return logger.SetupLogger(&sharedConfig.LogConfig{
		Level:      cfg.Level,
		Format:     cfg.Format,
		Output:     cfg.Output,
		Filename:   cfg.Filename,
		MaxSize:    cfg.MaxSize,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAge,
		Compress:   cfg.Compress,
	})
}
