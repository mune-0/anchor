package wal

import (
	"os"
	"testing"
	"time"
	"context"
)

// Tests that a LogEntry can be successfully be encoded, written to disk, and retrieved with all fields intact
func TestWAL_WriteRead(t *testing.T) {
	tmpFile, _ := os.CreateTemp("", "wal_test_*.log")
	defer os.Remove(tmpFile.Name())

	// 1. Initialize a Writer and write an entry
	writer, _ := NewWriter(tmpFile.Name())
	originalEntry := &LogEntry{
		Timestamp: time.Now().UnixNano(),
		Op: OpPut,
		Key: []byte("user:101"),
		Value: []byte("active"),
	}

	ctx := context.Background()

	// Using SyncWrite to ensure durability
	if err := writer.SyncWrite(ctx, originalEntry); err != nil {
		t.Fatalf("Failed to write: %v", err)
	}
	writer.Close()

	// 2. Initialize Reader and read it back
	reader, _ := NewReader(tmpFile.Name())
	recovered, err := reader.Next()
	if err != nil {
		t.Fatalf("Failed to read: %v", err)
	}

	// Verify fields
	if string(recovered.Key) != string(originalEntry.Key) {
		t.Errorf("Expected key %s, got %s", originalEntry.Key, recovered.Key)
	}
	
	if recovered.Op != originalEntry.Op {
		t.Errorf("Expected op %v, got %v", originalEntry.Op, recovered.Op)
	}
}

// Tests that the Reader correctly identifies if data has been tampered with or corrupted after being written
func TestWAL_ChecksumValidation(t *testing.T) {
	tmpFile, _ := os.CreateTemp("", "wal_corrupt_*.log")
	defer os.Remove(tmpFile.Name())

	ctx := context.Background()

	writer, _ := NewWriter(tmpFile.Name())
	writer.SyncWrite(ctx, &LogEntry{Key: []byte("secure"), Value: []byte("data")})
	writer.Close()

	// Manually flip a bit in the file to simulate corruption
	data, _ := os.ReadFile(tmpFile.Name())
	data[len(data)-1] ^= 0xFF // Flip the last bit of the value
	os.WriteFile(tmpFile.Name(), data, 0644)

	reader, _ := NewReader(tmpFile.Name())
	_, err := reader.Next()

	if err == nil {
		t.Error("Expected error due to checksum mismatch, but got nil")
	}
}

// Tests that the library does not crash (e.g. In the event of a power loss, the file might end abruptly)
func TestWAL_PartialWrite(t *testing.T) {
	tmpFile, _ := os.CreateTemp("", "wal_partial_*.log")
	defer os.Remove(tmpFile.Name())

	ctx := context.Background()

	writer, _ := NewWriter(tmpFile.Name())
	entry := &LogEntry{Key: []byte("partial"), Value: []byte("write_test")}
	writer.SyncWrite(ctx, entry)
	writer.Close()

	// Truncate the file to simulate a crash during the middle of a write
	info, _ := os.Stat(tmpFile.Name())
	os.Truncate(tmpFile.Name(), info.Size()-5)

	reader, _ := NewReader(tmpFile.Name())
	_, err := reader.Next()

	// We expect an EOF or UnexpectedEOF rather than a panic or valid entry
	if err == nil {
		t.Error("Reader should have failed to read a truncated entry")
	}
}
