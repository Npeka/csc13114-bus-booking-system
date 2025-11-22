package router

import (
	"bus-booking/payment-service/config"
	"bus-booking/payment-service/internal/handler"
	"bus-booking/shared/ginext"
	"bus-booking/shared/health"
	"bus-booking/shared/middleware"

	"github.com/gin-gonic/gin"
)

type RouterConfig struct {
	Config             *config.Config
	TransactionHandler handler.TransactionHandler
}

func SetupRoutes(router *gin.Engine, cfg *RouterConfig) {
	router.Use(middleware.Logger())
	router.Use(middleware.SetupCORS(&cfg.Config.CORS))
	router.Use(middleware.RequestContextMiddleware(cfg.Config.ServiceName))
	router.GET(health.Path, health.Handler(cfg.Config.ServiceName))

	v1 := router.Group("/api/v1")
	{
		payment := v1.Group("/transactions")
		{
			payment.POST("", ginext.WrapHandler(cfg.TransactionHandler.CreateTransaction))
		}
	}
}
