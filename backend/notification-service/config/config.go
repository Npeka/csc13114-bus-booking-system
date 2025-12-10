package config

import (
	sharedConfig "bus-booking/shared/config"
)

type Config struct {
	*sharedConfig.BaseConfig
	External     ExternalConfig `envPrefix:"EXTERNAL_"`
	BrevoAPIKey  string         `env:"BREVO_API_KEY,required"`
	FromEmail    string         `env:"FROM_EMAIL,required"`
	FromName     string         `env:"FROM_NAME" envDefault:"Bus Booking System"`
	TemplatePath string         `env:"TEMPLATE_PATH" envDefault:"templates"`
	LogoURL      string         `env:"LOGO_URL" envDefault:"https://csc13114-bus-booking-system.vercel.app/_next/image?url=%2Ffavicon.png&w=128&q=75"`
}

type ExternalConfig struct {
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
