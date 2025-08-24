package response

import (
	"github.com/gofiber/fiber/v2"
	"github.com/revandpratama/lognest/pkg/pagination"
)

type APIResponse struct {
	Status     string                `json:"status"`
	Message    string                `json:"message"`
	Data       any                   `json:"data,omitempty"`
	Error      any                   `json:"error,omitempty"`
	Pagination pagination.Pagination `json:"pagination,omitempty"`
}

type PaginatedResponse struct {
	APIResponse
	Pagination *pagination.Pagination `json:"pagination,omitempty"`
}

type ResponseParam struct {
	StatusCode int                   `json:"status_code"`
	Message    string                `json:"message"`
	Data       any                   `json:"data"`
	Pagination pagination.Pagination `json:"pagination"`
}

func Success(c *fiber.Ctx, statusCode int, message string, data any) error {
	return c.Status(statusCode).JSON(APIResponse{
		Status:  "success",
		Message: message,
		Data:    data,
	})
}

func Paginated(c *fiber.Ctx, statusCode int, message string, data any, pagination *pagination.Pagination) error {
	return c.Status(statusCode).JSON(
		PaginatedResponse{
			APIResponse: APIResponse{
				Status:  "success",
				Message: message,
				Data:    data,
			},
			Pagination: pagination,
		},
	)
}
