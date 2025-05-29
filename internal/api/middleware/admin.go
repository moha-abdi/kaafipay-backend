package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AdminAuthMiddleware(adminToken string) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("X-Admin-Token")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Admin token required"})
			c.Abort()
			return
		}

		if token != adminToken {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid admin token"})
			c.Abort()
			return
		}

		c.Next()
	}
}
