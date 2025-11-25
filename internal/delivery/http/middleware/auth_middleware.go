package middleware

import (
	"slido-clone-backend/internal/model"
	"slido-clone-backend/internal/usecase"
	"slido-clone-backend/internal/util"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func NewAuth(userUseCase *usecase.UserUseCase, tokenUtil *util.TokenUtil) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		//// sign token, and get token from header "Authorization"
		//request := &model.VerifyUserRequest{
		//	Token: ctx.Get("Authorization", "NOT_FOUND"),
		//}
		//
		//userUseCase.Log.Debugf("Authorization : %s", request.Token)
		//
		//// split "Bearer <token>"
		//parts := strings.Split(request.Token, " ")
		//if len(parts) != 2 || parts[0] != "Bearer" {
		//	userUseCase.Log.Warnf("Invalid token format")
		//	return fiber.ErrUnauthorized
		//}
		//
		//// parse token
		//auth, err := tokenUtil.ParseToken(ctx.UserContext(), parts[1])
		//if err != nil {
		//	userUseCase.Log.Warnf("Failed to parse token: %v", err)
		//	return fiber.ErrUnauthorized
		//}
		//
		//userUseCase.Log.Debugf("User : %+v", auth.Username)
		//
		//// set user to context
		//ctx.Locals("auth", auth)
		//return ctx.Next()

		var tokenString string

		// 1. Ambil dari header Authorization
		authHeader := ctx.Get("Authorization")
		if authHeader != "" {
			parts := strings.Split(authHeader, " ")
			if len(parts) == 2 && parts[0] == "Bearer" {
				tokenString = parts[1]
			} else {
				userUseCase.Log.Warn("Invalid Authorization header format")
				return fiber.ErrUnauthorized
			}
		}

		// 2. Jika token kosong, ambil dari query (WebSocket support)
		if tokenString == "" {
			tokenString = ctx.Query("token")
			if tokenString != "" {
				userUseCase.Log.Debug("Token found in query parameter (WebSocket)")
			}
		}

		// 3. Jika tetap tidak ada token -> unauthorized
		if tokenString == "" {
			userUseCase.Log.Warn("Token not found in header or query param")
			return fiber.ErrUnauthorized
		}

		// 4. Parse token
		auth, err := tokenUtil.ParseToken(ctx.UserContext(), tokenString)
		if err != nil {
			userUseCase.Log.Warnf("Failed to parse token: %v", err)
			return fiber.ErrUnauthorized
		}

		userUseCase.Log.Debugf("Authenticated User: %+v", auth.Username)

		// 5. Simpan user auth ke context
		ctx.Locals("auth", auth)

		return ctx.Next()
	}
}

// GetUser mengambil data user dari context
func GetUser(ctx *fiber.Ctx) *model.Auth {
	return ctx.Locals("auth").(*model.Auth)
}
