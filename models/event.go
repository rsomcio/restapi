package models

import (
	"time"
)

type Event struct {
	ID               string     `json:"id" db:"id"`
	Name             string     `json:"name" db:"name"`
	Description      *string    `json:"description" db:"description"`
	VenueName        string     `json:"venue_name" db:"venue_name"`
	Address          string     `json:"address" db:"address"`
	Date             string     `json:"date" db:"date"`
	Time             string     `json:"time" db:"time"`
	ContactMobile    *string    `json:"contact_mobile" db:"contact_mobile"`
	ContactEmail     *string    `json:"contact_email" db:"contact_email"`
	ContactInstagram *string    `json:"contact_instagram" db:"contact_instagram"`
	CreatedAt        time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at" db:"updated_at"`
}

type CreateEventRequest struct {
	Name             string  `json:"name"`
	Description      *string `json:"description"`
	VenueName        string  `json:"venue_name"`
	Address          string  `json:"address"`
	Date             string  `json:"date"`
	Time             string  `json:"time"`
	ContactMobile    *string `json:"contact_mobile"`
	ContactEmail     *string `json:"contact_email"`
	ContactInstagram *string `json:"contact_instagram"`
}

type UpdateEventRequest struct {
	Name             string  `json:"name"`
	Description      *string `json:"description"`
	VenueName        string  `json:"venue_name"`
	Address          string  `json:"address"`
	Date             string  `json:"date"`
	Time             string  `json:"time"`
	ContactMobile    *string `json:"contact_mobile"`
	ContactEmail     *string `json:"contact_email"`
	ContactInstagram *string `json:"contact_instagram"`
}