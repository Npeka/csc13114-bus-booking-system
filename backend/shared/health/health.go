package health

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const Path = "/health"

func Handler(serviceName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"service": serviceName,
		})
	}
}
