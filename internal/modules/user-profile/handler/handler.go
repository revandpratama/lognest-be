package handler

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/revandpratama/lognest/internal/modules/user-profile/entity"
	"github.com/revandpratama/lognest/internal/modules/user-profile/usecase"
	"github.com/revandpratama/lognest/pkg/errorhandler"
	"github.com/revandpratama/lognest/pkg/response"
)

// UserProfileHandler defines the HTTP handler interface for a UserProfile.
type UserProfileHandler interface {
	Create(c *fiber.Ctx) error
	FindByID(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	FindUser(c *fiber.Ctx) error
}

type userprofileHandler struct {
	usecase usecase.UserProfileUsecase
}

// NewUserProfileHandler creates a new instance of UserProfileHandler.
func NewUserProfileHandler(usecase usecase.UserProfileUsecase) UserProfileHandler {
	return &userprofileHandler{usecase: usecase}
}

func (h *userprofileHandler) Create(c *fiber.Ctx) error {
	var newUserProfile entity.UserProfile
	if err := c.BodyParser(&newUserProfile); err != nil {
		return errorhandler.BuildError(c, errorhandler.BadRequestError{Message: err.Error()}, nil)
	}
	userProfile, err := h.usecase.Create(c.Context(), &newUserProfile)
	if err != nil {
		return errorhandler.BuildError(c, err, nil)
	}
	return response.Success(c, fiber.StatusCreated, "user profile created", userProfile)
}

func (h *userprofileHandler) FindByID(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
	defer cancel()

	idStr := c.Params("id")
	if idStr == "" {
		return errorhandler.BuildError(c, errorhandler.BadRequestError{Message: "id is required"}, nil)
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		return errorhandler.BuildError(c, errorhandler.BadRequestError{Message: "invalid id format"}, nil)
	}

	userProfile, err := h.usecase.FindByID(ctx, id)
	if err != nil {
		return errorhandler.BuildError(c, err, nil)
	}

	return response.Success(c, fiber.StatusOK, "user profile found", userProfile)
}

func (h *userprofileHandler) Update(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
	defer cancel()

	idStr := c.Params("id")
	if idStr == "" {
		return errorhandler.BuildError(c, errorhandler.BadRequestError{Message: "id is required"}, nil)
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		return errorhandler.BuildError(c, errorhandler.BadRequestError{Message: "invalid id format"}, nil)
	}

	var updateUserProfile entity.UserProfile
	if err := c.BodyParser(&updateUserProfile); err != nil {
		return errorhandler.BuildError(c, errorhandler.BadRequestError{Message: err.Error()}, nil)
	}

	userProfile, err := h.usecase.Update(ctx, id, &updateUserProfile)
	if err != nil {
		return errorhandler.BuildError(c, err, nil)
	}

	return response.Success(c, fiber.StatusOK, "user profile updated", userProfile)
}

func (h *userprofileHandler) FindUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
	defer cancel()

	tokenStr := c.Cookies("access_token")
	if tokenStr == "" {
		return errorhandler.BuildError(c, errorhandler.UnauthorizedError{Message: "unauthorized, no access token provided"}, nil)
	}

	userProfile, err := h.usecase.FindUser(ctx, tokenStr)
	if err != nil {
		return errorhandler.BuildError(c, err, nil)
	}

	return response.Success(c, fiber.StatusOK, "user profile found", userProfile)
}
