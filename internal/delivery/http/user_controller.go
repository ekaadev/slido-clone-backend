package http

import (
	"slido-clone-backend/internal/model"
	"slido-clone-backend/internal/usecase"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type UserController struct {
	Log         *logrus.Logger
	UserUseCase *usecase.UserUseCase
}

// NewUserController create new instance of UserController
func NewUserController(log *logrus.Logger, userUseCase *usecase.UserUseCase) *UserController {
	return &UserController{
		Log:         log,
		UserUseCase: userUseCase,
	}
}

// Register handler yang digunakan untuk create user baru (call usecase create user)
func (c *UserController) Register(ctx *fiber.Ctx) error {
	// create model register user request
	request := new(model.RegisterUserRequest)

	// parsing body payload
	err := ctx.BodyParser(request)
	if err != nil {
		c.Log.Warnf("Body parse failed: %s", err)
		return fiber.ErrBadRequest
	}

	// call usecase to create user
	response, err := c.UserUseCase.Create(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Create failed: %s", err)
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(model.WebResponse{
		Data: response,
	})
}

// Login handler yang digunakan untuk login user (call usecase login user)
func (c *UserController) Login(ctx *fiber.Ctx) error {
	// create model login user request
	request := new(model.LoginUserRequest)

	// parsing body payload
	err := ctx.BodyParser(request)
	if err != nil {
		c.Log.Warnf("Body parse failed: %s", err)
		return fiber.ErrBadRequest
	}

	// call usecase to login user
	response, err := c.UserUseCase.Login(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Login failed: %s", err)
		return err
	}

	// return response
	return ctx.Status(fiber.StatusOK).JSON(model.WebResponse{
		Data: response,
	})
}

// Anon handler digunakan untuk anonymous user (call usecase anon user)
func (c *UserController) Anon(ctx *fiber.Ctx) error {
	// create model anon user request
	request := new(model.AnonymousUserRequest)

	// parsing body payload
	if err := ctx.BodyParser(request); err != nil {
		c.Log.Warnf("Body parse failed: %s", err)
		return fiber.ErrBadRequest
	}

	// call usecase to anon user
	response, err := c.UserUseCase.Anon(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Anon failed: %s", err)
		return err
	}

	// return response
	return ctx.Status(fiber.StatusOK).JSON(model.WebResponse{
		Data: response,
	})
}

// Logout handler untuk logout user dan invalidate token
func (c *UserController) Logout(ctx *fiber.Ctx) error {
	// extract token from Authorization header
	authHeader := ctx.Get("Authorization")
	if authHeader == "" {
		return fiber.ErrUnauthorized
	}

	// remove "Bearer " prefix
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader {
		// no Bearer prefix found
		return fiber.ErrUnauthorized
	}

	// call usecase to logout
	if err := c.UserUseCase.Logout(ctx.UserContext(), tokenString); err != nil {
		c.Log.Warnf("Logout failed: %s", err)
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(model.WebResponse{
		Data: map[string]string{
			"message": "Successfully logged out",
		},
	})
}
