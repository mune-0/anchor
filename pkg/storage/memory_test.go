package storage

import (
	"bytes"
	"testing"
	"context"
)

// Test put and get with context
func TestMemStore_PutGetWithContext(t *testing.T) {
	store := NewMemStore()
	defer store.Close()

	ctx := context.Background()

	key := "test-key"
	value := []byte("test-value")

	err := store.Put(ctx, key, value)
	if err != nil {
		t.Fatalf("Put failed: %v", err)
	}

	res, err := store.Get(ctx, key)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if !bytes.Equal(res, value) {
		t.Errorf("Got %v, expected %v", res, value)
	}
}

// Test context cancellation
func TestMemStore_ContextCancellation(t *testing.T) {
	store := NewMemStore()
	defer store.Close()

	// Create cancelled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately

	// Operations should fail with context.Cancelled
	err := store.Put(ctx, "key", []byte("value"))
	if err != context.Canceled {
		t.Errorf("Expected context.Canceled, got %v", err)
	}

	_, err = store.Get(ctx, "key")
	if err != context.Canceled {
		t.Errorf("Expected context.Canceled, got %v", err)
	}
}
