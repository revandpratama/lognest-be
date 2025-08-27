package route

import (
	"github.com/gofiber/fiber/v2"
	"github.com/revandpratama/lognest/internal/middlewares"
	"github.com/revandpratama/lognest/internal/modules/project/handler"
	"github.com/revandpratama/lognest/internal/modules/project/repository"
	"github.com/revandpratama/lognest/internal/modules/project/usecase"
	"gorm.io/gorm"
)

func initProjectHandler(db *gorm.DB) handler.ProjectHandler {
	projectRepo := repository.NewProjectRepository(db)
	projectusecase := usecase.NewProjectUsecase(projectRepo)
	projectHandler := handler.NewProjectHandler(projectusecase)
	return projectHandler
}

func InitProjectRoutes(api fiber.Router, db *gorm.DB) {
	projectHandler := initProjectHandler(db)

	projects := api.Group("/projects")

	projects.Use(middlewares.AuthMiddleware())

	projects.Get("/", projectHandler.FindAll)
	projects.Get("/:id", projectHandler.FindByID)
	projects.Get("/users/:userID", projectHandler.FindByPublicUserID)
	projects.Get("/me", projectHandler.FindByUserID)
	projects.Get("/slug/:slug", projectHandler.FindBySlug)
	projects.Post("/", projectHandler.Create)
	projects.Put("/:id", projectHandler.Update)
	projects.Delete("/:id", projectHandler.Delete)
}
