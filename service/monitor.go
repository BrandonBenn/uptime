package service

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"
	"uptime/models"

	"github.com/go-co-op/gocron"
	"github.com/uptrace/bun"
)

func PingMonitors(ctx context.Context, db *bun.DB, scheduler *gocron.Scheduler) error {
	var monitors []models.Monitor
	if err := db.NewSelect().
		Model(&monitors).
		Scan(ctx); err != nil {
		return err
	}

	for _, m := range monitors {
		if err := AddMonitor(ctx, db, scheduler, m); err != nil {
			log.Fatalln("Error scheduling job", err)
		}
	}

	return nil
}

func PingMonitorJob(db *bun.DB, mon models.Monitor, job gocron.Job) (*gocron.Job, error) {
	if err := PingMonitor(job.Context(), db, mon); err != nil {
		return &job, err
	}

	return &job, nil
}

func PingMonitor(ctx context.Context, db *bun.DB, mon models.Monitor) error {
	var (
		resp *http.Response
		data models.MonitorData
		err  error

		url   = strings.ToLower(mon.Protocol + "://" + mon.URL)
		start = time.Now()
	)

	resp, err = http.Get(url)
	if err != nil {
		return err
	}

	if resp != nil {
		defer resp.Body.Close()
	}

	data.MonitorID = mon.ID
	data.ResponseTime = time.Since(start).Seconds()
	data.StatusCode = resp.StatusCode

	log.Println("PING ", data.ResponseTime, " ", url, "...")
	if _, err = db.NewInsert().Model(&data).Exec(ctx); err != nil {
		return err
	}

	return nil
}

func AddMonitor(ctx context.Context, db *bun.DB, scheduler *gocron.Scheduler, mon models.Monitor) error {
	interval := int(mon.Interval)

	if _, err := scheduler.
		Every(interval).
		Second().
		Tag(mon.Tag()).
		DoWithJobDetails(PingMonitorJob, db, mon); err != nil {
		return err
	}

	if scheduler.IsRunning() {
		log.Println("Added monitor", mon.Tag())
	}

	return nil
}

func RemoveMonitor(ctx context.Context, scheduler *gocron.Scheduler, mon models.Monitor) error {
	if err := scheduler.RemoveByTag(mon.Tag()); err != nil {
		return err
	}

	if scheduler.IsRunning() {
		log.Println("Removed monitor", mon.Tag())
	}

	return nil
}
