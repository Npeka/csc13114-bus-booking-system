package logger

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	"bus-booking/shared/config"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/natefinch/lumberjack.v2"
)

// SetupLogger configures and initializes the zerolog logger
func SetupLogger(cfg *config.LogConfig) error {
	// Set log level
	level, err := parseLogLevel(cfg.Level)
	if err != nil {
		return err
	}
	zerolog.SetGlobalLevel(level)

	// Set output
	var writers []io.Writer

	switch cfg.Output {
	case "stdout":
		writers = append(writers, os.Stdout)
	case "stderr":
		writers = append(writers, os.Stderr)
	case "file":
		// Create logs directory if it doesn't exist
		logDir := filepath.Dir(cfg.Filename)
		if err := os.MkdirAll(logDir, 0755); err != nil {
			return err
		}

		// Configure log rotation
		fileWriter := &lumberjack.Logger{
			Filename:   cfg.Filename,
			MaxSize:    cfg.MaxSize, // megabytes
			MaxBackups: cfg.MaxBackups,
			MaxAge:     cfg.MaxAge, // days
			Compress:   cfg.Compress,
		}
		writers = append(writers, fileWriter)

		// Also write to stdout in development
		if level == zerolog.DebugLevel {
			writers = append(writers, os.Stdout)
		}
	default:
		writers = append(writers, os.Stdout)
	}

	// Set format and output
	var writer io.Writer
	if len(writers) > 1 {
		writer = io.MultiWriter(writers...)
	} else {
		writer = writers[0]
	}

	// Set format
	if cfg.Format == "text" || cfg.Format == "console" {
		writer = zerolog.ConsoleWriter{
			Out:        writer,
			TimeFormat: "2006-01-02 15:04:05",
		}
	}

	// Configure global logger
	log.Logger = zerolog.New(writer).With().
		Timestamp().
		Str("service", "bus-booking-system").
		Str("version", "1.0.0").
		Logger()

	// Add caller information in debug mode
	if cfg.Level == "debug" || cfg.Level == "trace" {
		log.Logger = log.Logger.With().Caller().Logger()
	}

	log.Info().Msg("Logger initialized successfully")

	return nil
}

// parseLogLevel converts string level to zerolog level
func parseLogLevel(level string) (zerolog.Level, error) {
	switch strings.ToLower(level) {
	case "trace":
		return zerolog.TraceLevel, nil
	case "debug":
		return zerolog.DebugLevel, nil
	case "info":
		return zerolog.InfoLevel, nil
	case "warn", "warning":
		return zerolog.WarnLevel, nil
	case "error":
		return zerolog.ErrorLevel, nil
	case "fatal":
		return zerolog.FatalLevel, nil
	case "panic":
		return zerolog.PanicLevel, nil
	default:
		return zerolog.InfoLevel, nil
	}
}

// GetLogger returns the global logger
func GetLogger() *zerolog.Logger {
	return &log.Logger
}

// WithFields creates a logger with structured fields
func WithFields(fields map[string]interface{}) *zerolog.Event {
	event := log.Info()
	for key, value := range fields {
		switch v := value.(type) {
		case string:
			event = event.Str(key, v)
		case int:
			event = event.Int(key, v)
		case int64:
			event = event.Int64(key, v)
		case float64:
			event = event.Float64(key, v)
		case bool:
			event = event.Bool(key, v)
		case error:
			event = event.AnErr(key, v)
		default:
			event = event.Interface(key, v)
		}
	}
	return event
}

// WithField creates a logger with a single field
func WithField(key string, value interface{}) *zerolog.Event {
	return WithFields(map[string]interface{}{key: value})
}

// WithError creates a logger with error field
func WithError(err error) *zerolog.Event {
	return log.Error().Err(err)
}

// HTTPLogger creates a logger for HTTP requests
func HTTPLogger(method, path, userAgent, clientIP string) *zerolog.Event {
	return log.Info().
		Str("type", "http").
		Str("method", method).
		Str("path", path).
		Str("user_agent", userAgent).
		Str("client_ip", clientIP)
}

// DatabaseLogger creates a logger for database operations
func DatabaseLogger(operation, table string) *zerolog.Event {
	return log.Debug().
		Str("type", "database").
		Str("operation", operation).
		Str("table", table)
}

// RedisLogger creates a logger for Redis operations
func RedisLogger(operation, key string) *zerolog.Event {
	return log.Debug().
		Str("type", "redis").
		Str("operation", operation).
		Str("key", key)
}

// ServiceLogger creates a logger for service operations
func ServiceLogger(service, operation string) *zerolog.Event {
	return log.Info().
		Str("type", "service").
		Str("service", service).
		Str("operation", operation)
}

// AuthLogger creates a logger for authentication operations
func AuthLogger(operation, userID string) *zerolog.Event {
	return log.Info().
		Str("type", "auth").
		Str("operation", operation).
		Str("user_id", userID)
}

// PaymentLogger creates a logger for payment operations
func PaymentLogger(operation, transactionID string) *zerolog.Event {
	return log.Info().
		Str("type", "payment").
		Str("operation", operation).
		Str("transaction_id", transactionID)
}

// ExternalAPILogger creates a logger for external API calls
func ExternalAPILogger(service, endpoint string) *zerolog.Event {
	return log.Info().
		Str("type", "external_api").
		Str("service", service).
		Str("endpoint", endpoint)
}

// ErrorLogger creates a logger for errors
func ErrorLogger() *zerolog.Event {
	return log.Error()
}

// InfoLogger creates a logger for info messages
func InfoLogger() *zerolog.Event {
	return log.Info()
}

// DebugLogger creates a logger for debug messages
func DebugLogger() *zerolog.Event {
	return log.Debug()
}

// WarnLogger creates a logger for warning messages
func WarnLogger() *zerolog.Event {
	return log.Warn()
}
