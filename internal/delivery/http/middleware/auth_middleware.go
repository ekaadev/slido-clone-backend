package middleware

import (
	"reisify/internal/model"
	"reisify/internal/usecase"
	"reisify/internal/util"

	"github.com/gofiber/fiber/v2"
)

func NewAuth(userUseCase *usecase.UserUseCase, tokenUtil *util.TokenUtil) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// extract token from HTTP-only cookie
		tokenString := ctx.Cookies("token")
		if tokenString == "" {
			userUseCase.Log.Warnf("Missing auth token cookie")
			return fiber.ErrUnauthorized
		}

		// parse and validate token
		auth, err := tokenUtil.ParseToken(ctx.UserContext(), tokenString)
		if err != nil {
			userUseCase.Log.Warnf("Failed to parse token: %v", err)
			return fiber.ErrUnauthorized
		}

		userUseCase.Log.Debugf("User : %+v", auth.Username)

		// set user to context
		ctx.Locals("auth", auth)
		return ctx.Next()
	}
}

// GetUser mengambil data user dari context
func GetUser(ctx *fiber.Ctx) *model.Auth {
	return ctx.Locals("auth").(*model.Auth)
}
