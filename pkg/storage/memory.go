package storage

import (
	"sync"
	"context"
	"strings"
	"time"
	"fmt"
	"com.github/mune-0/anchor/pkg/wal"
)

// MemStore is an in-memory implementation of Store
type MemStore struct {
	data map[string][]byte
	mut sync.RWMutex
	closed bool
	walWriter wal.WALWriter
}

// NewMemStore creates an new in-memory store
func NewMemStore (w wal.WALWriter) *MemStore {
	return &MemStore{
		data: make(map[string][]byte),
		walWriter : w,
	}
}

// Get returns a value by key
func  (mem *MemStore) Get (ctx context.Context, key string) ([]byte, error) {
	// Check context before acquiring lock
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	if strings.TrimSpace(key) == "" {
		return nil, ErrInvalidKey
	}

	mem.mut.RLock()
	defer mem.mut.RUnlock()

	if err := ctx.Err(); err != nil {
		return nil, err
	}

	// Checking if memory store is already closed
	if mem.closed {
		return nil, ErrStoreClosed 
	}

	if val, ok := mem.data[strings.TrimSpace(key)]; ok {
		// Defensive copy
		snapshot := make([]byte, len(val))
		copy(snapshot, val)
		return snapshot, nil
	} else {
		return nil, ErrKeyNotFound
	}
}

// Put stores a key-value pair
func (mem *MemStore) Put (ctx context.Context, key string, value []byte) error {
	// Check context before acquiring lock
	if err := ctx.Err(); err != nil {
		return err // Returns context.Canceled or context.DeadlineExceeded
	}

	if strings.TrimSpace(key) == "" {
		return ErrInvalidKey
	}

	// Defensive copy 
	snapshot := make([]byte, len(value))
	copy(snapshot, value)

	entry := &wal.LogEntry{
		Timestamp: time.Now().UnixNano(),
		Op: wal.OpPut,
		Key: []byte(key),
		Value: snapshot,
	}

	// Write to WAL (Durability)
	if err := mem.walWriter.SyncWrite(ctx, entry); err != nil {
		return fmt.Errorf("WAL failure (data safe, update aborted): %w", err)
	}

	mem.mut.Lock()
	defer mem.mut.Unlock()

	// Check context again after acquiring lock
	if err := ctx.Err(); err != nil {
		return err
	}
	
	if mem.closed {
		return ErrStoreClosed
	}
	
	mem.data[key] = snapshot

	return nil
}

// Delete removes a key
func (mem *MemStore) Delete(ctx context.Context, key string) error {
	// Check context before acquiring lock
	if err := ctx.Err(); err != nil {
		return err
	}

	if strings.TrimSpace(key) == "" {
		return ErrInvalidKey
	}

	mem.mut.Lock()
	defer mem.mut.Unlock()


	// Check context again after acquiring lock
	if err := ctx.Err(); err != nil {
		return err
	}

	delete(mem.data, key)
	return nil
}

// Close closes the store
func (mem *MemStore) Close () error {
	mem.mut.Lock()
	defer mem.mut.Unlock()

	if mem.closed {
		return ErrStoreClosed
	}

	mem.closed = true
	mem.data = nil
	return nil
}
