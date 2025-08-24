package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/revandpratama/lognest/internal/modules/user-profile/entity"
	"gorm.io/gorm"
)

// UserProfileRepository defines the interface for database operations for a UserProfile.
type UserProfileRepository interface {
	Create(ctx context.Context, newUserProfile *entity.UserProfile) (*entity.UserProfile, error)
	FindByID(ctx context.Context, id uuid.UUID) (*entity.UserProfile, error)
	Update(ctx context.Context, id uuid.UUID, updateUserProfile *entity.UserProfile) (*entity.UserProfile, error)
}

type userprofileRepository struct {
	db *gorm.DB
}

// NewUserProfileRepository creates a new instance of UserProfileRepository.
func NewUserProfileRepository(db *gorm.DB) UserProfileRepository {
	return &userprofileRepository{db: db}
}

func (r *userprofileRepository) Create(ctx context.Context, newUserProfile *entity.UserProfile) (*entity.UserProfile, error) {

	err := r.db.WithContext(ctx).Create(newUserProfile).Error
	if err != nil {
		return nil, err
	}

	return newUserProfile, nil
}

func (r *userprofileRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.UserProfile, error) {
	var userProfile entity.UserProfile
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&userProfile).Error; err != nil {
		return nil, err
	}
	return &userProfile, nil
}

func (r *userprofileRepository) Update(ctx context.Context, id uuid.UUID, updateUserProfile *entity.UserProfile) (*entity.UserProfile, error) {
	err := r.db.WithContext(ctx).Model(&entity.UserProfile{}).Where("id = ?", id).Updates(updateUserProfile).Error
	return updateUserProfile, err
}
