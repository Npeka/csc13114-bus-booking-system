package middleware

import (
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// RequestID adds a unique request ID to each request
func RequestID() gin.HandlerFunc {
	return requestid.New()
}

// Logger creates a structured logging middleware
func Logger() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		log.Info().
			Str("client_ip", param.ClientIP).
			Time("timestamp", param.TimeStamp).
			Str("method", param.Method).
			Str("path", param.Path).
			Str("protocol", param.Request.Proto).
			Int("status_code", param.StatusCode).
			Dur("latency", param.Latency).
			Str("user_agent", param.Request.UserAgent()).
			Str("error", param.ErrorMessage).
			Int("body_size", param.BodySize).
			Interface("request_id", param.Keys["request_id"]).
			Msg("HTTP Request")

		return ""
	})
}

// Recovery handles panics and returns a 500 error
func Recovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		log.Error().
			Interface("error", recovered).
			Str("request_id", c.GetString("request_id")).
			Str("path", c.Request.URL.Path).
			Str("method", c.Request.Method).
			Msg("Panic recovered")

		c.JSON(500, gin.H{
			"error": map[string]interface{}{
				"code":    "INTERNAL_SERVER_ERROR",
				"message": "Internal server error",
			},
		})
		c.Abort()
	})
}
