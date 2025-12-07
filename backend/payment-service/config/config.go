package config

import (
	sharedConfig "bus-booking/shared/config"
)

type Config struct {
	*sharedConfig.BaseConfig
	External ExternalConfig `envPrefix:"EXTERNAL_"`
	PayOS    PayOSConfig    `envPrefix:"PAYOS_"`
}

type ExternalConfig struct {
	BookingServiceURL string `env:"BOOKING_SERVICE_URL" envDefault:"http://localhost:8082"`
}

type PayOSConfig struct {
	ClientID    string `env:"CLIENT_ID" envDefault:""`
	APIKey      string `env:"API_KEY" envDefault:""`
	ChecksumKey string `env:"CHECKSUM_KEY" envDefault:""`
	ReturnURL   string `env:"RETURN_URL" envDefault:"http://localhost:3000/payment/success"`
	CancelURL   string `env:"CANCEL_URL" envDefault:"http://localhost:3000/payment/cancel"`
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
