package app

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/encryptcookie"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/revandpratama/lognest/config"
	route "github.com/revandpratama/lognest/internal/routes"

	// "github.com/revandpratama/lognest/internal/routes"
	"github.com/rs/zerolog/log"
)

func WithRESTServer() Option {
	return func(app *App) error {

		fiberApp := fiber.New(fiber.Config{
			DisableStartupMessage: true,
		})

		fiberApp.Use(func(c *fiber.Ctx) error {
			c.Set("Content-Type", "application/json")
			return c.Next()
		})

		fiberApp.Use(limiter.New(limiter.Config{
			Max:        50,
			Expiration: 5 * time.Second,
			KeyGenerator: func(c *fiber.Ctx) string {
				return c.IP()
			},
			LimitReached: func(c *fiber.Ctx) error {
				return c.Status(429).SendString("Too Many Requests")
			},
		}))

		fiberApp.Use(cors.New(cors.Config{
			AllowOrigins:     config.ENV.CORS_ALLOWED_ORIGINS,
			// AllowCredentials: true,
			AllowHeaders:     "Origin, Content-Type, Accept, Content-Length, Authorization, Accept-Encoding, X-CSRF-Token, X-Requested-With, X-Refresh-Token",
		}))

		fiberApp.Use(encryptcookie.New(encryptcookie.Config{
			Key: config.ENV.COOKIE_SECRET,
		}))

		fiberApp.Use(logger.New(logger.Config{
			Format:        "${time} | ${status} | ${latency} | ${ip} | ${method} | ${path} | ${error}\n",
			TimeFormat:    "15:04:05",
			TimeZone:      "Local",
			TimeInterval:  500 * time.Millisecond,
			Output:        os.Stdout,
			DisableColors: false,
		}))

		fiberApp.Get("/hello", func(c *fiber.Ctx) error {
			return c.SendString("Hello, World!")
		})

		api := fiberApp.Group("/api")

		api.Get("/test-700ms", func(c *fiber.Ctx) error {
			time.Sleep(500 * time.Millisecond)
			return c.SendString("Hello. 700ms delay!")
		})

		httpClient := &http.Client{}

		// * Initialize routes
		route.InitRoutes(api, app.DB, httpClient)

		app.fiberApp = fiberApp

		go func() {
			if err := fiberApp.Listen(fmt.Sprintf(":%s", config.ENV.REST_PORT)); err != nil {
				log.Error().Err(err).Msg("failed to start server")
				os.Exit(1)
			}
		}()

		log.Info().Msgf("REST server started on PORT %s", config.ENV.REST_PORT)

		return nil
	}
}
