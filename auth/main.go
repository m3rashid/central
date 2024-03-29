package main

import (
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	"github.com/m3rashid/central/auth/handlers"
	"github.com/m3rashid/central/auth/models"
	"github.com/m3rashid/central/auth/utils"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	app := fiber.New(fiber.Config{
		CaseSensitive:         true,
		PassLocalsToViews:     true,
		DisableStartupMessage: true,
		Concurrency:           256 * 1024 * 1024,
		AppName:               "Central Auth server",
		ServerHeader:          "Central Auth server",
		RequestMethods:        []string{"GET", "POST", "HEAD", "OPTIONS"},
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			log.Println(err.Error())
			return ctx.Status(code).JSON(fiber.Map{"success": false})
		},
	})

	app.Use(cors.New())

	app.Static("/public", "./public", fiber.Static{
		Compress:  true,
		ByteRange: true,
		Browse:    true,
	})

	app.Use(favicon.New(favicon.Config{
		File: "./public/icons/favicon.ico",
		URL:  "/favicon.ico",
	}))

	if os.Getenv("SERVER_MODE") == "production" {
		app.Use(limiter.New(limiter.Config{
			Max:               100,
			Expiration:        1 * time.Minute,
			LimiterMiddleware: limiter.SlidingWindow{},
		}))
	}

	if os.Getenv("SERVER_MODE") == "development" {
		app.Use(logger.New(logger.Config{
			Format: "${time} ${status} ${latency} ${method} ${path} ${body} ${query}\n",
		}))
	}

	app.Get("/init", handlers.InitFlow)

	app.Get("/login", handlers.RenderLoginScreen)
	app.Post("/login", handlers.HandleLogin)

	app.Get("/register", handlers.RenderRegisterScreen)
	app.Post("/register", handlers.HandleRegister)

	app.Get("/select-user", handlers.RenderSelectUserScreen)
	app.Get("/handle-select-user", handlers.HandleSelectUser)

	app.Get("/consent", handlers.RenderConsentScreen)
	app.Get("/handle-consent", handlers.HandleConsent)

	app.Get("/discover", handlers.DiscoverResourceServers)

	db, err := utils.GetDb()
	if err != nil {
		log.Fatal(err)
	}
	db.AutoMigrate(models.Models...)

	// // Run this script to seed the database
	// err = initialSeedDatabase()
	// if err != nil {
	// 	log.Fatal("database seeding failed")
	// }

	go handlers.PingResourceServers()

	log.Println("Server is running")
	app.Listen(":" + os.Getenv("AUTH_SERVER_PORT"))
}
