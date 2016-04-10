package quic

import (
	"bytes"
	"io"
)

func readUintN(b io.ByteReader, length uint8) (uint64, error) {
	var res uint64
	for i := uint8(0); i < length; i++ {
		bt, err := b.ReadByte()
		if err != nil {
			return 0, err
		}
		res = res<<8 + uint64(bt)
	}
	return res, nil
}

func readUint32(b io.ByteReader) (uint32, error) {
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

func readUint16(b io.ByteReader) (uint16, error) {
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

func writeUint64(b *bytes.Buffer, i uint64) {
	b.WriteByte(uint8(i & 0xff))
	b.WriteByte(uint8((i >> 8) & 0xff))
	b.WriteByte(uint8((i >> 16) & 0xff))
	b.WriteByte(uint8((i >> 24) & 0xff))
	b.WriteByte(uint8((i >> 32) & 0xff))
	b.WriteByte(uint8((i >> 40) & 0xff))
	b.WriteByte(uint8((i >> 48) & 0xff))
	b.WriteByte(uint8((i >> 56) & 0xff))
}

func writeUint32(b *bytes.Buffer, i uint32) {
	b.WriteByte(uint8(i & 0xff))
	b.WriteByte(uint8((i >> 8) & 0xff))
	b.WriteByte(uint8((i >> 16) & 0xff))
	b.WriteByte(uint8((i >> 24) & 0xff))
}

func writeUint16(b *bytes.Buffer, i uint16) {
	b.WriteByte(uint8(i & 0xff))
	b.WriteByte(uint8((i >> 8) & 0xff))
}

// Uint32Slice attaches the methods of sort.Interface to []uint32, sorting in increasing order.
type Uint32Slice []uint32

func (s Uint32Slice) Len() int           { return len(s) }
func (s Uint32Slice) Less(i, j int) bool { return s[i] < s[j] }
func (s Uint32Slice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
