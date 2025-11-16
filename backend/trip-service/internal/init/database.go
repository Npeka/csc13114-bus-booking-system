package appinit

import (
	sharedDB "bus-booking/shared/db"
	"bus-booking/trip-service/config"
	"fmt"
)

func InitDatabase(cfg *config.Config) (*sharedDB.DatabaseManager, error) {
	dbManager, err := sharedDB.NewPostgresConnection(&cfg.BaseConfig.Database, cfg.BaseConfig.Server.Environment)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	return dbManager, nil
}
