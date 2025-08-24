package cmd

import (
	"fmt"

	logEntity "github.com/revandpratama/lognest/internal/modules/log/entity"
	projectEntity "github.com/revandpratama/lognest/internal/modules/project/entity"
	userProfileEntity "github.com/revandpratama/lognest/internal/modules/user-profile/entity"
	"gorm.io/gorm"
)

func EnsureSchema(db *gorm.DB, schema string) error {
	return db.Exec(fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s", schema)).Error
}

func MigrateDatabase(db *gorm.DB) error {

	var models = []interface{}{
		&projectEntity.Project{},
		&logEntity.Log{},
		&userProfileEntity.UserProfile{},
	}

	return db.AutoMigrate(models...)
}
