package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/revandpratama/lognest/config"
	"github.com/revandpratama/lognest/internal/modules/user-profile/dto"
	"github.com/revandpratama/lognest/internal/modules/user-profile/entity"
	"github.com/revandpratama/lognest/internal/modules/user-profile/repository"
	"github.com/revandpratama/lognest/pkg/errorhandler"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

// UserProfileUsecase defines the business logic interface for a UserProfile.
type UserProfileUsecase interface {
	Create(ctx context.Context, newUserProfile *entity.UserProfile) (*entity.UserProfile, error)
	FindByID(ctx context.Context, id uuid.UUID) (*entity.UserProfile, error)
	Update(ctx context.Context, id uuid.UUID, updateUserProfile *entity.UserProfile) (*entity.UserProfile, error)
	FindUser(ctx context.Context, tokenStr string) (*entity.UserProfile, error)
}

type userprofileUsecase struct {
	repo       repository.UserProfileRepository
	httpClient *http.Client
}

// NewUserProfileUsecase creates a new instance of UserProfileUsecase.
func NewUserProfileUsecase(repo repository.UserProfileRepository, httpClient *http.Client) UserProfileUsecase {
	return &userprofileUsecase{
		repo:       repo,
		httpClient: httpClient,
	}
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

func (u *userprofileUsecase) FindUser(ctx context.Context, tokenStr string) (*entity.UserProfile, error) {

	var url = fmt.Sprintf("%s/api/auth/user", config.ENV.AUTH4ME_URL)

	log.Debug().Msg("tokenStr: " + tokenStr)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, errorhandler.InternalServerError{Message: err.Error()}
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", tokenStr)

	resp, err := u.httpClient.Do(req)
	if err != nil {
		return nil, errorhandler.InternalServerError{Message: err.Error()}
	}
	defer resp.Body.Close()

	reponseByte, _ := json.Marshal(resp.Body)
	log.Debug().Msg("response: " + string(reponseByte))

	log.Debug().Msg("response: " + resp.Status)

	if resp.StatusCode != http.StatusOK {
		return nil, errorhandler.InternalServerError{Message: "failed to get user"}
	}

	var user dto.MeResponse
	err = json.NewDecoder(resp.Body).Decode(&user)
	if err != nil {
		return nil, errorhandler.InternalServerError{Message: err.Error()}
	}

	log.Debug().Msg("user: " + user.Data.ID)

	userProfile, err := u.repo.FindByID(ctx, uuid.MustParse(user.Data.ID))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorhandler.NotFoundError{Message: "user profile not found"}
		}
		return nil, errorhandler.InternalServerError{Message: err.Error()}
	}

	userProfile.User = user.Data

	return userProfile, nil
}
