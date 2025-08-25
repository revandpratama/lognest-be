package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/revandpratama/lognest/internal/modules/tag/entity"
	"github.com/revandpratama/lognest/internal/modules/tag/repository"
	"github.com/revandpratama/lognest/pkg/pagination"
)

// TagUsecase defines the business logic interface for a Tag.
type TagUsecase interface {
	FindAll(ctx context.Context, paginationQuery *pagination.Pagination) ([]*entity.Tag, *pagination.Pagination, error)
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Tag, error)
	Create(ctx context.Context, newTag *entity.Tag) (*entity.Tag, error)
	Update(ctx context.Context, id uuid.UUID, updateTag *entity.Tag) (*entity.Tag, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type tagUsecase struct {
	repo repository.TagRepository
}

// NewTagUsecase creates a new instance of TagUsecase.
func NewTagUsecase(repo repository.TagRepository) TagUsecase {
	return &tagUsecase{repo: repo}
}

func (u *tagUsecase) FindAll(ctx context.Context, paginationQuery *pagination.Pagination) ([]*entity.Tag, *pagination.Pagination, error) {
	return u.repo.FindAll(ctx, paginationQuery)
}

func (u *tagUsecase) FindByID(ctx context.Context, id uuid.UUID) (*entity.Tag, error) {
	return u.repo.FindByID(ctx, id)
}

func (u *tagUsecase) Create(ctx context.Context, newTag *entity.Tag) (*entity.Tag, error) {
	return u.repo.Create(ctx, newTag)
}

func (u *tagUsecase) Update(ctx context.Context, id uuid.UUID, updateTag *entity.Tag) (*entity.Tag, error) {
	return u.repo.Update(ctx, id, updateTag)
}

func (u *tagUsecase) Delete(ctx context.Context, id uuid.UUID) error {
	return u.repo.Delete(ctx, id)
}
