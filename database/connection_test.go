package database

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConnectWithoutDatabaseURL(t *testing.T) {
	// Save original env var
	originalURL := os.Getenv("DATABASE_URL")
	defer func() {
		if originalURL != "" {
			os.Setenv("DATABASE_URL", originalURL)
		} else {
			os.Unsetenv("DATABASE_URL")
		}
	}()

	// Remove DATABASE_URL env var
	os.Unsetenv("DATABASE_URL")

	err := Connect()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "DATABASE_URL environment variable is required")
}

func TestConnectWithInvalidDatabaseURL(t *testing.T) {
	// Save original env var
	originalURL := os.Getenv("DATABASE_URL")
	defer func() {
		if originalURL != "" {
			os.Setenv("DATABASE_URL", originalURL)
		} else {
			os.Unsetenv("DATABASE_URL")
		}
	}()

	// Set invalid DATABASE_URL
	os.Setenv("DATABASE_URL", "invalid-url")

	err := Connect()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to connect to database")
}

func TestConnectWithValidDatabaseURL(t *testing.T) {
	// Skip this test if no DATABASE_URL is set in environment
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		t.Skip("DATABASE_URL not set, skipping integration test")
	}

	// Reset DB to nil to test fresh connection
	if DB != nil {
		DB.Close()
		DB = nil
	}

	err := Connect()
	require.NoError(t, err)
	assert.NotNil(t, DB)

	// Test that we can ping the database
	err = DB.Ping()
	assert.NoError(t, err)
}

func TestCreateTablesWithoutConnection(t *testing.T) {
	// Save original DB
	originalDB := DB
	defer func() {
		DB = originalDB
	}()

	// Set DB to nil to simulate no connection
	DB = nil

	// This will panic with nil pointer dereference, which is expected behavior
	// We test that it panics as expected
	assert.Panics(t, func() {
		CreateTables()
	}, "CreateTables should panic when DB is nil")
}

func TestCreateTablesWithValidConnection(t *testing.T) {
	// Skip this test if no DATABASE_URL is set in environment
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		t.Skip("DATABASE_URL not set, skipping integration test")
	}

	// Ensure we have a connection
	if DB == nil {
		err := Connect()
		require.NoError(t, err)
	}

	err := CreateTables()
	assert.NoError(t, err)

	// Verify table exists by querying it
	var count int
	err = DB.Get(&count, "SELECT COUNT(*) FROM events")
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, count, 0)
}

func TestCloseWithNilDB(t *testing.T) {
	// Save original DB
	originalDB := DB
	defer func() {
		DB = originalDB
	}()

	// Set DB to nil
	DB = nil

	err := Close()
	assert.NoError(t, err)
}

func TestCloseWithValidDB(t *testing.T) {
	// Skip this test if no DATABASE_URL is set in environment
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		t.Skip("DATABASE_URL not set, skipping integration test")
	}

	// Create a separate connection for this test
	err := Connect()
	require.NoError(t, err)
	require.NotNil(t, DB)

	err = Close()
	assert.NoError(t, err)

	// After closing, ping should fail
	if DB != nil {
		err = DB.Ping()
		assert.Error(t, err)
	}
}

// TestDatabaseSchemaIntegration tests the actual database schema creation
func TestDatabaseSchemaIntegration(t *testing.T) {
	// Skip this test if no DATABASE_URL is set in environment
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		t.Skip("DATABASE_URL not set, skipping integration test")
	}

	// Ensure we have a connection
	if DB == nil {
		err := Connect()
		require.NoError(t, err)
	}

	// Create tables
	err := CreateTables()
	require.NoError(t, err)

	// Test that we can insert and retrieve a test record
	testEventID := "550e8400-e29b-41d4-a716-446655440000"
	
	// Clean up any existing test data
	DB.Exec("DELETE FROM events WHERE id = $1", testEventID)

	// Insert test event
	_, err = DB.Exec(`
		INSERT INTO events (id, name, venue_name, address, date, time) 
		VALUES ($1, $2, $3, $4, $5, $6)`,
		testEventID, "Test Event", "Test Venue", "Test Address", "2024-03-15", "14:30:00")
	require.NoError(t, err)

	// Retrieve test event
	var name, venueName, address, date, time string
	err = DB.QueryRow("SELECT name, venue_name, address, date, time FROM events WHERE id = $1", testEventID).
		Scan(&name, &venueName, &address, &date, &time)
	require.NoError(t, err)

	assert.Equal(t, "Test Event", name)
	assert.Equal(t, "Test Venue", venueName)
	assert.Equal(t, "Test Address", address)
	assert.Equal(t, "2024-03-15", date)
	assert.Equal(t, "14:30:00", time)

	// Clean up test data
	_, err = DB.Exec("DELETE FROM events WHERE id = $1", testEventID)
	require.NoError(t, err)
}