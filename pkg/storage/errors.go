package storage

import "errors"

var (
	// ErrKeyNotFound is returned when provided key does not exist
	ErrKeyNotFound = errors.New("key not found")

	// ErrStoreClosed is returned when an operation is attempted on a closed store engine
	ErrStoreClosed = errors.New("store engine is closed")

	// ErrInvalidKey is returned when provided key is malformed or empty
	ErrInvalidKey = errors.New("key is malformed or empty")
)
