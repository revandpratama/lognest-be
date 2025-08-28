package dto

import (
	"time"

	"gorm.io/gorm"
)

type MeResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    User   `json:"data"`
}

type User struct {
	// ID         string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	ID         string `gorm:"primaryKey;type:uuid" json:"id"`
	Email      string `gorm:"uniqueIndex;not null" json:"email"`
	Password   string `gorm:"" json:"-"` // hashed password; can be empty for OAuth users
	FirstName  string `gorm:"column:first_name" json:"first_name"`
	LastName   string `gorm:"column:last_name" json:"last_name"`
	AvatarPath string `gorm:"size:500" json:"avatar_path"`

	Providers []OAuthProvider `gorm:"foreignKey:UserID" json:"providers,omitempty"`

	EmailVerified      bool      `gorm:"default:false" json:"email_verified"`
	VerificationToken  string    `gorm:"size:255" json:"-"`
	VerificationSentAt time.Time `json:"-"`

	MFAEnabled bool   `gorm:"default:false" json:"mfa_enabled"`
	MFASecret  string `gorm:"size:255" json:"-"`

	RoleID uint `gorm:"not null" json:"role_id"`

	// Timestamps
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type OAuthProvider struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	UserID       string    `gorm:"type:uuid;index" json:"user_id"`
	Provider     string    `gorm:"size:100;index" json:"provider"`
	ProviderID   string    `gorm:"size:255;index" json:"provider_id"`
	AccessToken  string    `gorm:"size:500" json:"-"`
	RefreshToken string    `gorm:"size:500" json:"-"`
	ExpiresAt    time.Time `json:"expires_at"`
}
