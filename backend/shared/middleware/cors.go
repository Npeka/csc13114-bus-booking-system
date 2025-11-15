package middleware

import (
	"time"

	"csc13114-bus-ticket-booking-system/shared/config"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// SetupCORS configures CORS middleware
func SetupCORS(cfg *config.CORSConfig) gin.HandlerFunc {
	corsConfig := cors.Config{
		AllowOrigins:     cfg.AllowOrigins,
		AllowMethods:     cfg.AllowMethods,
		AllowHeaders:     cfg.AllowHeaders,
		ExposeHeaders:    cfg.ExposeHeaders,
		AllowCredentials: cfg.AllowCredentials,
		MaxAge:           time.Duration(cfg.MaxAge) * time.Second,
	}

	// If allow origins contains "*", use AllowAllOrigins
	for _, origin := range cfg.AllowOrigins {
		if origin == "*" {
			corsConfig.AllowAllOrigins = true
			corsConfig.AllowOrigins = nil
			break
		}
	}

	return cors.New(corsConfig)
}
