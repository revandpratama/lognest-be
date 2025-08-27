package route

import (
	"github.com/gofiber/fiber/v2"
	"github.com/revandpratama/lognest/internal/middlewares"
	"github.com/revandpratama/lognest/internal/modules/tag/handler"
	"github.com/revandpratama/lognest/internal/modules/tag/repository"
	"github.com/revandpratama/lognest/internal/modules/tag/usecase"
	"gorm.io/gorm"
)

func initTagHandler(db *gorm.DB) handler.TagHandler {
	tagRepo := repository.NewTagRepository(db)
	tagUsecase := usecase.NewTagUsecase(tagRepo)
	tagHandler := handler.NewTagHandler(tagUsecase)

	return tagHandler
}

func InitTagRoutes(api fiber.Router, db *gorm.DB) {
	tagHandler := initTagHandler(db)

	tags := api.Group("/tags")

	tags.Use(middlewares.AuthMiddleware())

	tags.Get("/", tagHandler.FindAll)
	tags.Get("/:id", tagHandler.FindByID)
	tags.Post("/", tagHandler.Create)
	tags.Put("/:id", tagHandler.Update)
	tags.Delete("/:id", tagHandler.Delete)
}
