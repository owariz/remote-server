package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"

	"github.com/owariz/remote-server/internal/config"
	"github.com/owariz/remote-server/internal/handlers"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	cfg := config.New()

	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New())

	api := app.Group("/api/v1")

	// api.Use(middleware.Auth())

	api.Get("/status", handlers.GetStatus)
	api.Get("/metrics", handlers.GetMetrics)

	services := api.Group("/service/:name")
	services.Post("/restart", handlers.RestartService)
	services.Get("/logs", handlers.GetServiceLogs)

	api.Post("/exec", handlers.ExecuteCommand)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := app.Listen(":" + cfg.Port); err != nil {
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	log.Printf("Server started on port %s\n", cfg.Port)

	<-c
	log.Println("Gracefully shutting down...")
	_ = app.Shutdown()
}
