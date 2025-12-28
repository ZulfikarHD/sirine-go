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

	// Root endpoint untuk development mode info
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Sirine Go API Server",
			"version": "1.0.0",
			"endpoints": gin.H{
				"health": "/health",
				"api":    "/api/*",
				"docs":   "Lihat dokumentasi untuk endpoint lengkap",
			},
			"development": gin.H{
				"frontend_dev_server": "http://localhost:5173",
				"backend_api":         "http://localhost:8080/api",
				"note":                "Dalam development mode, akses frontend melalui Vite dev server (port 5173)",
			},
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

		// User Management routes (Admin/Manager only)
		passwordService := services.NewPasswordService()
		userService := services.NewUserService(db, passwordService)
		userHandler := handlers.NewUserHandler(userService)

		users := api.Group("/users")
		users.Use(middleware.AuthMiddleware(db, cfg))
		users.Use(middleware.RequireRole("ADMIN", "MANAGER"))
		users.Use(middleware.ActivityLogger(db))
		{
			users.GET("", userHandler.GetAllUsers)
			users.GET("/search", userHandler.SearchUsers)
			users.GET("/:id", userHandler.GetUserByID)
			users.POST("", middleware.RequireRole("ADMIN"), userHandler.CreateUser)
			users.PUT("/:id", middleware.RequireRole("ADMIN"), userHandler.UpdateUser)
			users.DELETE("/:id", middleware.RequireRole("ADMIN"), userHandler.DeleteUser)
			users.POST("/bulk-delete", middleware.RequireRole("ADMIN"), userHandler.BulkDeleteUsers)
			users.POST("/bulk-update-status", middleware.RequireRole("ADMIN"), userHandler.BulkUpdateStatus)
		}

		// Profile routes (Self-service untuk semua authenticated users)
		profileHandler := handlers.NewProfileHandler(userService)

		profile := api.Group("/profile")
		profile.Use(middleware.AuthMiddleware(db, cfg))
		profile.Use(middleware.ActivityLogger(db))
		{
			profile.GET("", profileHandler.GetProfile)
			profile.PUT("", profileHandler.UpdateProfile)
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

	// Serve static files for frontend (hanya untuk production)
	// Dalam development mode, frontend diakses melalui Vite dev server (port 5173)
	// r.Static("/assets", "./frontend/dist/assets")
	// r.NoRoute(func(c *gin.Context) {
	// 	c.File("./frontend/dist/index.html")
	// })
}
