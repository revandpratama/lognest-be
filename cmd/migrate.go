package cmd

import (
	"fmt"

	projectEntity "github.com/revandpratama/lognest/internal/modules/project/entity"
	logEntity "github.com/revandpratama/lognest/internal/modules/log/entity"
	"gorm.io/gorm"
)

func EnsureSchema(db *gorm.DB, schema string) error {
	return db.Exec(fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s", schema)).Error
}

func MigrateDatabase(db *gorm.DB) error {

	var models = []interface{}{
		&projectEntity.Project{},
		&logEntity.Log{},
	}

	return db.AutoMigrate(models...)
}
