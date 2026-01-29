package wal

import (
	"encoding/binary"
	"hash/crc32"
)

const (
	// HeaderSize = 4 (CRC) + 8 (TS) + 1 (OP) + 4 (KLen) + 4 (VLen)
	HeaderSize = 21
)

type OpType uint8

const (
	OpPut OpType = 0
	OpDelete OpType = 1
)

// LogEntry represents a single record in the WAL
type LogEntry struct {
	Checksum uint32
	Timestamp int64
	Op OpType
	Key []byte
	Value []byte
}

// Encode serializes the entry into a byte slice
func (e *LogEntry) Encode() []byte {
	buf := make([]byte, HeaderSize+len(e.Key) + len(e.Value))

	// Leave space for Checksum at buf[0:4]
	binary.LittleEndian.PutUint64(buf[4:12], uint64(e.Timestamp))
	buf[12] = uint8(e.Op)
	binary.LittleEndian.PutUint32(buf[13:17], uint32(len(e.Key)))
	binary.LittleEndian.PutUint32(buf[17:21], uint32(len(e.Value)))

	copy(buf[HeaderSize:], e.Key)
	copy(buf[HeaderSize+len(e.Key):], e.Value)

	// Calculate checksum of everything except the checksum field itself
	e.Checksum = crc32.ChecksumIEEE(buf[4:])
	binary.LittleEndian.PutUint32(buf[0:4], e.Checksum)

	return buf
}

// DecodeHeader parses the fixed-size portion to help manage memory allocation
func DecodeHeader(data []byte) (timestamp int64, op OpType, kLen, vLen uint32) {
	timestamp = int64(binary.LittleEndian.Uint64(data[4:12]))
	op = OpType(data[12])
	kLen = binary.LittleEndian.Uint32(data[13:17])
	vLen = binary.LittleEndian.Uint32(data[17:21])
	return
}
