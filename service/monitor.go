package service

import (
	"context"
	"fmt"
	"log"
	"math/rand"
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
	defer resp.Body.Close()

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
	var (
		interval = int(mon.Interval*1000) + int(rand.Intn(1000)) // To avoid all monitors pinging at the same time
		tag      = fmt.Sprint(mon.Name, "-", mon.ID)
	)

	if _, err := scheduler.
		Every(interval).
		Millisecond().
		Tag(tag).
		DoWithJobDetails(PingMonitorJob, db, mon); err != nil {
		return err
	}

	return nil
}

func RemoveMonitor(ctx context.Context, scheduler *gocron.Scheduler, mon models.Monitor) error {
	tag := fmt.Sprint(mon.Name, "-", mon.ID)
	if err := scheduler.RemoveByTag(tag); err != nil {
		return err
	}

	return nil
}
