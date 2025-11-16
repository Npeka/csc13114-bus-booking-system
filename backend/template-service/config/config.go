package config

import (
	sharedConfig "bus-booking/shared/config"
)

// Config for a simple service that only needs BaseConfig
type Config struct {
	*sharedConfig.BaseConfig
}

// LoadConfig loads configuration for this service
func LoadConfig(envFilePath ...string) (*Config, error) {
	return sharedConfig.LoadConfig[Config](envFilePath...)
}
