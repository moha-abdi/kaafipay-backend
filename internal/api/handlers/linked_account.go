package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/moha/kaafipay-backend/internal/models"
)

// LinkedAccountHandler handles operations on linked accounts
type LinkedAccountHandler struct {
	db *gorm.DB
}

// NewLinkedAccountHandler creates a new LinkedAccountHandler instance
func NewLinkedAccountHandler(db *gorm.DB) *LinkedAccountHandler {
	return &LinkedAccountHandler{db: db}
}

// Request/Response types
type linkAccountRequest struct {
	Provider      models.Provider `json:"provider" binding:"required"`
	AccountID     string          `json:"accountId" binding:"required"`
	AccountNumber string          `json:"accountNumber" binding:"required"`
	AccountTitle  string          `json:"accountTitle" binding:"required"`
	AccountType   string          `json:"accountType" binding:"required"`
	Currency      struct {
		Code   string `json:"code" binding:"required"`
		Name   string `json:"name" binding:"required"`
		Symbol string `json:"symbol" binding:"required"`
	} `json:"currency" binding:"required"`
	IsDefaultAccount bool `json:"isDefaultAccount"`
	Credentials      struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	} `json:"credentials" binding:"required"`
	DeviceInfo struct {
		DeviceID     string `json:"deviceId" binding:"required"`
		DeviceModel  string `json:"deviceModel" binding:"required"`
		Manufacturer string `json:"manufacturer" binding:"required"`
		OSVersion    string `json:"osVersion" binding:"required"`
	} `json:"deviceInfo" binding:"required"`
}

type accountResponse struct {
	ID            uuid.UUID       `json:"id"`
	Provider      models.Provider `json:"provider"`
	AccountID     string          `json:"accountId"`
	AccountNumber string          `json:"accountNumber"`
	AccountTitle  string          `json:"accountTitle"`
	AccountType   string          `json:"accountType"`
	Currency      struct {
		Code   string `json:"code"`
		Name   string `json:"name"`
		Symbol string `json:"symbol"`
	} `json:"currency"`
	IsDefaultAccount bool    `json:"isDefaultAccount"`
	CreatedAt        string  `json:"createdAt"`
	LastSyncAt       *string `json:"lastSyncAt,omitempty"`
}

func init() {
	// Configure logging to include timestamps and file information
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

// LinkAccount links a new account to a user
func (h *LinkedAccountHandler) LinkAccount(c *gin.Context) {
	var req linkAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[LINK-ACCOUNT] Invalid request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{
			"code":    "INVALID_REQUEST",
			"message": err.Error(),
		}})
		return
	}

	// Get user ID from context with detailed logging
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		log.Printf("[LINK-ACCOUNT] user_id not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": gin.H{
			"code":    "UNAUTHORIZED",
			"message": "User ID not found in context",
		}})
		return
	}

	log.Printf("[LINK-ACCOUNT] Retrieved user_id from context: %v (type: %T)", userIDInterface, userIDInterface)

	// Try to get UUID directly
	userID, ok := userIDInterface.(uuid.UUID)
	if !ok {
		// If not UUID, try string conversion as fallback
		if userIDStr, isStr := userIDInterface.(string); isStr {
			var err error
			userID, err = uuid.Parse(userIDStr)
			if err != nil {
				log.Printf("[LINK-ACCOUNT] Failed to parse user_id string as UUID: %v", err)
				c.JSON(http.StatusUnauthorized, gin.H{"error": gin.H{
					"code":    "UNAUTHORIZED",
					"message": "Invalid user ID format",
				}})
				return
			}
		} else {
			log.Printf("[LINK-ACCOUNT] user_id is neither UUID nor string: %v (type: %T)", userIDInterface, userIDInterface)
			c.JSON(http.StatusUnauthorized, gin.H{"error": gin.H{
				"code":    "UNAUTHORIZED",
				"message": "Invalid user ID format",
			}})
			return
		}
	}

	log.Printf("[LINK-ACCOUNT] Successfully got UUID: %s", userID)

	// Check if account already exists (excluding soft-deleted records)
	var existingAccount models.LinkedAccount
	err := h.db.Unscoped().
		Where("user_id = ? AND provider = ? AND account_number = ? AND deleted_at IS NULL",
			userID, req.Provider, req.AccountNumber).
		First(&existingAccount).Error
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": gin.H{
			"code":    "ACCOUNT_ALREADY_LINKED",
			"message": "This account is already linked to your profile",
		}})
		return
	}

	// If we found a soft-deleted account, we can reactivate it
	err = h.db.Unscoped().
		Where("user_id = ? AND provider = ? AND account_number = ? AND deleted_at IS NOT NULL",
			userID, req.Provider, req.AccountNumber).
		First(&existingAccount).Error
	if err == nil {
		// Reactivate the account with new details
		existingAccount.DeletedAt = gorm.DeletedAt{}
		existingAccount.AccountID = req.AccountID
		existingAccount.AccountTitle = req.AccountTitle
		existingAccount.AccountType = req.AccountType
		existingAccount.CurrencyCode = req.Currency.Code
		existingAccount.CurrencyName = req.Currency.Name
		existingAccount.CurrencySymbol = req.Currency.Symbol
		existingAccount.IsDefaultAccount = req.IsDefaultAccount
		existingAccount.ProviderUsername = req.Credentials.Username
		existingAccount.ProviderPassword = req.Credentials.Password
		existingAccount.DeviceID = req.DeviceInfo.DeviceID

		if err := h.db.Unscoped().Save(&existingAccount).Error; err != nil {
			log.Printf("[LINK-ACCOUNT] Failed to reactivate account: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to link account",
			}})
			return
		}

		log.Printf("[LINK-ACCOUNT] Successfully reactivated account for user %s", userID)
		c.JSON(http.StatusCreated, h.toAccountResponse(&existingAccount))
		return
	}

	// Create new linked account if no existing account found
	account := models.LinkedAccount{
		UserID:           userID,
		Provider:         req.Provider,
		AccountID:        req.AccountID,
		AccountNumber:    req.AccountNumber,
		AccountTitle:     req.AccountTitle,
		AccountType:      req.AccountType,
		CurrencyCode:     req.Currency.Code,
		CurrencyName:     req.Currency.Name,
		CurrencySymbol:   req.Currency.Symbol,
		IsDefaultAccount: req.IsDefaultAccount,
		ProviderUsername: req.Credentials.Username,
		ProviderPassword: req.Credentials.Password,
		DeviceID:         req.DeviceInfo.DeviceID,
	}

	if err := h.db.Create(&account).Error; err != nil {
		log.Printf("[LINK-ACCOUNT] Failed to create account: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{
			"code":    "INTERNAL_ERROR",
			"message": "Failed to link account",
		}})
		return
	}

	log.Printf("[LINK-ACCOUNT] Successfully created account for user %s", userID)
	c.JSON(http.StatusCreated, h.toAccountResponse(&account))
}

// GetLinkedAccounts returns all linked accounts for a user
func (h *LinkedAccountHandler) GetLinkedAccounts(c *gin.Context) {
	// Get user ID from context with detailed logging
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		log.Printf("[GET-ACCOUNTS] user_id not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": gin.H{
			"code":    "UNAUTHORIZED",
			"message": "User ID not found in context",
		}})
		return
	}

	log.Printf("[GET-ACCOUNTS] Retrieved user_id from context: %v (type: %T)", userIDInterface, userIDInterface)

	// Try to get UUID directly
	userID, ok := userIDInterface.(uuid.UUID)
	if !ok {
		// If not UUID, try string conversion as fallback
		if userIDStr, isStr := userIDInterface.(string); isStr {
			var err error
			userID, err = uuid.Parse(userIDStr)
			if err != nil {
				log.Printf("[GET-ACCOUNTS] Failed to parse user_id string as UUID: %v", err)
				c.JSON(http.StatusUnauthorized, gin.H{"error": gin.H{
					"code":    "UNAUTHORIZED",
					"message": "Invalid user ID format",
				}})
				return
			}
		} else {
			log.Printf("[GET-ACCOUNTS] user_id is neither UUID nor string: %v (type: %T)", userIDInterface, userIDInterface)
			c.JSON(http.StatusUnauthorized, gin.H{"error": gin.H{
				"code":    "UNAUTHORIZED",
				"message": "Invalid user ID format",
			}})
			return
		}
	}

	log.Printf("[GET-ACCOUNTS] Successfully got UUID: %s", userID)

	var accounts []models.LinkedAccount
	if err := h.db.Where("user_id = ?", userID).Find(&accounts).Error; err != nil {
		log.Printf("[GET-ACCOUNTS] Database query failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{
			"code":    "INTERNAL_ERROR",
			"message": "Failed to fetch accounts",
		}})
		return
	}

	log.Printf("[GET-ACCOUNTS] Found %d accounts for user %s", len(accounts), userID)

	response := make([]accountResponse, len(accounts))
	for i, account := range accounts {
		response[i] = *h.toAccountResponse(&account)
	}

	c.JSON(http.StatusOK, gin.H{"accounts": response})
}

// GetLinkedAccount returns a single linked account by ID
func (h *LinkedAccountHandler) GetLinkedAccount(c *gin.Context) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": gin.H{
			"code":    "UNAUTHORIZED",
			"message": "User ID not found in context",
		}})
		return
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
				return
			}
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": gin.H{
				"code":    "UNAUTHORIZED",
				"message": "Invalid user ID format",
			}})
			return
		}
	}

	accountID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{
			"code":    "INVALID_REQUEST",
			"message": "Invalid account ID",
		}})
		return
	}

	var account models.LinkedAccount
	if err := h.db.Where("id = ? AND user_id = ?", accountID, userID).First(&account).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": gin.H{
			"code":    "NOT_FOUND",
			"message": "Account not found",
		}})
		return
	}

	c.JSON(http.StatusOK, h.toAccountResponse(&account))
}

// UnlinkAccount removes a linked account
func (h *LinkedAccountHandler) UnlinkAccount(c *gin.Context) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": gin.H{
			"code":    "UNAUTHORIZED",
			"message": "User ID not found in context",
		}})
		return
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
				return
			}
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": gin.H{
				"code":    "UNAUTHORIZED",
				"message": "Invalid user ID format",
			}})
			return
		}
	}

	accountID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{
			"code":    "INVALID_REQUEST",
			"message": "Invalid account ID",
		}})
		return
	}

	result := h.db.Where("id = ? AND user_id = ?", accountID, userID).Delete(&models.LinkedAccount{})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{
			"code":    "INTERNAL_ERROR",
			"message": "Failed to unlink account",
		}})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": gin.H{
			"code":    "NOT_FOUND",
			"message": "Account not found",
		}})
		return
	}

	c.Status(http.StatusNoContent)
}

// SetDefaultAccount sets an account as the default for its provider
func (h *LinkedAccountHandler) SetDefaultAccount(c *gin.Context) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": gin.H{
			"code":    "UNAUTHORIZED",
			"message": "User ID not found in context",
		}})
		return
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
				return
			}
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": gin.H{
				"code":    "UNAUTHORIZED",
				"message": "Invalid user ID format",
			}})
			return
		}
	}

	accountID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{
			"code":    "INVALID_REQUEST",
			"message": "Invalid account ID",
		}})
		return
	}

	var account models.LinkedAccount
	if err := h.db.Where("id = ? AND user_id = ?", accountID, userID).First(&account).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": gin.H{
			"code":    "NOT_FOUND",
			"message": "Account not found",
		}})
		return
	}

	account.IsDefaultAccount = true
	if err := h.db.Save(&account).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{
			"code":    "INTERNAL_ERROR",
			"message": "Failed to update account",
		}})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":               account.ID,
		"isDefaultAccount": true,
	})
}

// RefreshAccount refreshes the account data from the provider
func (h *LinkedAccountHandler) RefreshAccount(c *gin.Context) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": gin.H{
			"code":    "UNAUTHORIZED",
			"message": "User ID not found in context",
		}})
		return
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
				return
			}
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": gin.H{
				"code":    "UNAUTHORIZED",
				"message": "Invalid user ID format",
			}})
			return
		}
	}

	accountID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{
			"code":    "INVALID_REQUEST",
			"message": "Invalid account ID",
		}})
		return
	}

	var account models.LinkedAccount
	if err := h.db.Where("id = ? AND user_id = ?", accountID, userID).First(&account).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": gin.H{
			"code":    "NOT_FOUND",
			"message": "Account not found",
		}})
		return
	}

	// TODO: Implement actual provider refresh logic here
	// For now, just update the LastSyncAt timestamp
	now := time.Now()
	account.LastSyncAt = &now
	if err := h.db.Save(&account).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{
			"code":    "INTERNAL_ERROR",
			"message": "Failed to refresh account",
		}})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":         account.ID,
		"lastSyncAt": account.LastSyncAt,
		"status":     "SUCCESS",
	})
}

// Helper function to convert LinkedAccount to accountResponse
func (h *LinkedAccountHandler) toAccountResponse(account *models.LinkedAccount) *accountResponse {
	response := &accountResponse{
		ID:               account.ID,
		Provider:         account.Provider,
		AccountID:        account.AccountID,
		AccountNumber:    account.AccountNumber,
		AccountTitle:     account.AccountTitle,
		AccountType:      account.AccountType,
		IsDefaultAccount: account.IsDefaultAccount,
		CreatedAt:        account.CreatedAt.Format(time.RFC3339),
	}

	response.Currency.Code = account.CurrencyCode
	response.Currency.Name = account.CurrencyName
	response.Currency.Symbol = account.CurrencySymbol

	if account.LastSyncAt != nil {
		lastSyncAt := account.LastSyncAt.Format(time.RFC3339)
		response.LastSyncAt = &lastSyncAt
	}

	return response
}
