package handler

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/revandpratama/lognest/internal/modules/interaction/entity"
	"github.com/revandpratama/lognest/internal/modules/interaction/usecase"
	"github.com/revandpratama/lognest/pkg/errorhandler"
	"github.com/revandpratama/lognest/pkg/response"
)

// InteractionHandler defines the HTTP handler interface for a Interaction.
type InteractionHandler interface {
	CreateLike(c *fiber.Ctx) error
	DeleteLike(c *fiber.Ctx) error
	FindLikeByLogID(c *fiber.Ctx) error
	CreateComment(c *fiber.Ctx) error
	UpdateComment(c *fiber.Ctx) error
	FindCommentByLogID(c *fiber.Ctx) error
	DeleteComment(c *fiber.Ctx) error
}

type interactionHandler struct {
	usecase usecase.InteractionUsecase
}

// NewInteractionHandler creates a new instance of InteractionHandler.
func NewInteractionHandler(usecase usecase.InteractionUsecase) InteractionHandler {
	return &interactionHandler{usecase: usecase}
}

func (h *interactionHandler) CreateLike(c *fiber.Ctx) error {

	ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
	defer cancel()

	var newLike entity.Like
	if err := c.BodyParser(&newLike); err != nil {
		return errorhandler.BuildError(c, errorhandler.BadRequestError{Message: err.Error()}, nil)
	}

	like, err := h.usecase.CreateLike(ctx, &newLike)
	if err != nil {
		return errorhandler.BuildError(c, err, nil)
	}

	return response.Success(c, fiber.StatusCreated, "like created", like)
}

func (h *interactionHandler) DeleteLike(c *fiber.Ctx) error {

	ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
	defer cancel()

	userIDStr, ok := c.Locals("userID").(string)
	if !ok {
		return errorhandler.BuildError(c, errorhandler.UnauthorizedError{Message: "unauthorized: userID not found"}, nil)
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return errorhandler.BuildError(c, errorhandler.BadRequestError{Message: "invalid userID format"}, nil)
	}

	logIDStr := c.Params("logID")
	if logIDStr == "" {
		return errorhandler.BuildError(c, errorhandler.BadRequestError{Message: "logID is required"}, nil)
	}

	logID, err := uuid.Parse(logIDStr)
	if err != nil {
		return errorhandler.BuildError(c, errorhandler.BadRequestError{Message: "invalid logID format"}, nil)
	}

	err = h.usecase.DeleteLike(ctx, userID, logID)
	if err != nil {
		return errorhandler.BuildError(c, err, nil)
	}

	return response.Success(c, fiber.StatusOK, "like deleted", nil)
}

func (h *interactionHandler) FindLikeByLogID(c *fiber.Ctx) error {

	ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
	defer cancel()

	logIDStr := c.Params("logID")
	if logIDStr == "" {
		return errorhandler.BuildError(c, errorhandler.BadRequestError{Message: "logID is required"}, nil)
	}

	logID, err := uuid.Parse(logIDStr)
	if err != nil {
		return errorhandler.BuildError(c, errorhandler.BadRequestError{Message: "invalid logID format"}, nil)
	}

	likes, err := h.usecase.FindLikeByLogID(ctx, logID)
	if err != nil {
		return errorhandler.BuildError(c, err, nil)
	}

	return response.Success(c, fiber.StatusOK, "likes found", likes)
}

func (h *interactionHandler) CreateComment(c *fiber.Ctx) error {

	ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
	defer cancel()

	var newComment entity.Comment
	if err := c.BodyParser(&newComment); err != nil {
		return errorhandler.BuildError(c, errorhandler.BadRequestError{Message: err.Error()}, nil)
	}

	comment, err := h.usecase.CreateComment(ctx, &newComment)
	if err != nil {
		return errorhandler.BuildError(c, err, nil)
	}

	return response.Success(c, fiber.StatusCreated, "comment created", comment)
}

func (h *interactionHandler) UpdateComment(c *fiber.Ctx) error {

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

	var newComment entity.Comment
	if err := c.BodyParser(&newComment); err != nil {
		return errorhandler.BuildError(c, errorhandler.BadRequestError{Message: err.Error()}, nil)
	}

	comment, err := h.usecase.UpdateComment(ctx, id, &newComment)
	if err != nil {
		return errorhandler.BuildError(c, err, nil)
	}

	return response.Success(c, fiber.StatusOK, "comment updated", comment)
}

func (h *interactionHandler) FindCommentByLogID(c *fiber.Ctx) error {

	ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
	defer cancel()

	logIDStr := c.Params("logID")
	if logIDStr == "" {
		return errorhandler.BuildError(c, errorhandler.BadRequestError{Message: "logID is required"}, nil)
	}

	logID, err := uuid.Parse(logIDStr)
	if err != nil {
		return errorhandler.BuildError(c, errorhandler.BadRequestError{Message: "invalid logID format"}, nil)
	}

	comments, err := h.usecase.FindCommentByLogID(ctx, logID)
	if err != nil {
		return errorhandler.BuildError(c, err, nil)
	}

	return response.Success(c, fiber.StatusOK, "comments found", comments)
}

func (h *interactionHandler) DeleteComment(c *fiber.Ctx) error {

	ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
	defer cancel()

	// userIDStr, ok := c.Locals("userID").(string)
	// if !ok {
	// 	return errorhandler.BuildError(c, errorhandler.UnauthorizedError{Message: "unauthorized: userID not found"}, nil)
	// }

	// userID, err := uuid.Parse(userIDStr)
	// if err != nil {
	// 	return errorhandler.BuildError(c, errorhandler.BadRequestError{Message: "invalid userID format"}, nil)
	// }

	idStr := c.Params("id")
	if idStr == "" {
		return errorhandler.BuildError(c, errorhandler.BadRequestError{Message: "id is required"}, nil)
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		return errorhandler.BuildError(c, errorhandler.BadRequestError{Message: "invalid id format"}, nil)
	}

	err = h.usecase.DeleteComment(ctx, id)
	if err != nil {
		return errorhandler.BuildError(c, err, nil)
	}

	return response.Success(c, fiber.StatusOK, "comment deleted", nil)
}
