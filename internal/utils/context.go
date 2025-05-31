package utils

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GetUserIDFromContext extracts and validates the user ID from the Gin context
func GetUserIDFromContext(c *gin.Context) (uuid.UUID, error) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": gin.H{
				"code":    "UNAUTHORIZED",
				"message": "User ID not found in context",
			},
		})
		return uuid.Nil, errors.New("user ID not found in context")
	}

	// Try to get UUID directly
	userID, ok := userIDInterface.(uuid.UUID)
	if !ok {
		// If not UUID, try string conversion as fallback
		if userIDStr, isStr := userIDInterface.(string); isStr {
			var err error
			userID, err = uuid.Parse(userIDStr)
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{
					"error": gin.H{
						"code":    "UNAUTHORIZED",
						"message": "Invalid user ID format",
					},
				})
				return uuid.Nil, err
			}
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": gin.H{
					"code":    "UNAUTHORIZED",
					"message": "Invalid user ID format",
				},
			})
			return uuid.Nil, errors.New("invalid user ID format")
		}
	}

	return userID, nil
}
