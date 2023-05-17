package handlers

import (
	"uptime/models"

	"github.com/gofiber/fiber/v2"
	"github.com/uptrace/bun"
)

type Root struct {
	DB *bun.DB
}

func (h *Root) Index(c *fiber.Ctx) error {
	var (
		currentUser = c.Locals("user").(models.User)
		data        = map[string]any{}
	)

	data["IsLoggedIn"] = currentUser.ID != 0
	return c.Render("views/index", data)

}

func (h *Root) Login(c *fiber.Ctx) error {
	var (
		currentUser = c.Locals("user").(models.User)
		data        = map[string]any{}
	)

	if currentUser.ID != 0 {
		return c.Redirect("/")
	}

	data["Title"] = "Login"
	return c.Render("views/login", data)
}
