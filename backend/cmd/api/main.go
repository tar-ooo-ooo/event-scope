package main

import (
	"log"

	"event-scope/backend/internal/config"
	"event-scope/backend/internal/router"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	cfg := config.Load()
	engine := router.Setup()

	log.Printf("API 啟動於 http://localhost:%s", cfg.Port)
	if err := engine.Run(":" + cfg.Port); err != nil {
		log.Fatalf("無法啟動 API 伺服器：%v", err)
	}
}
