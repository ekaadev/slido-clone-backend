package route

import (
	"slido-clone-backend/internal/delivery/http"

	"github.com/gofiber/fiber/v2"
)

type RouteConfig struct {
	App                   *fiber.App
	UserController        *http.UserController
	RoomController        *http.RoomController
	ParticipantController *http.ParticipantController
	AuthMiddleware        fiber.Handler
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
	c.App.Post("/api/v1/users/anonymous", c.UserController.Anon)

	c.App.Get("/api/v1/rooms/:room_code", c.RoomController.Get)
}

func (c *RouteConfig) SetupAuthRoute() {
	c.App.Use(c.AuthMiddleware)
	c.App.Post("/api/v1/rooms", c.RoomController.Create)
	c.App.Patch("/api/v1/rooms/:room_id/close", c.RoomController.UpdateToClosed)
	c.App.Post("/api/v1/rooms/:room_code/join", c.ParticipantController.Join)

	c.App.Get("/api/v1/users/me/rooms", c.RoomController.Search)
}
