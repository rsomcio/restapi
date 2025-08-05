package models

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEventJSONSerialization(t *testing.T) {
	description := "Test event description"
	contactEmail := "test@example.com"
	contactMobile := "+1234567890"
	contactInstagram := "testevent"

	event := Event{
		ID:               "123e4567-e89b-12d3-a456-426614174000",
		Name:             "Test Event",
		Description:      &description,
		VenueName:        "Test Venue",
		Address:          "123 Test Street",
		Date:             "2024-03-15",
		Time:             "14:30:00",
		ContactMobile:    &contactMobile,
		ContactEmail:     &contactEmail,
		ContactInstagram: &contactInstagram,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	jsonData, err := json.Marshal(event)
	require.NoError(t, err)

	var unmarshaled Event
	err = json.Unmarshal(jsonData, &unmarshaled)
	require.NoError(t, err)

	assert.Equal(t, event.ID, unmarshaled.ID)
	assert.Equal(t, event.Name, unmarshaled.Name)
	assert.Equal(t, event.VenueName, unmarshaled.VenueName)
	assert.Equal(t, event.Address, unmarshaled.Address)
	assert.Equal(t, event.Date, unmarshaled.Date)
	assert.Equal(t, event.Time, unmarshaled.Time)
	assert.Equal(t, *event.Description, *unmarshaled.Description)
	assert.Equal(t, *event.ContactEmail, *unmarshaled.ContactEmail)
	assert.Equal(t, *event.ContactMobile, *unmarshaled.ContactMobile)
	assert.Equal(t, *event.ContactInstagram, *unmarshaled.ContactInstagram)
}

func TestEventWithNilOptionalFields(t *testing.T) {
	event := Event{
		ID:        "123e4567-e89b-12d3-a456-426614174000",
		Name:      "Test Event",
		VenueName: "Test Venue",
		Address:   "123 Test Street",
		Date:      "2024-03-15",
		Time:      "14:30:00",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	jsonData, err := json.Marshal(event)
	require.NoError(t, err)

	var unmarshaled Event
	err = json.Unmarshal(jsonData, &unmarshaled)
	require.NoError(t, err)

	assert.Equal(t, event.ID, unmarshaled.ID)
	assert.Equal(t, event.Name, unmarshaled.Name)
	assert.Nil(t, unmarshaled.Description)
	assert.Nil(t, unmarshaled.ContactEmail)
	assert.Nil(t, unmarshaled.ContactMobile)
	assert.Nil(t, unmarshaled.ContactInstagram)
}

func TestCreateEventRequestValidation(t *testing.T) {
	tests := []struct {
		name    string
		request CreateEventRequest
		valid   bool
	}{
		{
			name: "valid request with all fields",
			request: CreateEventRequest{
				Name:      "Test Event",
				VenueName: "Test Venue",
				Address:   "123 Test Street",
				Date:      "2024-03-15",
				Time:      "14:30:00",
			},
			valid: true,
		},
		{
			name: "valid request with optional fields",
			request: CreateEventRequest{
				Name:         "Test Event",
				VenueName:    "Test Venue",
				Address:      "123 Test Street",
				Date:         "2024-03-15",
				Time:         "14:30:00",
				Description:  stringPtr("Description"),
				ContactEmail: stringPtr("test@example.com"),
			},
			valid: true,
		},
		{
			name: "invalid request missing name",
			request: CreateEventRequest{
				VenueName: "Test Venue",
				Address:   "123 Test Street",
				Date:      "2024-03-15",
				Time:      "14:30:00",
			},
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonData, err := json.Marshal(tt.request)
			require.NoError(t, err)

			var unmarshaled CreateEventRequest
			err = json.Unmarshal(jsonData, &unmarshaled)
			require.NoError(t, err)

			if tt.valid {
				assert.Equal(t, tt.request.Name, unmarshaled.Name)
				assert.Equal(t, tt.request.VenueName, unmarshaled.VenueName)
				assert.Equal(t, tt.request.Address, unmarshaled.Address)
			}
		})
	}
}

func stringPtr(s string) *string {
	return &s
}