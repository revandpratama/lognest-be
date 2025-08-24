package route

import (
	"github.com/gofiber/fiber/v2"
	"github.com/revandpratama/lognest/internal/modules/log/handler"
	"github.com/revandpratama/lognest/internal/modules/log/repository"
	"github.com/revandpratama/lognest/internal/modules/log/usecase"
	"gorm.io/gorm"
)

func InitLogHandlers(db *gorm.DB) handler.LogHandler {
	logRepository := repository.NewLogRepository(db)
	logUsecase := usecase.NewLogUsecase(logRepository)
	logHandler := handler.NewLogHandler(logUsecase)

	return logHandler
}

func InitLogRoutes(api fiber.Router, db *gorm.DB) {
	logHandler := InitLogHandlers(db)

	log := api.Group("/logs")

	log.Get("/projects/:projectID", logHandler.FindByProjectID)
	log.Get("/:id", logHandler.FindByID)
	log.Post("/", logHandler.Create)
	log.Put("/:id", logHandler.Update)
	log.Delete("/:id", logHandler.Delete)
}
