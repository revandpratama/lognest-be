package route

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func InitRoutes(api fiber.Router, db *gorm.DB) {

	InitProjectRoutes(api, db)

	InitLogRoutes(api, db)

	InitTagRoutes(api, db)

	InitUserProfileRoutes(api, db)

	InitInteractionRoutes(api, db)
}
