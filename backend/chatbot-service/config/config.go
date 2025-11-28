package config

import (
	sharedConfig "bus-booking/shared/config"
)

type Config struct {
	*sharedConfig.BaseConfig
	OpenAI   OpenAIConfig   `envPrefix:"OPENAI_"`
	External ExternalConfig `envPrefix:"EXTERNAL_"`
}

type OpenAIConfig struct {
	APIKey      string  `env:"API_KEY" envDefault:""`
	Model       string  `env:"MODEL" envDefault:"gpt-4-turbo-preview"`
	Temperature float32 `env:"TEMPERATURE" envDefault:"0.7"`
	MaxTokens   int     `env:"MAX_TOKENS" envDefault:"500"`
}

type ExternalConfig struct {
	TripServiceURL    string `env:"TRIP_SERVICE_URL" envDefault:"http://localhost:8083"`
	BookingServiceURL string `env:"BOOKING_SERVICE_URL" envDefault:"http://localhost:8082"`
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
