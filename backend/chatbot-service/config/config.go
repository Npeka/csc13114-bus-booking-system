package config

import (
	sharedConfig "bus-booking/shared/config"
)

type Config struct {
	*sharedConfig.BaseConfig
	Gemini   GeminiConfig   `envPrefix:"GEMINI_"`
	External ExternalConfig `envPrefix:"EXTERNAL_"`
}

type GeminiConfig struct {
	APIKey      string  `env:"API_KEY" envDefault:""`
	Model       string  `env:"MODEL" envDefault:"models/gemini-1.5-flash"`
	Temperature float32 `env:"TEMPERATURE" envDefault:"0.7"`
	MaxTokens   int     `env:"MAX_TOKENS" envDefault:"2048"`
}

type ExternalConfig struct {
	TripServiceURL    string `env:"TRIP_SERVICE_URL" envDefault:"http://localhost:8083"`
	BookingServiceURL string `env:"BOOKING_SERVICE_URL" envDefault:"http://localhost:8082"`
	PaymentServiceURL string `env:"PAYMENT_SERVICE_URL" envDefault:"http://localhost:8084"` // NEW
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
