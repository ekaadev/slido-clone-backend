package config

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

// NewFiber function untuk create fiber app yang digunakan untuk web server
func NewFiber(config *viper.Viper) *fiber.App {
	var app = fiber.New(fiber.Config{
		AppName:      config.GetString("app.name"),
		Prefork:      config.GetBool("app.prefork"),
		ErrorHandler: NewErrorHandler(),
	})

	return app
}

// NewErrorHandler function untuk membuat custom error handler di fiber
func NewErrorHandler() fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {
		code := fiber.StatusInternalServerError
		var e *fiber.Error
		if errors.As(err, &e) {
			code = e.Code
		}

		return c.Status(code).JSON(fiber.Map{
			"errors": err.Error(),
		})
	}
}
