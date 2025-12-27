package main

import (
	"log"
	"sirine-go/backend/config"
	"sirine-go/backend/database"
	"sirine-go/backend/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables dari file .env
	// Path relatif dari backend/ directory (working directory saat go run)
	if err := godotenv.Load("../.env"); err != nil {
		log.Println("Warning: .env file tidak ditemukan, menggunakan default values")
	}

	// Load configuration
	cfg := config.LoadConfig()

	// Set Gin mode
	gin.SetMode(cfg.GinMode)

	// Connect to database
	if err := database.Connect(cfg); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto migrate models menggunakan registry
	registry := database.NewModelsRegistry()
	if err := database.AutoMigrate(registry.GetModels()...); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
	log.Printf("Database tables synchronized (%d models)", registry.GetTableCount())

	// Initialize Gin router
	r := gin.Default()

	// Setup routes
	routes.SetupRoutes(r, cfg)

	// Start server
	log.Printf("Server berjalan di port %s", cfg.ServerPort)
	if err := r.Run(":" + cfg.ServerPort); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
