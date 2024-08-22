package main

import (
	"flag"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	session "github.com/nicknad/krankentransport/auth/session"
	"github.com/nicknad/krankentransport/db"
	"github.com/nicknad/krankentransport/pages"
)

func main() {
	var port string
	var databaseUrl string
	flag.StringVar(&port, "p", ":2209", "Exposed Port")
	flag.StringVar(&databaseUrl, "db", "./krankentransport.db", "Database Filename")
	flag.Parse()

	db.DATABASE_URL = databaseUrl

	go func() {
		for {
			session.ClearExpiredSessions()
			time.Sleep(5 * time.Minute)
		}
	}()

	app := fiber.New(fiber.Config{
		Network: fiber.NetworkTCP,
	})

	app.Use(logger.New())
	app.Static("/public", "./public")

	app.Use(func(c *fiber.Ctx) error {
		c.Set("Cache-Control", "no-store, no-cache, must-revalidate, proxy-revalidate")
		c.Set("Pragma", "no-cache")
		c.Set("Expires", "0")
		c.Set("Surrogate-Control", "no-store")
		return c.Next()
	})
	pages.RegisterRoutes(app)
	log.Fatal(app.Listen(port))
}
