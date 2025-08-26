package middlewares

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/revandpratama/lognest/pkg/errorhandler"
	"github.com/revandpratama/lognest/pkg/token"
)

func AuthMiddleware() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {

		authHeader := c.Get("Authorization")

		if c.Get("Authorization") == "" {
			return errorhandler.BuildError(c, errorhandler.UnauthorizedError{Message: "unauthorized, no token provided"}, nil)
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return errorhandler.BuildError(c, errorhandler.UnauthorizedError{Message: "unauthorized, invalid token format"}, nil)
		}

		encryptedToken := parts[1]
		user, err := token.ValidateToken(encryptedToken)
		if err != nil {
			return errorhandler.BuildError(c, errorhandler.UnauthorizedError{Message: "unauthorized, invalid token"}, nil)
		}

		c.Locals("userID", user.UserID)
		c.Locals("provider", user.Provider)
		c.Locals("email", user.Email)
		c.Locals("roleID", user.RoleID)
		c.Locals("sessionID", user.SessionID)
		c.Locals("mfaCompleted", user.MFACompleted)

		return c.Next()
	}
}
