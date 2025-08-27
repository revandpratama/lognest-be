package route

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func InitRoutes(api fiber.Router, db *gorm.DB, httpClient *http.Client) {

	InitProjectRoutes(api, db)

	InitLogRoutes(api, db)

	InitTagRoutes(api, db)

	InitUserProfileRoutes(api, db)

	InitInteractionRoutes(api, db)

	InitAuthRoute(api, db, httpClient)
}
