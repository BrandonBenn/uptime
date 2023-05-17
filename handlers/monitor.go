package handlers

import (
	"database/sql"
	"net/http"
	"uptime/models"
	"uptime/service"

	"github.com/go-co-op/gocron"
	"github.com/gofiber/fiber/v2"
	"github.com/uptrace/bun"
)

type Monitor struct {
	DB        *bun.DB
	Scheduler *gocron.Scheduler
}

func (h *Monitor) Index(c *fiber.Ctx) error {
	var (
		ctx         = c.Context()
		currentUser = c.Locals("user").(models.User)
		monitors    = []models.Monitor{}
		data        = map[string]any{}
	)

	if err := h.DB.NewSelect().
		Model(&monitors).
		Where("user_id = ?", currentUser.ID).
		Scan(ctx, &monitors); err != nil {
		return err
	}

	data = fiber.Map{"monitors": monitors}
	return c.JSON(data)
}

func (h *Monitor) New(c *fiber.Ctx) error {
	var (
		currentUser = c.Locals("user").(models.User)
	)

	monData := models.MonitorData{
		StatusCode: 404,
		Monitor: &models.Monitor{
			UserID:   currentUser.ID,
			Protocol: "HTTPS",
		},
	}

	return c.Render("partials/monitor_item", monData, "")
}

func (h *Monitor) Create(c *fiber.Ctx) error {
	var (
		ctx         = c.Context()
		currentUser = c.Locals("user").(models.User)
		monitor     = models.Monitor{}
		err         error
	)

	if err = c.BodyParser(&monitor); err != nil {
		return err
	}

	monitor.UserID = currentUser.ID
	service.AddMonitor(ctx, h.DB, h.Scheduler, monitor)

	if _, err = h.DB.NewInsert().
		Model(&monitor).
		Where("user_id = ?", currentUser.ID).
		Returning("*").
		Exec(ctx); err != nil {
		return err
	}

	monData := models.MonitorData{
		StatusCode: 404,
		Monitor:    &monitor,
	}

	return c.Render("partials/monitor_item", monData, "")
}

func (h *Monitor) Update(c *fiber.Ctx) error {
	var (
		currentUser = c.Locals("user").(models.User)
		monitor     = models.Monitor{}
		ctx         = c.Context()
		monitorID   int
		err         error
	)

	if monitorID, err = c.ParamsInt("id"); err != nil {
		return err
	}

	if err = c.BodyParser(&monitor); err != nil {
		return err
	}

	service.RemoveMonitor(ctx, h.Scheduler, monitor)
	monitor.UserID = currentUser.ID
	monitor.ID = int64(monitorID)
	if _, err := h.DB.NewUpdate().
		Model(&monitor).
		Where("id = ?", monitorID).
		Where("user_id = ?", currentUser.ID).
		Exec(ctx); err != nil {
		return err
	}
	service.AddMonitor(ctx, h.DB, h.Scheduler, monitor)

	return c.Render("partials/monitor_item", monitor, "")
}

func (h *Monitor) Delete(c *fiber.Ctx) error {
	var (
		currentUser = c.Locals("user").(models.User)
		monitor     = models.Monitor{}
		ctx         = c.Context()
		monitorID   int
		err         error
	)

	if monitorID, err = c.ParamsInt("id"); err != nil {
		return err
	}

	if monitorID == 0 {
		return c.SendStatus(http.StatusAccepted)
	}

	service.RemoveMonitor(ctx, h.Scheduler, monitor)
	var result sql.Result
	if result, err = h.DB.NewDelete().
		Model(&monitor).
		Where("id = ?", monitorID).
		Where("user_id = ?", currentUser.ID).
		Exec(ctx); err != nil {
		return err
	}

	if result, err = h.DB.NewDelete().
		Model((*models.MonitorData)(nil)).
		Where("monitor_id = ?", monitorID).
		Exec(ctx); err != nil {
		return err
	}

	var rows int64
	if rows, err = result.RowsAffected(); err != nil {
		return err
	} else if rows == 0 {
		return fiber.NewError(http.StatusNotFound, "Monitor not found")
	}

	if _, err := h.DB.NewSelect().
		Model(&monitor).
		Where("id = ?", monitorID).
		Where("user_id = ?", currentUser.ID).
		Exec(ctx); err != nil {
		return err
	}

	return c.SendStatus(http.StatusAccepted)
}
