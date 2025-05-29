package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Provider represents the available account providers
type Provider string

const (
	ProviderZaad     Provider = "ZAAD"
	ProviderEdahab   Provider = "EDAHAB"
	ProviderSahal    Provider = "SAHAL"
	ProviderEvcplus  Provider = "EVCPLUS"
	ProviderSomnet   Provider = "SOMNET"
	ProviderSoltelco Provider = "SOLTELCO"
)

// Currency represents the currency information for an account
type Currency struct {
	Code   string `json:"code"`
	Name   string `json:"name"`
	Symbol string `json:"symbol"`
}

// DeviceInfo represents the device information for account linking
type DeviceInfo struct {
	DeviceID     string `json:"deviceId"`
	DeviceModel  string `json:"deviceModel"`
	Manufacturer string `json:"manufacturer"`
	OSVersion    string `json:"osVersion"`
}

// LinkedAccount represents a user's linked third-party account
type LinkedAccount struct {
	ID               uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID           uuid.UUID `json:"userId" gorm:"type:uuid;index;not null"`
	Provider         Provider  `json:"provider" gorm:"type:account_provider;not null"`
	AccountID        string    `json:"accountId" gorm:"not null"`
	AccountNumber    string    `json:"accountNumber" gorm:"not null"`
	AccountTitle     string    `json:"accountTitle" gorm:"not null"`
	AccountType      string    `json:"accountType" gorm:"not null"`
	CurrencyCode     string    `json:"currencyCode" gorm:"not null"`
	CurrencyName     string    `json:"currencyName" gorm:"not null"`
	CurrencySymbol   string    `json:"currencySymbol" gorm:"not null"`
	IsDefaultAccount bool      `json:"isDefaultAccount" gorm:"default:false"`

	// Provider-specific authentication details (encrypted)
	ProviderUsername string `json:"-" gorm:"not null"`
	ProviderPassword string `json:"-" gorm:"not null"`
	DeviceID         string `json:"-" gorm:"not null"`

	// Additional provider-specific details
	CustomerID     string `json:"customerId,omitempty"`
	SubscriptionID string `json:"subscriptionId,omitempty"`

	// Metadata
	CreatedAt  time.Time      `json:"createdAt" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt  time.Time      `json:"updatedAt" gorm:"default:CURRENT_TIMESTAMP"`
	LastSyncAt *time.Time     `json:"lastSyncAt,omitempty"`
	IsActive   bool           `json:"isActive" gorm:"default:true"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`

	// Relations
	User        User          `json:"-" gorm:"foreignKey:UserID"`
	SyncHistory []AccountSync `json:"-" gorm:"foreignKey:LinkedAccountID"`
}

// AccountSync represents the sync history for a linked account
type AccountSync struct {
	ID              uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	LinkedAccountID uuid.UUID `json:"linkedAccountId" gorm:"type:uuid;index;not null"`
	SyncStatus      string    `json:"syncStatus" gorm:"not null"`
	ErrorMessage    string    `json:"errorMessage,omitempty"`
	CreatedAt       time.Time `json:"createdAt" gorm:"default:CURRENT_TIMESTAMP"`

	// Relations
	LinkedAccount LinkedAccount `json:"-" gorm:"foreignKey:LinkedAccountID"`
}

// BeforeCreate hook to ensure only one default account per user per provider
func (la *LinkedAccount) BeforeCreate(tx *gorm.DB) error {
	if la.IsDefaultAccount {
		// Set all other accounts of this user and provider to non-default
		tx.Model(&LinkedAccount{}).
			Where("user_id = ? AND provider = ? AND is_default_account = ?",
				la.UserID, la.Provider, true).
			Update("is_default_account", false)
	}
	return nil
}

// BeforeUpdate hook to maintain single default account constraint per provider
func (la *LinkedAccount) BeforeUpdate(tx *gorm.DB) error {
	if la.IsDefaultAccount {
		// Set all other accounts of this user and provider to non-default
		tx.Model(&LinkedAccount{}).
			Where("user_id = ? AND provider = ? AND id != ? AND is_default_account = ?",
				la.UserID, la.Provider, la.ID, true).
			Update("is_default_account", false)
	}
	return nil
}

// TableName specifies the table name for the LinkedAccount model
func (LinkedAccount) TableName() string {
	return "linked_accounts"
}

// TableName specifies the table name for the AccountSync model
func (AccountSync) TableName() string {
	return "linked_account_syncs"
}
