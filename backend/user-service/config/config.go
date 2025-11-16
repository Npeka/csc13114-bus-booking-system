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
	ServiceAccountKeyPath string `env:"SERVICE_ACCOUNT_KEY_PATH" envDefault:"config/csc13114-bus-booking-system-firebase-adminsdk-fbsvc.json"`
	ProjectID             string `env:"PROJECT_ID" envDefault:"csc13114-bus-booking-system"`
}

func LoadConfig(envFilePath ...string) (*Config, error) {
	return sharedConfig.LoadConfig[Config](envFilePath...)
}
