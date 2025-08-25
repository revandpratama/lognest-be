package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/revandpratama/lognest/internal/modules/tag/entity"
	"github.com/revandpratama/lognest/pkg/pagination"
	"gorm.io/gorm"
)

// TagRepository defines the interface for database operations for a Tag.
type TagRepository interface {
	FindAll(ctx context.Context, paginationQuery *pagination.Pagination) ([]*entity.Tag, *pagination.Pagination, error)
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Tag, error)
	Create(ctx context.Context, newTag *entity.Tag) (*entity.Tag, error)
	Update(ctx context.Context, id uuid.UUID, updateTag *entity.Tag) (*entity.Tag, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type tagRepository struct {
	db *gorm.DB
}

// NewTagRepository creates a new instance of TagRepository.
func NewTagRepository(db *gorm.DB) TagRepository {
	return &tagRepository{db: db}
}

func (r *tagRepository) FindAll(ctx context.Context, paginationQuery *pagination.Pagination) ([]*entity.Tag, *pagination.Pagination, error) {
	var tags []*entity.Tag
	
	allowedSortColumns := []string{
		"created_at",
	}
	
	paginatedDB := pagination.Paginate(r.db.WithContext(ctx), paginationQuery, &tags, allowedSortColumns)
	
	if err := paginatedDB.Find(&tags).Error; err != nil {
		return nil, nil, err
	}
	return tags, paginationQuery, nil
}

func (r *tagRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Tag, error) {
	var tag entity.Tag
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&tag).Error; err != nil {
		return nil, err
	}
	return &tag, nil
}


func (r *tagRepository) Create(ctx context.Context, newTag *entity.Tag) (*entity.Tag, error) {
	err := r.db.WithContext(ctx).Create(newTag).Error
	return newTag, err
}

func (r *tagRepository) Update(ctx context.Context, id uuid.UUID, updateTag *entity.Tag) (*entity.Tag, error) {
	err := r.db.WithContext(ctx).Model(&entity.Tag{}).Where("id = ?", id).Updates(updateTag).Error
	return updateTag, err
}

func (r *tagRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.Tag{}, "id = ?", id).Error
}

