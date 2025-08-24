package entity

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/revandpratama/lognest/config"
	"gorm.io/gorm"
)

// UserProfile represents the data structure for a userprofile.
type UserProfile struct {
	UserID uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	Bio    string    `gorm:"type:text" json:"bio"`

	// --- Counters (cached from Redis, synced periodically) ---
	FollowerCount  int `gorm:"default:0" json:"follower_count"`
	FollowingCount int `gorm:"default:0" json:"following_count"`

	CreatedAt time.Time      `gorm:"not null" json:"created_at"`
	UpdatedAt time.Time      `gorm:"not null" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Following []UserProfile `gorm:"many2many:user_followers;foreignKey:UserID;joinForeignKey:FollowerID;References:UserID;joinReferences:FollowingID" json:"following,omitempty"`

	// Users that follow this user
	Followers []UserProfile `gorm:"many2many:user_followers;foreignKey:UserID;joinForeignKey:FollowingID;References:UserID;joinReferences:FollowerID" json:"followers,omitempty"`

	User User  `gorm:"-" json:"user,omitzero"`
}

type User struct {
	ID            uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	Email         string    `gorm:"uniqueIndex;not null" json:"email"`
	FirstName     string    `gorm:"size:255" json:"first_name"`
	LastName      string    `gorm:"size:255" json:"last_name"`
	AvatarPath    string    `gorm:"size:500" json:"avatar_path"`
	EmailVerified *bool     `gorm:"default:false" json:"email_verified"`
}

// TableName sets the table name for the UserProfile.
func (UserProfile) TableName() string {
	return fmt.Sprintf("%s.%s", config.ENV.LOGNEST_SCHEMA, "user_profiles")
}

// func (User) TableName() string {
// 	return fmt.Sprintf("%s.%s", config.ENV.AUTH4ME_SCHEMA, "users")
// }
