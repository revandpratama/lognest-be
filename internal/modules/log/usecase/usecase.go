package usecase

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/revandpratama/lognest/internal/modules/log/entity"
	"github.com/revandpratama/lognest/internal/modules/log/repository"
	"github.com/revandpratama/lognest/pkg/pagination"
	"gorm.io/gorm"
)

// LogUsecase defines the business logic interface for a Log.
type LogUsecase interface {
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Log, error)
	FindByProjectID(ctx context.Context, projectID uuid.UUID, paginationQuery *pagination.Pagination) ([]entity.Log, *pagination.Pagination, error)
	Create(ctx context.Context, newLog *entity.Log) (*entity.Log, error)
	Update(ctx context.Context, id uuid.UUID, updateLog *entity.Log) (*entity.Log, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type logUsecase struct {
	repo repository.LogRepository
}

// NewLogUsecase creates a new instance of LogUsecase.
func NewLogUsecase(repo repository.LogRepository) LogUsecase {
	return &logUsecase{repo: repo}
}

func (u *logUsecase) FindByID(ctx context.Context, id uuid.UUID) (*entity.Log, error) {
	log, err := u.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("log not found")
		}
		return nil, err
	}

	return log, nil
}

func (u *logUsecase) FindByProjectID(ctx context.Context, projectID uuid.UUID, paginationQuery *pagination.Pagination) ([]entity.Log, *pagination.Pagination, error) {
	return u.repo.FindByProjectID(ctx, projectID, paginationQuery)
}

func (u *logUsecase) Create(ctx context.Context, newLog *entity.Log) (*entity.Log, error) {
	return u.repo.Create(ctx, newLog)
}

func (u *logUsecase) Update(ctx context.Context, id uuid.UUID, updateLog *entity.Log) (*entity.Log, error) {
	return u.repo.Update(ctx, id, updateLog)
}

func (u *logUsecase) Delete(ctx context.Context, id uuid.UUID) error {
	return u.repo.Delete(ctx, id)
}
