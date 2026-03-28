package config

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/helmet/v2"
	"github.com/spf13/viper"
)

// NewFiber function untuk create fiber app yang digunakan untuk web server
func NewFiber(config *viper.Viper) *fiber.App {
	var app = fiber.New(fiber.Config{
		AppName:      config.GetString("app.name"),
		Prefork:      config.GetBool("app.prefork"),
		ErrorHandler: NewErrorHandler(),
		BodyLimit:    256 * 1024, // 256 KB
	})

	app.Use(helmet.New())

	// CORS: restrict to configured origins; AllowCredentials required for HTTP-only cookies
	allowedOrigins := config.GetString("ALLOWED_ORIGINS")
	if allowedOrigins == "" {
		allowedOrigins = "http://localhost:3000"
	}
	app.Use(cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     "GET,POST,PUT,PATCH,DELETE,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept",
		AllowCredentials: true,
	}))

	return app
}

// NewErrorHandler function untuk membuat custom error handler di fiber
func NewErrorHandler() fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {
		code := fiber.StatusInternalServerError
		message := "Internal Server Error"

		var e *fiber.Error
		if errors.As(err, &e) {
			code = e.Code
			message = e.Message
		}

		return c.Status(code).JSON(fiber.Map{
			"errors": message,
		})
	}
}
