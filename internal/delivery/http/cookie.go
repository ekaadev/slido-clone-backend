package http

import "github.com/gofiber/fiber/v2"

const (
	authCookieName = "token"
	// 30 days in seconds, matching the JWT expiry duration
	authCookieMaxAge = 86400 * 30
)

// cookieSecure controls the Secure flag on auth cookies.
// Set to true in production (HTTPS). Defaults to false to allow HTTP in development.
var cookieSecure bool

// SetCookieSecure configures the Secure flag for all auth cookies.
// Call this once during application bootstrap with the value of COOKIE_SECURE env var.
func SetCookieSecure(secure bool) {
	cookieSecure = secure
}

// setAuthCookie sets the JWT token as an HTTP-only, SameSite=Lax cookie.
func setAuthCookie(ctx *fiber.Ctx, token string) {
	ctx.Cookie(&fiber.Cookie{
		Name:     authCookieName,
		Value:    token,
		HTTPOnly: true,
		Secure:   cookieSecure,
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
		Secure:   cookieSecure,
		SameSite: "Lax",
		Path:     "/",
		MaxAge:   -1,
	})
}
