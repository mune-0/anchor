package wal

import (
	"encoding/binary"
	"errors"
	"fmt"
	"hash/crc32"
	"io"
	"os"
)

var ErrCorruption = errors.New("wal: data corruption detected (checksum mismatch)")

type Reader struct {
	file *os.File
}

func NewReader(path string) (*Reader, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return &Reader{file: f}, nil
}

// Next reads the next entry from the log. Returns io.EOF at end of file.
func (r *Reader) Next() (*LogEntry, error) {
	// 1. Read the Fixed Header
	headerBuf := make([]byte, HeaderSize)
	if _, err := io.ReadFull(r.file, headerBuf); err != nil {
		return nil, err // Returns io.EOF or io.ErrUnexpectedEOF
	}

	// 2. Parse Header
	expectedCRC := binary.LittleEndian.Uint32(headerBuf[0:4])
	ts, op, kLen, vLen := DecodeHeader(headerBuf)
	
	// 3. Read Variable Data (Key + Value)
	payloadSize := int(kLen + vLen)
	payloadBuf := make([]byte, payloadSize)
	if _, err := io.ReadFull(r.file, payloadBuf); err != nil {
		return nil, err
	}

	// 4. Verify Integrity
	// Checksum was calculated on everything AFTER the CRC field
	h := crc32.NewIEEE()
	h.Write(headerBuf[4:]) // Header minus CRC
	h.Write(payloadBuf) // Key + Value

	if h.Sum32() != expectedCRC {
		return nil, fmt.Errorf("%w: at offset %d", ErrCorruption, r.CurrentOffset())
	}

	return &LogEntry {
		Checksum: expectedCRC,
		Timestamp: ts,
		Op: op,
		Key: payloadBuff[:kLen]
		Value: payloadBuf[kLen:],
	}, nil
}

func (r *Reader) CurrentOffset() int64 {
	offset, _ := r.file.Seek(0, io.SeekCurrent)
	return offset
}

func (r *Reader) Close() error {
	return r.file.Close()
}
