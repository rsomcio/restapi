// feature1
// feature2

package database

import (
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
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

var DB *sqlx.DB

func Connect() error {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return fmt.Errorf("DATABASE_URL environment variable is required")
	}

	var err error
	DB, err = sqlx.Connect("postgres", databaseURL)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := DB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	logInfo("Successfully connected to database")
	return nil
}

func CreateTables() error {
	schema := `
	CREATE TABLE IF NOT EXISTS events (
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
	);`

	_, err := DB.Exec(schema)
	if err != nil {
		return fmt.Errorf("failed to create tables: %w", err)
	}

	logInfo("Database tables created successfully")
	return nil
}

func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}
