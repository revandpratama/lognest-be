package errorhandler

import "github.com/gofiber/fiber/v2"

func BuildError(c *fiber.Ctx, err error, errors []string) error {
	var statusCode int
	var response ErrorResponse

	switch err.(type) {
	case NotFoundError:
		statusCode = fiber.StatusNotFound
	case BadRequestError:
		statusCode = fiber.StatusBadRequest
	case UnauthorizedError:
		statusCode = fiber.StatusUnauthorized
	case InternalServerError:
		statusCode = fiber.StatusInternalServerError
	case ConflictError:
		statusCode = fiber.StatusConflict
	default:
		statusCode = fiber.StatusInternalServerError
	}

	response = ErrorResponse{
		Status:  "error",
		Message: err.Error(),
		Errors:  errors,
	}

	return c.Status(statusCode).JSON(response)
}

type ErrorResponse struct {
	Status  string   `json:"status"`
	Message string   `json:"message"`
	Errors  []string `json:"errors,omitempty"`
}
