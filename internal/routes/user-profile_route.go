package route

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/revandpratama/lognest/internal/middlewares"
	"github.com/revandpratama/lognest/internal/modules/user-profile/handler"
	"github.com/revandpratama/lognest/internal/modules/user-profile/repository"
	"github.com/revandpratama/lognest/internal/modules/user-profile/usecase"
	"gorm.io/gorm"
)

func initUserProfileHandler(db *gorm.DB, httpClient *http.Client) handler.UserProfileHandler {
	userProfileRepo := repository.NewUserProfileRepository(db)
	userProfileUsecase := usecase.NewUserProfileUsecase(userProfileRepo, httpClient)
	userProfileHandler := handler.NewUserProfileHandler(userProfileUsecase)
	return userProfileHandler
}

func InitUserProfileRoutes(api fiber.Router, db *gorm.DB, httpClient *http.Client) {
	userProfileHandler := initUserProfileHandler(db, httpClient)

	profiles := api.Group("/profiles")

	profiles.Use(middlewares.AuthMiddleware())

	// profiles.Post("/", userProfileHandler.Create)
	// profiles.Get("/:id", userProfileHandler.FindByID)
	profiles.Get("/me", userProfileHandler.FindUser)
	// profiles.Put("/:id", userProfileHandler.Update)
}
