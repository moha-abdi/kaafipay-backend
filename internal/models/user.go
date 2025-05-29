package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID                uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Phone             string    `gorm:"type:varchar(50);unique;not null" json:"phone"`
	Name              string    `gorm:"type:varchar(100);not null" json:"name"`
	Password          string    `gorm:"type:varchar(255);not null" json:"-"`
	CountryCode       string    `gorm:"type:varchar(2)" json:"country_code"`
	PreferredCurrency string    `gorm:"type:varchar(3);default:USD" json:"preferred_currency"`
	IsActive          bool      `gorm:"default:true" json:"is_active"`
	CreatedAt         time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt         time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// BeforeCreate will set a UUID rather than numeric ID.
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

// TableName specifies the table name for the User model
func (User) TableName() string {
	return "users"
}
