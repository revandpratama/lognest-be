package handler

import (
	"context"
	"net/url"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/revandpratama/lognest/internal/modules/storage/entity"
	"github.com/revandpratama/lognest/internal/modules/storage/usecase"
	"github.com/revandpratama/lognest/pkg/errorhandler"
	"github.com/revandpratama/lognest/pkg/response"
)

// StorageHandler defines the HTTP handler interface for a Storage.
type StorageHandler interface {
	Upload(c *fiber.Ctx) error
	GetURL(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
}

type storageHandler struct {
	usecase usecase.StorageUsecase
}

// NewStorageHandler creates a new instance of StorageHandler.
func NewStorageHandler(usecase usecase.StorageUsecase) StorageHandler {
	return &storageHandler{usecase: usecase}
}

func (h *storageHandler) Upload(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var storage entity.Storage
	pathName := c.FormValue("path_name")
	if pathName == "" {
		return errorhandler.BuildError(c, errorhandler.BadRequestError{Message: "path_name is required"}, nil)
	}

	storage.PathName = pathName

	file, _ := c.FormFile("file")
	if file != nil {
		storage.File = file
	}

	image, _ := c.FormFile("image")
	if image != nil {
		storage.Image = image
	}

	res, err := h.usecase.Upload(ctx, &storage)
	if err != nil {
		return errorhandler.BuildError(c, err, nil)
	}

	return response.Success(c, fiber.StatusOK, "file uploaded", res)
}

func (h *storageHandler) GetURL(c *fiber.Ctx) error {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	encodedFilePath := c.Params("filePath")
	if encodedFilePath == "" {
		return errorhandler.BuildError(c, errorhandler.BadRequestError{Message: "filePath is required"}, nil)
	}

	// 2. Manually URL-decode the string.
	// This will convert "images%2Ffile.jpeg" into "images/file.jpeg"
	filePath, err := url.PathUnescape(encodedFilePath)
	if err != nil {
		// This happens if the URL has malformed encoding, e.g., "%2"
		return errorhandler.BuildError(c, errorhandler.BadRequestError{Message: "invalid file path encoding"}, nil)
	}

	res, err := h.usecase.GetURL(ctx, filePath)
	if err != nil {
		return errorhandler.BuildError(c, err, nil)
	}

	return response.Success(c, fiber.StatusOK, "get url success", res)
}

func (h *storageHandler) Delete(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filePath := c.Params("filePath")

	err := h.usecase.Delete(ctx, filePath)
	if err != nil {
		return errorhandler.BuildError(c, err, nil)
	}

	return response.Success(c, fiber.StatusOK, "file deleted", nil)
}
