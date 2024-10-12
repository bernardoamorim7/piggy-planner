package database

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"log"
	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"
)

//go:embed migrations/**/*.sql
var migrations embed.FS

// DbService represents a service that interacts with a database.
type DbService interface {
	// Health returns a map of health status information.
	// The keys and values in the map are service-specific.
	Health() map[string]string

	// Close terminates the database connection.
	// It returns an error if the connection cannot be closed.
	Close() error

	// QueryRow executes a query that is expected to return at most one row.
	QueryRow(query string, args ...any) *sql.Row

	// Query executes a query that returns rows, typically a SELECT.
	Query(query string, args ...any) (*sql.Rows, error)

	// Exec executes a query without returning any rows.
	Exec(query string, args ...any) (sql.Result, error)

	// Prepare creates a prepared statement for later queries or executions.
	Prepare(query string) (*sql.Stmt, error)

	// Begin starts a new transaction.
	Begin() (*sql.Tx, error)
}

type dbService struct {
	db *sql.DB
}

var (
	dbInstance *dbService
)

func New() (DbService, error) {
	// Reuse Connection
	if dbInstance != nil {
		return dbInstance, nil
	}

	db, err := sql.Open("sqlite3", "./piggy_planner.db")
	if err != nil {
		return nil, err
	}

	// Run migrations
	runMigrations(db)

	dbInstance = &dbService{db: db}
	return dbInstance, nil
}

// runMigrations runs the database migrations using the goose package.
func runMigrations(db *sql.DB) {
	goose.SetBaseFS(migrations)
	err := goose.SetDialect("sqlite3")
	if err != nil {
		log.Fatalf("Failed to set dialect: %v", err)
	}

	// Debugging: List embedded files
	entries, err := migrations.ReadDir("migrations")
	if err != nil {
		log.Fatalf("Failed to read embedded migrations: %v", err)
	}

	for _, entry := range entries {
		log.Printf("Found migration directory: %s", entry.Name())

		subEntries, err := migrations.ReadDir("migrations/" + entry.Name())
		if err != nil {
			log.Fatalf("Failed to read embedded migrations in %s: %v", entry.Name(), err)
		}

		for _, subEntry := range subEntries {
			log.Printf("Found migration file: %s/%s", entry.Name(), subEntry.Name())
		}

		migrationPath := "migrations/" + entry.Name()
		log.Printf("Running migrations in directory: %s", migrationPath)

		if err := goose.Up(db, migrationPath); err != nil {
			log.Fatalf("Failed to run migrations in %s: %v", migrationPath, err)
		}
	}
}

// Health checks the health of the database connection by pinging the database.
// It returns a map with keys indicating various health statistics.
func (s *dbService) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	stats := make(map[string]string)

	// Ping the database
	err := s.db.PingContext(ctx)
	if err != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("db down: %v", err)
		log.Fatalf("db down: %v", err) // Log the error and terminate the program
		return stats
	}

	// Database is up, add more statistics
	stats["status"] = "up"
	stats["message"] = "It's healthy"

	// Get database stats (like open connections, in use, idle, etc.)
	dbStats := s.db.Stats()
	stats["open_connections"] = strconv.Itoa(dbStats.OpenConnections)
	stats["in_use"] = strconv.Itoa(dbStats.InUse)
	stats["idle"] = strconv.Itoa(dbStats.Idle)
	stats["wait_count"] = strconv.FormatInt(dbStats.WaitCount, 10)
	stats["wait_duration"] = dbStats.WaitDuration.String()
	stats["max_idle_closed"] = strconv.FormatInt(dbStats.MaxIdleClosed, 10)
	stats["max_lifetime_closed"] = strconv.FormatInt(dbStats.MaxLifetimeClosed, 10)

	// Evaluate stats to provide a health message
	if dbStats.OpenConnections > 40 { // Assuming 50 is the max for this example
		stats["message"] = "The database is experiencing heavy load."
	}
	if dbStats.WaitCount > 1000 {
		stats["message"] = "The database has a high number of wait events, indicating potential bottlenecks."
	}

	if dbStats.MaxIdleClosed > int64(dbStats.OpenConnections)/2 {
		stats["message"] = "Many idle connections are being closed, consider revising the connection pool settings."
	}

	if dbStats.MaxLifetimeClosed > int64(dbStats.OpenConnections)/2 {
		stats["message"] = "Many connections are being closed due to max lifetime, consider increasing max lifetime or revising the connection usage pattern."
	}

	return stats
}

// Close closes the database connection.
// It logs a message indicating the disconnection from the specific database.
// If the connection is successfully closed, it returns nil.
// If an error occurs while closing the connection, it returns the error.
func (s *dbService) Close() error {
	log.Printf("Disconnected from database")
	return s.db.Close()
}

// QueryRow executes a query that is expected to return at most one row.
// QueryRow always returns a non-nil value. Errors are deferred until Row's Scan method is called.
func (s *dbService) QueryRow(query string, args ...any) *sql.Row {
	return s.db.QueryRow(query, args...)
}

// Query executes a query that returns rows, typically a SELECT.
// The args are for any placeholder parameters in the query.
func (s *dbService) Query(query string, args ...any) (*sql.Rows, error) {
	return s.db.Query(query, args...)
}

// Exec executes a query without returning any rows.
// The args are for any placeholder parameters in the query.
func (s *dbService) Exec(query string, args ...any) (sql.Result, error) {
	return s.db.Exec(query, args...)
}

// Prepare creates a prepared statement for later queries or executions.
// Multiple queries or executions may be run concurrently from the returned statement.
func (s *dbService) Prepare(query string) (*sql.Stmt, error) {
	return s.db.Prepare(query)
}

// Begin starts a new transaction.
// The default isolation level is dependent on the driver.
func (s *dbService) Begin() (*sql.Tx, error) {
	return s.db.Begin()
}
