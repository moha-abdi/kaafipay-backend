package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/moha/kaafipay-backend/internal/api/handlers"
	"github.com/moha/kaafipay-backend/internal/api/middleware"
	"github.com/moha/kaafipay-backend/internal/config"
	"github.com/moha/kaafipay-backend/internal/repository"
	"github.com/moha/kaafipay-backend/internal/services/whatsapp"
)

func SetupRouter(cfg *config.Config, db *gorm.DB) *gin.Engine {
	router := gin.Default()

	// Middleware
	router.Use(middleware.CORS())

	// Initialize providers
	whatsappProvider := whatsapp.NewWhatsAppProvider(
		db,
		cfg.WhatsAppAPIBaseURL,
		cfg.WhatsAppAPIKey,
		cfg.WhatsAppSessionID,
	)

	// Repositories
	userRepo := repository.NewUserRepository(db)

	// Handlers
	authHandler := handlers.NewAuthHandler(cfg, userRepo)
	verifyHandler := handlers.NewVerifyHandler(whatsappProvider)
	linkedAccountHandler := handlers.NewLinkedAccountHandler(db)
	budgetHandler := handlers.NewBudgetCategoryHandler(db)
	userHandler := handlers.NewUserHandler(userRepo)

	// Public routes
	v1 := router.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		verify := v1.Group("/verify")
		{
			verify.POST("/send-code", verifyHandler.SendCode)
			verify.POST("/verify-code", verifyHandler.VerifyCode)
			verify.POST("/verify-token", verifyHandler.VerifyToken)
		}

		// Protected routes
		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware(cfg))
		{
			// User profile routes
			user := protected.Group("/user")
			{
				user.GET("/profile", userHandler.GetProfile)
				user.PUT("/profile", userHandler.UpdateProfile)
				user.PUT("/password", userHandler.ChangePassword)
			}

			// Linked accounts routes
			accounts := protected.Group("/linked-accounts")
			{
				accounts.POST("", linkedAccountHandler.LinkAccount)
				accounts.GET("", linkedAccountHandler.GetLinkedAccounts)
				accounts.GET("/:id", linkedAccountHandler.GetLinkedAccount)
				accounts.DELETE("/:id", linkedAccountHandler.UnlinkAccount)
				accounts.PATCH("/:id/default", linkedAccountHandler.SetDefaultAccount)
				accounts.POST("/:id/refresh", linkedAccountHandler.RefreshAccount)
			}

			// Budget categories routes
			budgets := protected.Group("/budget-categories")
			{
				budgets.GET("", budgetHandler.GetBudgetCategories)
				budgets.POST("", budgetHandler.CreateBudgetCategory)
				budgets.PUT("/:id", budgetHandler.UpdateBudgetCategory)
				budgets.DELETE("/:id", budgetHandler.DeleteBudgetCategory)
			}
		}

		admin := v1.Group("/admin")
		admin.Use(middleware.AdminAuthMiddleware(cfg.AdminToken))
		{
			adminHandler := handlers.NewAdminHandler(whatsappProvider)
			whatsapp := admin.Group("/whatsapp")
			{
				whatsapp.GET("/sessions", adminHandler.ListSessions)
				whatsapp.GET("/sessions/:sessionId", adminHandler.GetSession)
				whatsapp.POST("/sessions", adminHandler.AddSession)
				whatsapp.DELETE("/sessions/:sessionId", adminHandler.DeleteSession)
			}
		}
	}

	return router
}
