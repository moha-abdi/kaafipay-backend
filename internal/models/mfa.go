package models

import (
	"time"

	"github.com/google/uuid"
)

type MFACode struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Code      string    `gorm:"type:varchar(6);not null" json:"code"`
	Phone     string    `gorm:"type:varchar(50);not null" json:"phone"`
	ExpiresAt time.Time `gorm:"not null" json:"expires_at"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
}

type MFAToken struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Token     string    `gorm:"type:varchar(64);not null;unique" json:"token"`
	Phone     string    `gorm:"type:varchar(50);not null" json:"phone"`
	ExpiresAt time.Time `gorm:"not null" json:"expires_at"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
}

type WhatsAppSession struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	SessionID string    `gorm:"type:varchar(100);not null;unique" json:"session_id"`
	Status    string    `gorm:"type:varchar(20);not null;default:'inactive'" json:"status"`
	Metadata  JSON      `gorm:"type:jsonb" json:"metadata,omitempty"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// TableName specifies the table names for GORM
func (MFACode) TableName() string {
	return "mfa_codes"
}

func (MFAToken) TableName() string {
	return "mfa_tokens"
}

func (WhatsAppSession) TableName() string {
	return "whatsapp_sessions"
}
