package appinit

import (
	"fmt"

	sharedDB "bus-booking/shared/db"
	"bus-booking/user-service/config"
	"bus-booking/user-service/internal/db"
	"bus-booking/user-service/internal/model"
)

func InitDatabase(cfg *config.Config) (*sharedDB.DatabaseManager, error) {
	dbManager, err := db.NewPostgresConnection(&cfg.Database, cfg.Server.Environment)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := dbManager.AutoMigrate(&model.User{}); err != nil {
		return nil, fmt.Errorf("failed to run database migrations: %w", err)
	}

	return dbManager, nil
}
