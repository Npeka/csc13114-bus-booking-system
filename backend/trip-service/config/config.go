package config

import (
	"fmt"

	sharedConfig "bus-booking/shared/config"
)

type Config struct {
	*sharedConfig.BaseConfig
	External ExternalConfig `envPrefix:"EXTERNAL_"`
}

type ExternalConfig struct {
	UserServiceURL    string `env:"USER_SERVICE_URL" envDefault:"http://localhost:8081"`
	BookingServiceURL string `env:"BOOKING_SERVICE_URL" envDefault:"http://localhost:8082"`
}

func (c *Config) IsProduction() bool {
	return c.BaseConfig.Server.Environment == "production"
}

func (c *Config) GetServerAddr() string {
	return fmt.Sprintf("%s:%d", c.BaseConfig.Server.Host, c.BaseConfig.Server.Port)
}

func LoadConfig(envFilePath ...string) (*Config, error) {
	return sharedConfig.LoadConfig[Config](envFilePath...)
}
