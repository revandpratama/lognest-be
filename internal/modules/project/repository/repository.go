package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/revandpratama/lognest/internal/modules/project/entity"
	"github.com/revandpratama/lognest/pkg/pagination"
	"gorm.io/gorm"
)

type ProjectRepository interface {
	FindBySlug(ctx context.Context, slug string) (*entity.Project, error)
	FindByUserID(ctx context.Context, userID uuid.UUID, paginationQuery *pagination.Pagination) ([]entity.Project, *pagination.Pagination, error)
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Project, error)
	FindAll(ctx context.Context, paginationQuery *pagination.Pagination) ([]entity.Project, *pagination.Pagination, error)
	Create(ctx context.Context, newProject *entity.Project) (*entity.Project, error)
	Update(ctx context.Context, id uuid.UUID, updateProject *entity.Project) (*entity.Project, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type projectRepository struct {
	db *gorm.DB
}

func NewProjectRepository(db *gorm.DB) ProjectRepository {
	return &projectRepository{db: db}
}

func (r *projectRepository) FindBySlug(ctx context.Context, slug string) (*entity.Project, error) {
	var project entity.Project
	if err := r.db.WithContext(ctx).Where("slug = ?", slug).First(&project).Error; err != nil {
		return nil, err
	}
	return &project, nil
}

func (r *projectRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Project, error) {
	var project entity.Project
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&project).Error; err != nil {
		return nil, err
	}
	return &project, nil
}

func (r *projectRepository) FindByUserID(ctx context.Context, userID uuid.UUID, paginationQuery *pagination.Pagination) ([]entity.Project, *pagination.Pagination, error) {
	var projects []entity.Project

	allowedSortColumns := []string{
		"created_at",
		"title",
		"is_public",
	}

	query := r.db.WithContext(ctx).Where("user_id = ?", userID)

	paginatedDB := pagination.Paginate(query, paginationQuery, &projects, allowedSortColumns)

	if err := paginatedDB.Find(&projects).Error; err != nil {
		return nil, nil, err
	}
	return projects, paginationQuery, nil
}

func (r *projectRepository) FindAll(ctx context.Context, paginationQuery *pagination.Pagination) ([]entity.Project, *pagination.Pagination, error) {
	var projects []entity.Project

	allowedSortColumns := []string{
		"created_at",
		"title",
		"is_public",
	}

	paginatedDB := pagination.Paginate(r.db, paginationQuery, &projects, allowedSortColumns)

	if err := paginatedDB.Find(&projects).Error; err != nil {
		return nil, nil, err
	}
	return projects, paginationQuery, nil
}

func (r *projectRepository) Create(ctx context.Context, newProject *entity.Project) (*entity.Project, error) {

	err := r.db.WithContext(ctx).Create(newProject).Error

	return newProject, err
}

func (r *projectRepository) Update(ctx context.Context, id uuid.UUID, updateProject *entity.Project) (*entity.Project, error) {
	err := r.db.WithContext(ctx).Where("id = ?", id).Updates(updateProject).Error

	return updateProject, err
}

func (r *projectRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.Project{}, "id = ?", id).Error
}
