package service

import "context"

// Database represents a service that interacts with a database.
type Database interface {
	// Health returns a map of health status information.
	// The keys and values in the map are service-specific.
	Health(context.Context) map[string]string

	// Close terminates the database connection.
	// It returns an error if the connection cannot be closed.
	Close() error
}
