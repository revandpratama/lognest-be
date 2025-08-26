package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/revandpratama/lognest/internal/modules/interaction/entity"
	"github.com/revandpratama/lognest/internal/modules/interaction/repository"
)

// InteractionUsecase defines the business logic interface for a Interaction.
type InteractionUsecase interface {
	CreateLike(ctx context.Context, newLike *entity.Like) (*entity.Like, error)
	DeleteLike(ctx context.Context, userProfileID uuid.UUID, logID uuid.UUID) error
	FindLikeByLogID(ctx context.Context, logID uuid.UUID) (*[]entity.Like, error)
	CreateComment(ctx context.Context, newComment *entity.Comment) (*entity.Comment, error)
	UpdateComment(ctx context.Context, id uuid.UUID, updateComment *entity.Comment) (*entity.Comment, error)
	DeleteComment(ctx context.Context, commentID uuid.UUID) error
	FindCommentByLogID(ctx context.Context, logID uuid.UUID) ([]entity.Comment, error)
}

type interactionUsecase struct {
	repo repository.InteractionRepository
}

// NewInteractionUsecase creates a new instance of InteractionUsecase.
func NewInteractionUsecase(repo repository.InteractionRepository) InteractionUsecase {
	return &interactionUsecase{repo: repo}
}

func (u *interactionUsecase) CreateLike(ctx context.Context, newLike *entity.Like) (*entity.Like, error) {
	return u.repo.CreateLike(ctx, newLike)
}

func (u *interactionUsecase) DeleteLike(ctx context.Context, userProfileID uuid.UUID, logID uuid.UUID) error {
	return u.repo.DeleteLike(ctx, userProfileID, logID)
}

func (u *interactionUsecase) FindLikeByLogID(ctx context.Context, logID uuid.UUID) (*[]entity.Like, error) {
	return u.repo.FindLikeByLogID(ctx, logID)
}

func (u *interactionUsecase) CreateComment(ctx context.Context, newComment *entity.Comment) (*entity.Comment, error) {
	return u.repo.CreateComment(ctx, newComment)
}

func (u *interactionUsecase) UpdateComment(ctx context.Context, id uuid.UUID, updateComment *entity.Comment) (*entity.Comment, error) {
	return u.repo.UpdateComment(ctx, id, updateComment)
}

func (u *interactionUsecase) DeleteComment(ctx context.Context, commentID uuid.UUID) error {
	return u.repo.DeleteComment(ctx, commentID)
}

func (u *interactionUsecase) FindCommentByLogID(ctx context.Context, logID uuid.UUID) ([]entity.Comment, error) {
	return u.repo.FindCommentByLogID(ctx, logID)
}
