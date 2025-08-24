package handler

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/revandpratama/lognest/internal/modules/log/entity"
	"github.com/revandpratama/lognest/internal/modules/log/usecase"
	"github.com/revandpratama/lognest/pkg/errorhandler"
	"github.com/revandpratama/lognest/pkg/pagination"
	"github.com/revandpratama/lognest/pkg/response"
)

// LogHandler defines the HTTP handler interface for a Log.
type LogHandler interface {
	FindByID(c *fiber.Ctx) error
	FindByProjectID(c *fiber.Ctx) error
	Create(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
}

type logHandler struct {
	usecase usecase.LogUsecase
}

// NewLogHandler creates a new instance of LogHandler.
func NewLogHandler(usecase usecase.LogUsecase) LogHandler {
	return &logHandler{usecase: usecase}
}

func (h *logHandler) FindByID(c *fiber.Ctx) error {
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

	log, err := h.usecase.FindByID(ctx, id)
	if err != nil {
		return errorhandler.BuildError(c, err, nil)
	}

	return response.Success(c, fiber.StatusOK, "log found", log)
}

func (h *logHandler) FindByProjectID(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
	defer cancel()

	projectIDStr := c.Params("projectID")
	if projectIDStr == "" {
		return errorhandler.BuildError(c, errorhandler.BadRequestError{Message: "id is required"}, nil)
	}

	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		return errorhandler.BuildError(c, errorhandler.BadRequestError{Message: "invalid id format"}, nil)
	}

	paginationQuery := new(pagination.Pagination)
	if err := c.QueryParser(paginationQuery); err != nil {
		return errorhandler.BuildError(c, errorhandler.BadRequestError{Message: err.Error()}, nil)
	}

	logs, pagination, err := h.usecase.FindByProjectID(ctx, projectID, paginationQuery)
	if err != nil {
		return errorhandler.BuildError(c, err, nil)
	}

	return response.Paginated(c, fiber.StatusOK, "logs found", logs, pagination)
}

func (h *logHandler) Create(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
	defer cancel()

	var newLog entity.Log
	if err := c.BodyParser(&newLog); err != nil {
		return errorhandler.BuildError(c, errorhandler.BadRequestError{Message: err.Error()}, nil)
	}

	log, err := h.usecase.Create(ctx, &newLog)
	if err != nil {
		return errorhandler.BuildError(c, err, nil)
	}

	return response.Success(c, fiber.StatusCreated, "log created", log)
}

func (h *logHandler) Update(c *fiber.Ctx) error {
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

	var updateLog entity.Log
	if err := c.BodyParser(&updateLog); err != nil {
		return errorhandler.BuildError(c, errorhandler.BadRequestError{Message: err.Error()}, nil)
	}

	log, err := h.usecase.Update(ctx, id, &updateLog)
	if err != nil {
		return errorhandler.BuildError(c, err, nil)
	}

	return response.Success(c, fiber.StatusOK, "log updated", log)
}

func (h *logHandler) Delete(c *fiber.Ctx) error {
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

	return response.Success(c, fiber.StatusOK, "log deleted", nil)
}
