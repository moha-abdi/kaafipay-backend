package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/moha/kaafipay-backend/internal/services/whatsapp"
)

type AdminHandler struct {
	whatsappProvider *whatsapp.WhatsAppProvider
}

func NewAdminHandler(whatsappProvider *whatsapp.WhatsAppProvider) *AdminHandler {
	return &AdminHandler{
		whatsappProvider: whatsappProvider,
	}
}

type AddSessionRequest struct {
	SessionID            string `json:"session_id" binding:"required"`
	ReadIncomingMessages bool   `json:"read_incoming_messages"`
	SyncFullHistory      bool   `json:"sync_full_history"`
}

type SessionResponse struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

// ListSessions returns all WhatsApp sessions
func (h *AdminHandler) ListSessions(c *gin.Context) {
	response, err := h.whatsappProvider.ListSessions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to list sessions: %v", err)})
		return
	}

	sessions := make([]SessionResponse, 0)
	for _, s := range response.Data {
		sessions = append(sessions, SessionResponse{
			ID:     s.ID,
			Status: s.Status,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"status":   response.Status,
		"message":  response.Message,
		"sessions": sessions,
	})
}

// GetSession returns details of a specific session
func (h *AdminHandler) GetSession(c *gin.Context) {
	sessionID := c.Param("sessionId")
	response, err := h.whatsappProvider.FindSession(sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to get session: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  response.Status,
		"message": response.Message,
		"session": SessionResponse{
			ID:     response.Data.ID,
			Status: response.Data.Status,
		},
	})
}

// AddSession creates a new WhatsApp session
func (h *AdminHandler) AddSession(c *gin.Context) {
	var req AddSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.whatsappProvider.AddSession(req.SessionID, req.ReadIncomingMessages, req.SyncFullHistory)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to add session: %v", err)})
		return
	}

	// If QR code is present, return it for scanning
	if response.QR != "" {
		c.JSON(http.StatusOK, gin.H{
			"status":  response.Status,
			"message": response.Message,
			"qr":      response.QR,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  response.Status,
		"message": response.Message,
	})
}

// DeleteSession removes a WhatsApp session
func (h *AdminHandler) DeleteSession(c *gin.Context) {
	sessionID := c.Param("sessionId")
	response, err := h.whatsappProvider.DeleteSession(sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to delete session: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  response.Status,
		"message": response.Message,
	})
}
