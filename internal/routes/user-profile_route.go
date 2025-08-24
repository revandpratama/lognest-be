package route

import (
	"github.com/gofiber/fiber/v2"
	"github.com/revandpratama/lognest/internal/modules/user-profile/handler"
	"github.com/revandpratama/lognest/internal/modules/user-profile/repository"
	"github.com/revandpratama/lognest/internal/modules/user-profile/usecase"
	"gorm.io/gorm"
)

func initUserProfileHandler(db *gorm.DB) handler.UserProfileHandler {
	userProfileRepo := repository.NewUserProfileRepository(db)
	userProfileUsecase := usecase.NewUserProfileUsecase(userProfileRepo)
	userProfileHandler := handler.NewUserProfileHandler(userProfileUsecase)
	return userProfileHandler
}

func InitUserProfileRoutes(api fiber.Router, db *gorm.DB) {
	userProfileHandler := initUserProfileHandler(db)

	profiles := api.Group("/profiles")

	profiles.Post("/", userProfileHandler.Create)
	profiles.Get("/:id", userProfileHandler.FindByID)
	profiles.Put("/:id", userProfileHandler.Update)
}
