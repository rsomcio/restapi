package handlers

import (
	"log"
	"regexp"
	"runtime"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rsomcio/restapi/database"
	"github.com/rsomcio/restapi/models"
	"gopkg.in/go-playground/validator.v9"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

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

func validateEmail(email string) bool {
	if email == "" {
		return true
	}
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

func validateDateFormat(date string) bool {
	_, err := time.Parse("2006-01-02", date)
	return err == nil
}

func validateTimeFormat(timeStr string) bool {
	_, err := time.Parse("15:04:05", timeStr)
	return err == nil
}

func CreateEvent(c *fiber.Ctx) error {
	var req models.CreateEventRequest
	if err := c.BodyParser(&req); err != nil {
		logError("Error parsing request body: %v", err)
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if req.Name == "" || req.VenueName == "" || req.Address == "" || req.Date == "" || req.Time == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Name, venue_name, address, date, and time are required"})
	}

	if !validateDateFormat(req.Date) {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid date format. Use YYYY-MM-DD format"})
	}

	if !validateTimeFormat(req.Time) {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid time format. Use HH:MM:SS format"})
	}

	if req.ContactEmail != nil && !validateEmail(*req.ContactEmail) {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid email format"})
	}

	query := `
		INSERT INTO events (name, description, venue_name, address, date, time, contact_mobile, contact_email, contact_instagram)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, name, description, venue_name, address, date, time, contact_mobile, contact_email, contact_instagram, created_at, updated_at`

	var event models.Event
	err := database.DB.QueryRow(query, req.Name, req.Description, req.VenueName, req.Address, req.Date, req.Time, req.ContactMobile, req.ContactEmail, req.ContactInstagram).Scan(
		&event.ID, &event.Name, &event.Description, &event.VenueName, &event.Address, &event.Date, &event.Time, &event.ContactMobile, &event.ContactEmail, &event.ContactInstagram, &event.CreatedAt, &event.UpdatedAt)

	if err != nil {
		logError("Error creating event: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create event"})
	}

	logInfo("Created event with ID: %s", event.ID)
	return c.Status(201).JSON(event)
}

func GetAllEvents(c *fiber.Ctx) error {
	var events []models.Event
	query := "SELECT id, name, description, venue_name, address, date, time, contact_mobile, contact_email, contact_instagram, created_at, updated_at FROM events ORDER BY date, time"

	err := database.DB.Select(&events, query)
	if err != nil {
		logError("Error fetching events: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch events"})
	}

	logInfo("Fetched %d events", len(events))
	return c.JSON(events)
}

func GetEventByID(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Event ID is required"})
	}

	var event models.Event
	query := "SELECT id, name, description, venue_name, address, date, time, contact_mobile, contact_email, contact_instagram, created_at, updated_at FROM events WHERE id = $1"

	err := database.DB.Get(&event, query, id)
	if err != nil {
		logError("Error fetching event %s: %v", id, err)
		return c.Status(404).JSON(fiber.Map{"error": "Event not found"})
	}

	logInfo("Fetched event with ID: %s", id)
	return c.JSON(event)
}

func UpdateEvent(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Event ID is required"})
	}

	var req models.UpdateEventRequest
	if err := c.BodyParser(&req); err != nil {
		logError("Error parsing request body: %v", err)
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if req.Name == "" || req.VenueName == "" || req.Address == "" || req.Date == "" || req.Time == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Name, venue_name, address, date, and time are required"})
	}

	if !validateDateFormat(req.Date) {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid date format. Use YYYY-MM-DD format"})
	}

	if !validateTimeFormat(req.Time) {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid time format. Use HH:MM:SS format"})
	}

	if req.ContactEmail != nil && !validateEmail(*req.ContactEmail) {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid email format"})
	}

	var existingEvent models.Event
	checkQuery := "SELECT id FROM events WHERE id = $1"
	err := database.DB.Get(&existingEvent, checkQuery, id)
	if err != nil {
		logError("Event %s not found: %v", id, err)
		return c.Status(404).JSON(fiber.Map{"error": "Event not found"})
	}

	query := `
		UPDATE events 
		SET name = $1, description = $2, venue_name = $3, address = $4, date = $5, time = $6, 
		    contact_mobile = $7, contact_email = $8, contact_instagram = $9, updated_at = CURRENT_TIMESTAMP
		WHERE id = $10
		RETURNING id, name, description, venue_name, address, date, time, contact_mobile, contact_email, contact_instagram, created_at, updated_at`

	var event models.Event
	err = database.DB.QueryRow(query, req.Name, req.Description, req.VenueName, req.Address, req.Date, req.Time, req.ContactMobile, req.ContactEmail, req.ContactInstagram, id).Scan(
		&event.ID, &event.Name, &event.Description, &event.VenueName, &event.Address, &event.Date, &event.Time, &event.ContactMobile, &event.ContactEmail, &event.ContactInstagram, &event.CreatedAt, &event.UpdatedAt)

	if err != nil {
		logError("Error updating event %s: %v", id, err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update event"})
	}

	logInfo("Updated event with ID: %s", id)
	return c.JSON(event)
}

func DeleteEvent(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Event ID is required"})
	}

	var existingEvent models.Event
	checkQuery := "SELECT id FROM events WHERE id = $1"
	err := database.DB.Get(&existingEvent, checkQuery, id)
	if err != nil {
		logError("Event %s not found: %v", id, err)
		return c.Status(404).JSON(fiber.Map{"error": "Event not found"})
	}

	query := "DELETE FROM events WHERE id = $1"
	result, err := database.DB.Exec(query, id)
	if err != nil {
		logError("Error deleting event %s: %v", id, err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete event"})
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "Event not found"})
	}

	logInfo("Deleted event with ID: %s", id)
	return c.SendStatus(204)
}