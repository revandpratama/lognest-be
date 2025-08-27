package entity

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/revandpratama/lognest/config"
	"gorm.io/gorm"
)

// Auth represents the data structure for a auth.
type Auth struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	CreatedAt   time.Time      `gorm:"not null" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"not null" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName sets the table name for the Auth.
func (Auth) TableName() string {
	return fmt.Sprintf("%s.%s", config.ENV.LOGNEST_SCHEMA, "auths")
}

func (p *Auth) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		uuidGenerated, err := uuid.NewV7()
		if err != nil {
			return err
		}
		p.ID = uuidGenerated
	}
	return nil
}
