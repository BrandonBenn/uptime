package middleware

import (
	"uptime/models"
	"uptime/service"

	"github.com/gofiber/fiber/v2"
	"github.com/uptrace/bun"
)

const SessionCookieName = "_uptime"

func UserProtected(db *bun.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		currentUser := c.Locals("user").(models.User)
		if currentUser.ID == 0 {
			return c.Redirect("/login", fiber.StatusFound)
		}

		return c.Next()
	}
}

func CurrentUser(db *bun.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Cookies(SessionCookieName)
		var user models.User
		_ = service.FindUserByToken(c.Context(), db, token, &user)
		c.Locals("user", user)
		return c.Next()
	}
}
