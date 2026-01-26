package storage

import (
	"sync"
	"context"
	"strings"
)

// MemStore is an in-memory implementation of Store
type MemStore struct {
	data map[string][]byte
	mut sync.RWMutex
	closed bool
}

// NewMemStore creates an new in-memory store
func NewMemStore () *MemStore {
	return &MemStore{
		data: make(map[string][]byte),
	}
}

// Get returns a value by key
func  (mem *MemStore) Get (ctx context.Context, key string) ([]byte, error) {
	if strings.TrimSpace(key) == "" {
		return nil, ErrInvalidKey
	}

	mem.mut.Lock()
	defer mem.mut.Unlock()

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
	if strings.TrimSpace(key) == "" {
		return ErrInvalidKey
	}

	mem.mut.RLock()
	defer mem.mut.RUnlock()
	
	if mem.closed {
		return ErrStoreClosed
	}
	
	// Defensive copy
	snapshot := make([]byte, len(value))
	copy(snapshot, value)
	mem.data[key] = snapshot
	return nil
}

// Delete removes a key
func (mem *MemStore) Delete(ctx context.Context, key string) error {
	if strings.TrimSpace(key) == "" {
		return ErrInvalidKey
	}

	mem.mut.Lock()
	defer mem.mut.Unlock()

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
