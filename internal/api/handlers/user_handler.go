package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"github.com/moha/kaafipay-backend/internal/repository"
	"github.com/moha/kaafipay-backend/internal/utils"
)

type UserHandler struct {
	userRepo repository.UserRepository
}

func NewUserHandler(userRepo repository.UserRepository) *UserHandler {
	return &UserHandler{
		userRepo: userRepo,
	}
}

// Profile response struct
type ProfileResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Phone string `json:"phone"`
}

// Update profile request struct
type UpdateProfileRequest struct {
	Name string `json:"name" binding:"required,min=1,max=100"`
}

// Change password request struct
type ChangePasswordRequest struct {
	CurrentPassword string `json:"currentPassword" binding:"required"`
	NewPassword     string `json:"newPassword" binding:"required,min=8"`
}

// GetProfile handles GET /user/profile
func (h *UserHandler) GetProfile(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		return
	}

	// Get user from database
	user, err := h.userRepo.FindByID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "DATABASE_ERROR",
				"message": "Error fetching user profile",
			},
		})
		return
	}

	// Create response
	response := ProfileResponse{
		ID:    user.ID.String(),
		Name:  user.Name,
		Phone: user.Phone,
	}

	c.JSON(http.StatusOK, response)
}

// UpdateProfile handles PUT /user/profile
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	// Get user ID from context
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		return
	}

	// Parse request body
	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "VALIDATION_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	// Trim whitespace
	req.Name = strings.TrimSpace(req.Name)

	// Get user from database
	user, err := h.userRepo.FindByID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "DATABASE_ERROR",
				"message": "Error fetching user profile",
			},
		})
		return
	}

	// Update user
	user.Name = req.Name
	if err := h.userRepo.Update(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "DATABASE_ERROR",
				"message": "Error updating user profile",
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Profile updated successfully",
	})
}

// ChangePassword handles PUT /user/password
func (h *UserHandler) ChangePassword(c *gin.Context) {
	// Get user ID from context
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		return
	}

	// Parse request body
	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "VALIDATION_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	// Get user from database
	user, err := h.userRepo.FindByID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "DATABASE_ERROR",
				"message": "Error fetching user profile",
			},
		})
		return
	}

	// Verify current password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.CurrentPassword)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": gin.H{
				"code":    "INVALID_CREDENTIALS",
				"message": "Current password is incorrect",
			},
		})
		return
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Error processing password change",
			},
		})
		return
	}

	// Update user password
	user.Password = string(hashedPassword)
	if err := h.userRepo.Update(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "DATABASE_ERROR",
				"message": "Error updating password",
			},
		})
		return
	}

	// TODO: Implement session invalidation for other devices
	// TODO: Send email notification about password change

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Password updated successfully",
	})
}
