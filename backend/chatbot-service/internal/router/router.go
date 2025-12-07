package router

import (
	"bus-booking/chatbot-service/config"
	"bus-booking/chatbot-service/internal/handler"
	"bus-booking/shared/ginext"
	"bus-booking/shared/health"
	"bus-booking/shared/middleware"
	"bus-booking/shared/swagger"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Handlers struct {
	ChatHandler handler.ChatHandler
}

func SetupRoutes(router *gin.Engine, cfg *config.Config, h *Handlers) {
	router.Use(middleware.Logger())
	router.Use(middleware.SetupCORS(&cfg.CORS))
	router.Use(middleware.RequestContext(cfg.ServiceName))
	router.GET(health.Path, health.Handler(cfg.ServiceName))
	router.GET(swagger.Path, ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := router.Group("/api/v1")
	{
		chat := v1.Group("/chat")
		{
			chat.POST("", ginext.WrapHandler(h.ChatHandler.Chat))
			chat.GET("/extract-search", ginext.WrapHandler(h.ChatHandler.ExtractSearchParams))
			chat.GET("/faq", ginext.WrapHandler(h.ChatHandler.GetFAQ))
		}
	}
}
