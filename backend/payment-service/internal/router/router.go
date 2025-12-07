package router

import (
	"bus-booking/payment-service/config"
	"bus-booking/payment-service/internal/handler"
	"bus-booking/shared/ginext"
	"bus-booking/shared/health"
	"bus-booking/shared/middleware"
	"bus-booking/shared/swagger"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Handlers struct {
	TransactionHandler handler.TransactionHandler
}

func SetupRoutes(router *gin.Engine, cfg *config.Config, h *Handlers) {
	router.Use(middleware.Logger())
	router.Use(middleware.SetupCORS(&cfg.CORS))
	router.Use(middleware.RequestContext(cfg.ServiceName))
	router.GET(health.Path, health.Handler(cfg.ServiceName))
	router.GET(swagger.Path, ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := router.Group("/api/v1")
	{
		transactions := v1.Group("/transactions")
		{
			transactions.POST("", ginext.WrapHandler(h.TransactionHandler.CreateTransaction))

			transactions.POST("/payment-link", ginext.WrapHandler(h.TransactionHandler.CreatePaymentLink))
			transactions.POST("/webhook", ginext.WrapHandler(h.TransactionHandler.HandlePaymentWebhook))
			transactions.GET("/return", ginext.WrapHandler(h.TransactionHandler.HandlePaymentReturn))
			transactions.GET("/cancel", ginext.WrapHandler(h.TransactionHandler.HandlePaymentCancel))
			transactions.GET("/:order_code", ginext.WrapHandler(h.TransactionHandler.GetTransactionByOrderCode))
		}
	}
}
