package initializer

import (
	"bus-booking/booking-service/config"
	sharedDB "bus-booking/shared/db"
	"fmt"
)

func InitDatabase(cfg *config.Config) (*sharedDB.DatabaseManager, error) {
	dbManager, err := sharedDB.NewPostgresConnection(&cfg.BaseConfig.Database)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	return dbManager, nil
}
