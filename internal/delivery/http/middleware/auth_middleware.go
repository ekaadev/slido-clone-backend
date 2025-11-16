package middleware

import (
	"slido-clone-backend/internal/model"
	"slido-clone-backend/internal/usecase"
	"slido-clone-backend/internal/util"

	"github.com/gofiber/fiber/v2"
)

func NewAuth(userUseCase *usecase.UserUseCase, tokenUtil *util.TokenUtil) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// sign token, and get token from header "Authorization"
		request := &model.VerifyUserRequest{
			Token: ctx.Get("Authorization", "NOT_FOUND"),
		}

		userUseCase.Log.Debugf("Authorization : %s", request.Token)

		// parse token
		auth, err := tokenUtil.ParseToken(ctx.UserContext(), request.Token)
		if err != nil {
			userUseCase.Log.Warnf("Failed to parse token: %v", err)
			return fiber.ErrUnauthorized
		}

		userUseCase.Log.Debugf("User : %+v", auth.ID)

		// set user to context
		ctx.Locals("auth", auth)
		return ctx.Next()
	}
}

// GetUser mengambil data user dari context
func GetUser(ctx *fiber.Ctx) *model.Auth {
	return ctx.Locals("auth").(*model.Auth)
}
