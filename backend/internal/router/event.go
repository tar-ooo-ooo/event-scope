package router

import (
	"github.com/gin-gonic/gin"

	"event-scope/backend/internal/handler"
)

func eventRoutes(rg *gin.RouterGroup) {
	rg.POST("", handler.CreateEventHandler)
}
