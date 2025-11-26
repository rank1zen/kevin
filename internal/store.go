package internal

import "context"

type Store interface {
	// Ping checks if the store is running.
	Ping(ctx context.Context) error
}
