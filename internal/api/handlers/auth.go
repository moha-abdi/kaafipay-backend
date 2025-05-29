package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/moha/kaafipay-backend/internal/config"
	"github.com/moha/kaafipay-backend/internal/models"
	"github.com/moha/kaafipay-backend/internal/repository"
	"github.com/moha/kaafipay-backend/internal/utils"
)

type AuthHandler struct {
	cfg      *config.Config
	userRepo repository.UserRepository
}

func NewAuthHandler(cfg *config.Config, userRepo repository.UserRepository) *AuthHandler {
	return &AuthHandler{
		cfg:      cfg,
		userRepo: userRepo,
	}
}

type RegisterRequest struct {
	Phone    string `json:"phone" binding:"required,min=9,max=15"`
	Name     string `json:"name" binding:"required,min=2,max=100"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserResponse struct {
	ID    string `json:"id"`
	Phone string `json:"phone"`
	Name  string `json:"name"`
}

type AuthResponse struct {
	User         UserResponse `json:"user"`
	Token        string       `json:"token"`
	RefreshToken string       `json:"refresh_token"`
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if user already exists
	if _, err := h.userRepo.FindByPhone(req.Phone); err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Phone number already registered"})
		return
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process password"})
		return
	}

	// Create user
	user := &models.User{
		Phone:    req.Phone,
		Name:     req.Name,
		Password: hashedPassword,
	}

	if err := h.userRepo.Create(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Generate tokens
	expiration, _ := time.ParseDuration(h.cfg.JWTExpiration)
	refreshExpiration, _ := time.ParseDuration(h.cfg.RefreshTokenExpiration)

	token, err := utils.GenerateToken(user.ID, user.Phone, h.cfg.JWTSecret, expiration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	refreshToken, err := utils.GenerateToken(user.ID, user.Phone, h.cfg.JWTSecret, refreshExpiration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate refresh token"})
		return
	}

	c.JSON(http.StatusCreated, AuthResponse{
		User: UserResponse{
			ID:    user.ID.String(),
			Phone: user.Phone,
			Name:  user.Name,
		},
		Token:        token,
		RefreshToken: refreshToken,
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find user
	user, err := h.userRepo.FindByPhone(req.Phone)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Verify password
	if !utils.CheckPassword(req.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate tokens
	expiration, _ := time.ParseDuration(h.cfg.JWTExpiration)
	refreshExpiration, _ := time.ParseDuration(h.cfg.RefreshTokenExpiration)

	token, err := utils.GenerateToken(user.ID, user.Phone, h.cfg.JWTSecret, expiration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	refreshToken, err := utils.GenerateToken(user.ID, user.Phone, h.cfg.JWTSecret, refreshExpiration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate refresh token"})
		return
	}

	c.JSON(http.StatusOK, AuthResponse{
		User: UserResponse{
			ID:    user.ID.String(),
			Phone: user.Phone,
			Name:  user.Name,
		},
		Token:        token,
		RefreshToken: refreshToken,
	})
}
