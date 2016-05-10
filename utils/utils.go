package utils

import (
	"bytes"
	"crypto/rand"
	"io"
	"time"
)

// ReadStream is the read part of a QUIC stream
type ReadStream interface {
	io.Reader
	io.ByteReader
}

// Stream is the interface for QUIC streams
type Stream interface {
	io.Reader
	io.ByteReader
	io.Writer
	io.Closer
}

// ReadUintN reads N bytes
func ReadUintN(b io.ByteReader, length uint8) (uint64, error) {
	var res uint64
	for i := uint8(0); i < length; i++ {
		bt, err := b.ReadByte()
		if err != nil {
			return 0, err
		}
		res ^= uint64(bt) << (i * 8)
	}
	return res, nil
}

// ReadUint64 reads a uint64
func ReadUint64(b io.ByteReader) (uint64, error) {
	var b1, b2, b3, b4, b5, b6, b7, b8 uint8
	var err error
	if b1, err = b.ReadByte(); err != nil {
		return 0, err
	}
	if b2, err = b.ReadByte(); err != nil {
		return 0, err
	}
	if b3, err = b.ReadByte(); err != nil {
		return 0, err
	}
	if b4, err = b.ReadByte(); err != nil {
		return 0, err
	}
	if b5, err = b.ReadByte(); err != nil {
		return 0, err
	}
	if b6, err = b.ReadByte(); err != nil {
		return 0, err
	}
	if b7, err = b.ReadByte(); err != nil {
		return 0, err
	}
	if b8, err = b.ReadByte(); err != nil {
		return 0, err
	}
	return uint64(b1) + uint64(b2)<<8 + uint64(b3)<<16 + uint64(b4)<<24 + uint64(b5)<<32 + uint64(b6)<<40 + uint64(b7)<<48 + uint64(b8)<<56, nil
}

// ReadUint32 reads a uint32
func ReadUint32(b io.ByteReader) (uint32, error) {
	var b1, b2, b3, b4 uint8
	var err error
	if b1, err = b.ReadByte(); err != nil {
		return 0, err
	}
	if b2, err = b.ReadByte(); err != nil {
		return 0, err
	}
	if b3, err = b.ReadByte(); err != nil {
		return 0, err
	}
	if b4, err = b.ReadByte(); err != nil {
		return 0, err
	}
	return uint32(b1) + uint32(b2)<<8 + uint32(b3)<<16 + uint32(b4)<<24, nil
}

// ReadUint16 reads a uint16
func ReadUint16(b io.ByteReader) (uint16, error) {
	var b1, b2 uint8
	var err error
	if b1, err = b.ReadByte(); err != nil {
		return 0, err
	}
	if b2, err = b.ReadByte(); err != nil {
		return 0, err
	}
	return uint16(b1) + uint16(b2)<<8, nil
}

// WriteUint64 writes a uint64
func WriteUint64(b *bytes.Buffer, i uint64) {
	b.WriteByte(uint8(i & 0xff))
	b.WriteByte(uint8((i >> 8) & 0xff))
	b.WriteByte(uint8((i >> 16) & 0xff))
	b.WriteByte(uint8((i >> 24) & 0xff))
	b.WriteByte(uint8((i >> 32) & 0xff))
	b.WriteByte(uint8((i >> 40) & 0xff))
	b.WriteByte(uint8((i >> 48) & 0xff))
	b.WriteByte(uint8(i >> 56))
}

// WriteUint48 writes 48 bit of a uint64
func WriteUint48(b *bytes.Buffer, i uint64) {
	b.WriteByte(uint8(i & 0xff))
	b.WriteByte(uint8((i >> 8) & 0xff))
	b.WriteByte(uint8((i >> 16) & 0xff))
	b.WriteByte(uint8((i >> 24) & 0xff))
	b.WriteByte(uint8((i >> 32) & 0xff))
	b.WriteByte(uint8((i >> 40)))
}

// WriteUint32 writes a uint32
func WriteUint32(b *bytes.Buffer, i uint32) {
	b.WriteByte(uint8(i & 0xff))
	b.WriteByte(uint8((i >> 8) & 0xff))
	b.WriteByte(uint8((i >> 16) & 0xff))
	b.WriteByte(uint8((i >> 24) & 0xff))
}

// WriteUint24 writes 24 bit of a uint32
func WriteUint24(b *bytes.Buffer, i uint32) {
	b.WriteByte(uint8(i & 0xff))
	b.WriteByte(uint8((i >> 8) & 0xff))
	b.WriteByte(uint8((i >> 16) & 0xff))
}

// WriteUint16 writes a uint16
func WriteUint16(b *bytes.Buffer, i uint16) {
	b.WriteByte(uint8(i & 0xff))
	b.WriteByte(uint8((i >> 8) & 0xff))
}

// Max returns the maximum of two Ints
func Max(a, b int) int {
	if a < b {
		return b
	}
	return a
}

// MaxUint32 returns the maximum of two uint32
func MaxUint32(a, b uint32) uint32 {
	if a < b {
		return b
	}
	return a
}

// MaxUint64 returns the maximum of two uint64
func MaxUint64(a, b uint64) uint64 {
	if a < b {
		return b
	}
	return a
}

// Min returns the minimum of two Ints
func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// MinInt64 returns the minimum of two int64
func MinInt64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

// MaxInt64 returns the minimum of two int64
func MaxInt64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

// MaxDuration returns the max duration
func MaxDuration(a, b time.Duration) time.Duration {
	if a > b {
		return a
	}
	return b
}

// AbsDuration returns the absolute value of a time duration
func AbsDuration(d time.Duration) time.Duration {
	if d >= 0 {
		return d
	}
	return -d
}

// RandomBit returns a cryptographically secure random bit (encoded as true / false)
func RandomBit() (bool, error) {
	// ToDo: it's probably more efficient to read a bigger slice of random numbers at once and to cache them somewhere
	b := make([]byte, 1)
	_, err := rand.Read(b)
	if err != nil {
		return false, err
	}
	if uint8(b[0])%2 == 0 {
		return false, nil
	}
	return true, nil
}

// Uint32Slice attaches the methods of sort.Interface to []uint32, sorting in increasing order.
type Uint32Slice []uint32

func (s Uint32Slice) Len() int           { return len(s) }
func (s Uint32Slice) Less(i, j int) bool { return s[i] < s[j] }
func (s Uint32Slice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
