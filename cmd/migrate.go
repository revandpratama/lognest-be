package cmd

import (
	"fmt"

	logEntity "github.com/revandpratama/lognest/internal/modules/log/entity"
	projectEntity "github.com/revandpratama/lognest/internal/modules/project/entity"
	tagEntity "github.com/revandpratama/lognest/internal/modules/tag/entity"
	userProfileEntity "github.com/revandpratama/lognest/internal/modules/user-profile/entity"
	"gorm.io/gorm"
)

func EnsureSchema(db *gorm.DB, schema string) error {
	return db.Exec(fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s", schema)).Error
}

var models = []interface{}{
	&projectEntity.Project{},
	&tagEntity.Tag{},
	&logEntity.Log{},
	&userProfileEntity.UserProfile{},
}

func MigrateDatabase(db *gorm.DB) error {

	if err := db.AutoMigrate(models...); err != nil {
		return err
	}

	return nil
}

func MigrateDatabaseFresh(db *gorm.DB) error {

	var conjunctionTables = []interface{}{
		&projectEntity.ProjectTag{},
	}

	models = append(models, conjunctionTables...)

	if err := db.Migrator().DropTable(models...); err != nil {
		return err
	}

	return MigrateDatabase(db)
}
