package whatsapp

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"gorm.io/gorm"

	"github.com/moha/kaafipay-backend/internal/models"
)

type WhatsAppProvider struct {
	db        *gorm.DB
	baseURL   string
	apiKey    string
	sessionID string
}

type SessionResponse struct {
	Status  string          `json:"status"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data,omitempty"`
}

type AddSessionResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	QR      string `json:"qr,omitempty"`
}

type SessionStatusResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		ID     string `json:"id"`
		Status string `json:"status"`
	} `json:"data"`
}

type SessionListResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    []struct {
		ID     string `json:"id"`
		Status string `json:"status"`
	} `json:"data"`
}

func NewWhatsAppProvider(db *gorm.DB, baseURL, apiKey, sessionID string) *WhatsAppProvider {
	return &WhatsAppProvider{
		db:        db,
		baseURL:   baseURL,
		apiKey:    apiKey,
		sessionID: sessionID,
	}
}

func (w *WhatsAppProvider) makeRequest(endpoint, method string, body interface{}) (*json.RawMessage, int, error) {
	var buf bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&buf).Encode(body); err != nil {
			return nil, 0, fmt.Errorf("failed to encode request body: %v", err)
		}
	}

	req, err := http.NewRequest(method, w.baseURL+endpoint, &buf)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", w.apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	var rawResponse json.RawMessage
	if err := json.NewDecoder(resp.Body).Decode(&rawResponse); err != nil {
		return nil, resp.StatusCode, fmt.Errorf("failed to decode response: %v", err)
	}

	return &rawResponse, resp.StatusCode, nil
}

func (w *WhatsAppProvider) ListSessions() (*SessionListResponse, error) {
	rawResp, _, err := w.makeRequest("/sessions", http.MethodGet, nil)
	if err != nil {
		return nil, err
	}

	var response SessionListResponse
	if err := json.Unmarshal(*rawResp, &response); err != nil {
		return nil, fmt.Errorf("failed to parse session list response: %v", err)
	}

	return &response, nil
}

func (w *WhatsAppProvider) FindSession(sessionID string) (*SessionStatusResponse, error) {
	rawResp, statusCode, err := w.makeRequest("/sessions/"+sessionID, http.MethodGet, nil)
	if err != nil {
		return nil, err
	}

	// Handle 404 Not Found
	if statusCode == http.StatusNotFound {
		return &SessionStatusResponse{
			Status:  "error",
			Message: "Session not found",
		}, nil
	}

	var response SessionStatusResponse
	if err := json.Unmarshal(*rawResp, &response); err != nil {
		return nil, fmt.Errorf("failed to parse session status response: %v", err)
	}

	return &response, nil
}

func (w *WhatsAppProvider) AddSession(sessionID string, readIncomingMessages, syncFullHistory bool) (*AddSessionResponse, error) {
	body := map[string]interface{}{
		"sessionId":            sessionID,
		"readIncomingMessages": readIncomingMessages,
		"syncFullHistory":      syncFullHistory,
	}

	rawResp, _, err := w.makeRequest("/sessions/add", http.MethodPost, body)
	if err != nil {
		return nil, err
	}

	var response AddSessionResponse
	if err := json.Unmarshal(*rawResp, &response); err != nil {
		return nil, fmt.Errorf("failed to parse add session response: %v", err)
	}

	return &response, nil
}

func (w *WhatsAppProvider) DeleteSession(sessionID string) (*SessionResponse, error) {
	rawResp, _, err := w.makeRequest("/sessions/"+sessionID, http.MethodDelete, nil)
	if err != nil {
		return nil, err
	}

	var response SessionResponse
	if err := json.Unmarshal(*rawResp, &response); err != nil {
		return nil, fmt.Errorf("failed to parse delete session response: %v", err)
	}

	return &response, nil
}

func (w *WhatsAppProvider) GenerateCode(phone string) (string, error) {
	// Generate 6-digit code
	code := fmt.Sprintf("%06d", time.Now().UnixNano()%1000000)
	expirationTime := time.Now().Add(5 * time.Minute)

	mfaCode := &models.MFACode{
		Code:      code,
		Phone:     phone,
		ExpiresAt: expirationTime,
	}

	if err := w.db.Create(mfaCode).Error; err != nil {
		return "", fmt.Errorf("failed to save MFA code: %v", err)
	}

	return code, nil
}

func (w *WhatsAppProvider) SendCode(code, phone string) error {
	jid := fmt.Sprintf("252%s@s.whatsapp.net", phone)
	message := fmt.Sprintf("Your KaafiPay verification code is: %s", code)

	body := map[string]interface{}{
		"jid": jid,
		"message": map[string]string{
			"text": message,
		},
	}

	_, _, err := w.makeRequest("/"+w.sessionID+"/messages/send", http.MethodPost, body)
	if err != nil {
		return fmt.Errorf("failed to send WhatsApp message: %v", err)
	}

	return nil
}

func (w *WhatsAppProvider) VerifyCode(code, phone string) (string, error) {
	var mfaCode models.MFACode
	err := w.db.Where("phone = ? AND code = ? AND expires_at > ?", phone, code, time.Now()).
		Order("created_at DESC").
		First(&mfaCode).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", fmt.Errorf("invalid or expired code")
		}
		return "", fmt.Errorf("failed to verify code: %v", err)
	}

	// Delete the used code
	if err := w.db.Delete(&mfaCode).Error; err != nil {
		return "", fmt.Errorf("failed to delete used code: %v", err)
	}

	// Generate MFA token
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return "", fmt.Errorf("failed to generate token: %v", err)
	}
	mfaToken := hex.EncodeToString(tokenBytes)

	// Save the token
	token := &models.MFAToken{
		Token:     mfaToken,
		Phone:     phone,
		ExpiresAt: time.Now().Add(2 * time.Minute),
	}

	if err := w.db.Create(token).Error; err != nil {
		return "", fmt.Errorf("failed to save token: %v", err)
	}

	return mfaToken, nil
}

func (w *WhatsAppProvider) VerifyToken(token string) (bool, error) {
	var mfaToken models.MFAToken
	err := w.db.Where("token = ? AND expires_at > ?", token, time.Now()).First(&mfaToken).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, fmt.Errorf("failed to verify token: %v", err)
	}

	// Delete the used token
	if err := w.db.Delete(&mfaToken).Error; err != nil {
		return false, fmt.Errorf("failed to delete used token: %v", err)
	}

	return true, nil
}
