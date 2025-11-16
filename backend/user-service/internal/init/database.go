package appinit

import (
	"fmt"

	sharedDB "bus-booking/shared/db"
	"bus-booking/user-service/config"
)

func InitDatabase(cfg *config.Config) (*sharedDB.DatabaseManager, error) {
	dbManager, err := sharedDB.NewPostgresConnection(&cfg.Database, cfg.Server.Environment)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	return dbManager, nil
}
