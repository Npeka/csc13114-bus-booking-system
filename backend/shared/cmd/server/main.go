package main

import (
	"flag"
	"fmt"
	"log"

	sharedConfig "bus-booking/shared/config"
	sharedDB "bus-booking/shared/db"
)

func main() {
	var (
		envFile = flag.String("env", ".env", "Path to environment file")
		service = flag.String("service", "", "Service name for migration")
	)
	flag.Parse()

	if *service == "" {
		log.Fatal("Service name is required. Use -service flag")
	}

	// Load configuration
	cfg, err := sharedConfig.LoadConfig[sharedConfig.BaseConfig](*envFile)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Create migration manager
	migrationManager, err := sharedDB.NewMigrationManager(&cfg.Database, cfg.Server.Environment)
	if err != nil {
		log.Fatalf("Failed to create migration manager: %v", err)
	}
	defer migrationManager.Close()

	// Get models for specific service
	models, err := getModelsForService(*service)
	if err != nil {
		log.Fatalf("Failed to get models for service %s: %v", *service, err)
	}

	// Run migrations
	if err := migrationManager.RunMigrations(models...); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	log.Printf("Migration completed successfully for service: %s", *service)
}

// getModelsForService returns models for a specific service
func getModelsForService(serviceName string) ([]interface{}, error) {
	switch serviceName {
	case "user-service":
		return getUserServiceModels()
	case "booking-service":
		return getBookingServiceModels()
	case "trip-service":
		return getTripServiceModels()
	case "payment-service":
		return getPaymentServiceModels()
	default:
		return nil, fmt.Errorf("unknown service: %s", serviceName)
	}
}

// Service-specific model functions
func getUserServiceModels() ([]interface{}, error) {
	// Import user service models dynamically or return error
	fmt.Println("To add user-service models, import the model package and return the models")
	return nil, fmt.Errorf("user-service models not implemented in shared migration")
}

func getBookingServiceModels() ([]interface{}, error) {
	fmt.Println("To add booking-service models, import the model package and return the models")
	return nil, fmt.Errorf("booking-service models not implemented in shared migration")
}

func getTripServiceModels() ([]interface{}, error) {
	fmt.Println("To add trip-service models, import the model package and return the models")
	return nil, fmt.Errorf("trip-service models not implemented in shared migration")
}

func getPaymentServiceModels() ([]interface{}, error) {
	fmt.Println("To add payment-service models, import the model package and return the models")
	return nil, fmt.Errorf("payment-service models not implemented in shared migration")
}
