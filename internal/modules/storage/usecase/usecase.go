package usecase

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/revandpratama/lognest/config"
	"github.com/revandpratama/lognest/internal/modules/storage/dto"
	"github.com/revandpratama/lognest/internal/modules/storage/entity"
	azurestorage "github.com/revandpratama/lognest/pkg/azure-storage"
	"github.com/revandpratama/lognest/pkg/errorhandler"
)

// StorageUsecase defines the business logic interface for a Storage.
type StorageUsecase interface {
	Upload(ctx context.Context, storage *entity.Storage) (string, error)
	GetURL(ctx context.Context, filePath string) (*dto.StorageURL, error)
	Delete(ctx context.Context, filePath string) error
}

type storageUsecase struct {
	azureClient *azblob.Client
}

// NewStorageUsecase creates a new instance of StorageUsecase.
func NewStorageUsecase(azureClient *azblob.Client) StorageUsecase {
	return &storageUsecase{
		azureClient: azureClient,
	}
}


func (u *storageUsecase) Upload(ctx context.Context, storage *entity.Storage) (string, error) {

	if storage == nil {
		return "", nil
	}

	if storage.File != nil && storage.Image != nil {
		return "", errorhandler.BadRequestError{Message: "file and image cannot be uploaded at the same time"}
	}

	var containerName = config.ENV.AZURE_STORAGE_CONTAINER_NAME

	var url string

	if storage.File != nil {
		openedFile, err := storage.File.Open()
		if err != nil {
			return "", errorhandler.InternalServerError{Message: err.Error()}
		}
		defer openedFile.Close()

		satinitzedFileName := azurestorage.SanitizeFileName(storage.File.Filename)
		url, err = azurestorage.UploadFile(ctx, u.azureClient, containerName, storage.PathName, satinitzedFileName, openedFile)
		if err != nil {
			return "", errorhandler.InternalServerError{Message: err.Error()}
		}
	}

	if storage.Image != nil {
		openedImage, err := storage.Image.Open()
		if err != nil {
			return "", errorhandler.InternalServerError{Message: err.Error()}
		}
		defer openedImage.Close()

		satinitzedFileName := azurestorage.SanitizeFileName(storage.Image.Filename)
		url, err = azurestorage.UploadFile(ctx, u.azureClient, containerName, storage.PathName, satinitzedFileName, openedImage)
		if err != nil {
			return "", errorhandler.InternalServerError{Message: err.Error()}
		}
	}

	return url, nil
}

func (u *storageUsecase) GetURL(ctx context.Context, filePath string) (*dto.StorageURL, error) {

	var storage dto.StorageURL

	var containerName = config.ENV.AZURE_STORAGE_CONTAINER_NAME

	url, err := azurestorage.GetFileURL(ctx, u.azureClient, containerName, filePath)
	if err != nil {
		return nil, errorhandler.InternalServerError{Message: err.Error()}
	}

	storage.URL = url

	return &storage, nil
}

func (u *storageUsecase) Delete(ctx context.Context, filePath string) error {

	var containerName = config.ENV.AZURE_STORAGE_CONTAINER_NAME

	if err := azurestorage.DeleteFile(ctx, u.azureClient, containerName, filePath); err != nil {
		return errorhandler.InternalServerError{Message: err.Error()}
	}

	return nil
}
