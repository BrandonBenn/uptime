package main

import (
	"context"
	"log"
	"os"
	"strings"
	"time"
	database "uptime/db"
	"uptime/handlers"
	"uptime/handlers/middleware"
	"uptime/service"

	"github.com/go-co-op/gocron"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/middleware/encryptcookie"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/storage/sqlite3"
	"github.com/gofiber/template/html"
	"github.com/uptrace/bun"
)

func main() {
	db, err := database.NewSqliteDB()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	app := fiber.New(fiber.Config{
		Views: html.New("templates", ".html").
			AddFunc("ToUpper", strings.ToUpper),
		ViewsLayout: "views/base",
	})
	app.Static("/static", "./templates/static")
	app.Use(
		logger.New(),
		encryptcookie.New(encryptcookie.Config{Key: os.Getenv("SECRET_KEY")}),
		csrf.New(csrf.Config{Storage: sqlite3.New()}),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	scheduler := gocron.NewScheduler(time.UTC)
	if err := setupJobs(ctx, db, scheduler); err != nil {
		panic(err)
	}

	setupRoutes(app, db, scheduler)
	scheduler.StartAsync()

	port := getEnv("PORT", "3000")
	host := getEnv("HOST", "127.0.0.1")
	log.Fatal(app.Listen(host + ":" + port))
}

func setupJobs(ctx context.Context, db *bun.DB, scheduler *gocron.Scheduler) error {
	scheduler.SingletonMode()
	if err := service.PingMonitors(ctx, db, scheduler); err != nil {
		return err
	}

	return nil
}

func setupRoutes(app *fiber.App, db *bun.DB, scheduler *gocron.Scheduler) {
	currentUser := middleware.CurrentUser(db)
	userProtected := middleware.UserProtected(db)

	app.Use(currentUser)
	app.Route("/", func(r fiber.Router) {
		root := handlers.Root{DB: db}
		r.Get("/", root.Index)
		r.Get("/login", root.Login)
	})

	app.Route("/", func(r fiber.Router) {
		handler := handlers.Session{DB: db}
		r.Post("/login", handler.Login)
		r.Get("/login/verify", handler.Verify)
		r.Get("/logout", handler.Logout)
	})

	app.Route("/monitors", func(r fiber.Router) {
		r.Use(userProtected)

		handler := handlers.Monitor{DB: db, Scheduler: scheduler}
		r.Post("/", handler.Create)
		r.Get("/new", handler.New)
		r.Put("/:id", handler.Update)
		r.Delete("/:id", handler.Delete)
	})

	app.Route("/dashboard", func(r fiber.Router) {
		r.Use(userProtected)

		handler := handlers.Dashboard{DB: db}
		r.Get("/", handler.Index)
	})
}

func getEnv(env, fallback string) string {
	value, exists := os.LookupEnv("HOST")
	if !exists {
		return fallback
	}

	return value
}
