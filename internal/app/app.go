package app

import (
	"context"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type App struct {
	fiberApp     *fiber.App
	DB           *gorm.DB
	AzblobClient *azblob.Client
}

type Option func(*App) error

func NewApp(opts ...Option) (*App, error) {
	app := &App{}
	for _, opt := range opts {
		if err := opt(app); err != nil {
			return nil, err
		}
	}
	return app, nil
}

func (a *App) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if a.fiberApp != nil {
		if err := a.fiberApp.ShutdownWithContext(ctx); err != nil {
			return err
		}
	}

	if a.DB != nil {
		sqlDb, _ := a.DB.DB()
		if err := sqlDb.Close(); err != nil {
			return err
		}
	}

	log.Info().Msg("resources cleaned up")

	return nil
}
