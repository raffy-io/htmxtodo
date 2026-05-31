package connection

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/tursodatabase/libsql-client-go/libsql" // Ensure the driver is blank imported
)

// Connect initializes and returns a long-lived database connection pool.
func Connect(dbURL string, authToken string) (*sql.DB, error) {
	if dbURL == "" || authToken == "" {
		return nil, fmt.Errorf("TURSO_DATABASE_URL and TURSO_AUTH_TOKEN must be set")
	}

	dataSourceName := fmt.Sprintf("%s?authToken=%s", dbURL, authToken)

	conn, err := sql.Open("libsql", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Configure pool settings
	conn.SetMaxOpenConns(10)
	conn.SetMaxIdleConns(5)
	conn.SetConnMaxLifetime(time.Hour)

	// Verify connection
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := conn.PingContext(ctx); err != nil {
		conn.Close() // Clean up if ping fails
		return nil, fmt.Errorf("database unreachable: %w", err)
	}

	return conn, nil
}