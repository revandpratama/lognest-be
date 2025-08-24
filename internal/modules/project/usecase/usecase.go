package usecase

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/revandpratama/lognest/internal/modules/project/entity"
	"github.com/revandpratama/lognest/internal/modules/project/repository"
	"github.com/revandpratama/lognest/pkg/errorhandler"
	"github.com/revandpratama/lognest/pkg/pagination"
	"github.com/revandpratama/lognest/pkg/slug"
	"gorm.io/gorm"
)

type ProjectUsecase interface {
	FindBySlug(ctx context.Context, slug string) (*entity.Project, error)
	FindByUserID(ctx context.Context, userID uuid.UUID, paginationQuery *pagination.Pagination) ([]entity.Project, *pagination.Pagination, error)
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Project, error)
	FindAll(ctx context.Context, paginationQuery *pagination.Pagination) ([]entity.Project, *pagination.Pagination, error)
	Create(ctx context.Context, newProject *entity.Project) (*entity.Project, error)
	Update(ctx context.Context, updateProject *entity.Project) (*entity.Project, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type projectUsecase struct {
	projectRepository repository.ProjectRepository
}

func NewProjectUsecase(projectRepository repository.ProjectRepository) ProjectUsecase {
	return &projectUsecase{projectRepository: projectRepository}
}

func (p *projectUsecase) FindBySlug(ctx context.Context, slug string) (*entity.Project, error) {

	project, err := p.projectRepository.FindBySlug(ctx, slug)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorhandler.NotFoundError{Message: "project not found"}
		}
		return nil, errorhandler.InternalServerError{Message: err.Error()}
	}

	return project, nil
}

func (p *projectUsecase) FindByUserID(ctx context.Context, userID uuid.UUID, paginationQuery *pagination.Pagination) ([]entity.Project, *pagination.Pagination, error) {
	return p.projectRepository.FindByUserID(ctx, userID, paginationQuery)
}

func (p *projectUsecase) FindByID(ctx context.Context, id uuid.UUID) (*entity.Project, error) {

	project, err := p.projectRepository.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorhandler.NotFoundError{Message: "project not found"}
		}
		return nil, errorhandler.InternalServerError{Message: err.Error()}
	}

	return project, nil
}

func (p *projectUsecase) FindAll(ctx context.Context, paginationQuery *pagination.Pagination) ([]entity.Project, *pagination.Pagination, error) {
	return p.projectRepository.FindAll(ctx, paginationQuery)
}

func (p *projectUsecase) Create(ctx context.Context, newProject *entity.Project) (*entity.Project, error) {

	slug := slug.ToSlug(newProject.Title)

	newProject.Slug = slug

	project, err := p.projectRepository.Create(ctx, newProject)
	if err != nil {
		return nil, err
	}
	return project, nil
}

func (p *projectUsecase) Update(ctx context.Context, updateProject *entity.Project) (*entity.Project, error) {
	return p.projectRepository.Update(ctx, updateProject)
}

func (p *projectUsecase) Delete(ctx context.Context, id uuid.UUID) error {
	return p.projectRepository.Delete(ctx, id)
}
