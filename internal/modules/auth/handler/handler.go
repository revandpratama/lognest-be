package handler

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/revandpratama/lognest/config"
	"github.com/revandpratama/lognest/internal/modules/auth/dto"
	"github.com/revandpratama/lognest/internal/modules/auth/usecase"
	"github.com/revandpratama/lognest/pkg/errorhandler"
	"github.com/revandpratama/lognest/pkg/response"
)

// AuthHandler defines the HTTP handler interface for a Auth.
type AuthHandler interface {
	Login(c *fiber.Ctx) error
	Register(c *fiber.Ctx) error
	RefreshToken(c *fiber.Ctx) error
	Logout(c *fiber.Ctx) error
}

type authHandler struct {
	usecase usecase.AuthUsecase
}

// NewAuthHandler creates a new instance of AuthHandler.
func NewAuthHandler(usecase usecase.AuthUsecase) AuthHandler {
	return &authHandler{usecase: usecase}
}

func (u *authHandler) Login(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
	defer cancel()

	var loginRequest dto.LoginRequest
	if err := c.BodyParser(&loginRequest); err != nil {
		return errorhandler.BuildError(c, errorhandler.BadRequestError{Message: err.Error()}, nil)
	}

	res, err := u.usecase.Login(ctx, &loginRequest)
	if err != nil {
		return errorhandler.BuildError(c, err, nil)
	}

	accessTokenCookie := &fiber.Cookie{
		Name:     "access_token",
		Value:    res.Data.AccessToken,
		Expires:  time.Now().Add(time.Minute * 5),
		HTTPOnly: true,
	}

	if config.ENV.APP_ENV == "production" {
		accessTokenCookie.SameSite = "None"
		accessTokenCookie.Domain = ".revandpratama.com"
		accessTokenCookie.Secure = true
	}

	c.Cookie(accessTokenCookie)

	refreshTokenCookie := &fiber.Cookie{
		Name:     "refresh_token",
		Value:    res.Data.RefreshToken,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	}

	if config.ENV.APP_ENV == "production" {
		refreshTokenCookie.SameSite = "None"
		refreshTokenCookie.Domain = ".revandpratama.com"
		refreshTokenCookie.Secure = true
	}

	c.Cookie(refreshTokenCookie)
	return response.Success(c, fiber.StatusOK, "login success", nil)
}

func (u *authHandler) Register(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
	defer cancel()

	var registerRequest dto.RegisterRequest
	if err := c.BodyParser(&registerRequest); err != nil {
		return errorhandler.BuildError(c, errorhandler.BadRequestError{Message: err.Error()}, nil)
	}

	err := u.usecase.Register(ctx, &registerRequest)
	if err != nil {
		return errorhandler.BuildError(c, err, nil)
	}

	return response.Success(c, fiber.StatusOK, "register success", nil)
}

func (u *authHandler) RefreshToken(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
	defer cancel()

	accessToken := c.Cookies("access_token")
	if accessToken == "" {
		return errorhandler.BuildError(c, errorhandler.UnauthorizedError{Message: "unauthorized, no access token provided"}, nil)
	}

	refreshToken := c.Cookies("refresh_token")
	if refreshToken == "" {
		return errorhandler.BuildError(c, errorhandler.UnauthorizedError{Message: "unauthorized, no refresh token provided"}, nil)
	}

	res, err := u.usecase.RefreshToken(ctx, accessToken, refreshToken)
	if err != nil {
		return errorhandler.BuildError(c, err, nil)
	}

	accessTokenCookie := &fiber.Cookie{
		Name:     "access_token",
		Value:    res.Data.AccessToken,
		Expires:  time.Now().Add(time.Minute * 5),
		HTTPOnly: true,
	}

	if config.ENV.APP_ENV == "production" {
		accessTokenCookie.SameSite = "None"
		accessTokenCookie.Domain = ".revandpratama.com"
		accessTokenCookie.Secure = true
	}

	c.Cookie(accessTokenCookie)

	refreshTokenCookie := &fiber.Cookie{
		Name:     "refresh_token",
		Value:    res.Data.RefreshToken,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	}

	if config.ENV.APP_ENV == "production" {
		refreshTokenCookie.SameSite = "None"
		refreshTokenCookie.Domain = ".revandpratama.com"
		refreshTokenCookie.Secure = true
	}

	c.Cookie(refreshTokenCookie)

	return response.Success(c, fiber.StatusOK, "refresh token success", nil)
}

func (u *authHandler) Logout(c *fiber.Ctx) error {
	accessTokenCookie := &fiber.Cookie{
		Name:     "access_token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour), // Set expiration to the past
		HTTPOnly: true,
		Path:     "/",
	}

	if config.ENV.APP_ENV == "production" {
		accessTokenCookie.SameSite = "None"
		accessTokenCookie.Secure = true
	}

	c.Cookie(accessTokenCookie)

	refreshTokenCookie := &fiber.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
		Path:     "/",
	}

	if config.ENV.APP_ENV == "production" {
		refreshTokenCookie.SameSite = "None"
		refreshTokenCookie.Secure = true
	}
	c.Cookie(refreshTokenCookie)
	return response.Success(c, fiber.StatusOK, "logout success", nil)
}
