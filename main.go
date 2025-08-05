package main

import (
	"log"
	"os"
	"runtime"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/rsomcio/restapi/database"
	"github.com/rsomcio/restapi/handlers"
)

func logError(msg string, args ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if ok {
		log.Printf("[%s:%d] "+msg, append([]interface{}{file, line}, args...)...)
	} else {
		log.Printf(msg, args...)
	}
}

func logInfo(msg string, args ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if ok {
		log.Printf("[%s:%d] "+msg, append([]interface{}{file, line}, args...)...)
	} else {
		log.Printf(msg, args...)
	}
}

func main() {
	if err := database.Connect(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer database.Close()

	if err := database.CreateTables(); err != nil {
		log.Fatal("Failed to create database tables:", err)
	}

	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			logError("Error: %v", err)
			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${method} ${path} (${latency})\n",
	}))
	app.Use(recover.New())
	app.Use(cors.New())

	api := app.Group("/api")
	events := api.Group("/events")

	events.Post("/", handlers.CreateEvent)
	events.Get("/", handlers.GetAllEvents)
	events.Get("/:id", handlers.GetEventByID)
	events.Put("/:id", handlers.UpdateEvent)
	events.Delete("/:id", handlers.DeleteEvent)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	logInfo("Server starting on port %s", port)
	log.Fatal(app.Listen(":" + port))
}