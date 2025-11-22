package initializer

import (
	"bus-booking/payment-service/config"
	"bus-booking/shared/logger"
)

func SetupLogger(cfg *config.Config) error {
	return logger.SetupLogger(&cfg.BaseConfig.Log)
}
