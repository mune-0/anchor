package storage

import (
	"bytes"
	"testing"
	"context"
	"time"
	"fmt"
	"sync"
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

// Test Get non-existent key
func TestMemStore_GetNotFound (t *testing.T) {
	store := NewMemStore()
	defer store.Close()

	ctx := context.Background()

	_, err := store.Get(ctx, "non-existent")
	if err != ErrKeyNotFound {
		t.Errorf("Expected ErrKeyNotFound, got %v", err)
	}
}

// Test Put with empty key
func TestMemStore_InvalidKey(t *testing.T) {
    store := NewMemStore()
    defer store.Close()
    
    ctx := context.Background()

    err := store.Put(ctx, "", []byte("value"))
    if err != ErrInvalidKey {
        t.Errorf("Expected ErrInvalidKey, got %v", err)
    }
}

// Test Delete operation
func TestMemStore_Delete(t *testing.T) {
    store := NewMemStore()
    defer store.Close()

    ctx := context.Background()

    key := "test-key"
    value := []byte("test-value")

    // Put then delete
    store.Put(ctx, key, value)
    err := store.Delete(ctx, key)
    if err != nil {
        t.Fatalf("Delete failed: %v", err)
    }

    // Verify it's gone
    _, err = store.Get(ctx, key)
    if err != ErrKeyNotFound {
        t.Errorf("Key should not exist after delete")
    }
}

// Test Update existing key
func TestMemStore_Update(t *testing.T) {
    store := NewMemStore()
    defer store.Close()

    ctx := context.Background()

    key := "test-key"
    value1 := []byte("first-value")
    value2 := []byte("second-value")

    // Put initial value
    store.Put(ctx, key, value1)

    // Update with new value
    store.Put(ctx, key, value2)

    // Verify we get the new value
    got, _ := store.Get(ctx, key)
    if !bytes.Equal(got, value2) {
        t.Errorf("Got %v, want %v", got, value2)
    }
}

// Test Defensive copying test - CRITICAL
func TestMemStore_DefensiveCopy(t *testing.T) {
    store := NewMemStore()
    defer store.Close()

    ctx := context.Background()

    key := "test-key"
    value := []byte("original")

    // Put value
    store.Put(ctx, key, value)

    // Modify the original slice
    value[0] = 'X'

    // Get should return original value, not modified
    got, _ := store.Get(ctx, key)
    if bytes.Equal(got, value) {
        t.Error("Store did not make defensive copy on Put!")
    }
    if !bytes.Equal(got, []byte("original")) {
        t.Error("Stored value was corrupted")
    }

    // Now modify what Get returned
    got[0] = 'Y'

    // Get again - should still be original
    got2, _ := store.Get(ctx, key)
    if !bytes.Equal(got2, []byte("original")) {
        t.Error("Store did not make defensive copy on Get!")
    }
}

// Test Concurrent reads (should work with RWMutex)
func TestMemStore_ConcurrentReads(t *testing.T) {
    store := NewMemStore()
    defer store.Close()

    ctx := context.Background()

    // Put some data
    for i := range 10 {
        key := fmt.Sprintf("key-%d", i)
        value := []byte(fmt.Sprintf("value-%d", i))
        store.Put(ctx, key, value)
    }

    // Launch 10 concurrent readers
    var wg sync.WaitGroup
    errors := make(chan error, 10)

    for i := range 10 {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()


            // Each goroutine reads all keys
	    for j := range 10 {
                key := fmt.Sprintf("key-%d", j)
                _, err := store.Get(ctx, key)
                if err != nil {
                    errors <- err
                    return
                }
            }
        }(i)
    }

    wg.Wait()
    close(errors)

    // Check for errors
    for err := range errors {
        t.Errorf("Concurrent read failed: %v", err)
    }
}


// Test Concurrent writes (should be safe with mutex)
func TestMemStore_ConcurrentWrites(t *testing.T) {
    store := NewMemStore()
    defer store.Close()
 
    ctx := context.Background()

    var wg sync.WaitGroup
    numGoroutines := 10
    writesPerGoroutine := 100
    
    // Each goroutine writes different keys
    for i := range numGoroutines {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            
            for j := range writesPerGoroutine {
                key := fmt.Sprintf("key-%d-%d", id, j)
                value := []byte(fmt.Sprintf("value-%d-%d", id, j))
		err := store.Put(ctx, key, value)
                if err != nil {
                    t.Errorf("Put failed: %v", err)
                }
            }
        }(i)
    }
    
    wg.Wait()
    
    // Verify all writes succeeded
    for i := range numGoroutines {
        for j := range writesPerGoroutine {
            key := fmt.Sprintf("key-%d-%d", i, j)
            expected := []byte(fmt.Sprintf("value-%d-%d", i, j))
            
            got, err := store.Get(ctx, key)
            if err != nil {
                t.Errorf("Get failed for %s: %v", key, err)
            }
            if !bytes.Equal(got, expected) {
                t.Errorf("Key %s: got %v, want %v", key, got, expected)
            }
        }
    }
}

// Test Close should prevent further operations
func TestMemStore_CloseStore(t *testing.T) {
    store := NewMemStore()

    ctx := context.Background()

    // Put some data
    store.Put(ctx, "key", []byte("value"))

    // Close the store
    err := store.Close()
    if err != nil {
        t.Fatalf("Close failed: %v", err)
    }

    // Operations after close should fail
    err = store.Put(ctx, "key2", []byte("value2"))
    if err != ErrStoreClosed {
        t.Errorf("Put after close should return ErrStoreClosed, got %v", err)
    }

    _, err = store.Get(ctx, "key")
    if err != ErrStoreClosed {
        t.Errorf("Get after close should return ErrStoreClosed, got %v", err)
    }
}





