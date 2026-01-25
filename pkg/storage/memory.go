package storage

import (
	"sync"
	"context"
	"strings"
)

type MemStore struct {
	kvstore map[string][]byte
	mut sync.RWMutex
	closed bool
}

func  (mem MemStore) Get(ctx context.Context, key string) ([]byte, error) {
	mem.mut.Lock()
	defer mem.mut.Unlock()

	// Checking if memory store is already closed
	if mem.closed {
		return nil, ErrStoreClosed 
	}

	if strings.TrimSpace(key) == "" {
		return nil, ErrInvalidKey
	} 

	if val, ok := mem.kvstore[strings.TrimSpace(key)]; ok {
		return val, nil
	} else {
		return nil, ErrKeyNotFound
	}
}

func (mem MemStore) Put(ctx context.Context, key string, value []byte) error {
	return nil
}
