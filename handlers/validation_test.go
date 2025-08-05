package handlers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		email string
		valid bool
	}{
		{"test@example.com", true},
		{"user.name@domain.co.uk", true},
		{"", true}, // empty is valid (optional field)
		{"invalid-email", false},
		{"@domain.com", false},
		{"user@", false},
		{"user@domain", false},
	}

	for _, tt := range tests {
		t.Run(tt.email, func(t *testing.T) {
			result := validateEmail(tt.email)
			assert.Equal(t, tt.valid, result)
		})
	}
}

func TestValidateDateFormat(t *testing.T) {
	tests := []struct {
		date  string
		valid bool
	}{
		{"2024-03-15", true},
		{"2024-12-31", true},
		{"2024-01-01", true},
		{"24-03-15", false},
		{"2024/03/15", false},
		{"2024-3-15", false},
		{"2024-03-32", false},
		{"invalid", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.date, func(t *testing.T) {
			result := validateDateFormat(tt.date)
			assert.Equal(t, tt.valid, result)
		})
	}
}

func TestValidateTimeFormat(t *testing.T) {
	tests := []struct {
		time  string
		valid bool
	}{
		{"14:30:00", true},
		{"00:00:00", true},
		{"23:59:59", true},
		{"14:30", false},
		{"14:30:00.000", true}, // Go's time.Parse actually accepts this format
		{"24:00:00", false},
		{"14:60:00", false},
		{"14:30:60", false},
		{"invalid", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.time, func(t *testing.T) {
			result := validateTimeFormat(tt.time)
			assert.Equal(t, tt.valid, result)
		})
	}
}