package routes

import (
	"sirine-go/backend/config"
	"sirine-go/backend/database"
	"sirine-go/backend/handlers"
	"sirine-go/backend/middleware"
	"sirine-go/backend/services"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, cfg *config.Config) {
	// Apply CORS middleware
	r.Use(middleware.CORS())

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"message": "Server berjalan dengan baik",
		})
	})

	// Get database instance
	db := database.GetDB()

	// API routes
	api := r.Group("/api")
	{
		// Authentication routes (public)
		authService := services.NewAuthService(db, cfg)
		authHandler := handlers.NewAuthHandler(authService)

		auth := api.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.RefreshToken)
		}

		// Protected authentication routes
		authProtected := api.Group("/auth")
		authProtected.Use(middleware.AuthMiddleware(db, cfg))
		{
			authProtected.POST("/logout", authHandler.Logout)
			authProtected.GET("/me", authHandler.GetCurrentUser)
		}

		// Example routes (protected) - Commented out for Sprint 1
		// Will be re-enabled when needed
		// exampleService := services.NewExampleService()
		// exampleHandler := handlers.NewExampleHandler(exampleService)
		// examples := api.Group("/examples")
		// examples.Use(middleware.AuthMiddleware(db, cfg))
		// {
		// 	examples.GET("", exampleHandler.GetAll)
		// 	examples.GET("/:id", exampleHandler.GetByID)
		// 	examples.POST("", exampleHandler.Create)
		// 	examples.PUT("/:id", exampleHandler.Update)
		// 	examples.DELETE("/:id", exampleHandler.Delete)
		// }
	}

	// Serve static files for frontend (untuk production)
	r.Static("/assets", "./frontend/dist/assets")
	r.StaticFile("/", "./frontend/dist/index.html")
	r.NoRoute(func(c *gin.Context) {
		c.File("./frontend/dist/index.html")
	})
}
