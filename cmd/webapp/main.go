package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	session "github.com/nicknad/krankentransport/auth/session"
	"github.com/nicknad/krankentransport/db"
	"github.com/nicknad/krankentransport/migrations"
	"github.com/nicknad/krankentransport/pages"
	"github.com/spf13/viper"
)

func main() {

	viper.SetConfigType("json")
	viper.SetConfigName("config")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	db.DATABASE_URL = viper.GetString("database.url")
	fmt.Println(db.DATABASE_URL)
	addr := viper.GetString("host.url")
	migrations.RunMigrations()

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
	log.Fatal(app.Listen(addr))
}
