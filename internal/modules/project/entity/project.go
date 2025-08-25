package entity

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/revandpratama/lognest/config"
	logEntity "github.com/revandpratama/lognest/internal/modules/log/entity"
	tagEntity "github.com/revandpratama/lognest/internal/modules/tag/entity"
	userProfileEntity "github.com/revandpratama/lognest/internal/modules/user-profile/entity"
	"gorm.io/gorm"
)

type Project struct {
	ID             uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	UserProfileID  uuid.UUID      `gorm:"not null" json:"user_profile_id"`
	Title          string         `gorm:"type:varchar(255);not null" json:"title" validate:"required,min=5,max=255"`
	Description    string         `gorm:"type:text" json:"description"`
	Slug           string         `gorm:"type:varchar(255);uniqueIndex;not null" json:"slug"`
	CoverImagePath string         `gorm:"type:varchar(255)" json:"cover_image_path"`
	IsPublic       *bool          `gorm:"default:true" json:"is_public"`
	CreatedAt      time.Time      `gorm:"not null" json:"created_at"`
	UpdatedAt      time.Time      `gorm:"not null" json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`

	// --- Relationships ---
	UserProfile userProfileEntity.UserProfile `gorm:"foreignKey:UserID;references:ID" json:"user_profile"`
	Logs        []logEntity.Log               `gorm:"foreignKey:ProjectID;references:ID;constraint:OnDelete:CASCADE;" json:"logs,omitempty"` // CASCADE means if project is deleted, its logs are too
	Tags        []tagEntity.Tag               `gorm:"many2many:lognest.project_tags;constraint:OnDelete:CASCADE;" json:"tags,omitempty"`
}

func (Project) TableName() string {
	return fmt.Sprintf("%s.%s", config.ENV.LOGNEST_SCHEMA, "projects")
}

func (p *Project) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		uuidGenerated, err := uuid.NewV7()
		if err != nil {
			return err
		}
		p.ID = uuidGenerated
	}
	return nil
}

type ProjectTag struct {
	ProjectID uuid.UUID `gorm:"type:uuid;" json:"project_id"`
	TagID     uuid.UUID `gorm:"type:uuid;" json:"tag_id"`
}

func (ProjectTag) TableName() string {
	return fmt.Sprintf("%s.%s", config.ENV.LOGNEST_SCHEMA, "project_tags")
}
