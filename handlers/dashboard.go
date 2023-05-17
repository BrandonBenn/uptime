package handlers

import (
	"uptime/models"

	"github.com/gofiber/fiber/v2"
	"github.com/uptrace/bun"
)

type Dashboard struct {
	DB *bun.DB
}

func (h *Dashboard) Index(c *fiber.Ctx) error {
	var (
		currentUser = c.Locals("user").(models.User)
		monData     = []models.MonitorData{}
		data        = map[string]any{}
		ctx         = c.Context()
	)

	if err := h.DB.NewSelect().
		Model(&monData).
		Relation("Monitor").
		Join(`join (
				select monitor_id, max(created_at) as max_created_at
				from monitor_data group by monitor_id
		) latest on md.monitor_id = latest.monitor_id and md.created_at = latest.max_created_at`).
		Where("monitor.user_id = ?", currentUser.ID).
		Scan(ctx); err != nil {
		return err
	}

	data["IsLoggedIn"] = currentUser.ID != 0
	data["Email"] = currentUser.Email
	data["MonitorData"] = monData
	data["Title"] = "Dashboard"
	return c.Render("views/dashboard", data)
}

func (h *Dashboard) Show(c *fiber.Ctx) error {
	var (
		currentUser = c.Locals("user").(models.User)
		monitors    = []models.Monitor{}
		data        = map[string]any{}
		ctx         = c.Context()
	)

	if err := h.DB.NewSelect().
		Model(&monitors).
		Where("user_id = ?", currentUser.ID).
		Scan(ctx); err != nil {
		return err
	}

	data["IsLoggedIn"] = currentUser.ID != 0
	data["Username"] = currentUser.Email
	data["Monitors"] = monitors
	data["Title"] = "Dashboard"
	return c.Render("views/dashboard", data)
}
