package appinit

import (
	"bus-booking/shared/logger"
	"bus-booking/user-service/config"
)

func SetupLogger(cfg *config.Config) error {
	return logger.SetupLogger(&cfg.Log)
}
