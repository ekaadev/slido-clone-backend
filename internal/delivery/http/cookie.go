package http

import "github.com/gofiber/fiber/v2"

const (
	authCookieName = "token"
	// 30 days in seconds, matching the JWT expiry duration
	authCookieMaxAge = 86400 * 30
)

// setAuthCookie sets the JWT token as an HTTP-only, Secure, SameSite=Lax cookie.
func setAuthCookie(ctx *fiber.Ctx, token string) {
	ctx.Cookie(&fiber.Cookie{
		Name:     authCookieName,
		Value:    token,
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Lax",
		Path:     "/",
		MaxAge:   authCookieMaxAge,
	})
}

// clearAuthCookie removes the auth cookie by setting MaxAge to -1.
func clearAuthCookie(ctx *fiber.Ctx) {
	ctx.Cookie(&fiber.Cookie{
		Name:     authCookieName,
		Value:    "",
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Lax",
		Path:     "/",
		MaxAge:   -1,
	})
}
