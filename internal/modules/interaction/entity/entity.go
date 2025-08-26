package entity

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/revandpratama/lognest/config"
	"gorm.io/gorm"
)

// Interaction represents the data structure for a interaction.
type Comment struct {
	ID            uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	UserProfileID uuid.UUID      `gorm:"type:uuid" json:"user_profile_id"`
	LogID         uuid.UUID      `gorm:"type:uuid" json:"log_id"`
	Body          string         `gorm:"type:text;not null" json:"body" validate:"required,min=1,max=255"`
	CreatedAt     time.Time      `gorm:"not null" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"not null" json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName sets the table name for the Interaction.
func (Comment) TableName() string {
	return fmt.Sprintf("%s.%s", config.ENV.LOGNEST_SCHEMA, "comments")
}

func (p *Comment) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		uuidGenerated, err := uuid.NewV7()
		if err != nil {
			return err
		}
		p.ID = uuidGenerated
	}
	return nil
}

type Like struct {
	UserProfileID uuid.UUID `gorm:"type:uuid" json:"user_profile_id"`
	LogID         uuid.UUID `gorm:"type:uuid" json:"log_id"`
}

// TableName sets the table name for the Interaction.
func (Like) TableName() string {
	return fmt.Sprintf("%s.%s", config.ENV.LOGNEST_SCHEMA, "likes")
}
