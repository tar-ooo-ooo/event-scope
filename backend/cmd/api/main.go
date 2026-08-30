package main

import (
	"log"

	"event-scope/backend/internal/config"
	"event-scope/backend/internal/kafka"
	"event-scope/backend/internal/router"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	cfg := config.Load()

	wr := kafka.NewWriter(cfg.Kafka.Broker, cfg.Kafka.Topics[0].Name)
	defer wr.Close()

	engine := router.Setup(cfg, wr)

	log.Printf("🚀 API 啟動於 http://localhost:%s", cfg.Port)
	if err := engine.Run(":" + cfg.Port); err != nil {
		log.Fatalf("🫠 無法啟動 API 伺服器：%v", err)
	}
}
