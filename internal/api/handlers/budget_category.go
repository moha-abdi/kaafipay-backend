package handlers

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/moha/kaafipay-backend/internal/models"
)

var (
	ErrUnauthorized = errors.New("unauthorized")
)

// BudgetCategoryHandler handles operations on budget categories
type BudgetCategoryHandler struct {
	db *gorm.DB
}

// NewBudgetCategoryHandler creates a new BudgetCategoryHandler instance
func NewBudgetCategoryHandler(db *gorm.DB) *BudgetCategoryHandler {
	return &BudgetCategoryHandler{db: db}
}

// GetBudgetCategories returns all budget categories for a user
func (h *BudgetCategoryHandler) GetBudgetCategories(c *gin.Context) {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		return
	}

	var categories []models.BudgetCategory
	if err := h.db.Where("user_id = ?", userID).Find(&categories).Error; err != nil {
		log.Printf("[GET-BUDGET-CATEGORIES] Database query failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{
			"code":    "INTERNAL_ERROR",
			"message": "Failed to fetch budget categories",
		}})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": categories})
}

// CreateBudgetCategory creates a new budget category
func (h *BudgetCategoryHandler) CreateBudgetCategory(c *gin.Context) {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		log.Printf("[CREATE-BUDGET-CATEGORY] Failed to get user ID from context: %v", err)
		return
	}
	log.Printf("[CREATE-BUDGET-CATEGORY] Processing request for user: %s", userID)

	var category models.BudgetCategory
	if err := c.ShouldBindJSON(&category); err != nil {
		log.Printf("[CREATE-BUDGET-CATEGORY] JSON binding failed: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{
			"code":    "VALIDATION_ERROR",
			"message": "Invalid request parameters",
		}})
		return
	}
	log.Printf("[CREATE-BUDGET-CATEGORY] Request payload: name=%s, icon=%s, budget=%.2f, rules_count=%d",
		category.Name, category.Icon, category.Budget, len(category.Rules))

	category.UserID = userID

	if err := category.Validate(); err != nil {
		if validationErr, ok := err.(models.ValidationError); ok {
			log.Printf("[CREATE-BUDGET-CATEGORY] Validation failed: field=%s, message=%s",
				validationErr.Field, validationErr.Message)
			c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{
				"code":    "VALIDATION_ERROR",
				"message": "Invalid request parameters",
				"details": map[string][]string{
					validationErr.Field: {validationErr.Message},
				},
			}})
			return
		}
		log.Printf("[CREATE-BUDGET-CATEGORY] Unexpected validation error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{
			"code":    "INTERNAL_ERROR",
			"message": "An unexpected error occurred",
		}})
		return
	}
	log.Printf("[CREATE-BUDGET-CATEGORY] Validation passed for category: %s", category.Name)

	if err := h.db.Create(&category).Error; err != nil {
		log.Printf("[CREATE-BUDGET-CATEGORY] Database creation failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{
			"code":    "INTERNAL_ERROR",
			"message": "Failed to create budget category",
		}})
		return
	}
	log.Printf("[CREATE-BUDGET-CATEGORY] Successfully created category: id=%s, name=%s",
		category.ID, category.Name)

	c.JSON(http.StatusCreated, gin.H{"data": category})
}

// UpdateBudgetCategory updates an existing budget category
func (h *BudgetCategoryHandler) UpdateBudgetCategory(c *gin.Context) {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		return
	}

	categoryID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{
			"code":    "INVALID_REQUEST",
			"message": "Invalid category ID",
		}})
		return
	}

	var existingCategory models.BudgetCategory
	if err := h.db.Where("id = ? AND user_id = ?", categoryID, userID).First(&existingCategory).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": gin.H{
				"code":    "NOT_FOUND",
				"message": "Budget category not found",
			}})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{
			"code":    "INTERNAL_ERROR",
			"message": "Failed to fetch budget category",
		}})
		return
	}

	var updateData models.BudgetCategory
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{
			"code":    "VALIDATION_ERROR",
			"message": "Invalid request parameters",
		}})
		return
	}

	// Update only provided fields
	if updateData.Name != "" {
		existingCategory.Name = updateData.Name
	}
	if updateData.Icon != "" {
		existingCategory.Icon = updateData.Icon
	}
	if updateData.Budget > 0 {
		existingCategory.Budget = updateData.Budget
	}
	if len(updateData.Rules) > 0 {
		existingCategory.Rules = updateData.Rules
	}

	if err := existingCategory.Validate(); err != nil {
		if validationErr, ok := err.(models.ValidationError); ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{
				"code":    "VALIDATION_ERROR",
				"message": "Invalid request parameters",
				"details": map[string][]string{
					validationErr.Field: {validationErr.Message},
				},
			}})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{
			"code":    "INTERNAL_ERROR",
			"message": "An unexpected error occurred",
		}})
		return
	}

	if err := h.db.Save(&existingCategory).Error; err != nil {
		log.Printf("[UPDATE-BUDGET-CATEGORY] Database error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{
			"code":    "INTERNAL_ERROR",
			"message": "Failed to update budget category",
		}})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": existingCategory})
}

// DeleteBudgetCategory deletes a budget category
func (h *BudgetCategoryHandler) DeleteBudgetCategory(c *gin.Context) {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		return
	}

	categoryID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{
			"code":    "INVALID_REQUEST",
			"message": "Invalid category ID",
		}})
		return
	}

	result := h.db.Where("id = ? AND user_id = ?", categoryID, userID).Delete(&models.BudgetCategory{})
	if result.Error != nil {
		log.Printf("[DELETE-BUDGET-CATEGORY] Database error: %v", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{
			"code":    "INTERNAL_ERROR",
			"message": "Failed to delete budget category",
		}})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": gin.H{
			"code":    "NOT_FOUND",
			"message": "Budget category not found",
		}})
		return
	}

	c.Status(http.StatusNoContent)
}

// Helper function to get user ID from context
func getUserIDFromContext(c *gin.Context) (uuid.UUID, error) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": gin.H{
			"code":    "UNAUTHORIZED",
			"message": "User ID not found in context",
		}})
		return uuid.Nil, ErrUnauthorized
	}

	// Try to get UUID directly
	userID, ok := userIDInterface.(uuid.UUID)
	if !ok {
		// If not UUID, try string conversion as fallback
		if userIDStr, isStr := userIDInterface.(string); isStr {
			var err error
			userID, err = uuid.Parse(userIDStr)
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": gin.H{
					"code":    "UNAUTHORIZED",
					"message": "Invalid user ID format",
				}})
				return uuid.Nil, ErrUnauthorized
			}
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": gin.H{
				"code":    "UNAUTHORIZED",
				"message": "Invalid user ID format",
			}})
			return uuid.Nil, ErrUnauthorized
		}
	}

	return userID, nil
}
