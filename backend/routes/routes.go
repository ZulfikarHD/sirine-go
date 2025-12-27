package routes

import (
	"sirine-go/backend/handlers"
	"sirine-go/backend/middleware"
	"sirine-go/backend/services"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	// Apply CORS middleware
	r.Use(middleware.CORS())

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"message": "Server berjalan dengan baik",
		})
	})

	// API routes
	api := r.Group("/api")
	{
		// Example routes
		exampleService := services.NewExampleService()
		exampleHandler := handlers.NewExampleHandler(exampleService)

		examples := api.Group("/examples")
		{
			examples.GET("", exampleHandler.GetAll)
			examples.GET("/:id", exampleHandler.GetByID)
			examples.POST("", exampleHandler.Create)
			examples.PUT("/:id", exampleHandler.Update)
			examples.DELETE("/:id", exampleHandler.Delete)
		}
	}

	// Serve static files for frontend (untuk production)
	r.Static("/assets", "./frontend/dist/assets")
	r.StaticFile("/", "./frontend/dist/index.html")
	r.NoRoute(func(c *gin.Context) {
		c.File("./frontend/dist/index.html")
	})
}
