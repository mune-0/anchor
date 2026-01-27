package storage

import (
	"bytes"
	"testing"
	"context"
	"time"
)

// Test put and get with context
func TestMemStore_PutGetWithContext (t *testing.T) {
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
func TestMemStore_ContextCancellation (t *testing.T) {
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

// Test context timeout
func TestMemStore_ContextTimeout (t *testing.T) {
	store := NewMemStore()
	defer store.Close()

	// Create context with very short timeout
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()

	// Wait for timeout
	time.Sleep(10 * time.Millisecond)

	err := store.Put(ctx, "key", []byte("value"))
	if err != context.DeadlineExceeded {
		t.Errorf("Expected context.DeadlineExceeded, got %v", err)
	}
}

// Test that valid context works
func TestMemStore_ValidContext (t *testing.T) {
	store := NewMemStore()
	defer store.Close()

	// Context with reasonable timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Should succeed
	err := store.Put(ctx, "key", []byte("value"))
	if err != nil {
		t.Fatalf("Get with valid context failed: %v", err)
	}

	val, err := store.Get(ctx, "key") 
	if err != nil {
		t.Fatalf("Get with valid context failed: %v", err)
	}


	if !bytes.Equal(val, []byte("value")) {
		t.Errorf("Got wrong value")
	}
}
