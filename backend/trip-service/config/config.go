package config

import (
	sharedConfig "bus-booking/shared/config"
)

type Config struct {
	*sharedConfig.BaseConfig
	External ExternalConfig `envPrefix:"EXTERNAL_"`
	Storage  sharedConfig.StorageConfig
}

type ExternalConfig struct {
	BookingServiceURL string `env:"BOOKING_SERVICE_URL" envDefault:"http://localhost:8080"`
	UserServiceURL    string `env:"USER_SERVICE_URL" envDefault:"http://localhost:8083"`
}

func LoadConfig(envFilePath ...string) (*Config, error) {
	return sharedConfig.LoadConfig[Config](envFilePath...)
}

func MustLoadConfig(envFilePath ...string) *Config {
	cfg, err := LoadConfig(envFilePath...)
	if err != nil {
		panic(err)
	}
	return cfg
}
