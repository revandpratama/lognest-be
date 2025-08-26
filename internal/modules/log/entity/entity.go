package entity

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/revandpratama/lognest/config"
	interactionEntity "github.com/revandpratama/lognest/internal/modules/interaction/entity"
	"gorm.io/gorm"
)

// Log represents the data structure for a log.
type Log struct {
	ID            uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	UserProfileID uuid.UUID      `gorm:"not null" json:"user_profile_id"`
	ProjectID     uuid.UUID      `gorm:"not null" json:"project_id"`
	Content       string         `gorm:"type:text;not null" json:"content" validate:"required"`
	LikeCount     int            `gorm:"default:0" json:"like_count"`
	CommentCount  int            `gorm:"default:0" json:"comment_count"`
	CreatedAt     time.Time      `gorm:"not null" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"not null" json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`

	Media    []Media                     `gorm:"foreignKey:LogID;references:ID;constraint:OnDelete:CASCADE;" json:"media,omitempty"`
	Comments []interactionEntity.Comment `gorm:"foreignKey:LogID;references:ID;constraint:OnDelete:CASCADE;" json:"comments,omitempty"`
}

type Media struct {
	ID            uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	LogID         uuid.UUID `gorm:"not null" json:"log_id"`
	FilePath      string    `gorm:"type:varchar(255);not null" json:"file_path"`
	ThumbnailPath string    `gorm:"type:varchar(255);not null" json:"thumbnail_path"`
	Type          string    `gorm:"type:enum('image', 'video');not null" json:"type"`
	SortOrder     int       `gorm:"default:0" json:"sort_order"`
}

// TableName sets the table name for the Log.
func (Log) TableName() string {
	return fmt.Sprintf("%s.%s", config.ENV.LOGNEST_SCHEMA, "logs")
}

func (p *Log) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		uuidGenerated, err := uuid.NewV7()
		if err != nil {
			return err
		}
		p.ID = uuidGenerated
	}
	return nil
}
