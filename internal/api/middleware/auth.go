package middleware

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/moha/kaafipay-backend/internal/config"
	"github.com/moha/kaafipay-backend/internal/utils"
)

func init() {
	// Configure logging to include timestamps and file information
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Log the request path for context
		log.Printf("[AUTH] New request to: %s", c.Request.URL.Path)

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			log.Printf("[AUTH] Missing authorization header for request to: %s", c.Request.URL.Path)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 || bearerToken[0] != "Bearer" {
			log.Printf("[AUTH] Invalid token format. Got: %s", authHeader)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			c.Abort()
			return
		}

		// Log the token for debugging (be careful with this in production!)
		log.Printf("[AUTH] Processing token: %s", bearerToken[1])

		claims, err := utils.ValidateToken(bearerToken[1], cfg.JWTSecret)
		if err != nil {
			log.Printf("[AUTH] Token validation failed: %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		log.Printf("[AUTH] Token validated successfully. UserID: %s, Phone: %s", claims.UserID, claims.Phone)

		// Set user information in context
		c.Set("user_id", claims.UserID)
		c.Set("phone", claims.Phone)

		// Verify the values were set correctly
		userID, exists := c.Get("user_id")
		if !exists {
			log.Printf("[AUTH] Failed to set user_id in context!")
		} else {
			log.Printf("[AUTH] Successfully set user_id in context: %v", userID)
		}

		c.Next()
	}
}
