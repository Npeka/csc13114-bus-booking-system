package logger

import (
	sharedConfig "bus-booking/shared/config"
	"bus-booking/shared/logger"
	"bus-booking/user-service/config"
)

// SetupLogger sets up logger using local config
func SetupLogger(cfg *config.LogConfig) error {
	// Convert local config to shared config
	sharedCfg := &sharedConfig.LogConfig{
		Level:      cfg.Level,
		Format:     cfg.Format,
		Output:     cfg.Output,
		Filename:   cfg.Filename,
		MaxSize:    cfg.MaxSize,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAge,
		Compress:   cfg.Compress,
	}

	return logger.SetupLogger(sharedCfg)
}
