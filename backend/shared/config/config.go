package config

import (
	"fmt"
	"reflect"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

type BaseConfig struct {
	Server    ServerConfig    `envPrefix:"SERVER_"`
	Database  DatabaseConfig  `envPrefix:"DATABASE_"`
	Redis     RedisConfig     `envPrefix:"REDIS_"`
	RateLimit RateLimitConfig `envPrefix:"RATE_LIMIT_"`
	CORS      CORSConfig      `envPrefix:"CORS_"`
	Log       LogConfig       `envPrefix:"LOG_"`
}

type ServerConfig struct {
	Port            int           `env:"PORT" envDefault:"8080"`
	Host            string        `env:"HOST" envDefault:"0.0.0.0"`
	Environment     string        `env:"ENVIRONMENT" envDefault:"development"`
	ReadTimeout     time.Duration `env:"READ_TIMEOUT" envDefault:"10s"`
	WriteTimeout    time.Duration `env:"WRITE_TIMEOUT" envDefault:"10s"`
	IdleTimeout     time.Duration `env:"IDLE_TIMEOUT" envDefault:"120s"`
	ShutdownTimeout time.Duration `env:"SHUTDOWN_TIMEOUT" envDefault:"30s"`
	MaxHeaderBytes  int           `env:"MAX_HEADER_BYTES" envDefault:"1048576"`
	IsProduction    bool          `env:"-"`
}

type DatabaseConfig struct {
	Host            string        `env:"HOST" envDefault:"localhost"`
	Port            int           `env:"PORT" envDefault:"5432"`
	Name            string        `env:"NAME" envDefault:"bus_booking"`
	Username        string        `env:"USERNAME" envDefault:"postgres"`
	Password        string        `env:"PASSWORD" envDefault:"postgres"`
	SSLMode         string        `env:"SSL_MODE" envDefault:"disable"`
	TimeZone        string        `env:"TIMEZONE" envDefault:"Asia/Ho_Chi_Minh"`
	MaxOpenConns    int           `env:"MAX_OPEN_CONNS" envDefault:"25"`
	MaxIdleConns    int           `env:"MAX_IDLE_CONNS" envDefault:"5"`
	ConnMaxLifetime time.Duration `env:"CONN_MAX_LIFETIME" envDefault:"300s"`
	ConnMaxIdleTime time.Duration `env:"CONN_MAX_IDLE_TIME" envDefault:"60s"`
}

type RedisConfig struct {
	Host         string        `env:"HOST" envDefault:"localhost"`
	Port         int           `env:"PORT" envDefault:"6379"`
	Username     string        `env:"USERNAME" envDefault:""`
	Password     string        `env:"PASSWORD" envDefault:""`
	DB           int           `env:"DB" envDefault:"0"`
	TLS          bool          `env:"TLS" envDefault:"false"`
	PoolSize     int           `env:"POOL_SIZE" envDefault:"10"`
	MinIdleConns int           `env:"MIN_IDLE_CONNS" envDefault:"2"`
	MaxRetries   int           `env:"MAX_RETRIES" envDefault:"3"`
	DialTimeout  time.Duration `env:"DIAL_TIMEOUT" envDefault:"60s"`
	ReadTimeout  time.Duration `env:"READ_TIMEOUT" envDefault:"60s"`
	WriteTimeout time.Duration `env:"WRITE_TIMEOUT" envDefault:"60s"`
}

type RateLimitConfig struct {
	RPS          int           `env:"RPS" envDefault:"100"`
	Burst        int           `env:"BURST" envDefault:"200"`
	Period       time.Duration `env:"PERIOD" envDefault:"1m"`
	CleanupAfter time.Duration `env:"CLEANUP_AFTER" envDefault:"10m"`
}

type CORSConfig struct {
	AllowOrigins     []string `env:"ALLOW_ORIGINS" envSeparator:"," envDefault:"*"`
	AllowMethods     []string `env:"ALLOW_METHODS" envSeparator:"," envDefault:"GET,POST,PUT,DELETE,OPTIONS"`
	AllowHeaders     []string `env:"ALLOW_HEADERS" envSeparator:"," envDefault:"Origin,Content-Type,Accept,Authorization,X-Requested-With"`
	ExposeHeaders    []string `env:"EXPOSE_HEADERS" envSeparator:","`
	AllowCredentials bool     `env:"ALLOW_CREDENTIALS" envDefault:"true"`
	MaxAge           int      `env:"MAX_AGE" envDefault:"86400"`
}

type LogConfig struct {
	Level      string `env:"LEVEL" envDefault:"info"`
	Format     string `env:"FORMAT" envDefault:"json"`
	Output     string `env:"OUTPUT" envDefault:"stdout"`
	Filename   string `env:"FILENAME" envDefault:"logs/app.log"`
	MaxSize    int    `env:"MAX_SIZE" envDefault:"100"`
	MaxBackups int    `env:"MAX_BACKUPS" envDefault:"3"`
	MaxAge     int    `env:"MAX_AGE" envDefault:"30"`
	Compress   bool   `env:"COMPRESS" envDefault:"true"`
}

type JWTConfig struct {
	SecretKey        string        `env:"SECRET_KEY"`
	RefreshSecretKey string        `env:"REFRESH_SECRET_KEY"`
	AccessTokenTTL   time.Duration `env:"ACCESS_TOKEN_TTL" envDefault:"15m"`
	RefreshTokenTTL  time.Duration `env:"REFRESH_TOKEN_TTL" envDefault:"168h"`
	Issuer           string        `env:"ISSUER" envDefault:"bus-booking-system"`
	Audience         string        `env:"AUDIENCE" envDefault:"bus-booking-users"`
}

type ExternalConfig struct {
	PaymentServiceURL string        `env:"PAYMENT_SERVICE_URL" envDefault:"https://api.payment-service.com"`
	NotificationURL   string        `env:"NOTIFICATION_URL" envDefault:"https://api.notification-service.com"`
	Timeout           time.Duration `env:"TIMEOUT" envDefault:"30s"`
	RetryAttempts     int           `env:"RETRY_ATTEMPTS" envDefault:"3"`
}

func (c *BaseConfig) GetServerAddr() string {
	return fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port)
}

func IsProduction(environment string) bool {
	return environment == "production"
}

func LoadConfig[T any](envFilePath ...string) (*T, error) {
	envFile := "./.env"
	if len(envFilePath) > 0 && envFilePath[0] != "" {
		envFile = envFilePath[0]
	}

	if err := godotenv.Load(envFile); err != nil {
		log.Debug().Str("file", envFile).Msg("No .env file found or error loading it, using environment variables")
	} else {
		log.Info().Str("file", envFile).Msg("Loaded configuration from .env file")
	}

	baseConfig := &BaseConfig{}
	if err := env.Parse(baseConfig); err != nil {
		return nil, fmt.Errorf("failed to parse base config environment variables: %w", err)
	}

	baseConfig.Server.IsProduction = IsProduction(baseConfig.Server.Environment)

	cfg := new(T)
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("failed to parse service config environment variables: %w", err)
	}

	if err := setBaseConfig(cfg, baseConfig); err != nil {
		return nil, fmt.Errorf("failed to set base config: %w", err)
	}

	return cfg, nil
}

func setBaseConfig(serviceConfig interface{}, baseConfig *BaseConfig) error {
	v := reflect.ValueOf(serviceConfig).Elem()
	t := v.Type()

	// Find the embedded BaseConfig field
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		// Check if this field is a pointer to BaseConfig
		if field.Type() == reflect.TypeOf((*BaseConfig)(nil)) {
			if !field.CanSet() {
				return fmt.Errorf("cannot set BaseConfig field %s", fieldType.Name)
			}
			field.Set(reflect.ValueOf(baseConfig))
			return nil
		}
	}

	return fmt.Errorf("no embedded *BaseConfig field found in service config")
}
