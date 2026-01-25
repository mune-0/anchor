package storage

import "context"

// Store defines the standard behavior for a Key-Value storage engine.
type Store interface {
	// Put inserts or updates the value associated with the given key.
	// Returns ErrInvalidKey if the key is empty.
	// Returns ErrStoreClosed if the store is no longer active.
	Put (ctx context.Context, key string, value []byte) error

	// Get retrieves the value associated with the given key.
	// Returns the value and nil on success.
	// Returns nil and ErrKeyNotFound if the key does not exist.
	// Returns nil and ErrStoreClosed if the store is no longer active.
	Get (ctx context.Context, key string) ([]byte, error)

	// Delete removes the value associated with the given key.
	// If the key does not exist, Delete should return nil (idempotent behavior).
	// Returns ErrStoreClosed if the store is no longer active.
	Delete (ctx context.Context, key string) error

	// Close gracefully shuts down the store, flushing any pending writes.
	// After Close is called, all other methods should return ErrStoreClosed.
	Close(ctx context.Context) error
}
