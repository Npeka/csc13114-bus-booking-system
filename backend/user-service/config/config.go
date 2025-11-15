package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

type Config struct {
	Server    ServerConfig    `envPrefix:"SERVER_"`
	Database  DatabaseConfig  `envPrefix:"DATABASE_"`
	Redis     RedisConfig     `envPrefix:"REDIS_"`
	JWT       JWTConfig       `envPrefix:"JWT_"`
	RateLimit RateLimitConfig `envPrefix:"RATE_LIMIT_"`
	CORS      CORSConfig      `envPrefix:"CORS_"`
	Log       LogConfig       `envPrefix:"LOG_"`
	External  ExternalConfig  `envPrefix:"EXTERNAL_"`
	Firebase  FirebaseConfig  `envPrefix:"FIREBASE_"`
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
	Password     string        `env:"PASSWORD" envDefault:""`
	DB           int           `env:"DB" envDefault:"0"`
	PoolSize     int           `env:"POOL_SIZE" envDefault:"10"`
	MinIdleConns int           `env:"MIN_IDLE_CONNS" envDefault:"2"`
	MaxRetries   int           `env:"MAX_RETRIES" envDefault:"3"`
	DialTimeout  time.Duration `env:"DIAL_TIMEOUT" envDefault:"5s"`
	ReadTimeout  time.Duration `env:"READ_TIMEOUT" envDefault:"3s"`
	WriteTimeout time.Duration `env:"WRITE_TIMEOUT" envDefault:"3s"`
}

type JWTConfig struct {
	SecretKey        string        `env:"SECRET_KEY"`
	RefreshSecretKey string        `env:"REFRESH_SECRET_KEY"`
	AccessTokenTTL   time.Duration `env:"ACCESS_TOKEN_TTL" envDefault:"15m"`
	RefreshTokenTTL  time.Duration `env:"REFRESH_TOKEN_TTL" envDefault:"168h"`
	Issuer           string        `env:"ISSUER" envDefault:"bus-booking-system"`
	Audience         string        `env:"AUDIENCE" envDefault:"bus-booking-users"`
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

type ExternalConfig struct {
	PaymentServiceURL string        `env:"PAYMENT_SERVICE_URL" envDefault:"https://api.payment-service.com"`
	NotificationURL   string        `env:"NOTIFICATION_URL" envDefault:"https://api.notification-service.com"`
	Timeout           time.Duration `env:"TIMEOUT" envDefault:"30s"`
	RetryAttempts     int           `env:"RETRY_ATTEMPTS" envDefault:"3"`
}

type FirebaseConfig struct {
	ProjectID    string `env:"PROJECT_ID"`
	PrivateKeyID string `env:"PRIVATE_KEY_ID"`
	PrivateKey   string `env:"PRIVATE_KEY"`
	ClientEmail  string `env:"CLIENT_EMAIL"`
	ClientID     string `env:"CLIENT_ID"`
}

// LoadConfig loads configuration from .env file and environment variables
func LoadConfig(envFilePath ...string) (*Config, error) {
	// Default .env file path
	envFile := "config/.env"
	if len(envFilePath) > 0 && envFilePath[0] != "" {
		envFile = envFilePath[0]
	}

	// Load .env file if exists
	if err := godotenv.Load(envFile); err != nil {
		log.Debug().Str("file", envFile).Msg("No .env file found or error loading it, using environment variables")
	} else {
		log.Info().Str("file", envFile).Msg("Loaded configuration from .env file")
	}

	// Parse environment variables into struct
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("failed to parse environment variables: %w", err)
	}

	// Validate required fields
	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	return cfg, nil
}

// validate validates the configuration
func (c *Config) validate() error {
	if c.JWT.SecretKey == "" {
		return fmt.Errorf("JWT_SECRET_KEY is required")
	}
	if len(c.JWT.SecretKey) < 32 {
		return fmt.Errorf("JWT_SECRET_KEY must be at least 32 characters long")
	}
	if c.JWT.RefreshSecretKey == "" {
		return fmt.Errorf("JWT_REFRESH_SECRET_KEY is required")
	}
	if len(c.JWT.RefreshSecretKey) < 32 {
		return fmt.Errorf("JWT_REFRESH_SECRET_KEY must be at least 32 characters long")
	}
	return nil
}

// IsProduction checks if the environment is production
func (c *Config) IsProduction() bool {
	return strings.ToLower(c.Server.Environment) == "production"
}

// IsDevelopment checks if the environment is development
func (c *Config) IsDevelopment() bool {
	return strings.ToLower(c.Server.Environment) == "development"
}

// GetServerAddr returns the server address
func (c *Config) GetServerAddr() string {
	return fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port)
}

// GetDSN returns the PostgreSQL connection string
func (c *Config) GetDSN() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
		c.Database.Host,
		c.Database.Username,
		c.Database.Password,
		c.Database.Name,
		c.Database.Port,
		c.Database.SSLMode,
		c.Database.TimeZone,
	)
}

// ToSharedLogConfig converts local LogConfig to shared LogConfig for compatibility
func (c *LogConfig) ToSharedLogConfig() interface{} {
	// Since we can't import shared config due to circular dependency,
	// we'll create a compatible struct with the same field structure
	return struct {
		Level      string `mapstructure:"level"`
		Format     string `mapstructure:"format"`
		Output     string `mapstructure:"output"`
		Filename   string `mapstructure:"filename"`
		MaxSize    int    `mapstructure:"max_size"`
		MaxBackups int    `mapstructure:"max_backups"`
		MaxAge     int    `mapstructure:"max_age"`
		Compress   bool   `mapstructure:"compress"`
	}{
		Level:      c.Level,
		Format:     c.Format,
		Output:     c.Output,
		Filename:   c.Filename,
		MaxSize:    c.MaxSize,
		MaxBackups: c.MaxBackups,
		MaxAge:     c.MaxAge,
		Compress:   c.Compress,
	}
}

// ToSharedCORSConfig converts local CORSConfig to shared CORSConfig for compatibility
func (c *CORSConfig) ToSharedCORSConfig() interface{} {
	// Since we can't import shared config due to circular dependency,
	// we'll create a compatible struct with the same field structure
	return struct {
		AllowOrigins     []string `mapstructure:"allow_origins"`
		AllowMethods     []string `mapstructure:"allow_methods"`
		AllowHeaders     []string `mapstructure:"allow_headers"`
		ExposeHeaders    []string `mapstructure:"expose_headers"`
		AllowCredentials bool     `mapstructure:"allow_credentials"`
		MaxAge           int      `mapstructure:"max_age"`
	}{
		AllowOrigins:     c.AllowOrigins,
		AllowMethods:     c.AllowMethods,
		AllowHeaders:     c.AllowHeaders,
		ExposeHeaders:    c.ExposeHeaders,
		AllowCredentials: c.AllowCredentials,
		MaxAge:           c.MaxAge,
	}
}
