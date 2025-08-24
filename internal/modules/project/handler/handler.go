package handler

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/revandpratama/lognest/internal/modules/project/entity"
	"github.com/revandpratama/lognest/internal/modules/project/repository"
	"github.com/revandpratama/lognest/pkg/errorhandler"
	"github.com/revandpratama/lognest/pkg/pagination"
	"github.com/revandpratama/lognest/pkg/response"
)

type ProjectHandler interface {
	FindBySlug(c *fiber.Ctx) error
	FindByUserID(c *fiber.Ctx) error
	FindAll(c *fiber.Ctx) error
	FindByID(c *fiber.Ctx) error
	Create(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
}

type projectHandler struct {
	ProjectRepository repository.ProjectRepository
}

func NewProjectHandler(projectRepository repository.ProjectRepository) ProjectHandler {
	return &projectHandler{ProjectRepository: projectRepository}
}

func (h *projectHandler) FindBySlug(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
	defer cancel()

	slug := c.Params("slug")
	if slug == "" {
		return errorhandler.BuildError(c, errorhandler.BadRequestError{Message: "slug is required"}, nil)
	}

	project, err := h.ProjectRepository.FindBySlug(ctx, slug)
	if err != nil {
		return errorhandler.BuildError(c, err, nil)
	}

	return response.Success(c, fiber.StatusOK, "project found", project)
}

func (h *projectHandler) FindByUserID(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
	defer cancel()

	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return errorhandler.BuildError(c, errorhandler.UnauthorizedError{Message: "unauthorized: user_id not found"}, nil)
	}

	paginationQuery := new(pagination.Pagination)
	if err := c.QueryParser(paginationQuery); err != nil {
		return errorhandler.BuildError(c, errorhandler.BadRequestError{Message: err.Error()}, nil)
	}

	projects, pagination, err := h.ProjectRepository.FindByUserID(ctx, userID, paginationQuery)
	if err != nil {
		return errorhandler.BuildError(c, err, nil)
	}

	return response.Paginated(c, fiber.StatusOK, "projects found", projects, pagination)
}

func (h *projectHandler) FindByID(c *fiber.Ctx) error {
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

	project, err := h.ProjectRepository.FindByID(ctx, id)
	if err != nil {
		return errorhandler.BuildError(c, err, nil)
	}

	return response.Success(c, fiber.StatusOK, "project found", project)
}

func (h *projectHandler) FindAll(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
	defer cancel()

	paginationQuery := new(pagination.Pagination)
	if err := c.QueryParser(paginationQuery); err != nil {
		return errorhandler.BuildError(c, errorhandler.BadRequestError{Message: err.Error()}, nil)
	}

	projects, pagination, err := h.ProjectRepository.FindAll(ctx, paginationQuery)
	if err != nil {
		return errorhandler.BuildError(c, err, nil)
	}

	return response.Paginated(c, fiber.StatusOK, "projects found", projects, pagination)
}

func (h *projectHandler) Create(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
	defer cancel()

	var newProject entity.Project

	if err := c.BodyParser(&newProject); err != nil {
		return errorhandler.BuildError(c, errorhandler.BadRequestError{Message: err.Error()}, nil)
	}

	project, err := h.ProjectRepository.Create(ctx, &newProject)
	if err != nil {
		return errorhandler.BuildError(c, err, nil)
	}

	return response.Success(c, fiber.StatusCreated, "project created", project)

}
func (h *projectHandler) Update(c *fiber.Ctx) error {
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

	var updateProject entity.Project

	if err := c.BodyParser(&updateProject); err != nil {
		return errorhandler.BuildError(c, errorhandler.BadRequestError{Message: err.Error()}, nil)
	}

	project, err := h.ProjectRepository.Update(ctx, id, &updateProject)
	if err != nil {
		return errorhandler.BuildError(c, err, nil)
	}

	return response.Success(c, fiber.StatusOK, "project updated", project)
}

func (h *projectHandler) Delete(c *fiber.Ctx) error {
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

	err = h.ProjectRepository.Delete(ctx, id)
	if err != nil {
		return errorhandler.BuildError(c, err, nil)
	}

	return response.Success(c, fiber.StatusOK, "project deleted", nil)
}
