package usecase

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/revandpratama/lognest/config"
	"github.com/revandpratama/lognest/internal/modules/auth/dto"
	"github.com/revandpratama/lognest/internal/modules/auth/repository"
	userProfileEntity "github.com/revandpratama/lognest/internal/modules/user-profile/entity"
	userProfileRepository "github.com/revandpratama/lognest/internal/modules/user-profile/repository"
	"github.com/revandpratama/lognest/pkg/errorhandler"
)

// AuthUsecase defines the business logic interface for a Auth.
type AuthUsecase interface {
	Login(ctx context.Context, loginRequest *dto.LoginRequest) (*dto.LoginResponse, error)
	Register(ctx context.Context, registerRequest *dto.RegisterRequest) (map[string]any, error)
	RefreshToken(ctx context.Context, accessToken string, refreshToken string) (*dto.LoginResponse, error)
}

type authUsecase struct {
	repo            repository.AuthRepository
	userProfileRepo userProfileRepository.UserProfileRepository
	httpClient      *http.Client
}

// NewAuthUsecase creates a new instance of AuthUsecase.
func NewAuthUsecase(repo repository.AuthRepository, userProfileRepo userProfileRepository.UserProfileRepository, httpClient *http.Client) AuthUsecase {
	return &authUsecase{
		repo:            repo,
		httpClient:      httpClient,
		userProfileRepo: userProfileRepo,
	}
}

func (u *authUsecase) Login(ctx context.Context, loginRequest *dto.LoginRequest) (*dto.LoginResponse, error) {

	jsonData, err := json.Marshal(loginRequest)
	if err != nil {
		return nil, errorhandler.InternalServerError{Message: err.Error()}
	}

	var url = fmt.Sprintf("%s/api/auth/login", config.ENV.AUTH4ME_URL)

	resp, err := u.httpClient.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, errorhandler.InternalServerError{Message: err.Error()}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errorhandler.InternalServerError{Message: "failed to login"}
	}

	var loginResponse dto.LoginResponse
	err = json.NewDecoder(resp.Body).Decode(&loginResponse)
	if err != nil {
		return nil, errorhandler.InternalServerError{Message: err.Error()}
	}

	return &loginResponse, nil
}

func (u *authUsecase) Register(ctx context.Context, registerRequest *dto.RegisterRequest) (map[string]any, error) {
	jsonData, err := json.Marshal(registerRequest)
	if err != nil {
		return nil, errorhandler.InternalServerError{Message: "failed to marshal request: " + err.Error()}
	}

	url := fmt.Sprintf("%s/api/auth/register", config.ENV.AUTH4ME_URL)
	resp, err := u.httpClient.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, errorhandler.InternalServerError{Message: "http request failed: " + err.Error()}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errorhandler.InternalServerError{Message: "failed to read response body: " + err.Error()}
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errorhandler.InternalServerError{Message: "failed to register with auth service"}
	}

	var result map[string]any
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, errorhandler.InternalServerError{Message: "could not parse auth service response: " + err.Error()}
	}

	dataField, dataExists := result["data"]
	if !dataExists {
		return nil, errorhandler.InternalServerError{Message: "auth service response is missing 'data' field"}
	}

	if dataField == nil {
		return nil, errorhandler.InternalServerError{Message: "auth service returned null 'data' field"}
	}

	data, ok := dataField.(map[string]any)
	if !ok {
		return nil, errorhandler.InternalServerError{Message: "auth service 'data' field is not a JSON object"}
	}

	userIDStr, ok := data["id"].(string)
	if !ok {
		return nil, errorhandler.InternalServerError{Message: "user ID not found or not a string in auth service response"}
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, errorhandler.InternalServerError{Message: "invalid user ID format from auth service"}
	}

	newUserProfile := &userProfileEntity.UserProfile{
		UserID:     userID,
		FirstName:  registerRequest.FirstName,
		Email:      registerRequest.Email,
		LastName:   registerRequest.LastName,
		AvatarPath: registerRequest.AvatarPath,
	}

	_, err = u.userProfileRepo.Create(ctx, newUserProfile)
	if err != nil {
		log.Printf("CRITICAL: Failed to create user profile for user ID %s after successful registration", userID)
		return nil, errorhandler.InternalServerError{Message: "failed to create user profile"}
	}

	return data, nil
}

func (u *authUsecase) RefreshToken(ctx context.Context, accessToken string, refreshToken string) (*dto.LoginResponse, error) {

	var url = fmt.Sprintf("%s/api/auth/refresh-token", config.ENV.AUTH4ME_URL)

	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return nil, errorhandler.InternalServerError{Message: err.Error()}
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", accessToken)
	req.Header.Set("X-Refresh-Token", refreshToken)

	resp, err := u.httpClient.Do(req)
	if err != nil {
		return nil, errorhandler.InternalServerError{Message: err.Error()}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errorhandler.InternalServerError{Message: "failed to refresh token"}
	}

	var loginResponse dto.LoginResponse
	err = json.NewDecoder(resp.Body).Decode(&loginResponse)
	if err != nil {
		return nil, errorhandler.InternalServerError{Message: err.Error()}
	}

	return &loginResponse, nil
}
