package main

import (
	"log"
	"os"

	"event-scope/backend/internal/router"
)

func main() {
	engine := router.Setup()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("API 啟動於 http://localhost:%s", port)
	if err := engine.Run(":" + port); err != nil {
		log.Fatalf("無法啟動 API 伺服器：%v", err)
	}
}
