package wal

import (
	"context"
)

// WAL File Writer interface
type WALWriter interface {
	// Log entry is not flushed to the file immediately
	// File is not always in sync
	Write(ctx context.Context, entry *LogEntry) error

	// Log entry directly to the file
	// File is always in sync
	SyncWrite(ctx context.Context, entry *LogEntry) error
}
