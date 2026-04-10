package main

import (
	"backend/internal/shared/config"
	"backend/internal/shared/database"
	"backend/internal/users/delivery"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
    cfg := config.Load()
    if err := database.Connect(cfg); err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }
    r := gin.Default()
    r.SetTrustedProxies([]string{"127.0.0.1", "172.18.0.0/16"})
    delivery.SetupRoutes(r)
    log.Printf("Server starting on port %s", cfg.Port)
    if err := r.Run(":" + cfg.Port); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}