package repository

import (
	"gorm.io/gorm"
)

// StorageRepository defines the interface for database operations for a Storage.
type StorageRepository interface {
	
}

type storageRepository struct {
	db *gorm.DB
}

// NewStorageRepository creates a new instance of StorageRepository.
func NewStorageRepository(db *gorm.DB) StorageRepository {
	return &storageRepository{db: db}
}

// NOTE: The following are example implementations. You will need to adjust them.
