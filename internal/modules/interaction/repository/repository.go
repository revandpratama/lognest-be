package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/revandpratama/lognest/internal/modules/interaction/entity"
	"gorm.io/gorm"
)

// InteractionRepository defines the interface for database operations for a Interaction.
type InteractionRepository interface {
	CreateLike(ctx context.Context, newLike *entity.Like) (*entity.Like, error)
	DeleteLike(ctx context.Context, userProfileID uuid.UUID, logID uuid.UUID) error
	FindLikeByLogID(ctx context.Context, logID uuid.UUID) (*[]entity.Like, error)
	CreateComment(ctx context.Context, newComment *entity.Comment) (*entity.Comment, error)
	UpdateComment(ctx context.Context, id uuid.UUID, updateComment *entity.Comment) (*entity.Comment, error)
	DeleteComment(ctx context.Context, commentID uuid.UUID) error
	FindCommentByLogID(ctx context.Context, logID uuid.UUID) ([]entity.Comment, error)
}

type interactionRepository struct {
	db *gorm.DB
}

// NewInteractionRepository creates a new instance of InteractionRepository.
func NewInteractionRepository(db *gorm.DB) InteractionRepository {
	return &interactionRepository{db: db}
}

func (r *interactionRepository) CreateLike(ctx context.Context, newLike *entity.Like) (*entity.Like, error) {
	err := r.db.WithContext(ctx).Create(newLike).Error
	return newLike, err
}

func (r *interactionRepository) DeleteLike(ctx context.Context, userProfileID uuid.UUID, logID uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.Like{}, "user_profile_id = ? AND log_id = ?", userProfileID, logID).Error
}

func (r *interactionRepository) FindLikeByLogID(ctx context.Context, logID uuid.UUID) (*[]entity.Like, error) {
	var likes []entity.Like
	if err := r.db.WithContext(ctx).Where("log_id = ?", logID).Find(&likes).Error; err != nil {
		return nil, err
	}
	return &likes, nil
}

func (r *interactionRepository) CreateComment(ctx context.Context, newComment *entity.Comment) (*entity.Comment, error) {
	err := r.db.WithContext(ctx).Create(newComment).Error
	return newComment, err
}

func (r *interactionRepository) UpdateComment(ctx context.Context, id uuid.UUID, updateComment *entity.Comment) (*entity.Comment, error) {
	err := r.db.WithContext(ctx).Model(&entity.Comment{}).Where("id = ?", id).Updates(updateComment).Error
	return updateComment, err
}

func (r *interactionRepository) DeleteComment(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.Comment{}, "id = ?", id).Error
}

func (r *interactionRepository) FindCommentByLogID(ctx context.Context, logID uuid.UUID) ([]entity.Comment, error) {
	var comments []entity.Comment
	if err := r.db.WithContext(ctx).Where("log_id = ?", logID).Find(&comments).Error; err != nil {
		return nil, err
	}
	return comments, nil
}
