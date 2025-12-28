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
		authHandler := handlers.NewAuthHandler(authService, db, cfg)

		auth := api.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.RefreshToken)
			auth.POST("/forgot-password", authHandler.ForgotPassword)
			auth.POST("/reset-password", authHandler.ResetPassword)
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
		passwordHandler := handlers.NewPasswordHandler(db, cfg)

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
			users.POST("/:id/reset-password", middleware.RequireRole("ADMIN"), passwordHandler.ForceResetPassword)
			users.POST("/import", middleware.RequireRole("ADMIN"), userHandler.ImportUsersFromCSV)
			users.GET("/export", userHandler.ExportUsersToCSV)
		}

		// File service untuk photo uploads
		fileService := services.NewFileService(db)

		// Achievement service untuk gamification
		notificationService := services.NewNotificationService(db)
		achievementService := services.NewAchievementService(db, notificationService)
		achievementHandler := handlers.NewAchievementHandler(achievementService)

		// Profile routes (Self-service untuk semua authenticated users)
		profileHandler := handlers.NewProfileHandler(userService, fileService, achievementService)

		profile := api.Group("/profile")
		profile.Use(middleware.AuthMiddleware(db, cfg))
		profile.Use(middleware.ActivityLogger(db))
		{
			profile.GET("", profileHandler.GetProfile)
			profile.PUT("", profileHandler.UpdateProfile)
			profile.PUT("/password", passwordHandler.ChangePassword)
			profile.POST("/photo", profileHandler.UploadProfilePhoto)
			profile.DELETE("/photo", profileHandler.DeleteProfilePhoto)
			profile.GET("/achievements", achievementHandler.GetUserAchievements)
			profile.GET("/stats", achievementHandler.GetUserStats)
		}

		// Achievement routes (Protected - All authenticated users)
		achievements := api.Group("/achievements")
		achievements.Use(middleware.AuthMiddleware(db, cfg))
		{
			achievements.GET("", achievementHandler.GetAllAchievements)
		}

		// Admin Achievement routes
		adminAchievements := api.Group("/admin/achievements")
		adminAchievements.Use(middleware.AuthMiddleware(db, cfg))
		adminAchievements.Use(middleware.RequireRole("ADMIN"))
		{
			adminAchievements.POST("/award", achievementHandler.AwardAchievement)
		}

		// Admin User Achievement routes
		adminUsers := api.Group("/admin/users")
		adminUsers.Use(middleware.AuthMiddleware(db, cfg))
		adminUsers.Use(middleware.RequireRole("ADMIN", "MANAGER"))
		{
			adminUsers.GET("/:id/achievements", achievementHandler.GetAchievementsByUserID)
		}

		// Notification routes (Protected - All authenticated users)
		notificationHandler := handlers.NewNotificationHandler(notificationService)

		notifications := api.Group("/notifications")
		notifications.Use(middleware.AuthMiddleware(db, cfg))
		{
			notifications.GET("", notificationHandler.GetUserNotifications)
			notifications.GET("/unread-count", notificationHandler.GetUnreadCount)
			notifications.GET("/recent", notificationHandler.GetRecentNotifications)
			notifications.PUT("/:id/read", notificationHandler.MarkAsRead)
			notifications.PUT("/read-all", notificationHandler.MarkAllAsRead)
			notifications.DELETE("/:id", notificationHandler.DeleteNotification)
		}

		// Activity Log routes (Admin only)
		activityLogService := services.NewActivityLogService(db)
		activityLogHandler := handlers.NewActivityLogHandler(activityLogService)

		activityLogs := api.Group("/admin/activity-logs")
		activityLogs.Use(middleware.AuthMiddleware(db, cfg))
		activityLogs.Use(middleware.RequireRole("ADMIN", "MANAGER"))
		{
			activityLogs.GET("", activityLogHandler.GetActivityLogs)
			activityLogs.GET("/stats", activityLogHandler.GetActivityStats)
			activityLogs.GET("/:id", activityLogHandler.GetActivityLogByID)
			activityLogs.GET("/user/:id", activityLogHandler.GetUserActivity)
		}

		// Profile activity logs (Self-service)
		profileActivity := api.Group("/profile")
		profileActivity.Use(middleware.AuthMiddleware(db, cfg))
		{
			profileActivity.GET("/activity", activityLogHandler.GetMyActivity)
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

	// Serve static files untuk uploads
	r.Static("/uploads", "./public/uploads")

	// Serve static files for frontend (hanya untuk production)
	// Dalam development mode, frontend diakses melalui Vite dev server (port 5173)
	// r.Static("/assets", "./frontend/dist/assets")
	// r.NoRoute(func(c *gin.Context) {
	// 	c.File("./frontend/dist/index.html")
	// })
}
