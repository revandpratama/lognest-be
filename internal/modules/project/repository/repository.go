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
	if err := r.db.WithContext(ctx).Where("id = ?", id).Preload("Tags").Preload("UserProfile").Preload("Logs").First(&project).Error; err != nil {
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

	if err := paginatedDB.Preload("Tags").Find(&projects).Error; err != nil {
		return nil, nil, err
	}
	return projects, paginationQuery, nil
}

func (r *projectRepository) Create(ctx context.Context, newProject *entity.Project) (*entity.Project, error) {

	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.WithContext(ctx).Omit("Tags").Create(newProject).Error; err != nil {
			return err
		}

		if len(newProject.Tags) > 0 {
			for _, tag := range newProject.Tags {
				joinRecord := map[string]interface{}{
					"project_id": newProject.ID,
					"tag_id":     tag.ID,
				}
				// Use .Table() to perform a direct, raw insert.
				if err := tx.Table("lognest.project_tags").Create(&joinRecord).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})

	return newProject, err
}

func (r *projectRepository) Update(ctx context.Context, id uuid.UUID, updateProject *entity.Project) (*entity.Project, error) {

	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.WithContext(ctx).Model(&entity.Project{}).Where("id = ?", id).Updates(updateProject).Error; err != nil {
			return err
		}

		if len(updateProject.Tags) > 0 {
			if err := tx.Exec("DELETE FROM lognest.project_tags WHERE project_id = ?", id).Error; err != nil {
				return err
			}
			if len(updateProject.Tags) > 0 {
				for _, tag := range updateProject.Tags {
					joinRecord := map[string]interface{}{
						"project_id": id,
						"tag_id":     tag.ID,
					}
					if err := tx.Table("lognest.project_tags").Create(&joinRecord).Error; err != nil {
						return err
					}
				}
			}
		}
		return nil
	})

	return updateProject, err
}

func (r *projectRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.Project{}, "id = ?", id).Error
}
