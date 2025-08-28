package route

import (
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func InitRoutes(api fiber.Router, db *gorm.DB, httpClient *http.Client, azureClient *azblob.Client) {

	InitProjectRoutes(api, db)

	InitLogRoutes(api, db)

	InitTagRoutes(api, db)

	InitUserProfileRoutes(api, db, httpClient)

	InitInteractionRoutes(api, db)

	InitStorageRoute(api, azureClient)

	InitAuthRoute(api, db, httpClient)
}
