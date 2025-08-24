package usecase

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/revandpratama/lognest/internal/modules/user-profile/entity"
	"github.com/revandpratama/lognest/internal/modules/user-profile/repository"
	"github.com/revandpratama/lognest/pkg/errorhandler"
	"gorm.io/gorm"
)

// UserProfileUsecase defines the business logic interface for a UserProfile.
type UserProfileUsecase interface {
	Create(ctx context.Context, newUserProfile *entity.UserProfile) (*entity.UserProfile, error)
	FindByID(ctx context.Context, id uuid.UUID) (*entity.UserProfile, error)
	Update(ctx context.Context, id uuid.UUID, updateUserProfile *entity.UserProfile) (*entity.UserProfile, error)
}

type userprofileUsecase struct {
	repo repository.UserProfileRepository
}

// NewUserProfileUsecase creates a new instance of UserProfileUsecase.
func NewUserProfileUsecase(repo repository.UserProfileRepository) UserProfileUsecase {
	return &userprofileUsecase{repo: repo}
}

func (u *userprofileUsecase) Create(ctx context.Context, newUserProfile *entity.UserProfile) (*entity.UserProfile, error) {

	return u.repo.Create(ctx, newUserProfile)
}

func (u *userprofileUsecase) FindByID(ctx context.Context, id uuid.UUID) (*entity.UserProfile, error) {
	userProfile, err := u.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorhandler.NotFoundError{Message: "user profile not found"}
		}
		return nil, errorhandler.InternalServerError{Message: err.Error()}
	}
	return userProfile, nil
}

func (u *userprofileUsecase) Update(ctx context.Context, id uuid.UUID, updateUserProfile *entity.UserProfile) (*entity.UserProfile, error) {
	return u.repo.Update(ctx, id, updateUserProfile)
}
