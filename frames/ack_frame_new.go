package frames

import (
	"bytes"
	"time"

	"github.com/lucas-clemente/quic-go/protocol"
	"github.com/lucas-clemente/quic-go/utils"
)

// An AckFrameNew is a ACK frame in QUIC c34
type AckFrameNew struct {
	// TODO: rename to LargestAcked
	LargestObserved protocol.PacketNumber
	NackRanges      []NackRange // has to be ordered. The NACK range with the highest FirstPacketNumber goes first, the NACK range with the lowest FirstPacketNumber goes last

	DelayTime          time.Duration
	PacketReceivedTime time.Time // only for received packets. Will not be modified for received ACKs frames
}

// ParseAckFrameNew reads an ACK frame
func ParseAckFrameNew(r *bytes.Reader, version protocol.VersionNumber) (*AckFrameNew, error) {
	frame := &AckFrameNew{}

	typeByte, err := r.ReadByte()
	if err != nil {
		return nil, err
	}

	hasNACK := false
	if typeByte&0x20 == 0x20 {
		hasNACK = true
	}

	if hasNACK {
		panic("NACKs not yet implemented")
	}

	largestObservedLen := 2 * ((typeByte & 0x0C) >> 2)
	if largestObservedLen == 0 {
		largestObservedLen = 1
	}

	missingSequenceNumberDeltaLen := 2 * (typeByte & 0x03)
	if missingSequenceNumberDeltaLen == 0 {
		missingSequenceNumberDeltaLen = 1
	}

	largestObserved, err := utils.ReadUintN(r, largestObservedLen)
	if err != nil {
		return nil, err
	}
	frame.LargestObserved = protocol.PacketNumber(largestObserved)

	delay, err := utils.ReadUfloat16(r)
	if err != nil {
		return nil, err
	}
	frame.DelayTime = time.Duration(delay) * time.Microsecond

	// TODO: read number of ACK blocks if n flag is set

	ackBlockLength, err := utils.ReadUintN(r, missingSequenceNumberDeltaLen)
	if err != nil {
		return nil, err
	}
	utils.Debugf("ackBlockLength: %d", ackBlockLength)

	// TODO: read ACK blocks

	var numTimestampByte byte
	numTimestampByte, err = r.ReadByte()
	if err != nil {
		return nil, err
	}
	numTimestamp := uint8(numTimestampByte)

	// Delta Largest observed
	_, err = r.ReadByte()
	if err != nil {
		return nil, err
	}
	// First Timestamp
	_, err = utils.ReadUint32(r)
	if err != nil {
		return nil, err
	}

	for i := 0; i < int(numTimestamp)-1; i++ {
		// Delta Largest observed
		_, err = r.ReadByte()
		if err != nil {
			return nil, err
		}

		// Time Since Previous Timestamp
		_, err = utils.ReadUint16(r)
		if err != nil {
			return nil, err
		}
	}

	return frame, nil
}

// Write writes an ACK frame.
func (f *AckFrameNew) Write(b *bytes.Buffer, version protocol.VersionNumber) error {
	largestObservedLen := protocol.GetPacketNumberLength(f.LargestObserved)

	typeByte := uint8(0x40)

	if largestObservedLen != protocol.PacketNumberLen1 {
		typeByte ^= (uint8(largestObservedLen / 2)) << 2
	}

	missingSequenceNumberDeltaLen := largestObservedLen
	if missingSequenceNumberDeltaLen != protocol.PacketNumberLen1 {
		typeByte ^= (uint8(missingSequenceNumberDeltaLen / 2))
	}

	f.DelayTime = time.Now().Sub(f.PacketReceivedTime)

	b.WriteByte(typeByte)

	switch largestObservedLen {
	case protocol.PacketNumberLen1:
		b.WriteByte(uint8(f.LargestObserved))
	case protocol.PacketNumberLen2:
		utils.WriteUint16(b, uint16(f.LargestObserved))
	case protocol.PacketNumberLen4:
		utils.WriteUint32(b, uint32(f.LargestObserved))
	case protocol.PacketNumberLen6:
		utils.WriteUint48(b, uint64(f.LargestObserved))
	}

	utils.WriteUfloat16(b, uint64(f.DelayTime/time.Microsecond))

	// TODO: write number of ACK blocks, if present

	switch missingSequenceNumberDeltaLen {
	case protocol.PacketNumberLen1:
		b.WriteByte(uint8(f.LargestObserved))
	case protocol.PacketNumberLen2:
		utils.WriteUint16(b, uint16(f.LargestObserved))
	case protocol.PacketNumberLen4:
		utils.WriteUint32(b, uint32(f.LargestObserved))
	case protocol.PacketNumberLen6:
		utils.WriteUint48(b, uint64(f.LargestObserved))
	}

	// TODO: write ACK blocks

	b.WriteByte(0x01)       // Just one timestamp
	b.WriteByte(0x00)       // Delta Largest observed
	utils.WriteUint32(b, 0) // First timestamp

	return nil
}

// MinLength of a written frame
func (f *AckFrameNew) MinLength(version protocol.VersionNumber) (protocol.ByteCount, error) {
	var length protocol.ByteCount
	length = 1 + 2 + 1 + 1 + 4 // 1 TypeByte, 2 ACK delay time, 1 Num Timestamp, 1 Delta Largest Observed, 4 FirstTimestamp
	length += protocol.ByteCount(protocol.GetPacketNumberLength(f.LargestObserved))
	// for the first ACK block length
	length += protocol.ByteCount(protocol.GetPacketNumberLength(f.LargestObserved))

	length += (1 + 2) * 0 /* TODO: num_timestamps */
	if f.HasNACK() {
		panic("NACKs not yet implemented")
	}
	return length, nil
}

// HasNACK returns if the frame has NACK ranges
func (f *AckFrameNew) HasNACK() bool {
	if len(f.NackRanges) > 0 {
		return true
	}
	return false
}
