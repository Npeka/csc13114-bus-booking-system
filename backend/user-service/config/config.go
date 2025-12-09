package config

import (
	"time"

	sharedConfig "bus-booking/shared/config"
)

type Config struct {
	*sharedConfig.BaseConfig
	JWT      JWTConfig                `envPrefix:"JWT_"`
	Redis    sharedConfig.RedisConfig `envPrefix:"REDIS_"`
	Firebase FirebaseConfig           `envPrefix:"FIREBASE_"`
	External ExternalConfig           `envPrefix:"EXTERNAL_"`
}

type JWTConfig struct {
	SecretKey        string        `env:"SECRET_KEY"`
	RefreshSecretKey string        `env:"REFRESH_SECRET_KEY"`
	AccessTokenTTL   time.Duration `env:"ACCESS_TOKEN_TTL" envDefault:"15m"`
	RefreshTokenTTL  time.Duration `env:"REFRESH_TOKEN_TTL" envDefault:"168h"`
	Issuer           string        `env:"ISSUER" envDefault:"bus-booking-system"`
	Audience         string        `env:"AUDIENCE" envDefault:"bus-booking-users"`
}

type FirebaseConfig struct {
	ServiceAccountKeyPath string `env:"SERVICE_ACCOUNT_KEY_PATH" envDefault:"config/fbsvc.json"`
	ProjectID             string `env:"PROJECT_ID" envDefault:"csc13114-bus-booking-system"`
}

type ExternalConfig struct {
	NotificationServiceURL string `env:"NOTIFICATION_SERVICE_URL" envDefault:"http://localhost:8085"`
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
