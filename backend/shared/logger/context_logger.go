package logger

import (
	"context"

	"github.com/rs/zerolog/log"
)

// ContextLogger provides simple contextual logging
type ContextLogger struct {
	service string
	handler string
}

// NewContextLogger creates a new context logger
func NewContextLogger(service, handler, method string) *ContextLogger {
	return &ContextLogger{
		service: service,
		handler: handler,
	}
}

// WithContext returns the same logger (simplified)
func (cl *ContextLogger) WithContext(ctx context.Context) *ContextLogger {
	return cl
}

// Error logs an error simply
func (cl *ContextLogger) Error(err error, msg string) {
	log.Error().Str("service", cl.service).Err(err).Msg(msg)
}
