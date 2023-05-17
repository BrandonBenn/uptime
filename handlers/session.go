package handlers

import (
	"fmt"
	"net/mail"
	"uptime/handlers/middleware"
	"uptime/models"
	"uptime/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/encryptcookie"
	"github.com/uptrace/bun"
)

type Session struct {
	DB *bun.DB
}

func (h *Session) Login(c *fiber.Ctx) error {
	email := c.FormValue("email")
	if _, err := mail.ParseAddress(email); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.Render("views/login", map[string]any{
			"Title": "Login",
			"Error": "Email is required.",
		})
	}

	if err := service.SendVerificationEmail(c.Context(), h.DB, email); err != nil {
		return err
	}
	return c.Render("views/login", map[string]any{
		"Title":   "Login",
		"Message": "Verification email has been sent to your email address.",
	})
}

func (h *Session) Verify(c *fiber.Ctx) error {
	var (
		user models.User
		err  error
		code string
	)
	code = c.Query("code")
	if code == "" {
		return fmt.Errorf("Code is required")
	}

	if user, err = service.ValidateEmailCode(c.Context(), h.DB, code); err != nil {
		return err
	}

	token := encryptcookie.GenerateKey()
	if err = h.DB.NewRaw(`
		insert into sessions (user_id, token)
		values (?, ?) on conflict (user_id) do update
		set token = ?, created_at = current_timestamp
		returning token
`, user.ID, token, token,
	).Scan(c.Context(), &token); err != nil {
		return err
	}

	c.Cookie(&fiber.Cookie{
		Name:     middleware.SessionCookieName,
		Value:    token,
		HTTPOnly: true,
	})

	return c.Redirect("/")
}

func (h *Session) Logout(c *fiber.Ctx) error {
	token := c.Cookies(middleware.SessionCookieName)
	if token == "" {
		return c.Redirect("/login")
	}

	if _, err := h.DB.Exec("delete from sessions where token = ?", token); err != nil {
		return err
	}

	c.ClearCookie(middleware.SessionCookieName)
	return c.Redirect("/")
}
