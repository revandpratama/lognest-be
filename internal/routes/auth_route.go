package route

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/revandpratama/lognest/internal/modules/auth/handler"
	"github.com/revandpratama/lognest/internal/modules/auth/repository"
	"github.com/revandpratama/lognest/internal/modules/auth/usecase"
	"gorm.io/gorm"
)

func initAuthHandler(db *gorm.DB, httpClient *http.Client) handler.AuthHandler {

	authRepo := repository.NewAuthRepository(db)
	authUsecase := usecase.NewAuthUsecase(authRepo, httpClient)
	authHandler := handler.NewAuthHandler(authUsecase)

	return authHandler
}

func InitAuthRoute(api fiber.Router, db *gorm.DB, httpClient *http.Client) {
	authHandler := initAuthHandler(db, httpClient)

	auth := api.Group("/auth")

	auth.Post("/login", authHandler.Login)
	auth.Post("/register", authHandler.Register)
	auth.Post("/refresh-token", authHandler.RefreshToken)
	auth.Post("/logout", authHandler.Logout)
}
