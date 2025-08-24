package app

import (
	"fmt"
	"time"

	"github.com/revandpratama/lognest/config"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func WithDB() Option {
	return func(app *App) error {

		dataSourceName := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable timezone=Asia/Jakarta",
			config.ENV.DB_HOST,
			config.ENV.DB_USER,
			config.ENV.DB_PASSWORD,
			config.ENV.DB_NAME,
			config.ENV.DB_PORT,
		)

		db, err := gorm.Open(postgres.New(postgres.Config{
			DSN: dataSourceName,
			PreferSimpleProtocol: true,
		}), &gorm.Config{})
		if err != nil {
			return err
		}

		sqlDB, err := db.DB()
		if err != nil {
			return fmt.Errorf("failed to get raw DB from GORM: %w", err)
		}

		sqlDB.SetMaxIdleConns(10)                  // number of idle connections
		sqlDB.SetMaxOpenConns(100)                 // total connections
		sqlDB.SetConnMaxLifetime(30 * time.Minute) // reuse a connection for 30 minutes
		sqlDB.SetConnMaxIdleTime(5 * time.Minute)  // keep idle for 5 minutes

		if err := sqlDB.Ping(); err != nil {
			return fmt.Errorf("database ping failed: %w", err)
		}

		app.DB = db

		log.Info().Msg("database connected")

		return nil
	}
}
