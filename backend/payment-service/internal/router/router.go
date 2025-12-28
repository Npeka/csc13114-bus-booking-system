package router

import (
	"bus-booking/payment-service/config"
	"bus-booking/payment-service/internal/handler"
	"bus-booking/shared/constants"
	"bus-booking/shared/ginext"
	"bus-booking/shared/health"
	"bus-booking/shared/middleware"
	"bus-booking/shared/swagger"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Handlers struct {
	TransactionHandler handler.TransactionHandler
	BankAccountHandler handler.BankAccountHandler
	ConstantsHandler   handler.ConstantsHandler
	RefundHandler      handler.RefundHandler
}

func SetupRoutes(router *gin.Engine, cfg *config.Config, h *Handlers) {
	router.Use(middleware.Logger())
	router.Use(middleware.SetupCORS(&cfg.CORS))
	router.Use(middleware.RequestContext(cfg.ServiceName))
	router.GET(health.Path, health.Handler(cfg.ServiceName))
	router.GET(swagger.Path, ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := router.Group("/api/v1")
	{
		constants := v1.Group("/constants")
		{
			constants.GET("", ginext.WrapHandler(h.ConstantsHandler.GetConstants))
		}

		transactions := v1.Group("/transactions")
		{
			transactions.POST("/webhook", ginext.WrapHandler(h.TransactionHandler.HandleWebhook))
		}
	}

	userV1 := router.Group("/api/v1")
	userV1.Use(middleware.RequireAuth())
	{
		bankAccounts := userV1.Group("/bank-accounts")
		{
			bankAccounts.POST("", ginext.WrapHandler(h.BankAccountHandler.CreateBankAccount))
			bankAccounts.GET("", ginext.WrapHandler(h.BankAccountHandler.GetMyBankAccounts))
			bankAccounts.PUT("/:id", ginext.WrapHandler(h.BankAccountHandler.UpdateBankAccount))
			bankAccounts.DELETE("/:id", ginext.WrapHandler(h.BankAccountHandler.DeleteBankAccount))
			bankAccounts.POST("/:id/set-primary", ginext.WrapHandler(h.BankAccountHandler.SetPrimaryBankAccount))
		}

		// Refund routes (for users)
		refunds := userV1.Group("/refunds")
		{
			refunds.POST("", ginext.WrapHandler(h.RefundHandler.Create))
			refunds.GET("/booking/:booking_id", ginext.WrapHandler(h.RefundHandler.GetByBookingID))
		}
	}

	adminV1 := router.Group("/api/v1")
	adminV1.Use(middleware.RequireAuth())
	adminV1.Use(middleware.RequireRole(constants.RoleAdmin))
	{
		transactions := adminV1.Group("/transactions")
		{
			transactions.GET("", ginext.WrapHandler(h.TransactionHandler.GetList))
			transactions.GET("/stats", ginext.WrapHandler(h.TransactionHandler.GetStats))
		}

		refunds := adminV1.Group("/refunds")
		{
			refunds.GET("", ginext.WrapHandler(h.RefundHandler.ListRefunds))
			refunds.PUT("/:id", ginext.WrapHandler(h.RefundHandler.UpdateRefundStatus))
			refunds.POST("/export", func(c *gin.Context) {
				req := &ginext.Request{GinCtx: c}
				if err := h.RefundHandler.ExportRefunds(req); err != nil {
					log.Error().Err(err).Msg("failed to export refunds")
				}
			})
		}
	}

	internalV1 := router.Group("/api/v1")
	{
		transactions := internalV1.Group("/transactions")
		{
			transactions.POST("", ginext.WrapHandler(h.TransactionHandler.Create))
			transactions.GET("/:id", ginext.WrapHandler(h.TransactionHandler.GetByID))
			transactions.POST("/:id/cancel", ginext.WrapHandler(h.TransactionHandler.Cancel))
		}
	}
}
