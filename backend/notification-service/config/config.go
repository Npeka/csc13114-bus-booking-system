package config

import (
	sharedConfig "bus-booking/shared/config"
)

type Config struct {
	*sharedConfig.BaseConfig
	External     ExternalConfig `envPrefix:"EXTERNAL_"`
	SMTP         SMTPConfig     `envPrefix:"SMTP_"`
	TemplatePath string         `env:"TEMPLATE_PATH" envDefault:"templates"`
	LogoURL      string         `env:"LOGO_URL" envDefault:"https://csc13114-bus-booking-system.vercel.app/_next/image?url=%2Ffavicon.png&w=128&q=75"`
}

type ExternalConfig struct {
	BookingServiceURL string `env:"BOOKING_SERVICE_URL" envDefault:"http://localhost:8082"`
}

type SMTPConfig struct {
	Host     string `env:"HOST" envDefault:"smtp.gmail.com"`
	Port     int    `env:"PORT" envDefault:"587"`
	Email    string `env:"EMAIL"`
	Password string `env:"PASSWORD"`
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
