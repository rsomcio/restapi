package handlers

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/rsomcio/restapi/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestApp() *fiber.App {
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

	api := app.Group("/api")
	events := api.Group("/events")

	events.Post("/", CreateEvent)
	events.Get("/", GetAllEvents)
	events.Get("/:id", GetEventByID)
	events.Put("/:id", UpdateEvent)
	events.Delete("/:id", DeleteEvent)

	return app
}

func TestCreateEventValidation(t *testing.T) {
	app := setupTestApp()

	tests := []struct {
		name           string
		payload        interface{}
		expectedStatus int
		expectedError  string
	}{
		{
			name: "invalid JSON",
			payload: `{
				"name": "Test Event",
				"venue_name": "Test Venue",
				"address": "123 Test Street",
				"date": "2024-03-15",
				"time": "14:30:00"
				// invalid JSON
			}`,
			expectedStatus: 400,
			expectedError:  "Invalid request body",
		},
		{
			name: "missing required fields",
			payload: models.CreateEventRequest{
				Name: "Test Event",
			},
			expectedStatus: 400,
			expectedError:  "Name, venue_name, address, date, and time are required",
		},
		{
			name: "invalid date format",
			payload: models.CreateEventRequest{
				Name:      "Test Event",
				VenueName: "Test Venue",
				Address:   "123 Test Street",
				Date:      "2024/03/15",
				Time:      "14:30:00",
			},
			expectedStatus: 400,
			expectedError:  "Invalid date format. Use YYYY-MM-DD format",
		},
		{
			name: "invalid time format",
			payload: models.CreateEventRequest{
				Name:      "Test Event",
				VenueName: "Test Venue",
				Address:   "123 Test Street",
				Date:      "2024-03-15",
				Time:      "14:30",
			},
			expectedStatus: 400,
			expectedError:  "Invalid time format. Use HH:MM:SS format",
		},
		{
			name: "invalid email format",
			payload: models.CreateEventRequest{
				Name:         "Test Event",
				VenueName:    "Test Venue",
				Address:      "123 Test Street",
				Date:         "2024-03-15",
				Time:         "14:30:00",
				ContactEmail: stringPtr("invalid-email"),
			},
			expectedStatus: 400,
			expectedError:  "Invalid email format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var body []byte
			var err error

			if str, ok := tt.payload.(string); ok {
				body = []byte(str)
			} else {
				body, err = json.Marshal(tt.payload)
				require.NoError(t, err)
			}

			req := httptest.NewRequest("POST", "/api/events", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			require.NoError(t, err)

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			if tt.expectedError != "" {
				var response map[string]string
				err = json.NewDecoder(resp.Body).Decode(&response)
				require.NoError(t, err)
				assert.Contains(t, response["error"], tt.expectedError)
			}
		})
	}
}

func TestUpdateEventValidation(t *testing.T) {
	app := setupTestApp()

	tests := []struct {
		name           string
		eventID        string
		payload        models.UpdateEventRequest
		expectedStatus int
		expectedError  string
	}{
		{
			name:    "missing required fields",
			eventID: "123e4567-e89b-12d3-a456-426614174000",
			payload: models.UpdateEventRequest{
				Name: "Updated Event",
			},
			expectedStatus: 400,
			expectedError:  "Name, venue_name, address, date, and time are required",
		},
		{
			name:    "invalid date format",
			eventID: "123e4567-e89b-12d3-a456-426614174000",
			payload: models.UpdateEventRequest{
				Name:      "Updated Event",
				VenueName: "Updated Venue",
				Address:   "456 Updated Street",
				Date:      "2024/03/15",
				Time:      "15:30:00",
			},
			expectedStatus: 400,
			expectedError:  "Invalid date format. Use YYYY-MM-DD format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, err := json.Marshal(tt.payload)
			require.NoError(t, err)

			url := "/api/events/" + tt.eventID
			req := httptest.NewRequest("PUT", url, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			require.NoError(t, err)

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			if tt.expectedError != "" {
				var response map[string]string
				err = json.NewDecoder(resp.Body).Decode(&response)
				require.NoError(t, err)
				assert.Contains(t, response["error"], tt.expectedError)
			}
		})
	}
}

// Test that invalid route returns 404
func TestInvalidRoute(t *testing.T) {
	app := setupTestApp()

	req := httptest.NewRequest("GET", "/api/invalid", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)

	assert.Equal(t, 404, resp.StatusCode)
}

// Test HTTP methods not allowed
func TestMethodNotAllowed(t *testing.T) {
	app := setupTestApp()

	req := httptest.NewRequest("PATCH", "/api/events", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)

	assert.Equal(t, 405, resp.StatusCode)
}

func stringPtr(s string) *string {
	return &s
}