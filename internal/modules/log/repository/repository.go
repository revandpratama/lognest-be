package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/revandpratama/lognest/internal/modules/log/entity"
	"github.com/revandpratama/lognest/pkg/pagination"
	"gorm.io/gorm"
)

// LogRepository defines the interface for database operations for a Log.
type LogRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Log, error)
	FindByProjectID(ctx context.Context, projectID uuid.UUID, paginationQuery *pagination.Pagination) ([]entity.Log, *pagination.Pagination, error)
	Create(ctx context.Context, newLog *entity.Log) (*entity.Log, error)
	Update(ctx context.Context, id uuid.UUID, updateLog *entity.Log) (*entity.Log, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type logRepository struct {
	db *gorm.DB
}

// NewLogRepository creates a new instance of LogRepository.
func NewLogRepository(db *gorm.DB) LogRepository {
	return &logRepository{db: db}
}

// NOTE: The following are example implementations. You will need to adjust them.

func (r *logRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Log, error) {
	var log entity.Log
	if err := r.db.WithContext(ctx).Where("id = ?", id).Preload("Comments").Preload("Media").First(&log).Error; err != nil {
		return nil, err
	}
	return &log, nil
}

func (r *logRepository) FindByProjectID(ctx context.Context, projectID uuid.UUID, paginationQuery *pagination.Pagination) ([]entity.Log, *pagination.Pagination, error) {
	var logs []entity.Log

	allowedSortColumns := []string{
		"created_at",
		"title",
		"is_public",
	}

	query := r.db.WithContext(ctx).Where("project_id = ?", projectID)

	paginatedDB := pagination.Paginate(query, paginationQuery, &logs, allowedSortColumns)

	if err := paginatedDB.Preload("Comments").Preload("Media").Find(&logs).Error; err != nil {
		return nil, nil, err
	}
	return logs, paginationQuery, nil
}

func (r *logRepository) Create(ctx context.Context, newLog *entity.Log) (*entity.Log, error) {
	err := r.db.WithContext(ctx).Create(newLog).Error
	return newLog, err
}

func (r *logRepository) Update(ctx context.Context, id uuid.UUID, updateLog *entity.Log) (*entity.Log, error) {
	err := r.db.WithContext(ctx).Model(&entity.Log{}).Where("id = ?", id).Updates(updateLog).Error
	return updateLog, err
}

func (r *logRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.Log{}, "id = ?", id).Error
}
