package wal

import (
	"bufio"
	"os"
	"sync"
	"context"
)

type Writer struct {
	file *os.File
	writer *bufio.Writer
	mut sync.RWMutex
}

func NewWriter(path string) (*Writer, error) {
	// Open for appending, create if missing
	f, err := os.OpenFile(path, os.O_APPEND | os.O_CREATE | os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	return &Writer {
		file: f,
		writer: bufio.NewWriterSize(f, 64*1024), // 64KB buffer
	}, nil
}

// Write appends an entry to the buffer
func (w *Writer) Write(ctx context.Context, entry *LogEntry) error {
	// Check context before acquiring lock
	if err := ctx.Err(); err != nil {
		return  err
	}

	w.mut.Lock()
	defer w.mut.Unlock()

	// Check context after acquiring lock
	if err := ctx.Err(); err != nil {
		return err
	}

	data := entry.Encode()
	_, err := w.writer.Write(data)
	return err
}

func (w *Writer) SyncWrite(ctx context.Context, entry *LogEntry) error {
	// Checking context before acquiring lock
	if err := ctx.Err(); err != nil {
		return err
	}

	w.mut.Lock()
	defer w.mut.Unlock()

	// Checking context after acquiring lock
	if err := ctx.Err(); err != nil {
		return err
	}

	data := entry.Encode()

	// Write to the OS buffer
	if _, err := w.writer.Write(data); err != nil {
		return err
	}

	// 1. Flush bufio to the OS
	if err := w.writer.Flush(); err != nil {
		return err
	}

	// 2. Fsync forces the disk controller to commit to physical media
	// This is the "Durability" in ACID.
	return w.file.Sync()
}

// Sync flushes the buffer to the OS and forces a disk write
func (w *Writer) Sync() error {
	if err := w.writer.Flush(); err != nil {
		return err
	}
	return w.file.Sync()
}

func (w *Writer) Close() error {
	w.Sync()
	return w.file.Close()
}
