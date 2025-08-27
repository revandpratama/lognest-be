package repository

import (
	"gorm.io/gorm"
)

// AuthRepository defines the interface for database operations for a Auth.
type AuthRepository interface {
	
}

type authRepository struct {
	db *gorm.DB
}

// NewAuthRepository creates a new instance of AuthRepository.
func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &authRepository{db: db}
}

// NOTE: The following are example implementations. You will need to adjust them.
