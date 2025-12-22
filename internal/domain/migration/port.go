package migration

import "context"

// Migrator is the Primary Port (Driver) or Secondary Port (Driven)?
// Here, it's the interface our service will use.
type Migrator interface {
	Up(ctx context.Context) error
	Status(ctx context.Context) error
}
