package router

import (
	"event-scope/backend/internal/handler"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Setup() *gin.Engine {
	router := gin.Default()

	router.GET("/healthz", func(context *gin.Context) {
		handler.Res(context, http.StatusOK, "API is healthy", nil)
	})

	eventRoutes(router.Group("/event"))

	return router
}
