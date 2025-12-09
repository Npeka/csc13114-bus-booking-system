package config

import (
	sharedConfig "bus-booking/shared/config"
)

type Config struct {
	*sharedConfig.BaseConfig
	External ExternalConfig `envPrefix:"EXTERNAL_"`
}

type ExternalConfig struct {
	PaymentServiceURL      string `env:"PAYMENT_SERVICE_URL" envDefault:"http://localhost:8081"`
	TripServiceURL         string `env:"TRIP_SERVICE_URL" envDefault:"http://localhost:8082"`
	UserServiceURL         string `env:"USER_SERVICE_URL" envDefault:"http://localhost:8083"`
	NotificationServiceURL string `env:"NOTIFICATION_SERVICE_URL" envDefault:"http://localhost:8085"`
}

func LoadConfig(envFilePath ...string) (*Config, error) {
	return sharedConfig.LoadConfig[Config](envFilePath...)
}

func MustLoadConfig(envFilePath ...string) *Config {
	return sharedConfig.MustLoadConfig[Config](envFilePath...)
}
