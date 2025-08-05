# Events REST API Specification

## Overview
A simple REST API service for managing events built with Go, Fiber framework, and SQLx.

## Technical Stack
- **Language**: Go
- **Framework**: Fiber v2
- **Database**: PostgreSQL
- **Database Library**: SQLx
- **Port**: 3000

## Database Schema

```sql
CREATE TABLE events (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    venue_name VARCHAR(255) NOT NULL,
    address TEXT NOT NULL,
    date DATE NOT NULL,
    time TIME NOT NULL,
    contact_mobile VARCHAR(20),
    contact_email VARCHAR(255),
    contact_instagram VARCHAR(100),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
```

## Data Model

### Event
```json
{
  "id": "uuid",
  "name": "string",
  "description": "string (optional)",
  "venue_name": "string",
  "address": "string",
  "date": "2024-03-15",
  "time": "14:30:00",
  "contact_mobile": "string (optional)",
  "contact_email": "string (optional)",
  "contact_instagram": "string (optional)",
  "created_at": "2024-03-15T10:30:00Z",
  "updated_at": "2024-03-15T10:30:00Z"
}
```

**Field Descriptions:**
- `id`: UUID, automatically generated
- `name`: Event name (required, max 255 chars)
- `description`: Optional detailed event description
- `venue_name`: Name of the venue (required, max 255 chars)
- `address`: Full address of the venue (required)
- `date`: Event date in YYYY-MM-DD format (required)
- `time`: Event time in HH:MM:SS format (required)
- `contact_mobile`: Optional contact phone number (max 20 chars)
- `contact_email`: Optional contact email (max 255 chars)
- `contact_instagram`: Optional Instagram handle (max 100 chars)
- `created_at`: Timestamp when record was created (auto-generated)
- `updated_at`: Timestamp when record was last updated (auto-generated)

## API Endpoints

### 1. Create Event
- **Method**: `POST`
- **Path**: `/api/events`
- **Request Body**: Event object (without `id`, `created_at`, `updated_at`)
- **Response**: Created event with generated fields
- **Status Codes**:
  - `201`: Created successfully
  - `400`: Invalid request body
  - `500`: Internal server error

### 2. Get All Events
- **Method**: `GET`
- **Path**: `/api/events`
- **Query Parameters**: None (for now)
- **Response**: Array of event objects
- **Status Codes**:
  - `200`: Success
  - `500`: Internal server error

### 3. Get Event by ID
- **Method**: `GET`
- **Path**: `/api/events/:id`
- **Parameters**: `id` (UUID, required)
- **Response**: Single event object
- **Status Codes**:
  - `200`: Success
  - `404`: Event not found
  - `500`: Internal server error

### 4. Update Event
- **Method**: `PUT`
- **Path**: `/api/events/:id`
- **Parameters**: `id` (UUID, required)
- **Request Body**: Event object (without `id`, `created_at`, `updated_at`)
- **Response**: Updated event object (with new `updated_at`)
- **Status Codes**:
  - `200`: Updated successfully
  - `400`: Invalid request body
  - `404`: Event not found
  - `500`: Internal server error

### 5. Delete Event
- **Method**: `DELETE`
- **Path**: `/api/events/:id`
- **Parameters**: `id` (UUID, required)
- **Response**: Empty body
- **Status Codes**:
  - `204`: Deleted successfully
  - `404`: Event not found
  - `500`: Internal server error

## Error Response Format

```json
{
  "error": "Error message describing what went wrong"
}
```

## Environment Variables

- `DATABASE_URL`: PostgreSQL connection string
- `PORT`: Server port (default: 3000)

## Project Structure
```
/
├── main.go
├── handlers/
│   └── events.go
├── models/
│   └── event.go
├── database/
│   └── connection.go
└── go.mod
```

## Requirements

1. Auto-generate UUID for event IDs using database default
2. Validate required fields: name, venue_name, address, date, time
3. Auto-update `updated_at` timestamp on record updates
4. Return appropriate HTTP status codes
5. Handle database connection errors gracefully
6. Use proper JSON serialization/deserialization
7. Include basic logging for requests
8. Validate email format for contact_email if provided
9. Validate time format (HH:MM:SS) and date format (YYYY-MM-DD)

## Dependencies

```go
// go.mod should include:
github.com/gofiber/fiber/v2
github.com/jmoiron/sqlx
github.com/lib/pq
github.com/google/uuid
```

## Notes

- Keep it simple for the initial implementation
- Focus on core CRUD operations
- Error handling should be consistent across all endpoints
- Database connection should be established once and reused
