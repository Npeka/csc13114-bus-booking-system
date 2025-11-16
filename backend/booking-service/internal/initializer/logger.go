package initializer

import (
	"bus-booking/shared/logger"
	"bus-booking/trip-service/config"
)

func SetupLogger(cfg *config.Config) error {
	return logger.SetupLogger(&cfg.BaseConfig.Log)
}
