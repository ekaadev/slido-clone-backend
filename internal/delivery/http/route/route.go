package route

import (
	"slido-clone-backend/internal/delivery/http"

	"github.com/gofiber/fiber/v2"
)

type RouteConfig struct {
	App            *fiber.App
	UserController *http.UserController
	RoomController *http.RoomController
	AuthMiddleware fiber.Handler
}

// Setup running all route setup here
func (c *RouteConfig) Setup() {
	c.SetupGuestRoute()
	c.SetupAuthRoute()
}

// SetupGuestRoute tambahkan route yang bisa diakses tanpa autentikasi
func (c *RouteConfig) SetupGuestRoute() {
	c.App.Post("/api/v1/users/register", c.UserController.Register)
	c.App.Post("/api/v1/users/login", c.UserController.Login)
}

func (c *RouteConfig) SetupAuthRoute() {
	c.App.Use(c.AuthMiddleware)
	c.App.Post("/api/v1/rooms", c.RoomController.Create)
}
