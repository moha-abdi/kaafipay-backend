package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/moha/kaafipay-backend/internal/services/whatsapp"
)

type VerifyHandler struct {
	whatsappProvider *whatsapp.WhatsAppProvider
}

func NewVerifyHandler(whatsappProvider *whatsapp.WhatsAppProvider) *VerifyHandler {
	return &VerifyHandler{
		whatsappProvider: whatsappProvider,
	}
}

type SendCodeRequest struct {
	Phone string `json:"phone" binding:"required,min=9,max=9"`
}

type VerifyCodeRequest struct {
	Phone string `json:"phone" binding:"required,min=9,max=9"`
	Code  string `json:"code" binding:"required,len=6"`
}

type VerifyTokenRequest struct {
	Token string `json:"token" binding:"required,len=64"`
}

func (h *VerifyHandler) SendCode(c *gin.Context) {
	var req SendCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	code, err := h.whatsappProvider.GenerateCode(req.Phone)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate verification code"})
		return
	}

	if err := h.whatsappProvider.SendCode(code, req.Phone); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send verification code"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Verification code sent"})
}

func (h *VerifyHandler) VerifyCode(c *gin.Context) {
	var req VerifyCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.whatsappProvider.VerifyCode(req.Code, req.Phone)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid verification code"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (h *VerifyHandler) VerifyToken(c *gin.Context) {
	var req VerifyTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	valid, err := h.whatsappProvider.VerifyToken(req.Token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify token"})
		return
	}

	if !valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"valid": true})
}
