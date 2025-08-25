package handler

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/revandpratama/lognest/internal/modules/tag/entity"
	"github.com/revandpratama/lognest/internal/modules/tag/usecase"
	"github.com/revandpratama/lognest/pkg/errorhandler"
	"github.com/revandpratama/lognest/pkg/pagination"
	"github.com/revandpratama/lognest/pkg/response"
)

// TagHandler defines the HTTP handler interface for a Tag.
type TagHandler interface {
	FindAll(c *fiber.Ctx) error
	FindByID(c *fiber.Ctx) error
	Create(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
}

type tagHandler struct {
	usecase usecase.TagUsecase
}

// NewTagHandler creates a new instance of TagHandler.
func NewTagHandler(usecase usecase.TagUsecase) TagHandler {
	return &tagHandler{usecase: usecase}
}

func (h *tagHandler) FindAll(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
	defer cancel()

	paginationQuery := new(pagination.Pagination)
	if err := c.QueryParser(paginationQuery); err != nil {
		return errorhandler.BuildError(c, errorhandler.BadRequestError{Message: err.Error()}, nil)
	}

	tags, pagination, err := h.usecase.FindAll(ctx, paginationQuery)
	if err != nil {
		return errorhandler.BuildError(c, err, nil)
	}

	return response.Paginated(c, fiber.StatusOK, "tags found", tags, pagination)

}

func (h *tagHandler) FindByID(c *fiber.Ctx) error {
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

	tag, err := h.usecase.FindByID(ctx, id)
	if err != nil {
		return errorhandler.BuildError(c, err, nil)
	}

	return response.Success(c, fiber.StatusOK, "tag found", tag)

}

func (h *tagHandler) Create(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
	defer cancel()

	var newTag entity.Tag
	if err := c.BodyParser(&newTag); err != nil {
		return errorhandler.BuildError(c, errorhandler.BadRequestError{Message: err.Error()}, nil)
	}

	tag, err := h.usecase.Create(ctx, &newTag)
	if err != nil {
		return errorhandler.BuildError(c, err, nil)
	}

	return response.Success(c, fiber.StatusCreated, "tag created", tag)

}


func (h *tagHandler) Update(c *fiber.Ctx) error {
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

	var updateTag entity.Tag
	if err := c.BodyParser(&updateTag); err != nil {
		return errorhandler.BuildError(c, errorhandler.BadRequestError{Message: err.Error()}, nil)
	}

	tag, err := h.usecase.Update(ctx, id, &updateTag)
	if err != nil {
		return errorhandler.BuildError(c, err, nil)
	}

	return response.Success(c, fiber.StatusOK, "tag updated", tag)
}

func (h *tagHandler) Delete(c *fiber.Ctx) error {
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

	err = h.usecase.Delete(ctx, id)
	if err != nil {
		return errorhandler.BuildError(c, err, nil)
	}

	return response.Success(c, fiber.StatusOK, "tag deleted", nil)
}

