package config

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type Config struct {
	Server    ServerConfig    `mapstructure:"server"`
	Database  DatabaseConfig  `mapstructure:"database"`
	Redis     RedisConfig     `mapstructure:"redis"`
	JWT       JWTConfig       `mapstructure:"jwt"`
	OAuth     OAuthConfig     `mapstructure:"oauth"`
	RateLimit RateLimitConfig `mapstructure:"rate_limit"`
	CORS      CORSConfig      `mapstructure:"cors"`
	Log       LogConfig       `mapstructure:"log"`
	External  ExternalConfig  `mapstructure:"external"`
}

type ServerConfig struct {
	Port            int           `mapstructure:"port" env:"SERVER_PORT" validate:"min=1,max=65535" default:"8080"`
	Host            string        `mapstructure:"host" env:"SERVER_HOST" default:"0.0.0.0"`
	Environment     string        `mapstructure:"environment" env:"ENVIRONMENT" validate:"oneof=development staging production" default:"development"`
	ReadTimeout     time.Duration `mapstructure:"read_timeout" env:"SERVER_READ_TIMEOUT" default:"10s"`
	WriteTimeout    time.Duration `mapstructure:"write_timeout" env:"SERVER_WRITE_TIMEOUT" default:"10s"`
	IdleTimeout     time.Duration `mapstructure:"idle_timeout" env:"SERVER_IDLE_TIMEOUT" default:"120s"`
	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout" env:"SERVER_SHUTDOWN_TIMEOUT" default:"30s"`
	MaxHeaderBytes  int           `mapstructure:"max_header_bytes" env:"SERVER_MAX_HEADER_BYTES" default:"1048576"`
}

type DatabaseConfig struct {
	Host            string        `mapstructure:"host" env:"DB_HOST" validate:"required" default:"localhost"`
	Port            int           `mapstructure:"port" env:"DB_PORT" validate:"min=1,max=65535" default:"5432"`
	Name            string        `mapstructure:"name" env:"DB_NAME" validate:"required" default:"bus_booking"`
	Username        string        `mapstructure:"username" env:"DB_USERNAME" validate:"required" default:"postgres"`
	Password        string        `mapstructure:"password" env:"DB_PASSWORD" validate:"required" default:"password"`
	SSLMode         string        `mapstructure:"ssl_mode" env:"DB_SSL_MODE" validate:"oneof=disable require verify-ca verify-full" default:"disable"`
	TimeZone        string        `mapstructure:"timezone" env:"DB_TIMEZONE" default:"Asia/Ho_Chi_Minh"`
	MaxOpenConns    int           `mapstructure:"max_open_conns" env:"DB_MAX_OPEN_CONNS" validate:"min=1" default:"25"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns" env:"DB_MAX_IDLE_CONNS" validate:"min=1" default:"5"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime" env:"DB_CONN_MAX_LIFETIME" default:"300s"`
	ConnMaxIdleTime time.Duration `mapstructure:"conn_max_idle_time" env:"DB_CONN_MAX_IDLE_TIME" default:"60s"`
}

type RedisConfig struct {
	Host         string        `mapstructure:"host" env:"REDIS_HOST" validate:"required" default:"localhost"`
	Port         int           `mapstructure:"port" env:"REDIS_PORT" validate:"min=1,max=65535" default:"6379"`
	Password     string        `mapstructure:"password" env:"REDIS_PASSWORD" default:""`
	DB           int           `mapstructure:"db" env:"REDIS_DB" validate:"min=0,max=15" default:"0"`
	PoolSize     int           `mapstructure:"pool_size" env:"REDIS_POOL_SIZE" validate:"min=1" default:"10"`
	MinIdleConns int           `mapstructure:"min_idle_conns" env:"REDIS_MIN_IDLE_CONNS" validate:"min=0" default:"2"`
	MaxRetries   int           `mapstructure:"max_retries" env:"REDIS_MAX_RETRIES" validate:"min=0" default:"3"`
	DialTimeout  time.Duration `mapstructure:"dial_timeout" env:"REDIS_DIAL_TIMEOUT" default:"5s"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout" env:"REDIS_READ_TIMEOUT" default:"3s"`
	WriteTimeout time.Duration `mapstructure:"write_timeout" env:"REDIS_WRITE_TIMEOUT" default:"3s"`
}

type JWTConfig struct {
	SecretKey        string        `mapstructure:"secret_key" env:"JWT_SECRET_KEY" validate:"required,min=32"`
	RefreshSecretKey string        `mapstructure:"refresh_secret_key" env:"JWT_REFRESH_SECRET_KEY" validate:"required,min=32"`
	AccessTokenTTL   time.Duration `mapstructure:"access_token_ttl" env:"JWT_ACCESS_TOKEN_TTL" default:"15m"`
	RefreshTokenTTL  time.Duration `mapstructure:"refresh_token_ttl" env:"JWT_REFRESH_TOKEN_TTL" default:"168h"`
	Issuer           string        `mapstructure:"issuer" env:"JWT_ISSUER" default:"bus-booking-system"`
	Audience         string        `mapstructure:"audience" env:"JWT_AUDIENCE" default:"bus-booking-users"`
}

type RateLimitConfig struct {
	RPS          int           `mapstructure:"rps" env:"RATE_LIMIT_RPS" validate:"min=1" default:"100"`
	Burst        int           `mapstructure:"burst" env:"RATE_LIMIT_BURST" validate:"min=1" default:"200"`
	Period       time.Duration `mapstructure:"period" env:"RATE_LIMIT_PERIOD" default:"1m"`
	CleanupAfter time.Duration `mapstructure:"cleanup_after" env:"RATE_LIMIT_CLEANUP_AFTER" default:"10m"`
}

type CORSConfig struct {
	AllowOrigins     []string `mapstructure:"allow_origins" env:"CORS_ALLOW_ORIGINS" envSeparator:"," default:"*"`
	AllowMethods     []string `mapstructure:"allow_methods" env:"CORS_ALLOW_METHODS" envSeparator:"," default:"GET,POST,PUT,DELETE,OPTIONS"`
	AllowHeaders     []string `mapstructure:"allow_headers" env:"CORS_ALLOW_HEADERS" envSeparator:"," default:"Origin,Content-Type,Accept,Authorization,X-Requested-With"`
	ExposeHeaders    []string `mapstructure:"expose_headers" env:"CORS_EXPOSE_HEADERS" envSeparator:","`
	AllowCredentials bool     `mapstructure:"allow_credentials" env:"CORS_ALLOW_CREDENTIALS" default:"true"`
	MaxAge           int      `mapstructure:"max_age" env:"CORS_MAX_AGE" default:"86400"`
}

type LogConfig struct {
	Level      string `mapstructure:"level" env:"LOG_LEVEL" validate:"oneof=trace debug info warn error fatal panic" default:"info"`
	Format     string `mapstructure:"format" env:"LOG_FORMAT" validate:"oneof=json text" default:"json"`
	Output     string `mapstructure:"output" env:"LOG_OUTPUT" validate:"oneof=stdout stderr file" default:"stdout"`
	Filename   string `mapstructure:"filename" env:"LOG_FILENAME" default:"app.log"`
	MaxSize    int    `mapstructure:"max_size" env:"LOG_MAX_SIZE" default:"100"`
	MaxBackups int    `mapstructure:"max_backups" env:"LOG_MAX_BACKUPS" default:"3"`
	MaxAge     int    `mapstructure:"max_age" env:"LOG_MAX_AGE" default:"30"`
	Compress   bool   `mapstructure:"compress" env:"LOG_COMPRESS" default:"true"`
}

type OAuthConfig struct {
	Firebase FirebaseConfig `mapstructure:"firebase"`
}

type FirebaseConfig struct {
	ProjectID           string `mapstructure:"project_id" env:"FIREBASE_PROJECT_ID" validate:"required"`
	PrivateKeyID        string `mapstructure:"private_key_id" env:"FIREBASE_PRIVATE_KEY_ID" validate:"required"`
	PrivateKey          string `mapstructure:"private_key" env:"FIREBASE_PRIVATE_KEY" validate:"required"`
	ClientEmail         string `mapstructure:"client_email" env:"FIREBASE_CLIENT_EMAIL" validate:"required,email"`
	ClientID            string `mapstructure:"client_id" env:"FIREBASE_CLIENT_ID" validate:"required"`
	AuthURI             string `mapstructure:"auth_uri" env:"FIREBASE_AUTH_URI" default:"https://accounts.google.com/o/oauth2/auth"`
	TokenURI            string `mapstructure:"token_uri" env:"FIREBASE_TOKEN_URI" default:"https://oauth2.googleapis.com/token"`
	AuthProviderCertURL string `mapstructure:"auth_provider_x509_cert_url" env:"FIREBASE_AUTH_PROVIDER_CERT_URL" default:"https://www.googleapis.com/oauth2/v1/certs"`
	ClientCertURL       string `mapstructure:"client_x509_cert_url" env:"FIREBASE_CLIENT_CERT_URL" validate:"required"`
}

type ExternalConfig struct {
	PaymentServiceURL string        `mapstructure:"payment_service_url" env:"PAYMENT_SERVICE_URL" validate:"url"`
	NotificationURL   string        `mapstructure:"notification_url" env:"NOTIFICATION_URL" validate:"url"`
	Timeout           time.Duration `mapstructure:"timeout" env:"EXTERNAL_TIMEOUT" default:"30s"`
	RetryAttempts     int           `mapstructure:"retry_attempts" env:"EXTERNAL_RETRY_ATTEMPTS" validate:"min=0" default:"3"`
}

// LoadConfig loads configuration from environment variables, config files, and defaults
func LoadConfig() (*Config, error) {
	cfg := &Config{}

	// Load .env file if exists
	if err := godotenv.Load(); err != nil {
		log.Debug().Msg("No .env file found or error loading it, continuing with environment variables")
	}

	// Parse environment variables into struct
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("failed to parse environment variables: %w", err)
	}

	// Setup Viper for config files
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("./configs")
	viper.AutomaticEnv()

	// Set default values
	setDefaults()

	// Read config file if exists
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
		log.Debug().Msg("No config file found, using environment variables and defaults")
	}

	// Unmarshal config file into struct (this will override env values if config file has values)
	if err := viper.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return cfg, nil
}

// setDefaults sets default values for Viper
func setDefaults() {
	// Server defaults
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("server.environment", "development")
	viper.SetDefault("server.read_timeout", "10s")
	viper.SetDefault("server.write_timeout", "10s")
	viper.SetDefault("server.idle_timeout", "120s")
	viper.SetDefault("server.shutdown_timeout", "30s")
	viper.SetDefault("server.max_header_bytes", 1048576)

	// Database defaults
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 5432)
	viper.SetDefault("database.name", "bus_booking")
	viper.SetDefault("database.username", "postgres")
	viper.SetDefault("database.password", "password")
	viper.SetDefault("database.ssl_mode", "disable")
	viper.SetDefault("database.timezone", "Asia/Ho_Chi_Minh")
	viper.SetDefault("database.max_open_conns", 25)
	viper.SetDefault("database.max_idle_conns", 5)
	viper.SetDefault("database.conn_max_lifetime", "300s")
	viper.SetDefault("database.conn_max_idle_time", "60s")

	// Redis defaults
	viper.SetDefault("redis.host", "localhost")
	viper.SetDefault("redis.port", 6379)
	viper.SetDefault("redis.db", 0)
	viper.SetDefault("redis.pool_size", 10)
	viper.SetDefault("redis.min_idle_conns", 2)
	viper.SetDefault("redis.max_retries", 3)
	viper.SetDefault("redis.dial_timeout", "5s")
	viper.SetDefault("redis.read_timeout", "3s")
	viper.SetDefault("redis.write_timeout", "3s")

	// JWT defaults
	viper.SetDefault("jwt.access_token_ttl", "15m")
	viper.SetDefault("jwt.refresh_token_ttl", "168h")
	viper.SetDefault("jwt.issuer", "bus-booking-system")
	viper.SetDefault("jwt.audience", "bus-booking-users")

	// Rate limit defaults
	viper.SetDefault("rate_limit.rps", 100)
	viper.SetDefault("rate_limit.burst", 200)
	viper.SetDefault("rate_limit.period", "1m")
	viper.SetDefault("rate_limit.cleanup_after", "10m")

	// CORS defaults
	viper.SetDefault("cors.allow_origins", []string{"*"})
	viper.SetDefault("cors.allow_methods", []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	viper.SetDefault("cors.allow_headers", []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"})
	viper.SetDefault("cors.allow_credentials", true)
	viper.SetDefault("cors.max_age", 86400)

	// Log defaults
	viper.SetDefault("log.level", "info")
	viper.SetDefault("log.format", "json")
	viper.SetDefault("log.output", "stdout")
	viper.SetDefault("log.filename", "app.log")
	viper.SetDefault("log.max_size", 100)
	viper.SetDefault("log.max_backups", 3)
	viper.SetDefault("log.max_age", 30)
	viper.SetDefault("log.compress", true)

	// External defaults
	viper.SetDefault("external.timeout", "30s")
	viper.SetDefault("external.retry_attempts", 3)
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

// GetRedisAddr returns the Redis connection address
func (c *Config) GetRedisAddr() string {
	return fmt.Sprintf("%s:%d", c.Redis.Host, c.Redis.Port)
}

// IsDevelopment returns true if the environment is development
func (c *Config) IsDevelopment() bool {
	return c.Server.Environment == "development"
}

// IsProduction returns true if the environment is production
func (c *Config) IsProduction() bool {
	return c.Server.Environment == "production"
}

// GetServerAddr returns the server address
func (c *Config) GetServerAddr() string {
	return fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port)
}
