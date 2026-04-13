package main

import (
	authDelivery "backend/internal/auth/delivery"
    assignmentsDelivery "backend/internal/assignments/delivery"
	periodDelivery "backend/internal/periods/delivery"
	"backend/internal/shared/config"
	"backend/internal/shared/database"
	usersDelivery "backend/internal/users/delivery"
	usersSeed "backend/internal/users/seed"
	weeksDelivery "backend/internal/weeks/delivery"
	workspacesDelivery "backend/internal/workspaces/delivery"
	"log"
	
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()
	if err := database.Connect(cfg); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err := usersSeed.SeedDefaultAdmin(); err != nil {
		log.Fatalf("Failed to seed default admin: %v", err)
	}

	r := gin.Default()
	r.SetTrustedProxies([]string{"127.0.0.1", "172.18.0.0/16"})

	api := r.Group("/api")
	api.GET("", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"project": "Seneprojects",
			"message": "Welcome to Seneprojects API",
		})
	})

	authHandler := authDelivery.SetupRoutes(api, cfg)
	usersDelivery.SetupRoutes(api, authHandler)
	periodDelivery.SetupRoutes(api)
	weeksDelivery.SetupRoutes(api)
	workspacesDelivery.SetupRoutes(api)
    assignmentsDelivery.SetupRoutes(api)

	log.Printf("Server starting on port %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
