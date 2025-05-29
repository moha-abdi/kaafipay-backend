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

		// Linked accounts routes
		accounts := v1.Group("/linked-accounts")
		{
			// Protected routes
			protected := accounts.Use(middleware.AuthMiddleware(cfg))
			{
				protected.POST("", linkedAccountHandler.LinkAccount)
				protected.GET("", linkedAccountHandler.GetLinkedAccounts)
				protected.GET("/:id", linkedAccountHandler.GetLinkedAccount)
				protected.DELETE("/:id", linkedAccountHandler.UnlinkAccount)
				protected.PATCH("/:id/default", linkedAccountHandler.SetDefaultAccount)
				protected.POST("/:id/refresh", linkedAccountHandler.RefreshAccount)
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
