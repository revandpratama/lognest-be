package route

import (
	"github.com/gofiber/fiber/v2"
	"github.com/revandpratama/lognest/internal/middlewares"
	"github.com/revandpratama/lognest/internal/modules/interaction/handler"
	"github.com/revandpratama/lognest/internal/modules/interaction/repository"
	"github.com/revandpratama/lognest/internal/modules/interaction/usecase"
	"gorm.io/gorm"
)

func initInteractionHandler(db *gorm.DB) handler.InteractionHandler {
	interactionRepo := repository.NewInteractionRepository(db)
	interactionUsecase := usecase.NewInteractionUsecase(interactionRepo)
	interactionHandler := handler.NewInteractionHandler(interactionUsecase)
	return interactionHandler
}

func InitInteractionRoutes(api fiber.Router, db *gorm.DB) {
	interactionHandler := initInteractionHandler(db)

	interaction := api.Group("/interactions")

	interaction.Use(middlewares.AuthMiddleware())

	interaction.Post("/likes", interactionHandler.CreateLike)
	interaction.Delete("/likes/:likeID", interactionHandler.DeleteLike)
	interaction.Get("/likes/logs/:logID", interactionHandler.FindLikeByLogID)
	interaction.Post("/comments", interactionHandler.CreateComment)
	interaction.Put("/comments/:commentID", interactionHandler.UpdateComment)
	interaction.Get("/comments/log/:logID", interactionHandler.FindCommentByLogID)
	interaction.Delete("/comments/:commentID", interactionHandler.DeleteComment)
}
