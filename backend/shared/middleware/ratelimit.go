package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"bus-booking/shared/config"
	"bus-booking/shared/response"

	"github.com/gin-gonic/gin"
	goredis "github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"github.com/ulule/limiter/v3"
	limiterstore "github.com/ulule/limiter/v3/drivers/store/redis"
)

// RateLimiter creates a rate limiting middleware
func RateLimiter(redisClient *goredis.Client, cfg *config.RateLimitConfig) gin.HandlerFunc {
	// Create rate limiter with Redis store
	store, err := limiterstore.NewStore(redisClient)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create Redis store for rate limiter")
	}

	// Create rate limiter instance
	rate := limiter.Rate{
		Period: cfg.Period,
		Limit:  int64(cfg.RPS),
	}

	instance := limiter.New(store, rate, limiter.WithTrustForwardHeader(true))

	return func(c *gin.Context) {
		// Get client IP for rate limiting
		key := getClientIP(c)

		// Get rate limit context
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Check rate limit
		limiterCtx, err := instance.Get(ctx, key)
		if err != nil {
			log.Error().Err(err).Msg("Failed to get rate limit context")
			response.InternalServerErrorResponse(c, "Rate limiter error")
			c.Abort()
			return
		}

		// Set rate limit headers
		c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", limiterCtx.Limit))
		c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", limiterCtx.Remaining))
		c.Header("X-RateLimit-Reset", fmt.Sprintf("%d", limiterCtx.Reset))

		// Check if rate limit exceeded
		if limiterCtx.Reached {
			log.Warn().
				Str("client_ip", key).
				Int64("limit", limiterCtx.Limit).
				Int64("remaining", limiterCtx.Remaining).
				Int64("reset_time", limiterCtx.Reset).
				Msg("Rate limit exceeded")

			response.TooManyRequestsResponse(c, "Rate limit exceeded. Please try again later.")
			c.Abort()
			return
		}

		c.Next()
	}
}

// getClientIP extracts client IP from request
func getClientIP(c *gin.Context) string {
	// Check X-Forwarded-For header first
	if xff := c.GetHeader("X-Forwarded-For"); xff != "" {
		return xff
	}

	// Check X-Real-IP header
	if xri := c.GetHeader("X-Real-IP"); xri != "" {
		return xri
	}

	// Fall back to remote address
	return c.ClientIP()
}

// IPWhitelistMiddleware creates an IP whitelist middleware
func IPWhitelistMiddleware(allowedIPs []string) gin.HandlerFunc {
	allowedIPMap := make(map[string]bool)
	for _, ip := range allowedIPs {
		allowedIPMap[ip] = true
	}

	return func(c *gin.Context) {
		clientIP := getClientIP(c)

		if !allowedIPMap[clientIP] {
			log.Warn().Str("client_ip", clientIP).Msg("IP not in whitelist")
			response.ForbiddenResponse(c, "Access denied from this IP address")
			c.Abort()
			return
		}

		c.Next()
	}
}

// TimeoutMiddleware creates a request timeout middleware
func TimeoutMiddleware(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()

		c.Request = c.Request.WithContext(ctx)

		done := make(chan struct{})
		go func() {
			defer close(done)
			c.Next()
		}()

		select {
		case <-done:
			return
		case <-ctx.Done():
			response.ErrorResponse(c, http.StatusRequestTimeout, "REQUEST_TIMEOUT", "Request timeout")
			c.Abort()
		}
	}
}
