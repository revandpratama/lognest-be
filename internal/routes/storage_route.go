package route

import (
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/gofiber/fiber/v2"
	"github.com/revandpratama/lognest/internal/modules/storage/handler"
	"github.com/revandpratama/lognest/internal/modules/storage/usecase"
)

func initStorageHandler(azureClient *azblob.Client) handler.StorageHandler {

	storageUsecase := usecase.NewStorageUsecase(azureClient)
	return handler.NewStorageHandler(storageUsecase)

}

func InitStorageRoute(api fiber.Router, azureClient *azblob.Client) {

	storageHandler := initStorageHandler(azureClient)

	storage := api.Group("/storage")
	storage.Get("/url/:filePath", storageHandler.GetURL)
	storage.Post("/upload", storageHandler.Upload)
	storage.Delete("/delete/:filePath", storageHandler.Delete)
}
