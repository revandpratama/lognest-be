package usecase

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/revandpratama/lognest/config"
	"github.com/revandpratama/lognest/internal/modules/auth/dto"
	"github.com/revandpratama/lognest/internal/modules/auth/repository"
	"github.com/revandpratama/lognest/pkg/errorhandler"
)

// AuthUsecase defines the business logic interface for a Auth.
type AuthUsecase interface {
	Login(ctx context.Context, loginRequest *dto.LoginRequest) (*dto.LoginResponse, error)
	Register(ctx context.Context, registerRequest *dto.RegisterRequest) error
	RefreshToken(ctx context.Context, accessToken string, refreshToken string) (*dto.LoginResponse, error)
}

type authUsecase struct {
	repo       repository.AuthRepository
	httpClient *http.Client
}

// NewAuthUsecase creates a new instance of AuthUsecase.
func NewAuthUsecase(repo repository.AuthRepository, httpClient *http.Client) AuthUsecase {
	return &authUsecase{
		repo:       repo,
		httpClient: httpClient,
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

func (u *authUsecase) Register(ctx context.Context, registerRequest *dto.RegisterRequest) error {

	jsonData, err := json.Marshal(registerRequest)
	if err != nil {
		return errorhandler.InternalServerError{Message: err.Error()}
	}

	var url = fmt.Sprintf("%s/api/auth/register", config.ENV.AUTH4ME_URL)

	resp, err := u.httpClient.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return errorhandler.InternalServerError{Message: err.Error()}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errorhandler.InternalServerError{Message: "failed to register"}
	}

	return nil
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
